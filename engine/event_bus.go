/**
  *
  * Pub-sub hub for game events
  *
**/

package engine

import (
	"sync"
)

type SubscriberList struct {
	// subscriberLists is a list of lists of EventChannels
	// where the outer list is indexed by the EventType (type aliased
	// to int). So you could have a list of queries on CollisionEvents, etc.
	// Each EventQuery's Predicate will be tested against the events
	// that are published for the matching type (and thus the predicates
	// can safely assert the type of the Data member of the event)
	channels [N_EVENT_TYPES][]EventChannel
	// Mutex to protect the modification of the above
	mutex sync.RWMutex
}

type EventBus struct {
	subscriberList SubscriberList
	publishChannel chan Event
}

func (ev *EventBus) Init() {
	ev.publishChannel = make(chan Event, EVENT_PUBLISH_CHANNEL_CAPACITY)
	// nothing for now
	go ev.pubSubLoop()
}

func (ev *EventBus) pubSubLoop() {
	for {
		e := <-ev.publishChannel
		ev.notifySubscribers(e)
	}
}

// Subscribe to listen for game events defined by a query
func (ev *EventBus) Subscribe(
	q EventQuery, name string) EventChannel {

	// Lock the subscriber slice while we modify it
	ev.subscriberList.mutex.Lock()
	defer ev.subscriberList.mutex.Unlock()

	// Create a channel to return to the user
	c := NewEventChannel(name, q)
	eventsDebug("Subscribe: %s on channel %v\n",
		name, c)
	// Add the channel to the subscriber list for its type
	ev.subscriberList.channels[q.Type] = append(
		ev.subscriberList.channels[q.Type], c)
	// return the channel to the caller
	return c
}

// Remove a subscriber
func (ev *EventBus) Unsubscribe(c EventChannel) {
	ev.subscriberList.mutex.Lock()
	defer ev.subscriberList.mutex.Unlock()
	eventsDebug("Unsubscribe on channel %v\n", c)
	removeEventChannelFromSlice(&ev.subscriberList.channels[c.Query.Type], c)
}

// Publish a game event for anyone listening
func (ev *EventBus) Publish(e Event) {
	if len(ev.publishChannel) == EVENT_PUBLISH_CHANNEL_CAPACITY {
		eventsDebug("⚠ ⚠ ⚠ publish channel full. how are you generating so " +
			"many events? Look at EVENT_PUBLISH_CHANNEL_CAPACITY")
	}
	ev.publishChannel <- e

}

// notify subscribers to a certain event
func (ev *EventBus) notifySubscribers(e Event) {
	ev.subscriberList.mutex.RLock()
	defer ev.subscriberList.mutex.RUnlock()
	eventsDebug("⚹: %s\n", e.Description)

	for _, c := range ev.subscriberList.channels[e.Type] {
		if !c.IsActive() {
			eventsDebug("%s channel inactive, "+
				"not sending event %s\n", c.Name, e.Description)
		}
		go func(c EventChannel) {
			// notify if the channel buffer is filled (we're in a
			// goroutine, so it's all good, but probably indicates
			// either the channel receiver is not right, or there's a
			// problem with the rate of events sent over the channel)
			if len(c.C) == EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
				eventsDebug("⚠ subscriber channel #%d for %s is full\n",
					c.Name)
			}
			if c.Query.Test(e) {
				eventsDebug("sending %s on %s",
					e.Description, c.Name)
				c.Send(e)
			}
		}(c)
	}
}
