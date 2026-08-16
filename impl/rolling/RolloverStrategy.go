package rolling

import "time"

type RolloverStrategy interface {
	Rollover(file, filePattern string, at time.Time)
}
