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
	// Each EventQuery's Predicate will be tested against the events
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

// Subscribe to listen for game events defined by a query
func (ev *EventBus) Subscribe(
	name string, q *EventQuery) *EventChannel {

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
	removeEventChannelFromSlice(&ev.subscriberList.channels[c.Query.Type], c)
}

// notify subscribers to a certain event
func (ev *EventBus) notifySubscribers(e Event) {

	// TODO: generate a means of printing events and create a special
	// system which listens on *all* events, printing them
	eventsLog("⚹: %s\n", EVENT_NAMES[e.Type])

	for _, c := range ev.subscriberList.channels[e.Type] {
		if !c.IsActive() {
			continue
		}
		// notify if the channel buffer is filled (we're in a
		// goroutine, so it's all good, but probably indicates
		// either the channel receiver is not right, or there's a
		// problem with the rate of events sent over the channel)
		if len(c.C) == EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
			Logger.Printf("⚠ event subscriber channel #%s for events of "+
				"type %s is full\n", c.Name, EVENT_NAMES[e.Type])
		}
		if c.Query.Test(e) {
			c.Send(e)
		}
	}
}
