package engine

type EventQuery struct {
	Type      EventType
	Predicate func(e Event) bool
}

// for simple event queries, predicate is never tested
func (q EventQuery) Test(e Event) bool {
	return q.Type == e.Type && (q.Predicate == nil || q.Predicate(e))
}

// Construct a new EventQuery which only asks about
// the Type of the event
func NewSimpleEventQuery(Type EventType) EventQuery {

	return EventQuery{Type, nil}
}

// Construct a new EventQuery which asks about Type and
// a user-given predicate
func NewPredicateEventQuery(
	Type EventType,
	predicate func(e Event) bool) EventQuery {

	return EventQuery{Type, predicate}
}
