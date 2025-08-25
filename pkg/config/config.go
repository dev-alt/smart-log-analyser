package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigFile = "config/app.yaml"
	ConfigVersion     = "1.0.0"
)

// ConfigManager handles configuration operations
type ConfigManager struct {
	configDir  string
	configFile string
	config     *AppConfig
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configDir string) *ConfigManager {
	if configDir == "" {
		configDir = "config"
	}
	
	return &ConfigManager{
		configDir:  configDir,
		configFile: filepath.Join(configDir, "app.yaml"),
	}
}

// Load loads the configuration from file
func (cm *ConfigManager) Load() error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(cm.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(cm.configFile); os.IsNotExist(err) {
		// Create default configuration
		cm.config = cm.createDefaultConfig()
		return cm.Save()
	}

	// Load existing configuration
	data, err := os.ReadFile(cm.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	config := &AppConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := cm.validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	cm.config = config
	return nil
}

// Save saves the current configuration to file
func (cm *ConfigManager) Save() error {
	if cm.config == nil {
		return fmt.Errorf("no configuration to save")
	}

	// Update version and timestamp
	cm.config.Version = ConfigVersion

	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *AppConfig {
	if cm.config == nil {
		cm.config = cm.createDefaultConfig()
	}
	return cm.config
}

// SetConfig sets the configuration
func (cm *ConfigManager) SetConfig(config *AppConfig) error {
	if err := cm.validateConfig(config); err != nil {
		return err
	}
	cm.config = config
	return nil
}

// UpdateAnalysisConfig updates analysis configuration
func (cm *ConfigManager) UpdateAnalysisConfig(analysis AnalysisConfig) error {
	config := cm.GetConfig()
	config.Analysis = analysis
	return cm.Save()
}

// AddPreset adds a new analysis preset
func (cm *ConfigManager) AddPreset(preset AnalysisPreset) error {
	config := cm.GetConfig()
	
	// Check for duplicate names
	for _, existingPreset := range config.Presets {
		if existingPreset.Name == preset.Name {
			return fmt.Errorf("preset with name '%s' already exists", preset.Name)
		}
	}

	preset.CreatedAt = time.Now()
	preset.UpdatedAt = time.Now()
	config.Presets = append(config.Presets, preset)
	
	return cm.Save()
}

// UpdatePreset updates an existing preset
func (cm *ConfigManager) UpdatePreset(name string, preset AnalysisPreset) error {
	config := cm.GetConfig()
	
	for i, existingPreset := range config.Presets {
		if existingPreset.Name == name {
			preset.CreatedAt = existingPreset.CreatedAt
			preset.UpdatedAt = time.Now()
			config.Presets[i] = preset
			return cm.Save()
		}
	}
	
	return fmt.Errorf("preset '%s' not found", name)
}

// DeletePreset removes a preset
func (cm *ConfigManager) DeletePreset(name string) error {
	config := cm.GetConfig()
	
	for i, preset := range config.Presets {
		if preset.Name == name {
			config.Presets = append(config.Presets[:i], config.Presets[i+1:]...)
			return cm.Save()
		}
	}
	
	return fmt.Errorf("preset '%s' not found", name)
}

// GetPreset retrieves a preset by name
func (cm *ConfigManager) GetPreset(name string) (*AnalysisPreset, error) {
	config := cm.GetConfig()
	
	for _, preset := range config.Presets {
		if preset.Name == name {
			return &preset, nil
		}
	}
	
	return nil, fmt.Errorf("preset '%s' not found", name)
}

// GetPresetsByCategory retrieves presets by category
func (cm *ConfigManager) GetPresetsByCategory(category string) []AnalysisPreset {
	config := cm.GetConfig()
	var presets []AnalysisPreset
	
	for _, preset := range config.Presets {
		if preset.Category == category {
			presets = append(presets, preset)
		}
	}
	
	return presets
}

// AddServerProfile adds a new server profile
func (cm *ConfigManager) AddServerProfile(profile ServerProfile) error {
	config := cm.GetConfig()
	
	// Check for duplicate names
	for _, existingProfile := range config.Servers {
		if existingProfile.Name == profile.Name {
			return fmt.Errorf("server profile with name '%s' already exists", profile.Name)
		}
	}

	// Set defaults
	if profile.Port == 0 {
		profile.Port = 22
	}
	if profile.LogPath == "" {
		profile.LogPath = "/var/log/nginx/access.log"
	}

	config.Servers = append(config.Servers, profile)
	return cm.Save()
}

// UpdateServerProfile updates an existing server profile
func (cm *ConfigManager) UpdateServerProfile(name string, profile ServerProfile) error {
	config := cm.GetConfig()
	
	for i, existingProfile := range config.Servers {
		if existingProfile.Name == name {
			profile.Name = name // Ensure name consistency
			config.Servers[i] = profile
			return cm.Save()
		}
	}
	
	return fmt.Errorf("server profile '%s' not found", name)
}

// DeleteServerProfile removes a server profile
func (cm *ConfigManager) DeleteServerProfile(name string) error {
	config := cm.GetConfig()
	
	for i, profile := range config.Servers {
		if profile.Name == name {
			config.Servers = append(config.Servers[:i], config.Servers[i+1:]...)
			return cm.Save()
		}
	}
	
	return fmt.Errorf("server profile '%s' not found", name)
}

// GetServerProfile retrieves a server profile by name
func (cm *ConfigManager) GetServerProfile(name string) (*ServerProfile, error) {
	config := cm.GetConfig()
	
	for _, profile := range config.Servers {
		if profile.Name == name {
			return &profile, nil
		}
	}
	
	return nil, fmt.Errorf("server profile '%s' not found", name)
}

// createDefaultConfig creates a default configuration
func (cm *ConfigManager) createDefaultConfig() *AppConfig {
	return &AppConfig{
		Version: ConfigVersion,
		Analysis: AnalysisConfig{
			DefaultTopIPs:    10,
			DefaultTopURLs:   10,
			DefaultTimeRange: "24h",
			AutoCharts:       true,
			ChartWidth:       80,
			NoColors:         false,
			ExportFormats:    []string{"json", "csv"},
			ShowDetails:      false,
			TrendAnalysis:    false,
		},
		Servers:   []ServerProfile{},
		Templates: []ReportTemplate{},
		Presets:   []AnalysisPreset{},
		Preferences: UserPreferences{
			DefaultExportDir: "output",
			DefaultConfigDir: "config",
			AutoSave:         true,
			ShowTips:         true,
			Theme:            "default",
			Language:         "en",
			Timezone:         "UTC",
			DateFormat:       "2006-01-02",
			TimeFormat:       "15:04:05",
		},
	}
}

// validateConfig validates the configuration structure
func (cm *ConfigManager) validateConfig(config *AppConfig) error {
	if config == nil {
		return fmt.Errorf("configuration is nil")
	}

	// Validate analysis config
	if config.Analysis.DefaultTopIPs < 1 {
		return ConfigValidationError{
			Field:   "analysis.default_top_ips",
			Message: "must be greater than 0",
		}
	}

	if config.Analysis.DefaultTopURLs < 1 {
		return ConfigValidationError{
			Field:   "analysis.default_top_urls", 
			Message: "must be greater than 0",
		}
	}

	if config.Analysis.ChartWidth < 20 || config.Analysis.ChartWidth > 200 {
		return ConfigValidationError{
			Field:   "analysis.chart_width",
			Message: "must be between 20 and 200",
		}
	}

	// Validate server profiles
	for i, server := range config.Servers {
		if server.Name == "" {
			return ConfigValidationError{
				Field:   fmt.Sprintf("servers[%d].name", i),
				Message: "server name is required",
			}
		}
		if server.Host == "" {
			return ConfigValidationError{
				Field:   fmt.Sprintf("servers[%d].host", i),
				Message: "server host is required",
			}
		}
		if server.Port < 1 || server.Port > 65535 {
			return ConfigValidationError{
				Field:   fmt.Sprintf("servers[%d].port", i),
				Message: "port must be between 1 and 65535",
			}
		}
	}

	// Validate presets
	for i, preset := range config.Presets {
		if preset.Name == "" {
			return ConfigValidationError{
				Field:   fmt.Sprintf("presets[%d].name", i),
				Message: "preset name is required",
			}
		}
		if preset.Category == "" {
			return ConfigValidationError{
				Field:   fmt.Sprintf("presets[%d].category", i),
				Message: "preset category is required",
			}
		}
	}

	return nil
}

// ConfigDir returns the configuration directory path
func (cm *ConfigManager) ConfigDir() string {
	return cm.configDir
}

// ConfigFile returns the main configuration file path
func (cm *ConfigManager) ConfigFile() string {
	return cm.configFile
}