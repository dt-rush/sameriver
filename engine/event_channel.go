package engine

import (
	"sync"
	"sync/atomic"
)

type EventChannel struct {
	active          uint32
	C               chan (Event)
	channelSendLock sync.Mutex
	Query           EventQuery
	Name            string
}

func NewEventChannel(name string, q EventQuery) EventChannel {

	return EventChannel{
		active: 1,
		C:      make(chan (Event), EVENT_CHANNEL_CAPACITY),
		Query:  q,
		Name:   name}
}

func (c *EventChannel) Activate() {
	atomic.StoreUint32(&c.active, 1)
}

func (c *EventChannel) Deactivate() {
	atomic.StoreUint32(&c.active, 0)
}

func (c *EventChannel) IsActive() bool {
	return atomic.LoadUint32(&c.active) == 1
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
