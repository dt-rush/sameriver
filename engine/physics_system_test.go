package engine

import (
	"testing"
	"time"
)

func TestPhysicsSystemMotion(t *testing.T) {
	w := NewWorld(1024, 1024)
	ps := NewPhysicsSystem()
	w.AddSystems(ps)
	e, _ := w.Em.Spawn(physicsSpawnRequestData())
	w.Em.Components.Velocity[e.ID] = Vec2D{1, 1}
	pos := w.Em.Components.Position[e.ID]
	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_SLEEP_MS / 2)
	time.Sleep(FRAME_SLEEP)
	w.Update(FRAME_SLEEP_MS / 2)
	if w.Em.Components.Position[e.ID] == pos {
		t.Fatal("failed to update position")
	}
}

func TestPhysicsSystemBounds(t *testing.T) {
	w := NewWorld(1024, 1024)
	ps := NewPhysicsSystem()
	w.AddSystems(ps)
	e, _ := w.Em.Spawn(physicsSpawnRequestData())
	directions := []Vec2D{
		Vec2D{1, 0},
		Vec2D{-1, 0},
		Vec2D{0, 1},
		Vec2D{0, -1},
	}
	pos := &w.Em.Components.Position[e.ID]
	box := &w.Em.Components.Box[e.ID]
	vel := &w.Em.Components.Velocity[e.ID]
	for _, d := range directions {
		*pos = Vec2D{512, 512}
		*vel = d
		for i := 0; i < 2048; i++ {
			w.Update(FRAME_SLEEP_MS / 2)
		}
		if !RectWithinRect(*pos, *box, Vec2D{0, 0}, Vec2D{1024, 1024}) {
			t.Fatal("Update() placed entity outside world")
		}
	}
}
