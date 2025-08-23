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
- [x] SSH remote log file download

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

### Local Analysis
```bash
# Analyse a local log file
./smart-log-analyser analyse /var/log/nginx/access.log

# Filter by time range
./smart-log-analyser analyse /var/log/nginx/access.log --since="2024-08-20 00:00:00" --until="2024-08-20 23:59:59"

# Get top 10 IPs and URLs
./smart-log-analyser analyse /var/log/nginx/access.log --top-ips=10 --top-urls=10

# Test with sample data
./smart-log-analyser analyse testdata/sample_access.log
```

### Remote Server Access
```bash
# Create SSH configuration file (only if doesn't exist)
./smart-log-analyser download --init

# Test SSH connections
./smart-log-analyser download --test

# List available log files without downloading
./smart-log-analyser download --list

# Download single log file (default: access.log only)
./smart-log-analyser download

# Download ALL access log files (current + rotated)
./smart-log-analyser download --all

# Download limited number of files
./smart-log-analyser download --all --max-files 5

# Download from specific server
./smart-log-analyser download --server your-server.com --all

# Analyse downloaded files
./smart-log-analyser analyse ./downloads/*.log
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
│   ├── analyse.go        # Analysis command
│   └── download.go       # Remote download command
├── pkg/
│   ├── parser/           # Log parsing logic
│   ├── analyser/         # Analysis algorithms
│   └── remote/           # SSH client and configuration
├── testdata/             # Sample log files for testing
├── servers.json.example  # Example SSH configuration
├── main.go               # Application entry point
└── README.md             # This file
```

## Command Line Options

### `analyse` command

- `--since`: Start time for analysis (format: "YYYY-MM-DD HH:MM:SS")
- `--until`: End time for analysis (format: "YYYY-MM-DD HH:MM:SS")
- `--top-ips`: Number of top IP addresses to display (default: 10)
- `--top-urls`: Number of top URLs to display (default: 10)

### `download` command

- `--config`: Path to SSH configuration file (default: "servers.json")
- `--server`: Specific server to download from (host name)
- `--output`: Directory to save downloaded files (default: "./downloads")
- `--test`: Test SSH connection without downloading
- `--init`: Create a sample configuration file (will not overwrite existing files)
- `--list`: List available log files without downloading
- `--all`: Download all access log files (current + rotated)
- `--max-files`: Maximum number of files to download when using --all (default: 10)

## SSH Configuration

The SSH configuration is stored in JSON format. Create it with:

```bash
# Creates servers.json only if it doesn't exist
./smart-log-analyser download --init

# Create with different filename
./smart-log-analyser download --init --config my-servers.json
```

**Note**: The `--init` command will **not overwrite** existing configuration files. If a file already exists, it will display the current configuration instead.

Example `servers.json`:
```json
{
  "servers": [
    {
      "host": "your-server.com",
      "port": 22,
      "username": "root",
      "password": "your-password",
      "log_path": "/var/log/nginx/access.log"
    }
  ]
}
```

⚠️ **Security Note**: Store the configuration file securely and restrict permissions (`chmod 600 servers.json`).

## Multi-File Download Features

### 🔍 Discovery and Listing
```bash
# List all available log files on server
./smart-log-analyser download --list
```

This will show you all access log files including:
- Current log files (`access.log`, `forum.access.log`, etc.)
- Rotated log files (`access.log.1`, `access.log.2`, etc.) 
- Compressed logs (`access.log.gz`, `access.log.12.gz`, etc.)

### 📦 Bulk Download Options
```bash
# Download all access log files (up to 10 by default)
./smart-log-analyser download --all

# Download more files (up to 20)
./smart-log-analyser download --all --max-files 20

# Download from specific server only
./smart-log-analyser download --all --server your-server.com
```

### 📊 Download Behavior
- **Single file mode** (default): Downloads only the configured `log_path` file
- **Multi-file mode** (`--all`): Downloads all access log files found in the log directory
- **Smart naming**: Files are saved as `hostname_timestamp_originalname` to avoid conflicts
- **Progress tracking**: Shows download progress for each file with size information

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
./smart-log-analyser download --help

# Test SSH configuration creation
./smart-log-analyser download --init
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
- **v0.1.1**: SSH remote log download ✅
- **v0.2.0**: Advanced analytics and export features
- **v0.3.0**: Real-time monitoring and alerting
- **v1.0.0**: Production-ready with comprehensive documentation

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) - SSH connectivity
- Go standard library for log parsing and analysis

## Security Notes

### 🔐 Credential Security
- SSH configuration files contain sensitive credentials and are **automatically excluded** from version control
- Use secure file permissions: `chmod 600 servers.json`
- Never commit real passwords, server IPs, or SSH keys to git
- Use the provided `servers.json.example` as a template

### 🛡️ Production Security Recommendations
- **Use SSH key authentication** instead of passwords in production
- **Restrict network access** to log servers (VPN, firewall rules)
- **Rotate credentials regularly** and use strong passwords
- **Monitor access logs** for unauthorized usage
- **Consider log aggregation systems** instead of direct server access

### ⚠️ Development Security
- Real server credentials in `servers.json` are excluded from git commits
- Test connections are logged - avoid using production servers for testing
- Downloaded log files may contain sensitive data - they are also excluded from git
- Review `.gitignore` regularly to ensure all sensitive patterns are covered

## Development Guidelines

### 📋 Development Workflow (Mandatory)
All contributors must follow these steps for every development session:

1. **Documentation First**: Update README.md and relevant docs for new features
2. **Security Review**: Check all changes for sensitive data before committing  
3. **Development Log**: Update `.development_log.md` with session details
4. **Testing**: Verify new features work and existing functionality remains intact
5. **Git Workflow**: Stage, commit with descriptive message, and push to GitHub

See `DEVELOPMENT_RULES.md` for comprehensive development standards and security practices.

### 🔍 Security Checklist
Before every commit, verify:
- [ ] No real passwords, API keys, or tokens in any files
- [ ] No SSH private keys or certificates committed
- [ ] No real server IPs or sensitive hostnames
- [ ] `.gitignore` updated for new sensitive file patterns  
- [ ] Documentation updated with security warnings
- [ ] Example files use placeholder values only