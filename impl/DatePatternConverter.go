package impl

import (
	"strings"
	"time"
)

type DatePatternConverter struct {
	AbstractPatternConverter
	layout   string
	location *time.Location
}

func NewDatePatternConverter(pattern string, location *time.Location, formatting FormattingInfo) *DatePatternConverter {
	result := &DatePatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formatting),
		location:                 location,
	}

	result.layout = result.toTimeLayout(pattern)
	return result
}

func (this *DatePatternConverter) Append(result []byte, event *LogEvent) []byte {
	value := event.Record.Time

	if this.location != nil {
		value = value.In(this.location)
	}

	start := len(result)
	result = value.AppendFormat(result, this.layout)

	return this.Format(result, start)
}

func (this *DatePatternConverter) toTimeLayout(pattern string) string {
	switch pattern {
	case "", "DEFAULT":
		return "2006-01-02 15:04:05.000"
	case "ISO8601":
		return "2006-01-02T15:04:05.000"
	case "ISO8601_BASIC":
		return "20060102T150405.000"
	case "ABSOLUTE":
		return "15:04:05.000"
	case "DATE":
		return "02 Jan 2006 15:04:05.000"
	default:
		return this.convertPattern(pattern)
	}
}

func (this *DatePatternConverter) convertPattern(pattern string) string {
	replacements := []struct {
		pattern string
		layout  string
	}{
		{"yyyy", "2006"},
		{"SSS", "000"},
		{"MM", "01"},
		{"dd", "02"},
		{"HH", "15"},
		{"mm", "04"},
		{"ss", "05"},
	}

	result := pattern

	for _, replacement := range replacements {
		result = strings.ReplaceAll(result, replacement.pattern, replacement.layout)
	}

	return result
}
