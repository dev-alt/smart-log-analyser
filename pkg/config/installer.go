package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	
	"gopkg.in/yaml.v3"
)

// Installer handles initial setup and configuration installation
type Installer struct {
	configManager *ConfigManager
}

// NewInstaller creates a new configuration installer
func NewInstaller(configDir string) *Installer {
	return &Installer{
		configManager: NewConfigManager(configDir),
	}
}

// Initialize sets up the configuration system with default values
func (i *Installer) Initialize() error {
	// Load or create configuration
	if err := i.configManager.Load(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create directory structure
	if err := i.createDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Install built-in presets
	if err := i.configManager.InstallBuiltinPresets(); err != nil {
		return fmt.Errorf("failed to install built-in presets: %w", err)
	}

	// Install built-in templates
	if err := i.configManager.InstallBuiltinTemplates(); err != nil {
		return fmt.Errorf("failed to install built-in templates: %w", err)
	}

	return nil
}

// createDirectoryStructure creates the necessary configuration directories
func (i *Installer) createDirectoryStructure() error {
	configDir := i.configManager.ConfigDir()
	
	directories := []string{
		configDir,
		filepath.Join(configDir, "presets"),
		filepath.Join(configDir, "templates"),
		filepath.Join(configDir, "profiles"),
		filepath.Join(configDir, "backup"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// Reset removes all configuration and recreates defaults
func (i *Installer) Reset() error {
	configDir := i.configManager.ConfigDir()
	
	// Remove existing configuration directory
	if err := os.RemoveAll(configDir); err != nil {
		return fmt.Errorf("failed to remove configuration directory: %w", err)
	}

	// Reinitialize
	return i.Initialize()
}

// Backup creates a backup of the current configuration
func (i *Installer) Backup() (string, error) {
	configDir := i.configManager.ConfigDir()
	backupDir := filepath.Join(configDir, "backup")
	
	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	backupFile := filepath.Join(backupDir, fmt.Sprintf("config-backup-%s.yaml", timestamp))

	// Copy main configuration file
	configFile := i.configManager.ConfigFile()
	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return "", fmt.Errorf("failed to read config file: %w", err)
		}

		if err := os.WriteFile(backupFile, data, 0644); err != nil {
			return "", fmt.Errorf("failed to write backup file: %w", err)
		}
	}

	return backupFile, nil
}

// Restore restores configuration from a backup file
func (i *Installer) Restore(backupFile string) error {
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupFile)
	}

	// Read backup file
	data, err := os.ReadFile(backupFile)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Write to main configuration file
	configFile := i.configManager.ConfigFile()
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to restore configuration: %w", err)
	}

	// Reload configuration
	return i.configManager.Load()
}

// GetStatus returns the current configuration status
func (i *Installer) GetStatus() (*ConfigStatus, error) {
	configDir := i.configManager.ConfigDir()
	configFile := i.configManager.ConfigFile()

	status := &ConfigStatus{
		ConfigDir:    configDir,
		ConfigFile:   configFile,
		Initialized:  false,
		Presets:      0,
		Templates:    0,
		Servers:      0,
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); err == nil {
		status.Initialized = true
		
		// Load configuration to get counts
		if err := i.configManager.Load(); err == nil {
			config := i.configManager.GetConfig()
			status.Presets = len(config.Presets)
			status.Templates = len(config.Templates)
			status.Servers = len(config.Servers)
		}
	}

	return status, nil
}

// ConfigStatus represents the current configuration status
type ConfigStatus struct {
	ConfigDir   string `json:"config_dir"`
	ConfigFile  string `json:"config_file"`
	Initialized bool   `json:"initialized"`
	Presets     int    `json:"presets"`
	Templates   int    `json:"templates"`
	Servers     int    `json:"servers"`
}

// ExportPresets exports presets to a separate file
func (i *Installer) ExportPresets(filename string) error {
	config := i.configManager.GetConfig()
	
	presetData := map[string]interface{}{
		"version": ConfigVersion,
		"presets": config.Presets,
	}

	data, err := yaml.Marshal(presetData)
	if err != nil {
		return fmt.Errorf("failed to marshal presets: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write presets file: %w", err)
	}

	return nil
}

// ImportPresets imports presets from a file
func (i *Installer) ImportPresets(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("presets file does not exist: %s", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read presets file: %w", err)
	}

	var presetData struct {
		Version string           `yaml:"version"`
		Presets []AnalysisPreset `yaml:"presets"`
	}

	if err := yaml.Unmarshal(data, &presetData); err != nil {
		return fmt.Errorf("failed to parse presets file: %w", err)
	}

	config := i.configManager.GetConfig()
	
	// Add imported presets (check for duplicates)
	existingNames := make(map[string]bool)
	for _, preset := range config.Presets {
		existingNames[preset.Name] = true
	}

	addedCount := 0
	for _, preset := range presetData.Presets {
		if !existingNames[preset.Name] {
			config.Presets = append(config.Presets, preset)
			addedCount++
		}
	}

	if addedCount > 0 {
		return i.configManager.Save()
	}

	return nil
}