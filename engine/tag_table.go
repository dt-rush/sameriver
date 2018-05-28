package engine

import (
	"sync"
)

// used by the EntityManager to tag entities
type TagTable struct {
	mutex           sync.RWMutex
	tagsOfEntity    [MAX_ENTITIES][]string
	entitiesWithTag map[string]*UpdatedEntityList
}

func (t *TagTable) Init() {
	t.entitiesWithTag = make(map[string]*UpdatedEntityList)
}

func (t *TagTable) NumEntitiesWithTag(tag string) int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if _, exists := t.entitiesWithTag[tag]; !exists {
		return 0
	} else {
		return t.entitiesWithTag[tag].Length()
	}
}
