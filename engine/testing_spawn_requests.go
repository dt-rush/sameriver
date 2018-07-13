package engine

func simpleSpawnRequestData() SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{},
	}
}

func positionSpawnRequestData(pos Vec2D) SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &pos,
		},
	}
}

func spatialSpawnRequestData(pos Vec2D, box Vec2D) SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &pos,
			Box:      &box,
		},
	}
}

func simpleTaggedSpawnRequestData(tag string) SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &Vec2D{0, 0},
		},
		Tags: []string{tag},
	}
}

func collisionSpawnRequestData() SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &Vec2D{0, 0},
			Box:      &Vec2D{10, 10},
		},
	}
}

func steeringSpawnRequestData() SpawnRequestData {
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

func physicsSpawnRequestData() SpawnRequestData {
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
