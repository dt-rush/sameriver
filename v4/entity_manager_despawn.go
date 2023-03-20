package sameriver

// User facing function which is used to drain the state of the
// entity manager, and will also kill any pending spawn requests
func (m *EntityManager) DespawnAll() {
	// drain the spawn request channel of pending spawns
	for len(m.spawnSubscription.C) > 0 {
		<-m.spawnSubscription.C
	}
	// iterate all entities which have been allocated and despawn them
	for e := range m.GetCurrentEntitiesSetCopy() {
		m.Despawn(e)
	}
}

// Despawn an entity
func (m *EntityManager) Despawn(e *Entity) {
	// guard against multiple logics per tick despawning an entity
	if !e.Despawned {
		e.Despawned = true
		m.entityIDAllocator.deallocate(e)
		e.RemoveAllLogics()
		m.setActiveState(e, false)
	}
}
