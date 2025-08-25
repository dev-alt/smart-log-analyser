package performance

import (
	"math"
	"time"
)

// CalculatePerformanceScore calculates an overall performance score based on analysis
func (a *Analyzer) CalculatePerformanceScore(analysis *PerformanceAnalysis) PerformanceScore {
	// Calculate individual component scores
	latencyScore := a.calculateLatencyScore(analysis.Summary.OverallLatency)
	throughputScore := a.calculateThroughputScore(analysis.Summary.PeakThroughput, analysis.Summary.TotalRequests)
	reliabilityScore := a.calculateReliabilityScore(analysis.Summary.ErrorRate, analysis.EndpointMetrics)
	efficiencyScore := a.calculateEfficiencyScore(analysis.Summary.AverageResponseSize, analysis.EndpointMetrics)

	// Weighted overall score
	// Latency: 35%, Reliability: 30%, Throughput: 20%, Efficiency: 15%
	overall := int(float64(latencyScore)*0.35 + 
		float64(reliabilityScore)*0.30 + 
		float64(throughputScore)*0.20 + 
		float64(efficiencyScore)*0.15)

	// Ensure score is within bounds
	if overall > 100 {
		overall = 100
	}
	if overall < 0 {
		overall = 0
	}

	return PerformanceScore{
		Overall:     overall,
		Latency:     latencyScore,
		Throughput:  throughputScore,
		Reliability: reliabilityScore,
		Efficiency:  efficiencyScore,
	}
}

// calculateLatencyScore scores based on response time performance
func (a *Analyzer) calculateLatencyScore(latency LatencyMetrics) int {
	if latency.P95 == 0 {
		return 50 // Neutral score if no data
	}

	// Score based on P95 latency using exponential decay
	p95Ms := float64(latency.P95) / float64(time.Millisecond)

	var score float64
	switch {
	case p95Ms <= 100:
		score = 100 // Excellent
	case p95Ms <= 200:
		score = 90 - (p95Ms-100)*0.2 // 90-70
	case p95Ms <= 500:
		score = 70 - (p95Ms-200)*0.1 // 70-40
	case p95Ms <= 1000:
		score = 40 - (p95Ms-500)*0.04 // 40-20
	case p95Ms <= 5000:
		score = 20 - (p95Ms-1000)*0.004 // 20-4
	default:
		score = 1 // Critical performance
	}

	// Adjust based on P99 latency (adds penalty for tail latency)
	p99Ms := float64(latency.P99) / float64(time.Millisecond)
	if p99Ms > p95Ms*2 {
		// High tail latency penalty
		penalty := math.Min(20, (p99Ms-p95Ms*2)/100)
		score -= penalty
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return int(score)
}

// calculateThroughputScore scores based on request handling capacity
func (a *Analyzer) calculateThroughputScore(peakThroughput float64, totalRequests int64) int {
	if peakThroughput == 0 || totalRequests == 0 {
		return 50 // Neutral score
	}

	// Base score on peak throughput (requests per second)
	var score float64
	switch {
	case peakThroughput >= 100:
		score = 100 // Excellent throughput
	case peakThroughput >= 50:
		score = 80 + (peakThroughput-50)*0.4 // 80-100
	case peakThroughput >= 20:
		score = 60 + (peakThroughput-20)*0.67 // 60-80
	case peakThroughput >= 10:
		score = 40 + (peakThroughput-10)*2 // 40-60
	case peakThroughput >= 5:
		score = 20 + (peakThroughput-5)*4 // 20-40
	case peakThroughput >= 1:
		score = 10 + (peakThroughput-1)*2.5 // 10-20
	default:
		score = peakThroughput * 10 // < 1 RPS
	}

	// Bonus for handling high request volumes
	if totalRequests > 10000 {
		bonus := math.Min(10, float64(totalRequests)/10000)
		score += bonus
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return int(score)
}

// calculateReliabilityScore scores based on error rates and stability
func (a *Analyzer) calculateReliabilityScore(overallErrorRate float64, endpointMetrics map[string]*EndpointPerformance) int {
	// Base score on overall error rate
	var score float64 = 100

	// Penalize based on error rate
	errorRatePercent := overallErrorRate * 100
	switch {
	case errorRatePercent <= 0.1:
		score = 100 // Excellent reliability
	case errorRatePercent <= 0.5:
		score = 95 - (errorRatePercent-0.1)*25 // 95-85
	case errorRatePercent <= 1.0:
		score = 85 - (errorRatePercent-0.5)*20 // 85-75
	case errorRatePercent <= 2.0:
		score = 75 - (errorRatePercent-1.0)*15 // 75-60
	case errorRatePercent <= 5.0:
		score = 60 - (errorRatePercent-2.0)*10 // 60-30
	case errorRatePercent <= 10.0:
		score = 30 - (errorRatePercent-5.0)*4 // 30-10
	default:
		score = 10 - math.Min(10, errorRatePercent-10) // <10
	}

	// Additional penalty for endpoints with very high error rates
	var criticalEndpoints int
	for _, metrics := range endpointMetrics {
		if metrics.ErrorRate > 0.10 { // >10% error rate
			criticalEndpoints++
		}
	}

	if criticalEndpoints > 0 {
		criticalPenalty := math.Min(20, float64(criticalEndpoints)*5)
		score -= criticalPenalty
	}

	// Bonus for consistent performance across endpoints
	if len(endpointMetrics) > 0 {
		var consistentEndpoints int
		for _, metrics := range endpointMetrics {
			if metrics.ErrorRate < 0.01 { // <1% error rate
				consistentEndpoints++
			}
		}

		consistency := float64(consistentEndpoints) / float64(len(endpointMetrics))
		if consistency > 0.8 { // >80% of endpoints are reliable
			consistencyBonus := (consistency - 0.8) * 25 // Up to 5 point bonus
			score += consistencyBonus
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return int(score)
}

// calculateEfficiencyScore scores based on resource utilization efficiency
func (a *Analyzer) calculateEfficiencyScore(avgResponseSize int64, endpointMetrics map[string]*EndpointPerformance) int {
	var score float64 = 100

	// Penalize large average response sizes
	avgSizeKB := float64(avgResponseSize) / 1024
	switch {
	case avgSizeKB <= 10:
		score = 100 // Very efficient
	case avgSizeKB <= 50:
		score = 95 - (avgSizeKB-10)*0.5 // 95-75
	case avgSizeKB <= 100:
		score = 75 - (avgSizeKB-50)*0.4 // 75-55
	case avgSizeKB <= 500:
		score = 55 - (avgSizeKB-100)*0.1 // 55-15
	case avgSizeKB <= 1000:
		score = 15 - (avgSizeKB-500)*0.02 // 15-5
	default:
		score = 5 // Very inefficient
	}

	// Analyze response size distribution
	if len(endpointMetrics) > 0 {
		var efficientEndpoints int
		var inefficientEndpoints int

		for _, metrics := range endpointMetrics {
			avgKB := float64(metrics.AverageSize) / 1024
			if avgKB <= 50 { // Efficient endpoints
				efficientEndpoints++
			} else if avgKB > 500 { // Inefficient endpoints
				inefficientEndpoints++
			}
		}

		// Bonus for having many efficient endpoints
		efficiencyRatio := float64(efficientEndpoints) / float64(len(endpointMetrics))
		if efficiencyRatio > 0.7 {
			bonus := (efficiencyRatio - 0.7) * 20 // Up to 6 point bonus
			score += bonus
		}

		// Penalty for having many inefficient endpoints
		inefficiencyRatio := float64(inefficientEndpoints) / float64(len(endpointMetrics))
		if inefficiencyRatio > 0.2 {
			penalty := (inefficiencyRatio - 0.2) * 25 // Up to 20 point penalty
			score -= penalty
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return int(score)
}

// GetScoreGrade returns a letter grade for a numeric score
func GetScoreGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

// GetScoreDescription returns a human-readable description of the score
func GetScoreDescription(score int) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 80:
		return "Good"
	case score >= 70:
		return "Fair"
	case score >= 60:
		return "Poor"
	default:
		return "Critical"
	}
}