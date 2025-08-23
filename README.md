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
- [x] Multi-file analysis support

### Phase 2 (Analytics) ✅
- [x] Enhanced statistics with percentages and visual formatting
- [x] HTTP method analysis (GET, POST, etc.)
- [x] Data transfer analytics (total bytes, average response size)
- [x] Unique visitor/resource counting
- [x] Improved console output with emojis and structured display
- [x] **Bot detection and traffic analysis** (human vs automated traffic)
- [x] **File type analysis** (CSS, JavaScript, images, dynamic content)
- [x] **Top bot/crawler identification** (Googlebot, curl, monitoring tools)
- [x] **Export functionality** (JSON and CSV formats with detailed breakdowns)
- [x] **Detailed drill-down analysis** (individual status codes, error URLs, large requests)
- [x] **Error pattern detection** (4xx/5xx URLs, failure analysis)
- [x] **Traffic pattern analysis** (hourly breakdowns, peak detection, visual charts)
- [x] **Peak traffic detection** (automatic identification of traffic spikes)
- [ ] Response time analysis and percentiles
- [ ] Geographic IP analysis

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
# Analyse a single log file
./smart-log-analyser analyse /var/log/nginx/access.log

# Analyse multiple log files together
./smart-log-analyser analyse /var/log/nginx/access.log /var/log/nginx/access.log.1

# Analyse all downloaded files using wildcard
./smart-log-analyser analyse ./downloads/*.log

# Filter by time range
./smart-log-analyser analyse /var/log/nginx/access.log --since="2024-08-20 00:00:00" --until="2024-08-20 23:59:59"

# Get top 10 IPs and URLs
./smart-log-analyser analyse /var/log/nginx/access.log --top-ips=10 --top-urls=10

# Test with sample data
./smart-log-analyser analyse testdata/sample_access.log

# Show detailed breakdown with individual status codes and error analysis
./smart-log-analyser analyse /var/log/nginx/access.log --details

# Export results for further analysis
./smart-log-analyser analyse ./downloads/*.log --export-json=detailed_report.json --export-csv=summary.csv

# Analyze traffic patterns and identify peak hours
./smart-log-analyser analyse /var/log/nginx/access.log* --details
```

### Remote Server Access
```bash
# Create SSH configuration file (only if doesn't exist)
./smart-log-analyser download --init

# Test SSH connections
./smart-log-analyser download --test

# List available log files without downloading
./smart-log-analyser download --list

# Download ALL access log files (default behavior)
./smart-log-analyser download

# Download single log file only
./smart-log-analyser download --single

# Download limited number of files  
./smart-log-analyser download --max-files 5

# Download from specific server
./smart-log-analyser download --server your-server.com

# Analyse downloaded files
./smart-log-analyser analyse ./downloads/*.log
```

## Example Output

```
📂 Analysing 3 log file(s)...

  [1/3] Processing: ./downloads/server1_20240823_access.log
    ✅ Parsed 1247 entries
  [2/3] Processing: ./downloads/server1_20240823_access.log.1
    ✅ Parsed 2156 entries
  [3/3] Processing: ./downloads/server1_20240823_access.log.2
    ✅ Parsed 893 entries

📊 Combined Analysis Results (4296 total entries):
╔════════════════════════════════════════════════════════════════╗
║                   Smart Log Analyser Results                  ║
╚════════════════════════════════════════════════════════════════╝

📊 Overview
├─ Total Requests: 4,296
├─ Unique IPs: 127
├─ Unique URLs: 48
├─ Data Transferred: 2.1 GB
├─ Average Response Size: 512.3 KB
└─ Date Range: 2024-08-22 10:15:30 to 2024-08-23 23:59:45

🤖 Traffic Analysis
├─ Human Traffic: 3,264 (76.0%)
├─ Bot/Automated: 1,032 (24.0%)

🔍 Top Bots/Crawlers
├─ Googlebot: 287 requests (6.7%)
├─ Bingbot: 156 requests (3.6%)
├─ Facebook Bot: 89 requests (2.1%)
├─ cURL: 67 requests (1.6%)
├─ Monitoring Bot: 43 requests (1.0%)

📁 File Type Analysis
├─ Dynamic/HTML: 2,847 requests (66.3%) - 1.8 GB total, 659.2 KB avg
├─ CSS: 542 requests (12.6%) - 85.4 MB total, 161.2 KB avg
├─ JavaScript: 398 requests (9.3%) - 124.7 MB total, 320.8 KB avg
├─ Images: 287 requests (6.7%) - 67.8 MB total, 241.7 KB avg
├─ Fonts: 89 requests (2.1%) - 15.2 MB total, 174.9 KB avg

📈 Traffic Patterns
├─ Average Requests/Hour: 179.0
├─ Peak Hour: 14:00 (Afternoon)  
├─ Quietest Hour: 03:00 (Night)
└─ Hourly Breakdown:
   ├─ 08:00: 98 requests (2.3%) [█░░░░░░░░░░░░░░░░░░░]
   ├─ 09:00: 156 requests (3.6%) [███░░░░░░░░░░░░░░░░░]
   ├─ 10:00: 234 requests (5.4%) [█████░░░░░░░░░░░░░░░]
   ├─ 11:00: 298 requests (6.9%) [██████░░░░░░░░░░░░░░]
   ├─ 12:00: 432 requests (10.1%) [██████████░░░░░░░░░░]
   ├─ 13:00: 578 requests (13.5%) [█████████████░░░░░░░]
   ├─ 14:00: 892 requests (20.8%) [████████████████████] ← Peak
   ├─ 15:00: 567 requests (13.2%) [█████████████░░░░░░░]
   ├─ 16:00: 341 requests (7.9%) [███████░░░░░░░░░░░░░]

🔥 Traffic Peaks Detected
├─ Peak #1: 2024-08-22 14:00 - 892 requests (1 hour)
├─ Peak #2: 2024-08-22 13:00 - 578 requests (1 hour)

🔧 HTTP Methods
├─ GET: 3,892 (90.6%)
├─ POST: 347 (8.1%)
├─ PUT: 42 (1.0%)
├─ DELETE: 15 (0.3%)

📈 Status Code Distribution
├─ 2xx Success: 3,847 (89.5%)
├─ 4xx Client Error: 312 (7.3%)
├─ 5xx Server Error: 137 (3.2%)

🌐 Top 5 IP Addresses
├─ 192.168.1.100: 247 requests (5.7%)
├─ 10.0.0.5: 198 requests (4.6%)
├─ 203.0.113.1: 156 requests (3.6%)
├─ 198.51.100.42: 143 requests (3.3%)
├─ 172.16.0.15: 98 requests (2.3%)

🔗 Top 5 URLs
├─ /index.html: 89 requests (2.1%)
├─ /api/status: 67 requests (1.6%)
├─ /assets/style.css: 54 requests (1.3%)
├─ /products.html: 43 requests (1.0%)
├─ /about.html: 38 requests (0.9%)
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

**Usage**: `./smart-log-analyser analyse [log-files...]`

Accepts one or more log files for analysis. When multiple files are provided, they are combined for comprehensive analysis.

**Options**:
- `--since`: Start time for analysis (format: "YYYY-MM-DD HH:MM:SS")
- `--until`: End time for analysis (format: "YYYY-MM-DD HH:MM:SS")
- `--top-ips`: Number of top IP addresses to display (default: 10)
- `--top-urls`: Number of top URLs to display (default: 10)
- `--details`: Show detailed breakdown (individual status codes, error URLs, large requests)
- `--export-json`: Export detailed results to JSON file (e.g., `--export-json=report.json`)
- `--export-csv`: Export detailed results to CSV file (e.g., `--export-csv=report.csv`)

### `download` command

- `--config`: Path to SSH configuration file (default: "servers.json")
- `--server`: Specific server to download from (host name)
- `--output`: Directory to save downloaded files (default: "./downloads")
- `--test`: Test SSH connection without downloading
- `--init`: Create a sample configuration file (will not overwrite existing files)
- `--list`: List available log files without downloading
- `--single`: Download only the main configured log file
- `--max-files`: Maximum number of files to download (default: 10)
- `--all`: Download all access log files (same as default behavior)

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

## Export and Analysis Features

### 📊 Export Formats

**JSON Export** (`--export-json`):
- Complete analysis results in structured JSON format
- Includes all metrics, detailed breakdowns, and raw data
- Perfect for programmatic processing or integration with other tools
- Contains individual status codes, bot details, file type statistics

**CSV Export** (`--export-csv`):
- Tabular format suitable for spreadsheet analysis
- Organized by sections (Overview, Status Codes, Top IPs, etc.)
- Includes percentages and detailed metrics
- Easy to import into Excel, Google Sheets, or database systems

### 🔍 Detailed Analysis Mode (`--details`)

When using the `--details` flag, you get additional insights:
- **Individual Status Codes**: See exact HTTP status codes (200, 404, 500, etc.)
- **Error Analysis**: URLs generating 4xx/5xx errors with frequency counts
- **Large Requests**: Biggest requests by response size to identify heavy resources
- **Enhanced Bot Detection**: More detailed bot breakdown and identification

### 📈 Use Cases

**For System Administrators**:
```bash
# Daily log analysis with export for reporting
./smart-log-analyser analyse /var/log/nginx/access.log* --export-csv=daily_report.csv --details

# Monitor error patterns and investigate issues
./smart-log-analyser analyse /var/log/nginx/access.log --details | grep "Error Analysis" -A 10
```

**For Security Analysis**:
```bash
# Export detailed data for security analysis tools
./smart-log-analyser analyse logs/*.log --export-json=security_analysis.json --details

# Focus on bot traffic and suspicious patterns
./smart-log-analyser analyse logs/*.log --details | grep -E "(Bot|Error)"
```

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
- **Multi-file mode** (default): Downloads all access log files found in the log directory (up to 10 files)
- **Single file mode** (`--single`): Downloads only the configured `log_path` file
- **Smart naming**: Files are saved as `hostname_timestamp_originalname` to avoid conflicts
- **Progress tracking**: Shows download progress for each file with size information
- **Configurable limit**: Use `--max-files` to control how many files to download

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