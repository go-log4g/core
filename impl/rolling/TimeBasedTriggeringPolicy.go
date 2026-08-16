package rolling

import (
	"time"

	"github.com/go-jang/go/lang"
	"github.com/go-log4g/core/impl/format"
)

type TimeBasedTriggeringPolicy struct {
	unit     format.Unit
	interval int
	modulate bool
	next     time.Time
}

func NewTimeBasedTriggeringPolicy(unit format.Unit, interval int, modulate bool) *TimeBasedTriggeringPolicy {
	lang.Assert(interval > 0, "interval must be positive")

	return &TimeBasedTriggeringPolicy{
		unit:     unit,
		interval: interval,
		modulate: modulate,
	}
}

// Implements TriggeringPolicy
func (this *TimeBasedTriggeringPolicy) IsTriggered(context TriggeringContext) bool {
	if this.next.IsZero() {
		this.next = this.unit.Next(context.Time, this.interval, this.modulate)
		return false
	}

	if context.Time.Before(this.next) {
		return false
	}

	this.next = this.unit.Next(context.Time, this.interval, this.modulate)
	return true
}
