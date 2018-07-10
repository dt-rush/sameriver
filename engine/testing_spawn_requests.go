package engine

func simpleSpawnRequestData() SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &Vec2D{0, 0},
		},
	}
}

func simpleTaggedSpawnRequestData() SpawnRequestData {
	return SpawnRequestData{
		Components: ComponentSet{
			Position: &Vec2D{0, 0},
			TagList:  &TagList{Tags: []string{"tag1"}},
		},
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
