package main

import (
	"github.com/dt-rush/sameriver/engine"
)

func simpleSpawnRequestData() engine.SpawnRequestData {
	return engine.SpawnRequestData{
		Components: engine.ComponentSet{
			Position: &engine.Vec2D{0, 0},
		},
	}
}

func simpleTaggedSpawnRequestData() engine.SpawnRequestData {
	return engine.SpawnRequestData{
		Components: engine.ComponentSet{
			Position: &engine.Vec2D{0, 0},
			TagList:  &engine.TagList{Tags: []string{"tag1"}},
		},
	}
}
