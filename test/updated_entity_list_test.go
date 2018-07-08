package main

import (
	"github.com/dt-rush/sameriver/engine"
	"testing"
	"time"
)

func TestUpdatedEntityListNoBacklog(t *testing.T) {
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
