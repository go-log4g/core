package impl

import (
	"log/slog"
	"strings"
)

type Configuration struct {
	Root         *LoggerConfig
	Loggers      []*LoggerConfig
	Appenders    map[string]Appender
	MinimumLevel slog.Level
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Root:         NewLoggerConfig("", slog.LevelError),
		Appenders:    make(map[string]Appender),
		MinimumLevel: slog.LevelError,
	}
}

func (this *Configuration) Resolve(logger string) *ResolvedLoggerConfig {
	loggerConfig := this.Root
	appenders := make([]Appender, 0, len(this.Appenders))
	added := make(map[string]bool)

	for _, candidate := range this.Loggers {
		if !this.matchesLogger(logger, candidate.Name) {
			continue
		}
		if loggerConfig == this.Root {
			loggerConfig = candidate
		}

		for _, name := range candidate.Appenders {
			if added[name] {
				continue
			}
			appender := this.Appenders[name]
			if appender == nil {
				continue
			}
			appenders = append(appenders, appender)
			added[name] = true
		}

		if !candidate.Additive {
			return NewResolvedLoggerConfig(loggerConfig.Level, appenders)
		}
	}

	for _, name := range this.Root.Appenders {
		if added[name] {
			continue
		}
		appender := this.Appenders[name]
		if appender == nil {
			continue
		}
		appenders = append(appenders, appender)
		added[name] = true
	}

	return NewResolvedLoggerConfig(loggerConfig.Level, appenders)
}

func (this *Configuration) matchesLogger(logger, configured string) bool {
	if logger == configured {
		return true
	}

	if len(logger) <= len(configured) {
		return false
	}

	return strings.HasPrefix(logger, configured) && logger[len(configured)] == '/'
}
