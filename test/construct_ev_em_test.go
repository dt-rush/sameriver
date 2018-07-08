package main

import (
	"github.com/dt-rush/sameriver/engine"
	"testing"
)

func TestConstructEvEm(t *testing.T) {
	ev := engine.NewEventBus()
	if ev == nil {
		t.Fatal("Could not construct NewEventBus()")
	}
	em := engine.NewEntityManager(ev)
	if em == nil {
		t.Fatal("Could not construct NewEntityManager()")
	}
}
