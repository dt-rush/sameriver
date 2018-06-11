package engine

// process the despawn requests in the channel buffer
func (m *EntityManager) processDespawnChannel() {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	n := len(m.despawnSubscription.C)
	for i := 0; i < n; i++ {
		e := <-m.despawnSubscription.C
		m.processDespawn(e)
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
		entity := m.entityTable.currentEntities[0]
		m.processDespawn(entity)
	}
	// drain the spawn channel
	for len(m.spawnSubscription) > 0 {
		// we're draining the channel, so do nothing
		_ = <-m.spawnSubscription
	}
}

func (m *EntityManager) Despawn(entity EntityToken) {
	// despawn is idempotent
	if m.entityTable.despawnFlags[entity.ID] == 0 {
		m.entityTable.despawnFlags[entity.ID] = 1
		m.ev.Publish(DESPAWN_EVENT, DespawnRequestData{entity})
	}
}

// internal despawn function which assumes the EntityTable is locked
func (m *EntityManager) processDespawn(entity EntityToken) {

	// if the gen doesn't match, another despawn for this same entity
	// has already been through here (if a regular Despawn() call and
	// a DespawnAll() were racing)
	if !m.entityTable.genValidate(entity) {
		return
	}
	// decrement the entity count
	m.entityTable.numEntities--
	// add the ID to the list of available IDs
	m.entityTable.availableIDs = append(m.entityTable.availableIDs, entity.ID)
	// remove the ID from the list of allocated IDs
	removeEntityTokenFromSlice(&m.entityTable.currentEntities, entity)
	// Increment the gen for the ID
	// NOTE: it's important that we increment gen before resetting the
	// locks, since any goroutines waiting for the
	// lock to be 0 so they can claim it in AtomicEntityModify() will then
	// immediately want to check if the gen of the entity still matches.
	m.entityTable.incrementGen(entity.ID)

	//  and notify
	m.setActiveState(entity, false)
	// delete the entity's logic (we have to do this *after* stopping it)
	m.EntityLogicTable.deleteLogic(entity)
}
