package engine

import (
	"testing"
)

func TestEntityQuery(t *testing.T) {
	w := NewWorld(1024, 1024)

	pos := Vec2D{0, 0}
	req := positionSpawnRequestData(pos)
	w.Em.Spawn(req)
	w.Em.Update()
	e := w.Em.Entities[0]
	q := EntityQuery{
		"positionQuery",
		func(e *EntityToken, em *EntityManager) bool {
			return w.Em.Components.Position[e.ID] == pos
		},
	}
	if !q.Test(e, w.Em) {
		t.Fatal("query did not return true")
	}
}

func TestEntityQueryFromTag(t *testing.T) {
	w := NewWorld(1024, 1024)

	tag := "tag1"
	req := simpleTaggedSpawnRequestData(tag)
	w.Em.Spawn(req)
	w.Em.Update()
	e := w.Em.Entities[0]
	q := EntityQueryFromTag(tag)
	if !q.Test(e, w.Em) {
		t.Fatal("query did not return true")
	}
}
