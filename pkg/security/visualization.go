package security

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"smart-log-analyser/pkg/charts"
)

// SecurityVisualizer implements security-focused visualization and reporting
type SecurityVisualizer struct {
	config SecurityConfig
}

// NewSecurityVisualizer creates a new security visualizer
func NewSecurityVisualizer(config SecurityConfig) *SecurityVisualizer {
	return &SecurityVisualizer{
		config: config,
	}
}

// GenerateSecurityDashboard creates a comprehensive ASCII security dashboard
func (sv *SecurityVisualizer) GenerateSecurityDashboard(analysis *EnhancedSecurityAnalysis) string {
	var output strings.Builder
	
	// Header
	output.WriteString("╔══════════════════════════════════════════════════════════════╗\n")
	output.WriteString("║                 🔐 SECURITY ANALYSIS DASHBOARD               ║\n")
	output.WriteString("╚══════════════════════════════════════════════════════════════╝\n\n")

	// Security Overview Card
	output.WriteString(sv.generateSecurityOverviewCard(analysis))
	
	// Risk Level Indicator
	output.WriteString(sv.generateRiskLevelIndicator(analysis.Summary.OverallRisk))
	
	// Security Dimensions Chart
	output.WriteString(sv.generateSecurityDimensionsChart(analysis.Summary.SecurityDimensions))
	
	// Threat Distribution Chart
	output.WriteString(sv.generateThreatDistributionChart(analysis.Threats))
	
	// High-Risk IPs Table
	if len(analysis.Summary.HighRiskIPs) > 0 {
		output.WriteString(sv.generateHighRiskIPsTable(analysis.IPProfiles, analysis.Summary.HighRiskIPs))
	}
	
	// Recent Incidents Summary
	if len(analysis.Incidents) > 0 {
		output.WriteString(sv.generateIncidentsSummary(analysis.Incidents))
	}
	
	// Security Recommendations
	if len(analysis.Summary.RecommendedActions) > 0 {
		output.WriteString(sv.generateRecommendationsCard(analysis.Summary.RecommendedActions))
	}
	
	return output.String()
}

// generateSecurityOverviewCard creates a security overview summary card
func (sv *SecurityVisualizer) generateSecurityOverviewCard(analysis *EnhancedSecurityAnalysis) string {
	var output strings.Builder
	
	output.WriteString("┌─ SECURITY OVERVIEW ─────────────────────────────────────────┐\n")
	
	// Security Score with color
	scoreColor := sv.getScoreColor(analysis.Summary.SecurityScore)
	output.WriteString(fmt.Sprintf("│ Security Score: %s%d/100%s", scoreColor, analysis.Summary.SecurityScore, charts.ColorReset))
	output.WriteString(strings.Repeat(" ", 39-len(fmt.Sprintf("%d/100", analysis.Summary.SecurityScore))))
	output.WriteString("│\n")
	
	// Risk Level
	riskColor := sv.getRiskColor(analysis.Summary.OverallRisk)
	output.WriteString(fmt.Sprintf("│ Risk Level:     %s%-12s%s", riskColor, analysis.Summary.OverallRisk.String(), charts.ColorReset))
	output.WriteString(strings.Repeat(" ", 36))
	output.WriteString("│\n")
	
	// Threats and Anomalies
	output.WriteString(fmt.Sprintf("│ Active Threats: %-8d", analysis.Summary.ActiveThreats))
	output.WriteString(fmt.Sprintf(" │ Critical Vulns: %-8d │\n", analysis.Summary.CriticalVulns))
	
	output.WriteString(fmt.Sprintf("│ Anomalies:      %-8d", len(analysis.Anomalies)))
	output.WriteString(fmt.Sprintf(" │ High-Risk IPs:  %-8d │\n", len(analysis.Summary.HighRiskIPs)))
	
	// Time Range
	timeRange := fmt.Sprintf("%s to %s", 
		analysis.Summary.TimeRange.Start.Format("Jan 02 15:04"),
		analysis.Summary.TimeRange.End.Format("Jan 02 15:04"))
	output.WriteString(fmt.Sprintf("│ Analysis Period: %-43s │\n", timeRange))
	
	// Total Entries
	output.WriteString(fmt.Sprintf("│ Log Entries:     %-43d │\n", analysis.TotalEntriesAnalyzed))
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// generateRiskLevelIndicator creates a visual risk level indicator
func (sv *SecurityVisualizer) generateRiskLevelIndicator(riskLevel RiskLevel) string {
	var output strings.Builder
	
	output.WriteString("┌─ RISK LEVEL INDICATOR ──────────────────────────────────────┐\n")
	
	levels := []struct{
		level RiskLevel
		name  string
		icon  string
	}{
		{RiskMinimal, "MINIMAL", "🟢"},
		{RiskLow, "LOW", "🟡"},
		{RiskMedium, "MEDIUM", "🟠"},
		{RiskHigh, "HIGH", "🔴"},
		{RiskCritical, "CRITICAL", "🚨"},
	}
	
	for _, level := range levels {
		indicator := "  "
		if level.level == riskLevel {
			indicator = "▶ "
		}
		
		color := sv.getRiskColor(level.level)
		output.WriteString(fmt.Sprintf("│ %s%s%-8s%s %s", indicator, color, level.name, charts.ColorReset, level.icon))
		output.WriteString(strings.Repeat(" ", 44-len(level.name)))
		output.WriteString("│\n")
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// generateSecurityDimensionsChart creates a chart showing security dimension scores
func (sv *SecurityVisualizer) generateSecurityDimensionsChart(dimensions SecurityDimensions) string {
	var output strings.Builder
	
	output.WriteString("┌─ SECURITY DIMENSIONS ───────────────────────────────────────┐\n")
	
	dimensionData := []struct{
		name  string
		score float64
		weight string
	}{
		{"Threat Detection", dimensions.ThreatDetection, "40%"},
		{"Anomaly Detection", dimensions.AnomalyDetection, "25%"},
		{"Traffic Integrity", dimensions.TrafficIntegrity, "20%"},
		{"Access Control", dimensions.AccessControl, "15%"},
	}
	
	for _, dim := range dimensionData {
		// Create bar visualization
		barLength := int(dim.score * 40 / 100) // Scale to 40 characters max
		bar := strings.Repeat("█", barLength)
		bar += strings.Repeat("░", 40-barLength)
		
		color := sv.getScoreColor(int(dim.score))
		output.WriteString(fmt.Sprintf("│ %-17s │%s%s%s│ %3.0f%% (%s) │\n", 
			dim.name, color, bar, charts.ColorReset, dim.score, dim.weight))
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// generateThreatDistributionChart creates a chart showing threat type distribution
func (sv *SecurityVisualizer) generateThreatDistributionChart(threats []EnhancedThreat) string {
	var output strings.Builder
	
	if len(threats) == 0 {
		return ""
	}
	
	output.WriteString("┌─ THREAT DISTRIBUTION ───────────────────────────────────────┐\n")
	
	// Count threats by type
	threatCounts := make(map[string]int)
	for _, threat := range threats {
		var threatType string
		switch t := threat.Type.(type) {
		case WebAttackType:
			threatType = t.String()
		case InfrastructureAttackType:
			threatType = t.String()
		default:
			threatType = "Unknown"
		}
		threatCounts[threatType]++
	}
	
	// Sort by count
	type threatCount struct {
		name  string
		count int
	}
	
	var sortedThreats []threatCount
	for name, count := range threatCounts {
		sortedThreats = append(sortedThreats, threatCount{name, count})
	}
	
	sort.Slice(sortedThreats, func(i, j int) bool {
		return sortedThreats[i].count > sortedThreats[j].count
	})
	
	// Display top 8 threats
	maxCount := 0
	if len(sortedThreats) > 0 {
		maxCount = sortedThreats[0].count
	}
	
	displayCount := 8
	if len(sortedThreats) < displayCount {
		displayCount = len(sortedThreats)
	}
	
	for i := 0; i < displayCount; i++ {
		threat := sortedThreats[i]
		
		// Create bar visualization
		barLength := 30
		if maxCount > 0 {
			barLength = int(float64(threat.count) * 30.0 / float64(maxCount))
		}
		
		bar := strings.Repeat("█", barLength)
		bar += strings.Repeat("░", 30-barLength)
		
		// Truncate long threat names
		name := threat.name
		if len(name) > 20 {
			name = name[:17] + "..."
		}
		
		output.WriteString(fmt.Sprintf("│ %-20s │%s%s%s│ %4d │\n", 
			name, charts.ColorRed, bar, charts.ColorReset, threat.count))
	}
	
	if len(sortedThreats) > displayCount {
		output.WriteString(fmt.Sprintf("│ ... and %d more threat types", len(sortedThreats)-displayCount))
		output.WriteString(strings.Repeat(" ", 60-len(fmt.Sprintf("... and %d more threat types", len(sortedThreats)-displayCount))))
		output.WriteString("│\n")
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// generateHighRiskIPsTable creates a table of high-risk IP addresses
func (sv *SecurityVisualizer) generateHighRiskIPsTable(profiles map[string]*IPBehaviorProfile, highRiskIPs []string) string {
	var output strings.Builder
	
	output.WriteString("┌─ HIGH-RISK IP ADDRESSES ────────────────────────────────────┐\n")
	output.WriteString("│ IP Address      │ Risk Level │ Requests │ Error Rate │ Score │\n")
	output.WriteString("├─────────────────┼────────────┼──────────┼────────────┼───────┤\n")
	
	displayCount := 10
	if len(highRiskIPs) < displayCount {
		displayCount = len(highRiskIPs)
	}
	
	for i := 0; i < displayCount; i++ {
		ip := highRiskIPs[i]
		profile, exists := profiles[ip]
		if !exists {
			continue
		}
		
		riskColor := sv.getRiskColor(profile.RiskLevel)
		
		output.WriteString(fmt.Sprintf("│ %-15s │ %s%-10s%s │ %8d │ %8.1f%% │ %5.2f │\n",
			ip,
			riskColor, profile.RiskLevel.String(), charts.ColorReset,
			profile.TotalRequests,
			profile.ErrorRate*100,
			profile.BehaviorScore))
	}
	
	if len(highRiskIPs) > displayCount {
		output.WriteString(fmt.Sprintf("│ ... and %d more high-risk IPs", len(highRiskIPs)-displayCount))
		output.WriteString(strings.Repeat(" ", 62-len(fmt.Sprintf("... and %d more high-risk IPs", len(highRiskIPs)-displayCount))))
		output.WriteString("│\n")
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// generateIncidentsSummary creates a summary of recent security incidents
func (sv *SecurityVisualizer) generateIncidentsSummary(incidents []IncidentData) string {
	var output strings.Builder
	
	output.WriteString("┌─ RECENT SECURITY INCIDENTS ─────────────────────────────────┐\n")
	
	displayCount := 5
	if len(incidents) < displayCount {
		displayCount = len(incidents)
	}
	
	// Sort incidents by severity and recency
	sortedIncidents := make([]IncidentData, len(incidents))
	copy(sortedIncidents, incidents)
	sort.Slice(sortedIncidents, func(i, j int) bool {
		if sortedIncidents[i].Severity != sortedIncidents[j].Severity {
			return sortedIncidents[i].Severity > sortedIncidents[j].Severity
		}
		return sortedIncidents[i].EndTime.After(sortedIncidents[j].EndTime)
	})
	
	for i := 0; i < displayCount; i++ {
		incident := sortedIncidents[i]
		
		severityColor := sv.getSeverityColor(incident.Severity)
		title := incident.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		
		timeStr := incident.EndTime.Format("Jan 02 15:04")
		
		output.WriteString(fmt.Sprintf("│ %s%-9s%s │ %-40s │ %s │\n",
			severityColor, incident.Severity.String(), charts.ColorReset,
			title, timeStr))
	}
	
	if len(incidents) > displayCount {
		output.WriteString(fmt.Sprintf("│ ... and %d more incidents", len(incidents)-displayCount))
		output.WriteString(strings.Repeat(" ", 62-len(fmt.Sprintf("... and %d more incidents", len(incidents)-displayCount))))
		output.WriteString("│\n")
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// generateRecommendationsCard creates a card showing top security recommendations
func (sv *SecurityVisualizer) generateRecommendationsCard(recommendations []SecurityRecommendation) string {
	var output strings.Builder
	
	output.WriteString("┌─ TOP SECURITY RECOMMENDATIONS ──────────────────────────────┐\n")
	
	displayCount := 5
	if len(recommendations) < displayCount {
		displayCount = len(recommendations)
	}
	
	for i := 0; i < displayCount; i++ {
		rec := recommendations[i]
		
		impactColor := sv.getSeverityColor(rec.Impact)
		title := rec.Title
		if len(title) > 45 {
			title = title[:42] + "..."
		}
		
		output.WriteString(fmt.Sprintf("│ %d. %s%-7s%s │ %-45s │\n",
			rec.Priority,
			impactColor, rec.Impact.String(), charts.ColorReset,
			title))
	}
	
	if len(recommendations) > displayCount {
		output.WriteString(fmt.Sprintf("│ ... and %d more recommendations", len(recommendations)-displayCount))
		output.WriteString(strings.Repeat(" ", 62-len(fmt.Sprintf("... and %d more recommendations", len(recommendations)-displayCount))))
		output.WriteString("│\n")
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// GenerateThreatTimelineChart creates a timeline visualization of threats
func (sv *SecurityVisualizer) GenerateThreatTimelineChart(threats []EnhancedThreat, timeWindow time.Duration) string {
	var output strings.Builder
	
	if len(threats) == 0 {
		return ""
	}
	
	output.WriteString("┌─ THREAT TIMELINE ───────────────────────────────────────────┐\n")
	
	// Group threats by time windows
	timeGroups := make(map[int64]int)
	var minTime, maxTime int64 = math.MaxInt64, 0
	
	for _, threat := range threats {
		windowTime := threat.Timestamp.Truncate(timeWindow).Unix()
		timeGroups[windowTime]++
		if windowTime < minTime {
			minTime = windowTime
		}
		if windowTime > maxTime {
			maxTime = windowTime
		}
	}
	
	// Create timeline
	if maxTime > minTime {
		windowCount := (maxTime - minTime) / int64(timeWindow.Seconds()) + 1
		maxWindowCount := 0
		for _, count := range timeGroups {
			if count > maxWindowCount {
				maxWindowCount = count
			}
		}
		
		// Display timeline (limit to reasonable number of windows)
		displayWindows := int(math.Min(float64(windowCount), 20))
		windowStep := windowCount / int64(displayWindows)
		if windowStep < 1 {
			windowStep = 1
		}
		
		for i := int64(0); i < int64(displayWindows); i++ {
			windowTime := minTime + i*windowStep*int64(timeWindow.Seconds())
			count := timeGroups[windowTime]
			
			// Create bar
			barLength := 40
			if maxWindowCount > 0 {
				barLength = int(float64(count) * 40.0 / float64(maxWindowCount))
			}
			
			bar := strings.Repeat("█", barLength)
			bar += strings.Repeat("░", 40-barLength)
			
			timeStr := time.Unix(windowTime, 0).Format("15:04")
			
			output.WriteString(fmt.Sprintf("│ %5s │%s%s%s│ %4d │\n",
				timeStr, charts.ColorRed, bar, charts.ColorReset, count))
		}
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// GenerateAnomalyHeatMap creates a heat map visualization of anomalies
func (sv *SecurityVisualizer) GenerateAnomalyHeatMap(anomalies []Anomaly) string {
	var output strings.Builder
	
	if len(anomalies) == 0 {
		return ""
	}
	
	output.WriteString("┌─ ANOMALY HEAT MAP ──────────────────────────────────────────┐\n")
	
	// Group anomalies by type and severity
	anomalyMatrix := make(map[AnomalyType]map[ThreatSeverity]int)
	
	for _, anomaly := range anomalies {
		if anomalyMatrix[anomaly.Type] == nil {
			anomalyMatrix[anomaly.Type] = make(map[ThreatSeverity]int)
		}
		anomalyMatrix[anomaly.Type][anomaly.Severity]++
	}
	
	// Display matrix
	severities := []ThreatSeverity{SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical}
	output.WriteString("│ Anomaly Type        │ Low │ Med │High│Crit│ Total │\n")
	output.WriteString("├─────────────────────┼─────┼─────┼────┼────┼───────┤\n")
	
	for anomalyType, severityMap := range anomalyMatrix {
		total := 0
		counts := make([]int, 4)
		
		for i, severity := range severities {
			count := severityMap[severity]
			counts[i] = count
			total += count
		}
		
		typeStr := anomalyType.String()
		if len(typeStr) > 19 {
			typeStr = typeStr[:16] + "..."
		}
		
		output.WriteString(fmt.Sprintf("│ %-19s │ %3d │ %3d │%4d│%4d│ %5d │\n",
			typeStr, counts[0], counts[1], counts[2], counts[3], total))
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// GenerateIPBehaviorChart creates a chart showing IP behavior analysis
func (sv *SecurityVisualizer) GenerateIPBehaviorChart(profiles map[string]*IPBehaviorProfile, topN int) string {
	var output strings.Builder
	
	if len(profiles) == 0 {
		return ""
	}
	
	output.WriteString("┌─ IP BEHAVIOR ANALYSIS ──────────────────────────────────────┐\n")
	output.WriteString("│ IP Address      │ Behavior Score │ Risk Level │ Requests │\n")
	output.WriteString("├─────────────────┼────────────────┼────────────┼──────────┤\n")
	
	// Sort profiles by behavior score
	type profileData struct {
		ip      string
		profile *IPBehaviorProfile
	}
	
	var sortedProfiles []profileData
	for ip, profile := range profiles {
		sortedProfiles = append(sortedProfiles, profileData{ip, profile})
	}
	
	sort.Slice(sortedProfiles, func(i, j int) bool {
		return sortedProfiles[i].profile.BehaviorScore > sortedProfiles[j].profile.BehaviorScore
	})
	
	displayCount := topN
	if len(sortedProfiles) < displayCount {
		displayCount = len(sortedProfiles)
	}
	
	for i := 0; i < displayCount; i++ {
		pd := sortedProfiles[i]
		
		// Create behavior score bar
		barLength := int(pd.profile.BehaviorScore * 10)
		bar := strings.Repeat("█", barLength)
		bar += strings.Repeat("░", 10-barLength)
		
		scoreColor := charts.ColorGreen
		if pd.profile.BehaviorScore > 0.7 {
			scoreColor = charts.ColorRed
		} else if pd.profile.BehaviorScore > 0.4 {
			scoreColor = charts.ColorYellow
		}
		
		riskColor := sv.getRiskColor(pd.profile.RiskLevel)
		
		output.WriteString(fmt.Sprintf("│ %-15s │ %s%s%s %.2f │ %s%-10s%s │ %8d │\n",
			pd.ip,
			scoreColor, bar, charts.ColorReset, pd.profile.BehaviorScore,
			riskColor, pd.profile.RiskLevel.String(), charts.ColorReset,
			pd.profile.TotalRequests))
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// GenerateSecurityTrendChart creates a trend analysis chart
func (sv *SecurityVisualizer) GenerateSecurityTrendChart(trends []ThreatTrend) string {
	var output strings.Builder
	
	if len(trends) == 0 {
		return ""
	}
	
	output.WriteString("┌─ SECURITY TRENDS ───────────────────────────────────────────┐\n")
	output.WriteString("│ Threat Type         │ Count │ Trend │ Direction       │\n")
	output.WriteString("├─────────────────────┼───────┼───────┼─────────────────┤\n")
	
	for _, trend := range trends {
		trendStr := fmt.Sprintf("%+.1f%%", trend.Trend*100)
		
		// Trend visualization
		direction := "Stable"
		directionColor := charts.ColorBlue
		if trend.Trend > 0.1 {
			direction = "↗ Increasing"
			directionColor = charts.ColorRed
		} else if trend.Trend < -0.1 {
			direction = "↘ Decreasing" 
			directionColor = charts.ColorGreen
		}
		
		threatType := trend.Type
		if len(threatType) > 19 {
			threatType = threatType[:16] + "..."
		}
		
		output.WriteString(fmt.Sprintf("│ %-19s │ %5d │ %5s │ %s%-15s%s │\n",
			threatType, trend.Count, trendStr,
			directionColor, direction, charts.ColorReset))
	}
	
	output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	
	return output.String()
}

// Color helper functions

func (sv *SecurityVisualizer) getScoreColor(score int) string {
	if score >= 80 {
		return charts.ColorGreen
	} else if score >= 60 {
		return charts.ColorYellow
	} else if score >= 40 {
		return charts.ColorRed
	}
	return charts.ColorRed + charts.ColorBold
}

func (sv *SecurityVisualizer) getRiskColor(risk RiskLevel) string {
	switch risk {
	case RiskMinimal:
		return charts.ColorGreen
	case RiskLow:
		return charts.ColorBlue
	case RiskMedium:
		return charts.ColorYellow
	case RiskHigh:
		return charts.ColorRed
	case RiskCritical:
		return charts.ColorRed + charts.ColorBold
	default:
		return charts.ColorReset
	}
}

func (sv *SecurityVisualizer) getSeverityColor(severity ThreatSeverity) string {
	switch severity {
	case SeverityInfo:
		return charts.ColorBlue
	case SeverityLow:
		return charts.ColorGreen
	case SeverityMedium:
		return charts.ColorYellow
	case SeverityHigh:
		return charts.ColorRed
	case SeverityCritical:
		return charts.ColorRed + charts.ColorBold
	default:
		return charts.ColorReset
	}
}

// GenerateDetailedThreatReport creates a detailed threat analysis report
func (sv *SecurityVisualizer) GenerateDetailedThreatReport(threats []EnhancedThreat) string {
	var output strings.Builder
	
	if len(threats) == 0 {
		return "No threats detected.\n"
	}
	
	output.WriteString("╔══════════════════════════════════════════════════════════════╗\n")
	output.WriteString("║                    DETAILED THREAT REPORT                   ║\n")
	output.WriteString("╚══════════════════════════════════════════════════════════════╝\n\n")
	
	// Group by severity
	severityGroups := make(map[ThreatSeverity][]EnhancedThreat)
	for _, threat := range threats {
		severityGroups[threat.Severity] = append(severityGroups[threat.Severity], threat)
	}
	
	// Display by severity (highest first)
	severityOrder := []ThreatSeverity{SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo}
	
	for _, severity := range severityOrder {
		threatList, exists := severityGroups[severity]
		if !exists || len(threatList) == 0 {
			continue
		}
		
		severityColor := sv.getSeverityColor(severity)
		output.WriteString(fmt.Sprintf("┌─ %s%s THREATS (%d)%s", 
			severityColor, severity.String(), len(threatList), charts.ColorReset))
		output.WriteString(strings.Repeat("─", 62-len(fmt.Sprintf("%s THREATS (%d)", severity.String(), len(threatList)))))
		output.WriteString("┐\n")
		
		// Show top 5 threats of this severity
		displayCount := 5
		if len(threatList) < displayCount {
			displayCount = len(threatList)
		}
		
		for i := 0; i < displayCount; i++ {
			threat := threatList[i]
			
			var threatType string
			switch t := threat.Type.(type) {
			case WebAttackType:
				threatType = t.String()
			case InfrastructureAttackType:
				threatType = t.String()
			default:
				threatType = "Unknown"
			}
			
			output.WriteString(fmt.Sprintf("│ %s from %s at %s\n",
				threatType, threat.IP, threat.Timestamp.Format("15:04:05")))
			
			if threat.URL != "" {
				url := threat.URL
				if len(url) > 55 {
					url = url[:52] + "..."
				}
				output.WriteString(fmt.Sprintf("│ Target: %s\n", url))
			}
			
			if threat.Payload != "" {
				payload := threat.Payload
				if len(payload) > 55 {
					payload = payload[:52] + "..."
				}
				output.WriteString(fmt.Sprintf("│ Payload: %s\n", payload))
			}
			
			output.WriteString(fmt.Sprintf("│ Confidence: %.0f%% │ Attack Vector: %s\n", 
				threat.Confidence*100, threat.AttackVector))
			
			if i < displayCount-1 {
				output.WriteString("├─────────────────────────────────────────────────────────────┤\n")
			}
		}
		
		if len(threatList) > displayCount {
			output.WriteString(fmt.Sprintf("│ ... and %d more %s threats\n", 
				len(threatList)-displayCount, strings.ToLower(severity.String())))
		}
		
		output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	}
	
	return output.String()
}

// GenerateAnomalyReport creates a detailed anomaly analysis report
func (sv *SecurityVisualizer) GenerateAnomalyReport(anomalies []Anomaly) string {
	var output strings.Builder
	
	if len(anomalies) == 0 {
		return "No anomalies detected.\n"
	}
	
	output.WriteString("╔══════════════════════════════════════════════════════════════╗\n")
	output.WriteString("║                    ANOMALY ANALYSIS REPORT                  ║\n")
	output.WriteString("╚══════════════════════════════════════════════════════════════╝\n\n")
	
	// Group by type
	typeGroups := make(map[AnomalyType][]Anomaly)
	for _, anomaly := range anomalies {
		typeGroups[anomaly.Type] = append(typeGroups[anomaly.Type], anomaly)
	}
	
	for anomalyType, anomalyList := range typeGroups {
		output.WriteString(fmt.Sprintf("┌─ %s (%d)%s┐\n",
			anomalyType.String(), len(anomalyList),
			strings.Repeat("─", 60-len(fmt.Sprintf("%s (%d)", anomalyType.String(), len(anomalyList))))))
		
		// Sort by severity and z-score
		sort.Slice(anomalyList, func(i, j int) bool {
			if anomalyList[i].Severity != anomalyList[j].Severity {
				return anomalyList[i].Severity > anomalyList[j].Severity
			}
			return math.Abs(anomalyList[i].ZScore) > math.Abs(anomalyList[j].ZScore)
		})
		
		displayCount := 3
		if len(anomalyList) < displayCount {
			displayCount = len(anomalyList)
		}
		
		for i := 0; i < displayCount; i++ {
			anomaly := anomalyList[i]
			
			severityColor := sv.getSeverityColor(anomaly.Severity)
			output.WriteString(fmt.Sprintf("│ %s%s%s │ IP: %s │ Z-Score: %.2f\n",
				severityColor, anomaly.Severity.String(), charts.ColorReset,
				anomaly.IP, anomaly.ZScore))
			
			output.WriteString(fmt.Sprintf("│ %s\n", anomaly.Description))
			
			output.WriteString(fmt.Sprintf("│ Expected: %.2f │ Actual: %.2f │ Confidence: %.0f%%\n",
				anomaly.ExpectedValue, anomaly.ActualValue, anomaly.Confidence*100))
			
			if i < displayCount-1 {
				output.WriteString("├─────────────────────────────────────────────────────────────┤\n")
			}
		}
		
		if len(anomalyList) > displayCount {
			output.WriteString(fmt.Sprintf("│ ... and %d more %s anomalies\n",
				len(anomalyList)-displayCount, strings.ToLower(anomalyType.String())))
		}
		
		output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	}
	
	return output.String()
}

// GenerateSecurityRecommendationReport creates a detailed recommendations report
func (sv *SecurityVisualizer) GenerateSecurityRecommendationReport(recommendations []SecurityRecommendation) string {
	var output strings.Builder
	
	if len(recommendations) == 0 {
		return "No specific recommendations at this time.\n"
	}
	
	output.WriteString("╔══════════════════════════════════════════════════════════════╗\n")
	output.WriteString("║              SECURITY RECOMMENDATIONS REPORT                ║\n")
	output.WriteString("╚══════════════════════════════════════════════════════════════╝\n\n")
	
	for i, rec := range recommendations {
		if i >= 10 { // Limit to top 10 recommendations
			break
		}
		
		impactColor := sv.getSeverityColor(rec.Impact)
		output.WriteString(fmt.Sprintf("┌─ RECOMMENDATION #%d ─ %s%s IMPACT%s", 
			rec.Priority, impactColor, rec.Impact.String(), charts.ColorReset))
		output.WriteString(strings.Repeat("─", 60-len(fmt.Sprintf("RECOMMENDATION #%d ─ %s IMPACT", rec.Priority, rec.Impact.String()))))
		output.WriteString("┐\n")
		
		output.WriteString(fmt.Sprintf("│ Category: %s\n", rec.Category))
		output.WriteString(fmt.Sprintf("│ Title: %s\n", rec.Title))
		output.WriteString(fmt.Sprintf("│ Effort Level: %s\n", rec.Effort))
		output.WriteString("│\n")
		output.WriteString(fmt.Sprintf("│ Description:\n"))
		
		// Word wrap description
		words := strings.Fields(rec.Description)
		line := "│ "
		for _, word := range words {
			if len(line)+len(word)+1 > 62 {
				output.WriteString(line + "\n")
				line = "│ " + word
			} else {
				if len(line) > 2 {
					line += " "
				}
				line += word
			}
		}
		if len(line) > 2 {
			output.WriteString(line + "\n")
		}
		
		output.WriteString("│\n")
		output.WriteString("│ Recommended Actions:\n")
		for j, action := range rec.Actions {
			if j >= 5 { // Limit to 5 actions
				break
			}
			output.WriteString(fmt.Sprintf("│ %d. %s\n", j+1, action))
		}
		
		output.WriteString("└─────────────────────────────────────────────────────────────┘\n\n")
	}
	
	return output.String()
}