package impl

import (
	"context"
	"log/slog"

	"github.com/go-errr/go/err"
)

type Handler struct {
	configuration  *Configuration
	callerContexts *CallerContextCache
	statusLogger   *StatusLogger
	attrs          []*HandlerAttr
	groups         []string
}

func NewHandler(configuration *Configuration, callerContexts *CallerContextCache, statusLogger *StatusLogger) *Handler {
	return &Handler{
		configuration:  configuration,
		callerContexts: callerContexts,
		statusLogger:   statusLogger,
	}
}

func (this *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= this.configuration.MinimumLevel
}

// Implements slog.Handler
func (this *Handler) Handle(ctx context.Context, record slog.Record) error {
	defer err.Recover(func(e any) {
		this.statusLogger.ErrorThrottled("handler", e, "Failed to process log event")
	})

	callerContext := this.callerContexts.Get(record.PC)
	if record.Level < callerContext.Config.Level {
		return nil
	}

	event := NewLogEvent(ctx, record, callerContext, this.attrs, this.groups)
	for _, appender := range callerContext.Config.Appenders {
		appender.Append(event)
	}

	return nil
}

func (this *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	result := *this
	result.attrs = append([]*HandlerAttr{}, this.attrs...)

	for _, attr := range attrs {
		result.attrs = append(result.attrs, NewHandlerAttr(attr, this.groups))
	}

	return &result
}

func (this *Handler) WithGroup(name string) slog.Handler {
	result := *this
	result.groups = append(append([]string{}, this.groups...), name)
	return &result
}
