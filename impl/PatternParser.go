package impl

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-log4g/core/impl/abbr"
)

type PatternParser struct {
	statusLogger *StatusLogger
}

func NewPatternParser(statusLogger *StatusLogger) *PatternParser {
	return &PatternParser{
		statusLogger: statusLogger,
	}
}

func (this *PatternParser) Parse(pattern string) []PatternConverter {
	result := make([]PatternConverter, 0, 16)
	literalStart := 0

	for index := 0; index < len(pattern); {
		if pattern[index] != '%' {
			index++
			continue
		}

		if literalStart < index {
			result = append(result, NewLiteralPatternConverter(pattern[literalStart:index]))
		}

		element, nextIndex := this.parseElement(pattern, index)
		result = append(result, element)

		index = nextIndex
		literalStart = index
	}

	if literalStart < len(pattern) {
		result = append(result, NewLiteralPatternConverter(pattern[literalStart:]))
	}

	return result
}

func (this *PatternParser) parseElement(pattern string, index int) (PatternConverter, int) {
	start := index
	index++

	if index >= len(pattern) {
		this.statusLogger.Error("Invalid pattern %q: incomplete conversion at position %d", pattern, start)
		return NewLiteralPatternConverter(pattern[start:]), len(pattern)
	}

	if pattern[index] == '%' {
		return NewLiteralPatternConverter("%"), index + 1
	}

	formatting, index := this.parseFormatting(pattern, index)

	if index >= len(pattern) {
		this.statusLogger.Error("Invalid pattern %q: incomplete conversion at position %d", pattern, start)
		return NewLiteralPatternConverter(pattern[start:]), len(pattern)
	}

	key, nextIndex := this.parseConverterKey(pattern, index)
	switch key {
	case "m", "msg", "message":
		return NewMessagePatternConverter(formatting), nextIndex
	case "n":
		return NewLineSeparatorPatternConverter(formatting), nextIndex
	case "p", "level":
		return NewLevelPatternConverter(formatting), nextIndex
	case "c", "logger":
		return this.parseLogger(pattern, start, nextIndex, formatting)
	case "F", "file":
		return NewFilePatternConverter(formatting), nextIndex
	case "L", "line":
		return NewLinePatternConverter(formatting), nextIndex
	case "M", "method":
		return NewMethodPatternConverter(formatting), nextIndex
	case "d", "date":
		return this.parseDate(pattern, start, nextIndex, formatting)
	case "X":
		return this.parseMdc(pattern, index+1, formatting)
	default:
		this.statusLogger.Error("Invalid pattern %q: unsupported conversion %%%s at position %d", pattern, key, start)
		return NewLiteralPatternConverter(pattern[start:nextIndex]), nextIndex
	}
}

func (this *PatternParser) parseFormatting(pattern string, index int) (FormattingInfo, int) {
	formatting := NewDefaultFormattingInfo()

	if index < len(pattern) && pattern[index] == '-' {
		formatting.LeftAlign = true
		index++
	}

	if index < len(pattern) && pattern[index] == '0' {
		formatting.ZeroPad = true
		index++
	}

	if value, nextIndex, ok := this.parseInt(pattern, index); ok {
		formatting.MinLength = value
		index = nextIndex
	}

	if index < len(pattern) && pattern[index] == '.' {
		index++

		if index < len(pattern) && pattern[index] == '-' {
			formatting.LeftTruncate = false
			index++
		}

		if value, nextIndex, ok := this.parseInt(pattern, index); ok {
			formatting.MaxLength = value
			index = nextIndex
		}
	}

	return formatting, index
}

func (this *PatternParser) parseLogger(pattern string, start, index int, formatting FormattingInfo) (PatternConverter, int) {
	if index >= len(pattern) || pattern[index] != '{' {
		return NewLoggerPatternConverter(formatting, abbr.NewNOPAbbreviator()), index
	}

	option, nextIndex, ok := this.parseOption(pattern, index)

	if !ok {
		this.statusLogger.Error("Invalid pattern %q: unclosed logger option at position %d", pattern, start)
		return NewLiteralPatternConverter(pattern[start:]), len(pattern)
	}

	abbreviator, e := NewNameAbbreviator(option)

	if e != nil {
		this.statusLogger.Error("Invalid pattern %q: %s at position %d", pattern, e, start)
		return NewLoggerPatternConverter(formatting, abbr.NewNOPAbbreviator()), nextIndex
	}

	return NewLoggerPatternConverter(formatting, abbreviator), nextIndex
}

func (this *PatternParser) parseConverterKey(pattern string, index int) (string, int) {
	start := index

	for index < len(pattern) {
		c := pattern[index]

		if c < 'a' || c > 'z' {
			break
		}

		index++
	}

	if start == index {
		return pattern[start : start+1], start + 1
	}

	return pattern[start:index], index
}

func (this *PatternParser) parseDate(pattern string, start, index int, formatting FormattingInfo) (PatternConverter, int) {
	datePattern := ""
	var location *time.Location

	if index < len(pattern) && pattern[index] == '{' {
		value, nextIndex, ok := this.parseOption(pattern, index)

		if !ok {
			this.statusLogger.Error("Invalid pattern %q: unclosed date pattern at position %d", pattern, start)
			return NewLiteralPatternConverter(pattern[start:]), len(pattern)
		}

		datePattern = value
		index = nextIndex
	}

	if index < len(pattern) && pattern[index] == '{' {
		value, nextIndex, ok := this.parseOption(pattern, index)

		if !ok {
			this.statusLogger.Error("Invalid pattern %q: unclosed timezone at position %d", pattern, start)
			return NewLiteralPatternConverter(pattern[start:]), len(pattern)
		}

		resolved, e := time.LoadLocation(value)

		if e != nil {
			this.statusLogger.Error("Invalid pattern %q: invalid timezone %q at position %d", pattern, value, start)
		} else {
			location = resolved
		}

		index = nextIndex
	}

	return NewDatePatternConverter(datePattern, location, formatting), index
}

func (this *PatternParser) parseOption(pattern string, index int) (string, int, bool) {
	start := index + 1

	for index = start; index < len(pattern); index++ {
		if pattern[index] == '}' {
			return pattern[start:index], index + 1, true
		}
	}

	return "", len(pattern), false
}

func (this *PatternParser) parseInt(pattern string, index int) (int, int, bool) {
	start := index

	for index < len(pattern) && pattern[index] >= '0' && pattern[index] <= '9' {
		index++
	}

	if start == index {
		return 0, index, false
	}

	value, _ := strconv.Atoi(pattern[start:index])
	return value, index, true
}
func (this *PatternParser) parseMdc(pattern string, index int, formatting FormattingInfo) (PatternConverter, int) {
	if index >= len(pattern) || pattern[index] != '{' {
		return NewMdcPatternConverter(formatting, ""), index
	}

	end := strings.IndexByte(pattern[index+1:], '}')
	if end < 0 {
		return NewMdcPatternConverter(formatting, ""), index
	}

	end += index + 1
	key := pattern[index+1 : end]

	return NewMdcPatternConverter(formatting, key), end + 1
}
