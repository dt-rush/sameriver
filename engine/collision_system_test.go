package engine

import (
	"testing"
	"time"
)

func TestCollisionSystem(t *testing.T) {
	w := NewWorld(1024, 1024)
	cs := NewCollisionSystem()
	w.AddSystems(
		NewSpatialHashSystem(10, 10),
		cs,
	)
	if cs.sh == nil {
		t.Fatal("failed to inject *SpatialHashSystem to CollisionSystem.sh")
	}
	ec := w.ev.Subscribe(
		"SimpleCollisionQuery",
		NewSimpleEventQuery(COLLISION_EVENT))
	w.em.Spawn(collisionSpawnRequestData())
	w.em.Spawn(collisionSpawnRequestData())
	w.Update(FRAME_SLEEP_MS)
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		break
	default:
		t.Fatal("collision event wasn't received within 16 ms")
	}
	w.em.Components.Position[w.em.Entities[0].ID] = Vec2D{100, 100}
	w.Update(FRAME_SLEEP_MS)
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
	cs := NewCollisionSystem()
	w.AddSystems(
		NewSpatialHashSystem(10, 10),
		cs,
	)
	ec := w.ev.Subscribe(
		"SimpleCollisionQuery",
		NewSimpleEventQuery(COLLISION_EVENT))
	e, _ := w.em.Spawn(collisionSpawnRequestData())
	w.em.Spawn(collisionSpawnRequestData())
	w.em.Update()
	w.Update(FRAME_SLEEP_MS)
	w.Update(FRAME_SLEEP_MS)
	time.Sleep(FRAME_SLEEP)
	if len(ec.C) != 1 {
		t.Fatal("collision rate-limiter didn't prevent collision duplication")
	}
	w.Update(FRAME_SLEEP_MS)
	cs.rateLimiterArray.Reset(e)
	w.Update(FRAME_SLEEP_MS)
	time.Sleep(FRAME_SLEEP)
	if len(ec.C) != 2 {
		t.Fatal("collision rate-limiter reset did not allow second collision")
	}

}
