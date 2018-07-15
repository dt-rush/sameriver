package engine

import (
	"fmt"
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
	w.ActivateAllLogics()
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
	name := "logic"
	w.AddLogic(name, func() { x += 1 })
	w.ActivateLogic(name)
	w.Update(FRAME_SLEEP_MS / 2)
	if ts.x != 1 {
		t.Fatal("failed to update system")
	}
	if x != 1 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldRemoveLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name := "l1"
	w.AddLogic(name, func() { x += 1 })
	w.RemoveLogic(name)
	for i := 0; i < 32; i++ {
		w.Update(FRAME_SLEEP_MS)
	}
	if x != 0 {
		t.Fatal("logic was removed but still ran during Update()")
	}
}

func TestWorldActivateLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name := "l1"
	w.AddLogic(name, func() { x += 1 })
	w.ActivateLogic(name)
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 1 {
		t.Fatal("logic should have been active and run - did not")
	}
}

func TestWorldActivateAllLogics(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("logic-%d", i)
		w.AddLogic(name, func() { x += 1 })
	}
	w.ActivateAllLogics()
	w.Update(FRAME_SLEEP_MS / 2)
	if x != n {
		t.Fatal("logics all should have been activated - some did not run")
	}
}

func TestWorldDeactivateLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name1 := "l1"
	w.AddLogic(name1, func() { x += 1 })
	w.DeactivateLogic(name1)
	if x != 0 {
		t.Fatal("deactivated logic ran")
	}
}

func TestWorldDeativateAllLogics(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("logic-%d", i)
		w.AddLogic(name, func() { x += 1 })
		w.ActivateLogic(name)
	}
	w.DeactivateAllLogics()
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 0 {
		t.Fatal("logics all should have been deactivated, but some ran")
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
