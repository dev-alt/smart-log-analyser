package html

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"smart-log-analyser/pkg/analyser"
)

//go:embed templates/*
var templateFS embed.FS

// ReportData contains all data needed for HTML report generation
type ReportData struct {
	// Header Information
	Title              string
	GeneratedAt        string
	DateRange          string
	Version            string
	ReportID           string
	AnalysisDuration   string

	// Overview Metrics
	TotalRequests       string
	UniqueIPs          string
	DataTransferred    string
	AverageResponseSize string

	// Traffic Analysis
	HumanTraffic int
	BotTraffic   int

	// Hourly Traffic Data
	HourlyLabels []string
	HourlyData   []int

	// Status Code Data
	StatusLabels []string
	StatusData   []int

	// Response Size Data (in KB)
	P50Size float64
	P95Size float64
	P99Size float64
	AvgSize float64

	// Geographic Data
	GeoLabels []string
	GeoData   []int

	// File Type Data
	FileTypeLabels []string
	FileTypeData   []int

	// Security Data
	SecurityScore  string
	SecurityClass  string
	TotalThreats   int
	SuspiciousIPs  int

	// Tables Data
	TopIPs   []IPRow
	TopURLs  []URLRow
	ErrorURLs []ErrorRow
}

// IPRow represents a row in the top IPs table
type IPRow struct {
	IP         string
	Count      int
	Percentage string
	Location   string
	Type       string
	TypeClass  string
}

// URLRow represents a row in the top URLs table
type URLRow struct {
	URL           string
	Count         int
	Percentage    string
	AverageSize   string
	TotalTransfer string
}

// ErrorRow represents a row in the error analysis table
type ErrorRow struct {
	URL         string
	ErrorCount  int
	StatusCodes string
	ErrorRate   string
}

// Generator handles HTML report generation
type Generator struct {
	template *template.Template
}

// NewGenerator creates a new HTML report generator
func NewGenerator() (*Generator, error) {
	// Create custom template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"formatBytes": formatBytes,
	}

	// Parse embedded template
	tmpl, err := template.New("report.html").Funcs(funcMap).ParseFS(templateFS, "templates/report.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Generator{
		template: tmpl,
	}, nil
}

// GenerateReport creates an HTML report from analysis results
func (g *Generator) GenerateReport(results *analyser.Results, outputPath string, title string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Transform analysis results to report data
	reportData := g.transformResults(results, title)

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := g.template.Execute(file, reportData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// transformResults converts analyser.Results to ReportData
func (g *Generator) transformResults(results *analyser.Results, title string) *ReportData {
	now := time.Now()

	// Generate unique report ID
	reportID := fmt.Sprintf("SLA-%d", now.Unix())

	// Calculate traffic data
	humanTraffic := results.HumanRequests
	botTraffic := results.BotRequests
	
	// Format date range
	dateRange := "N/A"
	if !results.TimeRange.Start.IsZero() && !results.TimeRange.End.IsZero() {
		dateRange = fmt.Sprintf("%s to %s", 
			results.TimeRange.Start.Format("2006-01-02 15:04"),
			results.TimeRange.End.Format("2006-01-02 15:04"))
	}

	// Prepare hourly data
	hourlyLabels := make([]string, 0)
	hourlyData := make([]int, 0)
	for _, hourly := range results.HourlyTraffic {
		hourlyLabels = append(hourlyLabels, fmt.Sprintf("%02d:00", hourly.Hour))
		hourlyData = append(hourlyData, hourly.RequestCount)
	}

	// Prepare status code data from map (only include non-zero values)
	statusLabels := make([]string, 0)
	statusData := make([]int, 0)
	
	statusCategories := map[string]string{
		"2": "2xx Success",
		"3": "3xx Redirect", 
		"4": "4xx Client Error",
		"5": "5xx Server Error",
	}
	
	for code, label := range statusCategories {
		count := getStatusCodeCount(results.StatusCodes, code+"xx")
		if count > 0 {
			statusLabels = append(statusLabels, label)
			statusData = append(statusData, count)
		}
	}

	// Prepare geographic data (only include non-zero values)
	geoLabels := make([]string, 0)
	geoData := make([]int, 0)
	
	if results.GeographicAnalysis.LocalTraffic > 0 {
		geoLabels = append(geoLabels, "Local/Private")
		geoData = append(geoData, results.GeographicAnalysis.LocalTraffic)
	}
	if results.GeographicAnalysis.CloudTraffic > 0 {
		geoLabels = append(geoLabels, "CDN/Cloud") 
		geoData = append(geoData, results.GeographicAnalysis.CloudTraffic)
	}
	if results.GeographicAnalysis.UnknownIPs > 0 {
		geoLabels = append(geoLabels, "Unknown")
		geoData = append(geoData, results.GeographicAnalysis.UnknownIPs)
	}

	// Prepare file type data
	fileTypeLabels := make([]string, 0)
	fileTypeData := make([]int, 0)
	for _, fileType := range results.FileTypes {
		if len(fileTypeLabels) < 6 { // Limit to top 6 file types
			fileTypeLabels = append(fileTypeLabels, fileType.FileType)
			fileTypeData = append(fileTypeData, fileType.Count)
		}
	}

	// Prepare top IPs
	topIPs := make([]IPRow, 0)
	for i, ip := range results.TopIPs {
		if i >= 10 { // Limit to top 10
			break
		}

		location := getLocationFromIP(ip.IP)
		ipType, typeClass := getIPTypeAndClass(ip.IP)
		
		topIPs = append(topIPs, IPRow{
			IP:         ip.IP,
			Count:      ip.Count,
			Percentage: fmt.Sprintf("%.1f", float64(ip.Count*100)/float64(results.TotalRequests)),
			Location:   location,
			Type:       ipType,
			TypeClass:  typeClass,
		})
	}

	// Prepare top URLs
	topURLs := make([]URLRow, 0)
	for i, url := range results.TopURLs {
		if i >= 10 { // Limit to top 10
			break
		}

		topURLs = append(topURLs, URLRow{
			URL:           truncateURL(url.URL, 80),
			Count:         url.Count,
			Percentage:    fmt.Sprintf("%.1f", float64(url.Count*100)/float64(results.TotalRequests)),
			AverageSize:   "N/A", // TODO: Calculate from results if available
			TotalTransfer: "N/A", // TODO: Calculate from results if available
		})
	}

	// Prepare error URLs from ErrorURLs field
	errorURLs := make([]ErrorRow, 0)
	for _, errorURL := range results.ErrorURLs {
		if len(errorURLs) >= 5 { // Limit to top 5 error URLs
			break
		}
		
		errorURLs = append(errorURLs, ErrorRow{
			URL:         truncateURL(errorURL.URL, 60),
			ErrorCount:  errorURL.Count,
			StatusCodes: errorURL.FormatStatusCodes(),
			ErrorRate:   fmt.Sprintf("%.1f", float64(errorURL.Count*100)/float64(results.TotalRequests)),
		})
	}

	// Determine security class
	securityClass := "security-low"
	if results.SecurityAnalysis.SecurityScore < 70 {
		securityClass = "security-high"
	} else if results.SecurityAnalysis.SecurityScore < 85 {
		securityClass = "security-medium"
	}

	return &ReportData{
		Title:              title,
		GeneratedAt:        now.Format("2006-01-02 15:04:05"),
		DateRange:          dateRange,
		Version:            "1.0.0", // TODO: Get from build info
		ReportID:           reportID,
		AnalysisDuration:   "N/A", // TODO: Add timing to results

		TotalRequests:       formatNumber(results.TotalRequests),
		UniqueIPs:          formatNumber(results.UniqueIPs),
		DataTransferred:    formatBytes(results.TotalBytes),
		AverageResponseSize: formatBytes(results.AverageSize),

		HumanTraffic: humanTraffic,
		BotTraffic:   botTraffic,

		HourlyLabels: hourlyLabels,
		HourlyData:   hourlyData,

		StatusLabels: statusLabels,
		StatusData:   statusData,

		P50Size: float64(results.ResponseTimeStats.MedianSize) / 1024,
		P95Size: float64(results.ResponseTimeStats.P95Size) / 1024,
		P99Size: float64(results.ResponseTimeStats.P99Size) / 1024,
		AvgSize: float64(results.ResponseTimeStats.AverageSize) / 1024,

		GeoLabels: geoLabels,
		GeoData:   geoData,

		FileTypeLabels: fileTypeLabels,
		FileTypeData:   fileTypeData,

		SecurityScore:  fmt.Sprintf("%d/100", getSecurityScore(results)),
		SecurityClass:  securityClass,
		TotalThreats:   getTotalThreats(results),
		SuspiciousIPs:  getSuspiciousIPCount(results),

		TopIPs:    topIPs,
		TopURLs:   topURLs,
		ErrorURLs: errorURLs,
	}
}

// Helper functions

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%.1f GB", float64(bytes)/(1024*1024*1024))
}

func formatNumber(num int) string {
	str := strconv.Itoa(num)
	if len(str) <= 3 {
		return str
	}

	result := ""
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(char)
	}
	return result
}

func truncateURL(url string, maxLength int) string {
	if len(url) <= maxLength {
		return url
	}
	return url[:maxLength-3] + "..."
}

func getLocationFromIP(ip string) string {
	// Simple pattern-based location detection
	if strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || 
	   strings.HasPrefix(ip, "172.") {
		return "Local Network"
	}
	if strings.HasPrefix(ip, "172.69.") || strings.HasPrefix(ip, "172.70.") ||
	   strings.HasPrefix(ip, "172.71.") {
		return "Cloudflare CDN"
	}
	return "Unknown"
}

func getIPTypeAndClass(ip string) (string, string) {
	if strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") {
		return "Private", "secondary"
	}
	if strings.HasPrefix(ip, "172.69.") || strings.HasPrefix(ip, "172.70.") {
		return "CDN", "info"
	}
	return "Public", "primary"
}

func getStatusCodeCount(statusCodes map[string]int, category string) int {
	count := 0
	for code, num := range statusCodes {
		if strings.HasPrefix(code, category[:1]) { // Match first character (2, 3, 4, 5)
			count += num
		}
	}
	return count
}

func getSecurityScore(results *analyser.Results) int {
	return results.SecurityAnalysis.SecurityScore
}

func getTotalThreats(results *analyser.Results) int {
	return results.SecurityAnalysis.TotalThreats
}

func getSuspiciousIPCount(results *analyser.Results) int {
	return len(results.SecurityAnalysis.SuspiciousIPs)
}