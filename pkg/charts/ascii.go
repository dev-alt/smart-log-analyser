package charts

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

// ChartConfig defines the configuration for ASCII charts
type ChartConfig struct {
	Width      int    // Terminal width for chart
	Height     int    // Chart height in lines
	Title      string // Chart title
	ShowColors bool   // Enable terminal colors
	ShowValues bool   // Show numeric values
	ShowPercent bool  // Show percentages
}

// BarData represents a single bar in a chart
type BarData struct {
	Label   string // Bar label
	Value   int64  // Bar value
	Percent float64 // Percentage of total
	Color   string // Terminal color code
}

// BarChart represents a horizontal bar chart
type BarChart struct {
	Config ChartConfig
	Data   []BarData
	Total  int64 // Total for percentage calculations
}

// NewBarChart creates a new horizontal bar chart
func NewBarChart(title string, width int) *BarChart {
	return &BarChart{
		Config: ChartConfig{
			Width:      width,
			Height:     20, // Default height
			Title:      title,
			ShowColors: true,
			ShowValues: true,
			ShowPercent: true,
		},
		Data: make([]BarData, 0),
	}
}

// AddBar adds a bar to the chart
func (c *BarChart) AddBar(label string, value int64, color string) {
	c.Data = append(c.Data, BarData{
		Label: label,
		Value: value,
		Color: color,
	})
	c.Total += value
}

// calculatePercentages calculates percentages for all bars
func (c *BarChart) calculatePercentages() {
	if c.Total == 0 {
		return
	}
	
	for i := range c.Data {
		c.Data[i].Percent = (float64(c.Data[i].Value) / float64(c.Total)) * 100
	}
}

// Render generates the ASCII bar chart as a string
func (c *BarChart) Render() string {
	if len(c.Data) == 0 {
		return "No data to display"
	}

	c.calculatePercentages()

	var result strings.Builder

	// Add title with emoji
	if c.Config.Title != "" {
		result.WriteString(fmt.Sprintf("📊 %s\n", c.Config.Title))
	}

	// Find the maximum value for scaling
	maxValue := int64(0)
	maxLabelWidth := 0
	for _, bar := range c.Data {
		if bar.Value > maxValue {
			maxValue = bar.Value
		}
		if len(bar.Label) > maxLabelWidth {
			maxLabelWidth = len(bar.Label)
		}
	}

	// Calculate available space for bars
	valueWidth := len(fmt.Sprintf("%d", maxValue))
	percentWidth := 6 // " 100.0%"
	
	// Calculate bar width
	totalFixedWidth := maxLabelWidth + 1 + valueWidth + percentWidth + 4 // spaces and formatting
	if c.Config.ShowValues && c.Config.ShowPercent {
		// Keep current calculation
	} else if c.Config.ShowValues || c.Config.ShowPercent {
		totalFixedWidth -= percentWidth // or valueWidth, depending on what we're not showing
	} else {
		totalFixedWidth -= valueWidth + percentWidth
	}

	barWidth := c.Config.Width - totalFixedWidth
	if barWidth < 20 {
		barWidth = 20 // Minimum bar width
	}

	// Generate each bar
	for _, bar := range c.Data {
		// Calculate bar length
		barLength := int(math.Round(float64(bar.Value) / float64(maxValue) * float64(barWidth)))
		if barLength == 0 && bar.Value > 0 {
			barLength = 1 // Ensure non-zero values get at least one character
		}

		// Generate the bar
		barStr := strings.Repeat("█", barLength)
		if barLength > 0 && barLength < barWidth {
			// Add partial character for more precision
			barStr += "▌"
		}

		// Apply colors if enabled
		if c.Config.ShowColors && bar.Color != "" {
			barStr = bar.Color + barStr + ColorReset
		}

		// Format the complete line
		line := fmt.Sprintf("%-*s ", maxLabelWidth, bar.Label)
		line += barStr

		// Add padding to align values
		currentBarLen := utf8.RuneCountInString(barStr)
		if c.Config.ShowColors && bar.Color != "" {
			// Account for color codes not being visible
			currentBarLen = barLength
			if barLength > 0 && barLength < barWidth {
				currentBarLen++ // for the partial character
			}
		}
		
		padding := barWidth - currentBarLen
		if padding > 0 {
			line += strings.Repeat(" ", padding)
		}

		// Add value and percentage
		if c.Config.ShowPercent && c.Config.ShowValues {
			line += fmt.Sprintf(" %*.1f%% (%*d)", 5, bar.Percent, valueWidth, bar.Value)
		} else if c.Config.ShowPercent {
			line += fmt.Sprintf(" %*.1f%%", 5, bar.Percent)
		} else if c.Config.ShowValues {
			line += fmt.Sprintf(" %*d", valueWidth, bar.Value)
		}

		result.WriteString(line + "\n")
	}

	return result.String()
}

// RenderSimple generates a simple bar chart without colors or detailed formatting
func (c *BarChart) RenderSimple() string {
	oldConfig := c.Config
	c.Config.ShowColors = false
	c.Config.ShowPercent = false
	result := c.Render()
	c.Config = oldConfig
	return result
}

// SetConfig updates the chart configuration
func (c *BarChart) SetConfig(config ChartConfig) {
	c.Config = config
}

// FormatNumber formats large numbers with appropriate units
func FormatNumber(num int64) string {
	if num < 1000 {
		return fmt.Sprintf("%d", num)
	} else if num < 1000000 {
		return fmt.Sprintf("%.1fK", float64(num)/1000)
	} else if num < 1000000000 {
		return fmt.Sprintf("%.1fM", float64(num)/1000000)
	}
	return fmt.Sprintf("%.1fB", float64(num)/1000000000)
}

// TruncateString truncates a string to a maximum length with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}