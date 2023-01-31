package sameriver

import (
	"go.uber.org/atomic"
)

type EventChannel struct {
	active *atomic.Uint32
	C      chan Event
	filter *EventFilter
	name   string
}

func NewEventChannel(name string, q *EventFilter) *EventChannel {
	return &EventChannel{
		active: atomic.NewUint32(1),
		C:      make(chan (Event), EVENT_SUBSCRIBER_CHANNEL_CAPACITY),
		filter: q,
		name:   name}
}

func (c *EventChannel) Activate() {
	c.active.Store(1)
}

func (c *EventChannel) Deactivate() {
	c.active.Store(0)
}

func (c *EventChannel) IsActive() bool {
	return c.active.Load() == 1
}

func (c *EventChannel) DrainChannel() {
	n := len(c.C)
	for i := 0; i < n; i++ {
		<-c.C
	}
}
