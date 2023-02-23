package sameriver

import (
	"math/rand"

	"testing"
)

func BenchmarkEntityManagerSpawnDespawn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewWorld(100, 100)
		w.RegisterComponents("Vec2D,Position", "Vec2D,Box", "Vec2D,Velocity")
		for i := 0; i < MAX_ENTITIES; i++ {
			w.Spawn(map[string]any{
				"Vec2D,Position": Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
				"Vec2D,Box":      Vec2D{5, 5},
				"Vec2D,Velocity": Vec2D{rand.Float64(), rand.Float64()},
			})
		}
	}
}
