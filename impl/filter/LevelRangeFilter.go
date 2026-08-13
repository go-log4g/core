package filter

import "log/slog"

type LevelRangeFilter struct {
	AbstractFilter
	minLevel slog.Level
	maxLevel slog.Level
}

func NewLevelRangeFilter(minLevel, maxLevel slog.Level, onMatch, onMismatch Result) *LevelRangeFilter {
	return &LevelRangeFilter{
		AbstractFilter: NewAbstractFilter(onMatch, onMismatch),
		minLevel:       minLevel,
		maxLevel:       maxLevel,
	}
}

func (this *LevelRangeFilter) Filter(level slog.Level) Result {
	if level >= this.minLevel && level <= this.maxLevel {
		return this.OnMatch
	}
	return this.OnMismatch
}
