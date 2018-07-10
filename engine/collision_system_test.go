package engine

import (
	"testing"
	"time"
)

func TestCollisionSystem(t *testing.T) {
	w := NewWorld(1024, 1024)
	ec := w.ev.Subscribe(
		"SimpleCollisionQuery",
		NewSimpleEventQuery(COLLISION_EVENT))
	w.em.Spawn(collisionSpawnRequestData())
	w.em.Spawn(collisionSpawnRequestData())
	w.em.Update()
	w.cs.Update()
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		break
	default:
		t.Fatal("collision event wasn't received within 16 ms")
	}
	w.em.Components.Position[w.em.Entities()[0].ID] = Vec2D{100000, 100000}
	w.cs.Update()
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		t.Fatal("collision event occurred but entities were not overlapping")
	default:
		break
	}
}

func TestCollisionRateLimit(t *testing.T) {
	w := NewWorld(1024, 1024)
	ec := w.ev.Subscribe(
		"SimpleCollisionQuery",
		NewSimpleEventQuery(COLLISION_EVENT))
	w.em.Spawn(collisionSpawnRequestData())
	w.em.Spawn(collisionSpawnRequestData())
	w.em.Update()
	w.cs.Update()
	w.cs.Update()
	time.Sleep(FRAME_SLEEP)
	if len(ec.C) != 1 {
		t.Fatal("collision rate-limiter didn't prevent collision duplication")
	}
}
