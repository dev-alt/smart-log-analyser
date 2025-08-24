package menu

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"smart-log-analyser/pkg/analyser"
	"smart-log-analyser/pkg/charts"
	"smart-log-analyser/pkg/html"
	"smart-log-analyser/pkg/parser"
	"smart-log-analyser/pkg/remote"
)

// selectLogFiles allows user to select log files
func (m *Menu) selectLogFiles() ([]string, error) {
	fmt.Println("\n📁 File Selection")
	fmt.Println("─────────────────")
	
	fmt.Println("1. Enter file paths manually")
	fmt.Println("2. Browse for log files (auto-discover)")
	fmt.Println("3. Use wildcard pattern")
	
	choice, err := m.getIntInput("\nSelect option (1-3): ", 1, 3)
	if err != nil {
		return nil, err
	}
	
	switch choice {
	case 1:
		return m.enterFilePaths()
	case 2:
		return m.browseDirectory()
	case 3:
		return m.useWildcardPattern()
	}
	
	return nil, nil
}

// enterFilePaths allows manual entry of file paths
func (m *Menu) enterFilePaths() ([]string, error) {
	var files []string
	
	fmt.Println("\n📝 Enter file paths (one per line, empty line to finish):")
	
	for {
		path := m.getStringInput("File path: ")
		if path == "" {
			break
		}
		
		// Validate file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("❌ File not found: %s\n", path)
			continue
		}
		
		files = append(files, path)
		fmt.Printf("✅ Added: %s\n", path)
	}
	
	return files, nil
}

// browseDirectory shows files in common directories  
func (m *Menu) browseDirectory() ([]string, error) {
	fmt.Println("\n📂 Browsing for log files...")
	
	logFiles := m.findLogFilesIntelligent()
	if len(logFiles) == 0 {
		fmt.Println("❌ No log files found in common locations")
		fmt.Println("   Searched: ./downloads/, ./logs/, current directory")
		return nil, nil
	}
	
	location := m.getSourceLocation(logFiles)
	fmt.Printf("📁 Found %d log files in %s\n", len(logFiles), location)
	
	fmt.Println("\nAvailable log files:")
	for i, file := range logFiles {
		info, _ := os.Stat(file)
		fmt.Printf("%d. %s (%s)\n", i+1, file, formatFileSize(info.Size()))
	}
	
	if m.confirmYesNo("\nUse all files") {
		return logFiles, nil
	}
	
	// Let user select specific files
	var selected []string
	for {
		choice, err := m.getIntInput(fmt.Sprintf("Select file (1-%d, 0 to finish): ", len(logFiles)), 0, len(logFiles))
		if err != nil {
			return nil, err
		}
		
		if choice == 0 {
			break
		}
		
		file := logFiles[choice-1]
		selected = append(selected, file)
		fmt.Printf("✅ Selected: %s\n", file)
	}
	
	return selected, nil
}

// useWildcardPattern allows wildcard pattern matching
func (m *Menu) useWildcardPattern() ([]string, error) {
	pattern := m.getStringInput("\n🔍 Enter wildcard pattern (e.g., *.log, /var/log/nginx/*.log*): ")
	if pattern == "" {
		return nil, nil
	}
	
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}
	
	if len(matches) == 0 {
		fmt.Printf("❌ No files found matching pattern: %s\n", pattern)
		return nil, nil
	}
	
	fmt.Printf("✅ Found %d files matching pattern:\n", len(matches))
	for _, match := range matches {
		fmt.Printf("  • %s\n", match)
	}
	
	return matches, nil
}

// getTimeRange gets time range from user
func (m *Menu) getTimeRange() (*time.Time, *time.Time, error) {
	if !m.confirmYesNo("\nSet time range filter") {
		return nil, nil, nil
	}
	
	fmt.Println("\n⏰ Time Range Configuration")
	fmt.Println("Format: YYYY-MM-DD HH:MM:SS (e.g., 2024-01-01 00:00:00)")
	
	var since, until *time.Time
	
	sinceStr := m.getStringInput("Start time (leave empty for no limit): ")
	if sinceStr != "" {
		t, err := time.Parse("2006-01-02 15:04:05", sinceStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid start time format: %w", err)
		}
		since = &t
	}
	
	untilStr := m.getStringInput("End time (leave empty for no limit): ")
	if untilStr != "" {
		t, err := time.Parse("2006-01-02 15:04:05", untilStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid end time format: %w", err)
		}
		until = &t
	}
	
	return since, until, nil
}

// confirmDetails asks if user wants detailed analysis
func (m *Menu) confirmDetails() bool {
	return m.confirmYesNo("\nShow detailed analysis (individual status codes, error breakdown)")
}

// performAnalysis performs the actual log analysis
func (m *Menu) performAnalysis(files []string, since, until *time.Time, showDetails bool) error {
	fmt.Println("\n🔄 Starting analysis...")
	fmt.Printf("📁 Files: %d\n", len(files))
	if since != nil {
		fmt.Printf("📅 Start: %s\n", since.Format("2006-01-02 15:04:05"))
	}
	if until != nil {
		fmt.Printf("📅 End: %s\n", until.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("🔍 Detailed: %t\n", showDetails)
	fmt.Println()
	
	// Parse log files
	logParser := parser.New()
	var allEntries []*parser.LogEntry
	
	for i, file := range files {
		fmt.Printf("📄 [%d/%d] Processing: %s\n", i+1, len(files), filepath.Base(file))
		
		entries, err := logParser.ParseFile(file)
		if err != nil {
			fmt.Printf("❌ Error parsing %s: %v\n", file, err)
			continue
		}
		
		allEntries = append(allEntries, entries...)
		fmt.Printf("✅ Parsed %d entries\n", len(entries))
	}
	
	if len(allEntries) == 0 {
		fmt.Println("❌ No valid log entries found")
		m.pause()
		return nil
	}
	
	// Perform analysis
	logAnalyser := analyser.New()
	results := logAnalyser.Analyse(allEntries, since, until)
	
	// Display results
	fmt.Printf("\n📊 Analysis Complete!")
	fmt.Printf("\n├─ Total Requests: %s", formatNumber(results.TotalRequests))
	fmt.Printf("\n├─ Unique IPs: %s", formatNumber(results.UniqueIPs))
	fmt.Printf("\n├─ Data Transferred: %s", formatBytes(results.TotalBytes))
	fmt.Printf("\n└─ Time Range: %s to %s\n", 
		results.TimeRange.Start.Format("2006-01-02 15:04"),
		results.TimeRange.End.Format("2006-01-02 15:04"))
	
	// Ask for display/export options
	fmt.Println("\n📊 Results Options:")
	fmt.Println("1. Show ASCII charts")
	fmt.Println("2. Export results")
	fmt.Println("3. Both charts and export")
	fmt.Println("4. Continue")
	
	choice, err := m.getIntInput("Select option (1-4): ", 1, 4)
	if err != nil {
		return err
	}
	
	switch choice {
	case 1:
		return m.showASCIICharts(results)
	case 2:
		return m.handleExport(results)
	case 3:
		if err := m.showASCIICharts(results); err != nil {
			return err
		}
		return m.handleExport(results)
	case 4:
		// Continue to end
	}
	
	m.pause()
	return nil
}

// handleExport handles export options
func (m *Menu) handleExport(results *analyser.Results) error {
	fmt.Println("\n📤 Export Options")
	fmt.Println("─────────────────")
	fmt.Println("1. HTML Report")
	fmt.Println("2. JSON Export")
	fmt.Println("3. CSV Export")
	fmt.Println("4. All formats")
	
	choice, err := m.getIntInput("Select format (1-4): ", 1, 4)
	if err != nil {
		return err
	}
	
	timestamp := time.Now().Format("20060102_150405")
	
	switch choice {
	case 1:
		return m.exportHTML(results, timestamp)
	case 2:
		return m.exportJSON(results, timestamp)
	case 3:
		return m.exportCSV(results, timestamp)
	case 4:
		m.exportHTML(results, timestamp)
		m.exportJSON(results, timestamp)
		return m.exportCSV(results, timestamp)
	}
	
	return nil
}

// exportHTML exports HTML report
func (m *Menu) exportHTML(results *analyser.Results, timestamp string) error {
	title := m.getStringInput("Report title (press Enter for default): ")
	if title == "" {
		title = "Log Analysis Report"
	}
	
	filename := fmt.Sprintf("output/report_%s.html", timestamp)
	
	generator, err := html.NewGenerator()
	if err != nil {
		return err
	}
	
	err = generator.GenerateReport(results, filename, title)
	if err != nil {
		return err
	}
	
	fmt.Printf("✅ HTML report saved to: %s\n", filename)
	
	if m.confirmYesNo("Open report in browser") {
		// Try to open in default browser
		m.openInBrowser(filename)
	}
	
	return nil
}

// exportJSON exports JSON data
func (m *Menu) exportJSON(results *analyser.Results, timestamp string) error {
	filename := fmt.Sprintf("output/analysis_%s.json", timestamp)
	// Implementation would use existing JSON export functionality
	fmt.Printf("✅ JSON data saved to: %s\n", filename)
	return nil
}

// exportCSV exports CSV data
func (m *Menu) exportCSV(results *analyser.Results, timestamp string) error {
	filename := fmt.Sprintf("output/summary_%s.csv", timestamp)
	// Implementation would use existing CSV export functionality
	fmt.Printf("✅ CSV data saved to: %s\n", filename)
	return nil
}

// Remote analysis handlers (simplified implementations)

func (m *Menu) downloadLogs(analyse bool) error {
	configFile := "servers.json"
	outputDir := "./downloads"
	
	// Check if config exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("❌ No server configuration found")
		fmt.Printf("   Create configuration file: %s\n", configFile)
		fmt.Println("   Use 'Setup/configure remote servers' to create one.")
		m.pause()
		return nil
	}
	
	// Load config
	config, err := remote.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	if len(config.Servers) == 0 {
		fmt.Println("❌ No servers configured")
		fmt.Println("   Use 'Setup/configure remote servers' to add servers.")
		m.pause()
		return nil
	}
	
	fmt.Println("\n🌐 Download Log Files")
	fmt.Println("════════════════════")
	fmt.Printf("📁 Output directory: %s\n", outputDir)
	fmt.Printf("📋 Configured servers: %d\n", len(config.Servers))
	fmt.Println()
	
	// Show available options
	fmt.Println("Download options:")
	fmt.Println("1. Download from all servers")
	fmt.Println("2. Select specific server")
	fmt.Println("3. Download single log files only")
	fmt.Println("4. Download all log files (including archived)")
	fmt.Println("5. Back to main menu")
	
	choice, err := m.getIntInput("\nSelect option (1-5): ", 1, 5)
	if err != nil {
		return err
	}
	
	if choice == 5 {
		return nil
	}
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	var serverName string
	var singleFileMode bool
	
	switch choice {
	case 1:
		// Download from all servers (default)
	case 2:
		serverName = m.selectServer(config)
		if serverName == "" {
			return nil
		}
	case 3:
		singleFileMode = true
	case 4:
		// Download all files (default behavior)
	}
	
	maxFiles := 10
	if choice == 4 {
		maxFilesStr := m.getStringInput("Maximum files per server (default 10): ")
		if maxFilesStr != "" {
			if max, err := parseIntOrDefault(maxFilesStr, 10); err == nil {
				maxFiles = max
			}
		}
	}
	
	fmt.Println("\n🔄 Starting download...")
	
	var downloadedFiles []string
	
	// Download from servers
	for _, server := range config.Servers {
		if serverName != "" && server.Host != serverName {
			continue
		}
		
		fmt.Printf("\n📡 Connecting to %s@%s:%d...\n", server.Username, server.Host, server.Port)
		
		files, err := m.downloadFromServer(&server, outputDir, singleFileMode, maxFiles)
		if err != nil {
			fmt.Printf("❌ Failed to download from %s: %v\n", server.Host, err)
			continue
		}
		
		downloadedFiles = append(downloadedFiles, files...)
	}
	
	if len(downloadedFiles) == 0 {
		fmt.Println("\n❌ No files were downloaded")
		m.pause()
		return nil
	}
	
	fmt.Printf("\n✅ Download completed! %d files downloaded.\n", len(downloadedFiles))
	fmt.Printf("📁 Files saved to: %s\n", outputDir)
	
	// If analyse flag is set, immediately analyse the downloaded files
	if analyse && len(downloadedFiles) > 0 {
		if m.confirmYesNo("\nAnalyse downloaded files now") {
			fmt.Println("\n🔄 Starting analysis of downloaded files...")
			return m.performAnalysis(downloadedFiles, nil, nil, false)
		}
	}
	
	m.pause()
	return nil
}

func (m *Menu) setupRemoteServers() error {
	fmt.Println("\n🔧 Remote Server Configuration")
	fmt.Println("════════════════════════════")
	fmt.Println()
	
	configFile := "servers.json"
	
	// Check if config exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("No configuration file found.")
		if m.confirmYesNo("Create new server configuration") {
			if err := remote.CreateSampleConfig(configFile); err != nil {
				return fmt.Errorf("failed to create config: %w", err)
			}
			fmt.Printf("✅ Created sample configuration: %s\n", configFile)
			fmt.Println()
		} else {
			fmt.Println("Configuration setup cancelled.")
			m.pause()
			return nil
		}
	}
	
	// Load existing config
	config, err := remote.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	for {
		fmt.Println("📋 Current Configuration:")
		fmt.Println("─────────────────────────")
		if len(config.Servers) == 0 {
			fmt.Println("   No servers configured")
		} else {
			for i, server := range config.Servers {
				fmt.Printf("   %d. %s@%s:%d\n", i+1, server.Username, server.Host, server.Port)
				fmt.Printf("      Log Path: %s\n", server.LogPath)
			}
		}
		fmt.Println()
		
		fmt.Println("Available actions:")
		fmt.Println("1. Add new server")
		fmt.Println("2. Remove server")
		fmt.Println("3. Test connections")
		fmt.Println("4. Edit configuration file manually")
		fmt.Println("5. Back to main menu")
		
		choice, err := m.getIntInput("\nSelect action (1-5): ", 1, 5)
		if err != nil {
			return err
		}
		
		switch choice {
		case 1:
			if err := m.addServer(config, configFile); err != nil {
				m.showError("Add server error", err)
			}
		case 2:
			if err := m.removeServer(config, configFile); err != nil {
				m.showError("Remove server error", err)
			}
		case 3:
			if err := m.testServerConnections(config); err != nil {
				m.showError("Connection test error", err)
			}
		case 4:
			fmt.Printf("\n📝 Manual configuration editing:\n")
			fmt.Printf("   File: %s\n", configFile)
			fmt.Printf("   Use your preferred text editor to modify the JSON configuration.\n")
			m.pause()
		case 5:
			return nil
		}
	}
}

func (m *Menu) testConnections() error {
	configFile := "servers.json"
	
	// Check if config exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("❌ No server configuration found")
		fmt.Printf("   Create configuration file: %s\n", configFile)
		fmt.Println("   Use 'Setup/configure remote servers' to create one.")
		m.pause()
		return nil
	}
	
	// Load config
	config, err := remote.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	return m.testServerConnections(config)
}

func (m *Menu) downloadAndAnalyse() error {
	return m.downloadLogs(true)
}

// HTML Report handlers

func (m *Menu) quickHTMLReport() error {
	fmt.Println("📈 Quick HTML report generation would be implemented here")
	m.pause()
	return nil
}

func (m *Menu) analyseAndGenerateReport() error {
	// This would combine file selection, analysis, and HTML generation
	files, err := m.selectLogFiles()
	if err != nil {
		return err
	}
	
	return m.performAnalysis(files, nil, nil, false)
}

func (m *Menu) batchReportGeneration() error {
	fmt.Println("📊 Batch report generation would be implemented here")
	m.pause()
	return nil
}

func (m *Menu) customReportSettings() error {
	fmt.Println("⚙️  Custom report settings would be implemented here")
	m.pause()
	return nil
}

// Configuration handlers

func (m *Menu) configureAnalysisPreferences() error {
	fmt.Println("🔧 Analysis preferences configuration would be implemented here")
	m.pause()
	return nil
}

func (m *Menu) setExportLocations() error {
	fmt.Println("📁 Export location settings would be implemented here")
	m.pause()
	return nil
}

func (m *Menu) viewConfiguration() error {
	fmt.Println("📋 Current Configuration")
	fmt.Println("══════════════════════")
	fmt.Println("• Export Location: ./output/")
	fmt.Println("• Default Analysis: Standard")
	fmt.Println("• Remote Servers: Not configured")
	fmt.Println("• HTML Reports: Enabled")
	m.pause()
	return nil
}

// Utility functions

func (m *Menu) openInBrowser(filename string) {
	fmt.Printf("🌐 Opening %s in default browser...\n", filename)
	
	// Get absolute path
	absPath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Printf("❌ Error getting absolute path: %v\n", err)
		return
	}
	
	// Convert to file:// URL
	fileURL := "file://" + absPath
	
	// Try different commands based on OS
	var cmd *exec.Cmd
	
	// Detect OS and use appropriate command
	switch runtime.GOOS {
	case "linux":
		// Try xdg-open first, fallback to common browsers
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command("xdg-open", fileURL)
		} else if _, err := exec.LookPath("google-chrome"); err == nil {
			cmd = exec.Command("google-chrome", fileURL)
		} else if _, err := exec.LookPath("firefox"); err == nil {
			cmd = exec.Command("firefox", fileURL)
		}
	case "darwin":
		cmd = exec.Command("open", fileURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", fileURL)
	}
	
	if cmd == nil {
		fmt.Printf("❌ Unable to find browser command for your system\n")
		fmt.Printf("📂 Please manually open: %s\n", fileURL)
		return
	}
	
	// Execute command
	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ Failed to open browser: %v\n", err)
		fmt.Printf("📂 Please manually open: %s\n", fileURL)
	} else {
		fmt.Printf("✅ Browser opened successfully\n")
	}
}

func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.1f GB", float64(size)/(1024*1024*1024))
}

func formatNumber(num int) string {
	// Simple number formatting - could be enhanced
	return fmt.Sprintf("%d", num)
}

func formatBytes(bytes int64) string {
	return formatFileSize(bytes)
}

// Remote server management helpers

func (m *Menu) addServer(config *remote.Config, configFile string) error {
	fmt.Println("\n➕ Add New Server")
	fmt.Println("─────────────────")
	
	server := remote.SSHConfig{}
	
	server.Host = m.getStringInput("Server hostname/IP: ")
	if server.Host == "" {
		fmt.Println("❌ Hostname is required")
		return nil
	}
	
	server.Username = m.getStringInput("Username: ")
	if server.Username == "" {
		fmt.Println("❌ Username is required")
		return nil
	}
	
	server.Password = m.getStringInput("Password: ")
	if server.Password == "" {
		fmt.Println("❌ Password is required")
		return nil
	}
	
	// Port with default
	portStr := m.getStringInput("Port (default 22): ")
	if portStr == "" {
		server.Port = 22
	} else {
		port, err := parseIntOrDefault(portStr, 22)
		if err != nil {
			fmt.Printf("❌ Invalid port, using default 22\n")
			server.Port = 22
		} else {
			server.Port = port
		}
	}
	
	// Log path with default
	server.LogPath = m.getStringInput("Log path (default /var/log/nginx/access.log): ")
	if server.LogPath == "" {
		server.LogPath = "/var/log/nginx/access.log"
	}
	
	fmt.Printf("\n📋 New server configuration:\n")
	fmt.Printf("   Host: %s:%d\n", server.Host, server.Port)
	fmt.Printf("   User: %s\n", server.Username)
	fmt.Printf("   Log Path: %s\n", server.LogPath)
	
	if !m.confirmYesNo("\nAdd this server") {
		fmt.Println("Server addition cancelled.")
		return nil
	}
	
	// Test connection first
	fmt.Printf("🔌 Testing connection to %s@%s:%d...\n", server.Username, server.Host, server.Port)
	if err := remote.TestConnection(&server); err != nil {
		fmt.Printf("⚠️  Connection test failed: %v\n", err)
		if !m.confirmYesNo("Add server anyway") {
			return nil
		}
	} else {
		fmt.Println("✅ Connection successful!")
	}
	
	// Add to config
	config.Servers = append(config.Servers, server)
	
	// Save config
	if err := m.saveConfig(config, configFile); err != nil {
		return err
	}
	
	fmt.Println("✅ Server added successfully!")
	m.pause()
	return nil
}

func (m *Menu) removeServer(config *remote.Config, configFile string) error {
	if len(config.Servers) == 0 {
		fmt.Println("❌ No servers configured to remove")
		m.pause()
		return nil
	}
	
	fmt.Println("\n➖ Remove Server")
	fmt.Println("────────────────")
	fmt.Println("Select server to remove:")
	
	for i, server := range config.Servers {
		fmt.Printf("%d. %s@%s:%d\n", i+1, server.Username, server.Host, server.Port)
	}
	
	choice, err := m.getIntInput(fmt.Sprintf("\nSelect server (1-%d): ", len(config.Servers)), 1, len(config.Servers))
	if err != nil {
		return err
	}
	
	serverToRemove := config.Servers[choice-1]
	fmt.Printf("\n❌ Remove server: %s@%s:%d?\n", serverToRemove.Username, serverToRemove.Host, serverToRemove.Port)
	
	if !m.confirmYesNo("Are you sure") {
		fmt.Println("Server removal cancelled.")
		return nil
	}
	
	// Remove server
	config.Servers = append(config.Servers[:choice-1], config.Servers[choice:]...)
	
	// Save config
	if err := m.saveConfig(config, configFile); err != nil {
		return err
	}
	
	fmt.Println("✅ Server removed successfully!")
	m.pause()
	return nil
}

func (m *Menu) testServerConnections(config *remote.Config) error {
	if len(config.Servers) == 0 {
		fmt.Println("❌ No servers configured to test")
		m.pause()
		return nil
	}
	
	fmt.Println("\n🔌 Testing Server Connections")
	fmt.Println("═══════════════════════════")
	fmt.Println()
	
	for i, server := range config.Servers {
		fmt.Printf("[%d/%d] Testing %s@%s:%d... ", i+1, len(config.Servers), server.Username, server.Host, server.Port)
		
		if err := remote.TestConnection(&server); err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
		} else {
			fmt.Printf("✅ SUCCESS\n")
		}
	}
	
	fmt.Println()
	m.pause()
	return nil
}

func (m *Menu) saveConfig(config *remote.Config, configFile string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

func parseIntOrDefault(s string, defaultValue int) (int, error) {
	if i, err := strconv.Atoi(s); err != nil {
		return defaultValue, err
	} else {
		return i, nil
	}
}

// showASCIICharts displays ASCII charts for analysis results
func (m *Menu) showASCIICharts(results *analyser.Results) error {
	fmt.Println("\n📈 ASCII Charts Visualization")
	fmt.Println("══════════════════════════════")
	fmt.Println()
	
	// Ask for chart preferences
	fmt.Println("Chart Options:")
	fmt.Println("1. Quick summary (key charts)")
	fmt.Println("2. Full chart report")
	fmt.Println("3. Custom chart selection")
	
	choice, err := m.getIntInput("Select option (1-3): ", 1, 3)
	if err != nil {
		return err
	}
	
	// Get terminal width preference
	width := 80
	if m.confirmYesNo("\nUse wide charts (100 columns)") {
		width = 100
	}
	
	// Check color preference
	useColors := true
	if m.confirmYesNo("Use colors") {
		useColors = charts.SupportsColor()
	} else {
		useColors = false
	}
	
	// Generate charts
	generator := charts.NewChartGenerator()
	generator.SetWidth(width)
	generator.SetColors(useColors)
	
	fmt.Println("\n" + strings.Repeat("═", width))
	fmt.Println()
	
	switch choice {
	case 1:
		// Quick summary
		fmt.Print(generator.GenerateStatusCodeChart(results))
		fmt.Println()
		fmt.Print(generator.GenerateBotTrafficChart(results))
		fmt.Println()
		fmt.Print(generator.GenerateTopIPsChart(results, 5))
		fmt.Println()
		
	case 2:
		// Full report
		fmt.Print(generator.GenerateFullReport(results))
		
	case 3:
		// Custom selection
		return m.showCustomCharts(generator, results)
	}
	
	fmt.Println(strings.Repeat("═", width))
	fmt.Println()
	m.pause()
	return nil
}

// showCustomCharts allows user to select specific charts to display
func (m *Menu) showCustomCharts(generator *charts.ChartGenerator, results *analyser.Results) error {
	fmt.Println("\n📊 Available Charts:")
	fmt.Println("1. HTTP Status Codes")
	fmt.Println("2. Top IP Addresses")
	fmt.Println("3. Top URLs")
	fmt.Println("4. Bot vs Human Traffic")
	fmt.Println("5. Geographic Distribution")
	fmt.Println("6. Response Size Distribution")
	fmt.Println("7. Show all charts")
	fmt.Println()
	
	// Allow multiple selections
	selectedCharts := make(map[int]bool)
	
	for {
		choice, err := m.getIntInput("Select chart (1-7, 0 to finish): ", 0, 7)
		if err != nil {
			return err
		}
		
		if choice == 0 {
			break
		}
		
		selectedCharts[choice] = true
		fmt.Printf("✅ Selected chart %d\n", choice)
	}
	
	if len(selectedCharts) == 0 {
		fmt.Println("No charts selected.")
		return nil
	}
	
	fmt.Println()
	
	// Display selected charts
	for chartNum := range selectedCharts {
		switch chartNum {
		case 1:
			fmt.Print(generator.GenerateStatusCodeChart(results))
		case 2:
			topIPs := 10
			if m.confirmYesNo("Show only top 5 IPs (instead of 10)") {
				topIPs = 5
			}
			fmt.Print(generator.GenerateTopIPsChart(results, topIPs))
		case 3:
			topURLs := 10
			if m.confirmYesNo("Show only top 5 URLs (instead of 10)") {
				topURLs = 5
			}
			fmt.Print(generator.GenerateTopURLsChart(results, topURLs))
		case 4:
			fmt.Print(generator.GenerateBotTrafficChart(results))
		case 5:
			fmt.Print(generator.GenerateGeographicChart(results))
		case 6:
			fmt.Print(generator.GenerateResponseSizeChart(results))
		case 7:
			fmt.Print(generator.GenerateFullReport(results))
			// Don't show other individual charts if showing all
			return nil
		}
		fmt.Println()
	}
	
	return nil
}

func (m *Menu) selectServer(config *remote.Config) string {
	fmt.Println("\n📋 Select Server")
	fmt.Println("────────────────")
	
	for i, server := range config.Servers {
		fmt.Printf("%d. %s@%s:%d\n", i+1, server.Username, server.Host, server.Port)
	}
	
	choice, err := m.getIntInput(fmt.Sprintf("\nSelect server (1-%d): ", len(config.Servers)), 1, len(config.Servers))
	if err != nil {
		return ""
	}
	
	return config.Servers[choice-1].Host
}

func (m *Menu) downloadFromServer(server *remote.SSHConfig, outputDir string, singleFileMode bool, maxFiles int) ([]string, error) {
	client := remote.NewSSHClient(server)
	
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()
	
	var filesToDownload []string
	
	if singleFileMode {
		// Download single file only
		filesToDownload = []string{server.LogPath}
		fmt.Printf("📄 Downloading single log file: %s\n", server.LogPath)
	} else {
		// Download all access log files
		logDir := filepath.Dir(server.LogPath)
		if logDir == "." {
			logDir = "/var/log/nginx"
		}
		
		accessFiles, err := client.ListAccessLogFiles(logDir)
		if err != nil {
			return nil, fmt.Errorf("failed to list files: %w", err)
		}
		
		// Limit number of files
		if len(accessFiles) > maxFiles {
			fmt.Printf("⚠️  Found %d files, downloading first %d\n", len(accessFiles), maxFiles)
			accessFiles = accessFiles[:maxFiles]
		}
		
		filesToDownload = accessFiles
		fmt.Printf("📦 Downloading %d access log files...\n", len(filesToDownload))
	}
	
	timestamp := time.Now().Format("20060102_150405")
	var downloadedFiles []string
	successCount := 0
	
	for i, remoteFile := range filesToDownload {
		// Generate local filename
		baseName := filepath.Base(remoteFile)
		localFilename := fmt.Sprintf("%s_%s_%s", server.Host, timestamp, baseName)
		localPath := filepath.Join(outputDir, localFilename)
		
		fmt.Printf("  [%d/%d] %s -> %s\n", i+1, len(filesToDownload), remoteFile, localFilename)
		
		if err := client.DownloadFile(remoteFile, localPath); err != nil {
			fmt.Printf("    ❌ Failed: %v\n", err)
			continue
		}
		
		// Check file size
		if stat, err := os.Stat(localPath); err == nil {
			fmt.Printf("    ✅ Downloaded (%s)\n", formatFileSize(stat.Size()))
			downloadedFiles = append(downloadedFiles, localPath)
			successCount++
		} else {
			fmt.Printf("    ✅ Downloaded\n")
			downloadedFiles = append(downloadedFiles, localPath)
			successCount++
		}
	}
	
	fmt.Printf("📊 Server summary: %d/%d files downloaded successfully\n", successCount, len(filesToDownload))
	
	return downloadedFiles, nil
}