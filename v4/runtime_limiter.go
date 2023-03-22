/*
RuntimeLimiter is used to store a set of logicUnits which want to be

executed together and frequently, but which are tolerant of being
partitioned in time in order to stay within a certain time constraint
(for example, running all the world logic we can within 4 milliseconds,
hopefully looping back around to wherever we started, but if not, picking
up where we left off next Run()

it traverses logics round-robin until it reaches the first one that can't fit
in the time remaining, then switches into opportunistic mode: it uses
logic.hotness to sort them in ascending frequency of invocation and execute
any as it iterates that sorted list which it can, excluding those
which can't run in the remaining_ms when we get to them
*/

package sameriver

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/TwiN/go-color"
)

// TODO: export fields so you can just poll the stats directly
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
	// parallel array to ascendingHotness which records how light is the lightest
	// logicunit in ascendingHotness[i:]
	ascendingHotnessLightestAfter []float64

	// used to retrieve units by string
	logicUnitsMap map[string]*LogicUnit
	// track which logics have been removed even during a Run()
	removed map[*LogicUnit]bool
	// used to queue add/remove events so we don't change the slice while iterating
	addRemoveChannel chan AddRemoveLogicEvent
	// used to estimate the time cost in milliseconds of running a function,
	// so that we can try to stay below the allowance_ms given to Run()
	runtimeEstimates map[*LogicUnit]float64
	// used to provide an accurate dt_ms to each logic unit so it can integrate
	// time smooth and proper
	lastRun map[*LogicUnit]time.Time
	// used to keep track of whether the schedule period has elapsed for each logic
	lastScheduleTick map[*LogicUnit]time.Time
	// we run a logic unit with a gap of at least x ms where it takes x ms
	// to run. so a function taking 4ms will have at least 4 ms after it finishes
	// til the next time it runs, so we need to keep track of when logicunits end.
	lastEnd map[*LogicUnit]time.Time
	// used to lookup the logicUnits slice index for access in the slice
	indexes map[*LogicUnit]int

	// used to keep a running average of the entire runtime
	totalRuntime_ms *float64
	// ran : number that ran by any means
	ran int
	// ranRobin : number that ran by round robin since last time bonsuTime = false
	// (that is, when loop = 0 in Share())
	ranRobin int
	// ranOpp : number that ran by opportunistic since last time bonsuTime = false
	ranOpp int
	// iterated : total number of logicunits considered
	iterated int
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
	// the amount of time we didn't use in the last Run() call (negative if
	// we overran the allowance)
	overunder_ms float64
}

func NewRuntimeLimiter() *RuntimeLimiter {
	return &RuntimeLimiter{
		logicUnits:                    make([]*LogicUnit, 0),
		logicUnitsMap:                 make(map[string]*LogicUnit),
		removed:                       make(map[*LogicUnit]bool),
		addRemoveChannel:              make(chan (AddRemoveLogicEvent), ADD_REMOVE_LOGIC_CHANNEL_CAPACITY),
		ascendingHotness:              make([]*LogicUnit, 0),
		ascendingHotnessLightestAfter: make([]float64, 0),
		runtimeEstimates:              make(map[*LogicUnit]float64),
		lastRun:                       make(map[*LogicUnit]time.Time),
		lastScheduleTick:              make(map[*LogicUnit]time.Time),
		lastEnd:                       make(map[*LogicUnit]time.Time),
		indexes:                       make(map[*LogicUnit]int),
	}
}

type IterMode int

const (
	RoundRobin IterMode = iota
	Opportunistic
)

func (r *RuntimeLimiter) Run(allowance_ms float64, shareLoop int) {
	tStart := time.Now()
	logRuntimeLimiter("Run(); allowance: %f ms", allowance_ms)
	poll_remaining_ms := func() float64 {
		return allowance_ms - float64(time.Since(tStart).Nanoseconds())/1e6
	}

	r.ProcessAddRemoveLogics()

	if len(r.logicUnits) == 0 {
		logRuntimeLimiter("no logic units to run")
		r.finished = true
		r.starvation = 0
		return
	}

	if shareLoop == 0 {
		r.loopZero()
	}

	mode := RoundRobin
	worstOverheadThisTime := 0.0
	remaining_ms := poll_remaining_ms()
	for remaining_ms > 0 && (shareLoop > 0 || !r.finished) {
		logRuntimeLimiter("[\\] remaining_ms: %f", remaining_ms)
		tLoop := time.Now()

		// select logic according to mode
		logic, bail, skip := r.iter(mode, remaining_ms, shareLoop > 0)
		if bail {
			break
		}

		// run function (if it should run)
		var func_ms float64
		if !skip && r.shouldRunOrSwitchMode(logic, &mode, poll_remaining_ms(), shareLoop > 0) {
			func_ms = r.run(logic, mode)
			logRuntimeLimiter("remaining after %s: %f", logic.name, remaining_ms)
		}
		remaining_ms = poll_remaining_ms()

		// step our iteration index according to mode
		r.advanceIter(mode, shareLoop > 0)

		// track worst overhead
		overhead := float64(time.Since(tLoop).Nanoseconds())/1e6 - func_ms
		if overhead > worstOverheadThisTime {
			worstOverheadThisTime = overhead
		}
	}

	// update run stats
	total_ms := float64(time.Since(tStart).Nanoseconds()) / 1.0e6
	r.updateState(worstOverheadThisTime, allowance_ms, total_ms)
}

func (r *RuntimeLimiter) loopZero() {
	r.startIx = r.runIx
	r.finished = false
	r.ranRobin = 0
	r.ranOpp = 0
	r.oppIx = 0
	r.iterated = 0
	r.starvation = 1
	r.initShouldRun()
}

func (r *RuntimeLimiter) initShouldRun() {
	for _, l := range r.logicUnits {
		durationHasElapsed := r.tick(l)
		// copy the schedule so we can Tick it without it being damaged for
		// when we want to check it in the loop again and again
		hasSchedule := l.runSchedule != nil
		var scheduled bool
		if hasSchedule {
			scheduleAccum := *l.runSchedule
			schedule_tick_ms := float64(time.Since(r.lastScheduleTick[l]).Nanoseconds()) / 1e6
			scheduled = hasSchedule && scheduleAccum.Tick(schedule_tick_ms)
		} else {
			scheduled = true
		}
		_, removed := r.removed[l]
		l.shouldRun = l.active && !removed && durationHasElapsed && (!hasSchedule || scheduled)
		if l.shouldRun {
			logRuntimeLimiter("%s should run", l.name)
		}
		l.ran = false
	}
}

func (r *RuntimeLimiter) iter(mode IterMode, remaining_ms float64, bonsuTime bool) (logic *LogicUnit, bail bool, skip bool) {
	if remaining_ms < r.loopOverhead_ms {
		logRuntimeLimiter("XXX RUN() OVERHEAD BAIL - overhead is %.5f XXX", r.loopOverhead_ms)
		return nil, true, true
	}
	switch mode {
	case RoundRobin:
		logic = r.logicUnits[r.runIx]
	case Opportunistic:
		logic = r.ascendingHotness[r.oppIx]
		if r.ascendingHotnessLightestAfter[r.oppIx] >= remaining_ms {
			logRuntimeLimiter("XXX OPPORTUNISTIC RUN() BAIL, nothing light ahead XXX")
			return nil, true, true
		}
	}
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
	r.iterated++
	logRuntimeLimiter(color.InWhiteOverBlack(logic.name))
	_, removed := r.removed[logic]
	logRuntimeLimiter("active: %t, removed: %t", logic.active, removed)
	skip = !logic.active || removed
	return logic, false, skip
}

func (r *RuntimeLimiter) shouldRunOrSwitchMode(logic *LogicUnit, mode *IterMode, remaining_ms float64, bonsuTime bool) bool {
	// check whether this logic has ever run
	_, hasRunBefore := r.lastRun[logic]
	// check its estimate
	estimate, hasEstimate := r.runtimeEstimates[logic]

	// estimate looks good if it's below allowance OR the estimate is above
	// allowance but we left off at this index last time; so we should get the
	// painful function over with rather than stall here forever or wait
	// to execute it when we get enough allowance (may never happen)
	// (first update remaining_ms so it's as accurate as possible)
	estimateLooksGood := estimate <= remaining_ms
	logRuntimeLimiter("estimateLooksGood: %t", estimateLooksGood)
	logRuntimeLimiter("estimate: %f", estimate)
	switch *mode {
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
				logRuntimeLimiter("Dropping into opportunistic")
				*mode = Opportunistic
				// we sort the logics by hotness only when opportunistic
				// needs it, so it always represents the state of
				// things just when we popped into it initially.
				r.refreshAscendingHotness()
				return false
			}
		}
	case Opportunistic:
		if hasEstimate && !estimateLooksGood {
			logRuntimeLimiter("opportunistic skipping bad estimate logic")
			return false
		}
	}

	// if the time since the last run of this logic is > the runtime estimate
	// (that is, a function taking 1ms to run on avg should run at most
	// every 1ms)
	durationHasElapsed := r.tick(logic)

	// tick schedule
	schedule_tick_ms := float64(time.Since(r.lastScheduleTick[logic]).Nanoseconds()) / 1e6
	r.lastScheduleTick[logic] = time.Now()
	hasSchedule := logic.runSchedule != nil
	scheduled := hasSchedule && logic.runSchedule.Tick(schedule_tick_ms)

	if DEBUG_RUNTIME_LIMITER {
		logRuntimeLimiter("hasRunBefore: %t", hasRunBefore)
		logRuntimeLimiter("durationHasElapsed: %t", durationHasElapsed)
		logRuntimeLimiter("hasSchedule: %t", hasSchedule)
		logRuntimeLimiter("scheduled: %t", scheduled)
	}

	return (!hasRunBefore && !hasSchedule) || (durationHasElapsed && (!hasSchedule || scheduled))
}

func (r *RuntimeLimiter) run(logic *LogicUnit, mode IterMode) (func_ms float64) {
	// note that we start lastrun from the moment the function starts, since
	// say it starts at t=0ms, takes 4 ms to run, then if it comes up to run
	// again at t=8ms (r.tick()), it will get dt_ms of 8 ms, the proper
	// intervening time since it last integrated a dt_ms increment.
	if DEBUG_RUNTIME_LIMITER {
		switch mode {
		case RoundRobin:
			logRuntimeLimiter("----------------------------------------- " + color.InGreen(fmt.Sprintf("ROUND_ROBIN: %s", logic.name)))
		case Opportunistic:
			logRuntimeLimiter("----------------------------------------- " + color.InCyan(fmt.Sprintf("OPPORTUNISTIC: %s", logic.name)))
		}
	}
	t0 := time.Now()
	// get real time since last run
	dt_ms := float64(time.Since(r.lastRun[logic]).Nanoseconds()) / 1e6
	r.lastRun[logic] = time.Now()
	logic.f(dt_ms)
	func_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
	logic.ran = true
	logic.hotness++
	r.normalizeHotness(logic.hotness)
	r.lastEnd[logic] = time.Now()
	r.updateEstimate(logic, func_ms)
	r.ran++
	switch mode {
	case RoundRobin:
		r.ranRobin++
	case Opportunistic:
		r.ranOpp++
	}
	return func_ms
}

func (r *RuntimeLimiter) tick(logic *LogicUnit) bool {
	if t, ok := r.lastRun[logic]; ok {
		return float64(time.Since(t).Nanoseconds())/1.0e6 > r.runtimeEstimates[logic]
	} else {
		return true
	}
}

func (r *RuntimeLimiter) advanceIter(mode IterMode, bonsuTime bool) {
	// end round-robin iteration if we reached back to where we started
	if mode == RoundRobin {
		r.runIx = (r.runIx + 1) % len(r.logicUnits)
		// on loop = 0, the initial share, we just try to run everything once.
		// but thereafter, we will loop in roundrobin as long as we can, not
		// breaking at runix = startix
		if r.runIx == r.startIx && !bonsuTime {
			r.finished = true
			return
		}
	}
	// just plain loop opportunistic, we will bail according to the result of
	// r.iter() if needed
	if mode == Opportunistic {
		r.oppIx = (r.oppIx + 1) % len(r.logicUnits)
	}
}

func (r *RuntimeLimiter) updateOverhead(worstThisTime float64) {
	if worstThisTime > r.loopOverhead_ms {
		r.loopOverhead_ms = worstThisTime
	} else {
		// else decay toward better worst overhead
		r.loopOverhead_ms = 0.5*r.loopOverhead_ms + 0.5*worstThisTime
	}
}

// update various online calculations that can be read by someone using the RuntimeLimiter
func (r *RuntimeLimiter) updateState(worstOverheadThisTime, allowance_ms, total_ms float64) {
	// update overhead
	r.updateOverhead(worstOverheadThisTime)
	// maintain moving average of totalRuntime_ms
	logRuntimeLimiter(color.InWhiteOverGreen(fmt.Sprintf("Run() total: %f ms", total_ms)))
	if r.totalRuntime_ms == nil {
		r.totalRuntime_ms = &total_ms
	} else {
		*r.totalRuntime_ms = (*r.totalRuntime_ms + total_ms) / 2.0
	}
	// calculate overunder
	r.overunder_ms = allowance_ms - total_ms
	// calculate starved
	starved := 0
	for _, l := range r.logicUnits {
		if l.shouldRun && !l.ran {
			starved++
		}
	}
	r.starvation = float64(starved) / float64(len(r.logicUnits))
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

func (r *RuntimeLimiter) ProcessAddRemoveLogics() {
	for len(r.addRemoveChannel) > 0 {
		ev := <-r.addRemoveChannel
		l := ev.l
		if ev.addRemove {
			r.addLogicImmediately(l)
		} else {
			r.removeLogicImmediately(l)
		}
	}
}

func (r *RuntimeLimiter) addLogicImmediately(l *LogicUnit) {
	// panic if adding duplicate by WorldID
	if _, ok := r.indexes[l]; ok {
		panic(fmt.Sprintf("Double-add of same logic unit to RuntimeLimiter "+
			"(name: %s)", l.name))
	}
	r.logicUnits = append(r.logicUnits, l)
	r.logicUnitsMap[l.name] = l
	r.lastScheduleTick[l] = time.Now()
	r.indexes[l] = len(r.logicUnits) - 1
	r.insertAscendingHotness(l)
}

func (r *RuntimeLimiter) removeLogicImmediately(l *LogicUnit) {
	// return early if nil
	if l == nil {
		return
	}
	// return early if not present
	index, ok := r.indexes[l]
	if !ok {
		return
	}

	// delete from logicUnits by replacing the last element into its spot,
	// updating the indexes entry for that element
	removeFromLogicSlice := func(ls []*LogicUnit, i int) {
		lastIndex := len(ls) - 1
		if len(r.logicUnits) > 1 {
			r.logicUnits[i] = r.logicUnits[lastIndex]
		}
		r.logicUnits = r.logicUnits[:lastIndex]
	}
	// remove from logicUnits
	removeFromLogicSlice(r.logicUnits, index)
	if len(r.logicUnits) != 0 {
		// update indexes for last-now-here element
		nowAtIndex := r.logicUnits[index]
		r.indexes[nowAtIndex] = index
	}
	// remove from ascending hotness slices
	// (first find lowest index with common hotness w binary search)
	left, right := 0, len(r.ascendingHotness)-1
	lowestIx := -1
	lucky := false
	for left <= right {
		mid := left + (right-left)/2
		// special case if we happen to hit it during iteration, bail early;
		// we don't need to iterate from lowest index with common hotness
		// to find it now
		if r.ascendingHotness[mid] == l {
			removeFromLogicSlice(r.ascendingHotness, mid)
			lucky = true
			break
		}
		if r.ascendingHotness[mid].hotness == l.hotness {
			lowestIx = mid
			right = mid - 1 // Continue searching the left side for the lowest index
		} else if r.ascendingHotness[mid].hotness < l.hotness {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	if !lucky {
		for i := lowestIx; i < len(r.ascendingHotness); i++ {
			if r.ascendingHotness[i] == l {
				removeFromLogicSlice(r.ascendingHotness, i)
				break
			}
		}
	}
	r.refreshAscendingHotnessLightestAfter()

	delete(r.logicUnitsMap, l.name)
	delete(r.removed, l)
	delete(r.runtimeEstimates, l)
	delete(r.lastRun, l)
	delete(r.lastEnd, l)
	delete(r.indexes, l)

	// update runIx - if we removed an entity earlier in the list,
	// we should subtract 1 to keep runIx at it's same position. If we
	// removed one later in the list or equal to the current position,
	// we do nothing
	if index < r.runIx {
		r.runIx--
	}
}

func (r *RuntimeLimiter) Add(l *LogicUnit) {
	do := func() {
		r.addRemoveChannel <- AddRemoveLogicEvent{addRemove: true, l: l}
	}
	if len(r.addRemoveChannel) >= ADD_REMOVE_LOGIC_CHANNEL_CAPACITY {
		logWarning("adding logic at such a rate the channel is at capacity. Spawning goroutines. If this continues to happen, the program might suffer.")
		go do()
	} else {
		do()
	}
}

func (r *RuntimeLimiter) Remove(l *LogicUnit) {
	do := func() {
		r.removed[l] = true
		r.addRemoveChannel <- AddRemoveLogicEvent{addRemove: false, l: l}
	}
	if len(r.addRemoveChannel) >= ADD_REMOVE_LOGIC_CHANNEL_CAPACITY {
		logWarning("removing logic at such a rate the channel is at capacity. Spawning goroutines. If this continues to happen, the program might suffer.")
		go do()
	} else {
		do()
	}
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

func (r *RuntimeLimiter) SetSchedule(logicName string, period_ms float64) {
	logic := r.logicUnitsMap[logicName]
	runSchedule := NewTimeAccumulator(period_ms)
	logic.runSchedule = &runSchedule
}

func (r *RuntimeLimiter) Finished() bool {
	return r.finished
}

func (r *RuntimeLimiter) refreshAscendingHotness() {
	sort.Slice(r.ascendingHotness, func(i, j int) bool {
		return r.ascendingHotness[i].hotness < r.ascendingHotness[j].hotness
	})
	r.refreshAscendingHotnessLightestAfter()
}

func (r *RuntimeLimiter) refreshAscendingHotnessLightestAfter() {
	if len(r.ascendingHotness) == 0 {
		r.ascendingHotnessLightestAfter = make([]float64, 0)
		return
	}
	// make sure the receiving array we will write to has the same length as
	// the ascendinghotness array
	if len(r.ascendingHotness) > len(r.ascendingHotnessLightestAfter) {
		// receiver smaller
		diff := len(r.ascendingHotness) - len(r.ascendingHotnessLightestAfter)
		emptySpace := make([]float64, diff)
		r.ascendingHotnessLightestAfter = append(r.ascendingHotnessLightestAfter, emptySpace...)
	} else if len(r.ascendingHotness) < len(r.ascendingHotnessLightestAfter) {
		// receiver bigger
		diff := len(r.ascendingHotnessLightestAfter) - len(r.ascendingHotness)
		endIx := len(r.ascendingHotnessLightestAfter) - diff
		r.ascendingHotnessLightestAfter = r.ascendingHotnessLightestAfter[0:endIx]
	}
	// compute lightest after by walking ascendingHotness backward
	lightest := math.Inf(1)
	for i := len(r.ascendingHotnessLightestAfter) - 1; i >= 0; i-- {
		weight := r.runtimeEstimates[r.ascendingHotness[i]]
		if weight < lightest {
			lightest = weight
		}
		r.ascendingHotnessLightestAfter[i] = lightest
	}
}

func (r *RuntimeLimiter) insertAscendingHotness(l *LogicUnit) {
	if len(r.ascendingHotness) == 0 {
		l.hotness = 0
		r.ascendingHotness = append(r.ascendingHotness, l)
	} else {
		// put it at [0] with the hotness of the old [0]
		r.ascendingHotness = append(r.ascendingHotness, nil)
		copy(r.ascendingHotness[1:], r.ascendingHotness[0:])
		l.hotness = r.ascendingHotness[1].hotness
		r.ascendingHotness[0] = l
	}
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
	stats["__numberOfLogicUnits"] = float64(len(r.logicUnits))
	stats["__starvation"] = r.starvation
	stats["__ran"] = float64(r.ran)
	stats["__ranRobin"] = float64(r.ranRobin)
	stats["__ranOpp"] = float64(r.ranOpp)
	stats["__overunder_ms"] = float64(r.overunder_ms)
	return
}
