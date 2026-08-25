package impl

import (
	"log/slog"

	"github.com/go-log4g/core/impl/filter"
)

type AppenderRef struct {
	Ref      string
	Appender Appender
	Level    *slog.Level
	Filter   filter.Filter
}

func (this *AppenderRef) Accept(event *LogEvent) bool {
	if this.Level != nil && event.Record.Level < *this.Level {
		return false
	}

	return this.Filter == nil || this.Filter.Filter(event.Record.Level) != filter.Deny
}
