package sameriver

import (
	"fmt"
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
	addRemoveChannel chan AddRemoveLogicEvent
}

func NewRuntimeLimitSharer() *RuntimeLimitSharer {
	r := &RuntimeLimitSharer{
		runners:          make([]*RuntimeLimiter, 0),
		runnerMap:        make(map[string]*RuntimeLimiter),
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
}

func (r *RuntimeLimitSharer) ProcessAddRemoveLogics() {
	for len(r.addRemoveChannel) > 0 {
		ev := <-r.addRemoveChannel
		l := ev.l
		runnerName := ev.runnerName
		if ev.addRemove {
			// add
			if _, ok := r.runnerMap[runnerName]; !ok {
				panic(fmt.Sprintf("Trying to add to runtimeLimiter with name %s - doesn't exist", runnerName))
			}
			r.runnerMap[runnerName].Add(l)
		} else {
			// remove
			if _, ok := r.runnerMap[runnerName]; !ok {
				panic(fmt.Sprintf("Trying to remove from runtimeLimiter with name %s - doesn't exist", runnerName))
			}
			r.runnerMap[runnerName].Remove(l)
		}
	}
}

func (r *RuntimeLimitSharer) AddLogic(runnerName string, l *LogicUnit) {
	r.addRemoveChannel <- AddRemoveLogicEvent{
		addRemove:  true,
		runnerName: runnerName,
		l:          l,
	}
}

func (r *RuntimeLimitSharer) RemoveLogic(runnerName string, l *LogicUnit) {

	r.addRemoveChannel <- AddRemoveLogicEvent{
		addRemove:  false,
		runnerName: runnerName,
		l:          l,
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

func (r *RuntimeLimitSharer) Share(allowance_ms float64) (remaining_ms float64, starved int) {
	// process addition and removal of logics (they get buffered in a channel
	// so we aren't adding logics while iterating logics)
	r.ProcessAddRemoveLogics()

	remaining_ms = allowance_ms
	ran := 0
	perRunner_ms := allowance_ms / float64(len(r.runners))
	// while we have allowance_ms, keep trying to run all runners
	// note: everybody gets firsts before anyone gets seconds
	// and, to avoid spinning way too many times when load is light,
	// we have MAX_LOOPS set to an arbitrary 20 (20 update cycles per
	// frame is not bad)
	MAX_LOOPS := 20
	loops := 0
	for allowance_ms >= 0 && loops < MAX_LOOPS {
		for allowance_ms >= 0 && ran < len(r.runners) {
			runner := r.runners[r.runIX]
			overunder_ms := runner.Run(perRunner_ms)
			used := perRunner_ms - overunder_ms
			allowance_ms -= used
			// increment to run next runner even if runner.Finished() isn't true
			// this means it will get another chance to finish itself when its
			// turn comes back around
			r.runIX = (r.runIX + 1) % len(r.runners)
			ran++
		}
		starved = len(r.runners) - ran
		loops++
	}
	return allowance_ms, starved
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
