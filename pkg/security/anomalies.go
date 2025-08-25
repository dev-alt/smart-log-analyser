package security

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"smart-log-analyser/pkg/parser"
)

// AnomalyDetector implements ML-based anomaly detection algorithms
type AnomalyDetector struct {
	config             SecurityConfig
	baselineEstablished bool
	behaviorProfiles   map[string]*IPBehaviorProfile
	globalBaseline     *GlobalBaseline
}

// GlobalBaseline represents normal system behavior patterns
type GlobalBaseline struct {
	AverageRequestsPerMinute float64
	AverageSize      int64
	CommonStatusCodes        map[int]float64
	CommonUserAgents         map[string]float64
	CommonEndpoints          map[string]float64
	TypicalRequestTimes      []time.Duration
	ErrorRateThreshold       float64
	SizeDistribution         SizeDistribution
	TimeDistribution         TimeDistribution
}

// SizeDistribution represents response size distribution
type SizeDistribution struct {
	Mean   float64
	StdDev float64
	P95    int64
	P99    int64
}

// TimeDistribution represents request timing distribution
type TimeDistribution struct {
	PeakHours    []int // Hours 0-23
	OffPeakHours []int
	WeekdayVsWeekend float64 // Ratio
}

// MetricSample represents a data point for statistical analysis
type MetricSample struct {
	Value     float64
	Timestamp time.Time
	IP        string
	Context   map[string]interface{}
}

// NewAnomalyDetector creates a new anomaly detector
func NewAnomalyDetector(config SecurityConfig) *AnomalyDetector {
	return &AnomalyDetector{
		config:           config,
		behaviorProfiles: make(map[string]*IPBehaviorProfile),
		globalBaseline: &GlobalBaseline{
			CommonStatusCodes: make(map[int]float64),
			CommonUserAgents:  make(map[string]float64),
			CommonEndpoints:   make(map[string]float64),
		},
	}
}

// DetectAnomalies identifies behavioral anomalies in log entries
func (ad *AnomalyDetector) DetectAnomalies(logs []*parser.LogEntry) ([]Anomaly, error) {
	var anomalies []Anomaly

	// Establish baseline if not already done
	if !ad.baselineEstablished {
		ad.establishBaseline(logs)
	}

	// Update IP behavior profiles
	ad.updateBehaviorProfiles(logs)

	// Detect various types of anomalies
	requestFreqAnomalies := ad.detectRequestFrequencyAnomalies(logs)
	anomalies = append(anomalies, requestFreqAnomalies...)

	requestSizeAnomalies := ad.detectRequestSizeAnomalies(logs)
	anomalies = append(anomalies, requestSizeAnomalies...)

	errorRateAnomalies := ad.detectErrorRateAnomalies(logs)
	anomalies = append(anomalies, errorRateAnomalies...)

	timingAnomalies := ad.detectRequestTimingAnomalies(logs)
	anomalies = append(anomalies, timingAnomalies...)

	userAgentAnomalies := ad.detectUserAgentAnomalies(logs)
	anomalies = append(anomalies, userAgentAnomalies...)

	geographicAnomalies := ad.detectGeographicAnomalies(logs)
	anomalies = append(anomalies, geographicAnomalies...)

	endpointAnomalies := ad.detectEndpointPatternAnomalies(logs)
	anomalies = append(anomalies, endpointAnomalies...)

	statusCodeAnomalies := ad.detectStatusAnomalies(logs)
	anomalies = append(anomalies, statusCodeAnomalies...)

	return anomalies, nil
}

// ProfileIPs creates behavioral profiles for IP addresses
func (ad *AnomalyDetector) ProfileIPs(logs []*parser.LogEntry) (map[string]*IPBehaviorProfile, error) {
	profiles := make(map[string]*IPBehaviorProfile)

	// Group entries by IP
	ipEntries := make(map[string][]*parser.LogEntry)
	for _, entry := range logs {
		ipEntries[entry.IP] = append(ipEntries[entry.IP], entry)
	}

	// Create profile for each IP
	for ip, entries := range ipEntries {
		profile := ad.createIPBehaviorProfile(ip, entries)
		profile.BehaviorScore = ad.calculateBehaviorScore(profile)
		profile.RiskLevel = ad.assessRiskLevel(profile.BehaviorScore)
		profile.Anomalies = ad.detectIPSpecificAnomalies(profile, entries)
		profiles[ip] = profile
	}

	return profiles, nil
}

// establishBaseline establishes normal behavior patterns from log data
func (ad *AnomalyDetector) establishBaseline(logs []*parser.LogEntry) {
	if len(logs) == 0 {
		return
	}

	// Calculate request frequency baseline
	duration := logs[len(logs)-1].Timestamp.Sub(logs[0].Timestamp)
	if duration > 0 {
		ad.globalBaseline.AverageRequestsPerMinute = float64(len(logs)) / duration.Minutes()
	}

	// Calculate response size baseline
	var totalSize int64
	var sizes []int64
	for _, entry := range logs {
		totalSize += entry.Size
		sizes = append(sizes, entry.Size)
	}
	if len(logs) > 0 {
		ad.globalBaseline.AverageSize = totalSize / int64(len(logs))
	}

	// Calculate size distribution
	sort.Slice(sizes, func(i, j int) bool { return sizes[i] < sizes[j] })
	if len(sizes) > 0 {
		ad.globalBaseline.SizeDistribution.Mean = float64(totalSize) / float64(len(sizes))
		
		// Calculate standard deviation
		var variance float64
		for _, size := range sizes {
			variance += math.Pow(float64(size)-ad.globalBaseline.SizeDistribution.Mean, 2)
		}
		ad.globalBaseline.SizeDistribution.StdDev = math.Sqrt(variance / float64(len(sizes)))

		// Calculate percentiles
		p95Index := int(0.95 * float64(len(sizes)))
		p99Index := int(0.99 * float64(len(sizes)))
		if p95Index < len(sizes) {
			ad.globalBaseline.SizeDistribution.P95 = sizes[p95Index]
		}
		if p99Index < len(sizes) {
			ad.globalBaseline.SizeDistribution.P99 = sizes[p99Index]
		}
	}

	// Analyze status codes
	statusCodeCounts := make(map[int]int)
	errorCount := 0
	for _, entry := range logs {
		statusCodeCounts[entry.Status]++
		if entry.Status >= 400 {
			errorCount++
		}
	}

	// Calculate status code frequencies
	for code, count := range statusCodeCounts {
		ad.globalBaseline.CommonStatusCodes[code] = float64(count) / float64(len(logs))
	}

	// Calculate error rate threshold
	ad.globalBaseline.ErrorRateThreshold = float64(errorCount) / float64(len(logs))
	if ad.globalBaseline.ErrorRateThreshold < 0.05 {
		ad.globalBaseline.ErrorRateThreshold = 0.05 // Minimum 5% threshold
	}

	// Analyze user agents
	userAgentCounts := make(map[string]int)
	for _, entry := range logs {
		userAgentCounts[entry.UserAgent]++
	}

	// Calculate user agent frequencies (top 20)
	type uaCount struct {
		ua    string
		count int
	}
	var uaCounts []uaCount
	for ua, count := range userAgentCounts {
		uaCounts = append(uaCounts, uaCount{ua, count})
	}
	sort.Slice(uaCounts, func(i, j int) bool { return uaCounts[i].count > uaCounts[j].count })

	for i, uac := range uaCounts {
		if i >= 20 { // Top 20 user agents
			break
		}
		ad.globalBaseline.CommonUserAgents[uac.ua] = float64(uac.count) / float64(len(logs))
	}

	// Analyze endpoints
	endpointCounts := make(map[string]int)
	for _, entry := range logs {
		// Extract path from URL
		path := strings.Split(entry.URL, "?")[0] // Remove query parameters
		endpointCounts[path]++
	}

	// Calculate endpoint frequencies (top 50)
	type epCount struct {
		endpoint string
		count    int
	}
	var epCounts []epCount
	for ep, count := range endpointCounts {
		epCounts = append(epCounts, epCount{ep, count})
	}
	sort.Slice(epCounts, func(i, j int) bool { return epCounts[i].count > epCounts[j].count })

	for i, epc := range epCounts {
		if i >= 50 { // Top 50 endpoints
			break
		}
		ad.globalBaseline.CommonEndpoints[epc.endpoint] = float64(epc.count) / float64(len(logs))
	}

	// Analyze time patterns
	hourCounts := make(map[int]int)
	for _, entry := range logs {
		hour := entry.Timestamp.Hour()
		hourCounts[hour]++
	}

	// Identify peak hours (above average)
	averagePerHour := float64(len(logs)) / 24.0
	for hour, count := range hourCounts {
		if float64(count) > averagePerHour*1.2 { // 20% above average
			ad.globalBaseline.TimeDistribution.PeakHours = append(ad.globalBaseline.TimeDistribution.PeakHours, hour)
		} else if float64(count) < averagePerHour*0.8 { // 20% below average
			ad.globalBaseline.TimeDistribution.OffPeakHours = append(ad.globalBaseline.TimeDistribution.OffPeakHours, hour)
		}
	}

	ad.baselineEstablished = true
}

// updateBehaviorProfiles updates IP behavior profiles with new data
func (ad *AnomalyDetector) updateBehaviorProfiles(logs []*parser.LogEntry) {
	ipEntries := make(map[string][]*parser.LogEntry)
	for _, entry := range logs {
		ipEntries[entry.IP] = append(ipEntries[entry.IP], entry)
	}

	for ip, entries := range ipEntries {
		if profile, exists := ad.behaviorProfiles[ip]; exists {
			ad.updateExistingProfile(profile, entries)
		} else {
			ad.behaviorProfiles[ip] = ad.createIPBehaviorProfile(ip, entries)
		}
	}
}

// createIPBehaviorProfile creates a new behavioral profile for an IP
func (ad *AnomalyDetector) createIPBehaviorProfile(ip string, entries []*parser.LogEntry) *IPBehaviorProfile {
	if len(entries) == 0 {
		return &IPBehaviorProfile{IP: ip}
	}

	profile := &IPBehaviorProfile{
		IP:                      ip,
		FirstSeen:               entries[0].Timestamp,
		LastSeen:                entries[len(entries)-1].Timestamp,
		TotalRequests:           int64(len(entries)),
		CommonUserAgents:        make(map[string]int),
		VisitedEndpoints:        make(map[string]int),
		HTTPMethods:             make(map[string]int),
		StatusCodeDistribution:  make(map[int]int),
		GeographicLocations:     []string{},
		AssociatedThreats:       []string{},
		Tags:                    []string{},
	}

	// Calculate request frequency
	duration := profile.LastSeen.Sub(profile.FirstSeen)
	if duration > 0 {
		profile.RequestFrequency = float64(len(entries)) / duration.Minutes()
	}

	// Calculate average request interval
	if len(entries) > 1 {
		var totalInterval time.Duration
		for i := 1; i < len(entries); i++ {
			totalInterval += entries[i].Timestamp.Sub(entries[i-1].Timestamp)
		}
		profile.AverageRequestInterval = totalInterval / time.Duration(len(entries)-1)
	}

	// Analyze user agents
	userAgentCounts := make(map[string]int)
	for _, entry := range entries {
		userAgentCounts[entry.UserAgent]++
	}
	profile.CommonUserAgents = userAgentCounts

	// Analyze visited endpoints
	endpointCounts := make(map[string]int)
	for _, entry := range entries {
		path := strings.Split(entry.URL, "?")[0]
		endpointCounts[path]++
	}
	profile.VisitedEndpoints = endpointCounts

	// Analyze HTTP methods
	methodCounts := make(map[string]int)
	for _, entry := range entries {
		methodCounts[entry.Method]++
	}
	profile.HTTPMethods = methodCounts

	// Analyze status codes
	statusCodeCounts := make(map[int]int)
	errorCount := 0
	var totalSize int64
	
	for _, entry := range entries {
		statusCodeCounts[entry.Status]++
		totalSize += entry.Size
		if entry.Status >= 400 {
			errorCount++
		}
	}
	profile.StatusCodeDistribution = statusCodeCounts
	profile.ErrorRate = float64(errorCount) / float64(len(entries))
	profile.AverageResponseSize = totalSize / int64(len(entries))

	// Analyze typical request times
	for _, entry := range entries {
		profile.TypicalRequestTimes = append(profile.TypicalRequestTimes, entry.Timestamp)
	}

	// Set geographic consistency (simplified - in production would use IP geolocation)
	profile.GeographicConsistency = true // Assume consistent for now

	return profile
}

// updateExistingProfile updates an existing profile with new entries
func (ad *AnomalyDetector) updateExistingProfile(profile *IPBehaviorProfile, entries []*parser.LogEntry) {
	if len(entries) == 0 {
		return
	}

	// Update timestamps
	if entries[0].Timestamp.Before(profile.FirstSeen) {
		profile.FirstSeen = entries[0].Timestamp
	}
	if entries[len(entries)-1].Timestamp.After(profile.LastSeen) {
		profile.LastSeen = entries[len(entries)-1].Timestamp
	}

	// Update request count
	profile.TotalRequests += int64(len(entries))

	// Recalculate request frequency
	duration := profile.LastSeen.Sub(profile.FirstSeen)
	if duration > 0 {
		profile.RequestFrequency = float64(profile.TotalRequests) / duration.Minutes()
	}

	// Update user agents
	for _, entry := range entries {
		profile.CommonUserAgents[entry.UserAgent]++
	}

	// Update endpoints
	for _, entry := range entries {
		path := strings.Split(entry.URL, "?")[0]
		profile.VisitedEndpoints[path]++
	}

	// Update HTTP methods
	for _, entry := range entries {
		profile.HTTPMethods[entry.Method]++
	}

	// Update status codes and calculate new error rate
	errorCount := 0
	for status, count := range profile.StatusCodeDistribution {
		if status >= 400 {
			errorCount += count
		}
	}
	
	for _, entry := range entries {
		profile.StatusCodeDistribution[entry.Status]++
		if entry.Status >= 400 {
			errorCount++
		}
	}
	
	profile.ErrorRate = float64(errorCount) / float64(profile.TotalRequests)
}

// detectRequestFrequencyAnomalies detects unusual request frequency patterns
func (ad *AnomalyDetector) detectRequestFrequencyAnomalies(logs []*parser.LogEntry) []Anomaly {
	var anomalies []Anomaly

	// Group by IP and time windows
	ipRequests := make(map[string][]time.Time)
	for _, entry := range logs {
		ipRequests[entry.IP] = append(ipRequests[entry.IP], entry.Timestamp)
	}

	// Analyze each IP's request frequency
	for ip, timestamps := range ipRequests {
		if len(timestamps) < 10 { // Need minimum requests for meaningful analysis
			continue
		}

		// Sort timestamps
		sort.Slice(timestamps, func(i, j int) bool {
			return timestamps[i].Before(timestamps[j])
		})

		// Calculate requests per minute in sliding windows
		windowSize := 5 * time.Minute
		var frequencies []float64

		for i := 0; i < len(timestamps); i++ {
			windowStart := timestamps[i]
			windowEnd := windowStart.Add(windowSize)
			
			count := 0
			for j := i; j < len(timestamps) && timestamps[j].Before(windowEnd); j++ {
				count++
			}
			
			frequency := float64(count) / windowSize.Minutes()
			frequencies = append(frequencies, frequency)
		}

		// Calculate z-scores
		if len(frequencies) > 5 {
			mean, stdDev := calculateStats(frequencies)
			for i, freq := range frequencies {
				zScore := (freq - mean) / stdDev
				if math.Abs(zScore) > ad.config.AnomalyThreshold {
					severity := SeverityLow
					if math.Abs(zScore) > 3.0 {
						severity = SeverityMedium
					}
					if math.Abs(zScore) > 5.0 {
						severity = SeverityHigh
					}

					anomaly := Anomaly{
						ID:            fmt.Sprintf("freq_%d_%s", time.Now().UnixNano(), ip),
						Type:          AnomalyRequestFrequency,
						Severity:      severity,
						Description:   fmt.Sprintf("Unusual request frequency pattern (%.2f req/min)", freq),
						Metric:        "requests_per_minute",
						ExpectedValue: mean,
						ActualValue:   freq,
						Deviation:     math.Abs(freq - mean),
						ZScore:        zScore,
						IP:            ip,
						Timestamp:     timestamps[i],
						TimeWindow:    windowSize,
						Confidence:    math.Min(math.Abs(zScore)/5.0, 1.0),
						Context:       map[string]interface{}{"window_size_minutes": windowSize.Minutes()},
					}
					anomalies = append(anomalies, anomaly)
				}
			}
		}
	}

	return anomalies
}

// detectRequestSizeAnomalies detects unusual response size patterns
func (ad *AnomalyDetector) detectRequestSizeAnomalies(logs []*parser.LogEntry) []Anomaly {
	var anomalies []Anomaly
	var sizes []float64

	// Collect response sizes
	for _, entry := range logs {
		sizes = append(sizes, float64(entry.Size))
	}

	if len(sizes) < 20 {
		return anomalies
	}

	mean, stdDev := calculateStats(sizes)
	
	// Check each entry for size anomalies
	for _, entry := range logs {
		size := float64(entry.Size)
		zScore := (size - mean) / stdDev
		
		if math.Abs(zScore) > ad.config.AnomalyThreshold {
			severity := SeverityLow
			if math.Abs(zScore) > 3.0 {
				severity = SeverityMedium
			}
			if math.Abs(zScore) > 5.0 {
				severity = SeverityHigh
			}

			anomaly := Anomaly{
				ID:            fmt.Sprintf("size_%d_%s", time.Now().UnixNano(), entry.IP),
				Type:          AnomalyRequestSize,
				Severity:      severity,
				Description:   fmt.Sprintf("Unusual response size (%d bytes)", entry.Size),
				Metric:        "response_size_bytes",
				ExpectedValue: mean,
				ActualValue:   size,
				Deviation:     math.Abs(size - mean),
				ZScore:        zScore,
				IP:            entry.IP,
				Timestamp:     entry.Timestamp,
				TimeWindow:    time.Hour, // Consider hourly window
				Confidence:    math.Min(math.Abs(zScore)/5.0, 1.0),
				Context:       map[string]interface{}{"url": entry.URL},
			}
			anomalies = append(anomalies, anomaly)
		}
	}

	return anomalies
}

// detectErrorRateAnomalies detects unusual error rate patterns
func (ad *AnomalyDetector) detectErrorRateAnomalies(logs []*parser.LogEntry) []Anomaly {
	var anomalies []Anomaly

	// Group by IP and time windows
	windowSize := 10 * time.Minute
	ipWindows := make(map[string]map[int64][]int) // ip -> timestamp -> status codes

	for _, entry := range logs {
		windowTime := entry.Timestamp.Truncate(windowSize).Unix()
		if ipWindows[entry.IP] == nil {
			ipWindows[entry.IP] = make(map[int64][]int)
		}
		ipWindows[entry.IP][windowTime] = append(ipWindows[entry.IP][windowTime], entry.Status)
	}

	// Analyze error rates for each IP
	for ip, windows := range ipWindows {
		var errorRates []float64
		var timestamps []int64

		for timestamp, statusCodes := range windows {
			if len(statusCodes) < 5 { // Need minimum requests
				continue
			}

			errorCount := 0
			for _, code := range statusCodes {
				if code >= 400 {
					errorCount++
				}
			}

			errorRate := float64(errorCount) / float64(len(statusCodes))
			errorRates = append(errorRates, errorRate)
			timestamps = append(timestamps, timestamp)
		}

		if len(errorRates) > 3 {
			mean, stdDev := calculateStats(errorRates)
			
			for i, rate := range errorRates {
				if stdDev > 0 {
					zScore := (rate - mean) / stdDev
					if math.Abs(zScore) > ad.config.AnomalyThreshold && rate > ad.globalBaseline.ErrorRateThreshold*2 {
						severity := SeverityMedium
						if rate > 0.5 { // More than 50% error rate
							severity = SeverityHigh
						}
						if rate > 0.8 { // More than 80% error rate
							severity = SeverityCritical
						}

						anomaly := Anomaly{
							ID:            fmt.Sprintf("error_%d_%s", time.Now().UnixNano(), ip),
							Type:          AnomalyErrorRate,
							Severity:      severity,
							Description:   fmt.Sprintf("Unusual error rate (%.1f%%)", rate*100),
							Metric:        "error_rate_percentage",
							ExpectedValue: mean,
							ActualValue:   rate,
							Deviation:     math.Abs(rate - mean),
							ZScore:        zScore,
							IP:            ip,
							Timestamp:     time.Unix(timestamps[i], 0),
							TimeWindow:    windowSize,
							Confidence:    math.Min(rate*2, 1.0),
							Context:       map[string]interface{}{"baseline_error_rate": ad.globalBaseline.ErrorRateThreshold},
						}
						anomalies = append(anomalies, anomaly)
					}
				}
			}
		}
	}

	return anomalies
}

// detectRequestTimingAnomalies detects unusual request timing patterns
func (ad *AnomalyDetector) detectRequestTimingAnomalies(logs []*parser.LogEntry) []Anomaly {
	var anomalies []Anomaly

	// Group by IP
	ipEntries := make(map[string][]*parser.LogEntry)
	for _, entry := range logs {
		ipEntries[entry.IP] = append(ipEntries[entry.IP], entry)
	}

	// Analyze timing patterns for each IP
	for ip, entries := range ipEntries {
		if len(entries) < 10 {
			continue
		}

		// Sort by timestamp
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp.Before(entries[j].Timestamp)
		})

		// Calculate intervals between requests
		var intervals []float64
		for i := 1; i < len(entries); i++ {
			interval := entries[i].Timestamp.Sub(entries[i-1].Timestamp)
			intervals = append(intervals, interval.Seconds())
		}

		if len(intervals) > 5 {
			mean, stdDev := calculateStats(intervals)
			
			// Check for too regular patterns (bot-like)
			if stdDev < mean*0.1 && mean < 5.0 { // Very regular and fast
				anomaly := Anomaly{
					ID:          fmt.Sprintf("timing_%d_%s", time.Now().UnixNano(), ip),
					Type:        AnomalyRequestTiming,
					Severity:    SeverityMedium,
					Description: fmt.Sprintf("Highly regular request timing (%.2fs intervals)", mean),
					Metric:      "request_interval_seconds",
					ExpectedValue: 10.0, // Expected human-like variance
					ActualValue:   stdDev,
					Deviation:     math.Abs(stdDev - 10.0),
					ZScore:        (10.0 - stdDev) / 5.0,
					IP:            ip,
					Timestamp:     entries[len(entries)-1].Timestamp,
					TimeWindow:    time.Duration(len(entries)) * time.Duration(mean) * time.Second,
					Confidence:    0.8,
					Context:       map[string]interface{}{"pattern": "regular_intervals", "avg_interval": mean},
				}
				anomalies = append(anomalies, anomaly)
			}
		}
	}

	return anomalies
}

// detectUserAgentAnomalies detects unusual user agent patterns
func (ad *AnomalyDetector) detectUserAgentAnomalies(logs []*parser.LogEntry) []Anomaly {
	var anomalies []Anomaly

	// Analyze user agents by IP
	ipUserAgents := make(map[string]map[string]int)
	for _, entry := range logs {
		if ipUserAgents[entry.IP] == nil {
			ipUserAgents[entry.IP] = make(map[string]int)
		}
		ipUserAgents[entry.IP][entry.UserAgent]++
	}

	// Check for suspicious patterns
	for ip, userAgents := range ipUserAgents {
		totalRequests := 0
		for _, count := range userAgents {
			totalRequests += count
		}

		if totalRequests < 5 {
			continue
		}

		// Check for suspicious patterns
		for userAgent, count := range userAgents {
			confidence := 0.0
			severity := SeverityLow
			description := ""

			// Check against common user agents baseline
			if frequency, exists := ad.globalBaseline.CommonUserAgents[userAgent]; exists {
				expectedCount := int(frequency * float64(totalRequests))
				if count > expectedCount*3 { // Much higher than expected
					confidence = 0.6
					description = "User agent frequency higher than baseline"
				}
			} else {
				// Unknown user agent
				confidence = 0.4
				description = "Unknown user agent"
			}

			// Suspicious patterns
			suspiciousPatterns := []string{
				"bot", "crawler", "spider", "scraper",
				"python", "java", "curl", "wget",
				"scanner", "test", "monitor",
			}

			for _, pattern := range suspiciousPatterns {
				if strings.Contains(strings.ToLower(userAgent), pattern) {
					confidence += 0.3
					severity = SeverityMedium
					description += fmt.Sprintf(", contains suspicious term: %s", pattern)
					break
				}
			}

			// Very generic or empty user agent
			if userAgent == "" || userAgent == "Mozilla/5.0" {
				confidence += 0.4
				severity = SeverityMedium
				description += ", generic/empty user agent"
			}

			if confidence > 0.5 {
				anomaly := Anomaly{
					ID:            fmt.Sprintf("ua_%d_%s", time.Now().UnixNano(), ip),
					Type:          AnomalyUserAgent,
					Severity:      severity,
					Description:   fmt.Sprintf("Suspicious user agent pattern%s", description),
					Metric:        "user_agent_suspicion_score",
					ExpectedValue: 0.0,
					ActualValue:   confidence,
					Deviation:     confidence,
					ZScore:        confidence * 5.0, // Convert to z-score like value
					IP:            ip,
					Timestamp:     time.Now(),
					TimeWindow:    time.Hour,
					Confidence:    confidence,
					Context:       map[string]interface{}{"user_agent": userAgent, "request_count": count},
				}
				anomalies = append(anomalies, anomaly)
			}
		}
	}

	return anomalies
}

// detectGeographicAnomalies detects unusual geographic patterns
func (ad *AnomalyDetector) detectGeographicAnomalies(logs []*parser.LogEntry) []Anomaly {
	var anomalies []Anomaly

	// Simple geographic anomaly detection
	// In production, this would use IP geolocation services
	
	ipCounts := make(map[string]int)
	for _, entry := range logs {
		ipCounts[entry.IP]++
	}

	// Check for IPs with unusually high activity from single source
	totalRequests := len(logs)
	for ip, count := range ipCounts {
		percentage := float64(count) / float64(totalRequests)
		
		if percentage > 0.1 { // Single IP accounts for more than 10% of traffic
			severity := SeverityMedium
			if percentage > 0.25 {
				severity = SeverityHigh
			}
			if percentage > 0.5 {
				severity = SeverityCritical
			}

			anomaly := Anomaly{
				ID:            fmt.Sprintf("geo_%d_%s", time.Now().UnixNano(), ip),
				Type:          AnomalyGeographic,
				Severity:      severity,
				Description:   fmt.Sprintf("Single IP accounts for %.1f%% of total traffic", percentage*100),
				Metric:        "traffic_percentage",
				ExpectedValue: 0.01, // Expected 1% or less per IP
				ActualValue:   percentage,
				Deviation:     percentage - 0.01,
				ZScore:        (percentage - 0.01) / 0.05,
				IP:            ip,
				Timestamp:     time.Now(),
				TimeWindow:    time.Hour,
				Confidence:    math.Min(percentage*5, 1.0),
				Context:       map[string]interface{}{"request_count": count, "total_requests": totalRequests},
			}
			anomalies = append(anomalies, anomaly)
		}
	}

	return anomalies
}

// detectEndpointPatternAnomalies detects unusual endpoint access patterns
func (ad *AnomalyDetector) detectEndpointPatternAnomalies(logs []*parser.LogEntry) []Anomaly {
	var anomalies []Anomaly

	// Group by IP and analyze endpoint access patterns
	ipEndpoints := make(map[string]map[string]int)
	for _, entry := range logs {
		path := strings.Split(entry.URL, "?")[0]
		if ipEndpoints[entry.IP] == nil {
			ipEndpoints[entry.IP] = make(map[string]int)
		}
		ipEndpoints[entry.IP][path]++
	}

	// Analyze patterns
	for ip, endpoints := range ipEndpoints {
		totalRequests := 0
		for _, count := range endpoints {
			totalRequests += count
		}

		if totalRequests < 10 {
			continue
		}

		// Check for excessive endpoint enumeration
		uniqueEndpoints := len(endpoints)
		enumerationRatio := float64(uniqueEndpoints) / float64(totalRequests)

		if enumerationRatio > 0.5 && uniqueEndpoints > 20 {
			severity := SeverityMedium
			if uniqueEndpoints > 100 {
				severity = SeverityHigh
			}

			anomaly := Anomaly{
				ID:            fmt.Sprintf("endpoint_%d_%s", time.Now().UnixNano(), ip),
				Type:          AnomalyEndpointPattern,
				Severity:      severity,
				Description:   fmt.Sprintf("Excessive endpoint enumeration (%d unique endpoints)", uniqueEndpoints),
				Metric:        "unique_endpoints_ratio",
				ExpectedValue: 0.1, // Expected 10% unique endpoints
				ActualValue:   enumerationRatio,
				Deviation:     enumerationRatio - 0.1,
				ZScore:        (enumerationRatio - 0.1) / 0.2,
				IP:            ip,
				Timestamp:     time.Now(),
				TimeWindow:    time.Hour,
				Confidence:    math.Min(enumerationRatio*2, 1.0),
				Context:       map[string]interface{}{"unique_endpoints": uniqueEndpoints, "total_requests": totalRequests},
			}
			anomalies = append(anomalies, anomaly)
		}

		// Check for access to uncommon endpoints
		for endpoint, count := range endpoints {
			if frequency, exists := ad.globalBaseline.CommonEndpoints[endpoint]; exists {
				expectedCount := int(frequency * float64(totalRequests))
				if count > expectedCount*5 { // Much higher than baseline
					anomaly := Anomaly{
						ID:            fmt.Sprintf("endpoint_freq_%d_%s", time.Now().UnixNano(), ip),
						Type:          AnomalyEndpointPattern,
						Severity:      SeverityLow,
						Description:   fmt.Sprintf("Unusual access frequency to endpoint %s", endpoint),
						Metric:        "endpoint_access_frequency",
						ExpectedValue: float64(expectedCount),
						ActualValue:   float64(count),
						Deviation:     float64(count - expectedCount),
						ZScore:        float64(count-expectedCount) / float64(expectedCount+1),
						IP:            ip,
						Timestamp:     time.Now(),
						TimeWindow:    time.Hour,
						Confidence:    0.6,
						Context:       map[string]interface{}{"endpoint": endpoint, "access_count": count},
					}
					anomalies = append(anomalies, anomaly)
				}
			}
		}
	}

	return anomalies
}

// detectStatusAnomalies detects unusual status code patterns
func (ad *AnomalyDetector) detectStatusAnomalies(logs []*parser.LogEntry) []Anomaly {
	var anomalies []Anomaly

	// Group by IP and analyze status code patterns
	ipStatusCodes := make(map[string]map[int]int)
	for _, entry := range logs {
		if ipStatusCodes[entry.IP] == nil {
			ipStatusCodes[entry.IP] = make(map[int]int)
		}
		ipStatusCodes[entry.IP][entry.Status]++
	}

	// Analyze patterns for each IP
	for ip, statusCodes := range ipStatusCodes {
		totalRequests := 0
		for _, count := range statusCodes {
			totalRequests += count
		}

		if totalRequests < 10 {
			continue
		}

		// Check each status code against baseline
		for statusCode, count := range statusCodes {
			if frequency, exists := ad.globalBaseline.CommonStatusCodes[statusCode]; exists {
				actualFreq := float64(count) / float64(totalRequests)
				
				zScore := (actualFreq - frequency) / (frequency + 0.01)
				
				if math.Abs(zScore) > ad.config.AnomalyThreshold {
					severity := SeverityLow
					if statusCode >= 500 && actualFreq > 0.1 {
						severity = SeverityMedium
					}
					if statusCode >= 500 && actualFreq > 0.3 {
						severity = SeverityHigh
					}

					anomaly := Anomaly{
						ID:            fmt.Sprintf("status_%d_%s", time.Now().UnixNano(), ip),
						Type:          AnomalyStatusCodePattern,
						Severity:      severity,
						Description:   fmt.Sprintf("Unusual status code %d frequency (%.1f%%)", statusCode, actualFreq*100),
						Metric:        "status_code_frequency",
						ExpectedValue: frequency,
						ActualValue:   actualFreq,
						Deviation:     math.Abs(actualFreq - frequency),
						ZScore:        zScore,
						IP:            ip,
						Timestamp:     time.Now(),
						TimeWindow:    time.Hour,
						Confidence:    math.Min(math.Abs(zScore)/3.0, 1.0),
						Context:       map[string]interface{}{"status_code": statusCode, "count": count},
					}
					anomalies = append(anomalies, anomaly)
				}
			}
		}
	}

	return anomalies
}

// detectIPSpecificAnomalies detects anomalies specific to an IP's profile
func (ad *AnomalyDetector) detectIPSpecificAnomalies(profile *IPBehaviorProfile, entries []*parser.LogEntry) []Anomaly {
	var anomalies []Anomaly

	// This would contain IP-specific anomaly detection logic
	// For now, return empty slice - to be expanded based on specific requirements
	
	return anomalies
}

// calculateBehaviorScore calculates a behavior score for an IP (0.0-1.0, higher = more suspicious)
func (ad *AnomalyDetector) calculateBehaviorScore(profile *IPBehaviorProfile) float64 {
	score := 0.0

	// High request frequency indicator
	if profile.RequestFrequency > ad.globalBaseline.AverageRequestsPerMinute*5 {
		score += 0.3
	}

	// High error rate indicator
	if profile.ErrorRate > ad.globalBaseline.ErrorRateThreshold*3 {
		score += 0.2
	}

	// Too regular intervals (bot-like)
	if profile.AverageRequestInterval > 0 && profile.AverageRequestInterval < 5*time.Second {
		score += 0.2
	}

	// Many unique endpoints (scanning behavior)
	if len(profile.VisitedEndpoints) > 50 {
		score += 0.1
	}

	// Suspicious user agents
	for userAgent := range profile.CommonUserAgents {
		if strings.Contains(strings.ToLower(userAgent), "bot") || 
		   strings.Contains(strings.ToLower(userAgent), "scanner") {
			score += 0.1
			break
		}
	}

	// Geographic inconsistency
	if !profile.GeographicConsistency {
		score += 0.1
	}

	return math.Min(score, 1.0)
}

// assessRiskLevel determines risk level based on behavior score
func (ad *AnomalyDetector) assessRiskLevel(behaviorScore float64) RiskLevel {
	if behaviorScore >= 0.8 {
		return RiskCritical
	} else if behaviorScore >= 0.6 {
		return RiskHigh
	} else if behaviorScore >= 0.4 {
		return RiskMedium
	} else if behaviorScore >= 0.2 {
		return RiskLow
	}
	return RiskMinimal
}

// Helper function to calculate mean and standard deviation
func calculateStats(data []float64) (mean, stdDev float64) {
	if len(data) == 0 {
		return 0, 0
	}

	// Calculate mean
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	mean = sum / float64(len(data))

	// Calculate standard deviation
	variance := 0.0
	for _, value := range data {
		variance += math.Pow(value-mean, 2)
	}
	variance /= float64(len(data))
	stdDev = math.Sqrt(variance)

	return mean, stdDev
}