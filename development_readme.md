# Development Guide

## Prerequisites

- Go 1.21 or higher
- Git

## Project Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/smart-log-analyzer.git
cd smart-log-analyzer
```

2. Initialize Go module:
```bash
go mod init smart-log-analyzer
```

3. Install dependencies:
```bash
go mod tidy
```

## Project Architecture

### Core Components

#### Parser Package (`pkg/parser/`)
Handles parsing of different Nginx log formats into structured data.

**Key files:**
- `nginx_parser.go` - Main parsing logic
- `log_entry.go` - Data structures for log entries
- `formats.go` - Supported log format definitions

#### Analyzer Package (`pkg/analyzer/`)
Contains analysis algorithms and statistical calculations.

**Key files:**
- `basic_stats.go` - Basic statistics (counts, distributions)
- `traffic_analyzer.go` - Traffic pattern analysis
- `error_analyzer.go` - Error detection and patterns

#### Output Package (`pkg/output/`)
Handles formatting and displaying results.

**Key files:**
- `console.go` - Console output formatting
- `json.go` - JSON export functionality
- `csv.go` - CSV export functionality

### Phase 1 Implementation Plan

#### Step 1: Basic Log Entry Structure
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

#### Step 2: Parser Implementation
- Support for combined log format first
- Regex-based parsing with proper error handling
- Validation of parsed entries

#### Step 3: Basic Statistics
- Request count by status code
- Top IPs by request count
- Top URLs by request count
- Time range filtering

#### Step 4: CLI Interface
- Use cobra for command structure
- Implement `analyze` command
- Add flags for filtering and output options

## Testing Strategy

### Unit Tests
- Test each parser format independently
- Test statistical calculations with known data
- Test CLI argument parsing

### Integration Tests
- Test with real Nginx log samples
- Test full analysis pipeline

### Test Data
Create sample log files in `testdata/` directory:
- `testdata/nginx_combined.log` - Combined format sample
- `testdata/nginx_common.log` - Common format sample
- `testdata/nginx_large.log` - Large file for performance testing

## Code Style

- Follow Go standard formatting (use `gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Handle errors explicitly
- Use interfaces for testability

## Performance Considerations

- Stream processing for large files
- Memory-efficient parsing
- Concurrent processing where appropriate
- Progress indicators for large file processing

## Dependencies

### Phase 1
```go
require (
    github.com/spf13/cobra v1.7.0
    github.com/stretchr/testify v1.8.4
)
```

### Phase 2 (Future)
```go
require (
    github.com/fsnotify/fsnotify v1.6.0  // Real-time monitoring
    gopkg.in/yaml.v3 v3.0.1             // Configuration files
)
```

## Build and Release

### Local Development
```bash
# Build for current platform
go build -o bin/smart-log-analyzer ./cmd

# Run tests
go test ./...

# Run with sample data
./bin/smart-log-analyzer analyze testdata/nginx_combined.log
```

### Cross-platform Builds
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/smart-log-analyzer-linux-amd64 ./cmd

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/smart-log-analyzer-windows-amd64.exe ./cmd

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/smart-log-analyzer-darwin-amd64 ./cmd
```

## Git Workflow

1. Create feature branches from `main`
2. Keep commits focused and atomic
3. Write descriptive commit messages
4. Test before committing
5. Create PR with description of changes

### Commit Message Format
```
feat: add basic nginx log parsing
fix: handle malformed timestamp entries
docs: update installation instructions
test: add parser unit tests
```

## Debugging

### Common Issues
- **Time parsing errors**: Check timezone handling in log timestamps
- **Memory usage**: Monitor for memory leaks with large files
- **Regex performance**: Profile regex patterns for efficiency

### Useful Debug Commands
```bash
# Run with verbose output
go run ./cmd analyze -v testdata/nginx_combined.log

# Profile memory usage
go run ./cmd analyze --cpuprofile=cpu.prof testdata/large.log

# Race condition detection
go run -race ./cmd analyze testdata/nginx_combined.log
```