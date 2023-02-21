/**
  *
  * Pub-sub hub for game events
  *
**/

package sameriver

import (
	"go.uber.org/atomic"
)

type SubscriberList struct {
	// subscriberLists is a list of lists of EventChannels
	// where the outer list is indexed by the EventType (type aliased
	// to int). So you could have a list of queries on CollisionEvents, etc.
	// Each EventFilter's Predicate will be tested against the events
	// that are published for the matching type (and thus the predicates
	// can safely assert the type of the Data member of the event)
	channels map[string][]*EventChannel
}

type Event struct {
	Type string
	Data interface{}
}

type EventBus struct {
	subscriberList SubscriberList
	publishChannel chan Event
	// number of goroutines spawned to publish events to subscriber channels
	// that are full
	nHanging atomic.Int32
}

func NewEventBus() *EventBus {
	b := &EventBus{}
	b.subscriberList.channels = make(map[string][]*EventChannel)
	return b
}

func (ev *EventBus) Publish(t string, data interface{}) {
	ev.notifySubscribers(Event{t, data})
}

// Subscribe to listen for game events defined by a Filter
func (ev *EventBus) Subscribe(q *EventFilter) *EventChannel {

	// Create a channel to return to the user
	c := NewEventChannel(q)
	// Add the channel to the subscriber list for its type
	ev.subscriberList.channels[q.eventType] = append(
		ev.subscriberList.channels[q.eventType], c)
	// return the channel to the caller
	return c
}

// Remove a subscriber
func (ev *EventBus) Unsubscribe(c *EventChannel) {
	eventType := c.filter.eventType
	channels, ok := ev.subscriberList.channels[eventType]
	if ok {
		channels = removeEventChannelFromSlice(channels, c)
		ev.subscriberList.channels[eventType] = channels
	}
}

// notify subscribers to a certain event
func (ev *EventBus) notifySubscribers(e Event) {
	// TODO: create a special system which listens on *all* events,
	// printing them if it's turned on
	eventsLog("⚹: %s\n", e.Type)

	var notifyFull = func(c *EventChannel) {
		Logger.Printf("[WARNING] event subscriber channel for events of type %s is full\n", e.Type)
	}

	var notifyExtraFull = func() {
		Logger.Printf("[WARNING] number of goroutines waiting for event channel (of event type %s) to go below max capacity is now greater than capacity (%d)", e.Type, EVENT_SUBSCRIBER_CHANNEL_CAPACITY)
	}

	for _, c := range ev.subscriberList.channels[e.Type] {
		if !c.IsActive() {
			continue
		}
		if c.filter.Test(e) {
			if len(c.C) >= EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
				notifyFull(c)
				// spawn a goroutine to do the channel send since we don't
				// want a hang here to affect other subscribers
				// (note: if you severely overrun, even these goroutines
				// will add up and cause problems)
				// TODO: count how many of these are waiting and warn if too high
				go func() {
					ev.nHanging.Add(1)
					if ev.nHanging.Load() > EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
						notifyExtraFull()
					}
					c.C <- e
					ev.nHanging.Add(-1)
				}()
			} else {
				c.C <- e
			}
		}
	}
}
