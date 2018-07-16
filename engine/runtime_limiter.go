package engine

import (
	"time"
)

// used to store a set of logicUnits which want to be executed together and
// frequently, but which are tolerant of being partitioned in time in order to
// stay within a certain time constraint (for example, running all the world
// logic we can within 4 milliseocnds, picking up where we left off next
// Update() loop)
type RuntimeLimiter struct {
	// used to degrade gracefully under time pressure, by picking up where we
	// left off in the iteration of logicUnits to run in the case that we can't
	// get to all of them within the milliseconds allotted
	startIX  int
	runIX    int
	finished bool
	// used so we can iterate the added logicUnits in order
	logicUnits []*LogicUnit
	// used to estimate the time cost in milliseconds of running a function,
	// so that we can try to stay below the limit provided
	runtimeEstimates map[*LogicUnit]float64
	// used to lookup the logicUnits slice index given an object to which
	// the LogicUnit is coupled, it's Parent (for System.Update() instances,
	// this is the System, for world LogicUnits this is the LogicUnit itself
	// This is needed to support efficient delete and activate/deactivate
	indexes map[int]int
}

func NewRuntimeLimiter() *RuntimeLimiter {
	return &RuntimeLimiter{
		logicUnits:       make([]*LogicUnit, 0),
		runtimeEstimates: make(map[*LogicUnit]float64),
		indexes:          make(map[int]int),
	}
}

func (r *RuntimeLimiter) Start() {
	r.startIX = r.runIX
	r.finished = false
}

func (r *RuntimeLimiter) Run(allowance float64) (remaining_ms float64) {
	remaining_ms = allowance
	if len(r.logicUnits) == 0 {
		r.finished = true
		return
	}
	for allowance > 0 && len(r.logicUnits) > 0 {
		logic := r.logicUnits[r.runIX]
		estimate, hasEstimate := r.runtimeEstimates[logic]
		var t0 time.Time
		var elapsed_ms float64
		if hasEstimate && (estimate > allowance) && (r.runIX != r.startIX) {
			return remaining_ms
		}
		if !hasEstimate ||
			(hasEstimate && estimate <= allowance) ||
			(hasEstimate && estimate > allowance && r.runIX == r.startIX) {
			t0 = time.Now()
			if logic.Active {
				logic.F()
			}
			elapsed_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
			if !hasEstimate {
				r.runtimeEstimates[logic] = elapsed_ms
			} else {
				r.runtimeEstimates[logic] =
					(r.runtimeEstimates[logic] + elapsed_ms) / 2.0
			}
		}
		allowance -= elapsed_ms
		r.runIX = (r.runIX + 1) % len(r.logicUnits)
		if r.runIX == r.startIX {
			r.finished = true
			break
		}
	}
	return allowance
}

func (r *RuntimeLimiter) Add(logic *LogicUnit) {
	// panic if adding duplicate by WorldID
	if _, ok := r.indexes[logic.WorldID]; ok {
		panic("Double-add of same logic unit to RuntimeLimiter")
	}
	r.logicUnits = append(r.logicUnits, logic)
	r.indexes[logic.WorldID] = len(r.logicUnits) - 1
}

func (r *RuntimeLimiter) Remove(WorldID int) bool {
	// return early if not present
	index, ok := r.indexes[WorldID]
	if !ok {
		return false
	}
	// delete from runtimeEstimates
	logicUnit := r.logicUnits[index]
	if _, ok := r.runtimeEstimates[logicUnit]; ok {
		delete(r.runtimeEstimates, logicUnit)
	}
	// delete from indexes
	delete(r.indexes, WorldID)
	// delete from logicUnits by replacing the last element into its spot,
	// updating the indexes entry for that element
	lastIndex := len(r.logicUnits) - 1
	if len(r.logicUnits) > 1 {
		r.logicUnits[index] = r.logicUnits[lastIndex]
		r.indexes[r.logicUnits[index].WorldID] = index
	}
	r.logicUnits = r.logicUnits[:lastIndex]
	return true
}

func (r *RuntimeLimiter) ActivateAll() {
	for _, l := range r.logicUnits {
		l.Active = true
	}
}

func (r *RuntimeLimiter) DeactivateAll() {
	for _, l := range r.logicUnits {
		l.Active = false
	}
}

func (r *RuntimeLimiter) Finished() bool {
	return r.finished
}

func RuntimeLimitShare(
	allowance float64, runners ...*RuntimeLimiter) (remaining_ms float64) {

	remaining_ms = allowance
	for _, r := range runners {
		r.Start()
	}
	finished := 0
	for allowance >= 0 && finished != len(runners) {
		perRunner := allowance / float64(len(runners)-finished)
		var remaining_ms float64
		for _, r := range runners {
			remaining_ms += r.Run(perRunner)
			if r.Finished() {
				finished++
			}
		}
		allowance = remaining_ms
	}
	return allowance
}
