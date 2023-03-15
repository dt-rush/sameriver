package sameriver

type LogicUnit struct {
	// TODO: export name and active
	name        string
	f           func(dt_ms float64)
	active      bool
	worldID     int
	runSchedule *TimeAccumulator
}

func (l *LogicUnit) Activate() {
	l.active = true
}

func (l *LogicUnit) Deactivate() {
	l.active = false
}
