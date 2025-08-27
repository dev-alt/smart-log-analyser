package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"smart-log-analyser/pkg/analyser"
	"smart-log-analyser/pkg/config"
	"smart-log-analyser/pkg/html"
	"smart-log-analyser/pkg/parser"
	"smart-log-analyser/pkg/query"
)

// Server represents the IPC server
type Server struct {
	listener   net.Listener
	clients    map[string]net.Conn
	clientsMux sync.RWMutex
	shutdown   chan struct{}
	analyzer   *analyser.Analyser
	configMgr  *config.ConfigManager
	htmlGen    *html.Generator
}

// NewServer creates a new IPC server
func NewServer() (*Server, error) {
	analyzer := analyser.New()
	configMgr := config.NewConfigManager("config")
	htmlGen, err := html.NewGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to create HTML generator: %v", err)
	}

	return &Server{
		clients:   make(map[string]net.Conn),
		shutdown:  make(chan struct{}),
		analyzer:  analyzer,
		configMgr: configMgr,
		htmlGen:   htmlGen,
	}, nil
}

// Start starts the IPC server
func (s *Server) Start() error {
	var err error
	
	if runtime.GOOS == "windows" {
		// Use Named Pipes on Windows
		pipeName := `\\.\pipe\SmartLogAnalyser`
		s.listener, err = net.Listen("pipe", pipeName)
		if err != nil {
			return fmt.Errorf("failed to create named pipe: %v", err)
		}
		log.Printf("IPC server listening on named pipe: %s", pipeName)
	} else {
		// Use Unix Domain Sockets on other platforms
		socketPath := filepath.Join(os.TempDir(), "smart-log-analyser.sock")
		
		// Remove existing socket file
		os.Remove(socketPath)
		
		s.listener, err = net.Listen("unix", socketPath)
		if err != nil {
			return fmt.Errorf("failed to create unix socket: %v", err)
		}
		log.Printf("IPC server listening on unix socket: %s", socketPath)
	}

	go s.acceptConnections()
	return nil
}

// Stop stops the IPC server
func (s *Server) Stop() error {
	close(s.shutdown)
	
	s.clientsMux.Lock()
	for _, conn := range s.clients {
		conn.Close()
	}
	s.clientsMux.Unlock()
	
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// acceptConnections accepts new client connections
func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					log.Printf("Accept error: %v", err)
				}
				continue
			}
			
			clientID := fmt.Sprintf("client_%d", len(s.clients))
			s.clientsMux.Lock()
			s.clients[clientID] = conn
			s.clientsMux.Unlock()
			
			log.Printf("Client connected: %s", clientID)
			go s.handleConnection(clientID, conn)
		}
	}
}

// handleConnection handles a client connection
func (s *Server) handleConnection(clientID string, conn net.Conn) {
	defer func() {
		conn.Close()
		s.clientsMux.Lock()
		delete(s.clients, clientID)
		s.clientsMux.Unlock()
		log.Printf("Client disconnected: %s", clientID)
	}()

	scanner := bufio.NewScanner(conn)
	encoder := json.NewEncoder(conn)

	for scanner.Scan() {
		var request IPCRequest
		if err := json.Unmarshal(scanner.Bytes(), &request); err != nil {
			response := IPCResponse{
				ID:      "unknown",
				Success: false,
				Error:   fmt.Sprintf("Invalid JSON: %v", err),
			}
			encoder.Encode(response)
			continue
		}

		response := s.processRequest(request)
		if err := encoder.Encode(response); err != nil {
			log.Printf("Failed to send response to %s: %v", clientID, err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error with %s: %v", clientID, err)
	}
}

// processRequest processes an IPC request
func (s *Server) processRequest(request IPCRequest) IPCResponse {
	log.Printf("Processing request: %s (action: %s)", request.ID, request.Action)
	
	switch request.Action {
	case ActionAnalyze:
		return s.handleAnalyze(request)
	case ActionQuery:
		return s.handleQuery(request)
	case ActionListPresets:
		return s.handleListPresets(request)
	case ActionRunPreset:
		return s.handleRunPreset(request)
	case ActionGetConfig:
		return s.handleGetConfig(request)
	case ActionUpdateConfig:
		return s.handleUpdateConfig(request)
	case ActionGetStatus:
		return s.handleGetStatus(request)
	case ActionShutdown:
		go func() {
			s.Stop()
		}()
		return IPCResponse{
			ID:      request.ID,
			Success: true,
			Data:    AnalysisResultData{Status: "shutting down"},
		}
	default:
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   ErrInvalidAction,
		}
	}
}

// handleAnalyze processes analyze requests
func (s *Server) handleAnalyze(request IPCRequest) IPCResponse {
	if request.LogFile == "" {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   ErrMissingLogFile,
		}
	}

	// Check if file exists
	if _, err := os.Stat(request.LogFile); os.IsNotExist(err) {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("Log file not found: %s", request.LogFile),
		}
	}

	// Parse the log file
	p := parser.New()
	logs, err := p.ParseFile(request.LogFile)
	if err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrAnalysisFailed, err),
		}
	}

	// Perform analysis
	results := s.analyzer.Analyse(logs, nil, nil)

	responseData := AnalysisResultData{
		Results: results,
	}

	// Generate HTML if requested
	if request.Options.GenerateHTML {
		htmlTitle := request.Options.HTMLTitle
		if htmlTitle == "" {
			htmlTitle = "Smart Log Analysis Report"
		}

		outputPath := request.Options.OutputPath
		if outputPath == "" {
			outputPath = "report.html"
		}

		var htmlErr error
		if request.Options.Interactive {
			htmlErr = s.htmlGen.GenerateInteractiveReport(results, outputPath, htmlTitle)
		} else {
			htmlErr = s.htmlGen.GenerateReport(results, outputPath, htmlTitle)
		}

		if htmlErr != nil {
			return IPCResponse{
				ID:      request.ID,
				Success: false,
				Error:   fmt.Sprintf("HTML generation failed: %v", htmlErr),
			}
		}

		responseData.HTMLPath = outputPath
	}

	return IPCResponse{
		ID:      request.ID,
		Success: true,
		Data:    responseData,
	}
}

// handleQuery processes query requests
func (s *Server) handleQuery(request IPCRequest) IPCResponse {
	if request.LogFile == "" || request.Query == "" {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   "Missing log file or query",
		}
	}

	// Parse the log file
	p := parser.New()
	logs, err := p.ParseFile(request.LogFile)
	if err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrQueryFailed, err),
		}
	}

	// Execute query
	results, err := query.ExecuteQuery(request.Query, logs)
	if err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrQueryFailed, err),
		}
	}

	return IPCResponse{
		ID:      request.ID,
		Success: true,
		Data: AnalysisResultData{
			QueryResults: results,
		},
	}
}

// handleListPresets processes list presets requests
func (s *Server) handleListPresets(request IPCRequest) IPCResponse {
	err := s.configMgr.Load()
	if err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrConfigFailed, err),
		}
	}

	appConfig := s.configMgr.GetConfig()
	return IPCResponse{
		ID:      request.ID,
		Success: true,
		Data: AnalysisResultData{
			Presets: appConfig.Presets,
		},
	}
}

// handleRunPreset processes run preset requests
func (s *Server) handleRunPreset(request IPCRequest) IPCResponse {
	if request.LogFile == "" || request.Preset == "" {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   "Missing log file or preset name",
		}
	}

	err := s.configMgr.Load()
	if err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrConfigFailed, err),
		}
	}

	appConfig := s.configMgr.GetConfig()

	// Find preset
	var selectedPreset config.AnalysisPreset
	var found bool
	for _, preset := range appConfig.Presets {
		if preset.Name == request.Preset {
			selectedPreset = preset
			found = true
			break
		}
	}

	if !found {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("Preset not found: %s", request.Preset),
		}
	}

	// Parse the log file
	p := parser.New()
	logs, err := p.ParseFile(request.LogFile)
	if err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrPresetFailed, err),
		}
	}

	// Execute preset query
	results, err := query.ExecuteQuery(selectedPreset.Query, logs)
	if err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrPresetFailed, err),
		}
	}

	return IPCResponse{
		ID:      request.ID,
		Success: true,
		Data: AnalysisResultData{
			QueryResults: results,
		},
	}
}

// handleGetConfig processes get config requests
func (s *Server) handleGetConfig(request IPCRequest) IPCResponse {
	err := s.configMgr.Load()
	if err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrConfigFailed, err),
		}
	}

	appConfig := s.configMgr.GetConfig()
	return IPCResponse{
		ID:      request.ID,
		Success: true,
		Data: AnalysisResultData{
			Config: appConfig,
		},
	}
}

// handleUpdateConfig processes update config requests
func (s *Server) handleUpdateConfig(request IPCRequest) IPCResponse {
	if request.Config == nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   "Missing configuration data",
		}
	}

	// Convert map to config structure
	configJSON, err := json.Marshal(request.Config)
	if err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("Invalid configuration format: %v", err),
		}
	}

	var newConfig config.AppConfig
	if err := json.Unmarshal(configJSON, &newConfig); err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("Invalid configuration structure: %v", err),
		}
	}

	if err := s.configMgr.SetConfig(&newConfig); err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrConfigFailed, err),
		}
	}

	if err := s.configMgr.Save(); err != nil {
		return IPCResponse{
			ID:      request.ID,
			Success: false,
			Error:   fmt.Sprintf("%s: %v", ErrConfigFailed, err),
		}
	}

	return IPCResponse{
		ID:      request.ID,
		Success: true,
		Data: AnalysisResultData{
			Status: "configuration updated",
		},
	}
}

// handleGetStatus processes get status requests
func (s *Server) handleGetStatus(request IPCRequest) IPCResponse {
	status := fmt.Sprintf("Smart Log Analyser IPC Server - Running on %s", runtime.GOOS)
	if runtime.GOOS == "windows" {
		status += " (Named Pipes)"
	} else {
		status += " (Unix Domain Sockets)"
	}

	clientCount := len(s.clients)
	status += fmt.Sprintf(" - %d active clients", clientCount)

	return IPCResponse{
		ID:      request.ID,
		Success: true,
		Data: AnalysisResultData{
			Status: status,
		},
	}
}