package sameriver

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
			f:       func(dt_ms float64) {},
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
		f:       func(dt_ms float64) {},
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
		f:       func(dt_ms float64) { x += 1 },
		active:  true})
	r.Run(1)
	for i := 0; i < 32; i++ {
		r.Run(FRAME_DURATION_INT)
	}
	Logger.Println(x)
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
		f:       func(dt_ms float64) { time.Sleep(150 * time.Millisecond) },
		active:  true})
	r.Run(1)
	remaining_ms := r.Run(100)
	if remaining_ms > 0 {
		t.Fatal("overrun time not calculated properly")
	}
	if !r.overrun {
		t.Fatal("didn't set overrun flag")
	}
}

func TestRuntimeLimiterUnderrun(t *testing.T) {
	r := NewRuntimeLimiter()
	r.Add(&LogicUnit{
		name:    "logic",
		worldID: 0,
		f:       func(dt_ms float64) { time.Sleep(100 * time.Millisecond) },
		active:  true})
	r.Run(1)
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
		f:       func(dt_ms float64) { time.Sleep(10 * time.Millisecond) },
		active:  true})
	r.Add(&LogicUnit{
		name:    "logic-slow",
		worldID: 1,
		f:       func(dt_ms float64) { fastRan = true },
		active:  true})
	r.Run(1)
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
		f:       func(dt_ms float64) {},
		active:  true})
	x := 0
	name := "l1"
	ms_duration := 100
	r.Add(&LogicUnit{
		name:    name,
		worldID: 1,
		f: func(dt_ms float64) {
			x += 1
			time.Sleep(time.Duration(ms_duration) * time.Millisecond)
		},
		active: true})
	// since it's never run before, running the logic will set its estimate
	r.Run(FRAME_DURATION_INT)
	// now we try to run it again, but give it no time to run (exceeds estimate)
	r.Run(float64(ms_duration) / 10.0)
	if x == 2 {
		t.Fatal("ran logic even though estimate should have prevented this")
	}
}

func TestRuntimeLimiterRemove(t *testing.T) {
	r := NewRuntimeLimiter()
	// test that we can remove a logic which doens't exist idempotently
	if r.Remove(nil) != false {
		t.Fatal("somehow removed a logic which doesn't exist")
	}
	x := 0
	name := "l1"
	logic := &LogicUnit{
		name:    name,
		worldID: 0,
		f:       func(dt_ms float64) { x += 1 },
		active:  true}
	r.Add(logic)
	// run logic a few times so that it has runtimeEstimate data
	for i := 0; i < 32; i++ {
		r.Run(FRAME_DURATION_INT)
	}
	// remove it
	Logger.Println(fmt.Sprintf("Removing logic: %s", logic.name))
	r.Remove(logic)
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
	w := testingWorld()
	sharer := NewRuntimeLimitSharer()
	counters := make([]int, 0)

	const N = 3
	const M = 3
	const LOOPS = 5
	const SLEEP = 16

	sharer.RegisterRunner("basic")
	for i := 0; i < N; i++ {
		func(i int) {
			counters = append(counters, 0) // jet fuel can't melt steel beams
			sharer.AddLogic("basic", &LogicUnit{
				name:    fmt.Sprintf("basic-%d", i),
				worldID: w.IdGen.Next(),
				f: func(dt_ms float64) {
					time.Sleep(SLEEP)
					counters[i] += 1
				},
				active: true})
		}(i)
	}
	sharer.RegisterRunner("extra")
	for i := 0; i < M; i++ {
		func(i int) {
			counters = append(counters, 0) // jet fuel can't melt steel beams
			sharer.AddLogic("extra", &LogicUnit{
				name:    fmt.Sprintf("extra-%d", i),
				worldID: w.IdGen.Next(),
				f: func(dt_ms float64) {
					time.Sleep(SLEEP)
					counters[i] += 1
				},
				active: true})
		}(i)
	}
	for i := 0; i < LOOPS; i++ {
		sharer.Share((N+M)*SLEEP + 100)
	}
	// -1 because the first run just sets the logicunits `lastRun` time.Time
	expected := N*(LOOPS-1) + M*(LOOPS-1)
	sum := 0
	for _, counter := range counters {
		sum += counter

	}
	if sum != expected {
		t.Fatal("didn't share runtime properly")
	}
}

func TestRuntimeLimitShareInsertWhileRunning(t *testing.T) {
	w := testingWorld()
	sharer := NewRuntimeLimitSharer()
	counter := 0

	const N = 3
	const LOOPS = 5
	const SLEEP = 16

	sharer.RegisterRunner("basic")
	insert := func(i int) {
		sharer.AddLogic("basic", &LogicUnit{
			name:    fmt.Sprintf("basic-%d", i),
			worldID: w.IdGen.Next(),
			f: func(dt_ms float64) {
				time.Sleep(SLEEP)
				counter += 1
			},
			active: true})
	}
	for i := 0; i < N; i++ {
		insert(i)
	}
	for i := 0; i < LOOPS; i++ {
		// insert with 3 loops left to go
		if i == LOOPS-3 {
			insert(N + i)
		}
		// ensure there's always enough time to run every one
		sharer.Share(5 * N * SLEEP)
	}
	Logger.Printf("Result: %d", counter)
	expected := N*(LOOPS-1) + (3 - 1)
	if counter != expected {
		t.Fatal("didn't share runtime properly")
	}
}

func TestRuntimeLimiterInsertAppending(t *testing.T) {
	r := NewRuntimeLimiter()
	for i := 0; i < 32; i++ {
		name := fmt.Sprintf("logic-%d", i)
		logic := &LogicUnit{
			name:    name,
			worldID: i,
			f:       func(dt_ms float64) {},
			active:  true}
		r.Add(logic)
		if !(len(r.logicUnits) > 0 &&
			r.indexes[logic.worldID] == len(r.logicUnits)-1) {
			t.Fatal("was not inserted properly")
		}
	}
}
