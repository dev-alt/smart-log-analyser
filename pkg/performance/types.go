package performance

import (
	"time"

	"smart-log-analyser/pkg/parser"
)

// PerformanceGrade represents the performance classification of an endpoint
type PerformanceGrade int

const (
	Excellent PerformanceGrade = iota // < 100ms estimated
	Good                              // < 500ms estimated  
	Fair                              // < 1s estimated
	Poor                              // < 5s estimated
	Critical                          // >= 5s estimated
)

// String returns the string representation of PerformanceGrade
func (pg PerformanceGrade) String() string {
	switch pg {
	case Excellent:
		return "Excellent"
	case Good:
		return "Good"
	case Fair:
		return "Fair"
	case Poor:
		return "Poor"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// BottleneckType represents different types of performance bottlenecks
type BottleneckType int

const (
	SlowEndpoint BottleneckType = iota
	HighErrorRate
	TrafficSpike
	ResourceExhaustion
	DatabaseBottleneck
	NetworkLatency
)

// String returns the string representation of BottleneckType
func (bt BottleneckType) String() string {
	switch bt {
	case SlowEndpoint:
		return "Slow Endpoint"
	case HighErrorRate:
		return "High Error Rate"
	case TrafficSpike:
		return "Traffic Spike"
	case ResourceExhaustion:
		return "Resource Exhaustion"
	case DatabaseBottleneck:
		return "Database Bottleneck"
	case NetworkLatency:
		return "Network Latency"
	default:
		return "Unknown"
	}
}

// OptimizationCategory represents different optimization areas
type OptimizationCategory int

const (
	CachingOptimization OptimizationCategory = iota
	DatabaseOptimization
	StaticAssetOptimization
	CodeOptimization
	InfrastructureScaling
	ContentDelivery
	ErrorReduction
)

// String returns the string representation of OptimizationCategory
func (oc OptimizationCategory) String() string {
	switch oc {
	case CachingOptimization:
		return "Caching Optimization"
	case DatabaseOptimization:
		return "Database Optimization"
	case StaticAssetOptimization:
		return "Static Asset Optimization"
	case CodeOptimization:
		return "Code Optimization"
	case InfrastructureScaling:
		return "Infrastructure Scaling"
	case ContentDelivery:
		return "Content Delivery"
	case ErrorReduction:
		return "Error Reduction"
	default:
		return "Unknown"
	}
}

// ImpactLevel represents the expected impact of an optimization
type ImpactLevel int

const (
	LowImpact ImpactLevel = iota
	MediumImpact
	HighImpact
	CriticalImpact
)

// String returns the string representation of ImpactLevel
func (il ImpactLevel) String() string {
	switch il {
	case LowImpact:
		return "Low"
	case MediumImpact:
		return "Medium"
	case HighImpact:
		return "High"
	case CriticalImpact:
		return "Critical"
	default:
		return "Unknown"
	}
}

// EffortLevel represents the implementation effort required
type EffortLevel int

const (
	LowEffort EffortLevel = iota
	MediumEffort
	HighEffort
)

// String returns the string representation of EffortLevel
func (el EffortLevel) String() string {
	switch el {
	case LowEffort:
		return "Low"
	case MediumEffort:
		return "Medium"
	case HighEffort:
		return "High"
	default:
		return "Unknown"
	}
}

// TimeRange represents a time window for analysis
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// LatencyMetrics represents latency percentile data
type LatencyMetrics struct {
	P50 time.Duration // Median latency
	P95 time.Duration // 95th percentile
	P99 time.Duration // 99th percentile
	Min time.Duration // Minimum latency
	Max time.Duration // Maximum latency
	Avg time.Duration // Average latency
}

// PerformanceScore represents overall performance metrics
type PerformanceScore struct {
	Overall     int // 0-100 (higher is better)
	Latency     int // Response time score
	Throughput  int // Request volume handling  
	Reliability int // Error rate and stability
	Efficiency  int // Resource utilization
}

// EndpointPerformance represents performance metrics for a specific endpoint
type EndpointPerformance struct {
	URL              string
	RequestCount     int64
	AverageSize      int64
	TotalSize        int64
	EstimatedLatency LatencyMetrics
	ErrorRate        float64
	PeakThroughput   float64
	Performance      PerformanceGrade
	StatusCodes      map[int]int64
	Methods          map[string]int64
}

// HourlyPerformance represents performance metrics for a specific hour
type HourlyPerformance struct {
	Hour           int
	RequestCount   int64
	AverageSize    int64
	ErrorRate      float64
	Throughput     float64
	LatencyMetrics LatencyMetrics
	TopErrors      []ErrorSummary
}

// ErrorSummary represents error information
type ErrorSummary struct {
	Status int
	Count  int64
	URL    string
}

// Bottleneck represents a performance bottleneck
type Bottleneck struct {
	Type        BottleneckType
	Severity    int           // 1-10 scale
	Affected    []string      // URLs or IPs affected
	TimeWindow  TimeRange     // When the bottleneck occurred
	Description string        // Human-readable description
	Impact      string        // Impact assessment
	Suggestions []string      // Optimization recommendations
}

// OptimizationRecommendation represents a performance improvement suggestion
type OptimizationRecommendation struct {
	Category    OptimizationCategory
	Priority    int         // 1-10
	Impact      ImpactLevel
	Effort      EffortLevel
	Title       string
	Description string
	Details     string
	Examples    []string
	EstimatedImprovementPercent int
}

// PerformanceSummary provides high-level performance overview
type PerformanceSummary struct {
	TotalRequests       int64
	AverageResponseSize int64
	OverallLatency      LatencyMetrics
	ErrorRate           float64
	PeakThroughput      float64
	PerformanceGrade    PerformanceGrade
	Score               PerformanceScore
	TopSlowEndpoints    []string
	CriticalIssues      int
	Recommendations     int
}

// PerformanceAnalysis represents the complete performance analysis result
type PerformanceAnalysis struct {
	Summary             PerformanceSummary
	EndpointMetrics     map[string]*EndpointPerformance
	TimeBasedMetrics    []HourlyPerformance
	Bottlenecks         []Bottleneck
	Recommendations     []OptimizationRecommendation
	Score               PerformanceScore
	AnalysisTimestamp   time.Time
	LogTimeRange        TimeRange
	TotalEntriesAnalyzed int64
}

// RequestContext provides context for performance analysis
type RequestContext struct {
	ConcurrentRequests map[time.Time]int
	BaselineLatency    time.Duration
	PeakHours          []int
	TrafficPatterns    map[int]float64 // Hour -> multiplier
}

// PerformanceThresholds defines performance evaluation criteria
type PerformanceThresholds struct {
	ExcellentLatency time.Duration // Default: 100ms
	GoodLatency      time.Duration // Default: 500ms
	FairLatency      time.Duration // Default: 1s
	PoorLatency      time.Duration // Default: 5s
	MaxErrorRate     float64       // Default: 0.1%
	MinThroughput    float64       // Default: 10 req/s
}

// DefaultThresholds returns standard performance thresholds
func DefaultThresholds() PerformanceThresholds {
	return PerformanceThresholds{
		ExcellentLatency: 100 * time.Millisecond,
		GoodLatency:      500 * time.Millisecond,
		FairLatency:      1 * time.Second,
		PoorLatency:      5 * time.Second,
		MaxErrorRate:     0.001, // 0.1%
		MinThroughput:    10.0,  // 10 req/s
	}
}

// PerformanceAnalyzer interface defines the main analysis capabilities
type PerformanceAnalyzer interface {
	Analyze(logs []*parser.LogEntry) (*PerformanceAnalysis, error)
	AnalyzeEndpoint(url string, entries []*parser.LogEntry) (*EndpointPerformance, error)
	DetectBottlenecks(analysis *PerformanceAnalysis) ([]Bottleneck, error)
	GenerateRecommendations(analysis *PerformanceAnalysis) ([]OptimizationRecommendation, error)
	CalculatePerformanceScore(analysis *PerformanceAnalysis) PerformanceScore
}