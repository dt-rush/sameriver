package engine

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCanConstructWorld(t *testing.T) {
	w := testingWorld()
	if w == nil {
		t.Fatal("NewWorld() was nil")
	}
}

func TestWorldRegisterSystems(t *testing.T) {
	w := testingWorld()
	w.RegisterSystems(newTestSystem())
}

func TestWorldRegisterSystemsDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := testingWorld()
	w.RegisterSystems(newTestSystem(), newTestSystem())
	t.Fatal("should have panic'd")
}

func TestWorldAddDependentSystems(t *testing.T) {
	w := testingWorld()
	dep := newTestDependentSystem()
	w.RegisterSystems(
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
	w := testingWorld()
	w.RegisterSystems(
		newTestDependentSystem(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldNonPointerReceiverSystem(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := testingWorld()
	w.RegisterSystems(
		newTestNonPointerReceiverSystem(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldMisnamedSystem(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := testingWorld()
	w.RegisterSystems(
		newTestSystemThatIsMisnamed(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldSystemDependencyNonPointer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := testingWorld()
	w.RegisterSystems(
		newTestDependentNonPointerSystem(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldSystemDependencyNonSystem(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := testingWorld()
	w.RegisterSystems(
		newTestDependentNonSystemSystem(),
	)
	t.Fatal("should have panic'd")
}

func TestWorldRunSystemsOnly(t *testing.T) {
	w := testingWorld()
	ts := newTestSystem()
	w.RegisterSystems(ts)
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	if ts.updates == 0 {
		t.Fatal("failed to update system")
	}
}

func TestWorldAddWorldLogicDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	w := testingWorld()
	w.AddWorldLogic("world-logic", func(dt_ms float64) {})
	w.AddWorldLogic("world-logic", func(dt_ms float64) {})
	t.Fatal("should have panic'd")
}

func TestWorldRunWorldLogicsOnly(t *testing.T) {
	w := testingWorld()
	x := 0
	w.AddWorldLogic("logic", func(dt_ms float64) { x += 1 })
	w.ActivateAllWorldLogics()
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	if x != 1 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldRunEntityLogicsOnly(t *testing.T) {
	w := testingWorld()
	x := 0
	e, _ := testingSpawnSimple(w)
	e.AddLogic("incrementer", func(dt_ms float64) { x += 1 })
	w.ActivateAllEntityLogics()
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	if x != 1 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldRunAllLogicTypes(t *testing.T) {
	w, ts, worldUpdates, entityUpdates := testingWorldWithAllLogicTypes()
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	if ts.updates != 1 {
		t.Fatal("failed to update system")
	}
	if *worldUpdates != 1 {
		t.Fatal("failed to run world logic")
	}
	if *entityUpdates != 1 {
		t.Fatal("failed to run entity logic")
	}
}

func TestWorldRemoveWorldLogic(t *testing.T) {
	w := testingWorld()
	x := 0
	name := "l1"
	w.AddWorldLogic(name, func(dt_ms float64) { x += 1 })
	w.RemoveWorldLogic(name)
	w.Update(1)
	for i := 0; i < 32; i++ {
		w.Update(FRAME_DURATION_INT)
	}
	if x != 0 {
		t.Fatal("logic was removed but still ran during Update()")
	}
}

func TestWorldRemoveEntityLogic(t *testing.T) {
	w := testingWorld()
	x := 0
	e, _ := testingSpawnSimple(w)
	e.AddLogic("incrementer", func(dt_ms float64) { x += 1 })
	e.RemoveLogic("incrementer")
	w.Update(1)
	for i := 0; i < 32; i++ {
		w.Update(FRAME_DURATION_INT)
	}
	if x != 0 {
		t.Fatal("logic was removed but still ran during Update()")
	}
}

func TestWorldActivateWorldLogic(t *testing.T) {
	w := testingWorld()
	x := 0
	name := "l1"
	w.AddWorldLogic(name, func(dt_ms float64) { x += 1 })
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	if x != 1 {
		t.Fatal("logic should have been active and run - did not")
	}
}

func TestWorldActivateAllWorldLogics(t *testing.T) {
	w := testingWorld()
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("logic-%d", i)
		w.AddWorldLogic(name, func(dt_ms float64) { x += 1 })
	}
	w.ActivateAllWorldLogics()
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	if x != n {
		t.Fatal("logics all should have been activated - some did not run")
	}
}

func TestWorldDeactivateWorldLogic(t *testing.T) {
	w := testingWorld()
	x := 0
	name1 := "l1"
	w.AddWorldLogic(name1, func(dt_ms float64) { x += 1 })
	w.DeactivateWorldLogic(name1)
	if x != 0 {
		t.Fatal("deactivated logic ran")
	}
}

func TestWorldDeativateAllWorldLogics(t *testing.T) {
	w := testingWorld()
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("logic-%d", i)
		w.AddWorldLogic(name, func(dt_ms float64) { x += 1 })
	}
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	w.DeactivateAllWorldLogics()
	w.Update(FRAME_DURATION_INT / 2)
	Logger.Println(x)
	if x != 16 {
		t.Fatal("logics all should have been deactivated, but some ran")
	}
}

func TestWorldEntityLogicActiveDefault(t *testing.T) {
	w := testingWorld()
	x := 0
	e, _ := testingSpawnSimple(w)
	e.AddLogic("incrementer", func(dt_ms float64) { x += 1 })
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	if x != 1 {
		t.Fatal("logic should have been active and run - did not")
	}
}

func TestWorldActivateAllEntityLogics(t *testing.T) {
	w := testingWorld()
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		e, _ := testingSpawnSimple(w)
		e.AddLogic("incrementer", func(dt_ms float64) { x += 1 })
	}
	w.ActivateAllEntityLogics()
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	if x != n {
		t.Fatal("logics all should have been activated - some did not run")
	}
}

func TestWorldDeactivateEntityLogic(t *testing.T) {
	w := testingWorld()
	x := 0
	e, _ := testingSpawnSimple(w)
	e.AddLogic("incrementer", func(dt_ms float64) { x += 1 })
	w.DeactivateEntityLogics(e)
	if x != 0 {
		t.Fatal("deactivated logic ran")
	}
}

func TestWorldDeativateAllEntityLogics(t *testing.T) {
	w := testingWorld()
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		e, _ := testingSpawnSimple(w)
		e.AddLogic("incrementer", func(dt_ms float64) { x += 1 })
	}
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	w.DeactivateAllEntityLogics()
	w.Update(FRAME_DURATION_INT / 2)
	if x != 16 {
		t.Fatal("logics all should have been deactivated, but some ran")
	}
}

func TestWorldDumpStats(t *testing.T) {
	w, _, _, _ := testingWorldWithAllLogicTypes()
	w.Update(1)
	w.Update(FRAME_DURATION_INT / 2)
	// test dump stats object
	stats := w.DumpStats()
	if len(stats) <= 1 {
		t.Fatal("stats not populated properly - should have keys for each " +
			"subsystem at least")
	}
	if _, ok := stats["totals"]["total"]; !ok {
		t.Fatal("no total update runtime stat included")
	}
	// test whether individual runtime limiter DumpStats() corresponds to
	// their entries in the overall DumpStats()
	systemStats, _ := w.runtimeSharer.runnerMap["systems"].DumpStats()
	worldStats, _ := w.runtimeSharer.runnerMap["world"].DumpStats()
	entityStats, _ := w.runtimeSharer.runnerMap["entities"].DumpStats()

	if !reflect.DeepEqual(stats["systems"], systemStats) {
		t.Fatal("system stats dump was not equal to systemsRunner stats dump")
	}
	if !reflect.DeepEqual(stats["world"], worldStats) {
		t.Fatal("world stats dump was not equal to worldLogicsRunner stats dump")
	}
	if !reflect.DeepEqual(stats["entities"], entityStats) {
		t.Fatal("primary entity stats dump was not equal to primaryEntityLogicsRunner stats dump")
	}
	// test stats string
	statsString := w.DumpStatsString()
	if statsString == "" || len(statsString) == 0 {
		t.Fatal("statsString is empty")
	}

}
