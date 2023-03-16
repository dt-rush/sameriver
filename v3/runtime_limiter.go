package sameriver

import (
	"fmt"
	"time"
)

// used to store a set of logicUnits which want to be executed together and
// frequently, but which are tolerant of being partitioned in time in order to
// stay within a certain time constraint (for example, running all the world
// logic we can within 4 milliseconds, hopefully looping back around to wherever
// we started, but if not, picking up where we left off next Run()
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
	// used to provide an accurate dt_ms to each logic unit so it can integrate
	// time smooth and proper
	lastRun map[*LogicUnit]time.Time
	// we run a logic unit with a gap of at least x ms where it takes x ms
	// to run. so a function taking 4ms will have at least 4 ms after it finishes
	// til the next time it runs, so we need to keep track of when logicunits end.
	lastEnd map[*LogicUnit]time.Time
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
	// coefficient 0.0 to 1.0, percentage of runners in the last Run() cycle
	// from startIx back around to itself that did not get to run. Used
	// to allot extra allowance time left over once all runners have run once
	// to let starving runners try to use the leftover time to run something
	// for example, say 12 ms are available. 2 ms are given to each of 6 runners.
	// the first 5 complete to 100% and don't even use their full budget, but the
	// 6th runner starves at 20% complete using its mere 2 ms. then, the remaining
	// time left over, say 7 ms, is portioned entirely to the 6th runner to try to
	// complete. If there were 2 of 6 that starved at 20% each, they would divide
	// the remaining ms in half. If one starved at 10% and the other at 30%, then
	// the 10% one would get (10 / (10 + 30))th of the time, and the other would
	// get (30 / (10 + 30))th. The division of the spoils proceeds like this, with
	// leftover time alloted proportional to starvation in this way, until the
	// total starve of those that ran is zero.
	starvation float64
}

func NewRuntimeLimiter() *RuntimeLimiter {
	return &RuntimeLimiter{
		logicUnits:       make([]*LogicUnit, 0),
		runtimeEstimates: make(map[*LogicUnit]float64),
		lastRun:          make(map[*LogicUnit]time.Time),
		lastEnd:          make(map[*LogicUnit]time.Time),
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
	ran := 0
	for remaining_ms > 0 {
		logic := r.logicUnits[r.runIX]
		logRuntimeLimiter("--- %s", logic.name)

		// check whether this logic has ever run
		_, hasRunBefore := r.lastRun[logic]
		// check its estimate
		estimate, hasEstimate := r.runtimeEstimates[logic]

		// estimate looks good if it's below allowance OR the estimate is above
		// allowance but we left off at this index last time; so we should get the
		// painful function over with rather than stall here forever or wait
		// to execute it when we get enough allowance (may never happen)
		estimateLooksGood := (hasEstimate && estimate <= remaining_ms) ||
			(hasEstimate && estimate > allowance_ms && r.runIX == r.startIX)
		logRuntimeLimiter("estimateLooksGood: %t", estimateLooksGood)
		if hasEstimate && !estimateLooksGood {
			// NOTE: exiting early when we hit our first bad estimate may be
			// suboptimal in the sense that it leaves extra unused time that
			// smaller logic units ahead might have used, but trying to
			// implement skipping of heavy logics that also doesn't leave them
			// behind, leading to them getting starved is very hard to manage.
			// when we exit early, our starvation is > 0.0 so we will receive
			// a proportional share of the total leftover time in the next
			// loop inside Share()
			return remaining_ms
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

		if DEBUG_RUNTIME_LIMITER {
			logRuntimeLimiter("hasRunBefore: %t", hasRunBefore)
			logRuntimeLimiter("isActive: %t", isActive)
			logRuntimeLimiter("durationHasElapsed: %t", durationHasElapsed)
			logRuntimeLimiter("scheduled: %t", scheduled)
		}
		if !hasRunBefore ||
			(isActive && durationHasElapsed && scheduled) {

			t0 := time.Now()
			// note that we start lastrun from the moment the function starts, since
			// say it starts at t=0ms, takes 4 ms to run, then if it comes up to run
			// again at t=8ms (r.tick()), it will get dt_ms of 8 ms, the proper
			// intervening time since it last integrated a dt_ms increment.
			r.lastRun[logic] = time.Now()
			logic.f(dt_ms)
			r.lastEnd[logic] = time.Now()
			elapsed_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
			r.updateEstimate(logic, elapsed_ms)
			ran++
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
	r.starvation = float64(len(r.logicUnits)-ran) / float64(len(r.logicUnits))
	return overunder_ms
}

func (r *RuntimeLimiter) tick(logic *LogicUnit) bool {
	if t, ok := r.lastEnd[logic]; ok {
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
