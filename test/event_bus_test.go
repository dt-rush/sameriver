package main

import (
	"github.com/dt-rush/sameriver/engine"
	"testing"
	"time"
)

func TestSimpleEventQueryMatching(t *testing.T) {
	ev := engine.NewEventBus()
	ec := ev.Subscribe(
		"SimpleCollisionQuery",
		engine.NewSimpleEventQuery(engine.COLLISION_EVENT))
	ev.Publish(engine.COLLISION_EVENT, nil)
	time.Sleep(16 * time.Millisecond)
	select {
	case _ = <-ec.C:
		break
	default:
		t.Fatal("simple event query wasn't received by subscriber channel " +
			"within 16 ms")
	}
}

func TestSimpleEventQueryNonMatching(t *testing.T) {
	ev := engine.NewEventBus()
	ec := ev.Subscribe(
		"SimpleCollisionQuery",
		engine.NewSimpleEventQuery(engine.COLLISION_EVENT))
	ev.Publish(engine.SPAWNREQUEST_EVENT, nil)
	time.Sleep(16 * time.Millisecond)
	select {
	case _ = <-ec.C:
		t.Fatal("simple event query sent event to wrong type channel")
	default:
		break
	}
}

func TestCollisionEventQueryMatching(t *testing.T) {
	ev := engine.NewEventBus()
	collision := engine.CollisionData{
		EntityA: &engine.EntityToken{ID: 0},
		EntityB: &engine.EntityToken{ID: 1},
	}
	ec := ev.Subscribe(
		"PredicateCollisionQuery",
		engine.NewPredicateEventQuery(
			engine.COLLISION_EVENT,
			func(e engine.Event) bool {
				return e.Data.(engine.CollisionData) == collision
			}),
	)
	ev.Publish(engine.COLLISION_EVENT, collision)
	time.Sleep(16 * time.Millisecond)
	select {
	case _ = <-ec.C:
		break
	default:
		t.Fatal("collision event query did not match")
	}
}

func TestCollisionEventQueryNonMatching(t *testing.T) {
	ev := engine.NewEventBus()
	collision := engine.CollisionData{
		EntityA: &engine.EntityToken{ID: 0},
		EntityB: &engine.EntityToken{ID: 1},
	}
	ec := ev.Subscribe(
		"PredicateCollisionQuery",
		engine.NewPredicateEventQuery(
			engine.COLLISION_EVENT,
			func(e engine.Event) bool {
				return e.Data.(engine.CollisionData) == collision
			}),
	)
	ev.Publish(engine.COLLISION_EVENT,
		engine.CollisionData{
			EntityA: &engine.EntityToken{ID: 7},
			EntityB: &engine.EntityToken{ID: 9},
		})
	time.Sleep(16 * time.Millisecond)
	select {
	case _ = <-ec.C:
		t.Fatal("collision event query matched for wrong event data")
	default:
		break
	}
}
