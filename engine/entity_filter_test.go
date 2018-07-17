package engine

import (
	"testing"
)

func TestEntityFilter(t *testing.T) {
	w := testingWorld()

	pos := Vec2D{0, 0}
	req := positionSpawnRequest(pos)
	w.em.spawn(req)
	w.em.Update()
	e := w.em.Entities[0]
	q := EntityFilter{
		"positionFilter",
		func(e *Entity) bool {
			return w.em.Components.Position[e.ID] == pos
		},
	}
	if !q.Test(e) {
		t.Fatal("Filter did not return true")
	}
}

func TestEntityFilterFromTag(t *testing.T) {
	w := testingWorld()

	tag := "tag1"
	req := simpleTaggedSpawnRequest(tag)
	w.em.spawn(req)
	w.em.Update()
	e := w.em.Entities[0]
	q := w.entityFilterFromTag(tag)
	if !q.Test(e) {
		t.Fatal("Filter did not return true")
	}
}
