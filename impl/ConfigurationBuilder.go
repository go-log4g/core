package impl

import (
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/go-log4g/core/impl/model"
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

	patternParser := NewPatternParser(this.statusLogger)
	for name, appenderDefinition := range definition.Appenders {
		appender := this.buildAppender(name, appenderDefinition, patternParser)
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
		this.statusLogger.Error("Unsupported log level %q, using INFO", value)
		return slog.LevelInfo
	}
}

func (this *ConfigurationBuilder) buildAppender(name string, definition model.AppenderDefinition, parser *PatternParser) Appender {
	switch strings.ToLower(definition.Type) {
	case "console":
		layout := this.buildLayout(definition.Layout, parser)

		switch strings.ToLower(definition.Target) {
		case "", "stdout":
			return NewConsoleAppender(os.Stdout, layout)
		case "stderr":
			return NewConsoleAppender(os.Stderr, layout)
		default:
			this.statusLogger.Error("Unsupported target %q for appender %q", definition.Target, name)
			return nil
		}

	default:
		this.statusLogger.Error("Unsupported appender type %q for appender %q", definition.Type, name)
		return nil
	}
}

func (this *ConfigurationBuilder) buildLayout(definition model.LayoutDefinition, parser *PatternParser) Layout {
	switch strings.ToLower(definition.Type) {
	case "", "pattern":
		pattern := definition.Pattern
		if pattern == "" {
			pattern = defaultPattern
		}
		return NewPatternLayout(pattern, parser)

	default:
		this.statusLogger.Error("Unsupported layout type %q, using pattern layout", definition.Type)
		return NewPatternLayout(defaultPattern, parser)
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
