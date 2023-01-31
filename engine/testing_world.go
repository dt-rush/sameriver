package engine

func testingWorld() *World {
	w := NewWorld(1024, 1024)
	w.RegisterComponents([]string{
		"Vec2D,Position",
	})
	return w
}

func testingWorldWithAllLogicTypes() (*World, *testSystem, *int, *int) {
	w := testingWorld()
	// add system
	ts := newTestSystem()
	w.AddSystems(ts)
	// add world logic
	worldUpdates := 0
	entityUpdates := 0
	name := "logic"
	w.AddWorldLogic(name, func() { worldUpdates += 1 })
	w.ActivateWorldLogic(name)
	// add entity logic
	e, _ := testingSpawnSimple(w)
	e.AddLogic("incrementer", func() { entityUpdates += 1 })
	return w, ts, &worldUpdates, &entityUpdates
}
