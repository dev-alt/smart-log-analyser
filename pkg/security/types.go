package security

import (
	"time"

	"smart-log-analyser/pkg/parser"
)

// WebAttackType represents different types of web application attacks
type WebAttackType int

const (
	SQLInjection WebAttackType = iota
	CrossSiteScripting
	CommandInjection
	DirectoryTraversal
	RemoteFileInclusion
	LocalFileInclusion
	XXEInjection
	DeserializationAttack
	HTTPHeaderInjection
	CSRFAttack
	AuthenticationBypass
	SessionHijacking
	Clickjacking
	CSPBypass
	HTTPSplitting
)

// String returns the string representation of WebAttackType
func (wat WebAttackType) String() string {
	switch wat {
	case SQLInjection:
		return "SQL Injection"
	case CrossSiteScripting:
		return "Cross-Site Scripting (XSS)"
	case CommandInjection:
		return "Command Injection"
	case DirectoryTraversal:
		return "Directory Traversal"
	case RemoteFileInclusion:
		return "Remote File Inclusion"
	case LocalFileInclusion:
		return "Local File Inclusion"
	case XXEInjection:
		return "XML External Entity (XXE)"
	case DeserializationAttack:
		return "Deserialization Attack"
	case HTTPHeaderInjection:
		return "HTTP Header Injection"
	case CSRFAttack:
		return "Cross-Site Request Forgery (CSRF)"
	case AuthenticationBypass:
		return "Authentication Bypass"
	case SessionHijacking:
		return "Session Hijacking"
	case Clickjacking:
		return "Clickjacking"
	case CSPBypass:
		return "Content Security Policy Bypass"
	case HTTPSplitting:
		return "HTTP Response Splitting"
	default:
		return "Unknown Attack"
	}
}

// InfrastructureAttackType represents different types of infrastructure attacks
type InfrastructureAttackType int

const (
	BruteForceLogin InfrastructureAttackType = iota
	PasswordSpray
	DDoSAttack
	PortScan
	VulnerabilityScanning
	WebShellAccess
	PrivilegeEscalation
	DataExfiltration
	BotnetActivity
	CryptoMining
	ResourceExhaustion
	ServiceEnumeration
	ForceBrowsing
	RateLimitEvasion
	CachePoison
)

// String returns the string representation of InfrastructureAttackType
func (iat InfrastructureAttackType) String() string {
	switch iat {
	case BruteForceLogin:
		return "Brute Force Login"
	case PasswordSpray:
		return "Password Spray Attack"
	case DDoSAttack:
		return "Distributed Denial of Service"
	case PortScan:
		return "Port Scanning"
	case VulnerabilityScanning:
		return "Vulnerability Scanning"
	case WebShellAccess:
		return "Web Shell Access"
	case PrivilegeEscalation:
		return "Privilege Escalation"
	case DataExfiltration:
		return "Data Exfiltration"
	case BotnetActivity:
		return "Botnet Activity"
	case CryptoMining:
		return "Cryptocurrency Mining"
	case ResourceExhaustion:
		return "Resource Exhaustion"
	case ServiceEnumeration:
		return "Service Enumeration"
	case ForceBrowsing:
		return "Forced Browsing"
	case RateLimitEvasion:
		return "Rate Limit Evasion"
	case CachePoison:
		return "Cache Poisoning"
	default:
		return "Unknown Infrastructure Attack"
	}
}

// ThreatSeverity represents the severity level of a threat
type ThreatSeverity int

const (
	SeverityInfo ThreatSeverity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// String returns the string representation of ThreatSeverity
func (ts ThreatSeverity) String() string {
	switch ts {
	case SeverityInfo:
		return "Info"
	case SeverityLow:
		return "Low"
	case SeverityMedium:
		return "Medium"
	case SeverityHigh:
		return "High"
	case SeverityCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// RiskLevel represents the overall risk level
type RiskLevel int

const (
	RiskMinimal RiskLevel = iota // 90-100 score
	RiskLow                      // 70-89 score
	RiskMedium                   // 50-69 score
	RiskHigh                     // 30-49 score
	RiskCritical                 // 0-29 score
)

// String returns the string representation of RiskLevel
func (rl RiskLevel) String() string {
	switch rl {
	case RiskMinimal:
		return "Minimal"
	case RiskLow:
		return "Low"
	case RiskMedium:
		return "Medium"
	case RiskHigh:
		return "High"
	case RiskCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// AnomalyType represents different types of behavioral anomalies
type AnomalyType int

const (
	AnomalyRequestFrequency AnomalyType = iota
	AnomalyRequestTiming
	AnomalyRequestSize
	AnomalyErrorRate
	AnomalyUserAgent
	AnomalyGeographic
	AnomalyEndpointPattern
	AnomalyStatusCodePattern
)

// String returns the string representation of AnomalyType
func (at AnomalyType) String() string {
	switch at {
	case AnomalyRequestFrequency:
		return "Unusual Request Frequency"
	case AnomalyRequestTiming:
		return "Unusual Request Timing"
	case AnomalyRequestSize:
		return "Unusual Request Size"
	case AnomalyErrorRate:
		return "Unusual Error Rate"
	case AnomalyUserAgent:
		return "Unusual User Agent"
	case AnomalyGeographic:
		return "Unusual Geographic Location"
	case AnomalyEndpointPattern:
		return "Unusual Endpoint Access Pattern"
	case AnomalyStatusCodePattern:
		return "Unusual Status Code Pattern"
	default:
		return "Unknown Anomaly"
	}
}

// EnhancedThreat represents a detected security threat with advanced attributes
type EnhancedThreat struct {
	ID                string
	Type              interface{} // WebAttackType or InfrastructureAttackType
	Severity          ThreatSeverity
	Confidence        float64 // 0.0-1.0
	Pattern           string
	URL               string
	IP                string
	UserAgent         string
	Timestamp         time.Time
	Method            string
	StatusCode        int
	ResponseSize      int64
	AttackVector      string
	Payload           string
	Context           map[string]interface{}
	RelatedThreats    []string // IDs of related threats
	IOCs              []string // Indicators of Compromise
	MitigationAdvice  []string
}

// Anomaly represents a behavioral anomaly detection
type Anomaly struct {
	ID            string
	Type          AnomalyType
	Severity      ThreatSeverity
	Description   string
	Metric        string
	ExpectedValue float64
	ActualValue   float64
	Deviation     float64
	ZScore        float64
	IP            string
	Timestamp     time.Time
	TimeWindow    time.Duration
	Confidence    float64
	Context       map[string]interface{}
}

// IPBehaviorProfile represents behavioral analysis of an IP address
type IPBehaviorProfile struct {
	IP                      string
	FirstSeen               time.Time
	LastSeen                time.Time
	TotalRequests           int64
	RequestFrequency        float64 // requests per minute
	AverageRequestInterval  time.Duration
	TypicalRequestTimes     []time.Time
	CommonUserAgents        map[string]int
	VisitedEndpoints        map[string]int
	HTTPMethods             map[string]int
	StatusCodeDistribution  map[int]int
	ErrorRate               float64
	AverageResponseSize     int64
	GeographicConsistency   bool
	GeographicLocations     []string
	BehaviorScore           float64 // 0.0-1.0 (higher = more suspicious)
	RiskLevel               RiskLevel
	Anomalies               []Anomaly
	AssociatedThreats       []string
	Tags                    []string // "bot", "scanner", "legitimate", etc.
}

// SecurityDimensions represents different aspects of security analysis
type SecurityDimensions struct {
	ThreatDetection   float64 // 40% weight - Direct threat identification
	AnomalyDetection  float64 // 25% weight - Behavioral anomalies
	TrafficIntegrity  float64 // 20% weight - Traffic pattern health
	AccessControl     float64 // 15% weight - Authentication/authorization issues
}

// ThreatIntelligence represents threat intelligence data
type ThreatIntelligence struct {
	MaliciousIPs     map[string]ThreatInfo
	AttackSignatures []AttackSignature
	KnownPayloads    map[string]PayloadInfo
	VulnerabilityDB  []VulnerabilityInfo
}

// ThreatInfo represents information about a known threat
type ThreatInfo struct {
	IP            string
	ThreatTypes   []string
	Severity      ThreatSeverity
	FirstSeen     time.Time
	LastSeen      time.Time
	Attribution   string
	IOCs          []string
	Description   string
	References    []string
}

// AttackSignature represents a known attack signature
type AttackSignature struct {
	ID          string
	Name        string
	Pattern     string
	Type        interface{} // WebAttackType or InfrastructureAttackType
	Severity    ThreatSeverity
	Description string
	References  []string
}

// PayloadInfo represents information about a known malicious payload
type PayloadInfo struct {
	Payload     string
	Type        interface{}
	Severity    ThreatSeverity
	Description string
	Variants    []string
}

// VulnerabilityInfo represents vulnerability information
type VulnerabilityInfo struct {
	CVE         string
	Description string
	Severity    ThreatSeverity
	CVSS        float64
	Patterns    []string
}

// SecurityRecommendation represents actionable security advice
type SecurityRecommendation struct {
	Priority    int         // 1-10
	Category    string
	Title       string
	Description string
	Impact      ThreatSeverity
	Effort      string // "Low", "Medium", "High"
	Actions     []string
	References  []string
}

// IncidentData represents data for incident response
type IncidentData struct {
	ID               string
	Title            string
	Severity         ThreatSeverity
	StartTime        time.Time
	EndTime          time.Time
	AffectedSystems  []string
	AttackVector     string
	ThreatActor      string
	IOCs             []string
	Timeline         []IncidentEvent
	Impact           string
	Recommendations  []SecurityRecommendation
	Evidence         []string
}

// IncidentEvent represents a single event in an incident timeline
type IncidentEvent struct {
	Timestamp   time.Time
	Description string
	Type        string
	Severity    ThreatSeverity
	Source      string
	Details     map[string]interface{}
}

// SecuritySummary represents high-level security overview
type SecuritySummary struct {
	OverallRisk         RiskLevel
	SecurityScore       int // 0-100
	SecurityDimensions  SecurityDimensions
	ActiveThreats       int
	CriticalVulns       int
	HighRiskIPs         []string
	TopAttackTypes      []string
	ThreatTrends        []ThreatTrend
	RecommendedActions  []SecurityRecommendation
	ComplianceScore     int
	IncidentCount       int
	TimeRange           TimeRange
}

// ThreatTrend represents threat trends over time
type ThreatTrend struct {
	Type      string
	Count     int
	Trend     float64 // positive = increasing, negative = decreasing
	TimeRange TimeRange
}

// TimeRange represents a time period
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// EnhancedSecurityAnalysis represents comprehensive security analysis results
type EnhancedSecurityAnalysis struct {
	Summary             SecuritySummary
	Threats             []EnhancedThreat
	Anomalies           []Anomaly
	IPProfiles          map[string]*IPBehaviorProfile
	ThreatIntelligence  *ThreatIntelligence
	Incidents           []IncidentData
	Recommendations     []SecurityRecommendation
	ComplianceData      map[string]interface{}
	AnalysisTimestamp   time.Time
	LogTimeRange        TimeRange
	TotalEntriesAnalyzed int64
}

// SecurityAnalyzer interface defines the main security analysis capabilities
type SecurityAnalyzer interface {
	Analyze(logs []*parser.LogEntry) (*EnhancedSecurityAnalysis, error)
	DetectWebAttacks(logs []*parser.LogEntry) ([]EnhancedThreat, error)
	DetectInfrastructureAttacks(logs []*parser.LogEntry) ([]EnhancedThreat, error)
	DetectAnomalies(logs []*parser.LogEntry) ([]Anomaly, error)
	ProfileIPs(logs []*parser.LogEntry) (map[string]*IPBehaviorProfile, error)
	CalculateSecurityScore(analysis *EnhancedSecurityAnalysis) int
	GenerateIncidents(threats []EnhancedThreat, anomalies []Anomaly) ([]IncidentData, error)
	GenerateRecommendations(analysis *EnhancedSecurityAnalysis) ([]SecurityRecommendation, error)
}

// Configuration for security analysis
type SecurityConfig struct {
	ThreatDetectionSensitivity float64 // 1.0-10.0
	AnomalyThreshold          float64 // Z-score threshold (default: 2.5)
	BehavioralAnalysisEnabled bool
	ThreatIntelligenceEnabled bool
	IncidentResponseEnabled   bool
	ComplianceReportingEnabled bool
}

// Default configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		ThreatDetectionSensitivity: 7.0,
		AnomalyThreshold:          2.5,
		BehavioralAnalysisEnabled: true,
		ThreatIntelligenceEnabled: true,
		IncidentResponseEnabled:   true,
		ComplianceReportingEnabled: true,
	}
}