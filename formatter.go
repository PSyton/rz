package rz

import (
	"io"

	"github.com/gookit/color"
)

// LogFormatter can be used to log to another format than JSON
type LogFormatter func(writer io.Writer, ev *Event) error

func levelColor(level string) color.Color {
	switch level {
	case "debug":
		return color.FgLightCyan
	case "info":
		return color.FgLightGreen
	case "warning":
		return color.FgLightYellow
	case "error", "fatal", "panic":
		return color.FgLightRed
	default:
		return color.FgDefault
	}
}

func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}
