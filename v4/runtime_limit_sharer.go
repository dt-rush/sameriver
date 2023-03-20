/*
RuntimeLimitSharer provides a mechanism to distribute available processing time
across multiple RuntimeLimiters (each of which tries to run a set of LogicUnits)
in a controlled and efficient manner.

The RuntimeLimitSharer manages the registration, sharing of processing time,
and monitoring of starvation for each RuntimeLimiter.

The main components of the package are:

In the Share() function, the available processing time (allowance_ms) is
distributed across the registered RuntimeLimiters in a loop, and the method
is allowed to a finite number of times (RUNTIME_LIMIT_SHARER_MAX_LOOPS).

The partitioning of remaining_ms per loop in Share() can happen in two ways:

Equal partitioning:

|____|____|____|____|

The remaining_ms is divided equally among all RuntimeLimiters.

Starvation-proportional partitioning:

|_||____|___________|

When a RuntimeLimiter has experienced starvation, the remaining_ms is
divided according to the proportion of the starvation experienced by
each RuntimeLimiter.
*/

package sameriver

import (
	"fmt"
	"time"

	"github.com/TwiN/go-color"
)

type RuntimeLimitSharer struct {
	runIX       int
	runners     []*RuntimeLimiter
	RunnerMap   map[string]*RuntimeLimiter
	runnerNames map[*RuntimeLimiter]string
	// normally we evenly divide, but this can be overridden
	InitialShareScale map[string]float64
	addRemoveChannel  chan AddRemoveLogicEvent
}

func NewRuntimeLimitSharer() *RuntimeLimitSharer {
	r := &RuntimeLimitSharer{
		runners:           make([]*RuntimeLimiter, 0),
		RunnerMap:         make(map[string]*RuntimeLimiter),
		runnerNames:       make(map[*RuntimeLimiter]string),
		InitialShareScale: make(map[string]float64),
		addRemoveChannel:  make(chan (AddRemoveLogicEvent), ADD_REMOVE_LOGIC_CHANNEL_CAPACITY),
	}
	return r
}

func (r *RuntimeLimitSharer) registerRunner(name string, p float64) *RuntimeLimiter {
	if _, ok := r.RunnerMap[name]; ok {
		panic(fmt.Sprintf("Trying to double-add RuntimeLimiter %s", name))
	}
	runner := NewRuntimeLimiter()
	r.runners = append(r.runners, runner)
	r.RunnerMap[name] = runner
	r.runnerNames[runner] = name
	r.InitialShareScale[name] = p
	return runner
}

func (r *RuntimeLimitSharer) RegisterRunners(spec map[string]float64) {
	sum := 0.0
	for _, k := range spec {
		sum += k
	}
	for name, k := range spec {
		r.registerRunner(name, k/sum)
	}
}

func (r *RuntimeLimitSharer) Share(allowance_ms float64) (overunder_ms float64, starved int) {
	tStart := time.Now()
	overunder_ms = allowance_ms
	// while we have allowance_ms, keep trying to run all runners
	// note: everybody gets firsts before anyone gets seconds; this is controlled
	// using starvedMode.
	// and, to avoid spinning way too many times when load is light,
	// we have MAX_LOOPS set to an arbitrary 8 (8 update cycles per
	// frame is not bad! haha)
	loop := 0
	remaining_ms := allowance_ms
	starvedMode := false
	var lastStarvation float64
	logRuntimeLimiter("\n====================\nshare loop\n====================\n")
	for remaining_ms >= 0 && loop < RUNTIME_LIMIT_SHARER_MAX_LOOPS {
		toShare_ms := remaining_ms
		logRuntimeLimiter("\n===\nloop = %d, total share = %f ms\n===\n", loop, toShare_ms)
		totalStarvation := 0.0
		considered := 0
		var ran int
		for ran = 0; remaining_ms >= 0 && considered < len(r.runners); {
			considered++
			runner := r.runners[r.runIX]
			var runnerAllowance float64
			// if not starved, divide according to initialsharescale
			if !starvedMode {
				p := r.InitialShareScale[r.runnerNames[runner]]
				runnerAllowance = toShare_ms * p
			} else {
				logRuntimeLimiter("|||||| %f * (%f / %f)", toShare_ms, runner.starvation, lastStarvation)
				runnerAllowance = toShare_ms * (runner.starvation / lastStarvation)
			}

			logRuntimeLimiter("%s.starvation = %f", r.runnerNames[runner], runner.starvation)
			logRuntimeLimiter("Run()? starvedMode: %t, starvedMode: %t, runner.starvation: %f", starvedMode, starvedMode, runner.starvation)
			if !starvedMode || (starvedMode && runner.starvation != 0) {
				logRuntimeLimiter(color.InWhiteOverBlue(fmt.Sprintf("|||||| sharing %f ms to %s", runnerAllowance, r.runnerNames[runner])))
				// loop > 0 is the parameter of Run(), bonsuTime (AKA bonusTime)
				runner.Run(runnerAllowance, loop > 0)
				totalStarvation += runner.starvation
				if runner.starvation != 0 {
					logRuntimeLimiter(color.InYellow(fmt.Sprintf("%s.starvation = %f", r.runnerNames[runner], runner.starvation)))
				}
				remaining_ms = allowance_ms - float64(time.Since(tStart).Nanoseconds())/1e6
				logRuntimeLimiter(color.InWhiteOverBlue(fmt.Sprintf("[remaining_ms: %f]", remaining_ms)))
				ran++
			}

			r.runIX = (r.runIX + 1) % len(r.runners)
		}
		if ran == 0 {
			break
		} else {
			starvedMode = (totalStarvation > 0)
			lastStarvation = totalStarvation
		}
		loop++
	}
	if DEBUG_RUNTIME_LIMITER && loop == RUNTIME_LIMIT_SHARER_MAX_LOOPS {
		logRuntimeLimiter("Reached MAX_LOOPS in RuntimeSharer with %f percent time remaining", 100*remaining_ms/allowance_ms)
	}
	// above we were concerned with starvation of logics inside runners, now
	// we are concerned with starvation of entire runners. This can happen
	// when a runner that we encounter as we iterate the runners uses up, in
	// one logic func, more than its own budget + another, so that we quit
	// the runner iteration at remaining <= 0 before the later runner(s) got
	// a chance to even run.
	// starved 0 means they all ran once (even if they didn't complete*)
	// starved < 0 means at least one ran more than once
	// starved > 0 means at least one didn't run
	starved = 0
	for i := 0; i < len(r.runners); i++ {
		if r.runners[i].starvation > 0 {
			starved++
		}
	}
	return remaining_ms, starved
}

func (r *RuntimeLimitSharer) DumpStats() map[string](map[string]float64) {
	stats := make(map[string](map[string]float64))
	stats["totals"] = make(map[string]float64)
	stats["starvation"] = make(map[string]float64)
	for name, r := range r.RunnerMap {
		runnerStats, totals := r.DumpStats()
		stats[name] = runnerStats
		stats["starvation"][name] = r.starvation
		stats["totals"][name] = totals
	}
	return stats
}
