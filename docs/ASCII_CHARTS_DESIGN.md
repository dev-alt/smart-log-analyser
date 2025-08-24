# ASCII Charts and Terminal Visualizations Design

**Created**: Session 19  
**Status**: Implementation Ready
**Complexity**: Medium ⭐⭐⭐

## Overview

Add terminal-based data visualization capabilities to complement the existing HTML reports. This feature provides immediate visual feedback in CLI environments where HTML reports aren't practical.

## Feature Requirements

### Core Visualizations
1. **Traffic Timeline** - Requests over time (hourly/daily)
2. **Status Code Distribution** - Bar chart of HTTP status codes  
3. **Top Resources** - Horizontal bar chart of most requested URLs
4. **Top IPs** - Horizontal bar chart of highest traffic sources
5. **Response Size Distribution** - Histogram of response sizes
6. **Geographic Distribution** - Simple text-based breakdown

### Display Features
- **Configurable Width**: Adapt to terminal width (80-120+ columns)
- **Color Support**: Use terminal colors for better readability
- **Scalable Data**: Handle datasets from dozens to millions of entries
- **Interactive Option**: Optional in menu system vs CLI flags

## Technical Architecture

### Package Structure
```
pkg/
├── charts/
│   ├── ascii.go          # Core ASCII chart engine
│   ├── timeline.go       # Traffic over time charts  
│   ├── bars.go          # Horizontal/vertical bar charts
│   ├── histogram.go     # Distribution charts
│   └── colors.go        # Terminal color support
```

### Integration Points
- **CLI Integration**: New `--ascii-charts` flag for analyse command
- **Menu Integration**: Charts display option in analysis results
- **Data Pipeline**: Uses existing `analyser.Results` struct
- **Output Options**: Standalone or combined with text reports

## Chart Types Implementation

### 1. Traffic Timeline
```
📈 Traffic Over Time (24 hours)
┌─────────────────────────────────────────────┐
│ 1.2k ┤                                 ╭─╮ │ 
│ 1.0k ┤                               ╭─╯ ╰╮│
│ 800  ┤                           ╭───╯    ││
│ 600  ┤                       ╭───╯        ││
│ 400  ┤                   ╭───╯            ││
│ 200  ┤               ╭───╯                ││
│   0  ┤───────────────╯                    ╰│
└─────────────────────────────────────────────┘
  00:00    06:00    12:00    18:00    24:00
```

### 2. Status Code Distribution  
```
📊 HTTP Status Codes
200 ████████████████████████████████████▌ 89.2% (2,827)
404 ████████▌ 7.1% (225)  
301 ██▌ 2.8% (89)
403 ▌ 0.6% (19)
500 ▌ 0.3% (9)
```

### 3. Top Resources
```
🔥 Top Requested URLs  
/api/users        ████████████████████████████████████▌ 1,234
/dashboard        ████████████████████████▌ 856  
/login           ████████████████▌ 567
/api/posts       ████████▌ 289
/settings        ███▌ 123
```

### 4. Geographic Distribution
```
🌍 Geographic Distribution
Local Networks    ████████████████████████████████████▌ 68.8% (2,180)
CDN/Proxy         ███████████▌ 25.2% (799)  
International     ██▌ 6.1% (193)
```

## Implementation Strategy

### Phase 1: Core Engine (Session 19)
1. **ASCII Chart Package** - Create base chart rendering engine
2. **Color Support** - Terminal color detection and styling
3. **Bar Charts** - Horizontal bar chart implementation
4. **Data Integration** - Connect with existing analysis results

### Phase 2: Advanced Charts (Future)
1. **Timeline Charts** - Traffic over time visualization
2. **Histograms** - Response size and time distributions  
3. **Interactive Features** - Scrolling, zooming for large datasets
4. **Export Options** - Save ASCII charts to files

### Phase 3: Menu Integration (Future)
1. **Analysis Results Display** - Show charts after analysis
2. **Chart Selection** - Choose which charts to display
3. **Configuration** - Terminal width, colors, detail level

## Technical Implementation Details

### Chart Engine Core
```go
type ChartConfig struct {
    Width       int    // Terminal width
    Height      int    // Chart height  
    Title       string // Chart title
    ShowColors  bool   // Enable terminal colors
    ShowValues  bool   // Show numeric values
}

type BarChart struct {
    Config ChartConfig
    Data   []BarData
}

type BarData struct {
    Label string
    Value int64
    Color string // Terminal color code
}

func (c *BarChart) Render() string {
    // Generate ASCII bar chart
}
```

### Color Support
```go
// Terminal color constants
const (
    ColorReset  = "\033[0m"
    ColorRed    = "\033[31m" 
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorBlue   = "\033[34m"
    ColorPurple = "\033[35m"
    ColorCyan   = "\033[36m"
    ColorWhite  = "\033[37m"
)
```

### Data Processing
```go
func GenerateStatusCodeChart(results *analyser.Results) *BarChart {
    // Process status code data into chart format
}

func GenerateTopIPsChart(results *analyser.Results, limit int) *BarChart {
    // Process IP data into horizontal bar chart
}
```

## CLI Integration

### New Command Flags
```bash
# Show ASCII charts with analysis
./smart-log-analyser analyse logs/ --ascii-charts

# Show specific chart types
./smart-log-analyser analyse logs/ --ascii-charts --charts=status,ips,urls

# Configure chart appearance  
./smart-log-analyser analyse logs/ --ascii-charts --chart-width=100 --no-colors
```

### Menu Integration
```
📊 Analysis Results Options:
1. Show detailed text report
2. Show ASCII charts  
3. Export HTML report
4. Export to files
5. Back to main menu
```

## Benefits

### User Experience
- **Immediate Visualization**: No need to open HTML files
- **SSH-Friendly**: Works perfectly over remote connections
- **Quick Analysis**: Fast visual overview of key metrics
- **Professional Output**: Clean, readable terminal graphics

### Technical Benefits
- **Low Dependencies**: Pure Go implementation, no external libraries
- **Fast Rendering**: Efficient ASCII generation
- **Flexible Output**: Configurable for different terminal sizes
- **Cross-Platform**: Works on any terminal that supports Unicode

## Success Criteria

### Functionality
- ✅ Generate bar charts for status codes, IPs, URLs
- ✅ Handle datasets of various sizes (10s to 10,000+ entries)  
- ✅ Automatic scaling and formatting
- ✅ Terminal color support with fallback

### Performance
- ✅ Chart generation under 100ms for typical datasets
- ✅ Memory usage under 10MB for chart rendering
- ✅ Clean terminal output without artifacts

### Integration
- ✅ Seamless CLI flag integration
- ✅ Optional menu system integration  
- ✅ Compatible with existing analysis pipeline
- ✅ Does not break existing functionality

## Future Enhancements

### Advanced Features
- **Interactive Charts**: Scroll through long datasets
- **Real-time Updates**: Live updating charts for ongoing analysis
- **Chart Export**: Save ASCII charts to text files
- **Custom Themes**: Different color schemes and styles

### Additional Chart Types
- **Sparklines**: Compact inline charts
- **Pie Charts**: Text-based pie chart representation
- **Scatter Plots**: For correlation analysis
- **Heat Maps**: Time-based activity visualization

---

*This design provides a solid foundation for terminal-based data visualization while maintaining the Smart Log Analyser's focus on performance and usability.*