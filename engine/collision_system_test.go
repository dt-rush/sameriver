package engine

import (
	"testing"
)

func TestCollision(t *testing.T) {
	ev := NewEventBus()
	em := NewEntityManager(ev)
	em.Spawn(collisionSpawnRequestData())
	em.Spawn(collisionSpawnRequestData())
	em.Update()
	if em.NumEntities() == 0 {
		t.Fatal("failed to spawn simple spawn request entity")
	}
	e := em.Entities()[0]
	if !e.Active {
		t.Fatal("spawned entity was not active")
	}
}
