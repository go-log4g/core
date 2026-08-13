package filter

import "log/slog"

type ThresholdFilter struct {
	AbstractFilter
	level slog.Level
}

func NewThresholdFilter(level slog.Level, onMatch, onMismatch Result) *ThresholdFilter {
	return &ThresholdFilter{
		AbstractFilter: NewAbstractFilter(onMatch, onMismatch),
		level:          level,
	}
}

func (this *ThresholdFilter) Filter(level slog.Level) Result {
	if level >= this.level {
		return this.OnMatch
	}
	return this.OnMismatch
}
