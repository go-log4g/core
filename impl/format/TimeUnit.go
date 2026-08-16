package format

import (
	"time"

	"github.com/go-jang/go/lang"
)

type Unit int

const (
	Year Unit = iota
	Month
	Day
	Hour
	Minute
	Second
)

func (this Unit) Add(value time.Time, interval int) time.Time {
	switch this {
	case Year:
		return value.AddDate(interval, 0, 0)
	case Month:
		return value.AddDate(0, interval, 0)
	case Day:
		return value.AddDate(0, 0, interval)
	case Hour:
		return value.Add(time.Duration(interval) * time.Hour)
	case Minute:
		return value.Add(time.Duration(interval) * time.Minute)
	case Second:
		return value.Add(time.Duration(interval) * time.Second)
	default:
		panic("unsupported time unit")
	}
}

func (this Unit) Truncate(value time.Time) time.Time {
	location := value.Location()

	switch this {
	case Year:
		return time.Date(value.Year(), 1, 1, 0, 0, 0, 0, location)
	case Month:
		return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, location)
	case Day:
		return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, location)
	case Hour:
		return time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), 0, 0, 0, location)
	case Minute:
		return time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), 0, 0, location)
	case Second:
		return time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), value.Second(), 0, location)
	default:
		panic("unsupported time unit")
	}
}

func (this Unit) Next(value time.Time, interval int, modulate bool) time.Time {
	lang.Assert(interval > 0, "interval must be positive")

	if modulate {
		value = this.Truncate(value)
	}

	return this.Add(value, interval)
}
