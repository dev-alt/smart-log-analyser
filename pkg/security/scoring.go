package security

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// SecurityScorer implements security scoring and risk assessment algorithms
type SecurityScorer struct {
	config SecurityConfig
}

// NewSecurityScorer creates a new security scorer
func NewSecurityScorer(config SecurityConfig) *SecurityScorer {
	return &SecurityScorer{
		config: config,
	}
}

// CalculateSecurityScore calculates overall security score (0-100)
func (ss *SecurityScorer) CalculateSecurityScore(analysis *EnhancedSecurityAnalysis) int {
	dimensions := ss.CalculateSecurityDimensions(analysis)
	
	// Weighted scoring system
	weights := SecurityDimensionWeights{
		ThreatDetection:  0.40, // 40% - Direct threats are most important
		AnomalyDetection: 0.25, // 25% - Behavioral anomalies
		TrafficIntegrity: 0.20, // 20% - Overall traffic health
		AccessControl:    0.15, // 15% - Authentication/authorization issues
	}
	
	weightedScore := (dimensions.ThreatDetection * weights.ThreatDetection) +
		(dimensions.AnomalyDetection * weights.AnomalyDetection) +
		(dimensions.TrafficIntegrity * weights.TrafficIntegrity) +
		(dimensions.AccessControl * weights.AccessControl)
	
	return int(math.Round(weightedScore))
}

// SecurityDimensionWeights defines the weight distribution for scoring
type SecurityDimensionWeights struct {
	ThreatDetection  float64
	AnomalyDetection float64
	TrafficIntegrity float64
	AccessControl    float64
}

// CalculateSecurityDimensions calculates security scores across different dimensions
func (ss *SecurityScorer) CalculateSecurityDimensions(analysis *EnhancedSecurityAnalysis) SecurityDimensions {
	return SecurityDimensions{
		ThreatDetection:  ss.calculateThreatDetectionScore(analysis),
		AnomalyDetection: ss.calculateAnomalyDetectionScore(analysis),
		TrafficIntegrity: ss.calculateTrafficIntegrityScore(analysis),
		AccessControl:    ss.calculateAccessControlScore(analysis),
	}
}

// calculateThreatDetectionScore calculates threat detection dimension score (0-100)
func (ss *SecurityScorer) calculateThreatDetectionScore(analysis *EnhancedSecurityAnalysis) float64 {
	if analysis.TotalEntriesAnalyzed == 0 {
		return 100.0 // No data means perfect score
	}

	baseScore := 100.0
	totalThreats := len(analysis.Threats)
	
	if totalThreats == 0 {
		return baseScore
	}

	// Calculate threat impact based on severity and frequency
	threatImpact := 0.0
	severityWeights := map[ThreatSeverity]float64{
		SeverityInfo:     1.0,
		SeverityLow:      2.0,
		SeverityMedium:   5.0,
		SeverityHigh:     10.0,
		SeverityCritical: 20.0,
	}

	for _, threat := range analysis.Threats {
		if weight, exists := severityWeights[threat.Severity]; exists {
			threatImpact += weight * threat.Confidence
		}
	}

	// Normalize threat impact per 1000 requests
	normalizedImpact := (threatImpact / float64(analysis.TotalEntriesAnalyzed)) * 1000

	// Apply logarithmic penalty (diminishing returns for high threat counts)
	penalty := math.Log10(normalizedImpact+1) * 15
	
	score := baseScore - penalty
	return math.Max(0, math.Min(100, score))
}

// calculateAnomalyDetectionScore calculates anomaly detection dimension score (0-100)
func (ss *SecurityScorer) calculateAnomalyDetectionScore(analysis *EnhancedSecurityAnalysis) float64 {
	if analysis.TotalEntriesAnalyzed == 0 {
		return 100.0
	}

	baseScore := 100.0
	totalAnomalies := len(analysis.Anomalies)
	
	if totalAnomalies == 0 {
		return baseScore
	}

	// Calculate anomaly impact
	anomalyImpact := 0.0
	severityWeights := map[ThreatSeverity]float64{
		SeverityInfo:     0.5,
		SeverityLow:      1.0,
		SeverityMedium:   2.5,
		SeverityHigh:     5.0,
		SeverityCritical: 10.0,
	}

	for _, anomaly := range analysis.Anomalies {
		if weight, exists := severityWeights[anomaly.Severity]; exists {
			// Weight by confidence and z-score magnitude
			impact := weight * anomaly.Confidence * (1 + math.Min(math.Abs(anomaly.ZScore)/5.0, 1.0))
			anomalyImpact += impact
		}
	}

	// Normalize per 1000 requests
	normalizedImpact := (anomalyImpact / float64(analysis.TotalEntriesAnalyzed)) * 1000

	// Apply penalty
	penalty := math.Log10(normalizedImpact+1) * 12
	
	score := baseScore - penalty
	return math.Max(0, math.Min(100, score))
}

// calculateTrafficIntegrityScore calculates traffic integrity dimension score (0-100)
func (ss *SecurityScorer) calculateTrafficIntegrityScore(analysis *EnhancedSecurityAnalysis) float64 {
	baseScore := 100.0
	penalties := 0.0

	// Analyze IP behavior profiles for integrity indicators
	if len(analysis.IPProfiles) > 0 {
		suspiciousIPs := 0
		totalIPs := len(analysis.IPProfiles)

		for _, profile := range analysis.IPProfiles {
			// High behavior score indicates suspicious activity
			if profile.BehaviorScore > 0.7 {
				suspiciousIPs++
			}
			
			// High error rate indicates potential issues
			if profile.ErrorRate > 0.3 {
				penalties += 5.0
			}
			
			// Geographic inconsistency
			if !profile.GeographicConsistency {
				penalties += 2.0
			}
		}

		// Penalty for high ratio of suspicious IPs
		suspiciousRatio := float64(suspiciousIPs) / float64(totalIPs)
		penalties += suspiciousRatio * 30.0
	}

	// Check for traffic distribution health
	if analysis.TotalEntriesAnalyzed > 0 {
		// Analyze request distribution across IPs
		ipRequestCounts := make(map[string]int64)
		for _, profile := range analysis.IPProfiles {
			ipRequestCounts[profile.IP] = profile.TotalRequests
		}

		if len(ipRequestCounts) > 0 {
			// Calculate traffic concentration
			var counts []int64
			var totalRequests int64
			for _, count := range ipRequestCounts {
				counts = append(counts, count)
				totalRequests += count
			}

			if totalRequests > 0 {
				sort.Slice(counts, func(i, j int) bool { return counts[i] > counts[j] })
				
				// Check if top IPs dominate traffic (indication of potential DDoS or bot activity)
				if len(counts) > 0 {
					topIPTraffic := float64(counts[0]) / float64(totalRequests)
					if topIPTraffic > 0.5 { // Single IP > 50% of traffic
						penalties += 20.0
					} else if topIPTraffic > 0.3 { // Single IP > 30% of traffic
						penalties += 10.0
					}
				}

				if len(counts) >= 3 {
					top3Traffic := float64(counts[0]+counts[1]+counts[2]) / float64(totalRequests)
					if top3Traffic > 0.8 { // Top 3 IPs > 80% of traffic
						penalties += 15.0
					}
				}
			}
		}
	}

	score := baseScore - penalties
	return math.Max(0, math.Min(100, score))
}

// calculateAccessControlScore calculates access control dimension score (0-100)
func (ss *SecurityScorer) calculateAccessControlScore(analysis *EnhancedSecurityAnalysis) float64 {
	baseScore := 100.0
	penalties := 0.0

	// Analyze authentication-related threats
	authThreats := 0
	for _, threat := range analysis.Threats {
		switch threat.Type {
		case BruteForceLogin:
			penalties += 15.0
			authThreats++
		case PasswordSpray:
			penalties += 12.0
			authThreats++
		case AuthenticationBypass:
			penalties += 20.0
			authThreats++
		case SessionHijacking:
			penalties += 18.0
			authThreats++
		}
	}

	// Analyze access-related anomalies
	for _, anomaly := range analysis.Anomalies {
		if anomaly.Type == AnomalyErrorRate {
			// High 401/403 error rates indicate access control issues
			if anomaly.ActualValue > 0.2 { // More than 20% auth errors
				penalties += 10.0
			}
		}
	}

	// Check for privilege escalation attempts
	escalationThreats := 0
	for _, threat := range analysis.Threats {
		if threat.Type == PrivilegeEscalation {
			penalties += 25.0
			escalationThreats++
		}
	}

	score := baseScore - penalties
	return math.Max(0, math.Min(100, score))
}

// CalculateRiskLevel determines overall risk level based on security score
func (ss *SecurityScorer) CalculateRiskLevel(securityScore int) RiskLevel {
	if securityScore >= 90 {
		return RiskMinimal
	} else if securityScore >= 70 {
		return RiskLow
	} else if securityScore >= 50 {
		return RiskMedium
	} else if securityScore >= 30 {
		return RiskHigh
	}
	return RiskCritical
}

// GenerateSecuritySummary creates a comprehensive security summary
func (ss *SecurityScorer) GenerateSecuritySummary(analysis *EnhancedSecurityAnalysis) SecuritySummary {
	securityScore := ss.CalculateSecurityScore(analysis)
	dimensions := ss.CalculateSecurityDimensions(analysis)
	riskLevel := ss.CalculateRiskLevel(securityScore)

	// Identify high-risk IPs
	var highRiskIPs []string
	for ip, profile := range analysis.IPProfiles {
		if profile.RiskLevel >= RiskHigh {
			highRiskIPs = append(highRiskIPs, ip)
		}
	}

	// Identify top attack types
	topAttackTypes := ss.identifyTopAttackTypes(analysis.Threats)

	// Calculate threat trends
	threatTrends := ss.calculateThreatTrends(analysis)

	// Generate recommended actions
	recommendedActions := ss.generateSecurityRecommendations(analysis, securityScore, dimensions)

	// Count critical vulnerabilities
	criticalVulns := 0
	for _, threat := range analysis.Threats {
		if threat.Severity == SeverityCritical {
			criticalVulns++
		}
	}

	return SecuritySummary{
		OverallRisk:         riskLevel,
		SecurityScore:       securityScore,
		SecurityDimensions:  dimensions,
		ActiveThreats:       len(analysis.Threats),
		CriticalVulns:       criticalVulns,
		HighRiskIPs:         highRiskIPs,
		TopAttackTypes:      topAttackTypes,
		ThreatTrends:        threatTrends,
		RecommendedActions:  recommendedActions,
		ComplianceScore:     ss.calculateComplianceScore(analysis),
		IncidentCount:       len(analysis.Incidents),
		TimeRange: TimeRange{
			Start: analysis.LogTimeRange.Start,
			End:   analysis.LogTimeRange.End,
		},
	}
}

// identifyTopAttackTypes identifies the most common attack types
func (ss *SecurityScorer) identifyTopAttackTypes(threats []EnhancedThreat) []string {
	attackCounts := make(map[string]int)
	
	for _, threat := range threats {
		var attackType string
		switch t := threat.Type.(type) {
		case WebAttackType:
			attackType = t.String()
		case InfrastructureAttackType:
			attackType = t.String()
		default:
			attackType = "Unknown"
		}
		attackCounts[attackType]++
	}

	// Sort by count
	type attackCount struct {
		name  string
		count int
	}
	
	var counts []attackCount
	for name, count := range attackCounts {
		counts = append(counts, attackCount{name, count})
	}
	
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})

	var topTypes []string
	for i, ac := range counts {
		if i >= 5 { // Top 5 attack types
			break
		}
		topTypes = append(topTypes, ac.name)
	}

	return topTypes
}

// calculateThreatTrends calculates threat trends over time
func (ss *SecurityScorer) calculateThreatTrends(analysis *EnhancedSecurityAnalysis) []ThreatTrend {
	var trends []ThreatTrend

	// Group threats by type and time windows
	timeWindow := time.Hour
	threatsByType := make(map[string]map[int64]int)

	for _, threat := range analysis.Threats {
		var attackType string
		switch t := threat.Type.(type) {
		case WebAttackType:
			attackType = t.String()
		case InfrastructureAttackType:
			attackType = t.String()
		default:
			attackType = "Unknown"
		}

		windowTime := threat.Timestamp.Truncate(timeWindow).Unix()
		if threatsByType[attackType] == nil {
			threatsByType[attackType] = make(map[int64]int)
		}
		threatsByType[attackType][windowTime]++
	}

	// Calculate trends for each attack type
	for attackType, timeWindows := range threatsByType {
		if len(timeWindows) < 2 {
			continue // Need at least 2 time windows for trend
		}

		var timestamps []int64
		var counts []int
		totalCount := 0

		for timestamp, count := range timeWindows {
			timestamps = append(timestamps, timestamp)
			counts = append(counts, count)
			totalCount += count
		}

		// Sort by timestamp
		sort.Slice(timestamps, func(i, j int) bool {
			return timestamps[i] < timestamps[j]
		})

		// Calculate simple trend (comparing first half vs second half)
		midpoint := len(timestamps) / 2
		firstHalfTotal := 0
		secondHalfTotal := 0

		for i, timestamp := range timestamps {
			count := timeWindows[timestamp]
			if i < midpoint {
				firstHalfTotal += count
			} else {
				secondHalfTotal += count
			}
		}

		trend := 0.0
		if firstHalfTotal > 0 {
			trend = float64(secondHalfTotal-firstHalfTotal) / float64(firstHalfTotal)
		}

		trends = append(trends, ThreatTrend{
			Type:  attackType,
			Count: totalCount,
			Trend: trend,
			TimeRange: TimeRange{
				Start: time.Unix(timestamps[0], 0),
				End:   time.Unix(timestamps[len(timestamps)-1], 0),
			},
		})
	}

	return trends
}

// generateSecurityRecommendations generates actionable security recommendations
func (ss *SecurityScorer) generateSecurityRecommendations(analysis *EnhancedSecurityAnalysis, securityScore int, dimensions SecurityDimensions) []SecurityRecommendation {
	var recommendations []SecurityRecommendation

	// Critical score recommendations
	if securityScore < 30 {
		recommendations = append(recommendations, SecurityRecommendation{
			Priority:    1,
			Category:    "Critical Security Alert",
			Title:       "Immediate Security Response Required",
			Description: "Security score is critically low. Immediate investigation and response required.",
			Impact:      SeverityCritical,
			Effort:      "High",
			Actions: []string{
				"Activate incident response procedures",
				"Review and block suspicious IPs immediately",
				"Analyze recent security logs in detail",
				"Consider temporary access restrictions",
			},
		})
	}

	// Threat detection recommendations
	if dimensions.ThreatDetection < 70 {
		criticalThreats := 0
		highThreats := 0
		for _, threat := range analysis.Threats {
			if threat.Severity == SeverityCritical {
				criticalThreats++
			} else if threat.Severity == SeverityHigh {
				highThreats++
			}
		}

		if criticalThreats > 0 {
			recommendations = append(recommendations, SecurityRecommendation{
				Priority:    2,
				Category:    "Threat Detection",
				Title:       "Critical Threats Detected",
				Description: fmt.Sprintf("%d critical threats require immediate attention", criticalThreats),
				Impact:      SeverityCritical,
				Effort:      "Medium",
				Actions: []string{
					"Review and investigate all critical threats",
					"Implement blocking rules for malicious IPs",
					"Update security signatures and rules",
					"Monitor for related attack patterns",
				},
			})
		}

		if highThreats > 5 {
			recommendations = append(recommendations, SecurityRecommendation{
				Priority:    3,
				Category:    "Threat Detection",
				Title:       "Multiple High-Severity Threats",
				Description: fmt.Sprintf("%d high-severity threats detected", highThreats),
				Impact:      SeverityHigh,
				Effort:      "Medium",
				Actions: []string{
					"Implement additional monitoring for affected endpoints",
					"Review and update security policies",
					"Consider rate limiting for suspicious patterns",
				},
			})
		}
	}

	// Anomaly detection recommendations
	if dimensions.AnomalyDetection < 70 {
		highAnomalies := 0
		for _, anomaly := range analysis.Anomalies {
			if anomaly.Severity >= SeverityHigh {
				highAnomalies++
			}
		}

		if highAnomalies > 0 {
			recommendations = append(recommendations, SecurityRecommendation{
				Priority:    4,
				Category:    "Anomaly Detection",
				Title:       "Significant Behavioral Anomalies",
				Description: fmt.Sprintf("%d high-severity anomalies detected in traffic patterns", highAnomalies),
				Impact:      SeverityHigh,
				Effort:      "Low",
				Actions: []string{
					"Investigate anomalous IP behavior patterns",
					"Implement behavioral-based blocking rules",
					"Review user agent and geographic patterns",
					"Consider implementing CAPTCHA for suspicious activity",
				},
			})
		}
	}

	// Traffic integrity recommendations
	if dimensions.TrafficIntegrity < 70 {
		recommendations = append(recommendations, SecurityRecommendation{
			Priority:    5,
			Category:    "Traffic Integrity",
			Title:       "Traffic Pattern Concerns",
			Description: "Unusual traffic patterns detected that may indicate bot activity or attacks",
			Impact:      SeverityMedium,
			Effort:      "Medium",
			Actions: []string{
				"Implement bot detection and mitigation",
				"Review traffic distribution across IPs",
				"Consider implementing rate limiting",
				"Monitor for distributed attack patterns",
			},
		})
	}

	// Access control recommendations
	if dimensions.AccessControl < 70 {
		authThreats := 0
		for _, threat := range analysis.Threats {
			switch threat.Type {
			case BruteForceLogin, PasswordSpray, AuthenticationBypass:
				authThreats++
			}
		}

		if authThreats > 0 {
			recommendations = append(recommendations, SecurityRecommendation{
				Priority:    6,
				Category:    "Access Control",
				Title:       "Authentication Security Issues",
				Description: fmt.Sprintf("%d authentication-related threats detected", authThreats),
				Impact:      SeverityHigh,
				Effort:      "Low",
				Actions: []string{
					"Implement account lockout policies",
					"Enable multi-factor authentication",
					"Review and strengthen password policies",
					"Monitor failed authentication attempts",
				},
			})
		}
	}

	// General security hardening
	if securityScore < 80 {
		recommendations = append(recommendations, SecurityRecommendation{
			Priority:    7,
			Category:    "Security Hardening",
			Title:       "General Security Improvements",
			Description: "Multiple areas for security improvement identified",
			Impact:      SeverityMedium,
			Effort:      "Medium",
			Actions: []string{
				"Review and update security policies",
				"Implement comprehensive logging and monitoring",
				"Regular security training for staff",
				"Schedule regular security assessments",
			},
		})
	}

	// Sort recommendations by priority
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Priority < recommendations[j].Priority
	})

	return recommendations
}

// calculateComplianceScore calculates compliance score based on security posture
func (ss *SecurityScorer) calculateComplianceScore(analysis *EnhancedSecurityAnalysis) int {
	baseScore := 100
	penalties := 0

	// Penalty for high-severity threats (compliance frameworks require threat management)
	criticalThreats := 0
	for _, threat := range analysis.Threats {
		if threat.Severity == SeverityCritical {
			criticalThreats++
		}
	}
	penalties += criticalThreats * 5

	// Penalty for security incidents
	penalties += len(analysis.Incidents) * 3

	// Penalty for high-risk IPs not being managed
	highRiskIPs := 0
	for _, profile := range analysis.IPProfiles {
		if profile.RiskLevel >= RiskHigh {
			highRiskIPs++
		}
	}
	penalties += highRiskIPs * 2

	score := baseScore - penalties
	return int(math.Max(0, math.Min(100, float64(score))))
}

// GenerateIncidents creates incident data from threats and anomalies
func (ss *SecurityScorer) GenerateIncidents(threats []EnhancedThreat, anomalies []Anomaly) ([]IncidentData, error) {
	var incidents []IncidentData

	// Group related threats into incidents
	incidentGroups := ss.groupThreatsIntoIncidents(threats)

	for i, group := range incidentGroups {
		if len(group) == 0 {
			continue
		}

		// Find the most severe threat in the group
		maxSeverity := SeverityInfo
		var primaryThreat EnhancedThreat
		for _, threat := range group {
			if threat.Severity > maxSeverity {
				maxSeverity = threat.Severity
				primaryThreat = threat
			}
		}

		// Create incident
		incident := IncidentData{
			ID:       fmt.Sprintf("INC-%d-%d", time.Now().Unix(), i+1),
			Title:    ss.generateIncidentTitle(group),
			Severity: maxSeverity,
			StartTime: group[0].Timestamp,
			EndTime:   group[len(group)-1].Timestamp,
			AffectedSystems: ss.extractAffectedSystems(group),
			AttackVector:    primaryThreat.AttackVector,
			ThreatActor:     ss.identifyThreatActor(group),
			IOCs:            ss.extractIOCs(group),
			Timeline:        ss.createIncidentTimeline(group),
			Impact:          ss.assessIncidentImpact(group),
			Recommendations: ss.generateIncidentRecommendations(group),
			Evidence:        ss.gatherEvidence(group),
		}

		incidents = append(incidents, incident)
	}

	return incidents, nil
}

// groupThreatsIntoIncidents groups related threats into incidents
func (ss *SecurityScorer) groupThreatsIntoIncidents(threats []EnhancedThreat) [][]EnhancedThreat {
	var groups [][]EnhancedThreat

	// Sort threats by timestamp
	sortedThreats := make([]EnhancedThreat, len(threats))
	copy(sortedThreats, threats)
	sort.Slice(sortedThreats, func(i, j int) bool {
		return sortedThreats[i].Timestamp.Before(sortedThreats[j].Timestamp)
	})

	// Group threats by IP and time proximity
	ipGroups := make(map[string][]EnhancedThreat)
	for _, threat := range sortedThreats {
		ipGroups[threat.IP] = append(ipGroups[threat.IP], threat)
	}

	// Create incident groups
	for _, ipThreats := range ipGroups {
		if len(ipThreats) >= 3 { // Minimum 3 threats for an incident
			groups = append(groups, ipThreats)
		}
	}

	return groups
}

// generateIncidentTitle creates a descriptive title for an incident
func (ss *SecurityScorer) generateIncidentTitle(threats []EnhancedThreat) string {
	if len(threats) == 0 {
		return "Security Incident"
	}

	// Count attack types
	attackTypes := make(map[string]int)
	for _, threat := range threats {
		var attackType string
		switch t := threat.Type.(type) {
		case WebAttackType:
			attackType = t.String()
		case InfrastructureAttackType:
			attackType = t.String()
		default:
			attackType = "Unknown Attack"
		}
		attackTypes[attackType]++
	}

	// Find most common attack type
	maxCount := 0
	primaryAttack := "Mixed Attacks"
	for attackType, count := range attackTypes {
		if count > maxCount {
			maxCount = count
			primaryAttack = attackType
		}
	}

	if len(attackTypes) == 1 {
		return fmt.Sprintf("%s from %s", primaryAttack, threats[0].IP)
	} else {
		return fmt.Sprintf("Multi-vector Attack (%s and %d others) from %s", primaryAttack, len(attackTypes)-1, threats[0].IP)
	}
}

// extractAffectedSystems identifies systems affected by threats
func (ss *SecurityScorer) extractAffectedSystems(threats []EnhancedThreat) []string {
	systems := make(map[string]bool)
	for _, threat := range threats {
		if threat.URL != "" {
			systems[threat.URL] = true
		}
	}

	var result []string
	for system := range systems {
		result = append(result, system)
	}
	return result
}

// identifyThreatActor attempts to identify the threat actor
func (ss *SecurityScorer) identifyThreatActor(threats []EnhancedThreat) string {
	if len(threats) == 0 {
		return "Unknown"
	}

	// Simple heuristic based on IP and attack patterns
	ip := threats[0].IP
	
	// Check for automated vs manual patterns
	automated := true
	if len(threats) > 1 {
		// Check time intervals between attacks
		var intervals []time.Duration
		for i := 1; i < len(threats); i++ {
			interval := threats[i].Timestamp.Sub(threats[i-1].Timestamp)
			intervals = append(intervals, interval)
		}
		
		// If intervals are too regular, likely automated
		if len(intervals) > 0 {
			var totalInterval time.Duration
			for _, interval := range intervals {
				totalInterval += interval
			}
			avgInterval := totalInterval / time.Duration(len(intervals))
			
			if avgInterval < 10*time.Second {
				return fmt.Sprintf("Automated Tool (%s)", ip)
			}
		}
	}

	if automated {
		return fmt.Sprintf("Automated Attacker (%s)", ip)
	}
	return fmt.Sprintf("Manual Attacker (%s)", ip)
}

// extractIOCs extracts Indicators of Compromise from threats
func (ss *SecurityScorer) extractIOCs(threats []EnhancedThreat) []string {
	iocs := make(map[string]bool)
	
	for _, threat := range threats {
		// Add IP as IOC
		iocs[fmt.Sprintf("IP: %s", threat.IP)] = true
		
		// Add user agent if suspicious
		if threat.UserAgent != "" && (strings.Contains(strings.ToLower(threat.UserAgent), "bot") ||
			strings.Contains(strings.ToLower(threat.UserAgent), "scanner")) {
			iocs[fmt.Sprintf("User-Agent: %s", threat.UserAgent)] = true
		}
		
		// Add payload patterns
		if threat.Payload != "" {
			iocs[fmt.Sprintf("Payload Pattern: %s", threat.Payload)] = true
		}
	}

	var result []string
	for ioc := range iocs {
		result = append(result, ioc)
	}
	return result
}

// createIncidentTimeline creates a timeline of incident events
func (ss *SecurityScorer) createIncidentTimeline(threats []EnhancedThreat) []IncidentEvent {
	var timeline []IncidentEvent

	for i, threat := range threats {
		var attackType string
		switch t := threat.Type.(type) {
		case WebAttackType:
			attackType = t.String()
		case InfrastructureAttackType:
			attackType = t.String()
		default:
			attackType = "Unknown Attack"
		}

		event := IncidentEvent{
			Timestamp:   threat.Timestamp,
			Description: fmt.Sprintf("%s detected targeting %s", attackType, threat.URL),
			Type:        "Attack",
			Severity:    threat.Severity,
			Source:      threat.IP,
			Details: map[string]interface{}{
				"attack_type":   attackType,
				"confidence":    threat.Confidence,
				"attack_vector": threat.AttackVector,
				"payload":       threat.Payload,
			},
		}

		if i == 0 {
			event.Description = fmt.Sprintf("Incident started: %s", event.Description)
		} else if i == len(threats)-1 {
			event.Description = fmt.Sprintf("Latest activity: %s", event.Description)
		}

		timeline = append(timeline, event)
	}

	return timeline
}

// assessIncidentImpact assesses the impact of an incident
func (ss *SecurityScorer) assessIncidentImpact(threats []EnhancedThreat) string {
	if len(threats) == 0 {
		return "Unknown impact"
	}

	criticalCount := 0
	highCount := 0
	for _, threat := range threats {
		if threat.Severity == SeverityCritical {
			criticalCount++
		} else if threat.Severity == SeverityHigh {
			highCount++
		}
	}

	if criticalCount > 0 {
		return fmt.Sprintf("CRITICAL: %d critical threats detected with potential for data breach or system compromise", criticalCount)
	} else if highCount > 0 {
		return fmt.Sprintf("HIGH: %d high-severity threats detected with potential for unauthorized access", highCount)
	} else {
		return fmt.Sprintf("MEDIUM: %d security threats detected with limited immediate impact", len(threats))
	}
}

// generateIncidentRecommendations generates specific recommendations for an incident
func (ss *SecurityScorer) generateIncidentRecommendations(threats []EnhancedThreat) []SecurityRecommendation {
	var recommendations []SecurityRecommendation

	if len(threats) == 0 {
		return recommendations
	}

	// Immediate blocking recommendation
	recommendations = append(recommendations, SecurityRecommendation{
		Priority:    1,
		Category:    "Immediate Response",
		Title:       "Block Malicious IP",
		Description: fmt.Sprintf("Immediately block IP %s to prevent further attacks", threats[0].IP),
		Impact:      SeverityHigh,
		Effort:      "Low",
		Actions: []string{
			fmt.Sprintf("Add firewall rule to block %s", threats[0].IP),
			"Monitor for additional attacks from related IPs",
			"Review logs for any successful attacks",
		},
	})

	// Attack-specific recommendations
	attackTypes := make(map[interface{}]bool)
	for _, threat := range threats {
		if !attackTypes[threat.Type] {
			attackTypes[threat.Type] = true
			
			var actions []string
			switch t := threat.Type.(type) {
			case WebAttackType:
				switch t {
				case SQLInjection:
					actions = []string{"Review SQL injection protections", "Implement parameterized queries", "Update WAF rules"}
				case CrossSiteScripting:
					actions = []string{"Implement output encoding", "Review CSP headers", "Update XSS protections"}
				case CommandInjection:
					actions = []string{"Review command execution code", "Implement input validation", "Apply principle of least privilege"}
				default:
					actions = []string{"Review application security controls", "Update security signatures"}
				}
			case InfrastructureAttackType:
				switch t {
				case BruteForceLogin:
					actions = []string{"Implement account lockout", "Enable MFA", "Review authentication logs"}
				case DDoSAttack:
					actions = []string{"Activate DDoS protection", "Scale infrastructure", "Monitor traffic patterns"}
				default:
					actions = []string{"Review infrastructure security", "Update monitoring rules"}
				}
			}

			if len(actions) > 0 {
				recommendations = append(recommendations, SecurityRecommendation{
					Priority:    2,
					Category:    "Attack Mitigation",
					Title:       fmt.Sprintf("Mitigate %v", threat.Type),
					Description: fmt.Sprintf("Address %v vulnerability", threat.Type),
					Impact:      threat.Severity,
					Effort:      "Medium",
					Actions:     actions,
				})
			}
		}
	}

	return recommendations
}

// gatherEvidence gathers evidence for an incident
func (ss *SecurityScorer) gatherEvidence(threats []EnhancedThreat) []string {
	var evidence []string

	for _, threat := range threats {
		evidence = append(evidence, fmt.Sprintf("Timestamp: %s", threat.Timestamp.Format(time.RFC3339)))
		evidence = append(evidence, fmt.Sprintf("Source IP: %s", threat.IP))
		evidence = append(evidence, fmt.Sprintf("Attack Type: %v", threat.Type))
		evidence = append(evidence, fmt.Sprintf("Target URL: %s", threat.URL))
		if threat.Payload != "" {
			evidence = append(evidence, fmt.Sprintf("Attack Payload: %s", threat.Payload))
		}
		evidence = append(evidence, "---")
	}

	return evidence
}