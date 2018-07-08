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
	em.Components.Position[em.Entities()[0].ID] = Vec2D{100000, 100000}
	cs.Update()
	time.Sleep(16 * time.Millisecond)
	select {
	case _ = <-ec.C:
		t.Fatal("collision event occurred but entities were not overlapping")
	default:
		break
	}
}
