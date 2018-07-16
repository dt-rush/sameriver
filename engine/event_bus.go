/**
  *
  * Pub-sub hub for game events
  *
**/

package engine

type SubscriberList struct {
	// subscriberLists is a list of lists of EventChannels
	// where the outer list is indexed by the EventType (type aliased
	// to int). So you could have a list of queries on CollisionEvents, etc.
	// Each EventFilter's Predicate will be tested against the events
	// that are published for the matching type (and thus the predicates
	// can safely assert the type of the Data member of the event)
	channels [N_EVENT_TYPES][]*EventChannel
}

type EventBus struct {
	subscriberList SubscriberList
	publishChannel chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (ev *EventBus) Publish(Type EventType, Data interface{}) {
	go ev.notifySubscribers(Event{Type, Data})
}

// Subscribe to listen for game events defined by a Filter
func (ev *EventBus) Subscribe(
	name string, q *EventFilter) *EventChannel {

	// Create a channel to return to the user
	c := NewEventChannel(name, q)
	// Add the channel to the subscriber list for its type
	ev.subscriberList.channels[q.Type] = append(
		ev.subscriberList.channels[q.Type], c)
	// return the channel to the caller
	return c
}

// Remove a subscriber
func (ev *EventBus) Unsubscribe(c *EventChannel) {
	removeEventChannelFromSlice(&ev.subscriberList.channels[c.Filter.Type], c)
}

// notify subscribers to a certain event
func (ev *EventBus) notifySubscribers(e Event) {
	// TODO: create a special system which listens on *all* events,
	// printing them if it's turned on
	eventsLog("⚹: %s\n", EVENT_NAMES[e.Type])

	var notifyFull = func(c *EventChannel) {
		Logger.Printf("⚠ event subscriber channel #%s for events of "+
			"type %s is full\n", c.Name, EVENT_NAMES[e.Type])
	}

	for _, c := range ev.subscriberList.channels[e.Type] {
		if !c.IsActive() {
			continue
		}
		if c.Filter.Test(e) {
			if len(c.C) >= EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
				notifyFull(c)
				// spawn a goroutine to do the channel send since we don't
				// want a hang here to affect other subscribers
				go func() {
					c.C <- e
				}()
			} else {
				c.C <- e
			}
		}
	}
}
