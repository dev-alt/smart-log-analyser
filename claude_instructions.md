# Claude CLI Project Instructions

## Project Overview
You are building a **Smart Log Analyzer** - a CLI tool in Go for analyzing Nginx access logs. This is a phased project starting simple and adding features incrementally.

## Phase 1 Goals (MVP)
Build a basic CLI tool that can:
1. Parse standard Nginx access log formats (combined format priority)
2. Generate basic statistics (request counts, status codes, top IPs, top URLs)
3. Filter by time ranges
4. Display results in a clean console format

## Technical Requirements

### Language & Dependencies
- **Language**: Go (1.21+)
- **CLI Framework**: Use `github.com/spf13/cobra` for commands
- **Testing**: Use `github.com/stretchr/testify` for tests
- **No heavy frameworks** - keep it lightweight

### Project Structure
```
smart-log-analyzer/
├── cmd/
│   ├── main.go          # Entry point
│   ├── root.go          # Root command setup
│   └── analyze.go       # Analyze command
├── pkg/
│   ├── parser/
│   │   ├── nginx_parser.go    # Main parsing logic
│   │   ├── log_entry.go      # LogEntry struct
│   │   └── formats.go        # Format definitions
│   ├── analyzer/
│   │   └── basic_stats.go    # Statistics calculations
│   └── output/
│       └── console.go        # Console formatting
├── testdata/
│   └── sample_nginx.log      # Test data
├── go.mod
├── go.sum
└── README.md
```

### Core Data Structure
```go
type LogEntry struct {
    IP        string    `json:"ip"`
    Timestamp time.Time `json:"timestamp"`
    Method    string    `json:"method"`
    URL       string    `json:"url"`
    Status    int       `json:"status"`
    Size      int64     `json:"size"`
    Referer   string    `json:"referer,omitempty"`
    UserAgent string    `json:"user_agent,omitempty"`
}
```

### Nginx Combined Log Format to Parse
```
$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"
```

Example log line:
```
192.168.1.1 - - [22/Aug/2024:10:30:45 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
```

## Implementation Steps

### Step 1: Project Initialization
1. Create the directory structure
2. Initialize Go module: `go mod init smart-log-analyzer`
3. Add cobra dependency: `go get github.com/spf13/cobra`
4. Create basic main.go and root command

### Step 2: Log Parser
1. Create LogEntry struct in `pkg/parser/log_entry.go`
2. Implement regex-based parser in `pkg/parser/nginx_parser.go`
3. Handle timestamp parsing (format: "02/Jan/2006:15:04:05 -0700")
4. Add error handling for malformed lines

### Step 3: Basic Statistics
1. Implement in `pkg/analyzer/basic_stats.go`:
   - Total request count
   - Status code distribution (2xx, 3xx, 4xx, 5xx)
   - Top 10 IP addresses by request count
   - Top 10 URLs by request count
   - Date range of logs

### Step 4: Console Output
1. Create `pkg/output/console.go`
2. Format statistics in a readable table format
3. Use color coding for different status code categories
4. Add progress indicators for large file processing

### Step 5: CLI Interface
1. Implement `analyze` command with flags:
   - `--since` and `--until` for time filtering
   - `--top-ips` (default 10) for number of top IPs to show
   - `--top-urls` (default 10) for number of top URLs to show
   - `--format` (future: json, csv support)

### Step 6: Testing
1. Create sample log files in `testdata/`
2. Write unit tests for parser
3. Write unit tests for statistics
4. Add integration tests

## Command Usage Examples
```bash
# Basic analysis
smart-log-analyzer analyze /var/log/nginx/access.log

# With time filtering
smart-log-analyzer analyze /var/log/nginx/access.log --since="2024-08-20 00:00:00" --until="2024-08-20 23:59:59"

# Custom top counts
smart-log-analyzer analyze /var/log/nginx/access.log --top-ips=5 --top-urls=20
```

## Code Quality Requirements
- Handle errors gracefully with descriptive messages
- Use Go best practices (proper error handling, meaningful names)
- Add progress indicators for large files
- Memory efficient processing (don't load entire file into memory)
- Add basic validation for log format detection

## Success Criteria for Phase 1
- [ ] Successfully parses Nginx combined format logs
- [ ] Generates accurate basic statistics
- [ ] Handles time range filtering correctly
- [ ] Provides clean, readable console output
- [ ] Includes comprehensive tests
- [ ] Handles errors gracefully
- [ ] Processes large files efficiently

## Future Phases (Don't implement yet)
- **Phase 2**: Advanced analytics, error pattern detection, export formats
- **Phase 3**: Real-time monitoring with fsnotify, configurable alerts

Start with a minimal working version and build up incrementally. Focus on code quality, error handling, and user experience.