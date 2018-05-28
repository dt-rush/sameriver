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

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) SetActiveStateAtomic(state bool) func(EntityToken) {
	return func(entity EntityToken) {
		// NOTE: we can access the active value directly since this is called
		// exclusively when the entityLock is set (will be reset at the end of
		// the loop iteration in processStateModificationChannel which called
		// this function via one of activate, deactivate, or despawn)
		if m.Components.Active.Data[entity.ID] != state {
			// setActiveState is only called when the entity is locked, so we're
			// good to write directly to the component
			m.Components.Active.Data[entity.ID] = state
			go m.activeEntityLists.notifyActiveState(entity, state)
		}
	}
}

// apply the given tag to the given entity
func (m *EntityManager) TagEntityAtomic(tag string) func(EntityToken) {
	return func(entity EntityToken) {

		// add the tag to the taglist component
		m.Components.TagList.Data[entity.ID].Add(tag)
		// if the entity is active, it has already been checked by all lists,
		// thus generate a new signal to add it to the list of the tag
		if m.Components.Active.Data[entity.ID] {
			m.tags.createTagListIfNeeded(tag)
			m.activeEntityLists.checkActiveEntity(entity)
		}
	}
}

// remove a tag from an entity
func (m *EntityManager) UntagEntityAtomic(tag string) func(EntityToken) {
	return func(entity EntityToken) {
		m.Components.TagList.Data[entity.ID].Remove(tag)
		m.activeEntityLists.checkActiveEntity(entity)
	}
}

// Tag each of the entities in the provided array of ID's with the given tag
func (m *EntityManager) TagEntitiesAtomic(tag string) func([]EntityToken) {
	return func(entities []EntityToken) {

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
