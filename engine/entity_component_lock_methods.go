package engine

// wait until we can lock an entity on a specific component and check
// gen-match after lock. If gen mismatch, release and return false.
// else return true
func (m *EntityManager) lockEntityComponent(
	entity EntityToken, component ComponentType) bool {

	// starts an RLock() on the accessLock for the component
	m.Components.accessStart(component)

	// Lock() the valueLock for this entity on this component
	m.Components.valueLocks[component][entity.ID].Lock()
	if !m.entityTable.genValidate(entity) {
		entityLocksDebug("GENMISMATCH entity %v", entity)
		m.Components.valueLocks[component][entity.ID].Unlock()
		return false
	}
	entityLocksDebug("LOCKED entity %v component %s",
		entity, COMPONENT_NAMES[component])
	return true
}

func (m *EntityManager) releaseEntityComponent(
	entity EntityToken, component ComponentType) {

	entityLocksDebug("RELEASING entity %v component %s",
		entity, COMPONENT_NAMES[component])
	m.Components[component].locks[entity.ID].Unlock()
	// ends the RLock() on the accessLock for the component
	m.Components.accessEnd(component)
}

func (m *EntityManager) releaseEntityComponents(
	entityComponents []EntityComponent) {
	for _, entityComponent := range entityComponents {
		m.releaseEntityComponent(
			entityComponent.entity,
			entityComponent.component)
	}
}
