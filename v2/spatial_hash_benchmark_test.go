package sameriver

import (
	"math"
	"testing"
)

func BenchmarkSpatialHashEntitiesWithinDistance(b *testing.B) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})

	e := testingSpawnSpatial(w, Vec2D{50, 50}, Vec2D{5, 5})

	near := make([]*Entity, 0)
	far := make([]*Entity, 0)
	for spawnRadius := 30.0; spawnRadius <= 38; spawnRadius += 8 {
		for i := 0.0; i < 720; i += 10 {
			theta := 2.0 * math.Pi * (i / 360)
			offset := Vec2D{
				spawnRadius * math.Cos(theta),
				spawnRadius * math.Sin(theta),
			}
			spawned := testingSpawnSpatial(w,
				e.GetVec2D("Position").Add(offset),
				Vec2D{5, 5})
			if spawnRadius == 30.0 {
				near = append(near, spawned)
			} else {
				far = append(far, spawned)
			}
		}
	}
	for i := 0; i < b.N; i++ {
		w.Update(FRAME_DURATION_INT / 2)
		w.EntitiesWithinDistance(
			*e.GetVec2D("Position"),
			*e.GetVec2D("Box"),
			30.0)
	}
}
