package engine

// process the despawn requests in the channel buffer
func (m *EntityManager) processDespawnChannel() {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	n := len(m.despawnSubscription.C)
	for i := 0; i < n; i++ {
		e := <-m.despawnSubscription.C
		m.processDespawn(e.Data.(DespawnRequestData).Entity)
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
		m.processDespawn(m.entityTable.currentEntities[0])
	}
	// drain the spawn channel
	for len(m.spawnSubscription.C) > 0 {
		// we're draining the channel, so do nothing
		_ = <-m.spawnSubscription.C
	}
}

func (m *EntityManager) Despawn(entity *EntityToken) {
	// despawn is idempotent
	if m.entityTable.despawnFlags[entity.ID] == 0 {
		m.entityTable.despawnFlags[entity.ID] = 1
		m.ev.Publish(DESPAWNREQUEST_EVENT, DespawnRequestData{entity})
	}
}

// internal despawn function which assumes the EntityTable is locked
func (m *EntityManager) processDespawn(entity *EntityToken) {

	// decrement the entity count
	m.entityTable.numEntities--
	// add the ID to the list of available IDs
	m.entityTable.availableIDs = append(m.entityTable.availableIDs, entity.ID)
	// remove the entityfrom the list of current entities
	removeEntityTokenFromSlice(&m.entityTable.currentEntities, entity)

	//  and notify
	m.setActiveState(entity, false)
	// delete the entity's logic
	if _, exists := m.Logics[entity]; exists {
		delete(m.Logics, entity)
	}
}
