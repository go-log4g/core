package impl

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/go-log4g/core/impl/abbr"
	"github.com/go-log4g/core/impl/substitution"
)

const defaultPattern = "%m%n"

var substitutor = substitution.NewSubstitutor()

func NewDefaultConfiguration() *Configuration {
	statusLogger := NewStatusLogger()
	patternParser := NewPatternParser(statusLogger)
	layout := NewPatternLayout(defaultPattern, patternParser)
	console := NewConsoleAppender(os.Stdout, layout, nil, statusLogger)

	rootLevel := slog.LevelError
	if value := substitutor.RootLevel(); value != "" {
		rootLevel = parseLevel(value)
	}

	configuration := NewConfiguration()
	configuration.Root = NewLoggerConfig("", rootLevel, "console")
	configuration.MinimumLevel = rootLevel
	configuration.Appenders["console"] = console

	return configuration
}

func NewNameAbbreviator(pattern string) (abbr.NameAbbreviator, error) {
	if pattern == "" {
		return abbr.NewNOPAbbreviator(), nil
	}

	precision, e := strconv.Atoi(pattern)
	if e == nil {
		return NewMaxElementAbbreviator(precision), nil
	}

	if abbreviator := newDynamicWordAbbreviator(pattern); abbreviator != nil {
		return abbreviator, nil
	}

	pattern = strings.TrimSuffix(pattern, ".")
	parts := strings.Split(pattern, ".")
	fragments := make([]*abbr.PatternAbbreviatorFragment, 0, len(parts))

	for _, part := range parts {
		if part == "*" {
			fragments = append(fragments, abbr.NewPatternAbbreviatorFragment(-1))
			continue
		}

		length, e := strconv.Atoi(part)
		if e != nil || length < 0 {
			return nil, fmt.Errorf("unsupported abbr pattern %q", pattern)
		}

		fragments = append(fragments, abbr.NewPatternAbbreviatorFragment(length))
	}

	return abbr.NewPatternAbbreviator(fragments...), nil
}

func newDynamicWordAbbreviator(pattern string) *abbr.DynamicWordAbbreviator {
	parts := strings.Split(pattern, ".")
	if len(parts) != 3 || parts[2] != "*" {
		return nil
	}

	charCount, e := strconv.Atoi(parts[0])
	if e != nil || charCount < 0 {
		return nil
	}

	retainCount, e := strconv.Atoi(parts[1])
	if e != nil || retainCount < 1 {
		return nil
	}

	return abbr.NewDynamicWordAbbreviator(charCount, retainCount)
}
