package sameriver

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
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
		if r := recover(); r == nil {
			t.Fatal("Should have panic'd")
		}
	}()
	w := testingWorld()
	w.RegisterSystems(newTestSystem(), newTestSystem())
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
		if r := recover(); r == nil {
			t.Fatal("Should have panic'd")
		}
	}()
	w := testingWorld()
	w.RegisterSystems(
		newTestDependentSystem(),
	)
}

func TestWorldNonPointerReceiverSystem(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Should have panic'd")
		}
	}()
	w := testingWorld()
	w.RegisterSystems(
		newTestNonPointerReceiverSystem(),
	)
}

func TestWorldMisnamedSystem(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should have panic'd")
		}
	}()
	w := testingWorld()
	w.RegisterSystems(
		newTestSystemThatIsMisnamed(),
	)
}

func TestWorldSystemDependencyNonPointer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should have panic'd")
		}
	}()
	w := testingWorld()
	s := newTestDependentNonPointerSystem()
	Logger.Println(s.ts)
	w.RegisterSystems(s)
}

func TestWorldSystemDependencyNonSystem(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should have panic'd")
		}
	}()
	w := testingWorld()
	s := newTestDependentNonSystemSystem()
	Logger.Println(s.ts)
	w.RegisterSystems(s)
}

func TestWorldRunSystemsOnly(t *testing.T) {
	w := testingWorld()
	ts := newTestSystem()
	w.RegisterSystems(ts)
	w.Update(FRAME_MS / 2)
	if ts.updates == 0 {
		t.Fatal("failed to update system")
	}
}

func TestWorldAddWorldLogicDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should have panic'd")
		}
	}()
	w := testingWorld()
	w.AddWorldLogic("world-logic", func(dt_ms float64) {})
	w.AddWorldLogic("world-logic", func(dt_ms float64) {})
}

func TestWorldRunWorldLogicsOnly(t *testing.T) {
	w := testingWorld()
	x := 0
	w.AddWorldLogic("logic", func(dt_ms float64) { x += 1 })
	w.ActivateAllWorldLogics()
	w.Update(FRAME_MS / 2)
	if x == 0 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldRunEntityLogicsOnly(t *testing.T) {
	w := testingWorld()
	x := 0
	e := testingSpawnSimple(w)
	e.AddLogic("incrementer", func(e *Entity, dt_ms float64) { x += 1 })
	w.ActivateAllEntityLogics()
	w.Update(FRAME_MS / 2)
	if x == 0 {
		t.Fatal("failed to run logic")
	}
}

func TestWorldRunAllLogicTypes(t *testing.T) {
	w, ts, worldUpdates, entityUpdates := testingWorldWithAllLogicTypes()
	w.Update(FRAME_MS / 2)
	if ts.updates == 0 {
		t.Fatal("failed to update system")
	}
	if *worldUpdates == 0 {
		t.Fatal("failed to run world logic")
	}
	if *entityUpdates == 0 {
		t.Fatal("failed to run entity logic")
	}
}

func TestWorldRemoveWorldLogic(t *testing.T) {
	w := testingWorld()
	x := 0
	name := "l1"
	w.AddWorldLogic(name, func(dt_ms float64) { x += 1 })
	w.RemoveWorldLogic(name)
	for i := 0; i < 32; i++ {
		w.Update(FRAME_MS)
	}
	if x != 0 {
		t.Fatal("logic was removed but still ran during Update()")
	}
}

func TestWorldRemoveEntityLogic(t *testing.T) {
	w := testingWorld()
	x := 0
	e := testingSpawnSimple(w)
	Logger.Println("adding incrementer")
	e.AddLogic("incrementer", func(e *Entity, dt_ms float64) { x += 1 })
	Logger.Println("removing incrementer")
	e.RemoveLogic("incrementer")
	for i := 0; i < 32; i++ {
		w.Update(FRAME_MS)
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
	w.Update(FRAME_MS / 2)
	if x == 0 {
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
	w.Update(FRAME_MS / 2)
	if x < n {
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
	w.Update(FRAME_MS / 2)
	w.DeactivateAllWorldLogics()
	w.Update(FRAME_MS / 2)
	Logger.Println(x)
	if x < n {
		t.Fatal("logics all should have been deactivated, but some ran")
	}
}

func TestWorldEntityLogicActiveDefault(t *testing.T) {
	w := testingWorld()
	x := 0
	e := testingSpawnSimple(w)
	e.AddLogic("incrementer", func(e *Entity, dt_ms float64) { x += 1 })
	w.Update(FRAME_MS / 2)
	if x == 0 {
		t.Fatal("logic should have been active and run - did not")
	}
}

func TestWorldActivateAllEntityLogics(t *testing.T) {
	w := testingWorld()
	x := 0
	n := 16
	for i := 0; i < n; i++ {
		e := testingSpawnSimple(w)
		e.AddLogic("incrementer", func(e *Entity, dt_ms float64) { x += 1 })
	}
	w.ActivateAllEntityLogics()
	w.Update(FRAME_MS / 2)
	if x < n {
		t.Fatal("logics all should have been activated - some did not run")
	}
}

func TestWorldDeactivateEntityLogic(t *testing.T) {
	w := testingWorld()
	x := 0
	e := testingSpawnSimple(w)
	e.AddLogic("incrementer", func(e *Entity, dt_ms float64) { x += 1 })
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
		e := testingSpawnSimple(w)
		e.AddLogic("incrementer", func(e *Entity, dt_ms float64) { x += 1 })
	}
	w.Update(FRAME_MS / 2)
	w.DeactivateAllEntityLogics()
	w.Update(FRAME_MS / 2)
	if x < n {
		t.Fatal("logics all should have been deactivated, but some ran")
	}
}

func TestWorldDumpStats(t *testing.T) {
	w, _, _, _ := testingWorldWithAllLogicTypes()
	w.Update(FRAME_MS / 2)
	// test dump stats object
	stats := w.DumpStats()
	if len(stats) <= 1 {
		t.Fatal("stats not populated properly - should have keys for each " +
			"subsystem at least")
	}
	if _, ok := stats["__totals"]["World.Update()"]; !ok {
		t.Fatal("no total update runtime stat included")
	}
	// test whether individual runtime limiter DumpStats() corresponds to
	// their entries in the overall DumpStats()
	systemStats, _ := w.RuntimeSharer.RunnerMap["systems"].DumpStats()
	worldStats, _ := w.RuntimeSharer.RunnerMap["world"].DumpStats()
	entityStats, _ := w.RuntimeSharer.RunnerMap["entities"].DumpStats()

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

func TestWorldPredicateEntities(t *testing.T) {
	w := NewWorld(map[string]any{
		"width":               100,
		"height":              100,
		"distanceHasherGridX": 10,
		"distanceHasherGridY": 10,
	})

	e := testingSpawnSpatial(w, Vec2D{50, 50}, Vec2D{5, 5})
	w.TagEntity(e, "tree")

	near := make([]*Entity, 0)
	nTrees := 0
	for spawnRadius := 30.0; spawnRadius <= 38; spawnRadius += 8 {
		for i := 0.0; i < 360; i += 10 {
			theta := 2.0 * math.Pi * (i / 360)
			offset := Vec2D{
				spawnRadius * math.Cos(theta),
				spawnRadius * math.Sin(theta),
			}
			spawned := testingSpawnSpatial(w,
				e.GetVec2D(POSITION).Add(offset),
				Vec2D{5, 5})
			if int(i)%20 == 0 {
				w.TagEntity(spawned, "tree")
				nTrees++
			}
			if spawnRadius == 30.0 {
				near = append(near, spawned)
			}
		}
	}
	isTree := func(e *Entity) bool {
		return e.HasTag("tree")
	}

	w.Update(FRAME_MS / 2)
	nearGot := w.EntitiesWithinDistance(
		*e.GetVec2D(POSITION),
		*e.GetVec2D(BOX),
		30.0)
	if len(nearGot) != len(near)+1 {
		t.Fatalf("Should be %d near entities; got %d", len(near)+1, len(nearGot))
	}

	treesFound := w.EntitiesWithinDistanceFilter(
		*e.GetVec2D(POSITION),
		*e.GetVec2D(BOX),
		30.0,
		isTree)
	if len(treesFound) != 19 {
		t.Fatalf("Should have found 19 near trees (1 original, 18 radial); got %d", len(treesFound))
	}
	allTrees := w.FilterAllEntities(isTree)
	if len(allTrees) != 37 {
		t.Fatal("Should have found 37 trees in the world (1 original, 18 radial x 2 layers")
	}
}

func TestWorldSetTimeout(t *testing.T) {
	w := testingWorld()
	x := 0
	w.SetTimeout(func() {
		x++
	}, 500)
	for i := 0; i < 516/FRAME_MS; i++ {
		w.Update(FRAME_MS)
		time.Sleep(FRAME_DURATION)
	}
	if x != 1 {
		t.Fatalf("Should've run settimeout func 1 time, ran %d times", x)
	}
}

func TestWorldSetInterval(t *testing.T) {
	w := testingWorld()
	x := 0
	w.SetInterval(func() {
		Logger.Println("run")
		x++
	}, 100)
	t0 := time.Now()
	// iterate for however many frames fit in 516 ms
	for i := 0; i < 516/FRAME_MS; i++ {
		t1 := time.Now()
		w.Update(FRAME_MS)
		frame_ms := float64(time.Since(t1).Nanoseconds()) / 1e6
		Logger.Printf("%f ms frame", frame_ms)
		time.Sleep(time.Duration(math.Max(0, float64(FRAME_MS)-frame_ms)*1e6) * time.Nanosecond)
	}
	elapsed := float64(time.Since(t0).Nanoseconds()) / 1e6
	Logger.Printf("elapsed: %f ms", elapsed)
	timesRun := int(elapsed / 100)
	if x != timesRun {
		t.Fatalf("Should've run setinterval func %d times, ran %d times", timesRun, x)
	}
}

func TestWorldSetNInterval(t *testing.T) {
	w := testingWorld()
	x := 0
	w.SetNInterval(func() {
		Logger.Println("run")
		x++
	}, 100, 3)
	for i := 0; i < 516/FRAME_MS; i++ {
		w.Update(FRAME_MS)
		time.Sleep(FRAME_DURATION)
	}
	if x != 3 {
		t.Fatalf("Should've run setninterval func 3 times, ran %d times", x)
	}
}
