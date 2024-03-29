package sameriver

// used to add/remove logics from either runtimelimiters or runtimelimitsharers
type AddRemoveLogicEvent struct {
	addRemove bool
	l         *LogicUnit
}

type LogicUnit struct {
	// TODO: export name and active
	name        string
	f           func(dt_ms float64)
	active      bool
	worldID     int
	runSchedule *TimeAccumulator
	// hotness increments every time this logic is run
	// note: this doesn't overflow since it gets normalised
	hotness int
	// flag set to true at the start of each Run() call in RuntimeLimiter
	// if the schedule has ticked, if active, etc.
	shouldRun bool
	// set when this logic unit is executed
	ran bool
}

func (l *LogicUnit) Activate() {
	l.active = true
}

func (l *LogicUnit) Deactivate() {
	l.active = false
}
