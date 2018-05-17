package engine

type GameEventQueryWatcher struct {
	Query   GameEventQuery
	Channel GameEventChannel
	Name    string
}

type GameEventQuery interface {
	Test(e GameEvent) bool
}

type simpleGameEventQuery struct {
	// the class of event
	class int
}

func (q simpleGameEventQuery) Test(e GameEvent) bool {
	return q.class == e.class
}

type predicateGameEventQuery struct {
	// the class of event
	class int
	// a function which predicates the game event more
	// specifically than its class
	predicate func(e GameEvent) bool
}

func (q predicateGameEventQuery) Test(e GameEvent) bool {
	return q.class == e.class && q.predicate(e)
}

// Construct a new game event query which only asks about
// the class of the event
func NewSimpleGameEventQuery(class int) GameEventQuery {

	return &simpleGameEventQuery{class}
}

// Construct a new game event query which asks about class and
// a user-given predicate
func NewPredicateGameEventQuery(
	class int, predicate func(e GameEvent) bool) GameEventQuery {

	return &predicateGameEventQuery{class, predicate}
}
