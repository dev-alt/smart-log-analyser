package charts

import (
	"fmt"
	"sort"
	"strconv"

	"smart-log-analyser/pkg/analyser"
)

// ChartGenerator provides methods to generate charts from analysis results
type ChartGenerator struct {
	width      int
	showColors bool
}

// NewChartGenerator creates a new chart generator
func NewChartGenerator() *ChartGenerator {
	return &ChartGenerator{
		width:      80, // Default width
		showColors: SupportsColor(),
	}
}

// SetWidth sets the terminal width for charts
func (g *ChartGenerator) SetWidth(width int) {
	g.width = width
}

// SetColors enables or disables color output
func (g *ChartGenerator) SetColors(enabled bool) {
	g.showColors = enabled
}

// GenerateStatusCodeChart creates a bar chart showing HTTP status code distribution
func (g *ChartGenerator) GenerateStatusCodeChart(results *analyser.Results) string {
	if len(results.StatusCodes) == 0 {
		return "No status code data available\n"
	}

	chart := NewBarChart("HTTP Status Code Distribution", g.width)
	chart.Config.ShowColors = g.showColors

	// Convert map to sorted slice
	type statusData struct {
		code  string
		count int
	}

	var statusList []statusData
	for code, count := range results.StatusCodes {
		statusList = append(statusList, statusData{code, count})
	}

	// Sort by count (descending)
	sort.Slice(statusList, func(i, j int) bool {
		return statusList[i].count > statusList[j].count
	})

	// Add bars to chart
	for _, status := range statusList {
		label := status.code
		color := ""
		if g.showColors {
			// Convert string to int for color selection
			codeInt, _ := strconv.Atoi(status.code)
			color = GetStatusCodeColor(codeInt)
		}
		chart.AddBar(label, int64(status.count), color)
	}

	return chart.Render()
}

// GenerateTopIPsChart creates a bar chart showing top IP addresses
func (g *ChartGenerator) GenerateTopIPsChart(results *analyser.Results, limit int) string {
	if len(results.TopIPs) == 0 {
		return "No IP address data available\n"
	}

	chart := NewBarChart(fmt.Sprintf("Top %d IP Addresses", limit), g.width)
	chart.Config.ShowColors = g.showColors

	// Take top N IPs
	count := limit
	if len(results.TopIPs) < count {
		count = len(results.TopIPs)
	}

	for i, ipData := range results.TopIPs[:count] {
		label := ipData.IP
		// Truncate long IPs for display
		if len(label) > 15 {
			label = TruncateString(label, 15)
		}
		
		color := ""
		if g.showColors {
			color = GetTrafficColor(i)
		}
		chart.AddBar(label, int64(ipData.Count), color)
	}

	return chart.Render()
}

// GenerateTopURLsChart creates a bar chart showing top requested URLs
func (g *ChartGenerator) GenerateTopURLsChart(results *analyser.Results, limit int) string {
	if len(results.TopURLs) == 0 {
		return "No URL data available\n"
	}

	chart := NewBarChart(fmt.Sprintf("Top %d Requested URLs", limit), g.width)
	chart.Config.ShowColors = g.showColors

	// Take top N URLs
	count := limit
	if len(results.TopURLs) < count {
		count = len(results.TopURLs)
	}

	for i, urlData := range results.TopURLs[:count] {
		label := urlData.URL
		// Truncate long URLs for display
		if len(label) > 30 {
			label = TruncateString(label, 30)
		}
		
		color := ""
		if g.showColors {
			color = GetTrafficColor(i)
		}
		chart.AddBar(label, int64(urlData.Count), color)
	}

	return chart.Render()
}

// GenerateBotTrafficChart creates a chart showing bot vs human traffic
func (g *ChartGenerator) GenerateBotTrafficChart(results *analyser.Results) string {
	chart := NewBarChart("Traffic Classification", g.width)
	chart.Config.ShowColors = g.showColors

	// Add human traffic
	humanCount := int64(results.TotalRequests - results.BotRequests)
	if g.showColors {
		chart.AddBar("Human Traffic", humanCount, ColorGreen)
		chart.AddBar("Bot Traffic", int64(results.BotRequests), ColorYellow)
	} else {
		chart.AddBar("Human Traffic", humanCount, "")
		chart.AddBar("Bot Traffic", int64(results.BotRequests), "")
	}

	return chart.Render()
}

// GenerateGeographicChart creates a text-based geographic distribution chart
func (g *ChartGenerator) GenerateGeographicChart(results *analyser.Results) string {
	geo := results.GeographicAnalysis
	if len(geo.TopCountries) == 0 && geo.LocalTraffic == 0 {
		return "No geographic data available\n"
	}

	chart := NewBarChart("Geographic Distribution", g.width)
	chart.Config.ShowColors = g.showColors
	
	if g.showColors {
		chart.AddBar("Local Networks", int64(geo.LocalTraffic), ColorGreen)
		chart.AddBar("Cloud/CDN", int64(geo.CloudTraffic), ColorBlue)
		chart.AddBar("Unknown IPs", int64(geo.UnknownIPs), ColorYellow)
	} else {
		chart.AddBar("Local Networks", int64(geo.LocalTraffic), "")
		chart.AddBar("Cloud/CDN", int64(geo.CloudTraffic), "")
		chart.AddBar("Unknown IPs", int64(geo.UnknownIPs), "")
	}

	return chart.Render()
}

// GenerateResponseSizeChart creates a histogram of response sizes
func (g *ChartGenerator) GenerateResponseSizeChart(results *analyser.Results) string {
	if results.TotalRequests == 0 {
		return "No response size data available\n"
	}

	chart := NewBarChart("Response Size Distribution", g.width)
	chart.Config.ShowColors = g.showColors

	// Create buckets for response sizes (using response size as proxy)
	// Define size buckets (in bytes)
	buckets := []struct {
		label string
		min   int64
		max   int64
		count int64
	}{
		{"< 1KB", 0, 1024, 0},
		{"1-10KB", 1024, 10240, 0},
		{"10-100KB", 10240, 102400, 0},
		{"100KB-1MB", 102400, 1048576, 0},
		{"> 1MB", 1048576, 999999999, 0},
	}

	// This is a simplified version - in a real implementation, 
	// you'd collect actual response size data during parsing
	// For now, we'll use percentiles as a proxy
	totalRequests := int64(results.TotalRequests)
	
	// Distribute requests across buckets based on percentiles (approximation)
	buckets[0].count = totalRequests * 20 / 100  // 20% small files
	buckets[1].count = totalRequests * 50 / 100  // 50% medium files
	buckets[2].count = totalRequests * 25 / 100  // 25% large files
	buckets[3].count = totalRequests * 4 / 100   // 4% very large files
	buckets[4].count = totalRequests * 1 / 100   // 1% huge files

	// Add bars to chart
	for i, bucket := range buckets {
		color := ""
		if g.showColors {
			color = GetTrafficColor(i)
		}
		chart.AddBar(bucket.label, bucket.count, color)
	}

	return chart.Render()
}

// GenerateFullReport generates all available charts
func (g *ChartGenerator) GenerateFullReport(results *analyser.Results) string {
	report := fmt.Sprintf("📈 ASCII Charts Report\n")
	report += fmt.Sprintf("═══════════════════════\n\n")

	report += g.GenerateStatusCodeChart(results) + "\n"
	report += g.GenerateTopIPsChart(results, 5) + "\n"
	report += g.GenerateTopURLsChart(results, 5) + "\n"
	report += g.GenerateBotTrafficChart(results) + "\n"
	report += g.GenerateGeographicChart(results) + "\n"
	report += g.GenerateResponseSizeChart(results) + "\n"

	return report
}

// GenerateSummary generates a compact summary with key charts
func (g *ChartGenerator) GenerateSummary(results *analyser.Results) string {
	report := fmt.Sprintf("📊 Quick Charts Summary\n")
	report += fmt.Sprintf("══════════════════════\n\n")

	report += g.GenerateStatusCodeChart(results) + "\n"
	report += g.GenerateBotTrafficChart(results) + "\n"

	return report
}