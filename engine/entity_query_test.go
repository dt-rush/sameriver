package engine

import (
	"testing"
)

func TestEntityQuery(t *testing.T) {
	w := NewWorld(1024, 1024)

	pos := Vec2D{0, 0}
	req := positionSpawnRequestData(pos)
	w.em.Spawn(req)
	w.em.Update()
	e := w.em.Entities[0]
	q := EntityQuery{
		"positionQuery",
		func(e *EntityToken, em *EntityManager) bool {
			return w.em.Components.Position[e.ID] == pos
		},
	}
	if !q.Test(e, w.em) {
		t.Fatal("query did not return true")
	}
}

func TestEntityQueryFromTag(t *testing.T) {
	w := NewWorld(1024, 1024)

	tag := "tag1"
	req := simpleTaggedSpawnRequestData(tag)
	w.em.Spawn(req)
	w.em.Update()
	e := w.em.Entities[0]
	q := EntityQueryFromTag(tag)
	if !q.Test(e, w.em) {
		t.Fatal("query did not return true")
	}
}
