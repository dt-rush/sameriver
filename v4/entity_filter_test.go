package sameriver

import (
	"testing"
)

func TestEntityFilter(t *testing.T) {
	w := testingWorld()

	pos := Vec2D{0, 0}
	e := testingSpawnPosition(w, pos)
	q := EntityFilter{
		"positionFilter",
		func(e *Entity) bool {
			return *e.GetVec2D(POSITION) == pos
		},
	}
	if !q.Test(e) {
		t.Fatal("Filter did not return true")
	}
}

func TestEntityFilterFromTag(t *testing.T) {
	w := testingWorld()

	tag := "tag1"
	e := testingSpawnTagged(w, tag)
	q := EntityFilterFromTag(tag)
	if !q.Test(e) {
		t.Fatal("Filter did not return true")
	}
}
