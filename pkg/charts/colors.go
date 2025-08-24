package charts

import (
	"os"
	"strings"
)

// Terminal color codes
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
	
	// Foreground colors
	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	
	// Bright colors
	ColorBrightBlack   = "\033[90m"
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightWhite   = "\033[97m"
	
	// Background colors  
	ColorBgBlack   = "\033[40m"
	ColorBgRed     = "\033[41m"
	ColorBgGreen   = "\033[42m"
	ColorBgYellow  = "\033[43m"
	ColorBgBlue    = "\033[44m"
	ColorBgMagenta = "\033[45m"
	ColorBgCyan    = "\033[46m"
	ColorBgWhite   = "\033[47m"
)

// Chart color palette for different data types
var (
	StatusCodeColors = map[string]string{
		"2xx": ColorGreen,      // Success - Green
		"3xx": ColorYellow,     // Redirect - Yellow  
		"4xx": ColorRed,        // Client Error - Red
		"5xx": ColorMagenta,    // Server Error - Magenta
		"1xx": ColorCyan,       // Info - Cyan
	}
	
	TrafficColors = []string{
		ColorBlue,
		ColorCyan,
		ColorGreen,
		ColorYellow,
		ColorMagenta,
		ColorRed,
		ColorBrightBlue,
		ColorBrightCyan,
	}
	
	DefaultBarColor = ColorBlue
)

// SupportsColor checks if the terminal supports color output
func SupportsColor() bool {
	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "" {
		return false
	}
	
	// Common terminals that support color
	colorTerms := []string{
		"xterm", "xterm-256color", "xterm-color",
		"screen", "screen-256color",
		"tmux", "tmux-256color",
		"linux", "cygwin",
	}
	
	for _, colorTerm := range colorTerms {
		if strings.Contains(term, colorTerm) {
			return true
		}
	}
	
	// Check for NO_COLOR environment variable (standard)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	
	// Check for FORCE_COLOR environment variable
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	
	return false
}

// GetStatusCodeColor returns the appropriate color for an HTTP status code
func GetStatusCodeColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return StatusCodeColors["2xx"]
	case statusCode >= 300 && statusCode < 400:
		return StatusCodeColors["3xx"]
	case statusCode >= 400 && statusCode < 500:
		return StatusCodeColors["4xx"]
	case statusCode >= 500 && statusCode < 600:
		return StatusCodeColors["5xx"]
	case statusCode >= 100 && statusCode < 200:
		return StatusCodeColors["1xx"]
	default:
		return DefaultBarColor
	}
}

// GetTrafficColor returns a color from the traffic color palette
func GetTrafficColor(index int) string {
	if index < 0 || index >= len(TrafficColors) {
		return DefaultBarColor
	}
	return TrafficColors[index%len(TrafficColors)]
}

// Colorize applies a color to text with automatic reset
func Colorize(text, color string) string {
	if !SupportsColor() || color == "" {
		return text
	}
	return color + text + ColorReset
}

// ColorizeBackground applies a background color to text
func ColorizeBackground(text, bgColor string) string {
	if !SupportsColor() || bgColor == "" {
		return text
	}
	return bgColor + text + ColorReset
}

// StripColors removes all ANSI color codes from a string
func StripColors(text string) string {
	// Simple regex would be better, but avoiding regex dependency
	result := text
	colorCodes := []string{
		ColorReset, ColorBold, ColorDim,
		ColorBlack, ColorRed, ColorGreen, ColorYellow, ColorBlue, ColorMagenta, ColorCyan, ColorWhite,
		ColorBrightBlack, ColorBrightRed, ColorBrightGreen, ColorBrightYellow,
		ColorBrightBlue, ColorBrightMagenta, ColorBrightCyan, ColorBrightWhite,
		ColorBgBlack, ColorBgRed, ColorBgGreen, ColorBgYellow, ColorBgBlue, ColorBgMagenta, ColorBgCyan, ColorBgWhite,
	}
	
	for _, code := range colorCodes {
		result = strings.ReplaceAll(result, code, "")
	}
	
	return result
}

// GetTerminalWidth attempts to determine the terminal width
func GetTerminalWidth() int {
	// Try to get from environment first
	if cols := os.Getenv("COLUMNS"); cols != "" {
		// Parse COLUMNS if available
	}
	
	// Default fallback width
	return 80
}