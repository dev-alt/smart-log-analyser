package ipc

import (
	"smart-log-analyser/pkg/analyser"
	"smart-log-analyser/pkg/config"
)

// IPCRequest represents a request from the client
type IPCRequest struct {
	ID      string                 `json:"id"`
	Action  string                 `json:"action"`
	LogFile string                 `json:"logFile,omitempty"`
	Options AnalysisOptions        `json:"options,omitempty"`
	Query   string                 `json:"query,omitempty"`
	Config  map[string]interface{} `json:"config,omitempty"`
	Preset  string                 `json:"preset,omitempty"`
}

// IPCResponse represents a response to the client
type IPCResponse struct {
	ID      string      `json:"id"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// AnalysisOptions represents options for log analysis
type AnalysisOptions struct {
	EnableSecurity    bool `json:"enableSecurity"`
	EnablePerformance bool `json:"enablePerformance"`
	EnableTrends      bool `json:"enableTrends"`
	GenerateHTML      bool `json:"generateHtml"`
	Interactive       bool `json:"interactive"`
	HTMLTitle         string `json:"htmlTitle,omitempty"`
	OutputPath        string `json:"outputPath,omitempty"`
}

// AnalysisResultData represents the data returned from analysis
type AnalysisResultData struct {
	Results      *analyser.Results `json:"results"`
	HTMLPath     string           `json:"htmlPath,omitempty"`
	QueryResults interface{}      `json:"queryResults,omitempty"`
	Presets      []config.AnalysisPreset `json:"presets,omitempty"`
	Config       *config.AppConfig       `json:"config,omitempty"`
	Status       string           `json:"status,omitempty"`
}

// Action constants
const (
	ActionAnalyze      = "analyze"
	ActionQuery        = "query"
	ActionListPresets  = "listPresets"
	ActionRunPreset    = "runPreset"
	ActionGetConfig    = "getConfig"
	ActionUpdateConfig = "updateConfig"
	ActionGetStatus    = "getStatus"
	ActionShutdown     = "shutdown"
)

// Error constants
const (
	ErrInvalidAction   = "invalid action"
	ErrMissingLogFile  = "missing log file"
	ErrAnalysisFailed  = "analysis failed"
	ErrConfigFailed    = "configuration operation failed"
	ErrPresetFailed    = "preset operation failed"
	ErrQueryFailed     = "query execution failed"
)