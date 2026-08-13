package impl

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/go-log4g/core/impl/filter"
	"github.com/go-log4g/core/impl/model"
	"github.com/go-log4g/core/impl/substitution"
)

type ConfigurationBuilder struct {
	statusLogger *StatusLogger
}

func NewConfigurationBuilder(statusLogger *StatusLogger) *ConfigurationBuilder {
	return &ConfigurationBuilder{
		statusLogger: statusLogger,
	}
}

func (this *ConfigurationBuilder) Build(definition *model.ConfigurationDefinition) *Configuration {
	if definition == nil {
		return NewDefaultConfiguration()
	}

	configuration := NewConfiguration()
	configuration.Root = NewLoggerConfig("", this.parseLevel(definition.Root.Level), definition.Root.Appenders...)
	for name, loggerDefinition := range definition.Loggers {
		loggerConfig := NewLoggerConfig(name, this.parseLevel(loggerDefinition.Level), loggerDefinition.Appenders...)
		if loggerDefinition.Additive != nil {
			loggerConfig.Additive = *loggerDefinition.Additive
		}
		configuration.Loggers = append(configuration.Loggers, loggerConfig)
	}

	sort.Slice(configuration.Loggers, func(i, j int) bool {
		left := configuration.Loggers[i]
		right := configuration.Loggers[j]
		if len(left.Name) != len(right.Name) {
			return len(left.Name) > len(right.Name)
		}
		return left.Name < right.Name
	})

	substitutor := substitution.NewSubstitutor(definition.Properties)
	patternParser := NewPatternParser(this.statusLogger)
	for name, appenderDefinition := range definition.Appenders {
		appender := this.buildAppender(name, appenderDefinition, patternParser, substitutor)
		if appender != nil {
			configuration.Appenders[name] = appender
		}
	}

	configuration.MinimumLevel = this.minimumLevel(configuration)
	return configuration
}

func (this *ConfigurationBuilder) parseLevel(value string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO", "":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		panic(fmt.Errorf("unsupported log level %q", value))
	}
}

func (this *ConfigurationBuilder) buildAppender(name string, definition model.AppenderDefinition, parser *PatternParser, substitutor *substitution.Substitutor) Appender {
	layout := this.buildLayout(definition.Layout, parser, substitutor)
	appenderFilter := this.buildFilter(definition.Filter)

	switch strings.ToLower(definition.Type) {
	case "console":
		switch strings.ToLower(definition.Target) {
		case "", "stdout":
			return NewConsoleAppender(os.Stdout, layout, appenderFilter)
		case "stderr":
			return NewConsoleAppender(os.Stderr, layout, appenderFilter)
		default:
			panic(fmt.Errorf("unsupported target %q for appender %q", definition.Target, name))
		}

	default:
		panic(fmt.Errorf("unsupported appender type %q for appender %q", definition.Type, name))
	}
}

func (this *ConfigurationBuilder) buildLayout(definition model.LayoutDefinition, parser *PatternParser, substitutor *substitution.Substitutor) Layout {
	switch strings.ToLower(definition.Type) {
	case "", "pattern":
		pattern := definition.Pattern
		if pattern == "" {
			pattern = defaultPattern
		} else {
			pattern = substitutor.Substitute(pattern)
		}
		return NewPatternLayout(pattern, parser)

	default:
		panic(fmt.Errorf("unsupported layout type %q", definition.Type))
	}
}

func (this *ConfigurationBuilder) buildFilter(definition model.FilterDefinition) filter.Filter {
	switch strings.ToLower(strings.TrimSpace(definition.Type)) {
	case "":
		return nil

	case "threshold":
		return filter.NewThresholdFilter(
			this.parseRequiredLevel(definition.Level, "filter.level"),
			filter.Neutral,
			filter.Deny,
		)

	case "levelrange":
		minLevel := this.parseRequiredLevel(definition.MinLevel, "filter.minLevel")
		maxLevel := this.parseRequiredLevel(definition.MaxLevel, "filter.maxLevel")

		if minLevel > maxLevel {
			panic(fmt.Errorf("invalid filter level range %q..%q", definition.MinLevel, definition.MaxLevel))
		}

		return filter.NewLevelRangeFilter(minLevel, maxLevel, filter.Neutral, filter.Deny)

	default:
		panic(fmt.Errorf("unsupported filter type %q", definition.Type))
	}
}

func (this *ConfigurationBuilder) parseRequiredLevel(value, property string) slog.Level {
	if strings.TrimSpace(value) == "" {
		panic(fmt.Errorf("missing required %s", property))
	}

	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		panic(fmt.Errorf("unsupported log level %q for %s", value, property))
	}
}
func (this *ConfigurationBuilder) minimumLevel(configuration *Configuration) slog.Level {
	result := configuration.Root.Level

	for _, loggerConfig := range configuration.Loggers {
		if loggerConfig.Level < result {
			result = loggerConfig.Level
		}
	}

	return result
}
