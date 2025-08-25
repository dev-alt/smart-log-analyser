package config

import (
	"fmt"
	"time"
)

// GetBuiltinTemplates returns built-in report templates
func GetBuiltinTemplates() []ReportTemplate {
	now := time.Now()
	
	return []ReportTemplate{
		{
			Name:        "security-report",
			Description: "Comprehensive security analysis report template",
			Category:    "security",
			Sections: []TemplateSection{
				{
					Name:    "Executive Summary",
					Type:    "text",
					Config:  map[string]interface{}{"content": "Security Analysis Overview"},
					Order:   1,
					Enabled: true,
				},
				{
					Name:    "Failed Login Attempts",
					Type:    "table",
					Query:   "SELECT ip, COUNT() as attempts FROM logs WHERE status IN (401, 403) GROUP BY ip HAVING attempts > 3 ORDER BY attempts DESC LIMIT 10",
					Order:   2,
					Enabled: true,
				},
				{
					Name:    "Attack Patterns Chart",
					Type:    "chart",
					Query:   "SELECT status, COUNT() as count FROM logs WHERE status >= 400 GROUP BY status ORDER BY count DESC",
					Config: map[string]interface{}{
						"chart_type": "bar",
						"title":      "Attack Patterns by Status Code",
						"color":      "#dc3545",
					},
					Order:   3,
					Enabled: true,
				},
				{
					Name:    "Suspicious IPs",
					Type:    "table",
					Query:   "SELECT ip, COUNT() as requests, COUNT(DISTINCT url) as unique_urls FROM logs WHERE status >= 400 GROUP BY ip ORDER BY requests DESC LIMIT 15",
					Order:   4,
					Enabled: true,
				},
				{
					Name:    "Security Recommendations",
					Type:    "text",
					Config: map[string]interface{}{
						"content": "Based on the analysis, consider implementing rate limiting and monitoring suspicious IP addresses.",
					},
					Order:   5,
					Enabled: true,
				},
			},
			Style: TemplateStyle{
				Theme:     "light",
				Colors:    map[string]string{"primary": "#dc3545", "secondary": "#6c757d"},
				Layout:    "single",
				ShowLogo:  true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "performance-report",
			Description: "Performance analysis and optimization report template",
			Category:    "performance",
			Sections: []TemplateSection{
				{
					Name:    "Performance Overview",
					Type:    "stats",
					Query:   "SELECT COUNT() as total_requests, AVG(size) as avg_response_size, COUNT(CASE WHEN status >= 400 THEN 1 END) as error_rate FROM logs",
					Order:   1,
					Enabled: true,
				},
				{
					Name:    "Response Time Analysis",
					Type:    "chart",
					Query:   "SELECT HOUR(timestamp) as hour, AVG(size) as avg_size FROM logs GROUP BY hour ORDER BY hour",
					Config: map[string]interface{}{
						"chart_type": "line",
						"title":      "Average Response Size by Hour",
						"color":      "#ffc107",
					},
					Order:   2,
					Enabled: true,
				},
				{
					Name:    "Slowest Endpoints",
					Type:    "table",
					Query:   "SELECT url, AVG(size) as avg_size, COUNT() as requests FROM logs WHERE status = 200 GROUP BY url ORDER BY avg_size DESC LIMIT 10",
					Order:   3,
					Enabled: true,
				},
				{
					Name:    "Error Analysis",
					Type:    "chart",
					Query:   "SELECT status, COUNT() as count FROM logs WHERE status >= 400 GROUP BY status ORDER BY count DESC",
					Config: map[string]interface{}{
						"chart_type": "pie",
						"title":      "Error Distribution",
						"color":      "#dc3545",
					},
					Order:   4,
					Enabled: true,
				},
				{
					Name:    "Performance Recommendations",
					Type:    "text",
					Config: map[string]interface{}{
						"content": "Consider optimizing endpoints with high response sizes and implementing caching strategies.",
					},
					Order:   5,
					Enabled: true,
				},
			},
			Style: TemplateStyle{
				Theme:     "light",
				Colors:    map[string]string{"primary": "#ffc107", "secondary": "#6c757d"},
				Layout:    "multi-column",
				ShowLogo:  true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "traffic-report",
			Description: "Traffic analysis and user behavior report template",
			Category:    "traffic",
			Sections: []TemplateSection{
				{
					Name:    "Traffic Summary",
					Type:    "stats",
					Query:   "SELECT COUNT() as total_requests, COUNT(DISTINCT ip) as unique_visitors, SUM(size) as total_bytes FROM logs",
					Order:   1,
					Enabled: true,
				},
				{
					Name:    "Hourly Traffic Pattern",
					Type:    "chart",
					Query:   "SELECT HOUR(timestamp) as hour, COUNT() as requests FROM logs GROUP BY hour ORDER BY hour",
					Config: map[string]interface{}{
						"chart_type": "line",
						"title":      "Requests by Hour",
						"color":      "#28a745",
					},
					Order:   2,
					Enabled: true,
				},
				{
					Name:    "Top Visitors",
					Type:    "table",
					Query:   "SELECT ip, COUNT() as requests, SUM(size) as total_bytes FROM logs GROUP BY ip ORDER BY requests DESC LIMIT 10",
					Order:   3,
					Enabled: true,
				},
				{
					Name:    "Content Type Distribution",
					Type:    "chart",
					Query:   "SELECT CASE WHEN url LIKE '%.css' THEN 'CSS' WHEN url LIKE '%.js' THEN 'JavaScript' WHEN url LIKE '%.png' OR url LIKE '%.jpg' THEN 'Images' ELSE 'Dynamic' END as type, COUNT() as requests FROM logs GROUP BY type ORDER BY requests DESC",
					Config: map[string]interface{}{
						"chart_type": "pie",
						"title":      "Content Type Distribution",
						"color":      "#28a745",
					},
					Order:   4,
					Enabled: true,
				},
				{
					Name:    "Top Requested Pages",
					Type:    "table",
					Query:   "SELECT url, COUNT() as requests FROM logs WHERE status = 200 GROUP BY url ORDER BY requests DESC LIMIT 15",
					Order:   5,
					Enabled: true,
				},
			},
			Style: TemplateStyle{
				Theme:     "light",
				Colors:    map[string]string{"primary": "#28a745", "secondary": "#6c757d"},
				Layout:    "single",
				ShowLogo:  true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "executive-summary",
			Description: "High-level executive summary report template",
			Category:    "general",
			Sections: []TemplateSection{
				{
					Name:    "Key Metrics",
					Type:    "stats",
					Query:   "SELECT COUNT() as total_requests, COUNT(DISTINCT ip) as unique_visitors, COUNT(CASE WHEN status >= 400 THEN 1 END) as errors, SUM(size) as total_bytes FROM logs",
					Order:   1,
					Enabled: true,
				},
				{
					Name:    "Traffic Trend",
					Type:    "chart",
					Query:   "SELECT DATE(timestamp) as date, COUNT() as requests FROM logs GROUP BY date ORDER BY date",
					Config: map[string]interface{}{
						"chart_type": "line",
						"title":      "Daily Traffic Trend",
						"color":      "#007bff",
					},
					Order:   2,
					Enabled: true,
				},
				{
					Name:    "Status Code Distribution",
					Type:    "chart",
					Query:   "SELECT CASE WHEN status < 300 THEN 'Success' WHEN status < 400 THEN 'Redirect' WHEN status < 500 THEN 'Client Error' ELSE 'Server Error' END as category, COUNT() as count FROM logs GROUP BY category ORDER BY count DESC",
					Config: map[string]interface{}{
						"chart_type": "pie",
						"title":      "Response Status Distribution",
						"color":      "#17a2b8",
					},
					Order:   3,
					Enabled: true,
				},
				{
					Name:    "Executive Summary",
					Type:    "text",
					Config: map[string]interface{}{
						"content": "This report provides a comprehensive overview of web server activity, including traffic patterns, error rates, and performance metrics.",
					},
					Order:   4,
					Enabled: true,
				},
			},
			Style: TemplateStyle{
				Theme:     "light",
				Colors:    map[string]string{"primary": "#007bff", "secondary": "#6c757d"},
				Layout:    "multi-column",
				ShowLogo:  true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Name:        "detailed-analysis",
			Description: "Comprehensive detailed analysis report template",
			Category:    "general",
			Sections: []TemplateSection{
				{
					Name:    "Analysis Overview",
					Type:    "stats",
					Query:   "SELECT COUNT() as total_requests, COUNT(DISTINCT ip) as unique_visitors, MIN(timestamp) as first_request, MAX(timestamp) as last_request FROM logs",
					Order:   1,
					Enabled: true,
				},
				{
					Name:    "Top IP Addresses",
					Type:    "table",
					Query:   "SELECT ip, COUNT() as requests, COUNT(DISTINCT url) as unique_pages, SUM(size) as total_bytes FROM logs GROUP BY ip ORDER BY requests DESC LIMIT 20",
					Order:   2,
					Enabled: true,
				},
				{
					Name:    "Most Requested URLs",
					Type:    "table",
					Query:   "SELECT url, COUNT() as requests, AVG(size) as avg_size FROM logs GROUP BY url ORDER BY requests DESC LIMIT 20",
					Order:   3,
					Enabled: true,
				},
				{
					Name:    "HTTP Methods Distribution",
					Type:    "chart",
					Query:   "SELECT method, COUNT() as count FROM logs GROUP BY method ORDER BY count DESC",
					Config: map[string]interface{}{
						"chart_type": "bar",
						"title":      "HTTP Methods Usage",
						"color":      "#6f42c1",
					},
					Order:   4,
					Enabled: true,
				},
				{
					Name:    "Error Analysis",
					Type:    "table",
					Query:   "SELECT url, status, COUNT() as occurrences FROM logs WHERE status >= 400 GROUP BY url, status ORDER BY occurrences DESC LIMIT 15",
					Order:   5,
					Enabled: true,
				},
				{
					Name:    "Bandwidth Usage Over Time",
					Type:    "chart",
					Query:   "SELECT HOUR(timestamp) as hour, SUM(size) as total_bytes FROM logs GROUP BY hour ORDER BY hour",
					Config: map[string]interface{}{
						"chart_type": "line",
						"title":      "Hourly Bandwidth Usage",
						"color":      "#fd7e14",
					},
					Order:   6,
					Enabled: true,
				},
			},
			Style: TemplateStyle{
				Theme:     "light",
				Colors:    map[string]string{"primary": "#6f42c1", "secondary": "#6c757d"},
				Layout:    "single",
				ShowLogo:  true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// InstallBuiltinTemplates installs built-in templates to the configuration
func (cm *ConfigManager) InstallBuiltinTemplates() error {
	builtinTemplates := GetBuiltinTemplates()
	config := cm.GetConfig()

	// Check for existing templates and only add new ones
	existingNames := make(map[string]bool)
	for _, template := range config.Templates {
		existingNames[template.Name] = true
	}

	addedCount := 0
	for _, template := range builtinTemplates {
		if !existingNames[template.Name] {
			config.Templates = append(config.Templates, template)
			addedCount++
		}
	}

	if addedCount > 0 {
		return cm.Save()
	}

	return nil // No new templates to add
}

// AddTemplate adds a new report template
func (cm *ConfigManager) AddTemplate(template ReportTemplate) error {
	config := cm.GetConfig()
	
	// Check for duplicate names
	for _, existingTemplate := range config.Templates {
		if existingTemplate.Name == template.Name {
			return fmt.Errorf("template with name '%s' already exists", template.Name)
		}
	}

	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	config.Templates = append(config.Templates, template)
	
	return cm.Save()
}

// UpdateTemplate updates an existing template
func (cm *ConfigManager) UpdateTemplate(name string, template ReportTemplate) error {
	config := cm.GetConfig()
	
	for i, existingTemplate := range config.Templates {
		if existingTemplate.Name == name {
			template.CreatedAt = existingTemplate.CreatedAt
			template.UpdatedAt = time.Now()
			config.Templates[i] = template
			return cm.Save()
		}
	}
	
	return fmt.Errorf("template '%s' not found", name)
}

// DeleteTemplate removes a template
func (cm *ConfigManager) DeleteTemplate(name string) error {
	config := cm.GetConfig()
	
	for i, template := range config.Templates {
		if template.Name == name {
			config.Templates = append(config.Templates[:i], config.Templates[i+1:]...)
			return cm.Save()
		}
	}
	
	return fmt.Errorf("template '%s' not found", name)
}

// GetTemplate retrieves a template by name
func (cm *ConfigManager) GetTemplate(name string) (*ReportTemplate, error) {
	config := cm.GetConfig()
	
	for _, template := range config.Templates {
		if template.Name == name {
			return &template, nil
		}
	}
	
	return nil, fmt.Errorf("template '%s' not found", name)
}

// GetTemplatesByCategory retrieves templates by category
func (cm *ConfigManager) GetTemplatesByCategory(category string) []ReportTemplate {
	config := cm.GetConfig()
	var templates []ReportTemplate
	
	for _, template := range config.Templates {
		if template.Category == category {
			templates = append(templates, template)
		}
	}
	
	return templates
}

// GetTemplateTemplate returns a template for creating new templates
func GetTemplateTemplate(category string) ReportTemplate {
	now := time.Now()
	
	return ReportTemplate{
		Name:        "",
		Description: "",
		Category:    category,
		Sections: []TemplateSection{
			{
				Name:    "Overview",
				Type:    "stats",
				Query:   "SELECT COUNT() as total_requests FROM logs",
				Order:   1,
				Enabled: true,
			},
			{
				Name:    "Analysis Chart",
				Type:    "chart",
				Query:   "SELECT status, COUNT() as count FROM logs GROUP BY status ORDER BY count DESC",
				Config: map[string]interface{}{
					"chart_type": "bar",
					"title":      "Analysis Results",
					"color":      "#007bff",
				},
				Order:   2,
				Enabled: true,
			},
			{
				Name:    "Detailed Results",
				Type:    "table",
				Query:   "SELECT * FROM logs LIMIT 10",
				Order:   3,
				Enabled: true,
			},
		},
		Style: TemplateStyle{
			Theme:     "light",
			Colors:    map[string]string{"primary": "#007bff", "secondary": "#6c757d"},
			Layout:    "single",
			ShowLogo:  true,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}