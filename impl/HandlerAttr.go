package impl

import "log/slog"

type HandlerAttr struct {
	Attr   slog.Attr
	Groups []string
}

func NewHandlerAttr(attr slog.Attr, groups []string) *HandlerAttr {
	return &HandlerAttr{
		Attr:   attr,
		Groups: groups,
	}
}
