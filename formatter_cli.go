package rz

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/gookit/color"
)

// FormatterCLI prettify output suitable for command-line interfaces.
func FormatterCLI() LogFormatter {
	return func(writer io.Writer, ev *Event) error {
		event, err := ev.Decode()
		if err != nil {
			return err
		}

		lvlColor := color.FgDefault
		level := ""
		if l, ok := event[ev.levelFieldName].(string); ok {
			lvlColor = levelColor(l)
			level = l
		}

		message := ""
		if m, ok := event[ev.messageFieldName].(string); ok {
			message = m
		}

		if level != "" {
			if _, err = fmt.Fprint(writer, colorize(lvlColor, levelSymbol(level))); err != nil {
				return err
			}
		}

		writer.Write([]byte(message))

		fields := make([]string, 0, len(event))
		for field := range event {
			switch field {
			case ev.timestampFieldName, ev.messageFieldName, ev.levelFieldName:
				continue
			}

			fields = append(fields, field)
		}

		sort.Strings(fields)
		for _, field := range fields {
			if needsQuote(field) {
				field = strconv.Quote(field)
			}
			if _, err := fmt.Fprintf(writer, " %s=", colorize(lvlColor, field)); err != nil {
				return err
			}

			switch value := event[field].(type) {
			case string:
				if len(value) == 0 {
					_, err = writer.Write([]byte("\"\""))
				} else if needsQuote(value) {
					_, err = writer.Write([]byte(strconv.Quote(value)))
				} else {
					_, err = writer.Write([]byte(value))
				}
			case time.Time:
				_, err = writer.Write([]byte(value.Format(time.RFC3339)))
			default:
				b, err := json.Marshal(value)
				if err != nil {
					return err
				}
				_, err = fmt.Fprint(writer, b)
			}

			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintln(writer)
		if err != nil {
			return err
		}

		return err
	}
}

func levelSymbol(level string) string {
	switch level {
	case "info":
		return "✔ "
	case "warning":
		return "⚠ "
	case "error", "fatal":
		return "✘ "
	default:
		return "• "
	}
}
