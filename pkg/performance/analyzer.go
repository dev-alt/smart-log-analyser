package performance

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"smart-log-analyser/pkg/parser"
)

// Analyzer implements the PerformanceAnalyzer interface
type Analyzer struct {
	thresholds PerformanceThresholds
}

// NewAnalyzer creates a new performance analyzer with default thresholds
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		thresholds: DefaultThresholds(),
	}
}

// NewAnalyzerWithThresholds creates a new performance analyzer with custom thresholds
func NewAnalyzerWithThresholds(thresholds PerformanceThresholds) *Analyzer {
	return &Analyzer{
		thresholds: thresholds,
	}
}

// Analyze performs comprehensive performance analysis on log entries
func (a *Analyzer) Analyze(logs []*parser.LogEntry) (*PerformanceAnalysis, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("no log entries provided for analysis")
	}

	analysis := &PerformanceAnalysis{
		EndpointMetrics:      make(map[string]*EndpointPerformance),
		AnalysisTimestamp:    time.Now(),
		TotalEntriesAnalyzed: int64(len(logs)),
	}

	// Set time range
	if len(logs) > 0 {
		analysis.LogTimeRange.Start = logs[0].Timestamp
		analysis.LogTimeRange.End = logs[len(logs)-1].Timestamp
	}

	// Group entries by endpoint for analysis
	endpointEntries := a.groupByEndpoint(logs)
	
	// Analyze each endpoint
	for endpoint, entries := range endpointEntries {
		perfMetrics, err := a.AnalyzeEndpoint(endpoint, entries)
		if err != nil {
			continue // Skip problematic endpoints
		}
		analysis.EndpointMetrics[endpoint] = perfMetrics
	}

	// Generate time-based metrics
	analysis.TimeBasedMetrics = a.generateTimeBasedMetrics(logs)

	// Generate summary
	analysis.Summary = a.generateSummary(analysis, logs)

	// Detect bottlenecks
	bottlenecks, err := a.DetectBottlenecks(analysis)
	if err == nil {
		analysis.Bottlenecks = bottlenecks
	}

	// Generate recommendations
	recommendations, err := a.GenerateRecommendations(analysis)
	if err == nil {
		analysis.Recommendations = recommendations
	}

	// Calculate performance score
	analysis.Score = a.CalculatePerformanceScore(analysis)

	return analysis, nil
}

// AnalyzeEndpoint analyzes performance metrics for a specific endpoint
func (a *Analyzer) AnalyzeEndpoint(endpoint string, entries []*parser.LogEntry) (*EndpointPerformance, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("no entries for endpoint %s", endpoint)
	}

	perf := &EndpointPerformance{
		URL:         endpoint,
		RequestCount: int64(len(entries)),
		StatusCodes: make(map[int]int64),
		Methods:     make(map[string]int64),
	}

	var sizes []int64
	var errorCount int64
	var totalSize int64

	// Process each entry
	for _, entry := range entries {
		// Track sizes
		sizes = append(sizes, entry.Size)
		totalSize += entry.Size

		// Track status codes
		perf.StatusCodes[entry.Status]++

		// Track methods
		perf.Methods[entry.Method]++

		// Count errors
		if entry.Status >= 400 {
			errorCount++
		}
	}

	// Calculate basic metrics
	perf.TotalSize = totalSize
	perf.AverageSize = totalSize / int64(len(entries))
	perf.ErrorRate = float64(errorCount) / float64(len(entries))

	// Calculate estimated latency metrics
	perf.EstimatedLatency = a.calculateLatencyMetrics(entries, sizes)

	// Determine performance grade
	perf.Performance = a.classifyPerformance(perf.EstimatedLatency.P95)

	// Calculate peak throughput (simplified)
	perf.PeakThroughput = a.calculatePeakThroughput(entries)

	return perf, nil
}

// calculateLatencyMetrics estimates latency metrics based on response sizes and patterns
func (a *Analyzer) calculateLatencyMetrics(entries []*parser.LogEntry, sizes []int64) LatencyMetrics {
	if len(sizes) == 0 {
		return LatencyMetrics{}
	}

	// Sort sizes for percentile calculation
	sortedSizes := make([]int64, len(sizes))
	copy(sortedSizes, sizes)
	sort.Slice(sortedSizes, func(i, j int) bool {
		return sortedSizes[i] < sortedSizes[j]
	})

	// Calculate percentile indices
	p50Idx := int(float64(len(sortedSizes)) * 0.5)
	p95Idx := int(float64(len(sortedSizes)) * 0.95)
	p99Idx := int(float64(len(sortedSizes)) * 0.99)

	// Ensure indices are within bounds
	if p50Idx >= len(sortedSizes) {
		p50Idx = len(sortedSizes) - 1
	}
	if p95Idx >= len(sortedSizes) {
		p95Idx = len(sortedSizes) - 1
	}
	if p99Idx >= len(sortedSizes) {
		p99Idx = len(sortedSizes) - 1
	}

	// Estimate latency based on response size (simplified model)
	// Base latency + size-dependent latency + URL complexity
	baseLatency := 50 * time.Millisecond

	return LatencyMetrics{
		P50: a.estimateLatencyFromSize(sortedSizes[p50Idx], entries[0]) + baseLatency,
		P95: a.estimateLatencyFromSize(sortedSizes[p95Idx], entries[0]) + baseLatency,
		P99: a.estimateLatencyFromSize(sortedSizes[p99Idx], entries[0]) + baseLatency,
		Min: a.estimateLatencyFromSize(sortedSizes[0], entries[0]) + baseLatency,
		Max: a.estimateLatencyFromSize(sortedSizes[len(sortedSizes)-1], entries[0]) + baseLatency,
		Avg: a.estimateLatencyFromSize(a.calculateAverage(sortedSizes), entries[0]) + baseLatency,
	}
}

// estimateLatencyFromSize estimates latency based on response size and URL characteristics
func (a *Analyzer) estimateLatencyFromSize(size int64, entry *parser.LogEntry) time.Duration {
	// Base formula: larger responses take longer
	sizeLatency := time.Duration(size/1024) * time.Millisecond // 1ms per KB

	// URL complexity factor
	complexityMultiplier := a.calculateURLComplexity(entry.URL)
	
	// Method factor
	methodMultiplier := a.getMethodMultiplier(entry.Method)

	// Status code factor (errors often take less time)
	statusMultiplier := a.getStatusMultiplier(entry.Status)

	estimatedLatency := time.Duration(float64(sizeLatency) * complexityMultiplier * methodMultiplier * statusMultiplier)

	// Cap at reasonable maximums
	if estimatedLatency > 30*time.Second {
		estimatedLatency = 30 * time.Second
	}
	if estimatedLatency < 10*time.Millisecond {
		estimatedLatency = 10 * time.Millisecond
	}

	return estimatedLatency
}

// calculateURLComplexity estimates URL processing complexity
func (a *Analyzer) calculateURLComplexity(urlPath string) float64 {
	complexity := 1.0

	// Parse URL to check for parameters
	if strings.Contains(urlPath, "?") {
		params := strings.Count(urlPath, "&") + 1
		complexity += float64(params) * 0.1 // Each parameter adds 10% overhead
	}

	// Path depth complexity
	pathDepth := strings.Count(urlPath, "/")
	if pathDepth > 3 {
		complexity += float64(pathDepth-3) * 0.05 // Deep paths add overhead
	}

	// Special endpoints that are typically slower
	lowerURL := strings.ToLower(urlPath)
	if strings.Contains(lowerURL, "search") || 
	   strings.Contains(lowerURL, "report") || 
	   strings.Contains(lowerURL, "export") ||
	   strings.Contains(lowerURL, "admin") {
		complexity *= 1.5
	}

	// API endpoints are often more complex
	if strings.Contains(lowerURL, "/api/") {
		complexity *= 1.2
	}

	return complexity
}

// getMethodMultiplier returns a multiplier based on HTTP method
func (a *Analyzer) getMethodMultiplier(method string) float64 {
	switch method {
	case "GET":
		return 1.0
	case "POST":
		return 1.3 // POST typically involves processing
	case "PUT", "PATCH":
		return 1.4 // Modification operations
	case "DELETE":
		return 1.2
	default:
		return 1.1
	}
}

// getStatusMultiplier returns a multiplier based on status code
func (a *Analyzer) getStatusMultiplier(status int) float64 {
	switch {
	case status >= 200 && status < 300:
		return 1.0 // Normal responses
	case status >= 300 && status < 400:
		return 0.8 // Redirects are typically fast
	case status >= 400 && status < 500:
		return 0.9 // Client errors, often cached responses
	case status >= 500:
		return 0.7 // Server errors often fail fast
	default:
		return 1.0
	}
}

// classifyPerformance determines performance grade based on estimated latency
func (a *Analyzer) classifyPerformance(p95Latency time.Duration) PerformanceGrade {
	switch {
	case p95Latency < a.thresholds.ExcellentLatency:
		return Excellent
	case p95Latency < a.thresholds.GoodLatency:
		return Good
	case p95Latency < a.thresholds.FairLatency:
		return Fair
	case p95Latency < a.thresholds.PoorLatency:
		return Poor
	default:
		return Critical
	}
}

// calculatePeakThroughput estimates peak throughput for an endpoint
func (a *Analyzer) calculatePeakThroughput(entries []*parser.LogEntry) float64 {
	if len(entries) == 0 {
		return 0
	}

	// Group by minute to find peak
	minuteCounts := make(map[string]int)
	
	for _, entry := range entries {
		minute := entry.Timestamp.Truncate(time.Minute).Format("2006-01-02 15:04")
		minuteCounts[minute]++
	}

	// Find maximum requests per minute
	maxPerMinute := 0
	for _, count := range minuteCounts {
		if count > maxPerMinute {
			maxPerMinute = count
		}
	}

	// Convert to requests per second
	return float64(maxPerMinute) / 60.0
}

// generateTimeBasedMetrics creates hourly performance metrics
func (a *Analyzer) generateTimeBasedMetrics(logs []*parser.LogEntry) []HourlyPerformance {
	hourlyData := make(map[int]*HourlyPerformance)

	// Initialize all hours (0-23)
	for i := 0; i < 24; i++ {
		hourlyData[i] = &HourlyPerformance{
			Hour:       i,
			TopErrors:  make([]ErrorSummary, 0),
		}
	}

	// Process each log entry
	for _, entry := range logs {
		hour := entry.Timestamp.Hour()
		hourData := hourlyData[hour]

		hourData.RequestCount++
		hourData.AverageSize = (hourData.AverageSize*int64(hourData.RequestCount-1) + entry.Size) / hourData.RequestCount

		// Count errors
		if entry.Status >= 400 {
			hourData.ErrorRate = (hourData.ErrorRate*float64(hourData.RequestCount-1) + 1.0) / float64(hourData.RequestCount)
		} else {
			hourData.ErrorRate = hourData.ErrorRate * float64(hourData.RequestCount-1) / float64(hourData.RequestCount)
		}
	}

	// Calculate throughput and convert to slice
	result := make([]HourlyPerformance, 0, 24)
	for i := 0; i < 24; i++ {
		hourData := hourlyData[i]
		if hourData.RequestCount > 0 {
			// Simplified throughput calculation (requests per second for that hour)
			hourData.Throughput = float64(hourData.RequestCount) / 3600.0
		}
		result = append(result, *hourData)
	}

	return result
}

// generateSummary creates a performance summary
func (a *Analyzer) generateSummary(analysis *PerformanceAnalysis, logs []*parser.LogEntry) PerformanceSummary {
	summary := PerformanceSummary{
		TotalRequests: int64(len(logs)),
	}

	var totalSize, errorCount int64

	// Calculate aggregate metrics
	for _, entry := range logs {
		totalSize += entry.Size
		if entry.Status >= 400 {
			errorCount++
		}
	}

	summary.AverageResponseSize = totalSize / int64(len(logs))
	summary.ErrorRate = float64(errorCount) / float64(len(logs))

	// Find overall latency metrics
	allSizes := make([]int64, len(logs))
	for i, entry := range logs {
		allSizes[i] = entry.Size
	}

	if len(logs) > 0 {
		summary.OverallLatency = a.calculateLatencyMetrics(logs, allSizes)
		summary.PerformanceGrade = a.classifyPerformance(summary.OverallLatency.P95)
	}

	// Find peak throughput across all endpoints
	for _, metrics := range analysis.EndpointMetrics {
		if metrics.PeakThroughput > summary.PeakThroughput {
			summary.PeakThroughput = metrics.PeakThroughput
		}
	}

	// Find top slow endpoints
	type endpointLatency struct {
		url     string
		latency time.Duration
	}

	var endpointLatencies []endpointLatency
	for url, metrics := range analysis.EndpointMetrics {
		endpointLatencies = append(endpointLatencies, endpointLatency{
			url:     url,
			latency: metrics.EstimatedLatency.P95,
		})
	}

	// Sort by latency (descending)
	sort.Slice(endpointLatencies, func(i, j int) bool {
		return endpointLatencies[i].latency > endpointLatencies[j].latency
	})

	// Take top 5 slowest endpoints
	for i, endpoint := range endpointLatencies {
		if i >= 5 {
			break
		}
		summary.TopSlowEndpoints = append(summary.TopSlowEndpoints, endpoint.url)
	}

	return summary
}

// groupByEndpoint groups log entries by their URL endpoints
func (a *Analyzer) groupByEndpoint(logs []*parser.LogEntry) map[string][]*parser.LogEntry {
	groups := make(map[string][]*parser.LogEntry)

	for _, entry := range logs {
		// Normalize URL by removing query parameters for grouping
		endpoint := a.normalizeEndpoint(entry.URL)
		groups[endpoint] = append(groups[endpoint], entry)
	}

	return groups
}

// normalizeEndpoint normalizes URL paths for grouping
func (a *Analyzer) normalizeEndpoint(urlPath string) string {
	// Remove query parameters
	if idx := strings.Index(urlPath, "?"); idx != -1 {
		urlPath = urlPath[:idx]
	}

	// Remove fragment
	if idx := strings.Index(urlPath, "#"); idx != -1 {
		urlPath = urlPath[:idx]
	}

	// Normalize trailing slashes
	if len(urlPath) > 1 && urlPath[len(urlPath)-1] == '/' {
		urlPath = urlPath[:len(urlPath)-1]
	}

	// Handle empty paths
	if urlPath == "" {
		urlPath = "/"
	}

	return urlPath
}

// calculateAverage calculates the average of int64 slice
func (a *Analyzer) calculateAverage(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}

	var sum int64
	for _, v := range values {
		sum += v
	}

	return sum / int64(len(values))
}