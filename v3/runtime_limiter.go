package sameriver

import (
	"fmt"
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
	//
	// key: logicunit worldID
	// value: index
	indexes map[int]int
	// used to keep a running average of the entire runtime
	totalRuntime_ms *float64
	// overrun flag gets set whenever we are exceeding the allowance_ms
	overrun bool
}

func NewRuntimeLimiter() *RuntimeLimiter {
	return &RuntimeLimiter{
		logicUnits:       make([]*LogicUnit, 0),
		runtimeEstimates: make(map[*LogicUnit]float64),
		indexes:          make(map[int]int),
	}
}

func (r *RuntimeLimiter) Run(allowance_ms float64) (remaining_ms float64) {
	r.startIX = r.runIX
	r.finished = false
	tStart := time.Now()
	r.overrun = false
	remaining_ms = allowance_ms
	if len(r.logicUnits) == 0 {
		r.finished = true
		return
	}
	for remaining_ms > 0 && len(r.logicUnits) > 0 {
		logic := r.logicUnits[r.runIX]
		if logic.lastRun.IsZero() {
			logic.lastRun = time.Now()
		}
		estimate, hasEstimate := r.runtimeEstimates[logic]
		var t0 time.Time
		var elapsed_ms float64
		// if we've already run at least one, quit early on estimated overrun
		if (r.runIX != r.startIX) && hasEstimate && (estimate > allowance_ms) {
			r.overrun = true
			return remaining_ms
		}
		// else, we're either at the first func, or we have no estimate for this
		// one, or the estimate is within allowance_ms. SO run it
		if !hasEstimate ||
			(hasEstimate && estimate <= allowance_ms) ||
			(hasEstimate && estimate > allowance_ms && r.runIX == r.startIX) {
			t0 = time.Now()
			dt_ms := float64(time.Since(logic.lastRun).Nanoseconds()) / 1.0e6
			if logic.active &&
				(logic.runSchedule == nil || logic.runSchedule.Tick(dt_ms)) {
				logic.f(dt_ms)
				logic.lastRun = time.Now()
			}
			elapsed_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
			// update estimate stat
			if !hasEstimate {
				r.runtimeEstimates[logic] = elapsed_ms
			} else {
				r.runtimeEstimates[logic] =
					(r.runtimeEstimates[logic] + elapsed_ms) / 2.0
			}
		}
		remaining_ms -= elapsed_ms
		r.runIX = (r.runIX + 1) % len(r.logicUnits)
		if r.runIX == r.startIX {
			r.finished = true
			break
		}
	}
	total_ms := float64(time.Since(tStart).Nanoseconds()) / 1.0e6
	// maintain moving average of totalRuntime_ms
	if r.totalRuntime_ms == nil {
		r.totalRuntime_ms = &total_ms
	} else {
		*r.totalRuntime_ms = (*r.totalRuntime_ms + total_ms) / 2.0
	}
	r.totalRuntime_ms = r.totalRuntime_ms
	// return overunder_ms
	overunder_ms := allowance_ms - total_ms
	if overunder_ms < 0 {
		r.overrun = true
	}
	return overunder_ms
}

func (r *RuntimeLimiter) Add(logic *LogicUnit) {
	// panic if adding duplicate by WorldID
	if _, ok := r.indexes[logic.worldID]; ok {
		panic(fmt.Sprintf("Double-add of same logic unit to RuntimeLimiter "+
			"(WorldID: %d)", logic.worldID))
	}
	r.logicUnits = append(r.logicUnits, logic)
	r.indexes[logic.worldID] = len(r.logicUnits) - 1
}

func (r *RuntimeLimiter) Remove(l *LogicUnit) bool {
	// return early if nil
	if l == nil {
		return false
	}
	// return early if not present
	index, ok := r.indexes[l.worldID]
	if !ok {
		return false
	}
	// delete from runtimeEstimates
	if _, ok := r.runtimeEstimates[l]; ok {
		delete(r.runtimeEstimates, l)
	}
	// delete from indexes
	delete(r.indexes, l.worldID)
	// delete from logicUnits by replacing the last element into its spot,
	// updating the indexes entry for that element
	lastIndex := len(r.logicUnits) - 1
	if len(r.logicUnits) > 1 && index != lastIndex {
		r.logicUnits[index] = r.logicUnits[lastIndex]
		// update the indexes array for the elemnt we put into the
		// place of the one we spliced out
		nowAtIndex := r.logicUnits[index]
		r.indexes[nowAtIndex.worldID] = index
	}
	r.logicUnits = r.logicUnits[:lastIndex]
	// update runIX - if we removed an entity earlier in the list,
	// we should subtract 1 to keep runIX at it's same position. If we
	// removed one later in the list or equal to the current position,
	// we do nothing
	if index < r.runIX {
		r.runIX--
	}
	// success!
	return true
}

func (r *RuntimeLimiter) ActivateAll() {
	for _, l := range r.logicUnits {
		l.active = true
	}
}

func (r *RuntimeLimiter) DeactivateAll() {
	for _, l := range r.logicUnits {
		l.active = false
	}
}

func (r *RuntimeLimiter) Finished() bool {
	return r.finished
}

func (r *RuntimeLimiter) DumpStats() (stats map[string]float64, total float64) {
	stats = make(map[string]float64)
	for _, l := range r.logicUnits {
		if est, ok := r.runtimeEstimates[l]; ok {
			stats[l.name] = est
		} else {
			stats[l.name] = 0.0
		}
	}
	if r.totalRuntime_ms != nil {
		total = *r.totalRuntime_ms
	}
	return
}
