package analyser

import (
	"sort"
	"strings"
	"time"

	"smart-log-analyser/pkg/parser"
)

type IPStat struct {
	IP    string
	Count int
}

type URLStat struct {
	URL   string
	Count int
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

type DetailedStatusCode struct {
	Code  int
	Count int
}

type Results struct {
	TotalRequests       int
	TimeRange           TimeRange
	StatusCodes         map[string]int
	DetailedStatusCodes []DetailedStatusCode
	TopIPs              []IPStat
	TopURLs             []URLStat
	HTTPMethods         []MethodStat
	TotalBytes          int64
	AverageSize         int64
	UniqueIPs           int
	UniqueURLs          int
	BotRequests         int
	HumanRequests       int
	TopBots             []BotStat
	FileTypes           []FileTypeStat
	ErrorURLs           []URLStat // URLs that generated errors
	LargeRequests       []URLStat // Largest requests by size
}

type Analyser struct{}

func New() *Analyser {
	return &Analyser{}
}

func (a *Analyser) Analyse(logs []*parser.LogEntry, since, until *time.Time) *Results {
	filtered := a.filterByTime(logs, since, until)
	
	if len(filtered) == 0 {
		return &Results{
			TotalRequests:       0,
			TimeRange:           TimeRange{},
			StatusCodes:         make(map[string]int),
			DetailedStatusCodes: []DetailedStatusCode{},
			TopIPs:              []IPStat{},
			TopURLs:             []URLStat{},
			HTTPMethods:         []MethodStat{},
			TotalBytes:          0,
			AverageSize:         0,
			UniqueIPs:           0,
			UniqueURLs:          0,
			BotRequests:         0,
			HumanRequests:       0,
			TopBots:             []BotStat{},
			FileTypes:           []FileTypeStat{},
			ErrorURLs:           []URLStat{},
			LargeRequests:       []URLStat{},
		}
	}

	botRequests, humanRequests := a.analyseBotTraffic(filtered)
	
	results := &Results{
		TotalRequests:       len(filtered),
		TimeRange:           a.calculateTimeRange(filtered),
		StatusCodes:         a.analyseStatusCodes(filtered),
		DetailedStatusCodes: a.analyseDetailedStatusCodes(filtered),
		TopIPs:              a.analyseTopIPs(filtered),
		TopURLs:             a.analyseTopURLs(filtered),
		HTTPMethods:         a.analyseHTTPMethods(filtered),
		TotalBytes:          a.calculateTotalBytes(filtered),
		AverageSize:         a.calculateAverageSize(filtered),
		UniqueIPs:           a.countUniqueIPs(filtered),
		UniqueURLs:          a.countUniqueURLs(filtered),
		BotRequests:         botRequests,
		HumanRequests:       humanRequests,
		TopBots:             a.analyseTopBots(filtered),
		FileTypes:           a.analyseFileTypes(filtered),
		ErrorURLs:           a.analyseErrorURLs(filtered),
		LargeRequests:       a.analyseLargeRequests(filtered),
	}

	return results
}

func (a *Analyser) filterByTime(logs []*parser.LogEntry, since, until *time.Time) []*parser.LogEntry {
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
		urlStats = append(urlStats, URLStat{URL: url, Count: count})
	}

	sort.Slice(urlStats, func(i, j int) bool {
		return urlStats[i].Count > urlStats[j].Count
	})

	return urlStats
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
	errorCounts := make(map[string]int)
	
	for _, log := range logs {
		if log.Status >= 400 { // 4xx and 5xx errors
			errorCounts[log.URL]++
		}
	}
	
	var errorStats []URLStat
	for url, count := range errorCounts {
		errorStats = append(errorStats, URLStat{URL: url, Count: count})
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
				URL:   req.url,
				Count: int(req.size), // Store size in count field for display
			})
			seen[req.url] = true
		}
	}
	
	return largeStats
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