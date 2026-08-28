package impl

import (
	"log/slog"
	"strings"

	"github.com/go-jang/go/lang"
)

type Configuration struct {
	Root             *LoggerConfig
	Loggers          []*LoggerConfig
	Appenders        map[string]Appender
	PendingAppenders map[string]*PendingAppender
	MinimumLevel     slog.Level
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Root:             NewLoggerConfig("", slog.LevelError),
		Appenders:        make(map[string]Appender),
		PendingAppenders: make(map[string]*PendingAppender),
		MinimumLevel:     slog.LevelError,
	}
}

func (this *Configuration) Resolve(logger string) *ResolvedLoggerConfig {
	loggerConfig := this.Root
	appenderRefs := make([]*AppenderRef, 0, len(this.Appenders))
	added := make(map[string]bool)

	for _, candidate := range this.Loggers {
		if !this.matchesLogger(logger, candidate.Name) {
			continue
		}
		if loggerConfig == this.Root {
			loggerConfig = candidate
		}

		for _, definition := range candidate.AppenderRefs {
			if added[definition.Ref] {
				continue
			}

			appenderRef := this.resolveAppenderRef(definition)
			if appenderRef == nil {
				continue
			}

			appenderRefs = append(appenderRefs, appenderRef)
			added[definition.Ref] = true
		}

		if !candidate.Additive {
			return NewResolvedLoggerConfig(loggerConfig.Level, appenderRefs)
		}
	}

	for _, definition := range this.Root.AppenderRefs {
		if added[definition.Ref] {
			continue
		}

		appenderRef := this.resolveAppenderRef(definition)
		if appenderRef == nil {
			continue
		}

		appenderRefs = append(appenderRefs, appenderRef)
		added[definition.Ref] = true
	}

	return NewResolvedLoggerConfig(loggerConfig.Level, appenderRefs)
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

func (this *Configuration) resolveAppenderRef(ref *AppenderRef) *AppenderRef {
	appender := this.Appenders[ref.Ref]
	if appender == nil {
		return nil
	}

	result := *ref
	result.Appender = appender
	return &result
}

func (this *Configuration) ActivatePendingAppenders(appenderType string, factory AppenderFactory) {
	for name, pending := range this.PendingAppenders {
		if pending.Type != appenderType {
			continue
		}

		appender := factory(name, pending)
		lang.Assert(appender != nil, "appender factory returned nil for %q", name)

		this.Appenders[name] = appender
		delete(this.PendingAppenders, name)
	}
}
