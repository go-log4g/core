package rolling

import "time"

type TriggeringContext struct {
	FileSize  int64
	EventSize int
	Time      time.Time
}
