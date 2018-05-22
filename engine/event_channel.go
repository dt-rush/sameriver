package engine

import (
	"sync"
)

type EventChannel struct {
	active          bool
	activeLock      sync.RWMutex
	C               chan (Event)
	channelSendLock sync.Mutex
	Query           EventQuery
	Name            string
}

func NewEventChannel(
	q EventQuery, name string) EventChannel {

	return EventChannel{
		active: true,
		C:      make(chan (Event), EVENT_CHANNEL_CAPACITY),
		Query:  q,
		Name:   name}
}

func (c *EventChannel) Activate() {
	c.activeLock.Lock()
	c.active = true
	c.activeLock.Unlock()
}

func (c *EventChannel) Deactivate() {
	c.activeLock.Lock()
	c.active = false
	c.activeLock.Unlock()
}

func (c *EventChannel) IsActive() bool {
	c.activeLock.RLock()
	defer c.activeLock.RLock()
	return c.active
}

func (c *EventChannel) Send(e Event) {
	c.channelSendLock.Lock()
	defer c.channelSendLock.Lock()
	c.C <- e
}

func (c *EventChannel) DrainChannel() {
	c.channelSendLock.Lock()
	defer c.channelSendLock.Lock()
	n := len(c.C)
	for i := 0; i < n; i++ {
		<-c.C
	}
}
