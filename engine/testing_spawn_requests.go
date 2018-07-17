package engine

func simpleSpawnRequest() SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{},
	}
}

func positionSpawnRequest(pos Vec2D) SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &pos,
		},
	}
}

func spatialSpawnRequest(pos Vec2D, box Vec2D) SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &pos,
			Box:      &box,
		},
	}
}

func simpleTaggedSpawnRequest(tag string) SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &Vec2D{0, 0},
		},
		Tags: []string{tag},
	}
}

func collisionSpawnRequest() SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &Vec2D{10, 10},
			Box:      &Vec2D{4, 4},
		},
	}
}

func steeringSpawnRequest() SpawnRequestData {
	mass := 3.0
	maxV := 3.0
	return SpawnRequestData{
		Components: ComponentSet{
			Position:       &Vec2D{0, 0},
			Velocity:       &Vec2D{0, 0},
			MaxVelocity:    &maxV,
			MovementTarget: &Vec2D{1, 1},
			Steer:          &Vec2D{0, 0},
			Mass:           &mass,
		},
	}
}

func physicsSpawnRequest() SpawnRequestData {
	mass := 3.0
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &Vec2D{0, 0},
			Velocity: &Vec2D{0, 0},
			Box:      &Vec2D{10, 10},
			Mass:     &mass,
		},
	}
}
