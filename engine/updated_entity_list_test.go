package engine

import (
	"testing"
)

func TestUpdatedEntityListAddRemove(t *testing.T) {
	list := NewUpdatedEntityList()
	e := &Entity{ID: 0, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e})
	if list.Length() != 1 {
		t.Fatal("entity was not added to list")
	}
	list.Signal(EntitySignal{ENTITY_REMOVE, e})
	if list.Length() != 0 {
		t.Fatal("entity was not removed from list")
	}
}

func TestSortedUpdatedEntityListAddRemove(t *testing.T) {
	list := NewSortedUpdatedEntityList()
	e8 := &Entity{ID: 8, Active: true, Despawned: false}
	e0 := &Entity{ID: 0, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e8})
	list.Signal(EntitySignal{ENTITY_ADD, e0})
	if list.entities[0].ID != 0 {
		t.Fatal("didn't insert in order")
	}
}

func TestUpdatedEntityListCallback(t *testing.T) {
	list := NewUpdatedEntityList()
	ran := false
	list.AddCallback(func(signal EntitySignal) {
		ran = true
	})
	e := &Entity{ID: 0, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e})
	if !ran {
		t.Fatal("callback didn't run")
	}
}

func TestUpdatedEntityListAccess(t *testing.T) {
	list := NewUpdatedEntityList()
	e0 := &Entity{ID: 0, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e0})
	entities := list.GetEntities()
	if len(entities) != 1 {
		t.Fatal("GetEntities() didn't contain the spawned entity")
	}
	e1 := &Entity{ID: 1, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e1})
	e, err := list.FirstEntity()
	if err != nil {
		t.Fatal("FirstEntity() returned err when there were 2 entities")
	}
	if e.ID != 0 {
		t.Fatal("FirstEntity() did not return first entity")
	}
	e, err = list.RandomEntity()
	if err != nil {
		t.Fatal("RandomEntity() returned err when there were 2 entities")
	}
	if !(e.ID == 0 || e.ID == 1) {
		t.Fatal("RandomEntity() did not return entity in list")
	}
	list.Signal(EntitySignal{ENTITY_REMOVE, e0})
	list.Signal(EntitySignal{ENTITY_REMOVE, e1})
	_, err = list.FirstEntity()
	if err == nil {
		t.Fatal("Should have returned err when list empty for FirstEntity()")
	}
	_, err = list.RandomEntity()
	if err == nil {
		t.Fatal("Should have returned err when list empty for RandomEntity()")
	}
}

func TestUpdatedEntityListToString(t *testing.T) {
	list := NewUpdatedEntityList()
	s0 := list.String()
	list.Signal(EntitySignal{
		ENTITY_ADD,
		&Entity{ID: 0, Active: true, Despawned: false}})
	s1 := list.String()
	if !(len(s0) < len(s1)) {
		t.Fatal("string doesn't seem to build when entitites added")
	}
}
