package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"smart-log-analyser/pkg/analyser"
	"smart-log-analyser/pkg/parser"
)

var (
	since      string
	until      string
	topIPs     int
	topURLs    int
	exportJSON string
	exportCSV  string
	showDetails bool
)

var analyseCmd = &cobra.Command{
	Use:   "analyse [log-files...]",
	Short: "Analyse Nginx access logs",
	Long:  `Parse and analyse Nginx access logs to provide statistical insights.
Accepts multiple log files to analyse together.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		p := parser.New()
		var allLogs []*parser.LogEntry
		
		fmt.Printf("📂 Analysing %d log file(s)...\n\n", len(args))
		
		for i, logFile := range args {
			fmt.Printf("  [%d/%d] Processing: %s\n", i+1, len(args), logFile)
			
			logs, err := p.ParseFile(logFile)
			if err != nil {
				fmt.Printf("    ❌ Failed to parse %s: %v\n", logFile, err)
				continue
			}
			
			fmt.Printf("    ✅ Parsed %d entries\n", len(logs))
			allLogs = append(allLogs, logs...)
		}
		
		if len(allLogs) == 0 {
			log.Fatal("No valid log entries found in any files")
		}
		
		fmt.Printf("\n📊 Combined Analysis Results (%d total entries):\n", len(allLogs))

		var sinceTime, untilTime *time.Time
		if since != "" {
			t, err := time.Parse("2006-01-02 15:04:05", since)
			if err != nil {
				log.Fatalf("Invalid since time format: %v", err)
			}
			sinceTime = &t
		}
		if until != "" {
			t, err := time.Parse("2006-01-02 15:04:05", until)
			if err != nil {
				log.Fatalf("Invalid until time format: %v", err)
			}
			untilTime = &t
		}

		a := analyser.New()
		results := a.Analyse(allLogs, sinceTime, untilTime)
		
		// Export to files if requested
		if exportJSON != "" {
			if err := exportToJSON(results, exportJSON); err != nil {
				fmt.Printf("❌ Failed to export JSON: %v\n", err)
			} else {
				fmt.Printf("📄 Exported detailed results to: %s\n", exportJSON)
			}
		}
		
		if exportCSV != "" {
			if err := exportToCSV(results, exportCSV); err != nil {
				fmt.Printf("❌ Failed to export CSV: %v\n", err)
			} else {
				fmt.Printf("📊 Exported detailed results to: %s\n", exportCSV)
			}
		}
		
		printResults(results)
	},
}

func init() {
	analyseCmd.Flags().StringVar(&since, "since", "", "Start time (YYYY-MM-DD HH:MM:SS)")
	analyseCmd.Flags().StringVar(&until, "until", "", "End time (YYYY-MM-DD HH:MM:SS)")
	analyseCmd.Flags().IntVar(&topIPs, "top-ips", 10, "Number of top IPs to show")
	analyseCmd.Flags().IntVar(&topURLs, "top-urls", 10, "Number of top URLs to show")
	analyseCmd.Flags().StringVar(&exportJSON, "export-json", "", "Export detailed results to JSON file")
	analyseCmd.Flags().StringVar(&exportCSV, "export-csv", "", "Export detailed results to CSV file")
	analyseCmd.Flags().BoolVar(&showDetails, "details", false, "Show detailed breakdown (individual status codes, etc.)")
}

func printResults(results *analyser.Results) {
	fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                   Smart Log Analyser Results                  ║\n") 
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n\n")
	
	// Overview Statistics
	fmt.Printf("📊 Overview\n")
	fmt.Printf("├─ Total Requests: %s\n", formatNumber(results.TotalRequests))
	fmt.Printf("├─ Unique IPs: %s\n", formatNumber(results.UniqueIPs))
	fmt.Printf("├─ Unique URLs: %s\n", formatNumber(results.UniqueURLs))
	fmt.Printf("├─ Data Transferred: %s\n", formatBytes(results.TotalBytes))
	fmt.Printf("├─ Average Response Size: %s\n", formatBytes(results.AverageSize))
	fmt.Printf("└─ Date Range: %s to %s\n\n", 
		results.TimeRange.Start.Format("2006-01-02 15:04:05"),
		results.TimeRange.End.Format("2006-01-02 15:04:05"))

	// Traffic Analysis (Bot vs Human)
	if results.BotRequests > 0 || results.HumanRequests > 0 {
		fmt.Printf("🤖 Traffic Analysis\n")
		botPercentage := float64(results.BotRequests) / float64(results.TotalRequests) * 100
		humanPercentage := float64(results.HumanRequests) / float64(results.TotalRequests) * 100
		fmt.Printf("├─ Human Traffic: %s (%.1f%%)\n", formatNumber(results.HumanRequests), humanPercentage)
		fmt.Printf("├─ Bot/Automated: %s (%.1f%%)\n", formatNumber(results.BotRequests), botPercentage)
		fmt.Println()
	}

	// Top Bots
	if len(results.TopBots) > 0 {
		fmt.Printf("🔍 Top Bots/Crawlers\n")
		count := 0
		for _, bot := range results.TopBots {
			if count >= 5 { // Show top 5 bots
				break
			}
			percentage := float64(bot.Count) / float64(results.TotalRequests) * 100
			fmt.Printf("├─ %s: %s requests (%.1f%%)\n", bot.BotName, formatNumber(bot.Count), percentage)
			count++
		}
		fmt.Println()
	}

	// File Types
	if len(results.FileTypes) > 0 {
		fmt.Printf("📁 File Type Analysis\n")
		count := 0
		for _, fileType := range results.FileTypes {
			if count >= 8 { // Show top 8 file types
				break
			}
			percentage := float64(fileType.Count) / float64(results.TotalRequests) * 100
			avgSize := fileType.Size / int64(fileType.Count)
			fmt.Printf("├─ %s: %s requests (%.1f%%) - %s total, %s avg\n", 
				fileType.FileType, formatNumber(fileType.Count), percentage, 
				formatBytes(fileType.Size), formatBytes(avgSize))
			count++
		}
		fmt.Println()
	}

	// Traffic Pattern Analysis
	if len(results.HourlyTraffic) > 0 {
		fmt.Printf("📈 Traffic Patterns\n")
		fmt.Printf("├─ Average Requests/Hour: %.1f\n", results.AverageRequestsPerHour)
		if results.PeakHour >= 0 {
			fmt.Printf("├─ Peak Hour: %02d:00 (%s)\n", results.PeakHour, getHourName(results.PeakHour))
		}
		if results.QuietestHour >= 0 {
			fmt.Printf("├─ Quietest Hour: %02d:00 (%s)\n", results.QuietestHour, getHourName(results.QuietestHour))
		}
		
		// Show hourly breakdown
		fmt.Printf("└─ Hourly Breakdown:\n")
		for _, hour := range results.HourlyTraffic {
			percentage := float64(hour.RequestCount) / float64(results.TotalRequests) * 100
			bar := createSimpleBar(percentage, 20)
			fmt.Printf("   ├─ %02d:00: %s requests (%.1f%%) %s\n", 
				hour.Hour, formatNumber(hour.RequestCount), percentage, bar)
		}
		fmt.Println()
	}

	// Traffic Peaks (only show if there are peaks and details requested)
	if showDetails && len(results.TrafficPeaks) > 0 {
		fmt.Printf("🔥 Traffic Peaks Detected\n")
		for i, peak := range results.TrafficPeaks {
			fmt.Printf("├─ Peak #%d: %s - %s requests (%s)\n", 
				i+1, peak.Time, formatNumber(peak.RequestCount), peak.Duration)
		}
		fmt.Println()
	}

	// HTTP Methods
	if len(results.HTTPMethods) > 0 {
		fmt.Printf("🔧 HTTP Methods\n")
		for _, method := range results.HTTPMethods {
			percentage := float64(method.Count) / float64(results.TotalRequests) * 100
			fmt.Printf("├─ %s: %s (%.1f%%)\n", method.Method, formatNumber(method.Count), percentage)
		}
		fmt.Println()
	}

	// Status Code Distribution
	fmt.Printf("📈 Status Code Distribution\n")
	statusOrder := []string{"2xx Success", "3xx Redirect", "4xx Client Error", "5xx Server Error", "1xx Informational"}
	for _, status := range statusOrder {
		if count, exists := results.StatusCodes[status]; exists {
			percentage := float64(count) / float64(results.TotalRequests) * 100
			fmt.Printf("├─ %s: %s (%.1f%%)\n", status, formatNumber(count), percentage)
		}
	}
	
	// Show detailed status codes if requested
	if showDetails && len(results.DetailedStatusCodes) > 0 {
		fmt.Printf("└─ Detailed Status Codes:\n")
		for i, status := range results.DetailedStatusCodes {
			if i >= 10 { break } // Show top 10 detailed codes
			percentage := float64(status.Count) / float64(results.TotalRequests) * 100
			fmt.Printf("   ├─ %d: %s requests (%.1f%%)\n", status.Code, formatNumber(status.Count), percentage)
		}
	}
	fmt.Println()

	// Top IPs
	fmt.Printf("🌐 Top %d IP Addresses\n", topIPs)
	count := 0
	for _, ip := range results.TopIPs {
		if count >= topIPs {
			break
		}
		percentage := float64(ip.Count) / float64(results.TotalRequests) * 100
		fmt.Printf("├─ %s: %s requests (%.1f%%)\n", ip.IP, formatNumber(ip.Count), percentage)
		count++
	}
	fmt.Println()

	// Top URLs
	fmt.Printf("🔗 Top %d URLs\n", topURLs)
	count = 0
	for _, url := range results.TopURLs {
		if count >= topURLs {
			break
		}
		percentage := float64(url.Count) / float64(results.TotalRequests) * 100
		// Truncate long URLs for display
		displayURL := url.URL
		if len(displayURL) > 60 {
			displayURL = displayURL[:57] + "..."
		}
		fmt.Printf("├─ %s: %s requests (%.1f%%)\n", displayURL, formatNumber(url.Count), percentage)
		count++
	}
	fmt.Println()
	
	// Error Analysis (only show if there are errors and details are requested)
	if showDetails && len(results.ErrorURLs) > 0 {
		fmt.Printf("⚠️  Error Analysis\n")
		fmt.Printf("├─ URLs with Errors (4xx/5xx):\n")
		for i, url := range results.ErrorURLs {
			if i >= 5 { break } // Show top 5 error URLs
			displayURL := url.URL
			if len(displayURL) > 50 {
				displayURL = displayURL[:47] + "..."
			}
			fmt.Printf("   ├─ %s: %d errors\n", displayURL, url.Count)
		}
		fmt.Println()
	}
	
	// Large Requests Analysis (only show if details are requested)
	if showDetails && len(results.LargeRequests) > 0 {
		fmt.Printf("📦 Largest Requests by Size\n")
		for i, url := range results.LargeRequests {
			if i >= 5 { break } // Show top 5 largest requests
			displayURL := url.URL
			if len(displayURL) > 50 {
				displayURL = displayURL[:47] + "..."
			}
			fmt.Printf("├─ %s: %s\n", displayURL, formatBytes(int64(url.Count))) // Count field contains size
		}
		fmt.Println()
	}
}

// Helper function to format numbers with commas
func formatNumber(num int) string {
	str := fmt.Sprintf("%d", num)
	if len(str) <= 3 {
		return str
	}
	
	result := ""
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(char)
	}
	return result
}

// Helper function to format bytes in human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Helper function to get hour name
func getHourName(hour int) string {
	switch {
	case hour >= 6 && hour < 12:
		return "Morning"
	case hour >= 12 && hour < 18:
		return "Afternoon"
	case hour >= 18 && hour < 22:
		return "Evening"
	default:
		return "Night"
	}
}

// Helper function to create a simple text-based bar chart
func createSimpleBar(percentage float64, maxWidth int) string {
	if percentage <= 0 {
		return ""
	}
	
	width := int(percentage / 100.0 * float64(maxWidth))
	if width == 0 && percentage > 0 {
		width = 1 // Ensure at least one character for non-zero values
	}
	
	bar := strings.Repeat("█", width)
	remaining := maxWidth - width
	if remaining > 0 {
		bar += strings.Repeat("░", remaining)
	}
	
	return fmt.Sprintf("[%s]", bar)
}

func exportToJSON(results *analyser.Results, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func exportToCSV(results *analyser.Results, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Write overview section
	writer.Write([]string{"Section", "Metric", "Value", "Percentage"})
	writer.Write([]string{"Overview", "Total Requests", strconv.Itoa(results.TotalRequests), "100.0"})
	writer.Write([]string{"Overview", "Unique IPs", strconv.Itoa(results.UniqueIPs), ""})
	writer.Write([]string{"Overview", "Unique URLs", strconv.Itoa(results.UniqueURLs), ""})
	writer.Write([]string{"Overview", "Total Bytes", strconv.FormatInt(results.TotalBytes, 10), ""})
	writer.Write([]string{"Overview", "Average Size", strconv.FormatInt(results.AverageSize, 10), ""})
	writer.Write([]string{"Overview", "Human Requests", strconv.Itoa(results.HumanRequests), fmt.Sprintf("%.1f", float64(results.HumanRequests)/float64(results.TotalRequests)*100)})
	writer.Write([]string{"Overview", "Bot Requests", strconv.Itoa(results.BotRequests), fmt.Sprintf("%.1f", float64(results.BotRequests)/float64(results.TotalRequests)*100)})
	
	// Write status codes
	for status, count := range results.StatusCodes {
		percentage := float64(count) / float64(results.TotalRequests) * 100
		writer.Write([]string{"Status Codes", status, strconv.Itoa(count), fmt.Sprintf("%.1f", percentage)})
	}
	
	// Write detailed status codes
	for _, status := range results.DetailedStatusCodes {
		percentage := float64(status.Count) / float64(results.TotalRequests) * 100
		writer.Write([]string{"Detailed Status", strconv.Itoa(status.Code), strconv.Itoa(status.Count), fmt.Sprintf("%.1f", percentage)})
	}
	
	// Write top IPs
	for i, ip := range results.TopIPs {
		if i >= 20 { break } // Limit to top 20 for CSV
		percentage := float64(ip.Count) / float64(results.TotalRequests) * 100
		writer.Write([]string{"Top IPs", ip.IP, strconv.Itoa(ip.Count), fmt.Sprintf("%.1f", percentage)})
	}
	
	// Write top URLs
	for i, url := range results.TopURLs {
		if i >= 20 { break } // Limit to top 20 for CSV
		percentage := float64(url.Count) / float64(results.TotalRequests) * 100
		writer.Write([]string{"Top URLs", url.URL, strconv.Itoa(url.Count), fmt.Sprintf("%.1f", percentage)})
	}
	
	// Write top bots
	for _, bot := range results.TopBots {
		percentage := float64(bot.Count) / float64(results.TotalRequests) * 100
		writer.Write([]string{"Top Bots", bot.BotName, strconv.Itoa(bot.Count), fmt.Sprintf("%.1f", percentage)})
	}
	
	// Write file types
	for _, ft := range results.FileTypes {
		percentage := float64(ft.Count) / float64(results.TotalRequests) * 100
		avgSize := ft.Size / int64(ft.Count)
		writer.Write([]string{"File Types", ft.FileType, strconv.Itoa(ft.Count), fmt.Sprintf("%.1f", percentage)})
		writer.Write([]string{"File Types Size", ft.FileType + " Total", strconv.FormatInt(ft.Size, 10), ""})
		writer.Write([]string{"File Types Size", ft.FileType + " Average", strconv.FormatInt(avgSize, 10), ""})
	}
	
	// Write error URLs
	for _, url := range results.ErrorURLs {
		writer.Write([]string{"Error URLs", url.URL, strconv.Itoa(url.Count), ""})
	}
	
	// Write large requests
	for _, url := range results.LargeRequests {
		writer.Write([]string{"Large Requests", url.URL, strconv.Itoa(url.Count), ""}) // Count field contains size
	}
	
	return nil
}