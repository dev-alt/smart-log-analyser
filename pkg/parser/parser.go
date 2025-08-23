package parser

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogEntry struct {
	IP        string
	Timestamp time.Time
	Method    string
	URL       string
	Protocol  string
	Status    int
	Size      int64
	Referer   string
	UserAgent string
}

type Parser struct {
	combinedRegex *regexp.Regexp
	commonRegex   *regexp.Regexp
}

func New() *Parser {
	combinedPattern := `^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d+) (\d+) "([^"]*)" "([^"]*)"$`
	commonPattern := `^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d+) (\d+)$`

	return &Parser{
		combinedRegex: regexp.MustCompile(combinedPattern),
		commonRegex:   regexp.MustCompile(commonPattern),
	}
}

func (p *Parser) ParseFile(filename string) ([]*LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a reader that handles compressed files
	reader, err := p.createReader(file, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader for %s: %w", filename, err)
	}
	defer func() {
		if closer, ok := reader.(io.Closer); ok {
			closer.Close()
		}
	}()

	var entries []*LogEntry
	scanner := bufio.NewScanner(reader)
	
	// Increase buffer size for potentially large compressed files
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 1MB buffer
	
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" {
			continue
		}

		entry, err := p.ParseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse line %d in %s: %v\n", lineNum, filepath.Base(filename), err)
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	return entries, nil
}

// createReader creates appropriate reader based on file extension
func (p *Parser) createReader(file *os.File, filename string) (io.Reader, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".gz":
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		return gzReader, nil
	case ".log":
		// Regular log file
		return file, nil
	default:
		// Check if it's a numbered log file (like .log.1, .log.2, etc.)
		if p.isRotatedLogFile(filename) {
			return file, nil
		}
		// For files without extension or unknown extensions, treat as regular text
		return file, nil
	}
}

// isRotatedLogFile checks if filename matches rotated log pattern
func (p *Parser) isRotatedLogFile(filename string) bool {
	// Match patterns like: access.log.1, error.log.12, site.access.log.5, etc.
	rotatedPattern := regexp.MustCompile(`\.(log|access|error)(\.\d+)?$`)
	return rotatedPattern.MatchString(strings.ToLower(filename))
}

func (p *Parser) ParseLine(line string) (*LogEntry, error) {
	if matches := p.combinedRegex.FindStringSubmatch(line); matches != nil {
		return p.parseCombinedFormat(matches)
	}
	
	if matches := p.commonRegex.FindStringSubmatch(line); matches != nil {
		return p.parseCommonFormat(matches)
	}

	return nil, fmt.Errorf("line does not match any known log format")
}

func (p *Parser) parseCombinedFormat(matches []string) (*LogEntry, error) {
	ip := matches[1]
	if !isValidIP(ip) {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	timestamp, err := parseTimestamp(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	status, err := strconv.Atoi(matches[6])
	if err != nil {
		return nil, fmt.Errorf("invalid status code: %w", err)
	}

	size, err := strconv.ParseInt(matches[7], 10, 64)
	if err != nil {
		size = 0
	}

	return &LogEntry{
		IP:        ip,
		Timestamp: timestamp,
		Method:    matches[3],
		URL:       matches[4],
		Protocol:  matches[5],
		Status:    status,
		Size:      size,
		Referer:   matches[8],
		UserAgent: matches[9],
	}, nil
}

func (p *Parser) parseCommonFormat(matches []string) (*LogEntry, error) {
	ip := matches[1]
	if !isValidIP(ip) {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	timestamp, err := parseTimestamp(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	status, err := strconv.Atoi(matches[6])
	if err != nil {
		return nil, fmt.Errorf("invalid status code: %w", err)
	}

	size, err := strconv.ParseInt(matches[7], 10, 64)
	if err != nil {
		size = 0
	}

	return &LogEntry{
		IP:        ip,
		Timestamp: timestamp,
		Method:    matches[3],
		URL:       matches[4],
		Protocol:  matches[5],
		Status:    status,
		Size:      size,
		Referer:   "",
		UserAgent: "",
	}, nil
}

func parseTimestamp(timestampStr string) (time.Time, error) {
	return time.Parse("02/Jan/2006:15:04:05 -0700", timestampStr)
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}