# Interactive HTML Reports

## Test the new interactive HTML reporting functionality

### Command line usage:
```bash
# Generate interactive HTML report (default)
./smart-log-analyser analyse access.log --export-html report.html

# Generate standard HTML report  
./smart-log-analyser analyse access.log --export-html report.html --interactive-html=false

# Custom title
./smart-log-analyser analyse access.log --export-html report.html --html-title "My Server Analysis"
```

### Interactive Menu usage:
1. Run `./smart-log-analyser menu`
2. Select analysis option
3. When exporting to HTML, you'll be prompted to choose:
   - Interactive Report (recommended) - Tabbed interface with drill-down
   - Standard Report - Simple static report

### Features:
- **Tabbed Interface**: Overview, Traffic Analysis, Error Analysis, Performance, Security, Geographic
- **Clickable Tables**: Click on any row to see detailed information
- **Real-time Filtering**: Filter IPs by type, errors by status code, etc.
- **Drill-down Information**: Detailed breakdowns for each data point
- **Professional Charts**: Interactive Chart.js visualizations
- **Status Code Analysis**: Click on status codes to see specific error details
- **IP Analysis**: View detailed IP information, geographic data, and traffic patterns
- **Error Analysis**: Comprehensive error breakdowns with fix suggestions

