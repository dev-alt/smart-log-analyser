package config

import (
	"time"
)

// GetBuiltinPresets returns the built-in analysis presets
func GetBuiltinPresets() []AnalysisPreset {
	now := time.Now()
	
	return []AnalysisPreset{
		// Security Analysis Presets
		{
			Name:        "security-failed-logins",
			Description: "Detect failed login attempts and suspicious authentication patterns",
			Category:    "security",
			Query:       "SELECT ip, COUNT() as attempts, url FROM logs WHERE status IN (401, 403) GROUP BY ip, url HAVING attempts > 3 ORDER BY attempts DESC",
			Filters: PresetFilters{
				StatusCodes: []int{401, 403},
			},
			Exports: []ExportConfig{
				{Format: "json", Filename: "security-failed-logins.json", AutoOpen: false},
				{Format: "csv", Filename: "security-failed-logins.csv", AutoOpen: false},
			},
			Charts: []ChartConfig{
				{Type: "bar", Title: "Failed Login Attempts by IP", Width: 80, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "security-attack-patterns",
			Description: "Identify potential attack patterns and malicious requests",
			Category:    "security",
			Query:       "SELECT ip, url, method, status, COUNT() as requests FROM logs WHERE (url LIKE '%admin%' OR url LIKE '%wp-admin%' OR url LIKE '%.php%' OR url LIKE '%sql%') GROUP BY ip, url, method, status ORDER BY requests DESC",
			Filters: PresetFilters{
				URLs: []string{"%admin%", "%wp-admin%", "%.php%", "%sql%"},
			},
			Exports: []ExportConfig{
				{Format: "json", Filename: "security-attacks.json", AutoOpen: false},
				{Format: "html", Template: "security-report", AutoOpen: true},
			},
			Charts: []ChartConfig{
				{Type: "bar", Title: "Attack Patterns by IP", Width: 80, Height: 20, Colors: true, Enabled: true},
				{Type: "pie", Title: "Attack Types Distribution", Width: 40, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "security-suspicious-ips",
			Description: "Find IP addresses with unusually high request rates or error patterns",
			Category:    "security",
			Query:       "SELECT ip, COUNT() as total_requests, COUNT(CASE WHEN status >= 400 THEN 1 END) as error_requests FROM logs GROUP BY ip HAVING total_requests > 100 OR error_requests > 10 ORDER BY total_requests DESC",
			Filters: PresetFilters{},
			Exports: []ExportConfig{
				{Format: "csv", Filename: "suspicious-ips.csv", AutoOpen: false},
			},
			Charts: []ChartConfig{
				{Type: "bar", Title: "Request Volume by Suspicious IPs", Width: 80, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},

		// Performance Analysis Presets
		{
			Name:        "performance-slow-endpoints",
			Description: "Identify slow-performing endpoints and large responses",
			Category:    "performance",
			Query:       "SELECT url, AVG(size) as avg_response_size, COUNT() as requests FROM logs WHERE status = 200 GROUP BY url HAVING avg_response_size > 50000 ORDER BY avg_response_size DESC LIMIT 20",
			Filters: PresetFilters{
				StatusCodes: []int{200},
				MinSize:     50000,
			},
			Exports: []ExportConfig{
				{Format: "json", Filename: "slow-endpoints.json", AutoOpen: false},
				{Format: "html", Template: "performance-report", AutoOpen: true},
			},
			Charts: []ChartConfig{
				{Type: "bar", Title: "Average Response Size by Endpoint", Width: 80, Height: 20, Colors: true, Enabled: true},
				{Type: "line", Title: "Performance Trend Over Time", Width: 100, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "performance-error-analysis",
			Description: "Analyze error patterns and response time issues",
			Category:    "performance",
			Query:       "SELECT status, url, COUNT() as error_count FROM logs WHERE status >= 400 GROUP BY status, url ORDER BY error_count DESC LIMIT 20",
			Filters: PresetFilters{
				StatusCodes: []int{400, 401, 403, 404, 500, 502, 503, 504},
			},
			Exports: []ExportConfig{
				{Format: "csv", Filename: "error-analysis.csv", AutoOpen: false},
			},
			Charts: []ChartConfig{
				{Type: "pie", Title: "Error Distribution by Status Code", Width: 60, Height: 20, Colors: true, Enabled: true},
				{Type: "bar", Title: "Top Error URLs", Width: 80, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "performance-resource-usage",
			Description: "Monitor bandwidth usage and resource consumption patterns",
			Category:    "performance",
			Query:       "SELECT HOUR(timestamp) as hour, SUM(size) as total_bytes, COUNT() as requests FROM logs GROUP BY hour ORDER BY hour",
			Filters: PresetFilters{},
			Exports: []ExportConfig{
				{Format: "json", Filename: "resource-usage.json", AutoOpen: false},
			},
			Charts: []ChartConfig{
				{Type: "line", Title: "Hourly Bandwidth Usage", Width: 100, Height: 20, Colors: true, Enabled: true},
				{Type: "bar", Title: "Requests per Hour", Width: 80, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},

		// Traffic Analysis Presets
		{
			Name:        "traffic-peak-analysis",
			Description: "Analyze traffic patterns and identify peak usage periods",
			Category:    "traffic",
			Query:       "SELECT HOUR(timestamp), COUNT() FROM logs GROUP BY HOUR(timestamp) ORDER BY HOUR(timestamp)",
			Filters: PresetFilters{},
			Exports: []ExportConfig{
				{Format: "json", Filename: "traffic-peaks.json", AutoOpen: false},
				{Format: "html", Template: "traffic-report", AutoOpen: true},
			},
			Charts: []ChartConfig{
				{Type: "line", Title: "Hourly Traffic Pattern", Width: 100, Height: 20, Colors: true, Enabled: true},
				{Type: "bar", Title: "Unique Visitors by Hour", Width: 80, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "traffic-user-agents",
			Description: "Analyze user agent patterns and bot traffic identification",
			Category:    "traffic",
			Query:       "SELECT user_agent, COUNT() as requests FROM logs GROUP BY user_agent ORDER BY requests DESC LIMIT 20",
			Filters: PresetFilters{},
			Exports: []ExportConfig{
				{Format: "csv", Filename: "user-agents.csv", AutoOpen: false},
			},
			Charts: []ChartConfig{
				{Type: "pie", Title: "Traffic by User Agent Type", Width: 60, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "traffic-geographic",
			Description: "Geographic distribution analysis of traffic sources",
			Category:    "traffic",
			Query:       "SELECT ip, COUNT() as requests FROM logs WHERE NOT IS_PRIVATE_IP(ip) GROUP BY ip ORDER BY requests DESC LIMIT 30",
			Filters: PresetFilters{},
			Exports: []ExportConfig{
				{Format: "json", Filename: "geographic-traffic.json", AutoOpen: false},
			},
			Charts: []ChartConfig{
				{Type: "bar", Title: "Traffic by Source IP", Width: 80, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "traffic-content-analysis",
			Description: "Analyze content type and resource access patterns",
			Category:    "traffic",
			Query:       "SELECT CASE WHEN url LIKE '%.css' THEN 'CSS' WHEN url LIKE '%.js' THEN 'JavaScript' WHEN url LIKE '%.png' OR url LIKE '%.jpg' OR url LIKE '%.gif' THEN 'Images' WHEN url LIKE '%.pdf' THEN 'Documents' ELSE 'Dynamic Content' END as content_type, COUNT() as requests, SUM(size) as total_bytes FROM logs GROUP BY content_type ORDER BY requests DESC",
			Filters: PresetFilters{},
			Exports: []ExportConfig{
				{Format: "json", Filename: "content-analysis.json", AutoOpen: false},
			},
			Charts: []ChartConfig{
				{Type: "pie", Title: "Requests by Content Type", Width: 60, Height: 20, Colors: true, Enabled: true},
				{Type: "bar", Title: "Bandwidth by Content Type", Width: 80, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		
		// Simple Test Presets (Compatible with current query language)
		{
			Name:        "simple-top-ips",
			Description: "Simple analysis of top requesting IP addresses",
			Category:    "traffic",
			Query:       "SELECT ip, COUNT() FROM logs GROUP BY ip ORDER BY COUNT() DESC LIMIT 10",
			Filters: PresetFilters{},
			Exports: []ExportConfig{
				{Format: "csv", Filename: "top-ips.csv", AutoOpen: false},
			},
			Charts: []ChartConfig{
				{Type: "bar", Title: "Top IP Addresses", Width: 80, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "simple-status-codes",
			Description: "Simple analysis of HTTP status code distribution",
			Category:    "performance",
			Query:       "SELECT status, COUNT() FROM logs GROUP BY status ORDER BY COUNT() DESC",
			Filters: PresetFilters{},
			Exports: []ExportConfig{
				{Format: "json", Filename: "status-codes.json", AutoOpen: false},
			},
			Charts: []ChartConfig{
				{Type: "bar", Title: "Status Code Distribution", Width: 80, Height: 20, Colors: true, Enabled: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// GetPresetCategories returns available preset categories with descriptions
func GetPresetCategories() []PresetCategory {
	return []PresetCategory{
		{
			Name:        "security",
			Description: "Security analysis, threat detection, and attack pattern identification",
			Icon:        "🔒",
			Color:       "#dc3545",
		},
		{
			Name:        "performance",
			Description: "Performance monitoring, bottleneck detection, and resource usage analysis",
			Icon:        "⚡",
			Color:       "#ffc107",
		},
		{
			Name:        "traffic",
			Description: "Traffic pattern analysis, user behavior, and geographic distribution",
			Icon:        "📊",
			Color:       "#28a745",
		},
		{
			Name:        "custom",
			Description: "User-created custom analysis presets",
			Icon:        "⚙️",
			Color:       "#6c757d",
		},
	}
}

// InstallBuiltinPresets installs built-in presets to the configuration
func (cm *ConfigManager) InstallBuiltinPresets() error {
	builtinPresets := GetBuiltinPresets()
	config := cm.GetConfig()

	// Check for existing presets and only add new ones
	existingNames := make(map[string]bool)
	for _, preset := range config.Presets {
		existingNames[preset.Name] = true
	}

	addedCount := 0
	for _, preset := range builtinPresets {
		if !existingNames[preset.Name] {
			config.Presets = append(config.Presets, preset)
			addedCount++
		}
	}

	if addedCount > 0 {
		return cm.Save()
	}

	return nil // No new presets to add
}

// GetPresetTemplate returns a template for creating new presets
func GetPresetTemplate(category string) AnalysisPreset {
	now := time.Now()
	
	template := AnalysisPreset{
		Name:        "",
		Description: "",
		Category:    category,
		Query:       "SELECT * FROM logs WHERE status = 200",
		Filters: PresetFilters{
			Since: "24h",
		},
		Exports: []ExportConfig{
			{Format: "json", AutoOpen: false},
		},
		Charts: []ChartConfig{
			{Type: "bar", Title: "Analysis Results", Width: 80, Height: 20, Colors: true, Enabled: true},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Customize template based on category
	switch category {
	case "security":
		template.Query = "SELECT ip, COUNT() as requests FROM logs WHERE status >= 400 GROUP BY ip ORDER BY requests DESC"
		template.Filters.StatusCodes = []int{400, 401, 403, 404, 500}
	case "performance":
		template.Query = "SELECT url, AVG(size) as avg_size, COUNT() as requests FROM logs WHERE status = 200 GROUP BY url ORDER BY avg_size DESC"
		template.Filters.StatusCodes = []int{200}
	case "traffic":
		template.Query = "SELECT HOUR(timestamp) as hour, COUNT() as requests FROM logs GROUP BY hour ORDER BY hour"
	}

	return template
}