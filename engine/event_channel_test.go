package engine

import (
	"testing"
)

func TestNewEventChannel(t *testing.T) {
	q := &EventFilter{Type: SPAWNREQUEST_EVENT}
	ec := NewEventChannel("testchannel", q)
	if !(ec.IsActive() &&
		ec.C != nil &&
		ec.Filter == q &&
		ec.Name == "testchannel") {
		t.Fatal("did not construct properly")
	}
}

func TestNewEventChannelActivateDeactivate(t *testing.T) {
	ec := NewEventChannel("testchannel", nil)
	ec.Deactivate()
	if ec.IsActive() {
		t.Fatal("Deactivate() didn't change result of IsActive()")
	}
	ec.Activate()
	if !ec.IsActive() {
		t.Fatal("Activate() didn't change result of IsActive()")
	}
}

func TestEventChannelDrain(t *testing.T) {
	ec := NewEventChannel("testchannel", nil)
	for i := 0; i < EVENT_SUBSCRIBER_CHANNEL_CAPACITY; i++ {
		ec.C <- Event{}
	}
	ec.DrainChannel()
	if len(ec.C) != 0 {
		t.Fatal("drain channel did not remove all events")
	}
}
