package engine

import (
	"testing"
)

func TestEntityFilter(t *testing.T) {
	w := NewWorld(1024, 1024)

	pos := Vec2D{0, 0}
	req := positionSpawnRequestData(pos)
	w.Em.Spawn(req)
	w.Em.Update()
	e := w.Em.Entities[0]
	q := EntityFilter{
		"positionFilter",
		func(e *EntityToken, em *EntityManager) bool {
			return w.Em.Components.Position[e.ID] == pos
		},
	}
	if !q.Test(e, w.Em) {
		t.Fatal("Filter did not return true")
	}
}

func TestEntityFilterFromTag(t *testing.T) {
	w := NewWorld(1024, 1024)

	tag := "tag1"
	req := simpleTaggedSpawnRequestData(tag)
	w.Em.Spawn(req)
	w.Em.Update()
	e := w.Em.Entities[0]
	q := EntityFilterFromTag(tag)
	if !q.Test(e, w.Em) {
		t.Fatal("Filter did not return true")
	}
}
