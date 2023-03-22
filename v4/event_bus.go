/**
  *
  * Pub-sub hub for game events
  *
**/

package sameriver

import (
	"go.uber.org/atomic"
)

type Event struct {
	Type string
	Data any
}

type EventBus struct {
	name string
	// channels is a map of lists of EventChannels
	// Each EventFilter's Predicate will be tested against the events
	// that are published for the matching type (and thus the predicates
	// can safely assert the type of the Data member of the event)
	channels map[string][]*EventChannel
	// number of goroutines spawned to publish events to subscriber channels
	// that are full
	nHanging atomic.Int32
}

func NewEventBus(name string) *EventBus {
	b := &EventBus{name: name}
	b.channels = make(map[string][]*EventChannel)
	return b
}

func (b *EventBus) Publish(t string, data any) {
	b.notifySubscribers(Event{t, data})
}

// Subscribe to listen for game events defined by a Filter
func (b *EventBus) Subscribe(q *EventFilter) *EventChannel {

	// Create a channel to return to the user
	c := NewEventChannel(q)
	// Add the channel to the subscriber list for its type
	b.channels[q.eventType] = append(
		b.channels[q.eventType], c)
	// return the channel to the caller
	return c
}

// Remove a subscriber
func (b *EventBus) Unsubscribe(c *EventChannel) {
	eventType := c.filter.eventType
	channels, ok := b.channels[eventType]
	if ok {
		channels = removeEventChannelFromSlice(channels, c)
		b.channels[eventType] = channels
	}
}

// notify subscribers to a certain event
func (b *EventBus) notifySubscribers(e Event) {
	// TODO: create a special system which listens on *all* events,
	// printing them if it's turned on
	logEvents("⚹: %s", e.Type)

	var notifyFull = func(c *EventChannel) {
		logWarning("event subscriber channel for events of type %s is full; possibly sending too many events; consider throttling or increase capacity\n", e.Type)
	}

	var notifyExtraFull = func() {
		logWarning("/!\\ /!\\ /!\\ number of goroutines waiting for an event channel (of event type %s) to go below max capacity is now greater than capacity (%d); you're sending too many events", e.Type, EVENT_SUBSCRIBER_CHANNEL_CAPACITY)
	}

	logEvents("len(b.channels[e.Type])=%d", len(b.channels[e.Type]))
	for _, c := range b.channels[e.Type] {
		logEvents("| Channel: %p", c)
		if !c.IsActive() {
			continue
		}
		logEvents("--> channel.filter.Test(e) = %t", c.filter.Test(e))
		if c.filter.Test(e) {
			if len(c.C) >= EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
				notifyFull(c)
				// spawn a goroutine to do the channel send since we don't
				// want a hang here to affect other subscribers
				// (note: if you severely overrun, even these goroutines
				// will add up and cause problems)
				// TODO: count how many of these are waiting and warn if too high
				go func() {
					b.nHanging.Add(1)
					if b.nHanging.Load() > EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
						notifyExtraFull()
					}
					logEvents("---- event channel put <- %s.%v", e.Type, e.Data)
					c.C <- e
					b.nHanging.Add(-1)
				}()
			} else {
				logEvents("---- event channel put <- %s.%v", e.Type, e.Data)
				c.C <- e
				logEvents("---- len(<%p>.C) = %d", c, len(c.C))
			}
		}
	}
}
