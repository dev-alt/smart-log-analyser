package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"smart-log-analyser/pkg/parser"
	"smart-log-analyser/pkg/performance"
)

// performanceCmd represents the performance command
var performanceCmd = &cobra.Command{
	Use:   "performance <logfile>",
	Short: "Analyze performance metrics and bottlenecks",
	Long: `Analyze performance metrics and bottlenecks from log files.

This command performs comprehensive performance analysis including:
- Response time estimation and latency analysis
- Throughput and capacity assessment  
- Bottleneck detection and identification
- Performance optimization recommendations
- Visual performance reporting

Examples:
  smart-log-analyser performance access.log
  smart-log-analyser performance server-logs/*.log --export-report
  smart-log-analyser performance logs.gz --latency-threshold 500ms
  smart-log-analyser performance access.log --bottleneck-sensitivity 8`,
	Args: cobra.ExactArgs(1),
	Run:  runPerformanceAnalysis,
}

var (
	latencyThreshold      string
	bottleneckSensitivity int
	exportPerfReport      bool
	perfReportFormat      string
	perfThresholds        struct {
		excellent string
		good      string  
		fair      string
		poor      string
	}
)

func init() {
	rootCmd.AddCommand(performanceCmd)

	// Performance-specific flags
	performanceCmd.Flags().StringVar(&latencyThreshold, "latency-threshold", "1s", 
		"Custom latency alert threshold (e.g., 500ms, 2s)")
	performanceCmd.Flags().IntVar(&bottleneckSensitivity, "bottleneck-sensitivity", 7, 
		"Bottleneck detection sensitivity (1-10, higher = more sensitive)")
	performanceCmd.Flags().BoolVar(&exportPerfReport, "export-report", false, 
		"Generate detailed performance report file")
	performanceCmd.Flags().StringVar(&perfReportFormat, "report-format", "html", 
		"Report format: text, html, json")

	// Custom threshold flags
	performanceCmd.Flags().StringVar(&perfThresholds.excellent, "excellent-threshold", "100ms", 
		"Threshold for excellent performance")
	performanceCmd.Flags().StringVar(&perfThresholds.good, "good-threshold", "500ms", 
		"Threshold for good performance")  
	performanceCmd.Flags().StringVar(&perfThresholds.fair, "fair-threshold", "1s", 
		"Threshold for fair performance")
	performanceCmd.Flags().StringVar(&perfThresholds.poor, "poor-threshold", "5s", 
		"Threshold for poor performance")
}

func runPerformanceAnalysis(cmd *cobra.Command, args []string) {
	logFile := args[0]

	// Check if file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		fmt.Printf("❌ Error: File '%s' does not exist\n", logFile)
		os.Exit(1)
	}

	fmt.Printf("🔍 Analyzing performance for: %s\n\n", logFile)

	// Parse log file
	p := parser.New()
	logs, err := p.ParseFile(logFile)
	if err != nil {
		fmt.Printf("❌ Error parsing log file: %v\n", err)
		os.Exit(1)
	}

	if len(logs) == 0 {
		fmt.Printf("⚠️  No valid log entries found in %s\n", logFile)
		os.Exit(1)
	}

	fmt.Printf("📊 Parsed %d log entries\n", len(logs))

	// Create performance analyzer with custom thresholds if provided
	analyzer := createPerformanceAnalyzer()

	// Perform analysis
	fmt.Printf("🔧 Performing performance analysis...\n\n")
	analysis, err := analyzer.Analyze(logs)
	if err != nil {
		fmt.Printf("❌ Error during performance analysis: %v\n", err)
		os.Exit(1)
	}

	// Display results
	visualizer := performance.NewPerformanceVisualizer()
	overview := visualizer.RenderPerformanceOverview(analysis)
	fmt.Print(overview)

	// Display detailed recommendations
	displayRecommendations(analysis)

	// Export report if requested
	if exportPerfReport {
		exportPerformanceReport(analysis, logFile)
	}

	// Summary
	displayPerformanceSummary(analysis)
}

func createPerformanceAnalyzer() *performance.Analyzer {
	// Parse custom thresholds
	thresholds := performance.DefaultThresholds()

	if perfThresholds.excellent != "" {
		if d, err := time.ParseDuration(perfThresholds.excellent); err == nil {
			thresholds.ExcellentLatency = d
		}
	}
	if perfThresholds.good != "" {
		if d, err := time.ParseDuration(perfThresholds.good); err == nil {
			thresholds.GoodLatency = d
		}
	}
	if perfThresholds.fair != "" {
		if d, err := time.ParseDuration(perfThresholds.fair); err == nil {
			thresholds.FairLatency = d
		}
	}
	if perfThresholds.poor != "" {
		if d, err := time.ParseDuration(perfThresholds.poor); err == nil {
			thresholds.PoorLatency = d
		}
	}

	return performance.NewAnalyzerWithThresholds(thresholds)
}

func displayRecommendations(analysis *performance.PerformanceAnalysis) {
	if len(analysis.Recommendations) == 0 {
		return
	}

	fmt.Printf("\n🎯 OPTIMIZATION RECOMMENDATIONS\n")
	fmt.Printf(strings.Repeat("=", 60) + "\n\n")

	for i, rec := range analysis.Recommendations {
		if i >= 5 { // Show top 5 recommendations
			break
		}

		// Priority indicator
		priorityIndicator := strings.Repeat("★", min(rec.Priority/2, 5))
		
		// Impact and effort indicators
		impactColor := getImpactColor(rec.Impact)
		effortColor := getEffortColor(rec.Effort)

		fmt.Printf("%d. %s\n", i+1, rec.Title)
		fmt.Printf("   Priority: %s (%d/10)\n", priorityIndicator, rec.Priority)
		fmt.Printf("   Impact: %s | Effort: %s\n", 
			impactColor(rec.Impact.String()), 
			effortColor(rec.Effort.String()))
		fmt.Printf("   Category: %s\n", rec.Category.String())
		
		if rec.EstimatedImprovementPercent > 0 {
			fmt.Printf("   Estimated Improvement: %d%%\n", rec.EstimatedImprovementPercent)
		}
		
		fmt.Printf("   %s\n", rec.Description)
		
		if len(rec.Examples) > 0 {
			fmt.Printf("   Examples:\n")
			for _, example := range rec.Examples {
				fmt.Printf("   • %s\n", example)
			}
		}
		
		fmt.Printf("\n")
	}
}

func displayPerformanceSummary(analysis *performance.PerformanceAnalysis) {
	fmt.Printf("\n📋 ANALYSIS SUMMARY\n")
	fmt.Printf(strings.Repeat("=", 30) + "\n")

	fmt.Printf("Overall Score: %d/100 (%s)\n", 
		analysis.Score.Overall, 
		performance.GetScoreGrade(analysis.Score.Overall))
	fmt.Printf("Performance Grade: %s\n", analysis.Summary.PerformanceGrade.String())
	
	if len(analysis.Bottlenecks) > 0 {
		fmt.Printf("Bottlenecks Found: %d\n", len(analysis.Bottlenecks))
		fmt.Printf("Critical Issues: %d\n", analysis.Summary.CriticalIssues)
	} else {
		fmt.Printf("✅ No significant bottlenecks detected\n")
	}
	
	fmt.Printf("Recommendations: %d\n", len(analysis.Recommendations))
	fmt.Printf("Analysis Duration: %v\n", 
		analysis.LogTimeRange.End.Sub(analysis.LogTimeRange.Start).Truncate(time.Minute))
}

func exportPerformanceReport(analysis *performance.PerformanceAnalysis, logFile string) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("output/performance_report_%s.%s", timestamp, perfReportFormat)

	fmt.Printf("📄 Exporting performance report to: %s\n", filename)

	// Create output directory
	if err := os.MkdirAll("output", 0755); err != nil {
		fmt.Printf("⚠️  Warning: Could not create output directory: %v\n", err)
		return
	}

	// Export based on format
	switch perfReportFormat {
	case "html":
		exportHTMLPerformanceReport(analysis, filename, logFile)
	case "json":
		exportJSONPerformanceReport(analysis, filename)
	default:
		exportTextPerformanceReport(analysis, filename, logFile)
	}
}

func exportHTMLPerformanceReport(analysis *performance.PerformanceAnalysis, filename, logFile string) {
	// Placeholder for HTML report generation
	content := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Performance Analysis Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .score { font-size: 2em; font-weight: bold; }
        .excellent { color: green; }
        .good { color: blue; }
        .fair { color: orange; }
        .poor { color: red; }
        .critical { color: darkred; }
    </style>
</head>
<body>
    <h1>Performance Analysis Report</h1>
    <p><strong>Log File:</strong> %s</p>
    <p><strong>Analysis Time:</strong> %s</p>
    
    <h2>Performance Score</h2>
    <div class="score %s">%d/100 (%s)</div>
    
    <h2>Summary</h2>
    <ul>
        <li>Total Requests: %d</li>
        <li>Error Rate: %.2f%%</li>
        <li>Average Response Size: %s</li>
        <li>Peak Throughput: %.1f req/s</li>
    </ul>
    
    <h2>Top Recommendations</h2>
    <ol>`, 
		logFile, 
		analysis.AnalysisTimestamp.Format("2006-01-02 15:04:05"),
		strings.ToLower(performance.GetScoreDescription(analysis.Score.Overall)),
		analysis.Score.Overall,
		performance.GetScoreGrade(analysis.Score.Overall),
		analysis.Summary.TotalRequests,
		analysis.Summary.ErrorRate*100,
		formatBytesInline(analysis.Summary.AverageResponseSize),
		analysis.Summary.PeakThroughput)

	// Add recommendations
	for i, rec := range analysis.Recommendations {
		if i >= 5 {
			break
		}
		content += fmt.Sprintf(`
        <li>
            <strong>%s</strong><br>
            <em>Priority: %d/10 | Impact: %s | Effort: %s</em><br>
            %s
        </li>`, rec.Title, rec.Priority, rec.Impact.String(), rec.Effort.String(), rec.Description)
	}

	content += `
    </ol>
</body>
</html>`

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		fmt.Printf("❌ Error writing HTML report: %v\n", err)
	} else {
		fmt.Printf("✅ HTML report exported successfully\n")
	}
}

func exportJSONPerformanceReport(analysis *performance.PerformanceAnalysis, filename string) {
	// This would require proper JSON marshaling of the analysis struct
	fmt.Printf("⚠️  JSON export not yet implemented\n")
}

func exportTextPerformanceReport(analysis *performance.PerformanceAnalysis, filename, logFile string) {
	visualizer := performance.NewPerformanceVisualizer()
	content := visualizer.RenderPerformanceOverview(analysis)
	
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		fmt.Printf("❌ Error writing text report: %v\n", err)
	} else {
		fmt.Printf("✅ Text report exported successfully\n")
	}
}

// Helper functions

func getImpactColor(impact performance.ImpactLevel) func(string) string {
	// Since we can't import the colors package here, use simple return
	return func(s string) string { return s }
}

func getEffortColor(effort performance.EffortLevel) func(string) string {
	// Since we can't import the colors package here, use simple return
	return func(s string) string { return s }
}

// formatBytesInline formats byte sizes for display (inline implementation to avoid duplicate)
func formatBytesInline(bytes int64) string {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}