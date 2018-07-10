package engine

import (
	"testing"
)

func TestPhysicsSystemMotion(t *testing.T) {
	w := NewWorld(1024, 1024)
	e, _ := w.em.spawn(physicsSpawnRequestData())
	w.em.Components.Velocity[e.ID] = Vec2D{1, 1}
	pos := w.em.Components.Position[e.ID]
	w.ps.Update(FRAME_SLEEP_MS)
	if w.em.Components.Position[e.ID] == pos {
		t.Fatal("failed to update position")
	}
}

func TestPhysicsSystemBounds(t *testing.T) {
	w := NewWorld(1024, 1024)
	e, _ := w.em.spawn(physicsSpawnRequestData())
	directions := []Vec2D{
		Vec2D{1, 0},
		Vec2D{-1, 0},
		Vec2D{0, 1},
		Vec2D{0, -1},
	}
	for _, d := range directions {
		w.em.Components.Position[e.ID] = Vec2D{512, 512}
		w.em.Components.Velocity[e.ID] = d
		for i := 0; i < 2048; i++ {
			w.ps.Update(FRAME_SLEEP_MS)
		}
		if !w.EntityIsWithinRect(e, Vec2D{0, 0}, Vec2D{1024, 1024}) {
			t.Fatal("Update() placed entity outside world")
		}
	}
}
