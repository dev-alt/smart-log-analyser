# Smart Log Analyser

A high-performance CLI tool for analysing Nginx access logs with real-time monitoring capabilities.

## Overview

Smart Log Analyser is designed to help system administrators and developers gain insights from their Nginx access logs. It provides statistical analysis, error pattern detection, traffic analysis, and real-time monitoring with configurable alerting.

## Features

### Phase 1 (MVP) ✅
- [x] Parse standard Nginx access log formats (common/combined)
- [x] Basic statistics: request counts, status code distribution, top IPs, top URLs
- [x] Time range filtering
- [x] Clean console output with formatting

### Phase 2 (Analytics) 🚧
- [ ] Error pattern detection and analysis
- [ ] Traffic analysis (requests per hour, peak detection)
- [ ] Response time analysis and percentiles
- [ ] Geographic IP analysis
- [ ] Export to JSON/CSV formats

### Phase 3 (Real-time & Alerts) 📋
- [ ] Real-time log file monitoring
- [ ] Configurable alert rules (YAML configuration)
- [ ] Custom alert thresholds
- [ ] Multiple output destinations

## Installation

### From Source
```bash
git clone https://github.com/dev-alt/smart-log-analyser.git
cd smart-log-analyser
go build -o smart-log-analyser
```

### Using Go Install
```bash
go install github.com/dev-alt/smart-log-analyser@latest
```

## Quick Start

```bash
# Analyse a log file
./smart-log-analyser analyse /var/log/nginx/access.log

# Filter by time range
./smart-log-analyser analyse /var/log/nginx/access.log --since="2024-08-20 00:00:00" --until="2024-08-20 23:59:59"

# Get top 10 IPs and URLs
./smart-log-analyser analyse /var/log/nginx/access.log --top-ips=10 --top-urls=10

# Test with sample data
./smart-log-analyser analyse testdata/sample_access.log
```

## Example Output

```
=== Smart Log Analyser Results ===

Total Requests: 10
Date Range: 2024-08-22 10:15:30 to 2024-08-22 10:24:30

=== Status Code Distribution ===
2xx Success: 7
4xx Client Error: 2
5xx Server Error: 1

=== Top 10 IP Addresses ===
192.168.1.100: 3 requests
10.0.0.5: 3 requests
203.0.113.1: 2 requests
198.51.100.42: 2 requests

=== Top 10 URLs ===
/index.html: 1 requests
/api/login: 1 requests
/about.html: 1 requests
/products.html: 1 requests
```

## Supported Log Formats

Currently supports standard Nginx access log formats:

- **Combined Log Format**: `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
- **Common Log Format**: `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`

## Project Structure

```
smart-log-analyser/
├── cmd/
│   ├── root.go           # CLI root command
│   └── analyse.go        # Analysis command
├── pkg/
│   ├── parser/           # Log parsing logic
│   └── analyser/         # Analysis algorithms
├── testdata/             # Sample log files for testing
├── main.go               # Application entry point
└── README.md             # This file
```

## Command Line Options

### `analyse` command

- `--since`: Start time for analysis (format: "YYYY-MM-DD HH:MM:SS")
- `--until`: End time for analysis (format: "YYYY-MM-DD HH:MM:SS")
- `--top-ips`: Number of top IP addresses to display (default: 10)
- `--top-urls`: Number of top URLs to display (default: 10)

## Development

### Prerequisites
- Go 1.18 or higher
- Git

### Building
```bash
go build -o smart-log-analyser
```

### Testing
```bash
# Test with sample data
./smart-log-analyser analyse testdata/sample_access.log

# Test help commands
./smart-log-analyser --help
./smart-log-analyser analyse --help
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Roadmap

- **v0.1.0**: Basic log parsing and statistics ✅
- **v0.2.0**: Advanced analytics and export features
- **v0.3.0**: Real-time monitoring and alerting
- **v1.0.0**: Production-ready with comprehensive documentation

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- Go standard library for log parsing and analysis