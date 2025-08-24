package trends

import (
	"time"

	"smart-log-analyser/pkg/analyser"
)

// PeriodMetrics contains key metrics for a specific time period
type PeriodMetrics struct {
	Period               string    // Human readable period description
	StartTime            time.Time // Start of period
	EndTime              time.Time // End of period
	TotalRequests        int       // Total number of requests
	AverageResponseSize  int64     // Average response size (proxy for response time)
	ErrorRate            float64   // Percentage of 4xx/5xx responses
	TrafficVolume        int64     // Total bytes transferred
	UniqueVisitors       int       // Unique IP addresses
	PeakHourRequests     int       // Requests during peak hour
	StatusCodeDistrib    map[string]int // Status code distribution
	TopErrorURLs         []analyser.URLStat // URLs with most errors
	BotTrafficPercent    float64   // Percentage of bot traffic
	GeographicDistrib    map[string]int // Country distribution
}

// TrendDirection indicates the direction of change between periods
type TrendDirection int

const (
	TrendStable TrendDirection = iota
	TrendImproving
	TrendDegrading
	TrendCritical // Significant degradation
)

// TrendChange represents a change in a metric between two periods
type TrendChange struct {
	MetricName     string         // Name of the metric
	OldValue       float64        // Value in previous period
	NewValue       float64        // Value in current period
	AbsoluteChange float64        // Absolute difference
	PercentChange  float64        // Percentage change
	Direction      TrendDirection // Direction of change
	Significance   string         // "low", "medium", "high"
	Description    string         // Human readable description
}

// DegradationAlert represents a detected performance degradation
type DegradationAlert struct {
	AlertID      string         // Unique alert identifier
	Severity     string         // "warning", "error", "critical"
	MetricName   string         // Affected metric
	CurrentValue float64        // Current metric value
	BaselineValue float64       // Expected/baseline value
	Threshold    float64        // Threshold that was exceeded
	Impact       string         // Description of potential impact
	Recommendation string       // Suggested action
	DetectedAt   time.Time      // When the degradation was detected
	Trend        TrendDirection // Overall trend direction
}

// PeriodComparison contains the results of comparing two time periods
type PeriodComparison struct {
	BaselinePeriod PeriodMetrics   // Earlier/baseline period
	CurrentPeriod  PeriodMetrics   // Later/current period
	TrendChanges   []TrendChange   // Changes in metrics
	OverallTrend   TrendDirection  // Overall trend direction
	RiskScore      int             // Risk score (0-100, higher is worse)
	Summary        string          // Human readable summary
}

// TrendAnalysis contains comprehensive trend analysis results
type TrendAnalysis struct {
	AnalysisType      string               // "comparison", "degradation", "historical"
	GeneratedAt       time.Time            // When analysis was performed
	PeriodComparisons []PeriodComparison   // Period-to-period comparisons
	DegradationAlerts []DegradationAlert   // Detected degradation issues
	OverallHealth     string               // "healthy", "warning", "critical"
	Recommendations   []string             // Actionable recommendations
	TrendSummary      string               // Executive summary of trends
}

// TrendConfiguration defines parameters for trend analysis
type TrendConfiguration struct {
	// Degradation thresholds
	ErrorRateThreshold      float64 // Error rate increase threshold (%)
	ResponseTimeThreshold   float64 // Response time increase threshold (%)
	TrafficDropThreshold    float64 // Traffic drop threshold (%)
	
	// Statistical parameters
	MinimumSampleSize       int     // Minimum requests needed for analysis
	SignificanceLevel       float64 // Statistical significance level
	
	// Period definitions
	DefaultComparisonPeriod string  // Default period to compare against
	
	// Alert settings
	EnableAlerts            bool    // Whether to generate alerts
	AlertCooldownHours      int     // Hours between similar alerts
}

// DefaultTrendConfiguration returns sensible default configuration
func DefaultTrendConfiguration() TrendConfiguration {
	return TrendConfiguration{
		ErrorRateThreshold:      10.0, // 10% increase triggers alert
		ResponseTimeThreshold:   20.0, // 20% increase triggers alert
		TrafficDropThreshold:    30.0, // 30% drop triggers alert
		MinimumSampleSize:       100,  // Need at least 100 requests
		SignificanceLevel:       0.05, // 95% confidence level
		DefaultComparisonPeriod: "previous-day",
		EnableAlerts:            true,
		AlertCooldownHours:      4, // 4 hours between similar alerts
	}
}

// String methods for enum types
func (td TrendDirection) String() string {
	switch td {
	case TrendStable:
		return "stable"
	case TrendImproving:
		return "improving"
	case TrendDegrading:
		return "degrading"
	case TrendCritical:
		return "critical"
	default:
		return "unknown"
	}
}