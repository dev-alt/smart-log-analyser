package trends

import (
	"fmt"
	"math"
	"sort"
	"time"

	"smart-log-analyser/pkg/analyser"
	"smart-log-analyser/pkg/parser"
)

// TrendAnalyser performs historical trend analysis on log data
type TrendAnalyser struct {
	config TrendConfiguration
}

// New creates a new TrendAnalyser with default configuration
func New() *TrendAnalyser {
	return &TrendAnalyser{
		config: DefaultTrendConfiguration(),
	}
}

// NewWithConfig creates a new TrendAnalyser with custom configuration
func NewWithConfig(config TrendConfiguration) *TrendAnalyser {
	return &TrendAnalyser{
		config: config,
	}
}

// ComparePeriods compares two time periods and identifies trends
func (ta *TrendAnalyser) ComparePeriods(baselineLogs, currentLogs []*parser.LogEntry) (*PeriodComparison, error) {
	// Analyze both periods
	analyzer := analyser.New()
	
	baselineResults := analyzer.Analyse(baselineLogs, nil, nil)
	currentResults := analyzer.Analyse(currentLogs, nil, nil)
	
	// Convert to PeriodMetrics
	baselineMetrics := ta.convertToPeriodMetrics("Baseline Period", baselineResults)
	currentMetrics := ta.convertToPeriodMetrics("Current Period", currentResults)
	
	// Calculate trend changes
	trendChanges := ta.calculateTrendChanges(baselineMetrics, currentMetrics)
	
	// Determine overall trend
	overallTrend := ta.calculateOverallTrend(trendChanges)
	
	// Calculate risk score
	riskScore := ta.calculateRiskScore(trendChanges)
	
	// Generate summary
	summary := ta.generateComparisonSummary(overallTrend, riskScore, trendChanges)
	
	return &PeriodComparison{
		BaselinePeriod: baselineMetrics,
		CurrentPeriod:  currentMetrics,
		TrendChanges:   trendChanges,
		OverallTrend:   overallTrend,
		RiskScore:      riskScore,
		Summary:        summary,
	}, nil
}

// DetectDegradation analyzes logs for performance degradation patterns
func (ta *TrendAnalyser) DetectDegradation(logs []*parser.LogEntry) (*TrendAnalysis, error) {
	// For degradation detection, we'll split the logs into periods and compare them
	if len(logs) < ta.config.MinimumSampleSize {
		return nil, fmt.Errorf("insufficient data: need at least %d log entries", ta.config.MinimumSampleSize)
	}
	
	// Split logs into two halves for before/after comparison
	midPoint := len(logs) / 2
	earlierLogs := logs[:midPoint]
	laterLogs := logs[midPoint:]
	
	comparison, err := ta.ComparePeriods(earlierLogs, laterLogs)
	if err != nil {
		return nil, err
	}
	
	// Detect specific degradation alerts
	alerts := ta.generateDegradationAlerts(comparison.TrendChanges)
	
	// Determine overall health
	overallHealth := ta.calculateOverallHealth(alerts, comparison.RiskScore)
	
	// Generate recommendations
	recommendations := ta.generateRecommendations(alerts, comparison.TrendChanges)
	
	// Create trend summary
	trendSummary := ta.generateTrendSummary(comparison, alerts)
	
	return &TrendAnalysis{
		AnalysisType:      "degradation",
		GeneratedAt:       time.Now(),
		PeriodComparisons: []PeriodComparison{*comparison},
		DegradationAlerts: alerts,
		OverallHealth:     overallHealth,
		Recommendations:   recommendations,
		TrendSummary:      trendSummary,
	}, nil
}

// convertToPeriodMetrics converts analyser.Results to PeriodMetrics
func (ta *TrendAnalyser) convertToPeriodMetrics(periodName string, results *analyser.Results) PeriodMetrics {
	// Calculate error rate
	totalErrors := results.StatusCodes["4xx"] + results.StatusCodes["5xx"]
	errorRate := 0.0
	if results.TotalRequests > 0 {
		errorRate = (float64(totalErrors) / float64(results.TotalRequests)) * 100
	}
	
	// Calculate bot traffic percentage
	botTrafficPercent := 0.0
	if results.TotalRequests > 0 {
		botTrafficPercent = (float64(results.BotRequests) / float64(results.TotalRequests)) * 100
	}
	
	// Get peak hour requests
	peakHourRequests := 0
	if len(results.HourlyTraffic) > 0 {
		for _, hourly := range results.HourlyTraffic {
			if hourly.RequestCount > peakHourRequests {
				peakHourRequests = hourly.RequestCount
			}
		}
	}
	
	// Create geographic distribution map
	geoDistrib := make(map[string]int)
	for _, geo := range results.GeographicAnalysis.TopCountries {
		geoDistrib[geo.Country] = geo.Count
	}
	
	return PeriodMetrics{
		Period:               periodName,
		StartTime:            results.TimeRange.Start,
		EndTime:              results.TimeRange.End,
		TotalRequests:        results.TotalRequests,
		AverageResponseSize:  results.AverageSize,
		ErrorRate:            errorRate,
		TrafficVolume:        results.TotalBytes,
		UniqueVisitors:       results.UniqueIPs,
		PeakHourRequests:     peakHourRequests,
		StatusCodeDistrib:    results.StatusCodes,
		TopErrorURLs:         results.ErrorURLs,
		BotTrafficPercent:    botTrafficPercent,
		GeographicDistrib:    geoDistrib,
	}
}

// calculateTrendChanges compares metrics between two periods
func (ta *TrendAnalyser) calculateTrendChanges(baseline, current PeriodMetrics) []TrendChange {
	var changes []TrendChange
	
	// Request volume change
	changes = append(changes, ta.calculateMetricChange(
		"Request Volume", 
		float64(baseline.TotalRequests), 
		float64(current.TotalRequests),
		"requests",
	))
	
	// Error rate change
	changes = append(changes, ta.calculateMetricChange(
		"Error Rate", 
		baseline.ErrorRate, 
		current.ErrorRate,
		"%",
	))
	
	// Average response size change (proxy for performance)
	changes = append(changes, ta.calculateMetricChange(
		"Average Response Size", 
		float64(baseline.AverageResponseSize), 
		float64(current.AverageResponseSize),
		"bytes",
	))
	
	// Traffic volume change
	changes = append(changes, ta.calculateMetricChange(
		"Traffic Volume", 
		float64(baseline.TrafficVolume), 
		float64(current.TrafficVolume),
		"bytes",
	))
	
	// Unique visitors change
	changes = append(changes, ta.calculateMetricChange(
		"Unique Visitors", 
		float64(baseline.UniqueVisitors), 
		float64(current.UniqueVisitors),
		"visitors",
	))
	
	// Bot traffic percentage change
	changes = append(changes, ta.calculateMetricChange(
		"Bot Traffic", 
		baseline.BotTrafficPercent, 
		current.BotTrafficPercent,
		"%",
	))
	
	return changes
}

// calculateMetricChange calculates the change between two metric values
func (ta *TrendAnalyser) calculateMetricChange(metricName string, oldValue, newValue float64, unit string) TrendChange {
	absoluteChange := newValue - oldValue
	percentChange := 0.0
	
	if oldValue != 0 {
		percentChange = (absoluteChange / oldValue) * 100
	}
	
	// Determine trend direction and significance
	direction := TrendStable
	significance := "low"
	
	absPercentChange := math.Abs(percentChange)
	
	// Determine direction based on metric type and change
	if metricName == "Error Rate" || metricName == "Average Response Size" {
		// Higher is worse for these metrics
		if percentChange > 5 {
			direction = TrendDegrading
			if percentChange > ta.config.ErrorRateThreshold {
				direction = TrendCritical
			}
		} else if percentChange < -5 {
			direction = TrendImproving
		}
	} else {
		// Higher is better for volume/traffic metrics
		if percentChange > 5 {
			direction = TrendImproving
		} else if percentChange < -5 {
			direction = TrendDegrading
			if absPercentChange > ta.config.TrafficDropThreshold {
				direction = TrendCritical
			}
		}
	}
	
	// Determine significance
	if absPercentChange > 50 {
		significance = "high"
	} else if absPercentChange > 15 {
		significance = "medium"
	}
	
	// Generate description
	description := ta.generateMetricDescription(metricName, percentChange, direction, unit)
	
	return TrendChange{
		MetricName:     metricName,
		OldValue:       oldValue,
		NewValue:       newValue,
		AbsoluteChange: absoluteChange,
		PercentChange:  percentChange,
		Direction:      direction,
		Significance:   significance,
		Description:    description,
	}
}

// generateMetricDescription creates a human-readable description of the metric change
func (ta *TrendAnalyser) generateMetricDescription(metricName string, percentChange float64, direction TrendDirection, unit string) string {
	absChange := math.Abs(percentChange)
	
	var changeVerb string
	if percentChange > 0 {
		changeVerb = "increased"
	} else {
		changeVerb = "decreased"
	}
	
	var severity string
	if absChange > 50 {
		severity = "significantly"
	} else if absChange > 15 {
		severity = "moderately"
	} else {
		severity = "slightly"
	}
	
	return fmt.Sprintf("%s %s %s by %.1f%%", metricName, severity, changeVerb, absChange)
}

// calculateOverallTrend determines the overall trend direction
func (ta *TrendAnalyser) calculateOverallTrend(changes []TrendChange) TrendDirection {
	var scores []int
	
	for _, change := range changes {
		switch change.Direction {
		case TrendCritical:
			scores = append(scores, -3)
		case TrendDegrading:
			scores = append(scores, -1)
		case TrendImproving:
			scores = append(scores, 1)
		case TrendStable:
			scores = append(scores, 0)
		}
		
		// Weight by significance
		if change.Significance == "high" {
			scores[len(scores)-1] *= 3
		} else if change.Significance == "medium" {
			scores[len(scores)-1] *= 2
		}
	}
	
	// Calculate weighted average
	totalScore := 0
	for _, score := range scores {
		totalScore += score
	}
	
	if totalScore <= -3 {
		return TrendCritical
	} else if totalScore < 0 {
		return TrendDegrading
	} else if totalScore > 0 {
		return TrendImproving
	}
	return TrendStable
}

// calculateRiskScore calculates an overall risk score (0-100)
func (ta *TrendAnalyser) calculateRiskScore(changes []TrendChange) int {
	riskScore := 0
	
	for _, change := range changes {
		if change.Direction == TrendDegrading || change.Direction == TrendCritical {
			risk := int(math.Abs(change.PercentChange))
			
			// Apply multipliers for critical metrics
			if change.MetricName == "Error Rate" {
				risk *= 2
			}
			
			// Apply significance multipliers
			if change.Significance == "high" {
				risk *= 2
			} else if change.Significance == "medium" {
				risk = risk * 3 / 2
			}
			
			riskScore += risk
		}
	}
	
	// Cap at 100
	if riskScore > 100 {
		riskScore = 100
	}
	
	return riskScore
}

// generateComparisonSummary creates a human-readable summary
func (ta *TrendAnalyser) generateComparisonSummary(trend TrendDirection, riskScore int, changes []TrendChange) string {
	var summary string
	
	switch trend {
	case TrendImproving:
		summary = "✅ Performance is improving compared to baseline period. "
	case TrendStable:
		summary = "📊 Performance remains stable with no significant changes. "
	case TrendDegrading:
		summary = "⚠️ Performance degradation detected. "
	case TrendCritical:
		summary = "🚨 Critical performance degradation requires immediate attention. "
	}
	
	// Add risk context
	if riskScore > 70 {
		summary += "High risk to system performance."
	} else if riskScore > 30 {
		summary += "Moderate risk to system performance."
	} else {
		summary += "Low risk to system performance."
	}
	
	// Mention most significant changes
	significantChanges := 0
	for _, change := range changes {
		if change.Significance == "high" && (change.Direction == TrendDegrading || change.Direction == TrendCritical) {
			significantChanges++
		}
	}
	
	if significantChanges > 0 {
		summary += fmt.Sprintf(" %d metrics show significant degradation.", significantChanges)
	}
	
	return summary
}

// generateDegradationAlerts creates alerts for detected degradation
func (ta *TrendAnalyser) generateDegradationAlerts(changes []TrendChange) []DegradationAlert {
	var alerts []DegradationAlert
	alertID := 1
	
	for _, change := range changes {
		if change.Direction == TrendDegrading || change.Direction == TrendCritical {
			severity := "warning"
			if change.Direction == TrendCritical || change.Significance == "high" {
				severity = "critical"
			}
			
			alert := DegradationAlert{
				AlertID:       fmt.Sprintf("TREND-%03d", alertID),
				Severity:      severity,
				MetricName:    change.MetricName,
				CurrentValue:  change.NewValue,
				BaselineValue: change.OldValue,
				Threshold:     ta.getThresholdForMetric(change.MetricName),
				Impact:        ta.getImpactDescription(change.MetricName, change.PercentChange),
				Recommendation: ta.getRecommendation(change.MetricName, change.Direction),
				DetectedAt:    time.Now(),
				Trend:         change.Direction,
			}
			alerts = append(alerts, alert)
			alertID++
		}
	}
	
	return alerts
}

// getThresholdForMetric returns the threshold for a given metric
func (ta *TrendAnalyser) getThresholdForMetric(metricName string) float64 {
	switch metricName {
	case "Error Rate":
		return ta.config.ErrorRateThreshold
	case "Average Response Size":
		return ta.config.ResponseTimeThreshold
	case "Traffic Volume", "Request Volume", "Unique Visitors":
		return ta.config.TrafficDropThreshold
	default:
		return 10.0 // Default 10% threshold
	}
}

// getImpactDescription describes the impact of the metric change
func (ta *TrendAnalyser) getImpactDescription(metricName string, percentChange float64) string {
	absChange := math.Abs(percentChange)
	
	switch metricName {
	case "Error Rate":
		if absChange > 50 {
			return "Severe impact on user experience and system reliability"
		} else if absChange > 20 {
			return "Significant impact on user experience"
		}
		return "Moderate impact on user experience"
	case "Average Response Size":
		if absChange > 50 {
			return "Severe performance degradation affecting all users"
		} else if absChange > 20 {
			return "Noticeable performance impact on users"
		}
		return "Minor performance impact"
	case "Traffic Volume", "Request Volume":
		if absChange > 50 {
			return "Major traffic loss potentially affecting business metrics"
		} else if absChange > 20 {
			return "Significant traffic reduction"
		}
		return "Moderate traffic changes"
	default:
		return "Impact requires investigation"
	}
}

// getRecommendation provides recommendations based on the metric and trend
func (ta *TrendAnalyser) getRecommendation(metricName string, direction TrendDirection) string {
	switch metricName {
	case "Error Rate":
		if direction == TrendCritical {
			return "Immediately investigate error logs and check application health"
		}
		return "Review error patterns and application logs"
	case "Average Response Size":
		return "Check server resources, database performance, and network conditions"
	case "Traffic Volume", "Request Volume":
		return "Investigate traffic sources, check for outages or routing issues"
	case "Unique Visitors":
		return "Analyze user behavior patterns and check marketing campaigns"
	default:
		return "Monitor metric closely and investigate root causes"
	}
}

// calculateOverallHealth determines system health based on alerts and risk score
func (ta *TrendAnalyser) calculateOverallHealth(alerts []DegradationAlert, riskScore int) string {
	criticalAlerts := 0
	for _, alert := range alerts {
		if alert.Severity == "critical" {
			criticalAlerts++
		}
	}
	
	if criticalAlerts > 0 || riskScore > 70 {
		return "critical"
	} else if len(alerts) > 0 || riskScore > 30 {
		return "warning"
	}
	return "healthy"
}

// generateRecommendations creates actionable recommendations
func (ta *TrendAnalyser) generateRecommendations(alerts []DegradationAlert, changes []TrendChange) []string {
	recommendations := make(map[string]bool) // Use map to avoid duplicates
	
	// Add recommendations from alerts
	for _, alert := range alerts {
		recommendations[alert.Recommendation] = true
	}
	
	// Add general recommendations based on trends
	if len(alerts) > 2 {
		recommendations["Perform comprehensive system health check"] = true
	}
	
	// Convert map to slice
	var result []string
	for rec := range recommendations {
		result = append(result, rec)
	}
	
	sort.Strings(result) // Sort for consistent output
	return result
}

// generateTrendSummary creates an executive summary
func (ta *TrendAnalyser) generateTrendSummary(comparison *PeriodComparison, alerts []DegradationAlert) string {
	summary := fmt.Sprintf("Analysis shows %s trend with risk score %d/100. ", 
		comparison.OverallTrend.String(), comparison.RiskScore)
	
	if len(alerts) > 0 {
		summary += fmt.Sprintf("%d degradation alerts generated. ", len(alerts))
	}
	
	// Add key findings
	significantChanges := 0
	for _, change := range comparison.TrendChanges {
		if change.Significance == "high" {
			significantChanges++
		}
	}
	
	if significantChanges > 0 {
		summary += fmt.Sprintf("%d metrics show significant changes requiring attention.", significantChanges)
	} else {
		summary += "No significant changes detected in key metrics."
	}
	
	return summary
}