package main

import (
	"fmt"
	"github.com/dt-rush/sameriver/engine"
	"testing"
	"time"
)

func TestUpdatedEntityList(t *testing.T) {
	c := make(chan engine.EntitySignal)
	list := engine.NewUpdatedEntityList(c)
	list.Start()
	e := &engine.EntityToken{ID: 0, Active: true, Despawned: false}
	c <- engine.EntitySignal{engine.ENTITY_ADD, e}
	time.Sleep(16 * time.Millisecond)
	if list.Length() != 1 {
		t.Fatal("entity was not added to list")
	}
	c <- engine.EntitySignal{engine.ENTITY_REMOVE, e}
	time.Sleep(16 * time.Millisecond)
	if list.Length() != 0 {
		t.Fatal("entity was not removed from list")
	}
}

func TestSortedUpdatedEntityList(t *testing.T) {
	c := make(chan engine.EntitySignal)
	list := engine.NewSortedUpdatedEntityList(c)
	list.Start()
	e8 := &engine.EntityToken{ID: 8, Active: true, Despawned: false}
	e0 := &engine.EntityToken{ID: 0, Active: true, Despawned: false}
	c <- engine.EntitySignal{engine.ENTITY_ADD, e8}
	c <- engine.EntitySignal{engine.ENTITY_ADD, e0}
	time.Sleep(16 * time.Millisecond)
	if list.Length() != 2 {
		t.Fatal(fmt.Sprintf("entities were not added to list "+
			"(size should be %d, was %d)", 2, list.Length()))
	}
	if list.Entities[0].ID != 0 {
		t.Fatal("didn't insert in order")
	}
}
