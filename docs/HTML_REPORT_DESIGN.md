# HTML Report Generation Design Document

## Overview
Implementation plan for HTML report generation with embedded charts for the Smart Log Analyser.

## Report Structure

### Page Layout
```
┌─────────────────────────────────────────────────┐
│                   HEADER                        │
│ Smart Log Analyser Report | Generated: DATE     │
├─────────────────────────────────────────────────┤
│                  OVERVIEW                       │
│ Total Requests | Unique IPs | Data Transfer     │
│ Date Range | Analysis Duration                  │
├─────────────────────────────────────────────────┤
│              TRAFFIC ANALYSIS                   │
│ ┌─────────────┐  ┌────────────────────────────┐ │
│ │ Human/Bot   │  │    Hourly Traffic Chart    │ │
│ │ Pie Chart   │  │    (Line/Bar Chart)        │ │
│ └─────────────┘  └────────────────────────────┘ │
├─────────────────────────────────────────────────┤
│             STATUS CODE ANALYSIS                │
│ ┌─────────────┐  ┌────────────────────────────┐ │
│ │ Status Code │  │    Response Time Chart     │ │
│ │ Donut Chart │  │    (Histogram/Box Plot)    │ │
│ └─────────────┘  └────────────────────────────┘ │
├─────────────────────────────────────────────────┤
│            GEOGRAPHIC ANALYSIS                  │
│ ┌─────────────┐  ┌────────────────────────────┐ │
│ │ Geographic  │  │    Top IPs/URLs Tables     │ │
│ │ Bar Chart   │  │    (Interactive Tables)    │ │
│ └─────────────┘  └────────────────────────────┘ │
├─────────────────────────────────────────────────┤
│             SECURITY ANALYSIS                   │
│ ┌─────────────┐  ┌────────────────────────────┐ │
│ │ Threat      │  │    Attack Timeline         │ │
│ │ Level Gauge │  │    (Timeline Chart)        │ │
│ └─────────────┘  └────────────────────────────┘ │
├─────────────────────────────────────────────────┤
│                  DETAILED TABLES                │
│ • Top URLs by Traffic                          │
│ • Error Analysis (4xx/5xx)                     │
│ • Bot Detection Results                        │
│ • File Type Breakdown                          │
└─────────────────────────────────────────────────┘
```

## Chart Libraries

### Chart.js Integration
- **Pros**: Lightweight, responsive, excellent documentation
- **Charts**: Pie, Bar, Line, Doughnut, Radar
- **Size**: ~60KB minified
- **CDN**: Available via CDN for easy integration

### Chart Types Mapping
```
Analysis Feature → Chart Type
├─ Bot vs Human Traffic → Pie Chart
├─ Status Code Distribution → Doughnut Chart  
├─ Hourly Traffic Pattern → Line Chart
├─ Response Size Distribution → Histogram
├─ Geographic Distribution → Bar Chart
├─ File Type Analysis → Stacked Bar Chart
├─ Security Threat Level → Gauge Chart
└─ Top IPs/URLs → Table with sorting
```

## Implementation Strategy

### Phase 1: Basic HTML Template
1. Create HTML template with Bootstrap CSS framework
2. Add Chart.js via CDN
3. Generate static HTML with embedded data
4. Basic styling and responsive design

### Phase 2: Interactive Charts
1. Dynamic chart generation from JSON data
2. Interactive tooltips and legends
3. Responsive chart resizing
4. Print-friendly CSS

### Phase 3: Advanced Features
1. Export chart images (PNG/SVG)
2. Interactive filtering/drill-down
3. Dark/light theme toggle
4. Report customization options

## File Structure
```
pkg/
├─ html/
│  ├─ templates/
│  │  ├─ report.html        # Main report template
│  │  ├─ charts.js          # Chart generation JavaScript
│  │  └─ styles.css         # Custom CSS styles
│  └─ generator.go          # HTML generation logic
```

## Data Flow
```
Analysis Results → HTML Generator → Template Engine → Static HTML File
                                   ↓
                              Chart.js Data → Interactive Charts
```

## CLI Integration
```bash
# Generate HTML report
./smart-log-analyser analyse logs/ --export-html=output/report.html

# Generate with custom title
./smart-log-analyser analyse logs/ --export-html=output/report.html --title="Production Server Analysis"

# Generate with charts disabled (faster)
./smart-log-analyser analyse logs/ --export-html=output/report.html --no-charts
```

## Security Considerations
- Sanitize all data before HTML output
- Use Content Security Policy headers
- Avoid inline JavaScript with user data
- Ensure no sensitive data leaks in HTML comments

## Browser Compatibility
- Modern browsers (Chrome, Firefox, Safari, Edge)
- Mobile responsive design
- Graceful degradation for older browsers
- Print stylesheet for PDF generation
- Chart.js UMD version for compatibility with script tags
- Fixed chart dimensions to prevent excessive growth

## Performance Targets
- HTML generation: < 2 seconds
- Page load time: < 3 seconds
- Chart rendering: < 1 second per chart
- File size: < 500KB (without log data)