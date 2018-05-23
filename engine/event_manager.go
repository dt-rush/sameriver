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

type EventManager struct {
	// subscriberLists is a list of lists of EventChannels
	// where the inner lists are indexed by the EventType (type aliased
	// to int). So you could have a list of queries on CollisionEvents, etc.
	// Each EventQuery's Predicate will be tested against the events
	// that are published for the matching type (and thus the predicates
	// can safely assert the type of the Data member of the event)
	subscriberLists [N_EVENT_TYPES][]EventChannel
	// Mutex to protect the modification of the above
	subscribeMutex sync.Mutex
}

func (m *EventManager) Init() {
	// nothing for now
}

// Subscribe to listen for game events defined by a query
func (m *EventManager) Subscribe(
	q EventQuery, name string) EventChannel {

	// Lock the subscriber slice while we modify it
	m.subscribeMutex.Lock()
	defer m.subscribeMutex.Unlock()

	// Create a channel to return to the user
	c := NewEventChannel(q, name)
	eventDebug("Subscribe: %s on channel %v\n",
		name, c)
	// Add the channel to the subscriber list for its type
	m.subscriberLists[q.Type] = append(m.subscriberLists[q.Type], c)
	// return the channel to the caller
	return c
}

// Remove a subscriber
func (m *EventManager) Unsubscribe(c EventChannel) {

	eventDebug("Unsubscribe on channel %v\n", c)

	// remove the query watcher from the subscriber list associated with
	// the given channel's
	removeEventChannelFromSlice(c, &m.subscriberLists[c.Query.Type])
}

// Publish a game event for anyone listening
func (m *EventManager) Publish(e Event) {

	eventDebug("⚹: %s\n", e.Description)

	// send e to all matching watchers
	for _, c := range m.subscriberLists[e.Type] {
		if len(c.C) == EVENT_CHANNEL_CAPACITY {
			Logger.Printf("⚠ event channel #%d "+
				"for %s is full - discarding event\n", c.Name)
		} else if c.IsActive() {
			eventDebug("sending %s on %s",
				e.Description, c.Name)
			c.Send(e)
		} else {
			Logger.Printf("%s channel inactive, "+
				"not sending event %s\n", c.Name, e.Description)
		}
	}
}
