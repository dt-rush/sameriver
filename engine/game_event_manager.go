/**
  *
  * Pub-sub hub for game events
  *
  *
**/

package engine

import (
	"sync"
)

type GameEventManager struct {
	// subscribers is a list of lists of GameEventQueryWatchers
	// where the inner lists are indexed by the GameEventType (type aliased
	// to int). So you could have a list of queries on CollisionEvents, etc.
	subscribers    [N_GAME_EVENT_TYPES][]GameEventQueryWatcher
	subscribeMutex sync.Mutex
}

func (m *GameEventManager) Init() {
	// nothing for now
}

// Subscribe to listen for game events defined by a query
func (m *GameEventManager) Subscribe(
	q GameEventQuery, name string) GameEventChannel {

	// Lock the subscriber slice while we modify it
	m.subscribeMutex.Lock()
	defer m.subscribeMutex.Unlock()
	// Create a channel to return to the user
	c := NewGameEventChannel()
	// Add a query watcher to the subscriber slice
	qw := GameEventQueryWatcher{q, c}
	m.subscribers[q.Type] = append(m.subscribers[q.Type], qw)
	// return the GameEventChannel to the caller
	return c
}

// Publish a game event for anyone listening
func (m *GameEventManager) Publish(e GameEvent) {
	if DEBUG_GAME_EVENTS {
		Logger.Printf("[Game event manager] ⚹: %s\n",
			e)
	}

	// send e to all matching watchers
	for _, qw := range m.subscribers[e.Type] {
		if len(qw.Channel.C) == GAME_EVENT_CHANNEL_CAPACITY {
			Logger.Printf("[Game event manager] ⚠  event channel #%d "+
				"for %s is full - discarding event\n", qw.Name)
		} else {
			qw.Channel.PushToChannel(e)
		}
	}
}
