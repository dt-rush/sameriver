package engine

import (
	"sync"
)

type GameEventChannel struct {
	active          bool
	activeLock      sync.RWMutex
	C               chan (GameEvent)
	channelSendLock sync.Mutex
	Query           GameEventQuery
	Name            string
}

func NewGameEventChannel(
	q GameEventQuery, name string) GameEventChannel {

	return GameEventChannel{
		active: true,
		C:      make(chan (GameEvent), GAME_EVENT_CHANNEL_CAPACITY),
		Query:  q,
		Name:   name}
}

func (c *GameEventChannel) Activate() {
	c.activeLock.Lock()
	c.active = true
	c.activeLock.Unlock()
}

func (c *GameEventChannel) Deactivate() {
	c.activeLock.Lock()
	c.active = false
	c.activeLock.Unlock()
}

func (c *GameEventChannel) IsActive() bool {
	c.activeLock.RLock()
	defer c.activeLock.RLock()
	return c.active
}

func (c *GameEventChannel) Send(e GameEvent) {
	c.channelSendLock.Lock()
	defer c.channelSendLock.Lock()
	c.C <- e
}

func (c *GameEventChannel) DrainChannel() {
	c.channelSendLock.Lock()
	defer c.channelSendLock.Lock()
	n := len(c.C)
	for i := 0; i < n; i++ {
		<-c.C
	}
}
