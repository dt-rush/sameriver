package engine

import (
	"fmt"
	"testing"
	"time"
)

func TestRuntimeLimiterAdd(t *testing.T) {
	r := NewRuntimeLimiter()
	for i := 0; i < 32; i++ {
		name := fmt.Sprintf("logic-%d", i)
		logic := &LogicUnit{
			name:    name,
			worldID: i,
			f:       func() {},
			active:  true}
		r.Add(logic)
		if !(len(r.logicUnits) > 0 &&
			r.indexes[logic.worldID] == len(r.logicUnits)-1) {
			t.Fatal("was not inserted properly")
		}
	}
}

func TestRuntimeLimiterAddDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	r := NewRuntimeLimiter()
	logic := &LogicUnit{
		name:    "logic",
		worldID: 0,
		f:       func() {},
		active:  true}
	r.Add(logic)
	r.Add(logic)
	t.Fatal("should have panic'd")
}

func TestRuntimeLimiterRun(t *testing.T) {
	r := NewRuntimeLimiter()
	x := 0
	name := "l1"
	r.Add(&LogicUnit{
		name:    name,
		worldID: 0,
		f:       func() { x += 1 },
		active:  true})
	for i := 0; i < 32; i++ {
		r.Start()
		r.Run(FRAME_SLEEP_MS)
	}
	if x != 32 {
		t.Fatal("didn't run logic")
	}
	if !r.Finished() {
		t.Fatal("should have returned finished = true when running sole " +
			"logic within time limit")
	}
}

func TestRuntimeLimiterOverrun(t *testing.T) {
	r := NewRuntimeLimiter()
	r.Add(&LogicUnit{
		name:    "logic",
		worldID: 0,
		f:       func() { time.Sleep(150 * time.Millisecond) },
		active:  true})
	r.Start()
	remaining_ms := r.Run(100)
	if remaining_ms > 0 {
		t.Fatal("overrun time not calculated properly")
	}
}

func TestRuntimeLimiterUnderrun(t *testing.T) {
	r := NewRuntimeLimiter()
	r.Add(&LogicUnit{
		name:    "logic",
		worldID: 0,
		f:       func() { time.Sleep(100 * time.Millisecond) },
		active:  true})
	r.Start()
	remaining_ms := r.Run(300)
	if !(remaining_ms > 0 && remaining_ms <= 200) {
		t.Fatal("underrun time not calculated properly")
	}
}

func TestRuntimeLimiterLimiting(t *testing.T) {
	r := NewRuntimeLimiter()
	fastRan := false
	r.Add(&LogicUnit{
		name:    "logic-slow",
		worldID: 0,
		f:       func() { time.Sleep(10 * time.Millisecond) },
		active:  true})
	r.Add(&LogicUnit{
		name:    "logic-slow",
		worldID: 1,
		f:       func() { fastRan = true },
		active:  true})
	r.Start()
	r.Run(2)
	if fastRan {
		t.Fatal("continued running logic despite using up allowed milliseconds")
	}
}

func TestRuntimeLimiterDoNotRunEstimatedSlow(t *testing.T) {
	r := NewRuntimeLimiter()
	r.Add(&LogicUnit{
		name:    "dummy",
		worldID: 0,
		f:       func() {},
		active:  true})
	x := 0
	name := "l1"
	ms_duration := 100
	r.Add(&LogicUnit{
		name:    name,
		worldID: 1,
		f: func() {
			x += 1
			time.Sleep(time.Duration(ms_duration) * time.Millisecond)
		},
		active: true})
	// since it's never run before, running the logic will set its estimate
	r.Start()
	r.Run(FRAME_SLEEP_MS)
	// now we try to run it again, but give it no time to run (exceeds estimate)
	r.Run(float64(ms_duration) / 10.0)
	if x == 2 {
		t.Fatal("ran logic even though estimate should have prevented this")
	}
}

func TestRuntimeLimiterRemove(t *testing.T) {
	r := NewRuntimeLimiter()
	// test that we can remove a logic which doens't exist idempotently
	if r.Remove(0) != false {
		t.Fatal("somehow removed a logic which doesn't exist")
	}
	x := 0
	name := "l1"
	logic := &LogicUnit{
		name:    name,
		worldID: 0,
		f:       func() { x += 1 },
		active:  true}
	r.Add(logic)
	// run logic a few times so that it has runtimeEstimate data
	for i := 0; i < 32; i++ {
		r.Start()
		r.Run(FRAME_SLEEP_MS)
	}
	// remove it
	r.Remove(0)
	// test if removed
	if _, ok := r.runtimeEstimates[logic]; ok {
		t.Fatal("did not delete runtimeEstimates data")
	}
	if _, ok := r.indexes[logic.worldID]; ok {
		t.Fatal("did not delete runtimeEstimates data")
	}
	if len(r.logicUnits) != 0 {
		t.Fatal("did not remove from logicUnits list")
	}
}

func TestRuntimeLimitShare(t *testing.T) {
	runners := make([]*RuntimeLimiter, 0)
	counters := make([]int, 0)
	for i := 0; i < 3; i++ {
		func(i int) {
			r := NewRuntimeLimiter()
			runners = append(runners, r)
			counters = append(counters, 0) // jet fuel can't melt steel beams
			r.Add(&LogicUnit{
				name:    "logic",
				worldID: 0,
				f:       func() { counters[i] += 1 },
				active:  true})
		}(i)
	}
	for i := 0; i < 32; i++ {
		RuntimeLimitShare(FRAME_SLEEP_MS, runners...)
	}
	for _, counter := range counters {
		if counter != 32 {
			t.Fatal("didn't share runtime properly")
		}
	}
}

func TestRuntimeLimiterRunDuringRemovals(t *testing.T) {
	r := NewRuntimeLimiter()
	for i := 0; i < 32; i++ {
		name := fmt.Sprintf("logic-%d", i)
		logic := &LogicUnit{
			name:    name,
			worldID: i,
			f:       func() {},
			active:  true}
		r.Add(logic)
		if !(len(r.logicUnits) > 0 &&
			r.indexes[logic.worldID] == len(r.logicUnits)-1) {
			t.Fatal("was not inserted properly")
		}
	}
}
