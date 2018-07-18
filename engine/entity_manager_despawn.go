package engine

// process the despawn requests in the channel buffer
func (m *EntityManager) processDespawnChannel() {
	n := len(m.despawnSubscription.C)
	for i := 0; i < n; i++ {
		e := <-m.despawnSubscription.C
		m.Despawn(e.Data.(DespawnRequestData).Entity)
	}
}

// User facing function which is used to drain the state of the
// entity manager, and will also kill any pending spawn requests
func (m *EntityManager) DespawnAll() {
	// drain the spawn request channel of pending spawns
	for len(m.spawnSubscription.C) > 0 {
		_ = <-m.spawnSubscription.C
	}
	// iterate all entities which have been allocated and despawn them
	for e, _ := range m.GetCurrentEntitiesSetCopy() {
		m.Despawn(e)
	}
}

// Despawn an entity
func (m *EntityManager) Despawn(e *Entity) {
	if !e.Despawned {
		e.Despawned = true
		m.entityTable.deallocate(e)
		m.entities[e.ID] = nil
		m.setActiveState(e, false)
	}
}

func (m *EntityManager) QueueDespawn(e *Entity) {
	req := DespawnRequestData{Entity: e}
	if len(m.despawnSubscription.C) >= EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
		go func() {
			m.despawnSubscription.C <- Event{DESPAWNREQUEST_EVENT, req}
		}()
	} else {
		m.despawnSubscription.C <- Event{DESPAWNREQUEST_EVENT, req}
	}
}
