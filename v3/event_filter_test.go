package sameriver

import (
	"testing"
)

func TestEventFilterSimple(t *testing.T) {
	f := SimpleEventFilter("spawn-request")
	if !f.Test(Event{Type: "spawn-request"}) {
		t.Fatal("did not test true for matching event type")
	}
	if f.Test(Event{Type: "collision"}) {
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
		Type: "collision",
		Data: CollisionData{This: entity}}
	f := PredicateEventFilter("collision", pred)
	if !f.Test(event) {
		t.Fatal("filter did not match, should have")
	}
}

func TestEventFilterCollision(t *testing.T) {
	this := &Entity{ID: 100}
	other := &Entity{ID: 100}
	predicate := func(ev Event) bool {
		c := ev.Data.(CollisionData)
		return c.This == this && c.Other == other
	}
	event := Event{
		Type: "collision",
		Data: CollisionData{This: this, Other: other}}
	f := PredicateEventFilter("collision", predicate)
	if !f.Test(event) {
		t.Fatal("filter did not match, should have")
	}
}
