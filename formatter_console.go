package rz

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/gookit/color"
)

func colorize(c color.Color, str string) string {
	code := c.String()
	if len(code) == 0 || str == "" {
		return str
	}

	return fmt.Sprintf(color.FullColorTpl, code, str)
}

// FormatterConsole prettify output for human cosumption
func FormatterConsole() LogFormatter {
	return func(writer io.Writer, ev *Event) error {
		event, err := ev.Decode()
		if err != nil {
			return err
		}

		lvlColor := color.FgDefault
		level := "????"
		if l, ok := event[DefaultLevelFieldName].(string); ok {
			lvlColor = levelColor(l)
			level = strings.ToUpper(l)[0:4]
		}

		message := ""
		if m, ok := event[DefaultMessageFieldName].(string); ok {
			message = m
		}

		timestamp := ""
		if t, ok := event[DefaultTimestampFieldName].(string); ok {
			timestamp = t
		}

		if _, err = fmt.Fprintf(writer, "%-20s", timestamp); err != nil {
			return err
		}

		if _, err = fmt.Fprintf(writer, " |%-4s|", colorize(lvlColor, level)); err != nil {
			return err
		}

		if message != "" {
			_, err = fmt.Fprint(writer, " "+message)
		}
		if err != nil {
			return err
		}

		fields := make([]string, 0, len(event))
		for field := range event {
			switch field {
			case DefaultTimestampFieldName, DefaultMessageFieldName, DefaultLevelFieldName:
				continue
			}

			fields = append(fields, field)
		}

		sort.Strings(fields)
		for _, field := range fields {
			if needsQuote(field) {
				field = strconv.Quote(field)
			}

			_, err = fmt.Fprintf(writer, " %s=", colorize(lvlColor, field))

			switch value := event[field].(type) {
			case string:
				if len(value) == 0 {
					_, err = fmt.Fprint(writer, "\"\"")
				} else if needsQuote(value) {
					_, err = fmt.Fprint(writer, strconv.Quote(value))
				} else {
					_, err = fmt.Fprint(writer, value)
				}
			default:
				if b, e := json.Marshal(value); e == nil {
					_, err = fmt.Fprint(writer, string(b))
				}
			}
			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintln(writer)
		if err != nil {
			return err
		}

		return nil
	}
}
