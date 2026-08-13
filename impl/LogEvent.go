package impl

import (
	"context"
	"log/slog"
)

type LogEvent struct {
	Context       context.Context
	Record        slog.Record
	CallerContext *CallerContext
	Attrs         []*HandlerAttr
	Groups        []string
}

func NewLogEvent(ctx context.Context, record slog.Record, callerContext *CallerContext, attrs []*HandlerAttr, groups []string) *LogEvent {
	return &LogEvent{
		Context:       ctx,
		Record:        record,
		CallerContext: callerContext,
		Attrs:         attrs,
		Groups:        groups,
	}
}
