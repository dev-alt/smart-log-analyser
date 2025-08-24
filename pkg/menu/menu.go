package menu

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Menu represents the interactive menu system
type Menu struct {
	scanner *bufio.Scanner
}

// New creates a new menu system
func New() *Menu {
	return &Menu{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Run starts the interactive menu system
func (m *Menu) Run() error {
	m.showWelcome()
	
	for {
		choice, err := m.showMainMenu()
		if err != nil {
			return err
		}
		
		switch choice {
		case 1:
			if err := m.handleLocalAnalysis(); err != nil {
				m.showError("Local analysis error", err)
			}
		case 2:
			if err := m.handleRemoteAnalysis(); err != nil {
				m.showError("Remote analysis error", err)
			}
		case 3:
			if err := m.handleHTMLReport(); err != nil {
				m.showError("HTML report error", err)
			}
		case 4:
			if err := m.handleConfiguration(); err != nil {
				m.showError("Configuration error", err)
			}
		case 5:
			m.showHelp()
		case 6:
			m.showGoodbye()
			return nil
		default:
			fmt.Println("❌ Invalid choice. Please try again.\n")
		}
	}
}

// showWelcome displays the welcome screen
func (m *Menu) showWelcome() {
	m.clearScreen()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                  Smart Log Analyser v1.0                    ║")
	fmt.Println("║                Interactive Menu System                       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("🚀 Welcome to the Smart Log Analyser!")
	fmt.Println("   Advanced Nginx log analysis with interactive reporting")
	fmt.Println()
	m.pauseForEffect()
}

// showMainMenu displays the main menu and returns the user's choice
func (m *Menu) showMainMenu() (int, error) {
	fmt.Println("📊 What would you like to do?")
	fmt.Println()
	fmt.Println("1. 📂 Analyse Local Log Files")
	fmt.Println("2. 🌐 Download & Analyse Remote Logs")
	fmt.Println("3. 📈 Generate HTML Report")
	fmt.Println("4. 🔧 Configuration & Setup")
	fmt.Println("5. 📚 Help & Documentation")
	fmt.Println("6. 🚪 Exit")
	fmt.Println()
	
	return m.getIntInput("Enter your choice (1-6): ", 1, 6)
}

// handleLocalAnalysis handles local log file analysis
func (m *Menu) handleLocalAnalysis() error {
	m.clearScreen()
	fmt.Println("📂 Local Log Analysis")
	fmt.Println("════════════════════")
	fmt.Println()
	fmt.Println("Available options:")
	fmt.Println("1. Quick analysis (current directory *.log files)")
	fmt.Println("2. Select specific files")
	fmt.Println("3. Analyse with time range filter")
	fmt.Println("4. Advanced analysis with all options")
	fmt.Println("5. Back to main menu")
	fmt.Println()
	
	choice, err := m.getIntInput("Enter choice (1-5): ", 1, 5)
	if err != nil {
		return err
	}
	
	if choice == 5 {
		return nil // Back to main menu
	}
	
	var files []string
	var since, until *time.Time
	showDetails := false
	
	switch choice {
	case 1:
		files = m.findLogFilesIntelligent()
		if len(files) == 0 {
			fmt.Println("❌ No log files found in common locations")
			fmt.Println("   Searched: ./downloads/, ./logs/, current directory")
			m.pause()
			return nil
		}
		location := m.getSourceLocation(files)
		fmt.Printf("📁 Found %d log files in %s\n", len(files), location)
		for i, file := range files {
			if i < 5 { // Show first 5 files
				fmt.Printf("  • %s\n", filepath.Base(file))
			} else if i == 5 {
				fmt.Printf("  • ... and %d more files\n", len(files)-5)
				break
			}
		}
		
	case 2:
		files, err = m.selectLogFiles()
		if err != nil {
			return err
		}
		
	case 3:
		files, err = m.selectLogFiles()
		if err != nil {
			return err
		}
		since, until, err = m.getTimeRange()
		if err != nil {
			return err
		}
		
	case 4:
		files, err = m.selectLogFiles()
		if err != nil {
			return err
		}
		since, until, err = m.getTimeRange()
		if err != nil {
			return err
		}
		showDetails = m.confirmDetails()
	}
	
	if len(files) == 0 {
		fmt.Println("❌ No files selected for analysis")
		m.pause()
		return nil
	}
	
	return m.performAnalysis(files, since, until, showDetails)
}

// handleRemoteAnalysis handles remote log analysis
func (m *Menu) handleRemoteAnalysis() error {
	m.clearScreen()
	fmt.Println("🌐 Remote Log Management")
	fmt.Println("═══════════════════════")
	fmt.Println()
	fmt.Println("Available options:")
	fmt.Println("1. Download logs from configured servers")
	fmt.Println("2. Setup/configure remote servers")
	fmt.Println("3. Test connections")
	fmt.Println("4. Download and analyse immediately")
	fmt.Println("5. Back to main menu")
	fmt.Println()
	
	choice, err := m.getIntInput("Enter choice (1-5): ", 1, 5)
	if err != nil {
		return err
	}
	
	switch choice {
	case 1:
		return m.downloadLogs(false)
	case 2:
		return m.setupRemoteServers()
	case 3:
		return m.testConnections()
	case 4:
		return m.downloadAndAnalyse()
	case 5:
		return nil // Back to main menu
	}
	
	return nil
}

// handleHTMLReport handles HTML report generation
func (m *Menu) handleHTMLReport() error {
	m.clearScreen()
	fmt.Println("📈 HTML Report Generation")
	fmt.Println("════════════════════════")
	fmt.Println()
	fmt.Println("Available options:")
	fmt.Println("1. Quick report from recent analysis")
	fmt.Println("2. Analyse files and generate report")
	fmt.Println("3. Batch report generation")
	fmt.Println("4. Custom report settings")
	fmt.Println("5. Back to main menu")
	fmt.Println()
	
	choice, err := m.getIntInput("Enter choice (1-5): ", 1, 5)
	if err != nil {
		return err
	}
	
	switch choice {
	case 1:
		return m.quickHTMLReport()
	case 2:
		return m.analyseAndGenerateReport()
	case 3:
		return m.batchReportGeneration()
	case 4:
		return m.customReportSettings()
	case 5:
		return nil // Back to main menu
	}
	
	return nil
}

// handleConfiguration handles configuration and setup
func (m *Menu) handleConfiguration() error {
	m.clearScreen()
	fmt.Println("🔧 Configuration & Setup")
	fmt.Println("═══════════════════════")
	fmt.Println()
	fmt.Println("Available options:")
	fmt.Println("1. Setup remote server connections")
	fmt.Println("2. Configure analysis preferences")
	fmt.Println("3. Set default export locations")
	fmt.Println("4. View current configuration")
	fmt.Println("5. Back to main menu")
	fmt.Println()
	
	choice, err := m.getIntInput("Enter choice (1-5): ", 1, 5)
	if err != nil {
		return err
	}
	
	switch choice {
	case 1:
		return m.setupRemoteServers()
	case 2:
		return m.configureAnalysisPreferences()
	case 3:
		return m.setExportLocations()
	case 4:
		return m.viewConfiguration()
	case 5:
		return nil // Back to main menu
	}
	
	return nil
}

// Helper functions

// clearScreen clears the terminal screen
func (m *Menu) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// pauseForEffect creates a small pause for better UX
func (m *Menu) pauseForEffect() {
	time.Sleep(500 * time.Millisecond)
}

// pause waits for user input to continue
func (m *Menu) pause() {
	fmt.Print("\nPress Enter to continue...")
	m.scanner.Scan()
}

// getIntInput gets integer input from user within a range
func (m *Menu) getIntInput(prompt string, min, max int) (int, error) {
	for {
		fmt.Print(prompt)
		if !m.scanner.Scan() {
			return 0, fmt.Errorf("failed to read input")
		}
		
		input := strings.TrimSpace(m.scanner.Text())
		if input == "" {
			continue
		}
		
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("❌ Please enter a number between %d and %d\n", min, max)
			continue
		}
		
		if num < min || num > max {
			fmt.Printf("❌ Please enter a number between %d and %d\n", min, max)
			continue
		}
		
		return num, nil
	}
}

// getStringInput gets string input from user
func (m *Menu) getStringInput(prompt string) string {
	fmt.Print(prompt)
	m.scanner.Scan()
	return strings.TrimSpace(m.scanner.Text())
}

// confirmYesNo gets yes/no confirmation from user
func (m *Menu) confirmYesNo(prompt string) bool {
	for {
		response := strings.ToLower(m.getStringInput(prompt + " (y/n): "))
		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("❌ Please enter 'y' for yes or 'n' for no")
		}
	}
}

// showError displays an error message
func (m *Menu) showError(context string, err error) {
	fmt.Printf("❌ %s: %v\n\n", context, err)
	m.pause()
}

// showSuccess displays a success message
func (m *Menu) showSuccess(message string) {
	fmt.Printf("✅ %s\n\n", message)
	m.pause()
}

// showGoodbye displays goodbye message
func (m *Menu) showGoodbye() {
	m.clearScreen()
	fmt.Println("👋 Thank you for using Smart Log Analyser!")
	fmt.Println("   Have a great day!")
	fmt.Println()
}

// showHelp displays help information
func (m *Menu) showHelp() {
	m.clearScreen()
	fmt.Println("📚 Help & Documentation")
	fmt.Println("══════════════════════")
	fmt.Println()
	fmt.Println("🎯 Smart Log Analyser Features:")
	fmt.Println("  • Advanced Nginx log analysis")
	fmt.Println("  • Interactive HTML reports with charts")
	fmt.Println("  • Bot detection and traffic analysis")
	fmt.Println("  • Security threat detection")
	fmt.Println("  • Geographic IP analysis")
	fmt.Println("  • Remote server log downloading")
	fmt.Println()
	fmt.Println("📖 Documentation:")
	fmt.Println("  • GitHub: https://github.com/dev-alt/smart-log-analyser")
	fmt.Println("  • CLI Help: ./smart-log-analyser --help")
	fmt.Println("  • Command Help: ./smart-log-analyser [command] --help")
	fmt.Println()
	fmt.Println("🔧 Common CLI Commands:")
	fmt.Println("  • ./smart-log-analyser analyse logs/ --details")
	fmt.Println("  • ./smart-log-analyser analyse logs/ --export-html=output/report.html")
	fmt.Println("  • ./smart-log-analyser download")
	fmt.Println()
	m.pause()
}

// findLogFiles finds log files in a directory
func (m *Menu) findLogFiles(dir string) []string {
	var files []string
	
	patterns := []string{"*.log", "*.log.*", "*.gz"}
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(dir, pattern))
		files = append(files, matches...)
	}
	
	return files
}

// findLogFilesIntelligent searches for log files in common locations
func (m *Menu) findLogFilesIntelligent() []string {
	// Priority order for searching log files
	searchDirs := []string{
		"./downloads/",
		"./logs/",
		".",
	}
	
	for _, dir := range searchDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue // Directory doesn't exist, skip
		}
		
		files := m.findLogFiles(dir)
		if len(files) > 0 {
			return files
		}
	}
	
	return []string{}
}

// getSourceLocation returns the directory name where files were found
func (m *Menu) getSourceLocation(files []string) string {
	if len(files) == 0 {
		return "unknown"
	}
	
	dir := filepath.Dir(files[0])
	switch dir {
	case "./downloads", "downloads":
		return "downloads folder"
	case "./logs", "logs":
		return "logs folder"
	case ".":
		return "current directory"
	default:
		return fmt.Sprintf("%s directory", dir)
	}
}