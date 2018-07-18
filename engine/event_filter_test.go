package engine

import (
	"testing"
)

func TestEventFilterSimple(t *testing.T) {
	f := SimpleEventFilter(SPAWNREQUEST_EVENT)
	if !f.Test(Event{Type: SPAWNREQUEST_EVENT}) {
		t.Fatal("did not test true for matching event type")
	}
	if f.Test(Event{Type: COLLISION_EVENT}) {
		t.Fatal("tested true for non-matching event type")
	}
}

func TestEventFilterPredicate(t *testing.T) {
	entity := &Entity{ID: 100}
	pred := func(e Event) bool {
		c := e.Data.(CollisionData)
		return c.This == entity
	}
	event := Event{
		Type: COLLISION_EVENT,
		Data: CollisionData{This: entity}}
	f := PredicateEventFilter(COLLISION_EVENT, pred)
	if !f.Test(event) {
		t.Fatal("filter did not match, should have")
	}
}

func TestEventFilterCollision(t *testing.T) {
	this := &Entity{ID: 100}
	other := &Entity{ID: 100}
	predicate := func(c CollisionData) bool {
		return c.This == this && c.Other == other
	}
	event := Event{
		Type: COLLISION_EVENT,
		Data: CollisionData{This: this, Other: other}}
	f := CollisionEventFilter(predicate)
	if !f.Test(event) {
		t.Fatal("filter did not match, should have")
	}
}
