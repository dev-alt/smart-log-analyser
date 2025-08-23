package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"smart-log-analyser/pkg/analyser"
	"smart-log-analyser/pkg/parser"
)

var (
	since   string
	until   string
	topIPs  int
	topURLs int
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
		
		printResults(results)
	},
}

func init() {
	analyseCmd.Flags().StringVar(&since, "since", "", "Start time (YYYY-MM-DD HH:MM:SS)")
	analyseCmd.Flags().StringVar(&until, "until", "", "End time (YYYY-MM-DD HH:MM:SS)")
	analyseCmd.Flags().IntVar(&topIPs, "top-ips", 10, "Number of top IPs to show")
	analyseCmd.Flags().IntVar(&topURLs, "top-urls", 10, "Number of top URLs to show")
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