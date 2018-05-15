/**
  *
  *
  *
  *
**/

// Game event manager is a bunch of channels used for
// systems to interact outside of reading
// and writing to components, using pub-sub pattern

package engine

import (
	"sync"
)

type GameEvent struct {
	Code int
	Name string
}

func (e GameEvent) String() string {
	return e.Name
}

type GameEventManager struct {
	channels    [](chan GameEvent)
	subscribers map[int]([](chan GameEvent))

	subscribe_mutex *sync.Mutex
}

func (m *GameEventManager) Init() {
	// capacity is arbitrary (can be tuned?) This should be expected to grow
	capacity := 4
	// be ready for about n^2 channels for n subscribers
	m.channels = make([](chan GameEvent), capacity*capacity)
	m.subscribers = make(map[int]([](chan GameEvent)), capacity)

	m.subscribe_mutex = &sync.Mutex{}
}

func (m *GameEventManager) Subscribe(e GameEvent) chan GameEvent {
	m.subscribe_mutex.Lock()
	// create subscribers array if DNE
	_, ok := m.subscribers[e.Code]
	if !ok {
		m.subscribers[e.Code] = make([](chan GameEvent), 0)
	}
	ch := make(chan GameEvent)
	m.subscribers[e.Code] = append(m.subscribers[e.Code], ch)
	m.channels = append(m.channels, ch)
	m.subscribe_mutex.Unlock()
	return ch
}

func (m *GameEventManager) Publish(e GameEvent) {
	if DEBUG_GAME_EVENTS {
		Logger.Printf("[Game event manager] E: %s\n", e)
	}
	// send e to all streams listening for GameEvent
	for _, ch := range m.subscribers[e.Code] {
		go func(ch chan GameEvent, e GameEvent) {
			ch <- e
		}(ch, e)
	}
}

func (m *GameEventManager) NumberOfActiveChannels() int {
	return len(m.channels)
}

func (m *GameEventManager) GetChannel() chan GameEvent {
	new_channel := make(chan GameEvent)
	m.channels = append(m.channels, new_channel)
	return new_channel
}
