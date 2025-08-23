package analyser

import (
	"sort"
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

type Results struct {
	TotalRequests int
	TimeRange     TimeRange
	StatusCodes   map[string]int
	TopIPs        []IPStat
	TopURLs       []URLStat
	HTTPMethods   []MethodStat
	TotalBytes    int64
	AverageSize   int64
	UniqueIPs     int
	UniqueURLs    int
}

type Analyser struct{}

func New() *Analyser {
	return &Analyser{}
}

func (a *Analyser) Analyse(logs []*parser.LogEntry, since, until *time.Time) *Results {
	filtered := a.filterByTime(logs, since, until)
	
	if len(filtered) == 0 {
		return &Results{
			TotalRequests: 0,
			TimeRange:     TimeRange{},
			StatusCodes:   make(map[string]int),
			TopIPs:        []IPStat{},
			TopURLs:       []URLStat{},
			HTTPMethods:   []MethodStat{},
			TotalBytes:    0,
			AverageSize:   0,
			UniqueIPs:     0,
			UniqueURLs:    0,
		}
	}

	results := &Results{
		TotalRequests: len(filtered),
		TimeRange:     a.calculateTimeRange(filtered),
		StatusCodes:   a.analyseStatusCodes(filtered),
		TopIPs:        a.analyseTopIPs(filtered),
		TopURLs:       a.analyseTopURLs(filtered),
		HTTPMethods:   a.analyseHTTPMethods(filtered),
		TotalBytes:    a.calculateTotalBytes(filtered),
		AverageSize:   a.calculateAverageSize(filtered),
		UniqueIPs:     a.countUniqueIPs(filtered),
		UniqueURLs:    a.countUniqueURLs(filtered),
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