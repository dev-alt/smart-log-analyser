package cmd

import (
	"fmt"
	"smart-log-analyser/pkg/config"
)

// createTestPreset creates a simple working test preset
func createTestPreset() {
	configManager := config.NewConfigManager("config")
	if err := configManager.Load(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Create a simple working preset
	testPreset := config.AnalysisPreset{
		Name:        "simple-traffic-test",
		Description: "Simple traffic analysis test preset",
		Category:    "test",
		Query:       "SELECT ip, COUNT() FROM logs GROUP BY ip ORDER BY COUNT() DESC LIMIT 10",
		Filters:     config.PresetFilters{},
		Exports: []config.ExportConfig{
			{Format: "csv", Filename: "simple-traffic.csv", AutoOpen: false},
		},
		Charts: []config.ChartConfig{
			{Type: "bar", Title: "Top IPs", Width: 80, Height: 20, Colors: true, Enabled: true},
		},
	}

	if err := configManager.AddPreset(testPreset); err != nil {
		fmt.Printf("Error adding preset: %v\n", err)
		return
	}

	fmt.Println("✅ Added simple test preset")
}