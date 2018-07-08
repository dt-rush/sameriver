package engine

import (
	"testing"
	"time"
)

func TestCollision(t *testing.T) {
	ev := NewEventBus()
	em := NewEntityManager(ev)
	cs := NewCollisionSystem(em, ev)
	ec := ev.Subscribe(
		"SimpleCollisionQuery",
		NewSimpleEventQuery(COLLISION_EVENT))
	em.Spawn(collisionSpawnRequestData())
	em.Spawn(collisionSpawnRequestData())
	em.Update()
	cs.Update()
	time.Sleep(16 * time.Millisecond)
	select {
	case _ = <-ec.C:
		break
	default:
		t.Fatal("collision event wasn't received within 16 ms")
	}
}
