package sameriver

import (
	"math"
	"math/rand"
	"testing"
)

/*
func BenchmarkSpatialHashUpdateParallelD(b *testing.B) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	for i := 0; i < 1024; i++ {
		testingSpawnSpatial(w,
			Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			Vec2D{5, 5})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.SpatialHasher.ParallelUpdateD()
	}
}

func BenchmarkSpatialHashUpdateParallelC(b *testing.B) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	for i := 0; i < 1024; i++ {
		testingSpawnSpatial(w,
			Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			Vec2D{5, 5})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.SpatialHasher.ParallelUpdateB()
	}
}

func BenchmarkSpatialHashUpdateParallelB(b *testing.B) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	for i := 0; i < 1024; i++ {
		testingSpawnSpatial(w,
			Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			Vec2D{5, 5})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.SpatialHasher.ParallelUpdateB()
	}
}

func BenchmarkSpatialHashUpdateParallelA(b *testing.B) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	for i := 0; i < 1024; i++ {
		testingSpawnSpatial(w,
			Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			Vec2D{5, 5})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.SpatialHasher.ParallelUpdateA()
	}
}
*/

func BenchmarkSpatialHashUpdateSingleThread(b *testing.B) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})
	for i := 0; i < 1024; i++ {
		testingSpawnSpatial(w,
			Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			Vec2D{5, 5})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.SpatialHasher.SingleThreadUpdate()
	}
}

func BenchmarkSpatialHashEntitiesWithinDistance(b *testing.B) {
	w := NewWorld(map[string]any{
		"width":  100,
		"height": 100,
	})

	e := testingSpawnSpatial(w, Vec2D{50, 50}, Vec2D{5, 5})

	for spawnRadius := 30.0; spawnRadius <= 38; spawnRadius += 8 {
		for i := 0.0; i < 720; i += 10 {
			theta := 2.0 * math.Pi * (i / 360)
			offset := Vec2D{
				spawnRadius * math.Cos(theta),
				spawnRadius * math.Sin(theta),
			}
			testingSpawnSpatial(w,
				e.GetVec2D("Position").Add(offset),
				Vec2D{5, 5})
		}
	}
	w.SpatialHasher.Update()
	for i := 0; i < b.N; i++ {
		w.EntitiesWithinDistance(
			*e.GetVec2D("Position"),
			*e.GetVec2D("Box"),
			30.0)
	}
}
