package engine

import (
	"sort"
)

func (m *EntityManager) lockEntityComponent(
	entity EntityToken, component ComponentType) {

	// starts an RLock() on the accessLock for the component
	m.Components.accessStart(component)
	// Lock() the valueLock for this entity on this component
	m.Components.valueLocks[component][entity.ID].Lock()
}

func (m *EntityManager) unlockEntityComponent(
	entity EntityToken, component ComponentType) {

	// Unlock() the valueLock for this entity on this component
	m.Components.valueLocks[component][entity.ID].Unlock()
	// ends the RLock() on the accessLock for the component
	m.Components.accessEnd(component)
}

func (m *EntityManager) lockEntityComponents(
	entityComponents EntityComponents) {

	// lock components in sorted order
	sort.Slice(entityComponents.Components, func(i int, j int) bool {
		return entityComponents.Components[i] < entityComponents.Components[j]
	})
	for _, component := range entityComponents.Components {
		m.lockEntityComponent(entityComponents.Entity, component)
	}
}

func (m *EntityManager) lockEntitiesComponents(
	entitiesComponents []EntityComponents) {

	// lock entities in sorted order
	sort.Slice(entitiesComponents, func(i int, j int) bool {
		return entitiesComponents[i].Entity.ID < entitiesComponents[j].Entity.ID
	})
	for _, entityComponents := range entitiesComponents {
		m.lockEntityComponents(entityComponents)
	}
}

func (m *EntityManager) unlockEntitiesComponents(
	entitiesComponents []EntityComponents) {
	for _, entityComponents := range entitiesComponents {
		for _, component := range entityComponents.Components {
			m.unlockEntityComponent(entityComponents.Entity, component)
		}
	}
}

func (m *EntityManager) unlockEntityComponents(
	entityComponents EntityComponents) {
	for _, component := range entityComponents.Components {
		m.unlockEntityComponent(
			entityComponents.Entity,
			component)
	}
}

func (m *EntityManager) rLockEntityComponent(
	entity EntityToken, component ComponentType) {

	// starts an RLock() on the accessLock for the component
	m.Components.accessStart(component)
	// RLock() the valueLock for this entity on this component
	m.Components.valueLocks[component][entity.ID].RLock()
}

func (m *EntityManager) rUnlockEntityComponent(
	entity EntityToken, component ComponentType) {

	// RUnlock() the valueLock for this entity on this component
	m.Components.valueLocks[component][entity.ID].RUnlock()
	// ends the RLock() on the accessLock for the component
	m.Components.accessEnd(component)
}
