package engine

import (
	"testing"
	"time"
)

func TestEventBusConstructEventBus(t *testing.T) {
	ev := NewEventBus()
	if ev == nil {
		t.Fatal("Could not construct NewEventBus()")
	}
}

func TestEventBusSimpleEventFilterMatching(t *testing.T) {
	ev := NewEventBus()
	ec := ev.Subscribe(
		"SimpleCollisionFilter",
		SimpleEventFilter(COLLISION_EVENT))
	ev.Publish(COLLISION_EVENT, nil)
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		break
	default:
		t.Fatal("simple event Filter wasn't received by subscriber channel " +
			"within 16 ms")
	}
}

func TestEventBusMaxCapacity(t *testing.T) {
	ev := NewEventBus()
	_ = ev.Subscribe(
		"SimpleCollisionFilter",
		SimpleEventFilter(COLLISION_EVENT))
	for i := 0; i < EVENT_SUBSCRIBER_CHANNEL_CAPACITY+4; i++ {
		ev.Publish(COLLISION_EVENT, nil)
	}
}

func TestEventBusDeactivatedSubscriber(t *testing.T) {
	ev := NewEventBus()
	ec := ev.Subscribe(
		"SimpleCollisionFilter",
		SimpleEventFilter(COLLISION_EVENT))
	ec.Deactivate()
	ev.Publish(COLLISION_EVENT, nil)
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		t.Fatal("event was received on deactivated EventChannel")
	default:
	}
}

func TestEventBusSimpleEventFilterNonMatching(t *testing.T) {
	ev := NewEventBus()
	ec := ev.Subscribe(
		"SimpleCollisionFilter",
		SimpleEventFilter(COLLISION_EVENT))
	ev.Publish(SPAWNREQUEST_EVENT, nil)
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		t.Fatal("simple event Filter sent event to wrong type channel")
	default:
		break
	}
}

func TestEventBusDataEventFilterMatching(t *testing.T) {
	ev := NewEventBus()
	collision := CollisionData{
		EntityA: &EntityToken{ID: 0},
		EntityB: &EntityToken{ID: 1},
	}
	ec := ev.Subscribe(
		"PredicateCollisionFilter",
		PredicateEventFilter(
			COLLISION_EVENT,
			func(e Event) bool {
				return e.Data.(CollisionData) == collision
			}),
	)
	ev.Publish(COLLISION_EVENT, collision)
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		break
	default:
		t.Fatal("collision event Filter did not match")
	}
}

func TestEventBusDataEventFilterNonMatching(t *testing.T) {
	ev := NewEventBus()
	collision := CollisionData{
		EntityA: &EntityToken{ID: 0},
		EntityB: &EntityToken{ID: 1},
	}
	ec := ev.Subscribe(
		"PredicateCollisionFilter",
		PredicateEventFilter(
			COLLISION_EVENT,
			func(e Event) bool {
				return e.Data.(CollisionData) == collision
			}),
	)
	ev.Publish(COLLISION_EVENT,
		CollisionData{
			EntityA: &EntityToken{ID: 7},
			EntityB: &EntityToken{ID: 9},
		})
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		t.Fatal("collision event Filter matched for wrong event data")
	default:
		break
	}
}

func TestEventBusUnsubscribe(t *testing.T) {
	ev := NewEventBus()
	ec := ev.Subscribe(
		"SimpleCollisionFilter",
		SimpleEventFilter(COLLISION_EVENT))
	ev.Unsubscribe(ec)
	ev.Publish(COLLISION_EVENT, nil)
	time.Sleep(FRAME_SLEEP)
	select {
	case _ = <-ec.C:
		t.Fatal("received event on unsubscribed channel")
	default:
		break
	}
}
