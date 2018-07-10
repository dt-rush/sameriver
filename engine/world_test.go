package engine

import (
	"testing"
)

func TestCanConstructWorld(t *testing.T) {
	w := NewWorld(1024, 1024)
	if w == nil {
		t.Fatal("NewWorld() was nil")
	}
}

func TestEntityWithinRect(t *testing.T) {
	w := NewWorld(1024, 1024)
	e, _ := w.em.spawn(physicsSpawnRequestData())
	w.em.Components.Position[e.ID] = Vec2D{10, 10}
	w.em.Components.Box[e.ID] = Vec2D{2, 2}
	if !w.EntityIsWithinRect(e, Vec2D{0, 0}, Vec2D{12, 12}) {
		t.Fatal("entity ought to be within rect, was not considered to be so")
	}
}
