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
	// iterate all IDs which could have been allocated and despawn them
	// (each time a despawn goes through, the entityTable.currentEntities list
	// will shrink)
	for len(m.entityTable.currentEntities) > 0 {
		// for each allocated ID, build a token based on the current gen
		// lockEntity here will always return true because
		// the gen will never mismatch, since only a processDespawn() call
		// could change that, and that occurs only in two places: here, where
		// we iterate one despawn for each entity while the entityTable is
		// locked, and in a call to Despawn() (in entity_modifications.go)
		// from a user which would only be able to proceed after we
		// released the entityTable lock, and would then exit since gen
		// had changed
		m.Despawn(m.entityTable.currentEntities[0])
	}
	// drain the spawn channel
	for len(m.spawnSubscription.C) > 0 {
		// we're draining the channel, so do nothing
		_ = <-m.spawnSubscription.C
	}
}

// internal despawn function processes the despawn
// (frees the ID and deactivates the entity)
func (m *EntityManager) Despawn(entity *EntityToken) {
	entity.Despawned = true
	m.entityTable.numEntities--
	m.entityTable.availableIDs = append(m.entityTable.availableIDs, entity.ID)
	removeEntityTokenFromSlice(&m.entityTable.currentEntities, entity)
	m.setActiveState(entity, false)
}
