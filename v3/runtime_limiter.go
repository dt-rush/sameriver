package sameriver

import (
	"fmt"
	"sort"
	"time"
)

// used to store a set of logicUnits which want to be executed together and
// frequently, but which are tolerant of being partitioned in time in order to
// stay within a certain time constraint (for example, running all the world
// logic we can within 4 milliseconds, hopefully looping back around to wherever
// we started, but if not, picking up where we left off next Run()
//
// it traverses logics round-robin until it reaches the first one that can't fit
// in the time remaining, then switches into opportunistic mode: it uses
// logic.hotness to sort them in ascending frequency of invocation and execute
// any as it iterates that sorted list which it can, excluding those
// which can't run in the remaining_ms when we get to them
type RuntimeLimiter struct {
	// used to degrade gracefully under time pressure, by picking up where we
	// left off in the iteration of logicUnits to run in the case that we can't
	// get to all of them within the milliseconds allotted; used for round-robin
	// iteration
	startIx  int
	runIx    int
	finished bool // whether we finished the round-robin, back to startIx
	// used for opportunistic time-fill of remaining once round-robin reaches
	// the first func too heavy to run
	oppIx int
	// used so we can iterate the added logicUnits in order
	logicUnits []*LogicUnit
	// logicUnits sorted by hotness ascending, which is an int incremented every time
	// the func gets run. this is used when, in round-robin scheduling according to
	// runIx, we reach the first unit that can't run in the budget. then we look
	// opportunistically to run any funcs which can in the time remaining sorted
	// ascending by hotness (hence we try to maintain a uniform distribution of
	// which funcs get called in opportunistic mode)
	ascendingHotness []*LogicUnit
	// used to estimate the time cost in milliseconds of running a function,
	// so that we can try to stay below the allowance_ms given to Run()
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
	// ranRobin : number that ran by round robin since last time bonsuTime = false
	// (that is, when loop = 0 in Share())
	ranRobin int
	// ranOpp : number that ran by opportunistic since last time bonsuTime = false
	ranOpp int
	// starvation : coefficient 0.0 to 1.0, percentage of runners in the last Run() cycle
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
	// keep track of the worst case loop overhead (time to iterate a logic unit
	// minus time it takes to execute)
	// after being set the first time, tracks a moving avg
	loopOverhead_ms float64
}

func NewRuntimeLimiter() *RuntimeLimiter {
	return &RuntimeLimiter{
		logicUnits:       make([]*LogicUnit, 0),
		ascendingHotness: make([]*LogicUnit, 0),
		runtimeEstimates: make(map[*LogicUnit]float64),
		lastRun:          make(map[*LogicUnit]time.Time),
		lastEnd:          make(map[*LogicUnit]time.Time),
		indexes:          make(map[int]int),
	}
}

func (r *RuntimeLimiter) Run(allowance_ms float64, bonsuTime bool) (remaining_ms float64) {
	if !bonsuTime {
		// if we're at loop = 0, the initial share from the RuntimeLimitSharer
		// this runner is in, we set startIx
		r.startIx = r.runIx
		// else, startix remains where it was on last time loop = 0,
		r.finished = false
	}
	tStart := time.Now()
	r.overrun = false
	remaining_ms = allowance_ms
	if len(r.logicUnits) == 0 {
		r.finished = true
		return
	}
	if !bonsuTime {
		r.ranRobin = 0
		r.ranOpp = 0
		r.starvation = 0
	}
	const (
		RoundRobin int = iota
		Opportunistic
	)
	mode := RoundRobin
	worstOverheadThisTime := 0.0
	logRuntimeLimiter("Run(); allowance: %f ms", allowance_ms)
	for remaining_ms > 0 {
		if remaining_ms < 3*r.loopOverhead_ms {
			logRuntimeLimiter("XXX RUN() OVERHEAD BAIL XXX")
			break
		}
		tLoop := time.Now()
		if DEBUG_RUNTIME_LIMITER {
			modeStr := ""
			if bonsuTime {
				modeStr += "BonsuTime "
			}
			switch mode {
			case RoundRobin:
				modeStr += "RoundRobin"
			case Opportunistic:
				modeStr += "Opportunistic"
			}
			logRuntimeLimiter(">>>iter: %s", modeStr)
		}
		// TODO: fetch in different way for opportunistic (uses sorted list)
		var logic *LogicUnit
		switch mode {
		case RoundRobin:
			logic = r.logicUnits[r.runIx]
		case Opportunistic:
			logic = r.ascendingHotness[r.oppIx]
		}

		var func_ms float64
		if logic.active {

			logRuntimeLimiter("--- %s", logic.name)

			// check whether this logic has ever run
			_, hasRunBefore := r.lastRun[logic]
			// check its estimate
			estimate, hasEstimate := r.runtimeEstimates[logic]

			// estimate looks good if it's below allowance OR the estimate is above
			// allowance but we left off at this index last time; so we should get the
			// painful function over with rather than stall here forever or wait
			// to execute it when we get enough allowance (may never happen)
			// (first update remaining_ms so it's as accurate as possible)
			remaining_ms = allowance_ms - float64(time.Since(tStart).Nanoseconds())/1e6
			estimateLooksGood := hasEstimate && estimate <= remaining_ms
			logRuntimeLimiter("estimateLooksGood: %t", estimateLooksGood)
			logRuntimeLimiter("estimate: %f", estimate)
			// used to skip past iteration of this element in opportunistic mode
			oppSkip := false
			switch mode {
			case RoundRobin:
				// pop into opportunistic at first bad estimate of roundrobin
				// running the first roundrobin element regardless of time
				// estimate if Share() loop == 0 (bonsuTime is true)
				//
				// if the estimate is bad and we've run at least one func
				// then pop into opportunistic. Note that the behaviour such that
				// if the estimate is bad in round robin and we're at the first
				// func of the Run(), then run it regardless. We can never expect
				// to get more allowance than we have right now, so we might as well
				// get the heavy func out of the way.
				// r.runIx != r.startIx
				if hasEstimate && !estimateLooksGood {
					// only drop into opportunistic when runIx > startIx or bonsuTime
					//
					// in other words, since Run() defaults to roundrobin to
					// begin with, we will - when bonsuTime is false
					// (Share() loop == 0) - run the func regardless of
					// hasEstimate && !estimateLooksGood.
					// conversely,
					// when bonsuTime is true, we will immediately drop
					// into opportunistic if the first roundrobin element is
					// too heavy, and not run it
					if r.runIx != r.startIx || bonsuTime {
						mode = Opportunistic
						r.oppIx = 0
						// we sort the logics by hotness only when opportunistic
						// needs it, so it always represents the state of
						// things just when we popped into it initially.
						sort.Slice(r.ascendingHotness, func(i, j int) bool {
							return r.ascendingHotness[i].hotness < r.ascendingHotness[j].hotness
						})
						continue
					}
				}
			case Opportunistic:
				if hasEstimate && !estimateLooksGood {
					oppSkip = true
				}
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

			if DEBUG_RUNTIME_LIMITER {
				logRuntimeLimiter("hasRunBefore: %t", hasRunBefore)
				logRuntimeLimiter("isActive: %t", isActive)
				logRuntimeLimiter("durationHasElapsed: %t", durationHasElapsed)
				logRuntimeLimiter("scheduled: %t", scheduled)
			}
			if !hasRunBefore ||
				(durationHasElapsed && scheduled && !oppSkip) {

				t0 := time.Now()
				// note that we start lastrun from the moment the function starts, since
				// say it starts at t=0ms, takes 4 ms to run, then if it comes up to run
				// again at t=8ms (r.tick()), it will get dt_ms of 8 ms, the proper
				// intervening time since it last integrated a dt_ms increment.
				r.lastRun[logic] = time.Now()
				switch mode {
				case RoundRobin:
					logRuntimeLimiter("----------------------------------------- ROUND_ROBIN: %s", logic.name)
				case Opportunistic:
					logRuntimeLimiter("----------------------------------------- OPPORTUNISTIC: %s", logic.name)
				}
				logic.f(dt_ms)
				func_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
				logic.hotness++
				r.normalizeHotness(logic.hotness)
				r.lastEnd[logic] = time.Now()
				r.updateEstimate(logic, func_ms)
				switch mode {
				case RoundRobin:
					r.ranRobin++
				case Opportunistic:
					r.ranOpp++
				}
				remaining_ms = allowance_ms - float64(time.Since(tStart).Nanoseconds())/1e6
				logRuntimeLimiter("remaining after %s: %f", logic.name, remaining_ms)
			}
		}

		// end round-robin iteration if we reached back to where we started
		if mode == RoundRobin {
			r.runIx = (r.runIx + 1) % len(r.logicUnits)
			if r.runIx == r.startIx {
				r.finished = true
				break
			}
		}
		// end opportunistic iteration if we've looked at all the funcs there are
		// to run
		if mode == Opportunistic {
			r.oppIx = (r.oppIx + 1) % len(r.logicUnits)
			if r.oppIx == 0 {
				r.finished = true
				break
			}
		}

		overhead := float64(time.Since(tLoop).Nanoseconds())/1e6 - func_ms
		if overhead > worstOverheadThisTime {
			worstOverheadThisTime = overhead
		}
	}
	r.updateOverhead(worstOverheadThisTime)
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
	if r.ranRobin == 0 {
		r.starvation = 1.0
	} else if r.ranRobin > 0 && r.ranRobin <= len(r.logicUnits) {
		r.starvation = float64(len(r.logicUnits)-r.ranRobin) / float64(len(r.logicUnits))
	} else if r.ranRobin > len(r.logicUnits) {
		r.starvation = 0.0
	}
	return overunder_ms
}

func (r *RuntimeLimiter) tick(logic *LogicUnit) bool {
	if t, ok := r.lastEnd[logic]; ok {
		return float64(time.Since(t).Nanoseconds())/1.0e6 > r.runtimeEstimates[logic]
	} else {
		return true
	}
}

func (r *RuntimeLimiter) updateOverhead(worstThisTime float64) {
	if worstThisTime > r.loopOverhead_ms {
		r.loopOverhead_ms = worstThisTime
	} else {
		// else decay toward better worst overhead
		r.loopOverhead_ms = 0.9*r.loopOverhead_ms + 0.1*worstThisTime
	}
}

// every time we increment a logic's hotness, we check if it is now max int
// if it is, we reset hotness of all funcs to 0, call it a debt jubilee
func (r *RuntimeLimiter) normalizeHotness(hot int) {
	maxInt := int(^uint(0) >> 1)
	if hot == maxInt {
		for _, logic := range r.logicUnits {
			logic.hotness = 0
		}
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
	r.insertAscendingHotness(logic)
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
	// update runIx - if we removed an entity earlier in the list,
	// we should subtract 1 to keep runIx at it's same position. If we
	// removed one later in the list or equal to the current position,
	// we do nothing
	if index < r.runIx {
		r.runIx--
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

func (r *RuntimeLimiter) insertAscendingHotness(logic *LogicUnit) {
	i := sort.Search(len(r.ascendingHotness),
		func(ix int) bool { return r.ascendingHotness[ix].hotness > logic.hotness })
	r.ascendingHotness = append(r.ascendingHotness, nil)
	copy(r.ascendingHotness[i+1:], r.ascendingHotness[i:])
	r.ascendingHotness[i] = logic
}
