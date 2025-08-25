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
- [x] **Response time analysis and percentiles** (P50, P95, P99 using response size as proxy)
- [x] **Geographic IP analysis** (country/region detection, private network identification)
- [x] **Advanced security analysis** (attack pattern detection, anomaly detection, threat scoring)
- [x] **Compressed file support** (automatic .gz decompression, rotated log files)

### Phase 3 (Advanced Analytics) 🚀
- [x] **HTML report generation with embedded charts** (Interactive reports with Chart.js visualizations)
- [x] **Interactive menu system** (User-friendly guided interface with dual-mode operation)
- [x] **ASCII charts and terminal visualizations** (Professional terminal-based charts with color support)
- [x] **Historical trend analysis** (Compare periods, track degradation, automated alerting)
- [x] **Advanced query language for complex filtering** (SQL-like query language with filtering, aggregation, and functions)
- [ ] Database integration (SQLite, PostgreSQL export)
- [ ] Plugin architecture for custom analyzers

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

## Project Structure

The Smart Log Analyser uses the following folder structure:

```
smart-log-analyser/
├── config/          # Configuration files (future use)
├── downloads/       # Downloaded log files from remote servers
├── output/          # Generated reports and export files
├── testdata/        # Sample log files for testing
├── pkg/            # Go packages (parser, analyser, remote)
├── cmd/            # CLI command implementations
└── scripts/        # Utility scripts (security checks, etc.)
```

**Important**: The `downloads/` and `output/` folders are excluded from git to prevent accidentally committing sensitive log files or large output files.

## Interactive Menu System 🎯

Smart Log Analyser now features an **interactive menu system** that launches when you run the program without arguments. This provides a user-friendly interface for all operations.

### Launching the Interactive Menu
```bash
# Simply run the program to launch the interactive menu
./smart-log-analyser
```

### Menu Features
- **📂 Analyse Local Log Files**: Browse and select log files with guided analysis
- **🌐 Download & Analyse Remote Logs**: Manage remote server connections and downloads
- **📈 Generate HTML Report**: Create interactive reports with custom settings
- **🔧 Configuration & Setup**: Configure servers, preferences, and export settings
- **📚 Help & Documentation**: Built-in help and guidance
- **🚪 Exit**: Clean exit from the application

### Interactive Workflows
The menu system guides you through:
- **File Selection**: Browse directories, use wildcards, or manually enter paths
- **Time Range Filtering**: Set custom date/time ranges for analysis
- **Advanced Analytics**: Access to all Phase 3 features through guided interface
- **Results Processing**: Choose from multiple analysis and visualization options
- **Export Options**: Choose from HTML, JSON, CSV formats with custom settings
- **Progress Tracking**: Real-time progress indicators for long operations

### Enhanced Results Menu
After log analysis completes, the interactive system offers comprehensive options:
```
📊 Results Options:
1. Show ASCII charts                              - Terminal visualizations
2. Export results                                 - HTML/JSON/CSV export  
3. Trend analysis & degradation detection         - Historical analysis
4. Combined analysis (charts + trends + export)   - All-in-one workflow
5. Continue                                       - Return to main menu
```

**Advanced Features Available:**
- **ASCII Visualizations**: Interactive charts with customizable width and colors
- **Trend Analysis**: Automated degradation detection with risk scoring
- **Combined Workflows**: Seamless integration of all analysis types
- **Smart Validation**: Helpful feedback about data requirements and limitations

## Quick Start

### Interactive Mode (Recommended)
```bash
# Launch interactive menu for guided experience
./smart-log-analyser

# Follow the prompts to:
# 1. Select "Analyse Local Log Files"
# 2. Choose your log files
# 3. Configure analysis options
# 4. Generate reports
```

### Command Line Mode
```bash
# Analyse a single log file
./smart-log-analyser analyse /var/log/nginx/access.log

# Analyse multiple log files together
./smart-log-analyser analyse /var/log/nginx/access.log /var/log/nginx/access.log.1

# Analyse compressed log files (.gz)
./smart-log-analyser analyse /var/log/nginx/access.log.1.gz

# Analyse all downloaded files using wildcard (supports compressed files)
./smart-log-analyser analyse ./downloads/*.log
./smart-log-analyser analyse ./downloads/*.gz

# Analyse all log files including compressed and rotated logs
./smart-log-analyser analyse /var/log/nginx/access.log* /var/log/nginx/error.log*

# Filter by time range
./smart-log-analyser analyse /var/log/nginx/access.log --since="2024-08-20 00:00:00" --until="2024-08-20 23:59:59"

# Get top 10 IPs and URLs
./smart-log-analyser analyse /var/log/nginx/access.log --top-ips=10 --top-urls=10

# Test with sample data
./smart-log-analyser analyse testdata/sample_access.log

# Show detailed breakdown with individual status codes and error analysis
./smart-log-analyser analyse /var/log/nginx/access.log --details

# Export results for further analysis (files saved to output/ folder)
./smart-log-analyser analyse ./downloads/*.log --export-json=output/detailed_report.json --export-csv=output/summary.csv

# Generate interactive HTML report with charts and visualizations
./smart-log-analyser analyse ./downloads/*.log --export-html=output/report.html --html-title="Production Server Analysis"

# Generate multiple export formats simultaneously
./smart-log-analyser analyse ./downloads/*.log --export-html=output/report.html --export-json=output/data.json --export-csv=output/summary.csv --details

# Display ASCII charts in terminal for immediate visual feedback
./smart-log-analyser analyse /var/log/nginx/access.log* --ascii-charts

# Customize ASCII chart display (width, colors, top results)
./smart-log-analyser analyse ./logs/*.log --ascii-charts --chart-width=100 --no-colors --top-ips=5

# Perform historical trend analysis and degradation detection
./smart-log-analyser analyse /var/log/nginx/access.log* --trend-analysis

# Combine trend analysis with visual charts for comprehensive insights
./smart-log-analyser analyse ./logs/*.log --trend-analysis --ascii-charts --chart-width=100

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

## ASCII Charts 📈

Visual terminal-based charts for immediate feedback without external tools. Perfect for SSH sessions and DevOps workflows.

### Features
- **Professional Terminal Charts**: Clean bar charts with proper scaling and labels
- **Color Intelligence**: Automatic color detection with graceful fallbacks for non-color terminals
- **SSH-Friendly**: Works perfectly over remote terminal connections
- **Flexible Sizing**: Adjustable chart width (60-100+ columns) for different terminal sizes
- **Multiple Chart Types**: Status codes, traffic analysis, geographic distribution, top IPs/URLs

### Chart Types
- **HTTP Status Code Distribution**: Color-coded by status type (2xx=green, 4xx=red, 5xx=magenta)
- **Human vs Bot Traffic**: Clear visualization of automated vs human requests
- **Top IP Addresses**: Traffic volume visualization with IP address display
- **Top URLs**: Request count charts with smart URL path truncation
- **Geographic Distribution**: Local/CDN/International traffic breakdown

### ASCII Chart Usage
```bash
# Basic ASCII charts with standard 80-column width
./smart-log-analyser analyse access.log --ascii-charts

# Wide charts for large terminals
./smart-log-analyser analyse access.log --ascii-charts --chart-width=100

# Disable colors for plain terminals or when piping output
./smart-log-analyser analyse access.log --ascii-charts --no-colors

# Combine with other options
./smart-log-analyser analyse access.log --ascii-charts --top-ips=10 --details
```

### Interactive Menu Integration
The ASCII charts are fully integrated into the interactive menu system:
1. Run analysis: `./smart-log-analyser analyse logs/`
2. Select: **[1] Analyse Local Log Files** → choose analysis option
3. After analysis completes, choose from **Results Options**:
   - **[1] Show ASCII Charts** - Full visual analysis with interactive options
   - **[3] Trend Analysis & Degradation Detection** - Historical trend analysis
   - **[4] Combined Analysis** - ASCII charts + trends + export options
4. Configure: Chart width (80/100) and color preferences in visualization menus

## Historical Trend Analysis 📈

Advanced trend analysis for detecting performance degradation and comparing different time periods automatically.

### Features
- **Automated Degradation Detection**: Identifies declining performance patterns using statistical analysis
- **Period-over-Period Comparison**: Compares metrics between different time segments automatically
- **Risk Assessment**: 0-100 risk scoring with actionable recommendations  
- **Smart Alerting**: Configurable thresholds with severity-based alert classification
- **Statistical Validation**: Minimum sample size requirements and confidence level assessment

### Analysis Types
- **Request Volume Trends**: Traffic pattern changes and volume fluctuations
- **Error Rate Monitoring**: 4xx/5xx response code trend analysis with threshold alerts
- **Performance Degradation**: Response time proxy analysis via response size tracking
- **Bot Traffic Analysis**: Automated vs human traffic pattern changes
- **Geographic Shifts**: Traffic source and distribution pattern analysis

### Visualization Integration
The trend analysis includes rich ASCII visualizations when combined with `--ascii-charts`:
- **Risk Score Gauge**: Horizontal gauge showing current risk level (0-100)
- **Degradation Alerts Chart**: Alert distribution by severity level
- **Metric Change Visualization**: Percentage changes across key performance indicators
- **Health Status Display**: Color-coded system health assessment

### Trend Analysis Usage
```bash
# Basic trend analysis with automatic period detection
./smart-log-analyser analyse access.log --trend-analysis

# Enhanced analysis with visual charts
./smart-log-analyser analyse access.log --trend-analysis --ascii-charts

# Comprehensive analysis with all features
./smart-log-analyser analyse access.log --trend-analysis --ascii-charts --details --chart-width=100

# Export results with trend analysis included
./smart-log-analyser analyse access.log --trend-analysis --export-html=report.html
```

### Sample Output
```
🏥 Overall Health: ⚠️ WARNING
📊 Analysis Type: degradation  
📈 Trend Summary: Analysis shows degrading trend with risk score 9/100

📋 Period Comparison:
├─ Overall Trend: 📉 degrading
├─ Risk Score: 9/100
└─ Key Changes: Bot Traffic decreased 5.1%

🚨 Degradation Alerts (1):
├─ Alert TREND-001: ⚠️ Bot Traffic
│  Impact: Impact requires investigation
└─ Recommendation: Monitor metric closely and investigate root causes

💡 Recommendations:
   1. Monitor metric closely and investigate root causes
```

### Configuration Parameters
- **Error Rate Threshold**: 10% increase triggers alerts (configurable)
- **Performance Threshold**: 20% response time increase detection
- **Traffic Drop Threshold**: 30% volume decrease monitoring
- **Minimum Sample Size**: 100 requests required for statistical validity
- **Risk Scoring**: Weighted analysis considering metric criticality and significance

### Interactive Menu Access
The trend analysis is seamlessly integrated into the interactive menu system:

**Menu Workflow:**
```
📊 Analysis Complete!
├─ Total Requests: 3,169
├─ Unique IPs: 1,413
└─ Time Range: 2025-08-23 00:00 to 2025-08-23 06:35

📊 Results Options:
1. Show ASCII charts
2. Export results
3. Trend analysis & degradation detection          ← Access trend analysis here
4. Combined analysis (charts + trends + export)    ← Comprehensive option
5. Continue
```

**Enhanced User Experience:**
- **Guided Interface**: No CLI knowledge required for advanced analysis
- **Smart Validation**: Automatic data sufficiency checking with helpful feedback
- **Progressive Disclosure**: Visualization options revealed after analysis completion
- **Professional Display**: Consistent formatting with emojis and structured output

## HTML Reports 📊

The Smart Log Analyser can generate beautiful, interactive HTML reports with embedded charts and visualizations.

### Features
- **Interactive Charts**: Powered by Chart.js with responsive design
- **Professional Layout**: Clean, modern interface with Bootstrap CSS
- **Comprehensive Analysis**: All analytics displayed with visual charts
- **Mobile Friendly**: Responsive design works on all devices
- **Print Ready**: Optimized styling for PDF generation

### Chart Types
- **Traffic Analysis**: Pie charts for human vs bot traffic
- **Status Code Distribution**: Doughnut charts for response codes
- **Hourly Traffic Patterns**: Line charts showing traffic over time
- **Response Size Analysis**: Bar charts for response time proxies
- **Geographic Distribution**: Bar charts for traffic by region
- **File Type Analysis**: Stacked bar charts for content types

### HTML Report Generation
```bash
# Basic HTML report
./smart-log-analyser analyse logs/ --export-html=output/report.html

# Custom title and detailed analysis
./smart-log-analyser analyse logs/ --export-html=output/report.html --html-title="Production Analysis" --details

# Multiple formats with HTML report
./smart-log-analyser analyse logs/ --export-html=output/report.html --export-json=output/data.json --details
```

### Opening HTML Reports
```bash
# Open in default browser (Linux)
xdg-open output/report.html

# Open in default browser (macOS)
open output/report.html

# Open in default browser (Windows)
start output/report.html
```

## Dual Operation Modes

Smart Log Analyser supports both **interactive menu mode** and **traditional CLI mode**:

### When to Use Each Mode

**🎯 Interactive Mode** - Best for:
- New users learning the system
- Complex analysis with multiple options
- Guided report generation
- Configuration and setup tasks
- Exploring available features

**⚡ CLI Mode** - Best for:
- Automation and scripting
- Batch processing workflows
- Integration with CI/CD pipelines
- Power users who know exact commands
- Remote server usage via SSH

### Switching Between Modes
```bash
# Interactive mode (menu-driven)
./smart-log-analyser

# CLI mode (direct commands)
./smart-log-analyser analyse logs/ --export-html=output/report.html
./smart-log-analyser download --test
./smart-log-analyser --help
```

Both modes provide access to the same powerful analysis features, just with different user interfaces optimized for different use cases.

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

⏱️  Response Size Analysis (Proxy for Response Time)
├─ Average Response: 512.3 KB
├─ Median (P50): 234.5 KB
├─ 95th Percentile: 1.2 MB
├─ 99th Percentile: 2.8 MB
├─ Range: 128 B - 3.4 MB
├─ Slowest Endpoints (by size):
│  ├─ /api/large-report: 3.4 MB
│  ├─ /downloads/document.pdf: 2.1 MB
└─ Fastest Endpoints (by size):
   ├─ /api/status: 128 B
   ├─ /health: 156 B

🌍 Geographic Distribution
├─ Local/Private: 1,247 (29.0%)
├─ Cloudflare CDN: 892 (20.8%)
├─ Countries (15 total):
│  ├─ United States: 1,156 requests (26.9%)
│  ├─ United Kingdom: 234 requests (5.4%)
│  ├─ Australia/NZ: 189 requests (4.4%)
│  ├─ Germany: 167 requests (3.9%)
│  ├─ Canada: 143 requests (3.3%)
└─ Regions:
   ├─ North America: 1,299 requests (30.2%)
   ├─ Europe: 567 requests (13.2%)
   ├─ Asia-Pacific: 234 requests (5.4%)
   ├─ Oceania: 189 requests (4.4%)

🔐 Security Analysis (Threat Level: LOW, Score: 92/100)
├─ Total Threats Detected: 12
├─ Suspicious IPs: 3
├─ Anomalies Detected: 1
├─ Attack Breakdown:
│  ├─ SQL Injection: 4 attempts
│  ├─ XSS Attempts: 2
│  ├─ Directory Traversal: 3 attempts
│  ├─ Brute Force: 2 attempts
│  ├─ Scanning Activity: 1 instances
├─ Top Threat IPs:
│  ├─ 203.0.113.1: 15 requests (Score: 45, sql_injection, scanner)
│  ├─ 198.51.100.42: 12 requests (Score: 35, xss, directory_traversal)
└─ Recent High-Severity Threats:
   ├─ [14:23:15] Sql Injection from 203.0.113.1
   │   URL: /search?q=admin' OR 1=1--
   │   Pattern: Boolean-based injection
   ├─ [14:25:42] Directory Traversal from 198.51.100.42
   │   URL: /files/../../../../etc/passwd
   │   Pattern: Unix-style traversal (../)

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

### File Format Support

The Smart Log Analyser can process various log file formats:

- **Regular log files**: `access.log`, `error.log`, `nginx.log`
- **Rotated log files**: `access.log.1`, `access.log.2`, `forum.access.log.5`
- **Compressed files**: `access.log.1.gz`, `error.log.14.gz`, `forum.access.log.5.gz`
- **Complex naming**: `example.com.access.log`, `site.error.log.11.gz`

**Supported compression formats:**
- ✅ **Gzip (.gz)** - Automatic decompression during analysis
- 📋 **Bzip2 (.bz2)** - Planned for future releases

**File detection patterns:**
- Files ending with `.log`, `.access`, `.error` (with optional numbers)
- Compressed variants with `.gz` extension
- Mixed patterns like `site.access.log.12.gz`

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

## Advanced Query Language (SLAQ)

The Smart Log Analyser Query Language provides powerful SQL-like capabilities for filtering and analyzing log data with custom queries.

### Basic Syntax
```sql
SELECT [fields] FROM logs WHERE [conditions] [GROUP BY field] [ORDER BY field] [LIMIT number]
```

### Query Examples

**Basic Filtering:**
```bash
# Find all 404 errors
./smart-log-analyser analyse access.log --query "SELECT * FROM logs WHERE status = 404"

# Filter by IP address
./smart-log-analyser analyse access.log --query "SELECT ip, url, status FROM logs WHERE ip = '192.168.1.100'"
```

**Aggregation Analysis:**
```bash  
# Top IP addresses by request count
./smart-log-analyser analyse access.log --query "SELECT ip, COUNT() FROM logs GROUP BY ip ORDER BY COUNT() DESC LIMIT 10"

# Error analysis by URL
./smart-log-analyser analyse access.log --query "SELECT url, status, COUNT() FROM logs WHERE status >= 400 GROUP BY url, status"
```

**Time-based Analysis:**
```bash
# Hourly traffic distribution
./smart-log-analyser analyse access.log --query "SELECT HOUR(timestamp), COUNT() FROM logs GROUP BY HOUR(timestamp) ORDER BY HOUR(timestamp)"

# Requests by day of week
./smart-log-analyser analyse access.log --query "SELECT WEEKDAY(timestamp), COUNT() FROM logs GROUP BY WEEKDAY(timestamp)"
```

**Complex Filtering:**
```bash
# Large API requests with errors
./smart-log-analyser analyse access.log --query "SELECT url, method, size, status FROM logs WHERE status >= 400 AND size > 10000 AND url LIKE '/api*'"

# Bot traffic analysis
./smart-log-analyser analyse access.log --query "SELECT user_agent, COUNT() FROM logs WHERE user_agent CONTAINS 'bot' GROUP BY user_agent"
```

### Available Fields
- `ip` - Client IP address
- `timestamp` - Request timestamp  
- `method` - HTTP method (GET, POST, etc.)
- `url` - Request URL/path
- `protocol` - HTTP protocol version
- `status` - HTTP status code
- `size` - Response size in bytes
- `referer` - HTTP referer header
- `user_agent` - User agent string

### Available Functions

**Aggregate Functions:**
- `COUNT()` - Count records
- `SUM(field)` - Sum numeric values
- `AVG(field)` - Average of numeric values
- `MIN(field)` - Minimum value
- `MAX(field)` - Maximum value

**Time Functions:**
- `HOUR(timestamp)` - Extract hour (0-23)
- `DAY(timestamp)` - Extract day of month
- `WEEKDAY(timestamp)` - Extract weekday (0=Sunday)
- `DATE(timestamp)` - Extract date part

**String Functions:**
- `UPPER(field)` - Convert to uppercase
- `LOWER(field)` - Convert to lowercase
- `LENGTH(field)` - String length

### Available Operators

**Comparison:**
- `=`, `!=`, `<>` - Equality/inequality
- `<`, `<=`, `>`, `>=` - Numeric comparisons

**String Matching:**
- `LIKE` - Pattern matching (`*` = multiple chars, `?` = single char)
- `CONTAINS` - String contains substring
- `STARTS_WITH` - String starts with prefix
- `ENDS_WITH` - String ends with suffix

**Logical:**
- `AND`, `OR`, `NOT` - Logical operations
- `IN` - Value in list: `status IN (200, 201, 202)`
- `BETWEEN` - Value in range: `size BETWEEN 1000 AND 5000`

### Output Formats

**Table Format (default):**
```bash
./smart-log-analyser analyse access.log --query "SELECT ip, COUNT() FROM logs GROUP BY ip LIMIT 3"
```

**CSV Format:**
```bash
./smart-log-analyser analyse access.log --query "SELECT ip, COUNT() FROM logs GROUP BY ip" --query-format csv
```

**JSON Format:**
```bash
./smart-log-analyser analyse access.log --query "SELECT ip, COUNT() FROM logs GROUP BY ip" --query-format json
```
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
- **Response Time Analysis**: Percentile analysis using response size as proxy (P50, P95, P99)
- **Geographic Distribution**: Country and region breakdown of IP addresses with CDN detection
- **Security Threat Analysis**: Detailed attack patterns, IP threat scoring, and anomaly detection
- **Compressed File Support**: Seamless processing of .gz files and rotated logs

### 📈 Use Cases

**For System Administrators**:
```bash
# Daily log analysis with export for reporting (including compressed files)
./smart-log-analyser analyse /var/log/nginx/access.log* --export-csv=daily_report.csv --details

# Analyse all log files including rotated and compressed
./smart-log-analyser analyse /var/log/nginx/*.log /var/log/nginx/*.log.* /var/log/nginx/*.gz --details

# Monitor error patterns and investigate issues
./smart-log-analyser analyse /var/log/nginx/access.log --details | grep "Error Analysis" -A 10

# Process downloaded compressed logs efficiently
./smart-log-analyser analyse ./downloads/*.gz --export-json=compressed_analysis.json
```

**For Security Analysis**:
```bash
# Comprehensive security analysis with threat detection
./smart-log-analyser analyse logs/*.log logs/*.gz --export-json=security_analysis.json --details

# Monitor for specific attack patterns and anomalies
./smart-log-analyser analyse /var/log/nginx/*.log* --details | grep -E "(Security|Threat|Anomaly)"

# Process compressed security logs for incident analysis
./smart-log-analyser analyse ./incident_logs/*.gz --details --export-csv=security_report.csv

# Real-world security monitoring example
./smart-log-analyser analyse /var/log/nginx/access.log* /var/log/nginx/error.log* --details
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

### Performance Features

**Compressed File Optimization**:
- Automatic detection of file types based on extensions
- Streaming decompression for memory efficiency
- Large buffer support (1MB) for processing big compressed files
- Concurrent processing capability for multiple files

**File Processing Capabilities**:
- Handles files of any size through streaming
- Memory-efficient processing of compressed archives
- Supports mixed file types in single analysis session
- Robust error handling for corrupted or incomplete files

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
- **v0.2.0**: Advanced analytics and export features ✅
- **v0.2.1**: Security analysis and compressed file support ✅  
- **v0.3.0**: Advanced analytics and visualizations
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
- Use `./scripts/check-sensitive-data.sh` before commits to scan for sensitive data

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