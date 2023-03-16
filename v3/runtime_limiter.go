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
	// keep track of whether we have skipped this logic unit once already due to
	// a bad estimate - if we did, we will override the consideration of the estimate
	skippedBadEstimate map[*LogicUnit]bool
	// we run a logic unit at most every x ms where x is the runtime estimate
	// so, we need to keep track of when it last ran
	lastRun map[*LogicUnit]time.Time
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
	// starved is a counter that gets updated each Run() to keep track of
	// how many logic units didn't get to run ( > 0 if we didn't iterate the full
	// list in the given allowance_ms)
	// TODO: take some kind of action on this in RuntimeLimitSharer
	starved int
}

func NewRuntimeLimiter() *RuntimeLimiter {
	return &RuntimeLimiter{
		logicUnits:         make([]*LogicUnit, 0),
		runtimeEstimates:   make(map[*LogicUnit]float64),
		skippedBadEstimate: make(map[*LogicUnit]bool),
		lastRun:            make(map[*LogicUnit]time.Time),
		indexes:            make(map[int]int),
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
	ran := 0
	for remaining_ms > 0 {
		logic := r.logicUnits[r.runIX]

		// check whether this logic has ever run
		_, hasRunBefore := r.lastRun[logic]
		// check its estimate
		estimate, hasEstimate := r.runtimeEstimates[logic]

		// if we've already run at least one, quit early on estimated overrun
		if (r.runIX != r.startIX) && hasRunBefore && (estimate > allowance_ms) {
			r.overrun = true
			return remaining_ms
		}

		// estimate looks good if it's below allowance OR the estimate is above
		// allowance but we left off at this index last time; so we should get the
		// painful function over with rather than stall here forever or wait
		// to execute it when we get enough allowance (may never happen)
		estimateLooksGood := (hasEstimate && estimate <= allowance_ms) ||
			(hasEstimate && estimate > allowance_ms && r.runIX == r.startIX)
		// override bad estimate if this logic unit was skipped for a bad estimate
		// once already
		if skipped, ok := r.skippedBadEstimate[logic]; ok && skipped {
			estimateLooksGood = true
			r.skippedBadEstimate[logic] = false
		}

		// if the time since the last run of this logic is > the runtime estimate
		// (that is, a function taking 1ms to run on avg should run at most
		// every 1ms)
		durationHasElapsed := r.tick(logic)

		// obviously the logic must be active
		isActive := logic.active

		// get real time since last run
		var dt_ms float64
		if hasRunBefore {
			dt_ms = float64(time.Since(r.lastRun[logic]).Nanoseconds()) / 1.0e6
		} else {
			dt_ms = 0
		}

		// finally, if it has a runschedule defined, we should also tick that
		// amount of time
		scheduled := logic.runSchedule == nil || logic.runSchedule.Tick(dt_ms)

		var elapsed_ms float64
		if !hasRunBefore ||
			(isActive && estimateLooksGood && durationHasElapsed && scheduled) {

			t0 := time.Now()
			// note that we start lastrun from the moment the function starts, since
			// say it starts at t=0ms, takes 4 ms to run, then if it comes up to run
			// again at t=8ms (r.tick()), it will be 8ms dt_ms since the last time
			// it ran
			r.lastRun[logic] = time.Now()
			logic.f(dt_ms)
			elapsed_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
			r.updateEstimate(logic, elapsed_ms)
			ran++
		} else if !estimateLooksGood {
			// if we didn't run because the estimate didn't look good, set a flag
			// so that next time we reach this, even if the estimate looks bad,
			// we run it anyway
			r.skippedBadEstimate[logic] = true
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
	// calculate overunder
	overunder_ms := allowance_ms - total_ms
	if overunder_ms < 0 {
		r.overrun = true
	}
	// calculate starved
	r.starved = len(r.logicUnits) - ran
	return overunder_ms
}

func (r *RuntimeLimiter) tick(logic *LogicUnit) bool {
	if t, ok := r.lastRun[logic]; ok {
		return float64(time.Since(t).Nanoseconds())/1.0e6 > r.runtimeEstimates[logic]
	} else {
		return true
	}
}

func (r *RuntimeLimiter) updateEstimate(logic *LogicUnit, elapsed_ms float64) {
	if _, ok := r.runtimeEstimates[logic]; !ok {
		r.runtimeEstimates[logic] = elapsed_ms
	} else {
		r.runtimeEstimates[logic] =
			(r.runtimeEstimates[logic] + elapsed_ms) / 2.0
	}
}

func (r *RuntimeLimiter) Add(logic *LogicUnit) {
	// panic if adding duplicate by WorldID
	if _, ok := r.indexes[logic.worldID]; ok {
		panic(fmt.Sprintf("Double-add of same logic unit to RuntimeLimiter "+
			"(WorldID: %d; name: %s)", logic.worldID, logic.name))
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
	delete(r.runtimeEstimates, l)
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
