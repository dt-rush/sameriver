package engine

type EventFilter struct {
	eventType EventType
	predicate func(e Event) bool
}

// for simple event queries, predicate is never tested
func (q *EventFilter) Test(e Event) bool {
	return q.eventType == e.Type && (q.predicate == nil || q.predicate(e))
}

// Construct a new EventFilter which only asks about
// the Type of the event
func SimpleEventFilter(Type EventType) *EventFilter {
	return &EventFilter{Type, nil}
}

// Construct a new EventFilter which asks about Type and
// a user-given predicate
func PredicateEventFilter(
	Type EventType, predicate func(e Event) bool) *EventFilter {

	return &EventFilter{Type, predicate}
}

// Construct a new EventFilter given a function from the user
func CollisionEventFilter(test func(c CollisionData) bool) *EventFilter {
	return &EventFilter{
		COLLISION_EVENT,
		func(e Event) bool {
			c := e.Data.(CollisionData)
			return test(c)
		},
	}
}
