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

func (t *TagTable) tagEntity(tag string, entity EntityToken) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	tagsDebug("adding tag \"%s\" to %v", tag, entity)

	t.tagsOfEntity[entity.ID] = append(t.tagsOfEntity[entity.ID], tag)
	t.createTagListIfNeeded(tag)
	t.entitiesWithTag[tag].actOnEntitySignal(EntitySignal{ENTITY_ADD, entity})
}

func (t *TagTable) untagEntity(tag string, entity EntityToken) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	tagsDebug("removing tag \"%s\" from %v", tag, entity)

	t.tagsOfEntity[entity.ID] = []string{}
	if _, exists := t.entitiesWithTag[tag]; exists {
		t.entitiesWithTag[tag].actOnEntitySignal(EntitySignal{ENTITY_REMOVE, entity})
	}
}

func (t *TagTable) entitiesWithTag(tag string) *UpdatedEntityList {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.createTagListIfNeeded(tag)
	return t.entitiesWithTag[tag]
}

func (t *TagTable) createTagListIfNeeded(tag string) {
	if _, exists := t.entitiesWithTag[tag]; !exists {
		t.entitiesWithTag[tag] = t.em.GetUpdatedActiveEntityList(
			tag, EntityQueryFromTag(tag))
	}
}
