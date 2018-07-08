package engine

import (
	"fmt"
	"testing"
)

func TestUpdatedEntityList(t *testing.T) {
	list := NewUpdatedEntityList()
	e := &EntityToken{ID: 0, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e})
	if list.Length() != 1 {
		t.Fatal("entity was not added to list")
	}
	list.Signal(EntitySignal{ENTITY_REMOVE, e})
	if list.Length() != 0 {
		t.Fatal("entity was not removed from list")
	}
}

func TestSortedUpdatedEntityList(t *testing.T) {
	list := NewSortedUpdatedEntityList()
	e8 := &EntityToken{ID: 8, Active: true, Despawned: false}
	e0 := &EntityToken{ID: 0, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e8})
	list.Signal(EntitySignal{ENTITY_ADD, e0})
	if list.Length() != 2 {
		t.Fatal(fmt.Sprintf("entities were not added to list "+
			"(size should be %d, was %d)", 2, list.Length()))
	}
	if list.Entities[0].ID != 0 {
		t.Fatal("didn't insert in order")
	}
}
