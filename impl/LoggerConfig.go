package impl

import "log/slog"

type LoggerConfig struct {
	Name         string
	Level        slog.Level
	AppenderRefs []*AppenderRef
	Additive     bool
}

func NewLoggerConfig(name string, level slog.Level, appenders ...string) *LoggerConfig {
	appenderRefs := make([]*AppenderRef, 0, len(appenders))
	for _, appender := range appenders {
		appenderRefs = append(appenderRefs, &AppenderRef{
			Ref: appender,
		})
	}

	return &LoggerConfig{
		Name:         name,
		Level:        level,
		AppenderRefs: appenderRefs,
		Additive:     true,
	}
}
