package impl

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-jang/go/lang"
	"github.com/go-log4g/core/impl/filter"
	"github.com/go-log4g/core/impl/format"
	"github.com/go-log4g/core/impl/model"
	"github.com/go-log4g/core/impl/rolling"
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

	substitutor.SetProperties(definition.Properties)
	configuration := NewConfiguration()
	configuration.Root = NewLoggerConfig("", this.rootLevel(definition, substitutor), definition.Root.Appenders...)
	for name, loggerDefinition := range definition.Loggers {
		loggerConfig := NewLoggerConfig(name, configuration.Root.Level, loggerDefinition.Appenders...)
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
	this.resolveLoggerLevels(configuration, definition.Loggers, substitutor)

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

func parseLevel(value string) slog.Level {
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
		panic(fmt.Errorf("unsupported log level %q", value))
	}
}

func (this *ConfigurationBuilder) buildAppender(name string, definition model.AppenderDefinition, parser *PatternParser, substitutor *substitution.Substitutor) Appender {
	appenderType := strings.ToLower(substitutor.Substitute(definition.Type))
	layout := this.buildLayout(definition.Layout, parser, substitutor)
	appenderFilter := this.buildFilter(definition.Filter, substitutor)

	switch appenderType {
	case "console":
		return this.buildConsoleAppender(name, definition, layout, appenderFilter, substitutor)
	case "file":
		return this.buildFileAppender(name, definition, layout, appenderFilter, substitutor)
	case "rollingfile":
		return this.buildRollingFileAppender(name, definition, layout, appenderFilter, substitutor)
	default:
		panic(fmt.Sprintf("unsupported appender type %q for appender %q", appenderType, name))
	}
}

func (this *ConfigurationBuilder) buildConsoleAppender(name string, definition model.AppenderDefinition, layout Layout, appenderFilter filter.Filter, substitutor *substitution.Substitutor) Appender {
	target := strings.ToLower(substitutor.Substitute(definition.Target))

	switch target {
	case "", "stdout":
		return NewConsoleAppender(os.Stdout, layout, appenderFilter, this.statusLogger)
	case "stderr":
		return NewConsoleAppender(os.Stderr, layout, appenderFilter, this.statusLogger)
	default:
		panic(fmt.Sprintf("unsupported target %q for appender %q", target, name))
	}
}

func (this *ConfigurationBuilder) buildFileAppender(name string, definition model.AppenderDefinition, layout Layout, appenderFilter filter.Filter, substitutor *substitution.Substitutor) Appender {
	file := substitutor.Substitute(definition.File)
	lang.Assert(strings.TrimSpace(file) != "", "file is required for appender %q", name)

	append := true
	if definition.Append != nil {
		append = *definition.Append
	}

	bufferSize := 8192
	if definition.BufferSize != nil {
		bufferSize = *definition.BufferSize
	}
	lang.Assert(bufferSize > 0, "bufferSize must be positive for appender %q", name)

	immediateFlush := true
	if definition.ImmediateFlush != nil {
		immediateFlush = *definition.ImmediateFlush
	}

	return NewFileAppender(file, append, bufferSize, immediateFlush, layout, appenderFilter, this.statusLogger)
}

func (this *ConfigurationBuilder) buildRollingFileAppender(name string, definition model.AppenderDefinition, layout Layout, appenderFilter filter.Filter, substitutor *substitution.Substitutor) Appender {
	file := substitutor.Substitute(definition.File)
	lang.Assert(strings.TrimSpace(file) != "", "file is required for appender %q", name)

	filePattern := substitutor.Substitute(definition.FilePattern)
	lang.Assert(strings.TrimSpace(filePattern) != "", "filePattern is required for appender %q", name)

	append := true
	if definition.Append != nil {
		append = *definition.Append
	}

	bufferSize := 8192
	if definition.BufferSize != nil {
		bufferSize = *definition.BufferSize
	}
	lang.Assert(bufferSize > 0, "bufferSize must be positive for appender %q", name)

	immediateFlush := true
	if definition.ImmediateFlush != nil {
		immediateFlush = *definition.ImmediateFlush
	}

	startupPolicy := this.buildStartupTriggeringPolicy(definition.Policies, substitutor)
	triggerPolicy := this.buildTriggeringPolicy(name, definition.Policies, filePattern, substitutor)
	strategy := this.buildRolloverStrategy(name, definition.DefaultRolloverStrategy, filePattern, substitutor)

	return NewRollingFileAppender(file, filePattern, append, bufferSize, immediateFlush, layout, appenderFilter, this.statusLogger, startupPolicy, triggerPolicy, strategy)
}

func (this *ConfigurationBuilder) buildStartupTriggeringPolicy(definition model.PoliciesDefinition, substitutor *substitution.Substitutor) *rolling.OnStartupTriggeringPolicy {
	if definition.OnStartupTriggeringPolicy == nil {
		return nil
	}

	minSize := substitutor.Substitute(definition.OnStartupTriggeringPolicy.MinSize)
	return rolling.NewOnStartupTriggeringPolicy(format.ParseFileSize(minSize, 1))
}

func (this *ConfigurationBuilder) buildTriggeringPolicy(name string, definition model.PoliciesDefinition, filePattern string, substitutor *substitution.Substitutor) rolling.TriggeringPolicy {
	policies := make([]rolling.TriggeringPolicy, 0, 2)

	if definition.TimeBasedTriggeringPolicy != nil {
		interval := 1
		if definition.TimeBasedTriggeringPolicy.Interval != nil {
			interval = *definition.TimeBasedTriggeringPolicy.Interval
		}
		lang.Assert(interval > 0, "timeBasedTriggeringPolicy.interval must be positive for appender %q", name)

		modulate := false
		if definition.TimeBasedTriggeringPolicy.Modulate != nil {
			modulate = *definition.TimeBasedTriggeringPolicy.Modulate
		}

		unit := rolling.TimeUnitFromFilePattern(filePattern)
		policies = append(policies, rolling.NewTimeBasedTriggeringPolicy(unit, interval, modulate))
	}

	if definition.SizeBasedTriggeringPolicy != nil {
		size := substitutor.Substitute(definition.SizeBasedTriggeringPolicy.Size)
		policies = append(policies, rolling.NewSizeBasedTriggeringPolicy(format.ParseFileSize(size, 10*1024*1024)))
	}

	if len(policies) == 0 {
		return nil
	} else if len(policies) == 1 {
		return policies[0]
	}
	return rolling.NewCompositeTriggeringPolicy(policies...)
}

func (this *ConfigurationBuilder) buildRolloverStrategy(name string, definition model.DefaultRolloverStrategyDefinition, filePattern string, substitutor *substitution.Substitutor) rolling.RolloverStrategy {
	max := 7
	if definition.Max != nil {
		max = *definition.Max
	}
	lang.Assert(max > 0, "defaultRolloverStrategy.max must be positive for appender %q", name)

	var deleteAction *rolling.DeleteAction
	if definition.Delete != nil {
		maxAge := time.Duration(0)
		if strings.TrimSpace(definition.Delete.MaxAge) != "" {
			maxAge = format.ParseDuration(substitutor.Substitute(definition.Delete.MaxAge))
		}

		maxFiles := 0
		if definition.Delete.MaxFiles != nil {
			maxFiles = *definition.Delete.MaxFiles
		}
		lang.Assert(maxFiles >= 0, "defaultRolloverStrategy.delete.maxFiles must not be negative for appender %q", name)

		maxTotalSize := int64(0)
		if strings.TrimSpace(definition.Delete.MaxTotalSize) != "" {
			maxTotalSize = format.ParseFileSize(substitutor.Substitute(definition.Delete.MaxTotalSize), 0)
		}

		deleteAction = rolling.NewDeleteAction(filePattern, maxAge, maxFiles, maxTotalSize)
	}

	return rolling.NewDefaultRolloverStrategy(max, deleteAction)
}

func (this *ConfigurationBuilder) buildLayout(definition model.LayoutDefinition, parser *PatternParser, substitutor *substitution.Substitutor) Layout {
	layoutType := strings.ToLower(substitutor.Substitute(definition.Type))
	switch layoutType {
	case "", "pattern":
		pattern := definition.Pattern
		if pattern == "" {
			pattern = defaultPattern
		} else {
			pattern = substitutor.Substitute(pattern)
		}
		return NewPatternLayout(pattern, parser)
	default:
		panic(fmt.Sprintf("unsupported layout type %q", layoutType))
	}
}

func (this *ConfigurationBuilder) buildFilter(definition model.FilterDefinition, substitutor *substitution.Substitutor) filter.Filter {
	filterType := strings.ToLower(substitutor.Substitute(definition.Type))
	switch filterType {
	case "":
		return nil
	case "threshold":
		level := this.parseRequiredLevel(substitutor.Substitute(definition.Level), "filter.level")
		return filter.NewThresholdFilter(level, filter.Neutral, filter.Deny)
	case "levelrange":
		minLevel := this.parseRequiredLevel(substitutor.Substitute(definition.MinLevel), "filter.minLevel")
		maxLevel := this.parseRequiredLevel(substitutor.Substitute(definition.MaxLevel), "filter.maxLevel")
		lang.Assert(minLevel <= maxLevel, "invalid filter level range %q..%q", definition.MinLevel, definition.MaxLevel)
		return filter.NewLevelRangeFilter(minLevel, maxLevel, filter.Neutral, filter.Deny)
	default:
		panic(fmt.Sprintf("unsupported filter type %q", filterType))
	}
}

func (this *ConfigurationBuilder) parseRequiredLevel(value, property string) slog.Level {
	lang.Assert(strings.TrimSpace(value) != "", "missing required %s", property)

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

func (this *ConfigurationBuilder) rootLevel(definition *model.ConfigurationDefinition, substitutor *substitution.Substitutor) slog.Level {
	if strings.TrimSpace(definition.Root.Level) != "" {
		return parseLevel(substitutor.Substitute(definition.Root.Level))
	} else if value := substitutor.RootLevel(); value != "" {
		return parseLevel(value)
	}
	return slog.LevelError
}

func (this *ConfigurationBuilder) resolveLoggerLevels(configuration *Configuration, definitions map[string]model.LoggerDefinition, substitutor *substitution.Substitutor) {
	for _, logger := range configuration.Loggers {
		definition := definitions[logger.Name]

		if strings.TrimSpace(definition.Level) != "" {
			logger.Level = parseLevel(substitutor.Substitute(definition.Level))
			continue
		}

		for _, parent := range configuration.Loggers {
			if parent.Name == logger.Name || !strings.HasPrefix(logger.Name, parent.Name+"/") {
				continue
			}
			if strings.TrimSpace(definitions[parent.Name].Level) != "" {
				logger.Level = parent.Level
				break
			}
		}
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
