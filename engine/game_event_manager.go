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

//
// TODO: refactor so that for each event, we iterate only those
// subscribers registered for the event type, and then we send
// only on those channels which match the query
//

type GameEventManager struct {
	eventQueryWatchers []GameEventQueryWatcher
	subscribeMutex     sync.Mutex
}

func (m *GameEventManager) Init() {
	// nothing for now
}

func (m *GameEventManager) Subscribe(
	q GameEventQuery, name string) GameEventChannel {

	m.subscribeMutex.Lock()
	defer m.subscribeMutex.Unlock()
	c := NewGameEventChannel()
	qw := GameEventQueryWatcher{q, c}
	m.eventQueryWatchers = append(m.eventQueryWatchers, qw)
	return c
}

func (m *GameEventManager) Publish(e GameEvent) {
	if DEBUG_GAME_EVENTS {
		Logger.Printf("[Game event manager] ᛤ: %s\n",
			e)
	}
	// send e to all streams listening for GameEvent
	for qw := range m.eventQueryWatchers {
		if len(qw.Channel.C) == GAME_EVENT_CHANNEL_CAPACITY {
			Logger.Printf("[Game event manager] ⚠  event channel #%d for %s is full - discarding event\n", qw.Name)
		} else {
			qw.Channel.PushToChannel(e)
		}
	}
}
