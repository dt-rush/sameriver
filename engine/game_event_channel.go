package engine

import (
	"sync"
)

type GameEventChannel struct {
	active      bool
	activeLock  sync.RWMutex
	C           chan (GameEvent)
	channelLock sync.Mutex
}

func NewGameEventChannel() {
	return GameEventChannel{
		active: true,
		C:      make(chan (GameEvent), GAME_EVENT_CHANNEL_CAPACITY)}
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

func (c *GameEventChannel) PushToChannel(e GameEvent) {
	c.channelLock.Lock()
	defer c.channelLock.Lock()
	c.C <- e
}

func (c *GameEventChannel) DrainChannel() {
	c.channelLock.Lock()
	defer c.channelLock.Lock()
	n := len(c.C)
	for i := 0; i < n; i++ {
		<-c.C
	}
}
