package engine

// process the despawn requests in the channel buffer
func (m *EntityManager) processDespawnChannel() {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	n := len(m.despawnSubscription.C)
	for i := 0; i < n; i++ {
		e := <-m.despawnSubscription.C
		m.Despawn(e.Data.(DespawnRequestData).Entity)
	}
}

// User facing function which is used to drain the state of the
// entity manager, and will also kill any pending spawn requests
func (m *EntityManager) DespawnAll() {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()
	// drain the spawn and despawn subscription channels of accumulated events
	for len(m.spawnSubscription.C) > 0 {
		_ = <-m.spawnSubscription.C
	}
	for len(m.despawnSubscription.C) > 0 {
		_ = <-m.despawnSubscription.C
	}
	// iterate all IDs which have been allocated and despawn them
	toDespawn := make([]*EntityToken, 0)
	for e, _ := range m.entityTable.currentEntities {
		toDespawn = append(toDespawn, e)
	}
	for _, e := range toDespawn {
		m.Despawn(e)
	}
}

// internal despawn function processes the despawn
// (frees the ID and deactivates the entity)
func (m *EntityManager) Despawn(entity *EntityToken) {
	entity.Despawned = true
	m.entityTable.deallocate(entity)
	m.Entities[entity.ID] = nil
	m.setActiveState(entity, false)
}
