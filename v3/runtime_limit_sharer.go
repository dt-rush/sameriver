package sameriver

import (
	"fmt"
	"time"

	"github.com/TwiN/go-color"
)

type AddRemoveLogicEvent struct {
	addRemove  bool
	runnerName string
	l          *LogicUnit
}

type RuntimeLimitSharer struct {
	runIX            int
	runners          []*RuntimeLimiter
	runnerMap        map[string]*RuntimeLimiter
	runnerNames      map[*RuntimeLimiter]string
	addRemoveChannel chan AddRemoveLogicEvent

	// used to keep track of expected worst case loop overhead
	innerLoopOverhead_ms float64
}

func NewRuntimeLimitSharer() *RuntimeLimitSharer {
	r := &RuntimeLimitSharer{
		runners:          make([]*RuntimeLimiter, 0),
		runnerMap:        make(map[string]*RuntimeLimiter),
		runnerNames:      make(map[*RuntimeLimiter]string),
		addRemoveChannel: make(chan (AddRemoveLogicEvent), ADD_REMOVE_LOGIC_CHANNEL_CAPACITY),
	}
	return r
}

func (r *RuntimeLimitSharer) RegisterRunner(name string) {
	if _, ok := r.runnerMap[name]; ok {
		panic(fmt.Sprintf("Trying to double-add RuntimeLimiter %s", name))
	}
	runner := NewRuntimeLimiter()
	r.runners = append(r.runners, runner)
	r.runnerMap[name] = runner
	r.runnerNames[runner] = name
}

func (r *RuntimeLimitSharer) ProcessAddRemoveLogics() {
	for len(r.addRemoveChannel) > 0 {
		ev := <-r.addRemoveChannel
		l := ev.l
		runnerName := ev.runnerName
		if ev.addRemove {
			r.addLogicImmediately(runnerName, l)
		} else {
			r.removeLogicImmediately(runnerName, l)
		}
	}
}

func (r *RuntimeLimitSharer) addLogicImmediately(runnerName string, l *LogicUnit) {
	// add
	if _, ok := r.runnerMap[runnerName]; !ok {
		panic(fmt.Sprintf("Trying to add to runtimeLimiter with name %s - doesn't exist", runnerName))
	}
	r.runnerMap[runnerName].Add(l)
}

func (r *RuntimeLimitSharer) removeLogicImmediately(runnerName string, l *LogicUnit) {
	// remove
	if _, ok := r.runnerMap[runnerName]; !ok {
		panic(fmt.Sprintf("Trying to remove from runtimeLimiter with name %s - doesn't exist", runnerName))
	}
	r.runnerMap[runnerName].Remove(l)
}

func (r *RuntimeLimitSharer) AddLogic(runnerName string, l *LogicUnit) {
	do := func() {
		r.addRemoveChannel <- AddRemoveLogicEvent{
			addRemove:  true,
			runnerName: runnerName,
			l:          l,
		}
	}
	if len(r.addRemoveChannel) >= ADD_REMOVE_LOGIC_CHANNEL_CAPACITY {
		logWarning("adding logic at such a rate the channel is at capacity. Spawning goroutines. If this continues to happen, the program might suffer.")
		go do()
	} else {
		do()
	}
}

func (r *RuntimeLimitSharer) RemoveLogic(runnerName string, l *LogicUnit) {
	do := func() {
		r.addRemoveChannel <- AddRemoveLogicEvent{
			addRemove:  false,
			runnerName: runnerName,
			l:          l,
		}
	}
	if len(r.addRemoveChannel) >= ADD_REMOVE_LOGIC_CHANNEL_CAPACITY {
		logWarning("removing logic at such a rate the channel is at capacity. Spawning goroutines. If this continues to happen, the program might suffer.")
		go do()
	} else {
		do()
	}
}

func (r *RuntimeLimitSharer) ActivateAll(runnerName string) {
	if _, ok := r.runnerMap[runnerName]; !ok {
		panic(fmt.Sprintf("Trying to activate all in runtimeLimiter with name %s - doesn't exist", runnerName))
	}
	r.runnerMap[runnerName].ActivateAll()
}

func (r *RuntimeLimitSharer) DeactivateAll(runnerName string) {
	if _, ok := r.runnerMap[runnerName]; !ok {
		panic(fmt.Sprintf("Trying to deactivate all in runtimeLimiter with name %s - doesn't exist", runnerName))
	}
	r.runnerMap[runnerName].DeactivateAll()
}

func (r *RuntimeLimitSharer) SetSchedule(runnerName string, logicWorldID int, period_ms float64) {
	runner := r.runnerMap[runnerName]
	logicIX := runner.indexes[logicWorldID]
	logic := runner.logicUnits[logicIX]
	runSchedule := NewTimeAccumulator(period_ms)
	logic.runSchedule = &runSchedule
}

func (r *RuntimeLimitSharer) Share(allowance_ms float64) (overunder_ms float64, starved int) {
	tStart := time.Now()
	// process addition and removal of logics (they get buffered in a channel
	// so we aren't adding logics while iterating logics)
	r.ProcessAddRemoveLogics()

	overunder_ms = allowance_ms
	// while we have allowance_ms, keep trying to run all runners
	// note: everybody gets firsts before anyone gets seconds; this is controlled
	// using starvedMode.
	// and, to avoid spinning way too many times when load is light,
	// we have MAX_LOOPS set to an arbitrary 8 (8 update cycles per
	// frame is not bad! haha)
	MAX_LOOPS := 8
	loop := 0
	remaining_ms := allowance_ms
	starvedMode := false
	var lastStarvation float64
	logRuntimeLimiter("\n====================\nshare loop\n====================\n")
	overheadBail := false
	for remaining_ms >= 0 && loop < MAX_LOOPS && !overheadBail {
		toShare_ms := remaining_ms
		logRuntimeLimiter("\n===\nloop = %d, total share = %f ms\n===\n", loop, toShare_ms)
		totalStarvation := 0.0
		considered := 0
		worstOverheadThisTime := 0.0
		for ran := 0; remaining_ms >= 0 && considered < len(r.runners); {
			if remaining_ms < r.innerLoopOverhead_ms {
				logRuntimeLimiter("XXX SHARE() OVERHEAD BAIL XXX")
				overheadBail = true
				break
			}
			tLoop := time.Now()
			var used float64
			considered++
			runner := r.runners[r.runIX]
			var runnerAllowance float64
			if !starvedMode {
				runnerAllowance = toShare_ms / float64(len(r.runners))
			} else {
				logRuntimeLimiter("|||||| %f * (%f / %f)", toShare_ms, runner.starvation, lastStarvation)
				runnerAllowance = toShare_ms * (runner.starvation / lastStarvation)
			}
			logRuntimeLimiter("%s.starvation = %f", r.runnerNames[runner], runner.starvation)
			if !starvedMode || (starvedMode && runner.starvation != 0) {
				logRuntimeLimiter(color.InWhiteOverBlue(fmt.Sprintf("|||||| sharing %f ms to %s", runnerAllowance, r.runnerNames[runner])))
				// loop > 0 is the parameter of Run(), bonsuTime (AKA bonusTime)
				t0 := time.Now()
				runner.Run(runnerAllowance, loop > 0)
				totalStarvation += runner.starvation
				if runner.starvation != 0 {
					logRuntimeLimiter(color.InYellow(fmt.Sprintf("%s.starvation = %f", r.runnerNames[runner], runner.starvation)))
				}
				used = float64(time.Since(t0).Nanoseconds()) / 1e6
				remaining_ms = allowance_ms - float64(time.Since(tStart).Nanoseconds())/1e6
				logRuntimeLimiter(color.InWhiteOverBlue(fmt.Sprintf("[remaining_ms: %f]", remaining_ms)))
				ran++
			}

			overhead := float64(time.Since(tLoop).Nanoseconds())/1e6 - used
			if overhead > worstOverheadThisTime {
				worstOverheadThisTime = overhead
			}
			r.runIX = (r.runIX + 1) % len(r.runners)
		}
		r.updateOverhead(worstOverheadThisTime)
		starvedMode = (totalStarvation > 0)
		lastStarvation = totalStarvation
		loop++
	}
	if DEBUG_RUNTIME_LIMITER && loop == MAX_LOOPS {
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

func (r *RuntimeLimitSharer) updateOverhead(worstThisTime float64) {
	if worstThisTime > r.innerLoopOverhead_ms {
		r.innerLoopOverhead_ms = worstThisTime
	} else {
		// else decay toward better worst overhead
		r.innerLoopOverhead_ms = 0.9*r.innerLoopOverhead_ms + 0.1*worstThisTime
	}

}

func (r *RuntimeLimitSharer) DumpStats() map[string](map[string]float64) {
	stats := make(map[string](map[string]float64))
	stats["totals"] = make(map[string]float64)
	for name, r := range r.runnerMap {
		runnerStats, totals := r.DumpStats()
		stats[name] = runnerStats
		stats["totals"][name] = totals
	}
	return stats
}
