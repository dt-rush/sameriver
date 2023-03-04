package sameriver

import (
	"math/rand"

	"testing"
)

func BenchmarkPhysicsSystemManySingleThreadUpdate(b *testing.B) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	for i := 0; i < 1000; i++ {
		e := testingSpawnPhysics(w)
		*e.GetVec2D("Velocity") = Vec2D{rand.Float64(), rand.Float64()}
	}
	// Update twice since physics system won't run the first time(needs a dt)
	for i := 0; i < b.N; i++ {
		ps.SingleThreadUpdate(FRAME_DURATION_INT / 2)
	}
}

func BenchmarkPhysicsSystemManyParallelUpdate(b *testing.B) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	for i := 0; i < 1000; i++ {
		e := testingSpawnPhysics(w)
		*e.GetVec2D("Velocity") = Vec2D{rand.Float64(), rand.Float64()}
	}
	// Update twice since physics system won't run the first time(needs a dt)
	for i := 0; i < b.N; i++ {
		ps.ParallelUpdate(FRAME_DURATION_INT / 2)
	}
}
