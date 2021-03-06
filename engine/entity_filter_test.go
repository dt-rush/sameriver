package engine

import (
	"testing"
)

func TestEntityFilter(t *testing.T) {
	w := testingWorld()

	pos := Vec2D{0, 0}
	e, _ := testingSpawnPosition(w, pos)
	q := EntityFilter{
		"positionFilter",
		func(e *Entity) bool {
			return *e.GetPosition() == pos
		},
	}
	if !q.Test(e) {
		t.Fatal("Filter did not return true")
	}
}

func TestEntityFilterFromTag(t *testing.T) {
	w := testingWorld()

	tag := "tag1"
	e, _ := testingSpawnTagged(w, tag)
	q := w.entityFilterFromTag(tag)
	if !q.Test(e) {
		t.Fatal("Filter did not return true")
	}
}
