package security

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"smart-log-analyser/pkg/parser"
)

// ThreatDetector implements advanced threat detection algorithms
type ThreatDetector struct {
	config              SecurityConfig
	webAttackPatterns   map[WebAttackType][]*regexp.Regexp
	infraAttackPatterns map[InfrastructureAttackType][]*regexp.Regexp
	knownBadIPs         map[string]ThreatInfo
	suspiciousPatterns  []*regexp.Regexp
	threatIntelligence  *ThreatIntelligence
}

// NewThreatDetector creates a new threat detector with configured patterns
func NewThreatDetector(config SecurityConfig) *ThreatDetector {
	td := &ThreatDetector{
		config:              config,
		webAttackPatterns:   make(map[WebAttackType][]*regexp.Regexp),
		infraAttackPatterns: make(map[InfrastructureAttackType][]*regexp.Regexp),
		knownBadIPs:         make(map[string]ThreatInfo),
		threatIntelligence:  &ThreatIntelligence{
			MaliciousIPs:     make(map[string]ThreatInfo),
			AttackSignatures: []AttackSignature{},
			KnownPayloads:    make(map[string]PayloadInfo),
			VulnerabilityDB:  []VulnerabilityInfo{},
		},
	}

	td.initializePatterns()
	td.loadThreatIntelligence()
	return td
}

// DetectWebAttacks identifies web application attacks in log entries
func (td *ThreatDetector) DetectWebAttacks(logs []*parser.LogEntry) ([]EnhancedThreat, error) {
	var threats []EnhancedThreat

	for _, entry := range logs {
		// SQL Injection Detection
		if sqlThreats := td.detectSQLInjection(entry); len(sqlThreats) > 0 {
			threats = append(threats, sqlThreats...)
		}

		// Cross-Site Scripting Detection
		if xssThreats := td.detectXSS(entry); len(xssThreats) > 0 {
			threats = append(threats, xssThreats...)
		}

		// Command Injection Detection
		if cmdThreats := td.detectCommandInjection(entry); len(cmdThreats) > 0 {
			threats = append(threats, cmdThreats...)
		}

		// Directory Traversal Detection
		if dirThreats := td.detectDirectoryTraversal(entry); len(dirThreats) > 0 {
			threats = append(threats, dirThreats...)
		}

		// File Inclusion Detection
		if fileThreats := td.detectFileInclusion(entry); len(fileThreats) > 0 {
			threats = append(threats, fileThreats...)
		}

		// XXE Injection Detection
		if xxeThreats := td.detectXXEInjection(entry); len(xxeThreats) > 0 {
			threats = append(threats, xxeThreats...)
		}

		// HTTP Header Injection Detection
		if headerThreats := td.detectHeaderInjection(entry); len(headerThreats) > 0 {
			threats = append(threats, headerThreats...)
		}
	}

	return threats, nil
}

// DetectInfrastructureAttacks identifies infrastructure-level attacks
func (td *ThreatDetector) DetectInfrastructureAttacks(logs []*parser.LogEntry) ([]EnhancedThreat, error) {
	var threats []EnhancedThreat

	// Group entries by IP for behavioral analysis
	ipEntries := make(map[string][]*parser.LogEntry)
	for _, entry := range logs {
		ipEntries[entry.IP] = append(ipEntries[entry.IP], entry)
	}

	for ip, entries := range ipEntries {
		// Brute Force Detection
		if bruteThreats := td.detectBruteForce(ip, entries); len(bruteThreats) > 0 {
			threats = append(threats, bruteThreats...)
		}

		// DDoS Detection
		if ddosThreats := td.detectDDoS(ip, entries); len(ddosThreats) > 0 {
			threats = append(threats, ddosThreats...)
		}

		// Port Scanning Detection
		if scanThreats := td.detectPortScanning(ip, entries); len(scanThreats) > 0 {
			threats = append(threats, scanThreats...)
		}

		// Vulnerability Scanning Detection
		if vulnThreats := td.detectVulnerabilityScanning(ip, entries); len(vulnThreats) > 0 {
			threats = append(threats, vulnThreats...)
		}

		// Bot Activity Detection
		if botThreats := td.detectBotActivity(ip, entries); len(botThreats) > 0 {
			threats = append(threats, botThreats...)
		}
	}

	return threats, nil
}

// detectSQLInjection detects SQL injection attempts
func (td *ThreatDetector) detectSQLInjection(entry *parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	// SQL injection patterns with different severity levels
	sqlPatterns := []struct {
		pattern  *regexp.Regexp
		severity ThreatSeverity
		desc     string
	}{
		{regexp.MustCompile(`(?i)(union\s+select|union\s+all\s+select)`), SeverityHigh, "UNION-based SQL injection"},
		{regexp.MustCompile(`(?i)(select\s+.*\s+from\s+|insert\s+into\s+|update\s+.*\s+set\s+|delete\s+from\s+)`), SeverityMedium, "SQL query injection"},
		{regexp.MustCompile(`(?i)(\'\s*or\s*\'1\'\s*=\s*\'1|\'\s*or\s*1\s*=\s*1|admin\'\s*--)`), SeverityHigh, "Boolean-based SQL injection"},
		{regexp.MustCompile(`(?i)(sleep\s*\(|benchmark\s*\(|pg_sleep\s*\()`), SeverityMedium, "Time-based SQL injection"},
		{regexp.MustCompile(`(?i)(drop\s+table|drop\s+database|truncate\s+table)`), SeverityCritical, "Destructive SQL injection"},
		{regexp.MustCompile(`(?i)(xp_cmdshell|sp_executesql|exec\s*\()`), SeverityCritical, "SQL command execution"},
		{regexp.MustCompile(`(?i)(\'\s*;\s*exec|\'\s*;\s*declare)`), SeverityHigh, "Stacked SQL injection"},
	}

	target := entry.URL + " " + entry.UserAgent + " " + entry.Referer

	for _, sqlPattern := range sqlPatterns {
		if sqlPattern.pattern.MatchString(target) {
			payload := sqlPattern.pattern.FindString(target)
			threat := EnhancedThreat{
				ID:               fmt.Sprintf("sql_%d_%s", time.Now().UnixNano(), entry.IP),
				Type:             SQLInjection,
				Severity:         sqlPattern.severity,
				Confidence:       td.calculateConfidence(sqlPattern.severity, payload),
				Pattern:          sqlPattern.pattern.String(),
				URL:              entry.URL,
				IP:               entry.IP,
				UserAgent:        entry.UserAgent,
				Timestamp:        entry.Timestamp,
				Method:           entry.Method,
				StatusCode:       entry.Status,
				ResponseSize:     entry.Size,
				AttackVector:     "HTTP Request",
				Payload:          payload,
				Context:          map[string]interface{}{"description": sqlPattern.desc},
				MitigationAdvice: []string{"Implement parameterized queries", "Use input validation", "Apply principle of least privilege"},
			}
			threats = append(threats, threat)
		}
	}

	return threats
}

// detectXSS detects Cross-Site Scripting attacks
func (td *ThreatDetector) detectXSS(entry *parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	xssPatterns := []struct {
		pattern  *regexp.Regexp
		severity ThreatSeverity
		desc     string
	}{
		{regexp.MustCompile(`(?i)(<script[^>]*>|</script>)`), SeverityHigh, "Script tag injection"},
		{regexp.MustCompile(`(?i)(javascript:|vbscript:|data:text/html)`), SeverityHigh, "Protocol-based XSS"},
		{regexp.MustCompile(`(?i)(onload\s*=|onclick\s*=|onerror\s*=|onmouseover\s*=)`), SeverityMedium, "Event handler injection"},
		{regexp.MustCompile(`(?i)(<iframe|<object|<embed|<applet)`), SeverityMedium, "Object embedding XSS"},
		{regexp.MustCompile(`(?i)(alert\s*\(|confirm\s*\(|prompt\s*\()`), SeverityMedium, "Dialog-based XSS"},
		{regexp.MustCompile(`(?i)(document\.cookie|document\.location|window\.location)`), SeverityHigh, "DOM manipulation XSS"},
		{regexp.MustCompile(`(?i)(<img[^>]*src\s*=\s*[\"']?javascript:)`), SeverityMedium, "Image-based XSS"},
	}

	target := entry.URL + " " + entry.UserAgent + " " + entry.Referer

	for _, xssPattern := range xssPatterns {
		if xssPattern.pattern.MatchString(target) {
			payload := xssPattern.pattern.FindString(target)
			threat := EnhancedThreat{
				ID:               fmt.Sprintf("xss_%d_%s", time.Now().UnixNano(), entry.IP),
				Type:             CrossSiteScripting,
				Severity:         xssPattern.severity,
				Confidence:       td.calculateConfidence(xssPattern.severity, payload),
				Pattern:          xssPattern.pattern.String(),
				URL:              entry.URL,
				IP:               entry.IP,
				UserAgent:        entry.UserAgent,
				Timestamp:        entry.Timestamp,
				Method:           entry.Method,
				StatusCode:       entry.Status,
				ResponseSize:     entry.Size,
				AttackVector:     "HTTP Request",
				Payload:          payload,
				Context:          map[string]interface{}{"description": xssPattern.desc},
				MitigationAdvice: []string{"Implement output encoding", "Use Content Security Policy", "Validate input data"},
			}
			threats = append(threats, threat)
		}
	}

	return threats
}

// detectCommandInjection detects command injection attempts
func (td *ThreatDetector) detectCommandInjection(entry *parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	cmdPatterns := []struct {
		pattern  *regexp.Regexp
		severity ThreatSeverity
		desc     string
	}{
		{regexp.MustCompile(`(?i)(;|\||&|&&|\$\(|` + "`" + `)`), SeverityMedium, "Command chaining operators"},
		{regexp.MustCompile(`(?i)(wget\s+|curl\s+|nc\s+|netcat\s+)`), SeverityHigh, "Network command injection"},
		{regexp.MustCompile(`(?i)(cat\s+/etc/passwd|cat\s+/etc/shadow)`), SeverityCritical, "System file access"},
		{regexp.MustCompile(`(?i)(rm\s+-rf|del\s+/|format\s+)`), SeverityCritical, "Destructive commands"},
		{regexp.MustCompile(`(?i)(whoami|id\s+|ps\s+|netstat\s+|ifconfig)`), SeverityMedium, "System reconnaissance"},
		{regexp.MustCompile(`(?i)(python\s+-c|perl\s+-e|ruby\s+-e|php\s+-r)`), SeverityHigh, "Script execution"},
		{regexp.MustCompile(`(?i)(/bin/bash|/bin/sh|cmd\.exe|powershell)`), SeverityHigh, "Shell execution"},
	}

	target := entry.URL + " " + entry.UserAgent

	for _, cmdPattern := range cmdPatterns {
		if cmdPattern.pattern.MatchString(target) {
			payload := cmdPattern.pattern.FindString(target)
			threat := EnhancedThreat{
				ID:               fmt.Sprintf("cmd_%d_%s", time.Now().UnixNano(), entry.IP),
				Type:             CommandInjection,
				Severity:         cmdPattern.severity,
				Confidence:       td.calculateConfidence(cmdPattern.severity, payload),
				Pattern:          cmdPattern.pattern.String(),
				URL:              entry.URL,
				IP:               entry.IP,
				UserAgent:        entry.UserAgent,
				Timestamp:        entry.Timestamp,
				Method:           entry.Method,
				StatusCode:       entry.Status,
				ResponseSize:     entry.Size,
				AttackVector:     "HTTP Request",
				Payload:          payload,
				Context:          map[string]interface{}{"description": cmdPattern.desc},
				MitigationAdvice: []string{"Use parameterized system calls", "Implement input sanitization", "Apply command whitelisting"},
			}
			threats = append(threats, threat)
		}
	}

	return threats
}

// detectDirectoryTraversal detects directory traversal attacks
func (td *ThreatDetector) detectDirectoryTraversal(entry *parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	traversalPatterns := []struct {
		pattern  *regexp.Regexp
		severity ThreatSeverity
		desc     string
	}{
		{regexp.MustCompile(`\.\.\/|\.\.\\`), SeverityMedium, "Basic directory traversal"},
		{regexp.MustCompile(`\.\.%2f|\.\.%5c|%2e%2e%2f|%2e%2e%5c`), SeverityMedium, "URL-encoded traversal"},
		{regexp.MustCompile(`\.\.\/\.\.\/\.\.\/|\.\.\\\.\.\\\.\.\\`), SeverityHigh, "Deep directory traversal"},
		{regexp.MustCompile(`(?i)(\/etc\/passwd|\/etc\/shadow|\/windows\/system32)`), SeverityCritical, "System file access attempt"},
		{regexp.MustCompile(`(?i)(\.\.\/)+.*(passwd|shadow|hosts|httpd\.conf)`), SeverityCritical, "Configuration file access"},
	}

	for _, traversalPattern := range traversalPatterns {
		if traversalPattern.pattern.MatchString(entry.URL) {
			payload := traversalPattern.pattern.FindString(entry.URL)
			threat := EnhancedThreat{
				ID:               fmt.Sprintf("traversal_%d_%s", time.Now().UnixNano(), entry.IP),
				Type:             DirectoryTraversal,
				Severity:         traversalPattern.severity,
				Confidence:       td.calculateConfidence(traversalPattern.severity, payload),
				Pattern:          traversalPattern.pattern.String(),
				URL:              entry.URL,
				IP:               entry.IP,
				UserAgent:        entry.UserAgent,
				Timestamp:        entry.Timestamp,
				Method:           entry.Method,
				StatusCode:       entry.Status,
				ResponseSize:     entry.Size,
				AttackVector:     "HTTP Request",
				Payload:          payload,
				Context:          map[string]interface{}{"description": traversalPattern.desc},
				MitigationAdvice: []string{"Implement path validation", "Use chroot jails", "Apply file access controls"},
			}
			threats = append(threats, threat)
		}
	}

	return threats
}

// detectFileInclusion detects Local/Remote File Inclusion attacks
func (td *ThreatDetector) detectFileInclusion(entry *parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	inclusionPatterns := []struct {
		pattern  *regexp.Regexp
		severity ThreatSeverity
		desc     string
		attackType interface{}
	}{
		{regexp.MustCompile(`(?i)(http://|https://|ftp://)`), SeverityHigh, "Remote File Inclusion", RemoteFileInclusion},
		{regexp.MustCompile(`(?i)(file://|php://|zip://|data://)`), SeverityHigh, "Protocol-based inclusion", RemoteFileInclusion},
		{regexp.MustCompile(`(?i)(\/proc\/|\/dev\/|\/sys\/)`), SeverityMedium, "System file inclusion", LocalFileInclusion},
		{regexp.MustCompile(`(?i)(\.log|\.txt|\.php|\.asp|\.jsp)($|\?|&)`), SeverityMedium, "Local file inclusion", LocalFileInclusion},
	}

	for _, inclusionPattern := range inclusionPatterns {
		if inclusionPattern.pattern.MatchString(entry.URL) {
			payload := inclusionPattern.pattern.FindString(entry.URL)
			threat := EnhancedThreat{
				ID:               fmt.Sprintf("inclusion_%d_%s", time.Now().UnixNano(), entry.IP),
				Type:             inclusionPattern.attackType,
				Severity:         inclusionPattern.severity,
				Confidence:       td.calculateConfidence(inclusionPattern.severity, payload),
				Pattern:          inclusionPattern.pattern.String(),
				URL:              entry.URL,
				IP:               entry.IP,
				UserAgent:        entry.UserAgent,
				Timestamp:        entry.Timestamp,
				Method:           entry.Method,
				StatusCode:       entry.Status,
				ResponseSize:     entry.Size,
				AttackVector:     "HTTP Request",
				Payload:          payload,
				Context:          map[string]interface{}{"description": inclusionPattern.desc},
				MitigationAdvice: []string{"Whitelist allowed files", "Disable remote includes", "Validate file paths"},
			}
			threats = append(threats, threat)
		}
	}

	return threats
}

// detectXXEInjection detects XML External Entity injection
func (td *ThreatDetector) detectXXEInjection(entry *parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	xxePatterns := []struct {
		pattern  *regexp.Regexp
		severity ThreatSeverity
		desc     string
	}{
		{regexp.MustCompile(`(?i)<!ENTITY.*SYSTEM`), SeverityHigh, "XXE with SYSTEM entity"},
		{regexp.MustCompile(`(?i)<!ENTITY.*PUBLIC`), SeverityMedium, "XXE with PUBLIC entity"},
		{regexp.MustCompile(`(?i)(file://|http://|ftp://).*>]>`), SeverityHigh, "XXE with external resource"},
		{regexp.MustCompile(`(?i)<!DOCTYPE.*\[.*ENTITY`), SeverityMedium, "DOCTYPE with entity declaration"},
	}

	target := entry.URL + " " + entry.UserAgent

	for _, xxePattern := range xxePatterns {
		if xxePattern.pattern.MatchString(target) {
			payload := xxePattern.pattern.FindString(target)
			threat := EnhancedThreat{
				ID:               fmt.Sprintf("xxe_%d_%s", time.Now().UnixNano(), entry.IP),
				Type:             XXEInjection,
				Severity:         xxePattern.severity,
				Confidence:       td.calculateConfidence(xxePattern.severity, payload),
				Pattern:          xxePattern.pattern.String(),
				URL:              entry.URL,
				IP:               entry.IP,
				UserAgent:        entry.UserAgent,
				Timestamp:        entry.Timestamp,
				Method:           entry.Method,
				StatusCode:       entry.Status,
				ResponseSize:     entry.Size,
				AttackVector:     "HTTP Request",
				Payload:          payload,
				Context:          map[string]interface{}{"description": xxePattern.desc},
				MitigationAdvice: []string{"Disable external entity processing", "Use secure XML parsers", "Validate XML input"},
			}
			threats = append(threats, threat)
		}
	}

	return threats
}

// detectHeaderInjection detects HTTP header injection
func (td *ThreatDetector) detectHeaderInjection(entry *parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	headerPatterns := []struct {
		pattern  *regexp.Regexp
		severity ThreatSeverity
		desc     string
	}{
		{regexp.MustCompile(`(%0d%0a|%0a%0d|\\r\\n|\\n\\r)`), SeverityHigh, "CRLF injection"},
		{regexp.MustCompile(`(?i)(set-cookie:|location:|content-type:)`), SeverityMedium, "Header manipulation"},
		{regexp.MustCompile(`(%20Set-Cookie:|%20Location:|%20Content-Type:)`), SeverityMedium, "URL-encoded header injection"},
	}

	target := entry.URL + " " + entry.UserAgent + " " + entry.Referer

	for _, headerPattern := range headerPatterns {
		if headerPattern.pattern.MatchString(target) {
			payload := headerPattern.pattern.FindString(target)
			threat := EnhancedThreat{
				ID:               fmt.Sprintf("header_%d_%s", time.Now().UnixNano(), entry.IP),
				Type:             HTTPHeaderInjection,
				Severity:         headerPattern.severity,
				Confidence:       td.calculateConfidence(headerPattern.severity, payload),
				Pattern:          headerPattern.pattern.String(),
				URL:              entry.URL,
				IP:               entry.IP,
				UserAgent:        entry.UserAgent,
				Timestamp:        entry.Timestamp,
				Method:           entry.Method,
				StatusCode:       entry.Status,
				ResponseSize:     entry.Size,
				AttackVector:     "HTTP Request",
				Payload:          payload,
				Context:          map[string]interface{}{"description": headerPattern.desc},
				MitigationAdvice: []string{"Validate header values", "Sanitize user input", "Use secure response handling"},
			}
			threats = append(threats, threat)
		}
	}

	return threats
}

// detectBruteForce detects brute force login attempts
func (td *ThreatDetector) detectBruteForce(ip string, entries []*parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	// Count failed authentication attempts
	failedAttempts := 0
	authPaths := []string{"/login", "/admin", "/wp-admin", "/auth", "/signin"}
	
	for _, entry := range entries {
		if entry.Status == 401 || entry.Status == 403 {
			for _, path := range authPaths {
				if strings.Contains(strings.ToLower(entry.URL), path) {
					failedAttempts++
					break
				}
			}
		}
	}

	// Threshold-based detection
	threshold := 10 // More than 10 failed attempts
	if failedAttempts > threshold {
		severity := SeverityMedium
		if failedAttempts > 50 {
			severity = SeverityHigh
		}
		if failedAttempts > 100 {
			severity = SeverityCritical
		}

		threat := EnhancedThreat{
			ID:               fmt.Sprintf("brute_%d_%s", time.Now().UnixNano(), ip),
			Type:             BruteForceLogin,
			Severity:         severity,
			Confidence:       float64(failedAttempts) / 100.0,
			Pattern:          "Multiple failed authentication attempts",
			URL:              "/auth-endpoints",
			IP:               ip,
			Timestamp:        entries[len(entries)-1].Timestamp,
			Method:           "POST",
			AttackVector:     "Authentication",
			Context:          map[string]interface{}{"failed_attempts": failedAttempts},
			MitigationAdvice: []string{"Implement account lockout", "Use CAPTCHA", "Enable rate limiting"},
		}
		threats = append(threats, threat)
	}

	return threats
}

// detectDDoS detects Distributed Denial of Service patterns
func (td *ThreatDetector) detectDDoS(ip string, entries []*parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	// Analyze request frequency
	if len(entries) < 50 {
		return threats // Not enough requests for DDoS detection
	}

	// Calculate requests per minute
	duration := entries[len(entries)-1].Timestamp.Sub(entries[0].Timestamp)
	if duration == 0 {
		return threats
	}

	requestsPerMinute := float64(len(entries)) / duration.Minutes()
	
	// DDoS thresholds
	if requestsPerMinute > 100 { // More than 100 requests per minute
		severity := SeverityMedium
		if requestsPerMinute > 500 {
			severity = SeverityHigh
		}
		if requestsPerMinute > 1000 {
			severity = SeverityCritical
		}

		threat := EnhancedThreat{
			ID:           fmt.Sprintf("ddos_%d_%s", time.Now().UnixNano(), ip),
			Type:         DDoSAttack,
			Severity:     severity,
			Confidence:   0.8,
			Pattern:      "High-frequency request pattern",
			IP:           ip,
			Timestamp:    entries[len(entries)-1].Timestamp,
			AttackVector: "Network flooding",
			Context:      map[string]interface{}{"requests_per_minute": requestsPerMinute, "total_requests": len(entries)},
			MitigationAdvice: []string{"Implement rate limiting", "Use DDoS protection service", "Block suspicious IPs"},
		}
		threats = append(threats, threat)
	}

	return threats
}

// detectPortScanning detects port scanning behavior
func (td *ThreatDetector) detectPortScanning(ip string, entries []*parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	// Track unique URLs accessed
	uniqueURLs := make(map[string]bool)
	for _, entry := range entries {
		parsedURL, err := url.Parse(entry.URL)
		if err == nil {
			uniqueURLs[parsedURL.Path] = true
		}
	}

	// Port scanning indicators
	scanningIndicators := []string{
		"/server-status", "/server-info", "/.well-known", "/robots.txt",
		"/sitemap.xml", "/admin", "/phpmyadmin", "/wp-admin",
	}

	scanCount := 0
	for path := range uniqueURLs {
		for _, indicator := range scanningIndicators {
			if strings.Contains(strings.ToLower(path), indicator) {
				scanCount++
				break
			}
		}
	}

	if scanCount >= 5 || len(uniqueURLs) > 50 {
		threat := EnhancedThreat{
			ID:           fmt.Sprintf("portscan_%d_%s", time.Now().UnixNano(), ip),
			Type:         PortScan,
			Severity:     SeverityMedium,
			Confidence:   0.7,
			Pattern:      "Multiple endpoint enumeration",
			IP:           ip,
			Timestamp:    entries[len(entries)-1].Timestamp,
			AttackVector: "Network reconnaissance",
			Context:      map[string]interface{}{"unique_paths": len(uniqueURLs), "scan_indicators": scanCount},
			MitigationAdvice: []string{"Hide server information", "Implement access controls", "Monitor for reconnaissance"},
		}
		threats = append(threats, threat)
	}

	return threats
}

// detectVulnerabilityScanning detects vulnerability scanning tools
func (td *ThreatDetector) detectVulnerabilityScanning(ip string, entries []*parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	// Known vulnerability scanner patterns
	scannerPatterns := []struct {
		pattern  *regexp.Regexp
		severity ThreatSeverity
		desc     string
	}{
		{regexp.MustCompile(`(?i)(nmap|nikto|nessus|openvas|sqlmap|dirb|dirbuster)`), SeverityHigh, "Known vulnerability scanner"},
		{regexp.MustCompile(`(?i)(masscan|zap|burp|acunetix|nuclei)`), SeverityHigh, "Security testing tool"},
		{regexp.MustCompile(`(?i)(gobuster|ffuf|wfuzz|dirfuzz)`), SeverityMedium, "Directory brute-force tool"},
		{regexp.MustCompile(`(?i)(python-requests|curl|wget)\/[\d.]+$`), SeverityLow, "Automated request tool"},
	}

	for _, entry := range entries {
		for _, scannerPattern := range scannerPatterns {
			if scannerPattern.pattern.MatchString(entry.UserAgent) {
				threat := EnhancedThreat{
					ID:               fmt.Sprintf("vulnscan_%d_%s", time.Now().UnixNano(), ip),
					Type:             VulnerabilityScanning,
					Severity:         scannerPattern.severity,
					Confidence:       0.9,
					Pattern:          scannerPattern.pattern.String(),
					URL:              entry.URL,
					IP:               ip,
					UserAgent:        entry.UserAgent,
					Timestamp:        entry.Timestamp,
					Method:           entry.Method,
					StatusCode:       entry.Status,
					AttackVector:     "Automated scanning",
					Payload:          entry.UserAgent,
					Context:          map[string]interface{}{"description": scannerPattern.desc},
					MitigationAdvice: []string{"Block scanner IPs", "Implement rate limiting", "Monitor for scanning patterns"},
				}
				threats = append(threats, threat)
				break // Only report once per IP
			}
		}
	}

	return threats
}

// detectBotActivity detects malicious bot activity
func (td *ThreatDetector) detectBotActivity(ip string, entries []*parser.LogEntry) []EnhancedThreat {
	var threats []EnhancedThreat

	// Bot detection patterns
	botPatterns := []struct {
		pattern  *regexp.Regexp
		severity ThreatSeverity
		desc     string
	}{
		{regexp.MustCompile(`(?i)(bot|crawler|spider|scraper)`), SeverityLow, "Generic bot pattern"},
		{regexp.MustCompile(`(?i)(malicious|badbot|evil|hack)`), SeverityHigh, "Malicious bot pattern"},
		{regexp.MustCompile(`^Mozilla\/5\.0$`), SeverityMedium, "Suspicious generic user agent"},
		{regexp.MustCompile(`(?i)(python|java|php)\/[\d.]+$`), SeverityMedium, "Scripted request pattern"},
	}

	// Check for bot indicators
	for _, entry := range entries {
		for _, botPattern := range botPatterns {
			if botPattern.pattern.MatchString(entry.UserAgent) {
				// Additional validation for bot behavior
				if td.isSuspiciousBotBehavior(entries) {
					threat := EnhancedThreat{
						ID:           fmt.Sprintf("bot_%d_%s", time.Now().UnixNano(), ip),
						Type:         BotnetActivity,
						Severity:     botPattern.severity,
						Confidence:   0.6,
						Pattern:      botPattern.pattern.String(),
						IP:           ip,
						UserAgent:    entry.UserAgent,
						Timestamp:    entry.Timestamp,
						AttackVector: "Automated activity",
						Context:      map[string]interface{}{"description": botPattern.desc},
						MitigationAdvice: []string{"Implement bot detection", "Use CAPTCHA", "Rate limit suspicious IPs"},
					}
					threats = append(threats, threat)
				}
				break
			}
		}
	}

	return threats
}

// Helper functions

// initializePatterns initializes regex patterns for threat detection
func (td *ThreatDetector) initializePatterns() {
	// Patterns are initialized in individual detection functions
	// This maintains flexibility and readability
}

// loadThreatIntelligence loads threat intelligence data
func (td *ThreatDetector) loadThreatIntelligence() {
	// Initialize with basic known bad patterns
	// In production, this would load from external threat feeds
	
	// Example known bad IPs (normally loaded from external sources)
	knownBadIPs := []string{
		"0.0.0.0", "127.0.0.1", // Localhost attacks
	}
	
	for _, badIP := range knownBadIPs {
		td.knownBadIPs[badIP] = ThreatInfo{
			IP:          badIP,
			ThreatTypes: []string{"known_malicious"},
			Severity:    SeverityHigh,
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
			Description: "Known malicious IP address",
		}
	}
}

// calculateConfidence calculates threat confidence based on severity and payload
func (td *ThreatDetector) calculateConfidence(severity ThreatSeverity, payload string) float64 {
	baseConfidence := 0.5
	
	// Adjust based on severity
	switch severity {
	case SeverityCritical:
		baseConfidence = 0.9
	case SeverityHigh:
		baseConfidence = 0.8
	case SeverityMedium:
		baseConfidence = 0.6
	case SeverityLow:
		baseConfidence = 0.4
	default:
		baseConfidence = 0.3
	}
	
	// Adjust based on payload complexity
	if len(payload) > 50 {
		baseConfidence += 0.1
	}
	if strings.Contains(payload, "script") || strings.Contains(payload, "union") {
		baseConfidence += 0.1
	}
	
	// Cap at 1.0
	if baseConfidence > 1.0 {
		baseConfidence = 1.0
	}
	
	return baseConfidence
}

// isSuspiciousBotBehavior analyzes bot behavior patterns
func (td *ThreatDetector) isSuspiciousBotBehavior(entries []*parser.LogEntry) bool {
	if len(entries) < 5 {
		return false
	}
	
	// Check for rapid sequential requests
	var intervals []time.Duration
	for i := 1; i < len(entries); i++ {
		interval := entries[i].Timestamp.Sub(entries[i-1].Timestamp)
		intervals = append(intervals, interval)
	}
	
	// Calculate average interval
	var totalInterval time.Duration
	for _, interval := range intervals {
		totalInterval += interval
	}
	avgInterval := totalInterval / time.Duration(len(intervals))
	
	// Suspicious if requests are too regular (bot-like) or too frequent
	return avgInterval < 5*time.Second || (avgInterval < 60*time.Second && len(entries) > 20)
}

// isPrivateIP checks if an IP is in a private range
func (td *ThreatDetector) isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	
	// Private IP ranges
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
	}
	
	for _, rangeStr := range privateRanges {
		_, network, err := net.ParseCIDR(rangeStr)
		if err == nil && network.Contains(ip) {
			return true
		}
	}
	
	return false
}