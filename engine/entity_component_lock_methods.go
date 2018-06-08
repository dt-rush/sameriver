package engine

func (m *EntityManager) lockEntityComponent(
	entity EntityToken, component ComponentType) {

	// starts an RLock() on the accessLock for the component
	m.Components.accessStart(component)
	// Lock() the valueLock for this entity on this component
	m.Components.valueLocks[component][entity.ID].Lock()
	entityLocksDebug("LOCKED entity %v component %s",
		entity, COMPONENT_NAMES[component])
}

func (m *EntityManager) unlockEntityComponent(
	entity EntityToken, component ComponentType) {

	entityLocksDebug("UNLOCK entity %v component %s",
		entity, COMPONENT_NAMES[component])
	// Unlock() the valueLock for this entity on this component
	m.Components.valueLocks[component][entity.ID].Unlock()
	// ends the RLock() on the accessLock for the component
	m.Components.accessEnd(component)
}

func (m *EntityManager) lockOneEntityComponents(
	entityComponents OneEntityComponents) {

	// lock components in sorted order
	sort.Slice(entityComponents.components, func(i int, j int) {
		return entityComponents.components[i] < entityComponents.components[j]
	})
	for _, component := range entityComponents.components {
		m.lockEntityComponent(entityComponents.entity, component)
	}
}

func (m *EntityManager) lockEntitiesComponents(
	entitiesComponents []OneEntityComponents) {

	// lock entities in sorted order
	sort.Slice(entitiesComponents, func(i int, j int) {
		return entitiesComponents[i].entity.ID < entitiesComponents[j].entity.ID
	})
	for _, entityComponents := range entitiesComponents {
		m.lockEntityComponents(entityComponents)
	}
}

func (m *EntityManager) unlockEntitiesComponents(
	entitiesComponents []OneEntityComponents) {
	for _, entityComponents := range entitiesComponents {
		for _, component := range entityComponents.components {
			m.unlockEntityComponent(entityComponents.entity, component)
		}
	}
}

func (m *EntityManager) unlockOneEntityComponents(
	entityComponents OneEntityComponents) {
	for _, entityComponent := range entityComponents {
		m.unlockEntityComponent(
			entityComponent.entity,
			entityComponent.component)
	}
}

func (m *EntityManager) unlockEntityComponents(
	entityComponents []EntityComponent) {
	for _, entityComponent := range entityComponents {
		m.unlockEntityComponent(
			entityComponent.entity,
			entityComponent.component)
	}
}

func (m *EntityManager) rLockEntityComponent(
	entity EntityToken, component ComponentType) {

	// starts an RLock() on the accessLock for the component
	m.Components.accessStart(component)
	// RLock() the valueLock for this entity on this component
	m.Components.valueLocks[component][entity.ID].RLock()
	entityLocksDebug("RLOCKED entity %v component %s",
		entity, COMPONENT_NAMES[component])
}

func (m *EntityManager) rUnlockEntityComponent(
	entity EntityToken, component ComponentType) {

	entityLocksDebug("RUNLOCK entity %v component %s",
		entity, COMPONENT_NAMES[component])
	// RUnlock() the valueLock for this entity on this component
	m.Components[component].locks[entity.ID].RUnlock()
	// ends the RLock() on the accessLock for the component
	m.Components.accessEnd(component)
}
