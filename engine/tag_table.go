package engine

import (
	"sync"
)

// used by the EntityManager to tag entities
type TagTable struct {
	// guards all state
	mutex sync.RWMutex
	// data members to support the entity tagging system, which allows us to
	// associate a set of strings with an entity
	// tag -> []IDs
	entitiesWithTag map[string]([]uint16)
	// ID -> []tag
	tagsOfEntity map[uint16]([]string)
}
