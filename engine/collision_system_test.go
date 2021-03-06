package engine

import (
	"testing"
	"time"
)

func testingSetupCollision() (*World, *CollisionSystem, *EventChannel, *Entity) {
	w := testingWorld()
	cs := NewCollisionSystem(FRAME_SLEEP / 2)
	w.AddSystems(
		NewSpatialHashSystem(1, 1),
		cs,
	)
	// spawn the colliding entities
	testingSpawnCollision(w)
	e, _ := testingSpawnCollision(w)
	// subscribe to collision events
	ec := w.Events.Subscribe(
		"SimpleCollisionFilter",
		SimpleEventFilter(COLLISION_EVENT))
	return w, cs, ec, e
}

func TestCollisionSystem(t *testing.T) {
	w, _, ec, e := testingSetupCollision()
	w.Update(FRAME_SLEEP_MS / 2)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		break
	default:
		t.Fatal("collision event wasn't received within 1 frame")
	}
	// move the enitity so it no longer collides
	*e.GetPosition() = Vec2D{100, 100}
	w.Update(FRAME_SLEEP_MS / 2)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_SLEEP)
	if len(ec.C) != 1 {
		t.Fatal("collision event occurred but entities were not overlapping")
	}
}

func TestCollisionRateLimit(t *testing.T) {
	w, cs, ec, e := testingSetupCollision()
	w.Update(FRAME_SLEEP_MS / 2)
	w.Update(FRAME_SLEEP_MS / 2)
	time.Sleep(FRAME_SLEEP)
	if len(ec.C) == 4 {
		t.Fatal("collision rate-limiter didn't prevent collision duplication")
	}
	for len(ec.C) > 0 {
		_ = <-ec.C
	}
	// wait for rate limit to die
	time.Sleep(FRAME_SLEEP)
	// check if we can reset the rate lmiiter
	w.Update(FRAME_SLEEP_MS / 2)
	cs.rateLimiterArray.Reset(e)
	w.Update(FRAME_SLEEP_MS / 2)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_SLEEP)
	if len(ec.C) != 4 {
		t.Fatal("collision rate-limiter reset did not allow second collision")
	}
}

func TestCollisionFilter(t *testing.T) {
	w, _, _, e := testingSetupCollision()
	coin, _ := w.Spawn([]string{"coin"},
		ComponentSet{
			Position: e.GetPosition(),
			Box:      e.GetBox(),
		},
	)
	predicate := func(c CollisionData) bool {
		return c.This == e && c.Other == coin
	}
	ec := w.Events.Subscribe("PredicateCollisionFilter",
		CollisionEventFilter(predicate))
	w.Update(FRAME_SLEEP_MS / 2)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_SLEEP)
	select {
	case ev := <-ec.C:
		if !predicate(ev.Data.(CollisionData)) {
			t.Fatal("CollisionEventFilter did not select correct event")
		}
	default:
		t.Fatal("CollisionEventFilter didn't pick up collision")
	}
}
