package sameriver

import (
	"math/rand"

	"testing"
)

func BenchmarkEntityManagerSpawnDespawn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewWorld(map[string]any{
			"width":  100,
			"height": 100,
		})
		w.RegisterComponents([]any{
			VELOCITY, VEC2D, "VELOCITY",
		})
		for i := 0; i < MAX_ENTITIES; i++ {
			w.Spawn(map[string]any{
				"components": map[ComponentID]any{
					POSITION: Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
					BOX:      Vec2D{5, 5},
					VELOCITY: Vec2D{rand.Float64(), rand.Float64()},
				}})
		}
	}
}
