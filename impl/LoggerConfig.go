package impl

import "log/slog"

type LoggerConfig struct {
	Name      string
	Level     slog.Level
	Appenders []string
	Additive  bool
}

func NewLoggerConfig(name string, level slog.Level, appenders ...string) *LoggerConfig {
	return &LoggerConfig{
		Name:      name,
		Level:     level,
		Appenders: appenders,
		Additive:  true,
	}
}
