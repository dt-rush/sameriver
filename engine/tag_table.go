package engine

import (
	"sync"
)

// used by the EntityManager to tag entities
type TagTable struct {
	em              *EntityManager
	mutex           sync.RWMutex
	entitiesWithTag map[string]*UpdatedEntityList
}

func (t *TagTable) Init(em *EntityManager) {
	t.em = em
	t.entitiesWithTag = make(map[string]*UpdatedEntityList)
}

func (t *TagTable) GetEntitiesWithTag(tag string) *UpdatedEntityList {
	t.createEntitiesWithTagListIfNeeded(tag)
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.entitiesWithTag[tag]
}

func (t *TagTable) createEntitiesWithTagListIfNeeded(tag string) {
	t.mutex.RLock()
	_, exists := t.entitiesWithTag[tag]
	t.mutex.RUnlock()
	if !exists {
		// NOTE: when we seize the lock below, another routine may have already
		// come through here since we hit RUnlock and tested the !exists condition.
		// thankfully GetUpdatedEntityList itself will return the same list if it
		// was already created, so we'll just write the same list to the map
		t.mutex.Lock()
		t.entitiesWithTag[tag] =
			t.em.GetUpdatedEntityList(EntityQueryFromTag(tag))
		t.mutex.Unlock()
	}
}
