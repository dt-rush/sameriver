package engine

import (
	"github.com/dt-rush/sameriver/engine/utils"
)

// used to communicate insert / remove events
type EntitySignalType int

const (
	ENTITY_ADD    = iota
	ENTITY_REMOVE = iota
)

type EntitySignal struct {
	SignalType EntitySignalType
	Entity     *EntityToken
}

type EntityQueryWatcher struct {
	Name    string
	ID      int
	Query   EntityQuery
	Channel chan EntitySignal
}

// Construct a new entity query watcher (its channel will be created at the
// capacity of ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY constant)
func NewEntityQueryWatcher(q EntityQuery) EntityQueryWatcher {

	return EntityQueryWatcher{
		q.Name,
		utils.IDGEN(),
		q,
		make(chan EntitySignal, ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY)}
}
