package engine

import (
	"errors"
	"fmt"
	"time"
)

// process the spawn requests in the channel buffer
func (m *EntityManager) processSpawnChannel() {
	// get the current number of requests in the channel and only process
	// them. More may continue to pile up. They'll get processed next time.
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	n := len(m.spawnSubscription.C)
	for i := 0; i < n; i++ {
		// get the request from the channel
		e := <-m.spawnSubscription.C
		_, err := m.Spawn(e.Data.(SpawnRequestData))
		if err != nil {
			go func() {
				time.Sleep(5 * time.Second)
				m.eventBus.Publish(SPAWNREQUEST_EVENT, e.Data)
			}()
		}
	}
}

func (m *EntityManager) Spawn(r SpawnRequestData) (*EntityToken, error) {
	return m.spawn(r, "")
}

func (m *EntityManager) SpawnUnique(
	tag string, r SpawnRequestData) (*EntityToken, error) {

	if _, ok := m.uniqueEntities[tag]; ok {
		return nil, errors.New(fmt.Sprintf("requested to spawn unique "+
			"entity for %s, but %s already exists", tag, tag))
	}
	e, err := m.spawn(r, tag)
	if err == nil {
		m.uniqueEntities[tag] = e
	}
	return e, err
}

// given a list of components, spawn an entity with the default values
// returns the EntityToken (used to spawn an entity for which we *want* the
// token back)

func (m *EntityManager) spawn(r SpawnRequestData, uniqueTag string) (
	*EntityToken, error) {

	// used if spawn is impossible for various reasons
	var fail = func(msg string) (*EntityToken, error) {
		return nil, errors.New(msg)
	}

	// get an ID for the entity
	entity, err := m.entityTable.allocateID()
	if err != nil {
		errorMsg := fmt.Sprintf("⚠ Error in allocateID(): %s. Will not spawn "+
			"entity with tags: %v\n", err, r.Components.TagList.Tags)
		return fail(errorMsg)
	}
	// print a debug message
	// set the bitarray for this entity
	entity.ComponentBitArray = r.Components.ToBitArray()
	// copy the data inNto the component storage for each component
	// (note: we dereference the pointers, this is a real copy, so it's good
	// that component values are either small pieces of data like [2]uint16
	// or a pointer to a func, etc.).
	// We don't "zero" the values of components not in the entity's set,
	// because if a system operating on the component data
	// expects to work on the data, it should be maintaining a list of
	// entities with the required components using an UpdatedEntityList
	// NOTE: we can directly set the Active component value since no other
	// goroutine could be also writing to this entity, due to the
	// AtomicEntityModify pattern
	m.ApplyComponentSet(r.Components)(entity)
	// create (if doesn't exist) entitiesWithTag lists for each tag
	if r.Components.TagList != nil {
		for _, tag := range r.Components.TagList.Tags {
			m.createEntitiesWithTagListIfNeeded(tag)
		}
	}
	// apply the unique tag if provided
	if uniqueTag != "" {
		m.createEntitiesWithTagListIfNeeded(uniqueTag)
	}
	// add the entity to the list of current entities
	m.entityTable.addToCurrentEntities(entity)
	// set entity active and notify entity is active
	m.setActiveState(entity, true)
	// return EntityToken
	return entity, nil
}
