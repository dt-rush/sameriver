package engine

func testingWorld() *World {
	return NewWorld(1024, 1024)
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
	w.SetPrimaryEntityLogic(e, func() { entityUpdates += 1 })
	w.ActivateEntityLogic(e)
	return w, ts, &worldUpdates, &entityUpdates
}
