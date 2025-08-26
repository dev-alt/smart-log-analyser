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
	"smart-log-analyser/pkg/charts"
	"smart-log-analyser/pkg/config"
	"smart-log-analyser/pkg/html"
	"smart-log-analyser/pkg/parser"
	"smart-log-analyser/pkg/query"
	"smart-log-analyser/pkg/trends"
)

var (
	since         string
	until         string
	topIPs        int
	topURLs       int
	exportJSON    string
	exportCSV     string
	exportHTML    string
	htmlTitle     string
	interactiveHTML bool
	showDetails   bool
	asciiCharts   bool
	chartWidth    int
	noColors      bool
	trendAnalysis bool
	comparePeriod string
	queryString   string
	queryFormat   string
	presetName    string
	analyseConfigDir string
)

var analyseCmd = &cobra.Command{
	Use:   "analyse [log-files...]",
	Short: "Analyse Nginx access logs",
	Long:  `Parse and analyse Nginx access logs to provide statistical insights.
Accepts multiple log files to analyse together.

Interactive HTML Reports:
Generate rich, interactive HTML reports with tabbed interfaces and drill-down capabilities:

  # Interactive report (default)
  --export-html report.html --html-title "My Analysis"
  
  # Standard static report
  --export-html report.html --interactive-html=false
  
Interactive reports include:
- Tabbed interface (Overview, Traffic, Errors, Performance, Security, Geographic)
- Clickable tables with detailed drill-down information
- Real-time filtering and search capabilities
- Professional charts and visualizations
- Error analysis with status code breakdowns
- IP analysis with geographic and type information

Advanced Query Language (SLAQ):
You can use SQL-like queries to filter and analyze your logs. Examples:

  # Basic filtering
  --query "SELECT * FROM logs WHERE status = 404"
  
  # Aggregation analysis  
  --query "SELECT ip, COUNT() FROM logs GROUP BY ip ORDER BY COUNT() DESC LIMIT 10"
  
  # Time-based analysis
  --query "SELECT HOUR(timestamp), COUNT() FROM logs GROUP BY HOUR(timestamp)"
  
  # Complex filtering
  --query "SELECT url, method FROM logs WHERE status >= 400 AND url LIKE '/api*'"

Available fields: ip, timestamp, method, url, protocol, status, size, referer, user_agent
Available functions: COUNT(), SUM(), AVG(), MIN(), MAX(), HOUR(), DAY(), UPPER(), LOWER()
Available operators: =, !=, <, >, <=, >=, LIKE, CONTAINS, STARTS_WITH, ENDS_WITH, IN, BETWEEN`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Handle preset loading if specified
		if presetName != "" {
			if err := applyPreset(presetName); err != nil {
				fmt.Printf("❌ Failed to apply preset '%s': %v\n", presetName, err)
				os.Exit(1)
			}
		}
		
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

		// Execute query if provided
		if queryString != "" {
			fmt.Printf("🔍 Executing query: %s\n", queryString)
			
			// Filter logs by time if specified
			var filteredLogs []*parser.LogEntry
			if sinceTime != nil || untilTime != nil {
				a := analyser.New()
				filteredLogs = a.FilterByTime(allLogs, sinceTime, untilTime)
			} else {
				filteredLogs = allLogs
			}
			
			// Execute the query
			engine := query.NewQueryEngine(filteredLogs)
			result, err := engine.Query(queryString, queryFormat)
			if err != nil {
				fmt.Printf("❌ Query error: %v\n", err)
				helper := query.NewQueryHelper()
				fmt.Printf("💡 %s\n", helper.SuggestCorrection(err))
				return
			}
			
			fmt.Printf("📊 Query Results:\n")
			fmt.Printf("%s", result)
			return
		}

		a := analyser.New()
		results := a.Analyse(allLogs, sinceTime, untilTime)
		
		// Perform trend analysis if requested
		if trendAnalysis {
			fmt.Printf("🔍 Performing trend analysis...\n")
			ta := trends.New()
			trendResults, err := ta.DetectDegradation(allLogs)
			if err != nil {
				fmt.Printf("❌ Failed to perform trend analysis: %v\n", err)
			} else {
				fmt.Printf("📈 Trend analysis completed\n")
				printTrendAnalysis(trendResults)
				
				// Display trend charts if ASCII charts are enabled
				if asciiCharts {
					fmt.Printf("\n")
					fmt.Print(trends.RenderTrendCharts(trendResults, chartWidth, !noColors))
				}
			}
		}
		
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
		
		if exportHTML != "" {
			title := htmlTitle
			if title == "" {
				title = "Log Analysis Report"
			}
			if err := exportToHTML(results, exportHTML, title, interactiveHTML); err != nil {
				fmt.Printf("❌ Failed to export HTML: %v\n", err)
			} else {
				reportType := "standard"
				if interactiveHTML {
					reportType = "interactive"
				}
				fmt.Printf("🌐 Exported %s HTML report to: %s\n", reportType, exportHTML)
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
	analyseCmd.Flags().StringVar(&exportHTML, "export-html", "", "Export HTML report")
	analyseCmd.Flags().StringVar(&htmlTitle, "html-title", "", "Custom title for HTML report")
	analyseCmd.Flags().BoolVar(&interactiveHTML, "interactive-html", true, "Generate interactive HTML report with tabs and drill-down (default: true)")
	analyseCmd.Flags().BoolVar(&showDetails, "details", false, "Show detailed breakdown (individual status codes, etc.)")
	analyseCmd.Flags().BoolVar(&asciiCharts, "ascii-charts", false, "Display ASCII charts with analysis results")
	analyseCmd.Flags().IntVar(&chartWidth, "chart-width", 80, "Width of ASCII charts (default: 80)")
	analyseCmd.Flags().BoolVar(&noColors, "no-colors", false, "Disable colors in ASCII charts")
	analyseCmd.Flags().BoolVar(&trendAnalysis, "trend-analysis", false, "Perform historical trend analysis and degradation detection")
	analyseCmd.Flags().StringVar(&comparePeriod, "compare-period", "", "Compare with specific period (e.g., 'previous-day', '2024-08-20')")
	analyseCmd.Flags().StringVar(&queryString, "query", "", "Execute a custom SQL-like query on log data")
	analyseCmd.Flags().StringVar(&queryFormat, "query-format", "table", "Output format for query results (table, csv, json)")
	analyseCmd.Flags().StringVar(&presetName, "preset", "", "Use a predefined analysis preset (security, performance, traffic)")
	analyseCmd.Flags().StringVar(&analyseConfigDir, "config-dir", "config", "Configuration directory path")
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

	// Response Time Analysis (only show if details requested)
	if showDetails && results.ResponseTimeStats.AverageSize > 0 {
		fmt.Printf("⏱️  Response Size Analysis (Proxy for Response Time)\n")
		fmt.Printf("├─ Average Response: %s\n", formatBytes(results.ResponseTimeStats.AverageSize))
		fmt.Printf("├─ Median (P50): %s\n", formatBytes(results.ResponseTimeStats.MedianSize))
		fmt.Printf("├─ 95th Percentile: %s\n", formatBytes(results.ResponseTimeStats.P95Size))
		fmt.Printf("├─ 99th Percentile: %s\n", formatBytes(results.ResponseTimeStats.P99Size))
		fmt.Printf("├─ Range: %s - %s\n", formatBytes(results.ResponseTimeStats.MinSize), formatBytes(results.ResponseTimeStats.MaxSize))
		
		if len(results.ResponseTimeStats.SlowRequests) > 0 {
			fmt.Printf("├─ Slowest Endpoints (by size):\n")
			for i, req := range results.ResponseTimeStats.SlowRequests {
				if i >= 3 { break } // Show top 3
				displayURL := req.URL
				if len(displayURL) > 40 {
					displayURL = displayURL[:37] + "..."
				}
				fmt.Printf("│  ├─ %s: %s\n", displayURL, formatBytes(int64(req.Count)))
			}
		}
		
		if len(results.ResponseTimeStats.FastRequests) > 0 {
			fmt.Printf("└─ Fastest Endpoints (by size):\n")
			for i, req := range results.ResponseTimeStats.FastRequests {
				if i >= 3 { break } // Show top 3
				displayURL := req.URL
				if len(displayURL) > 40 {
					displayURL = displayURL[:37] + "..."
				}
				fmt.Printf("   ├─ %s: %s\n", displayURL, formatBytes(int64(req.Count)))
			}
		}
		fmt.Println()
	}

	// Geographic Analysis
	if len(results.GeographicAnalysis.TopCountries) > 0 || results.GeographicAnalysis.LocalTraffic > 0 {
		fmt.Printf("🌍 Geographic Distribution\n")
		
		// Traffic sources breakdown
		if results.GeographicAnalysis.LocalTraffic > 0 {
			localPercent := float64(results.GeographicAnalysis.LocalTraffic) / float64(results.TotalRequests) * 100
			fmt.Printf("├─ Local/Private: %s (%.1f%%)\n", formatNumber(results.GeographicAnalysis.LocalTraffic), localPercent)
		}
		if results.GeographicAnalysis.CloudTraffic > 0 {
			cloudPercent := float64(results.GeographicAnalysis.CloudTraffic) / float64(results.TotalRequests) * 100
			fmt.Printf("├─ CDN/Cloud: %s (%.1f%%)\n", formatNumber(results.GeographicAnalysis.CloudTraffic), cloudPercent)
		}
		if results.GeographicAnalysis.UnknownIPs > 0 {
			unknownPercent := float64(results.GeographicAnalysis.UnknownIPs) / float64(results.TotalRequests) * 100
			fmt.Printf("├─ Unknown IPs: %s (%.1f%%)\n", formatNumber(results.GeographicAnalysis.UnknownIPs), unknownPercent)
		}
		
		// Top countries
		if len(results.GeographicAnalysis.TopCountries) > 0 {
			fmt.Printf("├─ Countries (%d total):\n", results.GeographicAnalysis.TotalCountries)
			for i, country := range results.GeographicAnalysis.TopCountries {
				if i >= 5 { break } // Show top 5 countries
				percentage := float64(country.Count) / float64(results.TotalRequests) * 100
				fmt.Printf("│  ├─ %s: %s requests (%.1f%%)\n", country.Country, formatNumber(country.Count), percentage)
			}
		}
		
		// Top regions (only show in details mode)
		if showDetails && len(results.GeographicAnalysis.TopRegions) > 0 {
			fmt.Printf("└─ Regions:\n")
			for i, region := range results.GeographicAnalysis.TopRegions {
				if i >= 4 { break } // Show top 4 regions
				percentage := float64(region.Count) / float64(results.TotalRequests) * 100
				fmt.Printf("   ├─ %s: %s requests (%.1f%%)\n", region.Country, formatNumber(region.Count), percentage)
			}
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
	
	// Security Analysis - show when details are requested or threats detected
	if showDetails || results.SecurityAnalysis.TotalThreats > 0 {
		threatEmoji := getThreatEmoji(results.SecurityAnalysis.ThreatLevel)
		fmt.Printf("%s Security Analysis (Threat Level: %s, Score: %d/100)\n", 
			threatEmoji, 
			strings.ToUpper(results.SecurityAnalysis.ThreatLevel), 
			results.SecurityAnalysis.SecurityScore)
		
		// Overall security metrics
		fmt.Printf("├─ Total Threats Detected: %s\n", formatNumber(results.SecurityAnalysis.TotalThreats))
		fmt.Printf("├─ Suspicious IPs: %s\n", formatNumber(len(results.SecurityAnalysis.SuspiciousIPs)))
		fmt.Printf("├─ Anomalies Detected: %s\n", formatNumber(len(results.SecurityAnalysis.AnomaliesDetected)))
		
		// Attack type breakdown
		if results.SecurityAnalysis.SQLInjectionAttempts > 0 ||
		   results.SecurityAnalysis.XSSAttempts > 0 ||
		   results.SecurityAnalysis.DirectoryTraversal > 0 ||
		   results.SecurityAnalysis.BruteForceAttempts > 0 ||
		   results.SecurityAnalysis.ScanningActivity > 0 {
			fmt.Printf("├─ Attack Breakdown:\n")
			
			if results.SecurityAnalysis.SQLInjectionAttempts > 0 {
				fmt.Printf("│  ├─ SQL Injection: %s attempts\n", formatNumber(results.SecurityAnalysis.SQLInjectionAttempts))
			}
			if results.SecurityAnalysis.XSSAttempts > 0 {
				fmt.Printf("│  ├─ XSS Attempts: %s\n", formatNumber(results.SecurityAnalysis.XSSAttempts))
			}
			if results.SecurityAnalysis.DirectoryTraversal > 0 {
				fmt.Printf("│  ├─ Directory Traversal: %s attempts\n", formatNumber(results.SecurityAnalysis.DirectoryTraversal))
			}
			if results.SecurityAnalysis.BruteForceAttempts > 0 {
				fmt.Printf("│  ├─ Brute Force: %s attempts\n", formatNumber(results.SecurityAnalysis.BruteForceAttempts))
			}
			if results.SecurityAnalysis.ScanningActivity > 0 {
				fmt.Printf("│  ├─ Scanning Activity: %s instances\n", formatNumber(results.SecurityAnalysis.ScanningActivity))
			}
		}
		
		// Show top attackers
		if len(results.SecurityAnalysis.TopAttackers) > 0 {
			fmt.Printf("├─ Top Threat IPs:\n")
			for i, attacker := range results.SecurityAnalysis.TopAttackers {
				if i >= 5 { break } // Show top 5 attackers
				fmt.Printf("│  ├─ %s: %s requests", attacker.IP, formatNumber(attacker.Count))
				
				// Find IP details for threat score and categories
				for _, suspiciousIP := range results.SecurityAnalysis.SuspiciousIPs {
					if suspiciousIP.IP == attacker.IP {
						fmt.Printf(" (Score: %d", suspiciousIP.ThreatScore)
						if len(suspiciousIP.ThreatCategories) > 0 {
							fmt.Printf(", %s", strings.Join(suspiciousIP.ThreatCategories, ", "))
						}
						fmt.Printf(")")
						break
					}
				}
				fmt.Printf("\n")
			}
		}
		
		// Show recent high-severity threats in details mode
		if showDetails && len(results.SecurityAnalysis.ThreatsDetected) > 0 {
			highSeverityThreats := []analyser.SecurityThreat{}
			for _, threat := range results.SecurityAnalysis.ThreatsDetected {
				if threat.Severity == "high" || threat.Severity == "critical" {
					highSeverityThreats = append(highSeverityThreats, threat)
				}
			}
			
			if len(highSeverityThreats) > 0 {
				fmt.Printf("├─ Recent High-Severity Threats:\n")
				for i, threat := range highSeverityThreats {
					if i >= 5 { break } // Show top 5 recent threats
					threatTime := threat.Timestamp.Format("15:04:05")
					threatType := strings.ReplaceAll(threat.Type, "_", " ")
					threatType = strings.Title(threatType)
					
					fmt.Printf("│  ├─ [%s] %s from %s\n", threatTime, threatType, threat.IP)
					if len(threat.URL) > 60 {
						fmt.Printf("│  │   URL: %s...\n", threat.URL[:57])
					} else {
						fmt.Printf("│  │   URL: %s\n", threat.URL)
					}
					fmt.Printf("│  │   Pattern: %s\n", threat.Pattern)
				}
			}
		}
		
		// Show anomalies if detected
		if len(results.SecurityAnalysis.AnomaliesDetected) > 0 {
			fmt.Printf("└─ Anomalies Detected:\n")
			for i, anomaly := range results.SecurityAnalysis.AnomaliesDetected {
				if i >= 3 { break } // Show top 3 anomalies
				fmt.Printf("   ├─ %s: %.1f%% (expected %.1f%%, +%.0f%% deviation)\n", 
					strings.ReplaceAll(anomaly.Description, "_", " "),
					anomaly.Value, 
					anomaly.Expected, 
					anomaly.Deviation)
			}
		}
		
		fmt.Println()
	}

	// ASCII Charts (if enabled)
	if asciiCharts {
		fmt.Printf("📈 ASCII Charts\n")
		fmt.Printf("═══════════════\n\n")
		
		// Generate charts
		generator := charts.NewChartGenerator()
		generator.SetWidth(chartWidth)
		generator.SetColors(!noColors && charts.SupportsColor())
		
		// Display selected charts
		fmt.Print(generator.GenerateStatusCodeChart(results))
		fmt.Println()
		
		fmt.Print(generator.GenerateTopIPsChart(results, topIPs))
		fmt.Println()
		
		fmt.Print(generator.GenerateTopURLsChart(results, topURLs))
		fmt.Println()
		
		fmt.Print(generator.GenerateBotTrafficChart(results))
		fmt.Println()
		
		fmt.Print(generator.GenerateGeographicChart(results))
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

// Helper function to get emoji for threat level
func getThreatEmoji(threatLevel string) string {
	switch strings.ToLower(threatLevel) {
	case "critical":
		return "🚨"
	case "high":
		return "⚠️ "
	case "medium":
		return "🔶"
	case "low":
		return "🔐"
	default:
		return "🔐"
	}
}

// exportToHTML generates an interactive HTML report
func exportToHTML(results *analyser.Results, filename string, title string, interactive bool) error {
	generator, err := html.NewGenerator()
	if err != nil {
		return fmt.Errorf("failed to create HTML generator: %w", err)
	}
	
	if interactive {
		return generator.GenerateInteractiveReport(results, filename, title)
	}
	return generator.GenerateReport(results, filename, title)
}

// printTrendAnalysis displays trend analysis results
func printTrendAnalysis(trendAnalysis *trends.TrendAnalysis) {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                    Trend Analysis Results                      ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n\n")

	// Overall health status
	healthEmoji := getHealthEmoji(trendAnalysis.OverallHealth)
	fmt.Printf("🏥 Overall Health: %s %s\n", healthEmoji, strings.ToUpper(trendAnalysis.OverallHealth))
	fmt.Printf("📊 Analysis Type: %s\n", trendAnalysis.AnalysisType)
	fmt.Printf("🕒 Generated: %s\n", trendAnalysis.GeneratedAt.Format("2006-01-02 15:04:05"))
	
	// Trend summary
	fmt.Printf("\n📈 Trend Summary:\n")
	fmt.Printf("   %s\n", trendAnalysis.TrendSummary)

	// Period comparisons
	if len(trendAnalysis.PeriodComparisons) > 0 {
		fmt.Printf("\n📋 Period Comparison:\n")
		for _, comparison := range trendAnalysis.PeriodComparisons {
			printPeriodComparison(&comparison)
		}
	}

	// Degradation alerts
	if len(trendAnalysis.DegradationAlerts) > 0 {
		fmt.Printf("\n🚨 Degradation Alerts (%d):\n", len(trendAnalysis.DegradationAlerts))
		for _, alert := range trendAnalysis.DegradationAlerts {
			printDegradationAlert(&alert)
		}
	}

	// Recommendations
	if len(trendAnalysis.Recommendations) > 0 {
		fmt.Printf("\n💡 Recommendations:\n")
		for i, rec := range trendAnalysis.Recommendations {
			fmt.Printf("   %d. %s\n", i+1, rec)
		}
	}
}

// printPeriodComparison displays period comparison details
func printPeriodComparison(comparison *trends.PeriodComparison) {
	trendEmoji := getTrendEmoji(comparison.OverallTrend)
	fmt.Printf("├─ Overall Trend: %s %s\n", trendEmoji, comparison.OverallTrend.String())
	fmt.Printf("├─ Risk Score: %d/100\n", comparison.RiskScore)
	fmt.Printf("├─ Summary: %s\n", comparison.Summary)
	
	// Show significant changes
	fmt.Printf("└─ Key Changes:\n")
	for _, change := range comparison.TrendChanges {
		if change.Significance == "high" || change.Direction == trends.TrendCritical {
			changeEmoji := getChangeEmoji(change.Direction)
			fmt.Printf("   %s %s\n", changeEmoji, change.Description)
		}
	}
}

// printDegradationAlert displays degradation alert details
func printDegradationAlert(alert *trends.DegradationAlert) {
	severityEmoji := getSeverityEmoji(alert.Severity)
	fmt.Printf("├─ Alert %s: %s %s\n", alert.AlertID, severityEmoji, alert.MetricName)
	fmt.Printf("│  Current Value: %.2f (was %.2f)\n", alert.CurrentValue, alert.BaselineValue)
	fmt.Printf("│  Impact: %s\n", alert.Impact)
	fmt.Printf("│  Recommendation: %s\n", alert.Recommendation)
}

// Helper functions for emojis
func getHealthEmoji(health string) string {
	switch strings.ToLower(health) {
	case "healthy":
		return "✅"
	case "warning":
		return "⚠️"
	case "critical":
		return "🚨"
	default:
		return "❓"
	}
}

func getTrendEmoji(trend trends.TrendDirection) string {
	switch trend {
	case trends.TrendImproving:
		return "📈"
	case trends.TrendStable:
		return "➡️"
	case trends.TrendDegrading:
		return "📉"
	case trends.TrendCritical:
		return "🚨"
	default:
		return "❓"
	}
}

func getChangeEmoji(direction trends.TrendDirection) string {
	switch direction {
	case trends.TrendImproving:
		return "✅"
	case trends.TrendStable:
		return "➡️"
	case trends.TrendDegrading:
		return "⚠️"
	case trends.TrendCritical:
		return "🚨"
	default:
		return "❓"
	}
}

func getSeverityEmoji(severity string) string {
	switch strings.ToLower(severity) {
	case "warning":
		return "⚠️"
	case "error":
		return "❌"
	case "critical":
		return "🚨"
	default:
		return "ℹ️"
	}
}

// applyPreset loads and applies a configuration preset
func applyPreset(presetName string) error {
	// Load configuration
	configManager := config.NewConfigManager(analyseConfigDir)
	if err := configManager.Load(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get the preset
	preset, err := configManager.GetPreset(presetName)
	if err != nil {
		return fmt.Errorf("preset not found: %w", err)
	}

	fmt.Printf("🎯 Applying preset: %s (%s)\n", preset.Name, preset.Description)
	fmt.Printf("📂 Category: %s\n\n", preset.Category)

	// Apply preset query if available
	if preset.Query != "" {
		queryString = preset.Query
		fmt.Printf("🔍 Using query: %s\n\n", preset.Query)
	}

	// Apply preset filters
	if preset.Filters.Since != "" {
		since = preset.Filters.Since
	}
	if preset.Filters.Until != "" {
		until = preset.Filters.Until
	}

	// Apply preset export configurations
	for _, exportConfig := range preset.Exports {
		switch exportConfig.Format {
		case "json":
			if exportConfig.Filename != "" {
				exportJSON = exportConfig.Filename
			} else {
				exportJSON = fmt.Sprintf("output/%s.json", presetName)
			}
		case "csv":
			if exportConfig.Filename != "" {
				exportCSV = exportConfig.Filename
			} else {
				exportCSV = fmt.Sprintf("output/%s.csv", presetName)
			}
		case "html":
			if exportConfig.Filename != "" {
				exportHTML = exportConfig.Filename
			} else {
				exportHTML = fmt.Sprintf("output/%s.html", presetName)
			}
			if exportConfig.Template != "" {
				htmlTitle = exportConfig.Template
			}
		}
	}

	// Apply chart configurations
	for _, chartConfig := range preset.Charts {
		if chartConfig.Enabled {
			asciiCharts = true
			if chartConfig.Width > 0 {
				chartWidth = chartConfig.Width
			}
			noColors = !chartConfig.Colors
		}
	}

	return nil
}