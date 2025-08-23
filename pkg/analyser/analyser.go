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

type Results struct {
	TotalRequests int
	TimeRange     TimeRange
	StatusCodes   map[string]int
	TopIPs        []IPStat
	TopURLs       []URLStat
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
		}
	}

	results := &Results{
		TotalRequests: len(filtered),
		TimeRange:     a.calculateTimeRange(filtered),
		StatusCodes:   a.analyseStatusCodes(filtered),
		TopIPs:        a.analyseTopIPs(filtered),
		TopURLs:       a.analyseTopURLs(filtered),
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