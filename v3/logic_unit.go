package sameriver

import (
	"time"
)

type LogicUnit struct {
	name        string
	f           func(dt_ms float64)
	active      bool
	worldID     int
	lastRun     time.Time
	runSchedule *TimeAccumulator
}

func (l *LogicUnit) Activate() {
	l.active = true
}

func (l *LogicUnit) Deactivate() {
	l.active = false
	// zero lastRun
	l.lastRun = time.Time{}
}
