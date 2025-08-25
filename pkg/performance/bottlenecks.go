package performance

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// DetectBottlenecks identifies performance bottlenecks in the analysis
func (a *Analyzer) DetectBottlenecks(analysis *PerformanceAnalysis) ([]Bottleneck, error) {
	var bottlenecks []Bottleneck

	// Detect slow endpoints
	slowEndpoints := a.detectSlowEndpoints(analysis.EndpointMetrics)
	bottlenecks = append(bottlenecks, slowEndpoints...)

	// Detect high error rate issues
	errorBottlenecks := a.detectHighErrorRates(analysis.EndpointMetrics)
	bottlenecks = append(bottlenecks, errorBottlenecks...)

	// Detect traffic spikes
	trafficSpikes := a.detectTrafficSpikes(analysis.TimeBasedMetrics)
	bottlenecks = append(bottlenecks, trafficSpikes...)

	// Detect resource exhaustion patterns
	resourceBottlenecks := a.detectResourceExhaustion(analysis.EndpointMetrics)
	bottlenecks = append(bottlenecks, resourceBottlenecks...)

	// Sort bottlenecks by severity (descending)
	sort.Slice(bottlenecks, func(i, j int) bool {
		return bottlenecks[i].Severity > bottlenecks[j].Severity
	})

	return bottlenecks, nil
}

// detectSlowEndpoints identifies endpoints with poor performance
func (a *Analyzer) detectSlowEndpoints(endpointMetrics map[string]*EndpointPerformance) []Bottleneck {
	var bottlenecks []Bottleneck

	// Calculate performance statistics
	var latencies []time.Duration
	for _, metrics := range endpointMetrics {
		latencies = append(latencies, metrics.EstimatedLatency.P95)
	}

	if len(latencies) == 0 {
		return bottlenecks
	}

	// Calculate statistics
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	mean := a.calculateMeanLatency(latencies)
	stdDev := a.calculateLatencyStdDev(latencies, mean)

	// Threshold for slow endpoints: mean + 2*stddev or > 2s, whichever is lower
	threshold := mean + time.Duration(2*float64(stdDev))
	if threshold > 2*time.Second {
		threshold = 2 * time.Second
	}
	if threshold < 500*time.Millisecond {
		threshold = 500 * time.Millisecond
	}

	// Find endpoints exceeding the threshold
	var slowEndpoints []string
	for endpoint, metrics := range endpointMetrics {
		if metrics.EstimatedLatency.P95 > threshold && metrics.RequestCount > 10 {
			slowEndpoints = append(slowEndpoints, endpoint)
		}
	}

	if len(slowEndpoints) > 0 {
		// Calculate severity based on how much they exceed the threshold
		severity := a.calculateSeverity(slowEndpoints, endpointMetrics, threshold)

		bottleneck := Bottleneck{
			Type:       SlowEndpoint,
			Severity:   severity,
			Affected:   slowEndpoints,
			TimeWindow: TimeRange{}, // Will be set if we had actual timestamps
			Description: fmt.Sprintf("Detected %d slow endpoints with P95 latency > %v", 
				len(slowEndpoints), threshold),
			Impact: fmt.Sprintf("These endpoints may be causing poor user experience and should be optimized"),
			Suggestions: []string{
				"Profile these endpoints to identify performance bottlenecks",
				"Consider implementing caching for frequently accessed resources",
				"Optimize database queries and reduce N+1 query problems",
				"Review code efficiency and algorithm complexity",
				"Consider implementing pagination for large result sets",
			},
		}

		bottlenecks = append(bottlenecks, bottleneck)
	}

	return bottlenecks
}

// detectHighErrorRates identifies endpoints with excessive error rates
func (a *Analyzer) detectHighErrorRates(endpointMetrics map[string]*EndpointPerformance) []Bottleneck {
	var bottlenecks []Bottleneck

	var problematicEndpoints []string
	var totalErrorRate float64
	endpointCount := 0

	for endpoint, metrics := range endpointMetrics {
		if metrics.RequestCount < 10 {
			continue // Skip low-traffic endpoints
		}

		totalErrorRate += metrics.ErrorRate
		endpointCount++

		// Flag endpoints with error rate > 5% as problematic
		if metrics.ErrorRate > 0.05 {
			problematicEndpoints = append(problematicEndpoints, endpoint)
		}
	}

	if len(problematicEndpoints) > 0 {
		severity := int(math.Min(10, float64(len(problematicEndpoints))*2))

		bottleneck := Bottleneck{
			Type:       HighErrorRate,
			Severity:   severity,
			Affected:   problematicEndpoints,
			Description: fmt.Sprintf("Detected %d endpoints with high error rates (>5%%)", 
				len(problematicEndpoints)),
			Impact: fmt.Sprintf("High error rates indicate reliability issues and poor user experience"),
			Suggestions: []string{
				"Review application logs for root cause of errors",
				"Implement proper error handling and graceful degradation",
				"Check database connectivity and query performance",
				"Validate input parameters and implement proper validation",
				"Consider implementing circuit breakers for external dependencies",
				"Monitor and alert on error rate spikes",
			},
		}

		bottlenecks = append(bottlenecks, bottleneck)
	}

	return bottlenecks
}

// detectTrafficSpikes identifies unusual traffic patterns
func (a *Analyzer) detectTrafficSpikes(hourlyMetrics []HourlyPerformance) []Bottleneck {
	var bottlenecks []Bottleneck

	if len(hourlyMetrics) == 0 {
		return bottlenecks
	}

	// Calculate traffic statistics
	var requestCounts []int64
	var totalRequests int64
	
	for _, metrics := range hourlyMetrics {
		requestCounts = append(requestCounts, metrics.RequestCount)
		totalRequests += metrics.RequestCount
	}

	if totalRequests == 0 {
		return bottlenecks
	}

	avgRequestsPerHour := float64(totalRequests) / 24.0
	
	// Find hours with traffic > 3x average
	var spikeHours []int
	var maxSpike int64 = 0

	for _, metrics := range hourlyMetrics {
		if float64(metrics.RequestCount) > 3*avgRequestsPerHour && metrics.RequestCount > 100 {
			spikeHours = append(spikeHours, metrics.Hour)
			if metrics.RequestCount > maxSpike {
				maxSpike = metrics.RequestCount
			}
		}
	}

	if len(spikeHours) > 0 {
		severity := int(math.Min(10, float64(maxSpike)/avgRequestsPerHour))
		if severity > 10 {
			severity = 10
		}

		hoursStr := make([]string, len(spikeHours))
		for i, hour := range spikeHours {
			hoursStr[i] = fmt.Sprintf("%02d:00", hour)
		}

		bottleneck := Bottleneck{
			Type:       TrafficSpike,
			Severity:   severity,
			Affected:   hoursStr,
			Description: fmt.Sprintf("Detected traffic spikes during %d hours with >3x average load", 
				len(spikeHours)),
			Impact: fmt.Sprintf("Traffic spikes can overwhelm server capacity and degrade performance"),
			Suggestions: []string{
				"Implement auto-scaling to handle traffic spikes",
				"Set up load balancing across multiple servers",
				"Implement rate limiting to prevent abuse",
				"Use CDN for static content to reduce server load",
				"Consider implementing request queuing during peak times",
				"Monitor and alert on traffic spike patterns",
			},
		}

		bottlenecks = append(bottlenecks, bottleneck)
	}

	return bottlenecks
}

// detectResourceExhaustion identifies patterns suggesting resource constraints
func (a *Analyzer) detectResourceExhaustion(endpointMetrics map[string]*EndpointPerformance) []Bottleneck {
	var bottlenecks []Bottleneck

	// Look for patterns suggesting memory/bandwidth exhaustion
	var largeResponseEndpoints []string
	var highVolumeEndpoints []string

	for endpoint, metrics := range endpointMetrics {
		// Large average response size (>1MB)
		if metrics.AverageSize > 1024*1024 && metrics.RequestCount > 5 {
			largeResponseEndpoints = append(largeResponseEndpoints, endpoint)
		}

		// High volume endpoints (>1000 requests)
		if metrics.RequestCount > 1000 {
			highVolumeEndpoints = append(highVolumeEndpoints, endpoint)
		}
	}

	// Bottleneck for large responses
	if len(largeResponseEndpoints) > 0 {
		severity := int(math.Min(8, float64(len(largeResponseEndpoints))*2))

		bottleneck := Bottleneck{
			Type:       ResourceExhaustion,
			Severity:   severity,
			Affected:   largeResponseEndpoints,
			Description: fmt.Sprintf("Detected %d endpoints serving large responses (>1MB average)", 
				len(largeResponseEndpoints)),
			Impact: fmt.Sprintf("Large responses consume bandwidth and memory, affecting server capacity"),
			Suggestions: []string{
				"Implement response compression (gzip/brotli)",
				"Optimize image and asset sizes",
				"Implement pagination for large data sets",
				"Use streaming responses for large files",
				"Consider implementing response caching",
				"Review if full response data is needed",
			},
		}

		bottlenecks = append(bottlenecks, bottleneck)
	}

	// Bottleneck for high volume endpoints
	if len(highVolumeEndpoints) > 3 {
		severity := int(math.Min(7, float64(len(highVolumeEndpoints))))

		bottleneck := Bottleneck{
			Type:       ResourceExhaustion,
			Severity:   severity,
			Affected:   highVolumeEndpoints,
			Description: fmt.Sprintf("Detected %d high-volume endpoints (>1000 requests each)", 
				len(highVolumeEndpoints)),
			Impact: fmt.Sprintf("High-volume endpoints may strain server resources during peak times"),
			Suggestions: []string{
				"Implement caching for frequently accessed endpoints",
				"Consider using a CDN for static or cacheable content",
				"Optimize database queries and indexes",
				"Implement connection pooling and resource management",
				"Consider horizontal scaling for high-traffic endpoints",
			},
		}

		bottlenecks = append(bottlenecks, bottleneck)
	}

	return bottlenecks
}

// calculateMeanLatency calculates the mean of latency values
func (a *Analyzer) calculateMeanLatency(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	var sum time.Duration
	for _, latency := range latencies {
		sum += latency
	}

	return sum / time.Duration(len(latencies))
}

// calculateLatencyStdDev calculates the standard deviation of latency values
func (a *Analyzer) calculateLatencyStdDev(latencies []time.Duration, mean time.Duration) float64 {
	if len(latencies) <= 1 {
		return 0
	}

	var sumSquaredDiff float64
	for _, latency := range latencies {
		diff := float64(latency - mean)
		sumSquaredDiff += diff * diff
	}

	variance := sumSquaredDiff / float64(len(latencies)-1)
	return math.Sqrt(variance)
}

// calculateSeverity determines bottleneck severity based on affected endpoints
func (a *Analyzer) calculateSeverity(endpoints []string, metrics map[string]*EndpointPerformance, threshold time.Duration) int {
	if len(endpoints) == 0 {
		return 1
	}

	// Base severity on number of affected endpoints
	baseSeverity := int(math.Min(8, float64(len(endpoints))))

	// Increase severity based on how much endpoints exceed threshold
	var excessFactor float64
	for _, endpoint := range endpoints {
		if metric, exists := metrics[endpoint]; exists {
			excess := float64(metric.EstimatedLatency.P95 - threshold) / float64(threshold)
			excessFactor += excess
		}
	}

	avgExcess := excessFactor / float64(len(endpoints))
	severityBonus := int(math.Min(2, avgExcess))

	severity := baseSeverity + severityBonus
	if severity > 10 {
		severity = 10
	}
	if severity < 1 {
		severity = 1
	}

	return severity
}