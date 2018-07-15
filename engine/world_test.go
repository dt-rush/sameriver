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

func TestWorldRunSystemsOnly(t *testing.T) {
	w := NewWorld(1024, 1024)
	ts := newTestSystem()
	w.AddSystems(ts)
	w.Update(FRAME_SLEEP_MS / 2)
	if ts.x == 0 {
		t.Fatal("failed to update system")
	}
}

func TestWorldRunWorldLogicsOnly(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	w.AddLogic("logic", func() { x += 1 })
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 1 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldRunSystemsAndLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	ts := newTestSystem()
	w.AddSystems(ts)
	x := 0
	w.AddLogic("logic", func() { x += 1 })
	w.Update(FRAME_SLEEP_MS / 2)
	if ts.x == 0 {
		t.Fatal("failed to update system")
	}
	if x != 1 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldActivateDeactivateLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name1 := "l1"
	w.AddLogic(name1, func() { x += 1 })
	// test Activate
	w.ActivateLogic(name1)
	if !w.worldLogicsRunner.byName[name1].Active {
		t.Fatal("failed to activate logic")
	}
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 1 {
		t.Fatal("active logic didn't run")
	}
	// test Deactivate
	x = 0
	w.DeactivateLogic(name1)
	if w.worldLogicsRunner.byName[name1].Active {
		t.Fatal("failed to deactivate logic")
	}
	if x != 0 {
		t.Fatal("deactivated logic ran")
	}
	// test ActivateAll/DeactivateAll
	name2 := "l2"
	w.AddLogic(name2, func() {})
	w.ActivateLogic(name2)
	w.ActivateLogic(name2)
	w.DeactivateAllLogics()
	for _, l := range w.worldLogicsRunner.logicUnits {
		if l.Active {
			t.Fatal("did not deactivate all logic")
		}
	}
	w.ActivateAllLogics()
	for _, l := range w.worldLogicsRunner.logicUnits {
		if !l.Active {
			t.Fatal("did not activate all logic")
		}
	}
}

func TestWorldRunLogicOverrun(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.AddLogic("logic", func() { time.Sleep(150 * time.Millisecond) })
	w.ActivateAllLogics()
	w.worldLogicsRunner.Start()
	remaining_ms := w.worldLogicsRunner.Run(100)
	if remaining_ms > 0 {
		t.Fatal("overrun time not calculated properly")
	}
	if !w.worldLogicsRunner.Finished() {
		t.Fatal("should have returned finished = true when running sole logic " +
			"within time limit")
	}
}

func TestWorldLogicUnderrun(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.AddLogic("logic", func() { time.Sleep(100 * time.Millisecond) })
	w.ActivateAllLogics()
	w.worldLogicsRunner.Start()
	remaining_ms := w.worldLogicsRunner.Run(300)
	if !(remaining_ms >= 0 && remaining_ms <= 200) {
		t.Fatal("remaining time not calculated properly")
	}
}

func TestWorldLogicTimeLimit(t *testing.T) {
	w := NewWorld(1024, 1024)
	w.AddLogic("slow", func() { time.Sleep(10 * time.Millisecond) })
	fastRan := false
	w.AddLogic("fast", func() { fastRan = true })
	w.ActivateAllLogics()
	w.worldLogicsRunner.Start()
	w.worldLogicsRunner.Run(2)
	if fastRan {
		t.Fatal("continued running logic despite using up allowed milliseconds")
	}
}
