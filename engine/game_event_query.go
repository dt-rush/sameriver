package engine

type GameEventQuery struct {
	Type int
	Predicate func(e GameEvent) bool
}

type GameEventQueryWatcher struct {
	Query   GameEventQuery
	Channel chan(GameEvent)
	ID      uint16
}
