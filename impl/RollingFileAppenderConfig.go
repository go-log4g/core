package impl

import (
	"github.com/go-log4g/core/impl/filter"
	"github.com/go-log4g/core/impl/rolling"
)

type RollingFileAppenderConfig struct {
	File           string
	FilePattern    string
	Append         bool
	BufferSize     int
	ImmediateFlush bool
	Layout         Layout
	Filter         filter.Filter
	StatusLogger   *StatusLogger
	StartupPolicy  *rolling.OnStartupTriggeringPolicy
	TriggerPolicy  rolling.TriggeringPolicy
	Strategy       rolling.RolloverStrategy
}
