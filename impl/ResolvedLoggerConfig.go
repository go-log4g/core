package impl

import "log/slog"

type ResolvedLoggerConfig struct {
	Level     slog.Level
	Appenders []Appender
}

func NewResolvedLoggerConfig(level slog.Level, appenders []Appender) *ResolvedLoggerConfig {
	return &ResolvedLoggerConfig{
		Level:     level,
		Appenders: appenders,
	}
}
