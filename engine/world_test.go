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

func TestWorldAddSystems(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.AddSystems(newTestSystem())
}

func TestWorldAddDependentSystems(t *testing.T) {
	w := NewWorld(1024, 1024)
	dep := newTestDependentSystem()
	w.AddSystems(
		newTestSystem(),
		dep,
	)
	if dep.ts == nil {
		t.Fatal("system dependency not injected")
	}
}

func TestWorldUnresolvedSystemDependency(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := NewWorld(1024, 1024)
	w.AddSystems(
		newTestDependentSystem(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldNonPointerReceiverSystem(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := NewWorld(1024, 1024)
	w.AddSystems(
		newTestNonPointerReceiverSystem(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldMisnamedSystem(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := NewWorld(1024, 1024)
	w.AddSystems(
		newTestSystemThatIsMisnamed(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldSystemDependencyNonPointer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := NewWorld(1024, 1024)
	w.AddSystems(
		newTestDependentNonPointerSystem(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldSystemDependencyNonSystem(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := NewWorld(1024, 1024)
	w.AddSystems(
		newTestDependentNonSystemSystem(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldUpdate(t *testing.T) {
	w := NewWorld(1024, 1024)
	ts := newTestSystem()
	w.AddSystems(ts)
	w.Update(FRAME_SLEEP_MS)
	if ts.x != FRAME_SLEEP_MS {
		t.Fatal("didn't update world.systems")
	}
}

func TestWorldActivateDeactivateLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	lu1 := LogicUnit{
		Name:   "l1",
		Active: false,
		F:      func() { x += 1 },
	}
	w.AddLogic(lu1)
	// test Activate
	w.ActivateLogic("l1")
	if !w.logics[0].Active {
		t.Fatal("failed to activate logic")
	}
	w.RunLogic(FRAME_SLEEP_MS)
	if x != 1 {
		t.Fatal("active logic didn't run")
	}
	// test Deactivate
	x = 0
	w.DeactivateLogic("l1")
	if w.logics[0].Active {
		t.Fatal("failed to deactivate logic")
	}
	if x != 0 {
		t.Fatal("deactivated logic ran")
	}
	// test ActivateAll/DeactivateAll
	lu2 := LogicUnit{
		Name:   "l2",
		Active: false,
		F:      func() {},
	}
	w.AddLogic(lu2)
	w.ActivateLogic("l1")
	w.ActivateLogic("l2")
	w.DeactivateAllLogic()
	for _, l := range w.logics {
		if l.Active {
			t.Fatal("did not deactivate all logic")
		}
	}
	w.ActivateAllLogic()
	for _, l := range w.logics {
		if !l.Active {
			t.Fatal("did not activate all logic")
		}
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
			F:      func() { time.Sleep(100 * time.Millisecond) },
		})
	}
	overrun_ms := w.RunLogic(150)
	if !(overrun_ms >= 50) {
		t.Fatal("overrun time not calculated")
	}
	overrun_ms = w.RunLogic(300)
	if !(overrun_ms >= 0 && overrun_ms <= 50) {
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
