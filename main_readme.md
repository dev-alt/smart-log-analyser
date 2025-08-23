# Smart Log Analyzer

A high-performance CLI tool for analyzing Nginx access logs with real-time monitoring capabilities.

## Overview

Smart Log Analyzer is designed to help system administrators and developers gain insights from their Nginx access logs. It provides statistical analysis, error pattern detection, traffic analysis, and real-time monitoring with configurable alerting.

## Features

### Phase 1 (MVP)
- [x] Parse standard Nginx access log formats (common/combined)
- [x] Basic statistics: request counts, status code distribution, top IPs, top URLs
- [x] Time range filtering
- [x] Clean console output with formatting

### Phase 2 (Analytics)
- [ ] Error pattern detection and analysis
- [ ] Traffic analysis (requests per hour, peak detection)
- [ ] Response time analysis and percentiles
- [ ] Geographic IP analysis
- [ ] Export to JSON/CSV formats

### Phase 3 (Real-time & Alerts)
- [ ] Real-time log file monitoring
- [ ] Configurable alert rules (YAML configuration)
- [ ] Custom alert thresholds
- [ ] Multiple output destinations

## Installation

```bash
go install github.com/yourusername/smart-log-analyzer@latest
```

Or download the latest binary from the [releases page](https://github.com/yourusername/smart-log-analyzer/releases).

## Quick Start

```bash
# Analyze a log file
smart-log-analyzer analyze /var/log/nginx/access.log

# Filter by time range
smart-log-analyzer analyze /var/log/nginx/access.log --since="2024-08-20 00:00:00" --until="2024-08-20 23:59:59"

# Get top 10 IPs
smart-log-analyzer analyze /var/log/nginx/access.log --top-ips=10

# Real-time monitoring (Phase 3)
smart-log-analyzer monitor /var/log/nginx/access.log --config=alerts.yaml
```

## Project Structure

```
smart-log-analyzer/
├── cmd/
│   └── root.go           # CLI commands and flags
├── pkg/
│   ├── parser/           # Log parsing logic
│   ├── analyzer/         # Analysis algorithms
│   ├── monitor/          # Real-time monitoring
│   └── output/           # Output formatters
├── configs/
│   └── alerts.yaml       # Sample alert configuration
├── testdata/             # Sample log files for testing
├── docs/                 # Additional documentation
└── scripts/              # Build and deployment scripts
```

## Supported Log Formats

Currently supports standard Nginx access log formats:

- **Combined Log Format**: `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`
- **Common Log Format**: `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`

## Development

See [DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed development instructions.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- **v0.1.0**: Basic log parsing and statistics
- **v0.2.0**: Advanced analytics and export features
- **v0.3.0**: Real-time monitoring and alerting
- **v1.0.0**: Production-ready with comprehensive documentation