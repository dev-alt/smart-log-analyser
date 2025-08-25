package analyser

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"smart-log-analyser/pkg/parser"
)

type IPStat struct {
	IP    string
	Count int
}

type URLStat struct {
	URL         string
	Count       int
	StatusCodes map[int]int // Maps status code to count (for error URLs)
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type MethodStat struct {
	Method string
	Count  int
}

type BotStat struct {
	BotName string
	Count   int
}

type FileTypeStat struct {
	FileType string
	Count    int
	Size     int64
}

type HourlyTraffic struct {
	Hour         int    // Hour of day (0-23)
	RequestCount int    // Number of requests in that hour
	Timestamp    string // Human readable timestamp for the hour
}

type TrafficPeak struct {
	Time         string // Timestamp of peak
	RequestCount int    // Number of requests during peak period
	Duration     string // Peak duration description
}

type ResponseTimeStats struct {
	AverageSize    int64   // Average response size (proxy for response time)
	MedianSize     int64   // 50th percentile
	P95Size        int64   // 95th percentile  
	P99Size        int64   // 99th percentile
	MinSize        int64   // Smallest response
	MaxSize        int64   // Largest response
	SlowRequests   []URLStat // URLs with largest response sizes
	FastRequests   []URLStat // URLs with smallest response sizes
}

type GeographicStat struct {
	Country string
	Count   int
	Region  string // Continent/region
}

type GeographicAnalysis struct {
	TopCountries     []GeographicStat
	TopRegions       []GeographicStat
	TotalCountries   int
	UnknownIPs       int
	LocalTraffic     int // Private IP ranges
	CloudTraffic     int // CDN/Cloud provider IPs
}

type SecurityThreat struct {
	Type        string // "sql_injection", "xss", "directory_traversal", "brute_force", etc.
	Pattern     string // The malicious pattern detected
	URL         string // The targeted URL
	IP          string // Source IP
	Timestamp   time.Time
	Severity    string // "low", "medium", "high", "critical"
	UserAgent   string // User agent string
}

type AnomalyDetection struct {
	Type          string  // Type of anomaly
	Description   string  // Human readable description
	Value         float64 // Actual value
	Expected      float64 // Expected/baseline value
	Deviation     float64 // How much it deviates (percentage)
	Significance  string  // "low", "medium", "high"
}

type IPThreatAnalysis struct {
	IP               string
	RequestCount     int
	ThreatScore      int    // 0-100 scale
	ThreatCategories []string // "brute_force", "scanner", "malicious_patterns", etc.
	FirstSeen        time.Time
	LastSeen         time.Time
	UniqueURLs       int
	ErrorRate        float64 // Percentage of requests resulting in errors
}

type SecurityAnalysis struct {
	ThreatLevel          string             // "low", "medium", "high", "critical"
	SecurityScore        int                // 0-100, higher is better
	TotalThreats         int
	ThreatsDetected      []SecurityThreat
	SuspiciousIPs        []IPThreatAnalysis
	AnomaliesDetected    []AnomalyDetection
	BruteForceAttempts   int
	SQLInjectionAttempts int
	XSSAttempts          int
	DirectoryTraversal   int
	ScanningActivity     int
	TopAttackers         []IPStat // IPs with most malicious activity
}

type DetailedStatusCode struct {
	Code  int
	Count int
}

type Results struct {
	TotalRequests          int
	TimeRange              TimeRange
	StatusCodes            map[string]int
	DetailedStatusCodes    []DetailedStatusCode
	TopIPs                 []IPStat
	TopURLs                []URLStat
	HTTPMethods            []MethodStat
	TotalBytes             int64
	AverageSize            int64
	UniqueIPs              int
	UniqueURLs             int
	BotRequests            int
	HumanRequests          int
	TopBots                []BotStat
	FileTypes              []FileTypeStat
	ErrorURLs              []URLStat // URLs that generated errors
	LargeRequests          []URLStat // Largest requests by size
	HourlyTraffic          []HourlyTraffic
	TrafficPeaks           []TrafficPeak
	AverageRequestsPerHour float64
	PeakHour               int
	QuietestHour           int
	ResponseTimeStats      ResponseTimeStats
	GeographicAnalysis     GeographicAnalysis
	SecurityAnalysis       SecurityAnalysis
}

type Analyser struct{}

func New() *Analyser {
	return &Analyser{}
}

func (a *Analyser) Analyse(logs []*parser.LogEntry, since, until *time.Time) *Results {
	filtered := a.FilterByTime(logs, since, until)
	
	if len(filtered) == 0 {
		return &Results{
			TotalRequests:          0,
			TimeRange:              TimeRange{},
			StatusCodes:            make(map[string]int),
			DetailedStatusCodes:    []DetailedStatusCode{},
			TopIPs:                 []IPStat{},
			TopURLs:                []URLStat{},
			HTTPMethods:            []MethodStat{},
			TotalBytes:             0,
			AverageSize:            0,
			UniqueIPs:              0,
			UniqueURLs:             0,
			BotRequests:            0,
			HumanRequests:          0,
			TopBots:                []BotStat{},
			FileTypes:              []FileTypeStat{},
			ErrorURLs:              []URLStat{},
			LargeRequests:          []URLStat{},
			HourlyTraffic:          []HourlyTraffic{},
			TrafficPeaks:           []TrafficPeak{},
			AverageRequestsPerHour: 0,
			PeakHour:               -1,
			QuietestHour:           -1,
			ResponseTimeStats:      ResponseTimeStats{},
			GeographicAnalysis:     GeographicAnalysis{},
			SecurityAnalysis:       SecurityAnalysis{},
		}
	}

	botRequests, humanRequests := a.analyseBotTraffic(filtered)
	hourlyTraffic := a.analyseHourlyTraffic(filtered)
	trafficPeaks := a.detectTrafficPeaks(hourlyTraffic)
	avgPerHour, peakHour, quietestHour := a.calculateTrafficStats(hourlyTraffic)
	responseTimeStats := a.analyseResponseTimes(filtered)
	geographicAnalysis := a.analyseGeographicDistribution(filtered)
	securityAnalysis := a.analyseSecurityThreats(filtered)
	
	results := &Results{
		TotalRequests:          len(filtered),
		TimeRange:              a.calculateTimeRange(filtered),
		StatusCodes:            a.analyseStatusCodes(filtered),
		DetailedStatusCodes:    a.analyseDetailedStatusCodes(filtered),
		TopIPs:                 a.analyseTopIPs(filtered),
		TopURLs:                a.analyseTopURLs(filtered),
		HTTPMethods:            a.analyseHTTPMethods(filtered),
		TotalBytes:             a.calculateTotalBytes(filtered),
		AverageSize:            a.calculateAverageSize(filtered),
		UniqueIPs:              a.countUniqueIPs(filtered),
		UniqueURLs:             a.countUniqueURLs(filtered),
		BotRequests:            botRequests,
		HumanRequests:          humanRequests,
		TopBots:                a.analyseTopBots(filtered),
		FileTypes:              a.analyseFileTypes(filtered),
		ErrorURLs:              a.analyseErrorURLs(filtered),
		LargeRequests:          a.analyseLargeRequests(filtered),
		HourlyTraffic:          hourlyTraffic,
		TrafficPeaks:           trafficPeaks,
		AverageRequestsPerHour: avgPerHour,
		PeakHour:               peakHour,
		QuietestHour:           quietestHour,
		ResponseTimeStats:      responseTimeStats,
		GeographicAnalysis:     geographicAnalysis,
		SecurityAnalysis:       securityAnalysis,
	}

	return results
}

func (a *Analyser) FilterByTime(logs []*parser.LogEntry, since, until *time.Time) []*parser.LogEntry {
	var filtered []*parser.LogEntry

	for _, log := range logs {
		if since != nil && log.Timestamp.Before(*since) {
			continue
		}
		if until != nil && log.Timestamp.After(*until) {
			continue
		}
		filtered = append(filtered, log)
	}

	return filtered
}

func (a *Analyser) calculateTimeRange(logs []*parser.LogEntry) TimeRange {
	if len(logs) == 0 {
		return TimeRange{}
	}

	start := logs[0].Timestamp
	end := logs[0].Timestamp

	for _, log := range logs {
		if log.Timestamp.Before(start) {
			start = log.Timestamp
		}
		if log.Timestamp.After(end) {
			end = log.Timestamp
		}
	}

	return TimeRange{Start: start, End: end}
}

func (a *Analyser) analyseStatusCodes(logs []*parser.LogEntry) map[string]int {
	statusCodes := make(map[string]int)

	for _, log := range logs {
		status := getStatusClass(log.Status)
		statusCodes[status]++
	}

	return statusCodes
}

func (a *Analyser) analyseTopIPs(logs []*parser.LogEntry) []IPStat {
	ipCounts := make(map[string]int)

	for _, log := range logs {
		ipCounts[log.IP]++
	}

	var ipStats []IPStat
	for ip, count := range ipCounts {
		ipStats = append(ipStats, IPStat{IP: ip, Count: count})
	}

	sort.Slice(ipStats, func(i, j int) bool {
		return ipStats[i].Count > ipStats[j].Count
	})

	return ipStats
}

func (a *Analyser) analyseTopURLs(logs []*parser.LogEntry) []URLStat {
	urlCounts := make(map[string]int)

	for _, log := range logs {
		urlCounts[log.URL]++
	}

	var urlStats []URLStat
	for url, count := range urlCounts {
		urlStats = append(urlStats, URLStat{
			URL:         url, 
			Count:       count,
			StatusCodes: nil, // Not applicable for top URLs (not error-specific)
		})
	}

	sort.Slice(urlStats, func(i, j int) bool {
		return urlStats[i].Count > urlStats[j].Count
	})

	return urlStats
}

// FormatStatusCodes formats status codes from a URLStat for display
func (u *URLStat) FormatStatusCodes() string {
	if u.StatusCodes == nil || len(u.StatusCodes) == 0 {
		return "N/A"
	}
	
	var codes []string
	for status := range u.StatusCodes {
		codes = append(codes, fmt.Sprintf("%d", status))
	}
	
	// Sort status codes numerically
	sort.Slice(codes, func(i, j int) bool {
		a, _ := strconv.Atoi(codes[i])
		b, _ := strconv.Atoi(codes[j])
		return a < b
	})
	
	return strings.Join(codes, "/")
}

func (a *Analyser) analyseHTTPMethods(logs []*parser.LogEntry) []MethodStat {
	methodCounts := make(map[string]int)

	for _, log := range logs {
		methodCounts[log.Method]++
	}

	var methodStats []MethodStat
	for method, count := range methodCounts {
		methodStats = append(methodStats, MethodStat{Method: method, Count: count})
	}

	sort.Slice(methodStats, func(i, j int) bool {
		return methodStats[i].Count > methodStats[j].Count
	})

	return methodStats
}

func (a *Analyser) calculateTotalBytes(logs []*parser.LogEntry) int64 {
	var total int64
	for _, log := range logs {
		total += log.Size
	}
	return total
}

func (a *Analyser) calculateAverageSize(logs []*parser.LogEntry) int64 {
	if len(logs) == 0 {
		return 0
	}
	return a.calculateTotalBytes(logs) / int64(len(logs))
}

func (a *Analyser) countUniqueIPs(logs []*parser.LogEntry) int {
	unique := make(map[string]bool)
	for _, log := range logs {
		unique[log.IP] = true
	}
	return len(unique)
}

func (a *Analyser) countUniqueURLs(logs []*parser.LogEntry) int {
	unique := make(map[string]bool)
	for _, log := range logs {
		unique[log.URL] = true
	}
	return len(unique)
}

func (a *Analyser) analyseBotTraffic(logs []*parser.LogEntry) (int, int) {
	botCount := 0
	humanCount := 0
	
	for _, log := range logs {
		if isBot(log.UserAgent) {
			botCount++
		} else {
			humanCount++
		}
	}
	
	return botCount, humanCount
}

func (a *Analyser) analyseTopBots(logs []*parser.LogEntry) []BotStat {
	botCounts := make(map[string]int)
	
	for _, log := range logs {
		if botName := getBotName(log.UserAgent); botName != "" {
			botCounts[botName]++
		}
	}
	
	var botStats []BotStat
	for bot, count := range botCounts {
		botStats = append(botStats, BotStat{BotName: bot, Count: count})
	}
	
	sort.Slice(botStats, func(i, j int) bool {
		return botStats[i].Count > botStats[j].Count
	})
	
	return botStats
}

func (a *Analyser) analyseFileTypes(logs []*parser.LogEntry) []FileTypeStat {
	fileTypeCounts := make(map[string]int)
	fileTypeSizes := make(map[string]int64)
	
	for _, log := range logs {
		fileType := getFileType(log.URL)
		fileTypeCounts[fileType]++
		fileTypeSizes[fileType] += log.Size
	}
	
	var fileTypeStats []FileTypeStat
	for fileType, count := range fileTypeCounts {
		fileTypeStats = append(fileTypeStats, FileTypeStat{
			FileType: fileType,
			Count:    count,
			Size:     fileTypeSizes[fileType],
		})
	}
	
	sort.Slice(fileTypeStats, func(i, j int) bool {
		return fileTypeStats[i].Count > fileTypeStats[j].Count
	})
	
	return fileTypeStats
}

func isBot(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	botKeywords := []string{
		"bot", "crawler", "spider", "scraper", "parser",
		"googlebot", "bingbot", "slurp", "facebookexternalhit",
		"twitterbot", "linkedinbot", "whatsapp", "telegram",
		"curl", "wget", "python", "go-http-client", "java",
		"monitoring", "uptime", "check", "test", "scan",
	}
	
	for _, keyword := range botKeywords {
		if strings.Contains(ua, keyword) {
			return true
		}
	}
	
	return false
}

func getBotName(userAgent string) string {
	if !isBot(userAgent) {
		return ""
	}
	
	ua := strings.ToLower(userAgent)
	
	// Common bot patterns
	botPatterns := map[string]string{
		"googlebot":              "Googlebot",
		"bingbot":                "Bingbot", 
		"slurp":                  "Yahoo Slurp",
		"facebookexternalhit":    "Facebook Bot",
		"twitterbot":             "Twitter Bot",
		"linkedinbot":            "LinkedIn Bot",
		"whatsapp":               "WhatsApp Bot",
		"telegram":               "Telegram Bot",
		"curl":                   "cURL",
		"wget":                   "Wget",
		"python":                 "Python Script",
		"go-http-client":         "Go HTTP Client",
		"java":                   "Java Client",
		"monitoring":             "Monitoring Bot",
		"uptime":                 "Uptime Monitor",
		"check":                  "Health Check",
		"scan":                   "Security Scanner",
	}
	
	for pattern, name := range botPatterns {
		if strings.Contains(ua, pattern) {
			return name
		}
	}
	
	return "Unknown Bot"
}

func getFileType(url string) string {
	// Remove query parameters
	url = strings.Split(url, "?")[0]
	
	// Get file extension
	parts := strings.Split(url, ".")
	if len(parts) < 2 {
		return "Dynamic/HTML"
	}
	
	ext := strings.ToLower(parts[len(parts)-1])
	
	// Group by file type categories
	switch ext {
	case "css":
		return "CSS"
	case "js":
		return "JavaScript"
	case "jpg", "jpeg", "png", "gif", "webp", "ico", "svg":
		return "Images"
	case "pdf":
		return "PDF"
	case "xml":
		return "XML"
	case "txt":
		return "Text Files"
	case "zip", "tar", "gz", "rar":
		return "Archives"
	case "mp4", "avi", "mov", "wmv":
		return "Videos"
	case "mp3", "wav", "ogg":
		return "Audio"
	case "woff", "woff2", "ttf", "eot":
		return "Fonts"
	default:
		return "Dynamic/HTML"
	}
}

func (a *Analyser) analyseDetailedStatusCodes(logs []*parser.LogEntry) []DetailedStatusCode {
	statusCounts := make(map[int]int)
	
	for _, log := range logs {
		statusCounts[log.Status]++
	}
	
	var statusStats []DetailedStatusCode
	for status, count := range statusCounts {
		statusStats = append(statusStats, DetailedStatusCode{Code: status, Count: count})
	}
	
	sort.Slice(statusStats, func(i, j int) bool {
		return statusStats[i].Count > statusStats[j].Count
	})
	
	return statusStats
}

func (a *Analyser) analyseErrorURLs(logs []*parser.LogEntry) []URLStat {
	// Map from URL to status code counts
	errorData := make(map[string]map[int]int)
	
	for _, log := range logs {
		if log.Status >= 400 { // 4xx and 5xx errors
			if errorData[log.URL] == nil {
				errorData[log.URL] = make(map[int]int)
			}
			errorData[log.URL][log.Status]++
		}
	}
	
	var errorStats []URLStat
	for url, statusCodes := range errorData {
		// Calculate total count for this URL
		totalCount := 0
		for _, count := range statusCodes {
			totalCount += count
		}
		
		errorStats = append(errorStats, URLStat{
			URL:         url,
			Count:       totalCount,
			StatusCodes: statusCodes,
		})
	}
	
	sort.Slice(errorStats, func(i, j int) bool {
		return errorStats[i].Count > errorStats[j].Count
	})
	
	// Return top 10 error URLs
	if len(errorStats) > 10 {
		errorStats = errorStats[:10]
	}
	
	return errorStats
}

func (a *Analyser) analyseLargeRequests(logs []*parser.LogEntry) []URLStat {
	type urlSize struct {
		url  string
		size int64
	}
	
	var requests []urlSize
	for _, log := range logs {
		requests = append(requests, urlSize{url: log.URL, size: log.Size})
	}
	
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].size > requests[j].size
	})
	
	// Convert to URLStat format (using size as count for sorting)
	var largeStats []URLStat
	seen := make(map[string]bool)
	
	for _, req := range requests {
		if !seen[req.url] && len(largeStats) < 10 {
			largeStats = append(largeStats, URLStat{
				URL:         req.url,
				Count:       int(req.size), // Store size in count field for display
				StatusCodes: nil,           // Not applicable for large requests
			})
			seen[req.url] = true
		}
	}
	
	return largeStats
}

func (a *Analyser) analyseHourlyTraffic(logs []*parser.LogEntry) []HourlyTraffic {
	if len(logs) == 0 {
		return []HourlyTraffic{}
	}
	
	// Count requests per hour
	hourlyCounts := make(map[int]int)
	hourTimestamps := make(map[int]string)
	
	for _, log := range logs {
		hour := log.Timestamp.Hour()
		hourlyCounts[hour]++
		
		// Store a representative timestamp for this hour (first occurrence)
		if _, exists := hourTimestamps[hour]; !exists {
			hourTimestamps[hour] = log.Timestamp.Format("2006-01-02 15:00")
		}
	}
	
	// Convert to slice and sort by hour
	var hourlyTraffic []HourlyTraffic
	for hour, count := range hourlyCounts {
		hourlyTraffic = append(hourlyTraffic, HourlyTraffic{
			Hour:         hour,
			RequestCount: count,
			Timestamp:    hourTimestamps[hour],
		})
	}
	
	sort.Slice(hourlyTraffic, func(i, j int) bool {
		return hourlyTraffic[i].Hour < hourlyTraffic[j].Hour
	})
	
	return hourlyTraffic
}

func (a *Analyser) detectTrafficPeaks(hourlyTraffic []HourlyTraffic) []TrafficPeak {
	if len(hourlyTraffic) < 3 {
		return []TrafficPeak{}
	}
	
	var peaks []TrafficPeak
	
	// Calculate average requests per hour
	totalRequests := 0
	for _, traffic := range hourlyTraffic {
		totalRequests += traffic.RequestCount
	}
	avgRequestsPerHour := float64(totalRequests) / float64(len(hourlyTraffic))
	
	// Define peak threshold as 150% of average
	peakThreshold := avgRequestsPerHour * 1.5
	
	for i, traffic := range hourlyTraffic {
		if float64(traffic.RequestCount) > peakThreshold {
			// Check if this is a local maximum
			isPeak := true
			
			// Check previous hour
			if i > 0 && hourlyTraffic[i-1].RequestCount >= traffic.RequestCount {
				isPeak = false
			}
			
			// Check next hour
			if i < len(hourlyTraffic)-1 && hourlyTraffic[i+1].RequestCount >= traffic.RequestCount {
				isPeak = false
			}
			
			if isPeak {
				peaks = append(peaks, TrafficPeak{
					Time:         traffic.Timestamp,
					RequestCount: traffic.RequestCount,
					Duration:     "1 hour", // For now, consider each peak as 1 hour
				})
			}
		}
	}
	
	// Sort peaks by request count (highest first)
	sort.Slice(peaks, func(i, j int) bool {
		return peaks[i].RequestCount > peaks[j].RequestCount
	})
	
	// Limit to top 5 peaks
	if len(peaks) > 5 {
		peaks = peaks[:5]
	}
	
	return peaks
}

func (a *Analyser) calculateTrafficStats(hourlyTraffic []HourlyTraffic) (float64, int, int) {
	if len(hourlyTraffic) == 0 {
		return 0, -1, -1
	}
	
	totalRequests := 0
	peakHour := -1
	quietestHour := -1
	maxRequests := -1
	minRequests := int(^uint(0) >> 1) // Max int value
	
	for _, traffic := range hourlyTraffic {
		totalRequests += traffic.RequestCount
		
		if traffic.RequestCount > maxRequests {
			maxRequests = traffic.RequestCount
			peakHour = traffic.Hour
		}
		
		if traffic.RequestCount < minRequests {
			minRequests = traffic.RequestCount
			quietestHour = traffic.Hour
		}
	}
	
	avgRequestsPerHour := float64(totalRequests) / float64(len(hourlyTraffic))
	
	return avgRequestsPerHour, peakHour, quietestHour
}

func (a *Analyser) analyseResponseTimes(logs []*parser.LogEntry) ResponseTimeStats {
	if len(logs) == 0 {
		return ResponseTimeStats{}
	}
	
	// Collect all response sizes for percentile calculation
	sizes := make([]int64, len(logs))
	totalSize := int64(0)
	minSize := int64(^uint64(0) >> 1) // Max int64
	maxSize := int64(0)
	
	for i, log := range logs {
		sizes[i] = log.Size
		totalSize += log.Size
		
		if log.Size < minSize {
			minSize = log.Size
		}
		if log.Size > maxSize {
			maxSize = log.Size
		}
	}
	
	// Sort sizes for percentile calculation
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] < sizes[j]
	})
	
	// Calculate percentiles
	p50Index := len(sizes) * 50 / 100
	p95Index := len(sizes) * 95 / 100
	p99Index := len(sizes) * 99 / 100
	
	// Ensure indices are within bounds
	if p50Index >= len(sizes) { p50Index = len(sizes) - 1 }
	if p95Index >= len(sizes) { p95Index = len(sizes) - 1 }
	if p99Index >= len(sizes) { p99Index = len(sizes) - 1 }
	
	avgSize := totalSize / int64(len(logs))
	
	// Find slowest and fastest requests (by size as proxy)
	slowRequests := a.analyseLargeRequests(logs)  // Reuse existing logic
	fastRequests := a.analyseSmallRequests(logs)
	
	return ResponseTimeStats{
		AverageSize:  avgSize,
		MedianSize:   sizes[p50Index],
		P95Size:      sizes[p95Index],
		P99Size:      sizes[p99Index],
		MinSize:      minSize,
		MaxSize:      maxSize,
		SlowRequests: slowRequests,
		FastRequests: fastRequests,
	}
}

func (a *Analyser) analyseSmallRequests(logs []*parser.LogEntry) []URLStat {
	type urlSize struct {
		url  string
		size int64
	}
	
	var requests []urlSize
	for _, log := range logs {
		requests = append(requests, urlSize{url: log.URL, size: log.Size})
	}
	
	// Sort by size (smallest first)
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].size < requests[j].size
	})
	
	// Convert to URLStat format
	var smallStats []URLStat
	seen := make(map[string]bool)
	
	for _, req := range requests {
		if !seen[req.url] && len(smallStats) < 10 {
			smallStats = append(smallStats, URLStat{
				URL:   req.url,
				Count: int(req.size), // Store size in count field for display
			})
			seen[req.url] = true
		}
	}
	
	return smallStats
}

func (a *Analyser) analyseGeographicDistribution(logs []*parser.LogEntry) GeographicAnalysis {
	countryCounts := make(map[string]int)
	regionCounts := make(map[string]int)
	
	localTraffic := 0
	cloudTraffic := 0
	unknownIPs := 0
	
	for _, log := range logs {
		country, region := a.getIPLocation(log.IP)
		
		if country == "Local" {
			localTraffic++
		} else if country == "Cloud" {
			cloudTraffic++
		} else if country == "Unknown" {
			unknownIPs++
		} else {
			countryCounts[country]++
			regionCounts[region]++
		}
	}
	
	// Convert to sorted slices
	var topCountries []GeographicStat
	for country, count := range countryCounts {
		topCountries = append(topCountries, GeographicStat{
			Country: country,
			Count:   count,
			Region:  a.getRegionForCountry(country),
		})
	}
	
	var topRegions []GeographicStat
	for region, count := range regionCounts {
		topRegions = append(topRegions, GeographicStat{
			Country: region,
			Count:   count,
			Region:  region,
		})
	}
	
	// Sort by count
	sort.Slice(topCountries, func(i, j int) bool {
		return topCountries[i].Count > topCountries[j].Count
	})
	
	sort.Slice(topRegions, func(i, j int) bool {
		return topRegions[i].Count > topRegions[j].Count
	})
	
	return GeographicAnalysis{
		TopCountries:   topCountries,
		TopRegions:     topRegions,
		TotalCountries: len(countryCounts),
		UnknownIPs:     unknownIPs,
		LocalTraffic:   localTraffic,
		CloudTraffic:   cloudTraffic,
	}
}

func (a *Analyser) getIPLocation(ip string) (string, string) {
	// Simple IP-based location detection using common patterns
	
	// Private IP ranges
	if strings.HasPrefix(ip, "192.168.") || 
	   strings.HasPrefix(ip, "10.") || 
	   strings.HasPrefix(ip, "172.") {
		return "Local", "Private Network"
	}
	
	// Common cloud/CDN providers (based on known ranges)
	if strings.HasPrefix(ip, "172.69.") || strings.HasPrefix(ip, "172.71.") ||
	   strings.HasPrefix(ip, "162.158.") || strings.HasPrefix(ip, "104.") {
		return "Cloud", "CDN/Cloud"
	}
	
	// Simple geographic patterns (very basic, real implementation would use GeoIP database)
	switch {
	case strings.HasPrefix(ip, "203."):
		return "Australia/NZ", "Oceania"
	case strings.HasPrefix(ip, "202."):
		return "Asia", "Asia"
	case strings.HasPrefix(ip, "80.") || strings.HasPrefix(ip, "81."):
		return "Europe", "Europe"
	case strings.HasPrefix(ip, "24.") || strings.HasPrefix(ip, "76."):
		return "United States", "North America"
	case strings.HasPrefix(ip, "201."):
		return "Brazil", "South America"
	default:
		return "Unknown", "Unknown"
	}
}

func (a *Analyser) getRegionForCountry(country string) string {
	switch country {
	case "United States", "Canada", "Mexico":
		return "North America"
	case "Brazil", "Argentina", "Chile":
		return "South America"
	case "Germany", "France", "UK", "Spain", "Italy":
		return "Europe"
	case "China", "Japan", "India", "Korea", "Asia":
		return "Asia"
	case "Australia/NZ", "Australia", "New Zealand":
		return "Oceania"
	case "South Africa", "Nigeria", "Egypt":
		return "Africa"
	default:
		return "Unknown"
	}
}

func getStatusClass(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "2xx Success"
	case status >= 300 && status < 400:
		return "3xx Redirect"
	case status >= 400 && status < 500:
		return "4xx Client Error"
	case status >= 500:
		return "5xx Server Error"
	default:
		return "1xx Informational"
	}
}

// Security Analysis Methods
func (a *Analyser) analyseSecurityThreats(logs []*parser.LogEntry) SecurityAnalysis {
	var threats []SecurityThreat
	var suspiciousIPs []IPThreatAnalysis
	var anomalies []AnomalyDetection
	
	// Counters for different attack types
	sqlInjection := 0
	xssAttempts := 0
	directoryTraversal := 0
	bruteForce := 0
	scanningActivity := 0
	
	// Track IP behavior for threat analysis
	ipStats := make(map[string]*IPThreatAnalysis)
	
	// Analyze each log entry for security threats
	for _, log := range logs {
		// Initialize IP stats if not exists
		if _, exists := ipStats[log.IP]; !exists {
			ipStats[log.IP] = &IPThreatAnalysis{
				IP:               log.IP,
				RequestCount:     0,
				ThreatScore:      0,
				ThreatCategories: []string{},
				FirstSeen:        log.Timestamp,
				LastSeen:         log.Timestamp,
				UniqueURLs:       0,
				ErrorRate:        0,
			}
		}
		
		ipStat := ipStats[log.IP]
		ipStat.RequestCount++
		ipStat.LastSeen = log.Timestamp
		
		// Check for SQL injection patterns
		if a.detectSQLInjection(log.URL) {
			threats = append(threats, SecurityThreat{
				Type:      "sql_injection",
				Pattern:   a.extractSQLPattern(log.URL),
				URL:       log.URL,
				IP:        log.IP,
				Timestamp: log.Timestamp,
				Severity:  "high",
				UserAgent: log.UserAgent,
			})
			sqlInjection++
			a.updateThreatScore(ipStat, "sql_injection", 30)
		}
		
		// Check for XSS attempts
		if a.detectXSS(log.URL) {
			threats = append(threats, SecurityThreat{
				Type:      "xss",
				Pattern:   a.extractXSSPattern(log.URL),
				URL:       log.URL,
				IP:        log.IP,
				Timestamp: log.Timestamp,
				Severity:  "medium",
				UserAgent: log.UserAgent,
			})
			xssAttempts++
			a.updateThreatScore(ipStat, "xss", 20)
		}
		
		// Check for directory traversal
		if a.detectDirectoryTraversal(log.URL) {
			threats = append(threats, SecurityThreat{
				Type:      "directory_traversal",
				Pattern:   a.extractTraversalPattern(log.URL),
				URL:       log.URL,
				IP:        log.IP,
				Timestamp: log.Timestamp,
				Severity:  "high",
				UserAgent: log.UserAgent,
			})
			directoryTraversal++
			a.updateThreatScore(ipStat, "directory_traversal", 25)
		}
		
		// Check for brute force attempts (multiple failed logins)
		if a.detectBruteForce(log.URL, log.Status) {
			bruteForce++
			a.updateThreatScore(ipStat, "brute_force", 15)
		}
		
		// Check for scanning activity
		if a.detectScanning(log.UserAgent, log.URL) {
			scanningActivity++
			a.updateThreatScore(ipStat, "scanner", 10)
		}
		
		// Track error rates for IP reputation
		if log.Status >= 400 {
			// Will calculate error rate later
		}
	}
	
	// Calculate IP threat scores and error rates
	for ip, stat := range ipStats {
		errorCount := 0
		uniqueURLs := make(map[string]bool)
		
		for _, log := range logs {
			if log.IP == ip {
				uniqueURLs[log.URL] = true
				if log.Status >= 400 {
					errorCount++
				}
			}
		}
		
		stat.UniqueURLs = len(uniqueURLs)
		if stat.RequestCount > 0 {
			stat.ErrorRate = float64(errorCount) / float64(stat.RequestCount) * 100
		}
		
		// Only include IPs with suspicious activity
		if stat.ThreatScore > 0 {
			suspiciousIPs = append(suspiciousIPs, *stat)
		}
	}
	
	// Sort suspicious IPs by threat score
	sort.Slice(suspiciousIPs, func(i, j int) bool {
		return suspiciousIPs[i].ThreatScore > suspiciousIPs[j].ThreatScore
	})
	
	// Generate anomaly detection
	anomalies = a.detectAnomalies(logs)
	
	// Calculate overall threat level and security score
	threatLevel := a.calculateThreatLevel(threats, suspiciousIPs)
	securityScore := a.calculateSecurityScore(len(logs), len(threats), len(suspiciousIPs))
	
	// Create top attackers list
	topAttackers := []IPStat{}
	for i, ip := range suspiciousIPs {
		if i >= 10 { // Top 10 attackers
			break
		}
		topAttackers = append(topAttackers, IPStat{
			IP:    ip.IP,
			Count: ip.RequestCount,
		})
	}
	
	return SecurityAnalysis{
		ThreatLevel:          threatLevel,
		SecurityScore:        securityScore,
		TotalThreats:         len(threats),
		ThreatsDetected:      threats,
		SuspiciousIPs:        suspiciousIPs,
		AnomaliesDetected:    anomalies,
		BruteForceAttempts:   bruteForce,
		SQLInjectionAttempts: sqlInjection,
		XSSAttempts:          xssAttempts,
		DirectoryTraversal:   directoryTraversal,
		ScanningActivity:     scanningActivity,
		TopAttackers:         topAttackers,
	}
}

// Attack pattern detection methods
func (a *Analyser) detectSQLInjection(url string) bool {
	sqlPatterns := []string{
		"'", "\"", ";", "--", "/*", "*/",
		"union", "select", "insert", "update", "delete", "drop",
		"exec", "execute", "sp_", "xp_",
		"or 1=1", "or 1=1--", "or 'a'='a", "1' or '1'='1",
		"admin'--", "admin'/*", "' or 1=1#", "' or 1=1--",
	}
	
	urlLower := strings.ToLower(url)
	for _, pattern := range sqlPatterns {
		if strings.Contains(urlLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func (a *Analyser) detectXSS(url string) bool {
	xssPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:",
		"onload=", "onerror=", "onclick=", "onmouseover=",
		"<img", "<iframe", "<object", "<embed",
		"alert(", "document.cookie", "document.write",
		"eval(", "setTimeout(", "setInterval(",
	}
	
	urlLower := strings.ToLower(url)
	for _, pattern := range xssPatterns {
		if strings.Contains(urlLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func (a *Analyser) detectDirectoryTraversal(url string) bool {
	traversalPatterns := []string{
		"../", "..\\", "....//", "....\\\\",
		"%2e%2e/", "%2e%2e\\", "%2e%2e%2f", "%2e%2e%5c",
		"..%2f", "..%5c", "..\\/", "../\\",
		"/etc/passwd", "/etc/shadow", "\\windows\\system32",
		"boot.ini", "win.ini",
	}
	
	urlLower := strings.ToLower(url)
	for _, pattern := range traversalPatterns {
		if strings.Contains(urlLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func (a *Analyser) detectBruteForce(url string, status int) bool {
	// Look for login/admin URLs with failed status codes
	loginPaths := []string{
		"login", "admin", "signin", "auth", "wp-admin",
		"administrator", "panel", "dashboard",
	}
	
	urlLower := strings.ToLower(url)
	for _, path := range loginPaths {
		if strings.Contains(urlLower, path) && (status == 401 || status == 403 || status == 404) {
			return true
		}
	}
	return false
}

func (a *Analyser) detectScanning(userAgent string, url string) bool {
	scannerPatterns := []string{
		"nmap", "nikto", "sqlmap", "burp", "owasp zap",
		"nessus", "openvas", "acunetix", "qualys",
		"masscan", "zap", "w3af", "skipfish",
		"gobuster", "dirb", "dirbuster", "wfuzz",
	}
	
	agentLower := strings.ToLower(userAgent)
	for _, pattern := range scannerPatterns {
		if strings.Contains(agentLower, pattern) {
			return true
		}
	}
	
	// Check for common scanning URLs
	scanUrls := []string{
		"/admin", "/test", "/backup", "/.git", "/.svn",
		"/config", "/database", "/db", "/phpmyadmin",
		"/wp-config", "/robots.txt", "/sitemap.xml",
	}
	
	urlLower := strings.ToLower(url)
	for _, scanUrl := range scanUrls {
		if strings.Contains(urlLower, scanUrl) {
			return true
		}
	}
	
	return false
}

// Helper methods for pattern extraction
func (a *Analyser) extractSQLPattern(url string) string {
	if strings.Contains(strings.ToLower(url), "union") {
		return "UNION-based injection"
	}
	if strings.Contains(url, "' or 1=1") {
		return "Boolean-based injection"
	}
	if strings.Contains(url, "'--") || strings.Contains(url, "'/*") {
		return "Comment-based injection"
	}
	return "Generic SQL injection pattern"
}

func (a *Analyser) extractXSSPattern(url string) string {
	if strings.Contains(strings.ToLower(url), "script") {
		return "Script injection"
	}
	if strings.Contains(strings.ToLower(url), "javascript:") {
		return "JavaScript protocol"
	}
	if strings.Contains(strings.ToLower(url), "onerror") || strings.Contains(strings.ToLower(url), "onload") {
		return "Event handler injection"
	}
	return "Generic XSS pattern"
}

func (a *Analyser) extractTraversalPattern(url string) string {
	if strings.Contains(url, "../") {
		return "Unix-style traversal (../)"
	}
	if strings.Contains(url, "..\\") {
		return "Windows-style traversal (..\\)"
	}
	if strings.Contains(url, "%2e%2e") {
		return "URL-encoded traversal"
	}
	return "Generic directory traversal"
}

// Threat scoring and reputation
func (a *Analyser) updateThreatScore(ipStat *IPThreatAnalysis, threatType string, score int) {
	ipStat.ThreatScore += score
	
	// Add threat category if not already present
	found := false
	for _, category := range ipStat.ThreatCategories {
		if category == threatType {
			found = true
			break
		}
	}
	if !found {
		ipStat.ThreatCategories = append(ipStat.ThreatCategories, threatType)
	}
}

// Anomaly detection
func (a *Analyser) detectAnomalies(logs []*parser.LogEntry) []AnomalyDetection {
	var anomalies []AnomalyDetection
	
	if len(logs) == 0 {
		return anomalies
	}
	
	// Calculate baseline metrics
	totalRequests := len(logs)
	errorCount := 0
	statusCodes := make(map[int]int)
	
	for _, log := range logs {
		statusCodes[log.Status]++
		if log.Status >= 400 {
			errorCount++
		}
	}
	
	// Check for anomalous error rates
	errorRate := float64(errorCount) / float64(totalRequests) * 100
	expectedErrorRate := 5.0 // 5% is typical baseline
	
	if errorRate > expectedErrorRate*2 { // 2x expected rate
		anomalies = append(anomalies, AnomalyDetection{
			Type:          "high_error_rate",
			Description:   "Unusually high error rate detected",
			Value:         errorRate,
			Expected:      expectedErrorRate,
			Deviation:     (errorRate - expectedErrorRate) / expectedErrorRate * 100,
			Significance:  a.getSignificance(errorRate, expectedErrorRate, 2.0),
		})
	}
	
	// Check for anomalous 404 rates
	notFoundCount := statusCodes[404]
	notFoundRate := float64(notFoundCount) / float64(totalRequests) * 100
	expectedNotFoundRate := 2.0 // 2% is typical
	
	if notFoundRate > expectedNotFoundRate*3 {
		anomalies = append(anomalies, AnomalyDetection{
			Type:          "high_404_rate",
			Description:   "Unusually high 404 Not Found rate - possible scanning activity",
			Value:         notFoundRate,
			Expected:      expectedNotFoundRate,
			Deviation:     (notFoundRate - expectedNotFoundRate) / expectedNotFoundRate * 100,
			Significance:  a.getSignificance(notFoundRate, expectedNotFoundRate, 3.0),
		})
	}
	
	return anomalies
}

func (a *Analyser) getSignificance(actual, expected, threshold float64) string {
	ratio := actual / expected
	if ratio > threshold*2 {
		return "high"
	} else if ratio > threshold*1.5 {
		return "medium"
	}
	return "low"
}

// Calculate overall threat level and security score
func (a *Analyser) calculateThreatLevel(threats []SecurityThreat, suspiciousIPs []IPThreatAnalysis) string {
	highSeverityCount := 0
	mediumSeverityCount := 0
	
	for _, threat := range threats {
		switch threat.Severity {
		case "critical", "high":
			highSeverityCount++
		case "medium":
			mediumSeverityCount++
		}
	}
	
	topThreatIPs := 0
	for _, ip := range suspiciousIPs {
		if ip.ThreatScore > 50 {
			topThreatIPs++
		}
	}
	
	if highSeverityCount > 10 || topThreatIPs > 5 {
		return "critical"
	} else if highSeverityCount > 5 || mediumSeverityCount > 10 || topThreatIPs > 2 {
		return "high"
	} else if highSeverityCount > 0 || mediumSeverityCount > 0 || topThreatIPs > 0 {
		return "medium"
	}
	
	return "low"
}

func (a *Analyser) calculateSecurityScore(totalRequests, threatCount, suspiciousIPCount int) int {
	if totalRequests == 0 {
		return 100
	}
	
	// Start with perfect score
	score := 100
	
	// Deduct points for threats
	threatRate := float64(threatCount) / float64(totalRequests) * 100
	score -= int(threatRate * 2) // Each 1% threat rate costs 2 points
	
	// Deduct points for suspicious IPs
	suspiciousRate := float64(suspiciousIPCount) / float64(totalRequests) * 100
	score -= int(suspiciousRate * 1.5) // Each 1% suspicious IP rate costs 1.5 points
	
	// Minimum score is 0
	if score < 0 {
		score = 0
	}
	
	return score
}