package engine

import (
	"fmt"
	"testing"
	"time"
)

func TestPhysicsSystemMotion(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.AddSystems(ps)
	e, _ := testingSpawnPhysics(w)
	w.em.components.Velocity[e.ID] = Vec2D{1, 1}
	pos := w.em.components.Position[e.ID]
	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_SLEEP_MS / 2)
	time.Sleep(FRAME_SLEEP)
	w.Update(FRAME_SLEEP_MS / 2)
	if w.em.components.Position[e.ID] == pos {
		t.Fatal("failed to update position")
	}
}

func TestPhysicsSystemBounds(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.AddSystems(ps)
	e, _ := testingSpawnPhysics(w)
	directions := []Vec2D{
		Vec2D{100, 0},
		Vec2D{-100, 0},
		Vec2D{0, 100},
		Vec2D{0, -100},
	}
	worldCenter := Vec2D{w.Width / 2, w.Height / 2}
	worldTopRight := Vec2D{w.Width, w.Height}
	pos := &w.em.components.Position[e.ID]
	box := &w.em.components.Box[e.ID]
	vel := &w.em.components.Velocity[e.ID]
	for _, d := range directions {
		*pos = Vec2D{512, 512}
		*vel = d
		for i := 0; i < 64; i++ {
			w.Update(FRAME_SLEEP_MS / 2)
			time.Sleep(1 * time.Millisecond)
		}
		if !RectWithinRect(pos, box, &worldCenter, &worldTopRight) {
			t.Fatal(fmt.Sprintf("traveling with velocity %v placed entity "+
				"outside world (at position %v, box %v)", *vel, *pos, *box))
		}
	}
}
