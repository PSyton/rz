package rz

import (
	"fmt"
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

func colorize(c color.Color, str string) string {
	code := c.String()
	if len(code) == 0 || str == "" {
		return str
	}

	if !color.Enable || !color.SupportColor() {
		return str
	}

	return fmt.Sprintf(color.FullColorTpl, code, str)
}
