/*
 * Functions which produce functions that can be passed to
 * AtomicEntityModify or AtomicEntitiesModify, that is, which produce functions
 * of type func(EntityToken) or func([]EntityToken)
 * Special case for functions which operate directly using the EntityToken,
 * they simply *are* the function which would be output, to prevent patterns
 * like DespawnAtomic()(entity), which just looks cluttered
 *
 * To help reading comprehension, they are always suffixed with "Atomic"
 *
 */

package engine

import ()

// user-facing despawn function which locks the EntityTable for a single
// despawn
func (m *EntityManager) DespawnAtomic(entity EntityToken) {
	m.despawnInternal(entity)
}

// apply the given tag to the given entity
func (m *EntityManager) TagEntityAtomic(tag string) func(EntityToken) {
	return func(entity EntityToken) {
		m.tagTable.tagEntity(entity, tag)
	}
}

// remove a tag from an entity
func (m *EntityManager) UntagEntityAtomic(tag string) func(EntityToken) {
	return func(entity EntityToken) {
		m.tagTable.untagEntity(tag, entity)
	}
}

// Tag each of the entities in the provided array of ID's with the given tag
func (m *EntityManager) TagEntitiesAtomic(tag string) func([]EntityToken) {
	return func(entities []EntityToken) {
		m.tagTable.mutex.Lock()
		defer m.tagTable.mutex.Unlock()

		for _, entity := range entities {
			m.TagEntityAtomic(tag)(entity)
		}
	}
}

// a function which generates an entity modification to augment the health
// of an entity
func (m *EntityManager) IncrementHealthAtomic(change int) func(EntityToken) {
	return func(e EntityToken) {
		healthNow := int(m.Components.Health.Data[e.ID]) + change
		if healthNow > 255 {
			healthNow = 255
		} else if healthNow < 0 {
			healthNow = 0
		}
		m.Components.Health.Data[e.ID] = uint8(healthNow)
	}
}
