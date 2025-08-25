package performance

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"smart-log-analyser/pkg/charts"
)

// ColorScheme provides color functions for performance visualization
type ColorScheme struct{}

// NewColorScheme creates a new color scheme
func NewColorScheme() *ColorScheme {
	return &ColorScheme{}
}

// Color functions using the charts package
func (cs *ColorScheme) Title(s string) string    { return charts.Colorize(s, charts.ColorBold+charts.ColorBlue) }
func (cs *ColorScheme) Header(s string) string   { return charts.Colorize(s, charts.ColorBold+charts.ColorCyan) }
func (cs *ColorScheme) Success(s string) string  { return charts.Colorize(s, charts.ColorGreen) }
func (cs *ColorScheme) Error(s string) string    { return charts.Colorize(s, charts.ColorRed) }
func (cs *ColorScheme) Warning(s string) string  { return charts.Colorize(s, charts.ColorYellow) }
func (cs *ColorScheme) Info(s string) string     { return charts.Colorize(s, charts.ColorBlue) }
func (cs *ColorScheme) Critical(s string) string { return charts.Colorize(s, charts.ColorBrightRed+charts.ColorBold) }
func (cs *ColorScheme) Value(s string) string    { return charts.Colorize(s, charts.ColorWhite) }
func (cs *ColorScheme) Dim(s string) string      { return charts.Colorize(s, charts.ColorDim) }

// PerformanceVisualizer handles ASCII visualization of performance data
type PerformanceVisualizer struct {
	width  int
	height int
	colors *ColorScheme
}

// NewPerformanceVisualizer creates a new performance visualizer
func NewPerformanceVisualizer() *PerformanceVisualizer {
	return &PerformanceVisualizer{
		width:  80,
		height: 20,
		colors: NewColorScheme(),
	}
}

// RenderPerformanceOverview creates a comprehensive performance overview
func (pv *PerformanceVisualizer) RenderPerformanceOverview(analysis *PerformanceAnalysis) string {
	var output strings.Builder

	// Header
	output.WriteString(pv.colors.Title("📊 PERFORMANCE ANALYSIS OVERVIEW"))
	output.WriteString("\n")
	output.WriteString(strings.Repeat("=", pv.width))
	output.WriteString("\n\n")

	// Performance Score Card
	output.WriteString(pv.renderScoreCard(analysis.Score))
	output.WriteString("\n")

	// Summary Statistics
	output.WriteString(pv.renderSummaryStats(analysis.Summary))
	output.WriteString("\n")

	// Latency Distribution Chart
	if len(analysis.EndpointMetrics) > 0 {
		output.WriteString(pv.renderLatencyDistribution(analysis.EndpointMetrics))
		output.WriteString("\n")
	}

	// Traffic Pattern Chart
	if len(analysis.TimeBasedMetrics) > 0 {
		output.WriteString(pv.renderTrafficPattern(analysis.TimeBasedMetrics))
		output.WriteString("\n")
	}

	// Top Slow Endpoints
	if len(analysis.Summary.TopSlowEndpoints) > 0 {
		output.WriteString(pv.renderTopSlowEndpoints(analysis.Summary.TopSlowEndpoints, analysis.EndpointMetrics))
		output.WriteString("\n")
	}

	// Bottlenecks Summary
	if len(analysis.Bottlenecks) > 0 {
		output.WriteString(pv.renderBottlenecksSummary(analysis.Bottlenecks))
		output.WriteString("\n")
	}

	return output.String()
}

// renderScoreCard displays the overall performance score
func (pv *PerformanceVisualizer) renderScoreCard(score PerformanceScore) string {
	var output strings.Builder

	output.WriteString(pv.colors.Header("🎯 Performance Score Card"))
	output.WriteString("\n")

	// Overall score with visual bar
	overallBar := pv.createScoreBar(score.Overall, 40)
	gradeColor := pv.getScoreColor(score.Overall)
	output.WriteString(fmt.Sprintf("Overall:     %s %s (%s)\n", 
		gradeColor(overallBar), 
		gradeColor(fmt.Sprintf("%3d", score.Overall)),
		gradeColor(GetScoreGrade(score.Overall))))

	// Individual component scores
	components := []struct {
		name  string
		value int
	}{
		{"Latency:    ", score.Latency},
		{"Throughput: ", score.Throughput},
		{"Reliability:", score.Reliability},
		{"Efficiency: ", score.Efficiency},
	}

	for _, comp := range components {
		bar := pv.createScoreBar(comp.value, 30)
		color := pv.getScoreColor(comp.value)
		output.WriteString(fmt.Sprintf("%s %s %s\n", 
			comp.name, 
			color(bar), 
			color(fmt.Sprintf("%3d", comp.value))))
	}

	return output.String()
}

// renderSummaryStats displays key summary statistics
func (pv *PerformanceVisualizer) renderSummaryStats(summary PerformanceSummary) string {
	var output strings.Builder

	output.WriteString(pv.colors.Header("📈 Summary Statistics"))
	output.WriteString("\n")

	// Format statistics
	stats := []struct {
		label string
		value string
		color func(string) string
	}{
		{"Total Requests:", fmt.Sprintf("%d", summary.TotalRequests), pv.colors.Value},
		{"Avg Response Size:", pv.formatBytes(summary.AverageResponseSize), pv.colors.Value},
		{"Error Rate:", fmt.Sprintf("%.2f%%", summary.ErrorRate*100), pv.getErrorRateColor(summary.ErrorRate)},
		{"Peak Throughput:", fmt.Sprintf("%.1f req/s", summary.PeakThroughput), pv.colors.Value},
		{"P95 Latency:", pv.formatDuration(summary.OverallLatency.P95), pv.getLatencyColor(summary.OverallLatency.P95)},
		{"Performance Grade:", summary.PerformanceGrade.String(), pv.getPerformanceGradeColor(summary.PerformanceGrade)},
	}

	for _, stat := range stats {
		output.WriteString(fmt.Sprintf("%-20s %s\n", stat.label, stat.color(stat.value)))
	}

	return output.String()
}

// renderLatencyDistribution creates a latency distribution histogram
func (pv *PerformanceVisualizer) renderLatencyDistribution(endpointMetrics map[string]*EndpointPerformance) string {
	var output strings.Builder

	output.WriteString(pv.colors.Header("⏱️  Latency Distribution (P95)"))
	output.WriteString("\n")

	// Collect P95 latencies
	var latencies []time.Duration
	for _, metrics := range endpointMetrics {
		latencies = append(latencies, metrics.EstimatedLatency.P95)
	}

	if len(latencies) == 0 {
		return output.String()
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	// Create histogram buckets
	buckets := pv.createLatencyBuckets(latencies)
	
	// Render histogram
	maxCount := 0
	for _, count := range buckets {
		if count > maxCount {
			maxCount = count
		}
	}

	if maxCount == 0 {
		return output.String()
	}

	bucketLabels := []string{"<100ms", "<500ms", "<1s", "<2s", "<5s", ">=5s"}
	
	for i, label := range bucketLabels {
		if i >= len(buckets) {
			break
		}
		
		count := buckets[i]
		barLength := int(float64(count) / float64(maxCount) * 50)
		bar := strings.Repeat("█", barLength)
		
		// Color based on latency bucket
		var colorFunc func(string) string
		switch i {
		case 0:
			colorFunc = pv.colors.Success
		case 1:
			colorFunc = pv.colors.Info
		case 2:
			colorFunc = pv.colors.Warning
		default:
			colorFunc = pv.colors.Error
		}
		
		output.WriteString(fmt.Sprintf("%-8s │%s %d endpoints\n", 
			label, 
			colorFunc(fmt.Sprintf("%-50s", bar)), 
			count))
	}

	return output.String()
}

// renderTrafficPattern creates a traffic pattern visualization
func (pv *PerformanceVisualizer) renderTrafficPattern(hourlyMetrics []HourlyPerformance) string {
	var output strings.Builder

	output.WriteString(pv.colors.Header("📊 24-Hour Traffic Pattern"))
	output.WriteString("\n")

	// Find max throughput for scaling
	maxThroughput := 0.0
	for _, metrics := range hourlyMetrics {
		if metrics.Throughput > maxThroughput {
			maxThroughput = metrics.Throughput
		}
	}

	if maxThroughput == 0 {
		return output.String()
	}

	// Render hourly bars
	for _, metrics := range hourlyMetrics {
		hour := fmt.Sprintf("%02d:00", metrics.Hour)
		barLength := int(metrics.Throughput / maxThroughput * 40)
		bar := strings.Repeat("▆", barLength)
		
		// Color based on traffic level
		var colorFunc func(string) string
		switch {
		case metrics.Throughput > maxThroughput*0.7:
			colorFunc = pv.colors.Error // Peak traffic
		case metrics.Throughput > maxThroughput*0.4:
			colorFunc = pv.colors.Warning // High traffic
		default:
			colorFunc = pv.colors.Info // Normal traffic
		}

		output.WriteString(fmt.Sprintf("%s │%s %.1f req/s\n", 
			hour, 
			colorFunc(fmt.Sprintf("%-40s", bar)), 
			metrics.Throughput))
	}

	return output.String()
}

// renderTopSlowEndpoints displays the slowest endpoints
func (pv *PerformanceVisualizer) renderTopSlowEndpoints(topSlow []string, endpointMetrics map[string]*EndpointPerformance) string {
	var output strings.Builder

	output.WriteString(pv.colors.Header("🐌 Top Slow Endpoints"))
	output.WriteString("\n")

	for i, endpoint := range topSlow {
		if i >= 5 { // Show top 5
			break
		}

		if metrics, exists := endpointMetrics[endpoint]; exists {
			latency := pv.formatDuration(metrics.EstimatedLatency.P95)
			grade := metrics.Performance.String()
			
			// Truncate long URLs
			displayURL := endpoint
			if len(displayURL) > 50 {
				displayURL = displayURL[:47] + "..."
			}

			gradeColor := pv.getPerformanceGradeColor(metrics.Performance)
			latencyColor := pv.getLatencyColor(metrics.EstimatedLatency.P95)

			output.WriteString(fmt.Sprintf("%d. %-50s %s [%s]\n", 
				i+1, 
				displayURL, 
				latencyColor(latency), 
				gradeColor(grade)))
		}
	}

	return output.String()
}

// renderBottlenecksSummary displays a summary of detected bottlenecks
func (pv *PerformanceVisualizer) renderBottlenecksSummary(bottlenecks []Bottleneck) string {
	var output strings.Builder

	output.WriteString(pv.colors.Header("⚠️  Performance Bottlenecks"))
	output.WriteString("\n")

	for i, bottleneck := range bottlenecks {
		if i >= 5 { // Show top 5
			break
		}

		severityBar := strings.Repeat("●", bottleneck.Severity)
		severityColor := pv.getSeverityColor(bottleneck.Severity)

		output.WriteString(fmt.Sprintf("%s [%s] %s\n", 
			bottleneck.Type.String(),
			severityColor(fmt.Sprintf("%-10s", severityBar)),
			bottleneck.Description))
	}

	return output.String()
}

// Helper methods

// createScoreBar creates a visual bar for scores
func (pv *PerformanceVisualizer) createScoreBar(score int, maxLength int) string {
	barLength := int(float64(score) / 100.0 * float64(maxLength))
	bar := strings.Repeat("█", barLength)
	padding := strings.Repeat("░", maxLength-barLength)
	return bar + padding
}

// createLatencyBuckets creates histogram buckets for latency distribution
func (pv *PerformanceVisualizer) createLatencyBuckets(latencies []time.Duration) []int {
	buckets := make([]int, 6) // 6 buckets

	for _, latency := range latencies {
		switch {
		case latency < 100*time.Millisecond:
			buckets[0]++
		case latency < 500*time.Millisecond:
			buckets[1]++
		case latency < 1*time.Second:
			buckets[2]++
		case latency < 2*time.Second:
			buckets[3]++
		case latency < 5*time.Second:
			buckets[4]++
		default:
			buckets[5]++
		}
	}

	return buckets
}

// Color functions

// getScoreColor returns appropriate color function based on score
func (pv *PerformanceVisualizer) getScoreColor(score int) func(string) string {
	switch {
	case score >= 80:
		return pv.colors.Success
	case score >= 60:
		return pv.colors.Warning
	default:
		return pv.colors.Error
	}
}

// getLatencyColor returns appropriate color function based on latency
func (pv *PerformanceVisualizer) getLatencyColor(latency time.Duration) func(string) string {
	switch {
	case latency < 200*time.Millisecond:
		return pv.colors.Success
	case latency < 1*time.Second:
		return pv.colors.Warning
	default:
		return pv.colors.Error
	}
}

// getErrorRateColor returns appropriate color function based on error rate
func (pv *PerformanceVisualizer) getErrorRateColor(errorRate float64) func(string) string {
	switch {
	case errorRate < 0.01:
		return pv.colors.Success
	case errorRate < 0.05:
		return pv.colors.Warning
	default:
		return pv.colors.Error
	}
}

// getPerformanceGradeColor returns appropriate color function based on performance grade
func (pv *PerformanceVisualizer) getPerformanceGradeColor(grade PerformanceGrade) func(string) string {
	switch grade {
	case Excellent:
		return pv.colors.Success
	case Good:
		return pv.colors.Info
	case Fair:
		return pv.colors.Warning
	case Poor:
		return pv.colors.Error
	case Critical:
		return pv.colors.Critical
	default:
		return pv.colors.Dim
	}
}

// getSeverityColor returns appropriate color function based on severity
func (pv *PerformanceVisualizer) getSeverityColor(severity int) func(string) string {
	switch {
	case severity >= 8:
		return pv.colors.Critical
	case severity >= 6:
		return pv.colors.Error
	case severity >= 4:
		return pv.colors.Warning
	default:
		return pv.colors.Info
	}
}

// formatBytes formats byte sizes for display
func (pv *PerformanceVisualizer) formatBytes(bytes int64) string {
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

// formatDuration formats durations for display
func (pv *PerformanceVisualizer) formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return "0ms"
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.0fμs", float64(d)/float64(time.Microsecond))
	} else if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d)/float64(time.Millisecond))
	} else {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
}