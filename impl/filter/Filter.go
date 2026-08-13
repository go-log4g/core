package filter

import "log/slog"

type Filter interface {
	Filter(level slog.Level) Result
}
