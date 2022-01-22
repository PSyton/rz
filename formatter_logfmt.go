package rz

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
)

// FormatterLogfmt prettify output for human consumption, using the logfmt format.
func FormatterLogfmt() LogFormatter {
	return func(writer io.Writer, ev *Event) error {
		event, err := ev.Decode()
		if err != nil {
			return err
		}

		fields := make([]string, 0, len(event))
		for field := range event {
			fields = append(fields, field)
		}

		sort.Strings(fields)
		for _, field := range fields {
			if needsQuote(field) {
				field = strconv.Quote(field)
			}

			if _, err = fmt.Fprintf(writer, " %s=", field); err != nil {
				return err
			}

			switch value := event[field].(type) {
			case string:
				if len(value) == 0 {
					_, err = fmt.Fprint(writer, "\"\"")
				} else if needsQuote(value) {
					_, err = fmt.Fprint(writer, strconv.Quote(value))
				} else {
					_, err = fmt.Fprint(writer, value)
				}
			case time.Time:
				_, err = fmt.Fprint(writer, value.Format(time.RFC3339))
			default:
				b, err := json.Marshal(value)
				if err != nil {
					return err
				}
				_, err = fmt.Fprint(writer, strconv.Quote(string(b)))
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
