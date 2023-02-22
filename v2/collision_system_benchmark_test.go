package sameriver

import (
	"math/rand"

	"testing"
)

func BenchmarkCollisionMany(b *testing.B) {
	w := testingWorld()
	sh := NewSpatialHashSystem(10, 10)
	cs := NewCollisionSystem(FRAME_DURATION)
	p := NewPhysicsSystem()
	w.RegisterSystems(sh, cs, p)

	w.SetSystemSchedule("CollisionSystem", 5)
	for i := 0; i < MAX_ENTITIES; i++ {
		w.Spawn(map[string]any{
			"Vec2D,Position": Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			"Vec2D,Box":      Vec2D{5, 5},
			"Vec2D,Velocity": Vec2D{rand.Float64(), rand.Float64()},
		})
	}
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100; i++ {
			p.Update(200)
			cs.Update(500)
		}
	}
}
