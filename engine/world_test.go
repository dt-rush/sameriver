package engine

import (
	"fmt"
	"testing"
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

func TestWorldAddSystemsDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := NewWorld(1024, 1024)
	w.AddSystems(newTestSystem(), newTestSystem())
	t.Fatal("should have panic'd")
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

func TestWorldAddWorldLogicDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := NewWorld(1024, 1024)
	w.AddWorldLogic("world-logic", func() {})
	w.AddWorldLogic("world-logic", func() {})
	t.Fatal("should have panic'd")
}

func TestWorldRunWorldLogicsOnly(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	w.AddWorldLogic("logic", func() { x += 1 })
	w.ActivateAllWorldLogics()
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 1 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldAddEntityLogicDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := NewWorld(1024, 1024)
	w.AddEntityLogic("entity-logic", func() {})
	w.AddEntityLogic("entity-logic", func() {})
	t.Fatal("should have panic'd")
}

func TestWorldRunEntityLogicsOnly(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	w.AddEntityLogic("logic", func() { x += 1 })
	w.ActivateAllEntityLogics()
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 1 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldRunSystemsAndWorldLogicsAndEntityLogics(t *testing.T) {
	w := NewWorld(1024, 1024)
	ts := newTestSystem()
	w.AddSystems(ts)
	x := 0
	y := 0
	name := "logic"
	w.AddWorldLogic(name, func() { x += 1 })
	w.ActivateWorldLogic(name)
	w.AddEntityLogic(name, func() { y += 1 })
	w.ActivateEntityLogic(name)
	w.Update(FRAME_SLEEP_MS / 2)
	if ts.x != 1 {
		t.Fatal("failed to update system")
	}
	if x != 1 {
		t.Fatal("failed to run world logic")
	}
	if y != 1 {
		t.Fatal("failed to run entity logic")
	}
}

func TestWorldRemoveWorldLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name := "l1"
	w.AddWorldLogic(name, func() { x += 1 })
	w.RemoveWorldLogic(name)
	for i := 0; i < 32; i++ {
		w.Update(FRAME_SLEEP_MS)
	}
	if x != 0 {
		t.Fatal("logic was removed but still ran during Update()")
	}
}

func TestWorldRemoveEntityLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name := "l1"
	w.AddEntityLogic(name, func() { x += 1 })
	w.RemoveEntityLogic(name)
	for i := 0; i < 32; i++ {
		w.Update(FRAME_SLEEP_MS)
	}
	if x != 0 {
		t.Fatal("logic was removed but still ran during Update()")
	}
}

func TestWorldActivateWorldLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name := "l1"
	w.AddWorldLogic(name, func() { x += 1 })
	w.ActivateWorldLogic(name)
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 1 {
		t.Fatal("logic should have been active and run - did not")
	}
}

func TestWorldActivateAllWorldLogics(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("logic-%d", i)
		w.AddWorldLogic(name, func() { x += 1 })
	}
	w.ActivateAllWorldLogics()
	w.Update(FRAME_SLEEP_MS / 2)
	if x != n {
		t.Fatal("logics all should have been activated - some did not run")
	}
}

func TestWorldDeactivateWorldLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name1 := "l1"
	w.AddWorldLogic(name1, func() { x += 1 })
	w.DeactivateWorldLogic(name1)
	if x != 0 {
		t.Fatal("deactivated logic ran")
	}
}

func TestWorldDeativateAllWorldLogics(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("logic-%d", i)
		w.AddWorldLogic(name, func() { x += 1 })
		w.ActivateWorldLogic(name)
	}
	w.DeactivateAllWorldLogics()
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 0 {
		t.Fatal("logics all should have been deactivated, but some ran")
	}
}

func TestWorldActivateEntityLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name := "l1"
	w.AddEntityLogic(name, func() { x += 1 })
	w.ActivateEntityLogic(name)
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 1 {
		t.Fatal("logic should have been active and run - did not")
	}
}

func TestWorldActivateAllEntityLogics(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("logic-%d", i)
		w.AddEntityLogic(name, func() { x += 1 })
	}
	w.ActivateAllEntityLogics()
	w.Update(FRAME_SLEEP_MS / 2)
	if x != n {
		t.Fatal("logics all should have been activated - some did not run")
	}
}

func TestWorldDeactivateEntityLogic(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	name1 := "l1"
	w.AddEntityLogic(name1, func() { x += 1 })
	w.DeactivateEntityLogic(name1)
	if x != 0 {
		t.Fatal("deactivated logic ran")
	}
}

func TestWorldDeativateAllEntityLogics(t *testing.T) {
	w := NewWorld(1024, 1024)
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("logic-%d", i)
		w.AddEntityLogic(name, func() { x += 1 })
		w.ActivateEntityLogic(name)
	}
	w.DeactivateAllEntityLogics()
	w.Update(FRAME_SLEEP_MS / 2)
	if x != 0 {
		t.Fatal("logics all should have been deactivated, but some ran")
	}
}
