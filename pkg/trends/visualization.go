package trends

import (
	"fmt"
	"math"
	"strings"

	"smart-log-analyser/pkg/charts"
)

// RenderTrendCharts renders ASCII charts for trend analysis
func RenderTrendCharts(analysis *TrendAnalysis, width int, useColors bool) string {
	var output strings.Builder
	
	output.WriteString(renderSectionHeader("📈 Trend Analysis Visualizations", width))
	
	// Render period comparison charts
	if len(analysis.PeriodComparisons) > 0 {
		for i, comparison := range analysis.PeriodComparisons {
			output.WriteString(fmt.Sprintf("\n--- Period Comparison %d ---\n", i+1))
			output.WriteString(renderPeriodComparisonChart(&comparison, width, useColors))
		}
	}
	
	// Render degradation alerts summary
	if len(analysis.DegradationAlerts) > 0 {
		output.WriteString("\n--- Degradation Alerts by Severity ---\n")
		output.WriteString(renderDegradationAlertsChart(analysis.DegradationAlerts, width, useColors))
	}
	
	// Render risk score visualization
	if len(analysis.PeriodComparisons) > 0 {
		output.WriteString("\n--- Risk Score Assessment ---\n")
		output.WriteString(renderRiskScoreChart(analysis.PeriodComparisons[0].RiskScore, width, useColors))
	}
	
	return output.String()
}

// renderPeriodComparisonChart creates a chart comparing metrics between periods
func renderPeriodComparisonChart(comparison *PeriodComparison, width int, useColors bool) string {
	var output strings.Builder
	
	chart := charts.NewBarChart("Metric Changes (% Change)", width-10)
	
	// Add bars for significant changes
	for _, change := range comparison.TrendChanges {
		if change.Significance != "low" {
			color := getChangeColorCode(change.Direction, useColors)
			
			// Convert percentage to display value (cap at reasonable range)
			displayValue := change.PercentChange
			if displayValue > 100 {
				displayValue = 100
			} else if displayValue < -100 {
				displayValue = -100
			}
			
			// Make negative values positive for bar display, but keep sign in label
			barValue := math.Abs(displayValue)
			label := fmt.Sprintf("%s %+.1f%%", change.MetricName, change.PercentChange)
			
			chart.AddBar(label, int64(barValue), color)
		}
	}
	
	output.WriteString(chart.Render())
	
	// Add trend direction summary
	trendEmoji := getTrendDirectionEmoji(comparison.OverallTrend)
	output.WriteString(fmt.Sprintf("\nOverall Trend: %s %s (Risk Score: %d/100)\n", 
		trendEmoji, comparison.OverallTrend.String(), comparison.RiskScore))
	
	return output.String()
}

// renderDegradationAlertsChart shows alerts by severity
func renderDegradationAlertsChart(alerts []DegradationAlert, width int, useColors bool) string {
	// Count alerts by severity
	severityCounts := make(map[string]int)
	for _, alert := range alerts {
		severityCounts[alert.Severity]++
	}
	
	if len(severityCounts) == 0 {
		return "✅ No degradation alerts detected\n"
	}
	
	chart := charts.NewBarChart("Alerts by Severity", width-10)
	
	// Add bars for each severity level
	severityOrder := []string{"critical", "error", "warning"}
	for _, severity := range severityOrder {
		if count, exists := severityCounts[severity]; exists {
			color := getSeverityColorCode(severity, useColors)
			chart.AddBar(strings.Title(severity), int64(count), color)
		}
	}
	
	return chart.Render()
}

// renderRiskScoreChart visualizes the risk score as a gauge
func renderRiskScoreChart(riskScore int, width int, useColors bool) string {
	var output strings.Builder
	
	// Create a horizontal gauge
	gaugeWidth := width - 20
	if gaugeWidth < 20 {
		gaugeWidth = 20
	}
	
	filled := (riskScore * gaugeWidth) / 100
	empty := gaugeWidth - filled
	
	// Choose color based on risk level
	var colorCode string
	if useColors {
		if riskScore > 70 {
			colorCode = charts.ColorRed
		} else if riskScore > 30 {
			colorCode = charts.ColorYellow
		} else {
			colorCode = charts.ColorGreen
		}
	}
	
	// Build gauge
	output.WriteString(fmt.Sprintf("Risk Score: %d/100\n", riskScore))
	output.WriteString("┌" + strings.Repeat("─", gaugeWidth+2) + "┐\n")
	
	filledBar := strings.Repeat("█", filled)
	if useColors && colorCode != "" {
		filledBar = charts.Colorize(filledBar, colorCode)
	}
	emptyBar := strings.Repeat("░", empty)
	
	output.WriteString("│ " + filledBar + emptyBar + " │\n")
	output.WriteString("└" + strings.Repeat("─", gaugeWidth+2) + "┘\n")
	
	// Add risk level description
	riskLevel := getRiskLevel(riskScore)
	output.WriteString(fmt.Sprintf("Risk Level: %s\n", riskLevel))
	
	return output.String()
}

// RenderQuickTrendSummary provides a concise trend overview
func RenderQuickTrendSummary(analysis *TrendAnalysis, width int, useColors bool) string {
	var output strings.Builder
	
	output.WriteString(renderSectionHeader("📊 Quick Trend Summary", width))
	
	if len(analysis.PeriodComparisons) > 0 {
		comparison := analysis.PeriodComparisons[0]
		
		// Health status
		healthText := strings.ToUpper(analysis.OverallHealth)
		if useColors {
			healthColor := getHealthColorCode(analysis.OverallHealth)
			healthText = charts.Colorize(healthText, healthColor)
		}
		output.WriteString(fmt.Sprintf("Health Status: %s\n", healthText))
		
		// Overall trend
		trendEmoji := getTrendDirectionEmoji(comparison.OverallTrend)
		output.WriteString(fmt.Sprintf("Overall Trend: %s %s\n", trendEmoji, comparison.OverallTrend.String()))
		
		// Risk score
		riskText := fmt.Sprintf("%d/100", comparison.RiskScore)
		if useColors {
			riskColor := getRiskScoreColorCode(comparison.RiskScore)
			riskText = charts.Colorize(riskText, riskColor)
		}
		output.WriteString(fmt.Sprintf("Risk Score: %s\n", riskText))
		
		// Alert count
		alertCount := len(analysis.DegradationAlerts)
		if alertCount > 0 {
			alertText := fmt.Sprintf("%d alerts", alertCount)
			if useColors {
				alertText = charts.Colorize(alertText, charts.ColorRed)
			}
			output.WriteString(fmt.Sprintf("Alerts: %s\n", alertText))
		} else {
			output.WriteString("Alerts: ✅ None\n")
		}
	}
	
	output.WriteString(fmt.Sprintf("\n%s\n", analysis.TrendSummary))
	
	return output.String()
}

// Helper functions
func renderSectionHeader(title string, width int) string {
	if width < len(title)+4 {
		width = len(title) + 4
	}
	
	padding := (width - len(title)) / 2
	header := strings.Repeat("═", width) + "\n"
	header += strings.Repeat(" ", padding) + title + strings.Repeat(" ", width-len(title)-padding) + "\n"
	header += strings.Repeat("═", width) + "\n"
	
	return header
}

func getChangeColorCode(direction TrendDirection, useColors bool) string {
	if !useColors {
		return ""
	}
	
	switch direction {
	case TrendImproving:
		return charts.ColorGreen
	case TrendDegrading:
		return charts.ColorYellow
	case TrendCritical:
		return charts.ColorRed
	default:
		return charts.ColorBlue
	}
}

func getSeverityColorCode(severity string, useColors bool) string {
	if !useColors {
		return ""
	}
	
	switch strings.ToLower(severity) {
	case "critical":
		return charts.ColorRed
	case "error":
		return charts.ColorMagenta
	case "warning":
		return charts.ColorYellow
	default:
		return charts.ColorBlue
	}
}

func getHealthColorCode(health string) string {
	switch strings.ToLower(health) {
	case "healthy":
		return charts.ColorGreen
	case "warning":
		return charts.ColorYellow
	case "critical":
		return charts.ColorRed
	default:
		return charts.ColorBlue
	}
}

func getRiskScoreColorCode(score int) string {
	if score > 70 {
		return charts.ColorRed
	} else if score > 30 {
		return charts.ColorYellow
	} else {
		return charts.ColorGreen
	}
}

func getTrendDirectionEmoji(direction TrendDirection) string {
	switch direction {
	case TrendImproving:
		return "📈"
	case TrendStable:
		return "➡️"
	case TrendDegrading:
		return "📉"
	case TrendCritical:
		return "🚨"
	default:
		return "❓"
	}
}

func getRiskLevel(score int) string {
	if score > 70 {
		return "🚨 HIGH RISK - Immediate action required"
	} else if score > 30 {
		return "⚠️ MEDIUM RISK - Monitor closely"
	} else {
		return "✅ LOW RISK - System operating normally"
	}
}