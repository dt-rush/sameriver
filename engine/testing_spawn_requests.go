package engine

func testingSimpleSpawnRequest() SpawnRequestData {
	return SpawnRequestData{Tags: []string{}, Components: ComponentSet{}}
}

func testingQueueSpawnSimple(em EntityManagerInterface) {
	em.QueueSpawn([]string{}, ComponentSet{})
}

func testingQueueSpawnUnique(em EntityManagerInterface) {
	em.QueueSpawnUnique("unique", []string{}, ComponentSet{})
}

func testingSpawnSimple(em EntityManagerInterface) (*Entity, error) {
	return em.Spawn([]string{}, ComponentSet{})
}

func testingSpawnPosition(
	em EntityManagerInterface, pos Vec2D) (*Entity, error) {
	return em.Spawn([]string{}, ComponentSet{Position: &pos})
}

func testingSpawnTagged(
	em EntityManagerInterface, tag string) (*Entity, error) {
	return em.Spawn([]string{tag}, ComponentSet{})
}

func testingSpawnSpatial(
	em EntityManagerInterface, pos Vec2D, box Vec2D) (*Entity, error) {
	return em.Spawn([]string{},
		ComponentSet{
			Position: &pos,
			Box:      &box,
		},
	)
}

func testingSpawnCollision(em EntityManagerInterface) (*Entity, error) {
	return em.Spawn([]string{},
		ComponentSet{
			Position: &Vec2D{10, 10},
			Box:      &Vec2D{4, 4},
		},
	)
}

func testingSpawnSteering(em EntityManagerInterface) (*Entity, error) {
	mass := 3.0
	maxV := 3.0
	return em.Spawn([]string{},
		ComponentSet{
			Position:       &Vec2D{0, 0},
			Velocity:       &Vec2D{0, 0},
			MaxVelocity:    &maxV,
			MovementTarget: &Vec2D{1, 1},
			Steer:          &Vec2D{0, 0},
			Mass:           &mass,
		},
	)
}

func testingSpawnPhysics(em EntityManagerInterface) (*Entity, error) {
	mass := 3.0
	return em.Spawn([]string{},
		ComponentSet{
			Position: &Vec2D{10, 10},
			Velocity: &Vec2D{0, 0},
			Box:      &Vec2D{1, 1},
			Mass:     &mass,
		},
	)
}
