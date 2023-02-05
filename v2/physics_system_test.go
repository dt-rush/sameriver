package sameriver

import (
	"fmt"
	"testing"
	"time"
)

func TestPhysicsSystemMotion(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)
	*e.GetVec2D("Velocity") = Vec2D{1, 1}
	pos := *e.GetVec2D("Position")
	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_DURATION_INT / 2)
	time.Sleep(FRAME_DURATION)
	w.Update(FRAME_DURATION_INT / 2)
	if *e.GetVec2D("Position") == pos {
		t.Fatal("failed to update position")
	}
}

func TestPhysicsSystemMany(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	for i := 0; i < 500; i++ {
		testingSpawnPhysics(w)
	}
	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_DURATION_INT / 2)
	time.Sleep(FRAME_DURATION)
	w.Update(FRAME_DURATION_INT / 2)
}

func TestPhysicsSystemBounds(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)
	directions := []Vec2D{
		Vec2D{100, 0},
		Vec2D{-100, 0},
		Vec2D{0, 100},
		Vec2D{0, -100},
	}
	worldCenter := Vec2D{w.Width / 2, w.Height / 2}
	worldTopRight := Vec2D{w.Width, w.Height}
	pos := e.GetVec2D("Position")
	box := e.GetVec2D("Box")
	vel := e.GetVec2D("Velocity")
	for _, d := range directions {
		*pos = Vec2D{512, 512}
		*vel = d
		for i := 0; i < 64; i++ {
			w.Update(FRAME_DURATION_INT / 2)
			time.Sleep(1 * time.Millisecond)
		}
		if !RectWithinRect(*pos, *box, worldCenter, worldTopRight) {
			t.Fatal(fmt.Sprintf("traveling with velocity %v placed entity "+
				"outside world (at position %v, box %v)", *vel, *pos, *box))
		}
	}
}
