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
