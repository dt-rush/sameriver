package main

import (
	"fmt"
	"github.com/dt-rush/sameriver/engine"
	"testing"
)

func TestUpdatedEntityList(t *testing.T) {
	list := engine.NewUpdatedEntityList()
	e := &engine.EntityToken{ID: 0, Active: true, Despawned: false}
	list.Signal(engine.EntitySignal{engine.ENTITY_ADD, e})
	if list.Length() != 1 {
		t.Fatal("entity was not added to list")
	}
	list.Signal(engine.EntitySignal{engine.ENTITY_REMOVE, e})
	if list.Length() != 0 {
		t.Fatal("entity was not removed from list")
	}
}

func TestSortedUpdatedEntityList(t *testing.T) {
	list := engine.NewSortedUpdatedEntityList()
	e8 := &engine.EntityToken{ID: 8, Active: true, Despawned: false}
	e0 := &engine.EntityToken{ID: 0, Active: true, Despawned: false}
	list.Signal(engine.EntitySignal{engine.ENTITY_ADD, e8})
	list.Signal(engine.EntitySignal{engine.ENTITY_ADD, e0})
	if list.Length() != 2 {
		t.Fatal(fmt.Sprintf("entities were not added to list "+
			"(size should be %d, was %d)", 2, list.Length()))
	}
	if list.Entities[0].ID != 0 {
		t.Fatal("didn't insert in order")
	}
}
