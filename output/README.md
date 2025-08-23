# Output Directory

This directory contains generated reports and export files from the Smart Log Analyser.

## Generated Files

The following types of files are created in this directory:

### Export Files
- `*.csv` - Summary data in CSV format for spreadsheet analysis
- `*.json` - Detailed analysis results in JSON format for programmatic use
- `detailed_report.json` - Comprehensive analysis with all metrics and breakdowns
- `summary.csv` - High-level statistics and summaries

### Report Files (Future)
- `*.html` - HTML reports with embedded charts and visualizations
- `*.pdf` - Printable PDF reports
- `analysis_*.txt` - Plain text analysis summaries

## Usage Examples

```bash
# Export detailed JSON report
./smart-log-analyser analyse logs/ --export-json=output/analysis_$(date +%Y%m%d).json

# Export CSV summary
./smart-log-analyser analyse logs/ --export-csv=output/summary_$(date +%Y%m%d).csv

# Export both formats
./smart-log-analyser analyse logs/ --export-json=output/report.json --export-csv=output/report.csv --details
```

## Security Note

This directory is excluded from git (.gitignore) to prevent accidentally committing:
- Sensitive log analysis data
- Large output files
- Reports containing IP addresses or other potentially sensitive information

Always review export files before sharing to ensure they don't contain sensitive data.