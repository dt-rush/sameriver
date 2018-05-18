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

func gameEventDebug(s string, params ...interface{}) {
	switch {
	case DEBUG_GAME_EVENTS:
		return
	case len(params) == 0:
		Logger.Printf(s)
	default:
		Logger.Printf(s, params)
	}
}

type GameEventManager struct {
	// subscriberLists is a list of lists of GameEventChannels
	// where the inner lists are indexed by the GameEventType (type aliased
	// to int). So you could have a list of queries on CollisionEvents, etc.
	// Each GameEventQuery's Predicate will be tested against the events
	// that are published for the matching type (and thus the predicates
	// can safely assert the type of the Data member of the event)
	subscriberLists [N_GAME_EVENT_TYPES][]GameEventChannel
	// Mutex to protect the modification of the above
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

	gameEventDebug("[Game event manager] Subscribe: %s on channel %v\n",
		name, c)

	// Create a channel to return to the user
	c := NewGameEventChannel(q, name)
	// Add the channel to the subscriber list for its type
	m.subscriberLists[q.Type] = append(m.subscriberLists[q.Type], c)
	// return the channel to the caller
	return c
}

// Remove a subscriber
func (m *GameEventManager) Unsubscribe(c GameEventChannel) {

	gameEventDebug("[Game event manager] Unsubscribe on channel %v\n", c)

	// remove the query watcher from the subscriber list associated with
	// the given channel's
	removeGameEventChannelFromSlice(c, &m.subscriberLists[c.Type])
}

// Publish a game event for anyone listening
func (m *GameEventManager) Publish(e GameEvent) {

	gameEventDebug("[Game event manager] ⚹: %s\n", e.Description)

	// send e to all matching watchers
	for _, c := range m.subscriberLists[e.Type] {
		if len(c.C) == GAME_EVENT_CHANNEL_CAPACITY {
			Logger.Printf("[Game event manager] ⚠ event channel #%d "+
				"for %s is full - discarding event\n", c.Name)
		} else if c.IsActive() {
			gameEventDebug("[Game event manager] sending %s on %s",
				e.Description, c.Name)
			c.Send(e)
		} else {
			Logger.Printf("[Game event manager] %s channel inactive, "+
				"not sending event %s\n", c.Name, e.Description)
		}
	}
}
