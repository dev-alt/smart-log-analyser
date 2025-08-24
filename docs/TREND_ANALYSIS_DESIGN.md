# Historical Trend Analysis Design Document

## Overview
Implementation design for historical trend analysis and degradation detection capabilities in the Smart Log Analyser. This feature enables comparison of different time periods to identify performance trends and automatically detect system degradation patterns.

## Architecture Components

### Core Package Structure
```
pkg/trends/
├── types.go          # Data structures and enums
├── analyser.go       # Core trend analysis engine
└── visualization.go  # ASCII chart visualizations
```

## Data Structures

### TrendDirection Enum
- `TrendStable`: No significant change detected
- `TrendImproving`: Metrics showing positive improvement  
- `TrendDegrading`: Metrics showing concerning decline
- `TrendCritical`: Metrics showing severe degradation requiring immediate attention

### Core Analysis Types

#### PeriodMetrics
Contains key performance metrics for a specific time period:
- Request volume and traffic patterns
- Error rates and response times (via response size proxy)
- Bot traffic percentage and geographic distribution
- Status code distribution and top error URLs

#### TrendChange
Represents a change in a specific metric between periods:
- Absolute and percentage change calculations
- Statistical significance assessment (low/medium/high)
- Trend direction classification
- Human-readable descriptions

#### DegradationAlert  
Alert structure for detected performance issues:
- Severity levels (warning/error/critical)
- Impact descriptions and actionable recommendations
- Threshold-based triggering with configurable parameters
- Unique alert IDs for tracking and deduplication

## Analysis Algorithms

### Period Comparison Engine
1. **Data Segmentation**: Splits log data into comparison periods
2. **Metric Extraction**: Converts raw logs to structured metrics using existing analyser
3. **Statistical Analysis**: Calculates percentage changes and significance levels
4. **Trend Classification**: Determines overall trend direction using weighted scoring

### Degradation Detection
1. **Threshold Monitoring**: Configurable thresholds for key metrics
   - Error rate threshold: 10% increase (configurable)
   - Response time threshold: 20% increase (configurable) 
   - Traffic drop threshold: 30% decrease (configurable)

2. **Statistical Significance**: 
   - Minimum sample size validation (100 requests default)
   - Confidence level assessment (95% default)
   - Change magnitude evaluation

3. **Risk Scoring**: 0-100 scale risk assessment
   - Weighted by metric criticality (error rates weighted 2x)
   - Adjusted by significance level (high significance weighted 2x)
   - Capped at 100 for maximum risk

### Smart Alerting System
1. **Alert Prioritization**: Critical alerts for severe degradation
2. **Deduplication**: Prevents alert flooding with cooldown periods
3. **Actionable Recommendations**: Context-aware suggestions for each metric type
4. **Impact Assessment**: Business impact descriptions based on degradation severity

## CLI Integration

### Command-Line Flags
```bash
--trend-analysis              # Enable historical trend analysis
--compare-period "period"     # Specify comparison period (future enhancement)
```

### Usage Examples
```bash
# Basic trend analysis
./smart-log-analyser analyse logs/*.log --trend-analysis

# With ASCII visualizations  
./smart-log-analyser analyse logs/*.log --trend-analysis --ascii-charts

# Combined with other exports
./smart-log-analyser analyse logs/*.log --trend-analysis --export-html=report.html
```

## Visualization Components

### ASCII Chart Types
1. **Period Comparison Chart**: Horizontal bar chart showing percentage changes
2. **Degradation Alerts Chart**: Alert count visualization by severity
3. **Risk Score Gauge**: Horizontal gauge with color-coded risk levels
4. **Quick Summary**: Compact overview with health status and key metrics

### Color Coding System
- **Green**: Healthy/improving metrics
- **Yellow**: Warning/moderate degradation  
- **Red**: Critical/severe degradation
- **Blue**: Informational/stable metrics

## Configuration Management

### TrendConfiguration Structure
```go
type TrendConfiguration struct {
    ErrorRateThreshold      float64  // Error rate increase threshold (%)
    ResponseTimeThreshold   float64  // Response time increase threshold (%)
    TrafficDropThreshold    float64  // Traffic drop threshold (%)
    MinimumSampleSize       int      // Minimum requests for analysis
    SignificanceLevel       float64  // Statistical significance level
    DefaultComparisonPeriod string   // Default comparison period
    EnableAlerts            bool     // Alert generation toggle
    AlertCooldownHours      int      // Hours between similar alerts
}
```

### Default Configuration
- **Error Rate Threshold**: 10% (moderate sensitivity)
- **Response Time Threshold**: 20% (performance focus)
- **Traffic Drop Threshold**: 30% (availability concern)
- **Minimum Sample Size**: 100 requests (statistical validity)
- **Significance Level**: 0.05 (95% confidence)

## Output Format

### Console Display
```
╔════════════════════════════════════════════════════════════════╗
║                    Trend Analysis Results                      ║
╚════════════════════════════════════════════════════════════════╝

🏥 Overall Health: ⚠️ WARNING
📊 Analysis Type: degradation
🕒 Generated: 2025-08-24 20:01:23

📈 Trend Summary:
   Analysis shows degrading trend with risk score 9/100...

📋 Period Comparison:
├─ Overall Trend: 📉 degrading
├─ Risk Score: 9/100
└─ Key Changes: [significant changes listed]

🚨 Degradation Alerts:
├─ Alert TREND-001: ⚠️ Bot Traffic
│  Impact: Impact requires investigation
└─ Recommendation: Monitor metric closely...

💡 Recommendations:
   1. Monitor metric closely and investigate root causes
```

### ASCII Visualization Integration
- **Risk Score Gauge**: Horizontal progress bar with color coding
- **Alert Severity Chart**: Bar chart showing alert distribution
- **Metric Change Chart**: Percentage change visualization
- **Quick Summary**: Compact health overview

## Implementation Strategy

### Phase 1: Core Analysis Engine ✅
- Basic data structures and trend direction classification
- Period comparison algorithms with statistical analysis
- Degradation detection with configurable thresholds
- Risk scoring and alert generation

### Phase 2: CLI Integration ✅  
- Command-line flag implementation
- Console output formatting with emojis and structure
- Integration with existing analysis workflow
- Error handling and validation

### Phase 3: Visualization ✅
- ASCII chart generation for trends
- Risk gauge and alert severity visualization
- Color-coded output with terminal compatibility
- Quick summary and detailed breakdown options

### Future Enhancements
- **Period Specification**: Custom date ranges for comparison
- **Historical Database**: Store analysis results for long-term trending
- **Machine Learning**: Anomaly detection using historical patterns
- **Export Integration**: Include trend data in HTML/JSON/CSV exports

## Performance Considerations

### Memory Efficiency
- Stream processing for large log files
- Minimal data retention (only metrics, not raw logs)
- Efficient data structures for statistical calculations

### Processing Speed
- Leverages existing analysis engine for base metrics
- Parallel processing where applicable
- Minimal computational overhead for trend calculations

### Scalability
- Handles datasets from 100 requests to millions
- Configurable sample size thresholds
- Memory usage scales with metrics, not raw data volume

## Error Handling

### Data Validation
- Minimum sample size validation (prevents false alerts)
- Time range validation and overlap detection
- Log format compatibility checking

### Graceful Degradation
- Falls back to standard analysis if trend analysis fails
- Clear error messages for configuration issues
- Continues operation with reduced functionality on errors

## Security Considerations

### Data Privacy
- Processes only statistical metrics, not sensitive log content
- No persistent storage of analyzed data
- Same security model as base analyser component

### Alert Security
- No sensitive information in alert descriptions
- Generic recommendations to avoid information disclosure
- Configurable alert verbosity levels

## Testing Strategy

### Unit Testing
- Statistical calculation accuracy
- Threshold-based alert triggering  
- Risk score computation validation
- Configuration parameter handling

### Integration Testing
- End-to-end CLI workflow testing
- ASCII chart rendering verification
- Large dataset performance testing
- Error condition handling validation

### Real-World Testing
- Production log analysis validation
- Performance degradation scenario testing
- False positive/negative rate assessment
- User experience and output clarity evaluation

## Success Metrics

### Functional Success
- ✅ Accurate degradation detection with configurable sensitivity
- ✅ Clear, actionable alert generation
- ✅ Professional visualization output
- ✅ Seamless CLI integration

### Performance Success
- ✅ Processing 3000+ log entries in seconds
- ✅ Memory usage independent of raw log volume
- ✅ Minimal overhead on standard analysis workflow

### User Experience Success
- ✅ Intuitive command-line interface
- ✅ Clear, professionally formatted output
- ✅ Actionable recommendations for detected issues
- ✅ Visual clarity with ASCII charts and color coding

## Future Development

### Planned Enhancements
1. **Advanced Period Comparison**: Custom date range specification
2. **Baseline Learning**: Automatic baseline establishment from historical data
3. **Predictive Analysis**: Forecast future degradation based on trends
4. **Integration Enhancements**: Include trend data in HTML reports and exports

### Extensibility Points
- **Custom Metrics**: Plugin architecture for domain-specific metrics
- **Alert Channels**: Integration with monitoring systems (PagerDuty, Slack)
- **Database Storage**: Persistent trend history for long-term analysis
- **Machine Learning**: Advanced anomaly detection algorithms

This trend analysis feature represents a significant advancement in the Smart Log Analyser's capabilities, providing proactive monitoring and intelligent alerting for system performance degradation.