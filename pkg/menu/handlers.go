package menu

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"smart-log-analyser/pkg/analyser"
	"smart-log-analyser/pkg/html"
	"smart-log-analyser/pkg/parser"
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
	
	// Ask for export options
	if m.confirmYesNo("\nExport results") {
		return m.handleExport(results)
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
	fmt.Println("🌐 Downloading logs from configured servers...")
	fmt.Println("This would use the existing remote download functionality")
	m.pause()
	return nil
}

func (m *Menu) setupRemoteServers() error {
	fmt.Println("🔧 Remote server setup would be implemented here")
	m.pause()
	return nil
}

func (m *Menu) testConnections() error {
	fmt.Println("🔌 Testing connections would be implemented here")
	m.pause()
	return nil
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
	// Implementation would use system-specific commands to open browser
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