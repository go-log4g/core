package impl

import "log/slog"

type ResolvedLoggerConfig struct {
	Level        slog.Level
	AppenderRefs []*AppenderRef
}

func NewResolvedLoggerConfig(level slog.Level, appenderRefs []*AppenderRef) *ResolvedLoggerConfig {
	return &ResolvedLoggerConfig{
		Level:        level,
		AppenderRefs: appenderRefs,
	}
}
