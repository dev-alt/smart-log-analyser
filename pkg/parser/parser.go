package parser

import (
	"bufio"
	"fmt"
	"net"
	"os"
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

	var entries []*LogEntry
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" {
			continue
		}

		entry, err := p.ParseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse line %d: %v\n", lineNum, err)
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return entries, nil
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