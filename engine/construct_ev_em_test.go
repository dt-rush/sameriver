package engine

import (
	"testing"
)

func TestConstructEvEm(t *testing.T) {
	ev := NewEventBus()
	if ev == nil {
		t.Fatal("Could not construct NewEventBus()")
	}
	em := NewEntityManager(ev)
	if em == nil {
		t.Fatal("Could not construct NewEntityManager()")
	}
}
