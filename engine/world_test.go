package engine

import (
	"testing"
	"time"
)

func TestCanConstructWorld(t *testing.T) {
	w := NewWorld(1024, 1024)
	if w == nil {
		t.Fatal("NewWorld() was nil")
	}
}

type testSystem struct {
	x float64
}

func newTestSystem() *testSystem {
	return &testSystem{}
}
func (s *testSystem) LinkWorld(w *World) {}
func (s *testSystem) Update(dt_ms float64) {
	s.x += dt_ms
}

func TestWorldAddSystem(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.AddSystem(newTestSystem())
}

func TestWorldUpdate(t *testing.T) {
	w := NewWorld(1024, 1024)
	ts := newTestSystem()
	w.AddSystem(ts)
	w.Update(FRAME_SLEEP_MS)
	if ts.x != FRAME_SLEEP_MS {
		t.Fatal("didn't update world.systems")
	}
}

func TestWorldActivateDeactivateLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	w.AddLogic(LogicUnit{
		Name:   "logic",
		Active: false,
		F:      func() { x += 1 },
	})
	// test Activate
	w.ActivateLogic("logic")
	if !w.logics[0].Active {
		t.Fatal("failed to activate logic")
	}
	w.RunLogic(FRAME_SLEEP_MS)
	if x != 1 {
		t.Fatal("active logic didn't run")
	}
	// test Deactivate
	x = 0
	w.DeactivateLogic("logic")
	if w.logics[0].Active {
		t.Fatal("failed to deactivate logic")
	}
	if x != 0 {
		t.Fatal("deactivated logic ran")
	}
}

func TestWorldRunLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	w.AddLogic(LogicUnit{
		Name:   "logic",
		Active: true,
		F:      func() { x += 1 },
	})
	w.RunLogic(FRAME_SLEEP_MS / 5)
	if x != 1 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldRunLogicTiming(t *testing.T) {
	w := NewWorld(1024, 1024)
	for i := 0; i < 3; i++ {
		w.AddLogic(LogicUnit{
			Name:   "logic",
			Active: true,
			F:      func() { time.Sleep(4 * time.Millisecond) },
		})
	}
	overrun_ms := w.RunLogic(2)
	if overrun_ms != 2 {
		t.Fatal("overrun time not calculated")
	}
	overrun_ms = w.RunLogic(12)
	if overrun_ms != 0 {
		t.Fatal("overrun time not calculated")
	}
	w = NewWorld(1024, 1024)
	w.AddLogic(LogicUnit{
		Name:   "slow",
		Active: true,
		F:      func() { time.Sleep(10 * time.Millisecond) },
	})
	fastRan := false
	w.AddLogic(LogicUnit{
		Name:   "fast",
		Active: true,
		F:      func() { fastRan = true },
	})
	w.RunLogic(8)
	if fastRan {
		t.Fatal("continued running logic despite using up allowed milliseconds")
	}
}
