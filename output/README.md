# Output Directory

This directory contains generated reports and export files from the Smart Log Analyser.

## Generated Files

The following types of files are created in this directory:

### Export Files
- `*.csv` - Summary data in CSV format for spreadsheet analysis
- `*.json` - Detailed analysis results in JSON format for programmatic use
- `detailed_report.json` - Comprehensive analysis with all metrics and breakdowns
- `summary.csv` - High-level statistics and summaries

### HTML Report Files
- `*.html` - Interactive HTML reports with embedded Chart.js visualizations
- `report.html` - Standard HTML report with all analytics and charts
- Professional responsive design suitable for stakeholders and presentations
- Self-contained files with embedded CSS/JavaScript for offline viewing

### Report Files (Future)
- `*.pdf` - Printable PDF reports (export from HTML reports)
- `analysis_*.txt` - Plain text analysis summaries

## Usage Examples

```bash
# Export interactive HTML report
./smart-log-analyser analyse logs/ --export-html=output/report.html --html-title="Production Analysis"

# Export detailed JSON report
./smart-log-analyser analyse logs/ --export-json=output/analysis_$(date +%Y%m%d).json

# Export CSV summary
./smart-log-analyser analyse logs/ --export-csv=output/summary_$(date +%Y%m%d).csv

# Export all formats simultaneously
./smart-log-analyser analyse logs/ --export-html=output/report.html --export-json=output/data.json --export-csv=output/summary.csv --details

# Interactive mode (guided HTML report generation)
./smart-log-analyser
# Then select: "3. 📈 Generate HTML Report"
```

## Security Note

This directory is excluded from git (.gitignore) to prevent accidentally committing:
- Sensitive log analysis data
- Large output files
- Reports containing IP addresses or other potentially sensitive information

Always review export files before sharing to ensure they don't contain sensitive data.