package config

import (
	"time"
)

// AppConfig represents the main application configuration
type AppConfig struct {
	Analysis    AnalysisConfig    `yaml:"analysis"`
	Servers     []ServerProfile   `yaml:"servers"`
	Templates   []ReportTemplate  `yaml:"templates"`
	Presets     []AnalysisPreset  `yaml:"presets"`
	Preferences UserPreferences   `yaml:"preferences"`
	Version     string            `yaml:"version"`
}

// AnalysisConfig holds default analysis settings
type AnalysisConfig struct {
	DefaultTopIPs    int      `yaml:"default_top_ips"`
	DefaultTopURLs   int      `yaml:"default_top_urls"`
	DefaultTimeRange string   `yaml:"default_time_range"`
	AutoCharts       bool     `yaml:"auto_charts"`
	ChartWidth       int      `yaml:"chart_width"`
	NoColors         bool     `yaml:"no_colors"`
	ExportFormats    []string `yaml:"export_formats"`
	ShowDetails      bool     `yaml:"show_details"`
	TrendAnalysis    bool     `yaml:"trend_analysis"`
}

// ServerProfile represents a server connection configuration
type ServerProfile struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password,omitempty"`
	KeyFile  string `yaml:"key_file,omitempty"`
	LogPath  string `yaml:"log_path"`
	Tags     []string `yaml:"tags,omitempty"`
}

// AnalysisPreset represents a saved analysis configuration
type AnalysisPreset struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Category    string        `yaml:"category"`
	Query       string        `yaml:"query,omitempty"`
	Filters     PresetFilters `yaml:"filters"`
	Exports     []ExportConfig `yaml:"exports"`
	Charts      []ChartConfig  `yaml:"charts"`
	CreatedAt   time.Time     `yaml:"created_at"`
	UpdatedAt   time.Time     `yaml:"updated_at"`
}

// PresetFilters holds filtering configuration for presets
type PresetFilters struct {
	Since         string   `yaml:"since,omitempty"`
	Until         string   `yaml:"until,omitempty"`
	StatusCodes   []int    `yaml:"status_codes,omitempty"`
	Methods       []string `yaml:"methods,omitempty"`
	IPs           []string `yaml:"ips,omitempty"`
	URLs          []string `yaml:"urls,omitempty"`
	ExcludeIPs    []string `yaml:"exclude_ips,omitempty"`
	ExcludeURLs   []string `yaml:"exclude_urls,omitempty"`
	MinSize       int64    `yaml:"min_size,omitempty"`
	MaxSize       int64    `yaml:"max_size,omitempty"`
}

// ExportConfig defines export settings for presets
type ExportConfig struct {
	Format   string `yaml:"format"` // json, csv, html
	Filename string `yaml:"filename,omitempty"`
	Template string `yaml:"template,omitempty"`
	AutoOpen bool   `yaml:"auto_open"`
}

// ChartConfig defines chart settings for presets
type ChartConfig struct {
	Type     string `yaml:"type"`     // bar, line, pie
	Title    string `yaml:"title"`
	Width    int    `yaml:"width"`
	Height   int    `yaml:"height"`
	Colors   bool   `yaml:"colors"`
	Enabled  bool   `yaml:"enabled"`
}

// ReportTemplate represents a custom report template
type ReportTemplate struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Category    string            `yaml:"category"`
	Sections    []TemplateSection `yaml:"sections"`
	Style       TemplateStyle     `yaml:"style"`
	CreatedAt   time.Time         `yaml:"created_at"`
	UpdatedAt   time.Time         `yaml:"updated_at"`
}

// TemplateSection defines a section in a report template
type TemplateSection struct {
	Name    string              `yaml:"name"`
	Type    string              `yaml:"type"` // stats, chart, table, text
	Query   string              `yaml:"query,omitempty"`
	Config  map[string]interface{} `yaml:"config,omitempty"`
	Order   int                 `yaml:"order"`
	Enabled bool                `yaml:"enabled"`
}

// TemplateStyle defines styling options for templates
type TemplateStyle struct {
	Theme       string            `yaml:"theme"`       // light, dark, minimal
	Colors      map[string]string `yaml:"colors,omitempty"`
	Fonts       map[string]string `yaml:"fonts,omitempty"`
	Layout      string            `yaml:"layout"`      // single, multi-column
	ShowLogo    bool              `yaml:"show_logo"`
	CustomCSS   string            `yaml:"custom_css,omitempty"`
}

// UserPreferences holds user-specific settings
type UserPreferences struct {
	DefaultExportDir string `yaml:"default_export_dir"`
	DefaultConfigDir string `yaml:"default_config_dir"`
	AutoSave         bool   `yaml:"auto_save"`
	ShowTips         bool   `yaml:"show_tips"`
	Theme            string `yaml:"theme"`
	Language         string `yaml:"language"`
	Timezone         string `yaml:"timezone"`
	DateFormat       string `yaml:"date_format"`
	TimeFormat       string `yaml:"time_format"`
}

// PresetCategory represents preset categories
type PresetCategory struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Icon        string `yaml:"icon,omitempty"`
	Color       string `yaml:"color,omitempty"`
}

// ConfigValidationError represents configuration validation errors
type ConfigValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ConfigValidationError) Error() string {
	return e.Field + ": " + e.Message
}