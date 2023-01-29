package engine

type EventFilter struct {
	eventType string
	predicate func(e Event) bool
}

// for simple event queries, predicate is never tested
func (q *EventFilter) Test(e Event) bool {
	return q.eventType == e.Type && (q.predicate == nil || q.predicate(e))
}

// Construct a new EventFilter which only asks about
// the Type of the event
func SimpleEventFilter(Type string) *EventFilter {
	return &EventFilter{Type, nil}
}

// Construct a new EventFilter which asks about Type and
// a user-given predicate
func PredicateEventFilter(
	Type string, predicate func(e Event) bool) *EventFilter {

	return &EventFilter{Type, predicate}
}
