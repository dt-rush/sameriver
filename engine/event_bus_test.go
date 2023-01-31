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
		SimpleEventFilter("collision"))
	ev.Publish("collision", nil)
	time.Sleep(FRAME_DURATION)
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
		SimpleEventFilter("collision"))
	for i := 0; i < EVENT_SUBSCRIBER_CHANNEL_CAPACITY+4; i++ {
		ev.Publish("collision", nil)
	}
}

func TestEventBusDeactivatedSubscriber(t *testing.T) {
	ev := NewEventBus()
	ec := ev.Subscribe(
		"SimpleCollisionFilter",
		SimpleEventFilter("collision"))
	ec.Deactivate()
	ev.Publish("collision", nil)
	time.Sleep(FRAME_DURATION)
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
		SimpleEventFilter("collision"))
	ev.Publish("spawn-request", nil)
	time.Sleep(FRAME_DURATION)
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
		This:  &Entity{ID: 0},
		Other: &Entity{ID: 1},
	}
	ec := ev.Subscribe(
		"PredicateCollisionFilter",
		PredicateEventFilter(
			"collision",
			func(e Event) bool {
				return e.Data.(CollisionData) == collision
			}),
	)
	ev.Publish("collision", collision)
	time.Sleep(FRAME_DURATION)
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
		This:  &Entity{ID: 0},
		Other: &Entity{ID: 1},
	}
	ec := ev.Subscribe(
		"PredicateCollisionFilter",
		PredicateEventFilter(
			"collision",
			func(e Event) bool {
				return e.Data.(CollisionData) == collision
			}),
	)
	ev.Publish("collision",
		CollisionData{
			This:  &Entity{ID: 7},
			Other: &Entity{ID: 9},
		})
	time.Sleep(FRAME_DURATION)
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
		SimpleEventFilter("collision"))
	ev.Unsubscribe(ec)
	ev.Publish("collision", nil)
	time.Sleep(FRAME_DURATION)
	select {
	case _ = <-ec.C:
		t.Fatal("received event on unsubscribed channel")
	default:
		break
	}
}
