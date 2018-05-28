package engine

import (
	"sync"
)

type EntityQueryWatcherList struct {
	// used to allow systems to keep an updated list of entities which have
	// components they're interested in operating on (eg. physics watches
	// for entities with position, velocity, and hitbox)
	watchers []EntityQueryWatcher
	// to protect modifying the above slice
	mutex sync.RWMutex
}
