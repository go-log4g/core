package impl

import (
	"context"
	"log/slog"
)

type Handler struct {
	configuration  *Configuration
	callerContexts *CallerContextCache
	attrs          []*HandlerAttr
	groups         []string
}

func NewHandler(configuration *Configuration, callerContexts *CallerContextCache) *Handler {
	return &Handler{
		configuration:  configuration,
		callerContexts: callerContexts,
	}
}

func (this *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= this.configuration.MinimumLevel
}

func (this *Handler) Handle(ctx context.Context, record slog.Record) error {
	callerContext := this.callerContexts.Get(record.PC)

	if record.Level < callerContext.Config.Level {
		return nil
	}

	event := NewLogEvent(ctx, record, callerContext, this.attrs, this.groups)
	for _, appender := range callerContext.Config.Appenders {
		if e := appender.Append(event); e != nil {
			return e
		}
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
