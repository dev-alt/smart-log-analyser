package performance

import (
	"fmt"
	"sort"
	"strings"
)

// GenerateRecommendations creates optimization recommendations based on performance analysis
func (a *Analyzer) GenerateRecommendations(analysis *PerformanceAnalysis) ([]OptimizationRecommendation, error) {
	var recommendations []OptimizationRecommendation

	// Generate recommendations based on performance score
	scoreRecommendations := a.generateScoreBasedRecommendations(analysis.Score)
	recommendations = append(recommendations, scoreRecommendations...)

	// Generate recommendations based on bottlenecks
	bottleneckRecommendations := a.generateBottleneckRecommendations(analysis.Bottlenecks)
	recommendations = append(recommendations, bottleneckRecommendations...)

	// Generate recommendations based on endpoint analysis
	endpointRecommendations := a.generateEndpointRecommendations(analysis.EndpointMetrics)
	recommendations = append(recommendations, endpointRecommendations...)

	// Generate recommendations based on traffic patterns
	trafficRecommendations := a.generateTrafficRecommendations(analysis.TimeBasedMetrics)
	recommendations = append(recommendations, trafficRecommendations...)

	// Sort recommendations by priority (descending)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Priority > recommendations[j].Priority
	})

	// Remove duplicates and limit to top 10
	recommendations = a.deduplicateRecommendations(recommendations)
	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	return recommendations, nil
}

// generateScoreBasedRecommendations creates recommendations based on overall performance score
func (a *Analyzer) generateScoreBasedRecommendations(score PerformanceScore) []OptimizationRecommendation {
	var recommendations []OptimizationRecommendation

	// Overall performance recommendations
	if score.Overall < 70 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Category:    CodeOptimization,
			Priority:    9,
			Impact:      HighImpact,
			Effort:      MediumEffort,
			Title:       "Critical Performance Issues Detected",
			Description: "Overall performance score is below acceptable levels",
			Details:     "Your application's performance score is critically low, indicating significant user experience issues.",
			Examples: []string{
				"Implement comprehensive performance monitoring",
				"Conduct thorough code review and optimization",
				"Consider performance testing and load testing",
			},
			EstimatedImprovementPercent: 40,
		})
	}

	// Latency-specific recommendations
	if score.Latency < 60 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Category:    CachingOptimization,
			Priority:    8,
			Impact:      HighImpact,
			Effort:      MediumEffort,
			Title:       "High Latency Issues",
			Description: "Response times are significantly impacting user experience",
			Details:     "Implement caching strategies to reduce response times and improve user satisfaction.",
			Examples: []string{
				"Implement Redis or Memcached for frequent queries",
				"Add browser caching headers for static assets",
				"Implement database query result caching",
				"Use CDN for static content delivery",
			},
			EstimatedImprovementPercent: 35,
		})
	}

	// Throughput recommendations
	if score.Throughput < 60 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Category:    InfrastructureScaling,
			Priority:    7,
			Impact:      MediumImpact,
			Effort:      HighEffort,
			Title:       "Throughput Capacity Issues",
			Description: "System struggling to handle request volume efficiently",
			Details:     "Consider scaling infrastructure to handle higher request volumes.",
			Examples: []string{
				"Implement horizontal scaling with load balancers",
				"Optimize server configuration and resource allocation",
				"Consider upgrading server hardware or cloud instance types",
				"Implement request queuing for peak load handling",
			},
			EstimatedImprovementPercent: 30,
		})
	}

	// Reliability recommendations
	if score.Reliability < 70 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Category:    ErrorReduction,
			Priority:    8,
			Impact:      HighImpact,
			Effort:      MediumEffort,
			Title:       "Reliability Issues Detected",
			Description: "High error rates are affecting system reliability",
			Details:     "Focus on error reduction and implementing proper error handling mechanisms.",
			Examples: []string{
				"Implement comprehensive error logging and monitoring",
				"Add graceful error handling and fallback mechanisms",
				"Implement health checks and circuit breakers",
				"Review and fix common error patterns",
			},
			EstimatedImprovementPercent: 25,
		})
	}

	return recommendations
}

// generateBottleneckRecommendations creates recommendations based on detected bottlenecks
func (a *Analyzer) generateBottleneckRecommendations(bottlenecks []Bottleneck) []OptimizationRecommendation {
	var recommendations []OptimizationRecommendation

	for _, bottleneck := range bottlenecks {
		switch bottleneck.Type {
		case SlowEndpoint:
			recommendations = append(recommendations, OptimizationRecommendation{
				Category:    CodeOptimization,
				Priority:    bottleneck.Severity,
				Impact:      a.mapSeverityToImpact(bottleneck.Severity),
				Effort:      MediumEffort,
				Title:       "Optimize Slow Endpoints",
				Description: fmt.Sprintf("Detected %d slow endpoints requiring optimization", len(bottleneck.Affected)),
				Details:     "Focus on optimizing the slowest performing endpoints first for maximum impact.",
				Examples: []string{
					"Profile slow endpoints: " + strings.Join(bottleneck.Affected[:min(3, len(bottleneck.Affected))], ", "),
					"Optimize database queries and reduce N+1 problems",
					"Implement endpoint-specific caching",
					"Review algorithm efficiency and data structures",
				},
				EstimatedImprovementPercent: min(50, bottleneck.Severity*5),
			})

		case HighErrorRate:
			recommendations = append(recommendations, OptimizationRecommendation{
				Category:    ErrorReduction,
				Priority:    bottleneck.Severity,
				Impact:      HighImpact,
				Effort:      MediumEffort,
				Title:       "Reduce Error Rates",
				Description: fmt.Sprintf("Address high error rates affecting %d endpoints", len(bottleneck.Affected)),
				Details:     "High error rates significantly impact user experience and should be prioritized.",
				Examples: []string{
					"Review error logs for endpoints: " + strings.Join(bottleneck.Affected[:min(3, len(bottleneck.Affected))], ", "),
					"Implement proper input validation",
					"Add error handling and graceful degradation",
					"Monitor dependencies and external services",
				},
				EstimatedImprovementPercent: min(40, bottleneck.Severity*4),
			})

		case TrafficSpike:
			recommendations = append(recommendations, OptimizationRecommendation{
				Category:    InfrastructureScaling,
				Priority:    bottleneck.Severity,
				Impact:      MediumImpact,
				Effort:      HighEffort,
				Title:       "Handle Traffic Spikes",
				Description: fmt.Sprintf("Prepare for traffic spikes during %d peak hours", len(bottleneck.Affected)),
				Details:     "Implement scaling strategies to handle traffic spikes without performance degradation.",
				Examples: []string{
					"Peak hours: " + strings.Join(bottleneck.Affected, ", "),
					"Implement auto-scaling policies",
					"Use load balancing for traffic distribution",
					"Consider CDN for static content during peaks",
				},
				EstimatedImprovementPercent: min(35, bottleneck.Severity*3),
			})

		case ResourceExhaustion:
			recommendations = append(recommendations, OptimizationRecommendation{
				Category:    StaticAssetOptimization,
				Priority:    bottleneck.Severity,
				Impact:      MediumImpact,
				Effort:      LowEffort,
				Title:       "Optimize Resource Usage",
				Description: fmt.Sprintf("Optimize resource-heavy endpoints affecting %d URLs", len(bottleneck.Affected)),
				Details:     "Large responses and high-volume endpoints are consuming excessive resources.",
				Examples: []string{
					"Optimize: " + strings.Join(bottleneck.Affected[:min(3, len(bottleneck.Affected))], ", "),
					"Implement response compression",
					"Optimize image and asset sizes",
					"Implement pagination for large datasets",
				},
				EstimatedImprovementPercent: min(30, bottleneck.Severity*3),
			})
		}
	}

	return recommendations
}

// generateEndpointRecommendations creates recommendations based on endpoint performance
func (a *Analyzer) generateEndpointRecommendations(endpointMetrics map[string]*EndpointPerformance) []OptimizationRecommendation {
	var recommendations []OptimizationRecommendation

	// Find endpoints that would benefit from caching
	var cachingCandidates []string
	var compressionCandidates []string

	for endpoint, metrics := range endpointMetrics {
		// High-traffic endpoints are good caching candidates
		if metrics.RequestCount > 100 && metrics.Performance != Excellent {
			cachingCandidates = append(cachingCandidates, endpoint)
		}

		// Large response endpoints are good compression candidates
		if metrics.AverageSize > 50*1024 { // > 50KB
			compressionCandidates = append(compressionCandidates, endpoint)
		}
	}

	// Caching recommendations
	if len(cachingCandidates) > 0 {
		priority := min(8, len(cachingCandidates))
		recommendations = append(recommendations, OptimizationRecommendation{
			Category:    CachingOptimization,
			Priority:    priority,
			Impact:      HighImpact,
			Effort:      MediumEffort,
			Title:       "Implement Caching Strategy",
			Description: fmt.Sprintf("Implement caching for %d high-traffic endpoints", len(cachingCandidates)),
			Details:     "High-traffic endpoints would benefit significantly from caching implementation.",
			Examples: []string{
				"Cache candidates: " + strings.Join(cachingCandidates[:min(3, len(cachingCandidates))], ", "),
				"Implement in-memory caching (Redis/Memcached)",
				"Add HTTP caching headers",
				"Consider database query result caching",
			},
			EstimatedImprovementPercent: min(45, len(cachingCandidates)*5),
		})
	}

	// Compression recommendations
	if len(compressionCandidates) > 0 {
		priority := min(7, len(compressionCandidates))
		recommendations = append(recommendations, OptimizationRecommendation{
			Category:    StaticAssetOptimization,
			Priority:    priority,
			Impact:      MediumImpact,
			Effort:      LowEffort,
			Title:       "Enable Response Compression",
			Description: fmt.Sprintf("Enable compression for %d endpoints with large responses", len(compressionCandidates)),
			Details:     "Large responses can be significantly reduced with proper compression.",
			Examples: []string{
				"Compression targets: " + strings.Join(compressionCandidates[:min(3, len(compressionCandidates))], ", "),
				"Enable gzip/brotli compression in web server",
				"Optimize image formats and sizes",
				"Minify CSS, JavaScript, and HTML",
			},
			EstimatedImprovementPercent: min(25, len(compressionCandidates)*3),
		})
	}

	return recommendations
}

// generateTrafficRecommendations creates recommendations based on traffic patterns
func (a *Analyzer) generateTrafficRecommendations(timeMetrics []HourlyPerformance) []OptimizationRecommendation {
	var recommendations []OptimizationRecommendation

	if len(timeMetrics) == 0 {
		return recommendations
	}

	// Find peak hours
	var peakHours []int
	var avgThroughput float64
	var totalThroughput float64

	for _, metrics := range timeMetrics {
		totalThroughput += metrics.Throughput
	}
	avgThroughput = totalThroughput / float64(len(timeMetrics))

	for _, metrics := range timeMetrics {
		if metrics.Throughput > avgThroughput*2 {
			peakHours = append(peakHours, metrics.Hour)
		}
	}

	// Recommendations for peak hours
	if len(peakHours) > 0 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Category:    InfrastructureScaling,
			Priority:    6,
			Impact:      MediumImpact,
			Effort:      MediumEffort,
			Title:       "Optimize for Peak Traffic Hours",
			Description: fmt.Sprintf("Prepare infrastructure for %d peak traffic hours", len(peakHours)),
			Details:     "Identified specific hours with significantly higher traffic that may require special handling.",
			Examples: []string{
				fmt.Sprintf("Peak hours: %v", formatHours(peakHours)),
				"Implement scheduled auto-scaling",
				"Pre-warm caches before peak hours",
				"Consider traffic shaping during peaks",
			},
			EstimatedImprovementPercent: 20,
		})
	}

	return recommendations
}

// deduplicateRecommendations removes duplicate recommendations
func (a *Analyzer) deduplicateRecommendations(recommendations []OptimizationRecommendation) []OptimizationRecommendation {
	seen := make(map[string]bool)
	var result []OptimizationRecommendation

	for _, rec := range recommendations {
		key := fmt.Sprintf("%s-%s", rec.Category.String(), rec.Title)
		if !seen[key] {
			seen[key] = true
			result = append(result, rec)
		}
	}

	return result
}

// mapSeverityToImpact converts bottleneck severity to impact level
func (a *Analyzer) mapSeverityToImpact(severity int) ImpactLevel {
	switch {
	case severity >= 8:
		return CriticalImpact
	case severity >= 6:
		return HighImpact
	case severity >= 4:
		return MediumImpact
	default:
		return LowImpact
	}
}

// formatHours formats hour numbers for display
func formatHours(hours []int) string {
	var formatted []string
	for _, hour := range hours {
		formatted = append(formatted, fmt.Sprintf("%02d:00", hour))
	}
	return strings.Join(formatted, ", ")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}