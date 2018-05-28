package engine

// used to communicate insert / remove events
type EntitySignalType int

const (
	ENTITY_ADD    = iota
	ENTITY_REMOVE = iota
)

type EntitySignal struct {
	signalType EntitySignalType
	entity     EntityToken
}

type EntityQueryWatcher struct {
	Name  string
	ID    int
	Query EntityQuery
	// A channel along which entity signals will be sent
	Channel chan EntitySignal
}

// Construct a new entity query watcher (its channel will be created at the
// capacity of ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY constant)
func NewEntityQueryWatcher(q EntityQuery) EntityQueryWatcher {

	return EntityQueryWatcher{
		q.Name,
		IDGEN(),
		q,
		make(chan EntitySignal, ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY)}
}
