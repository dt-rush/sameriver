package engine

import (
	"errors"
	"fmt"
)

// get the current number of requests in the channel and only process
// them. More may continue to pile up. They'll get processed next time.
func (m *EntityManager) processSpawnChannel() {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	n := len(m.spawnSubscription.C)
	for i := 0; i < n; i++ {
		// get the request from the channel
		e := <-m.spawnSubscription.C
		_, err := m.spawn(e.Data.(SpawnRequestData))
		if err != nil {
			Logger.Println(err)
		}
	}
}

func (m *EntityManager) spawn(r SpawnRequestData) (*Entity, error) {
	return m.doSpawn(r, "")
}

func (m *EntityManager) spawnUnique(
	tag string, r SpawnRequestData) (*Entity, error) {

	if _, ok := m.uniqueEntities[tag]; ok {
		return nil, errors.New(fmt.Sprintf("requested to spawn unique "+
			"entity for %s, but %s already exists", tag, tag))
	}
	e, err := m.doSpawn(r, tag)
	if err == nil {
		m.uniqueEntities[tag] = e
	}
	return e, err
}

// given a list of components, spawn an entity with the default values
// returns the Entity (used to spawn an entity for which we *want* the
// token back)

func (m *EntityManager) doSpawn(r SpawnRequestData, uniqueTag string) (
	*Entity, error) {

	// used if spawn is impossible for various reasons
	var fail = func(msg string) (*Entity, error) {
		return nil, errors.New(msg)
	}

	// get an ID for the entity
	e, err := m.entityTable.allocateID()
	if err != nil {
		errorMsg := fmt.Sprintf("⚠ Error in allocateID(): %s. Will not spawn "+
			"entity with tags: %v\n", err, r.Tags)
		return fail(errorMsg)
	}
	e.World = m.w
	// add the entity to the list of current entities
	m.Entities[e.ID] = e
	// set the bitarray for this entity
	e.ComponentBitArray = r.Components.ToBitArray()
	// copy the data inNto the component storage for each component
	// (note: we dereference the pointers, this is copy operation, so it's good
	// that component values are either small pieces of data like [2]uint16
	// or a pointer to a func, etc.).
	// We don't "zero" the values of components not in the entity's set,
	// because if a system operating on the component data
	// expects to work on the data, it should be maintaining a list of
	// entities with the required components using an UpdatedEntityList
	// NOTE: we can directly set the Active component value since no other
	// goroutine could be also writing to this entity, due to the
	// AtomicEntityModify pattern
	m.ApplyComponentSet(r.Components)(e)
	// create (if doesn't exist) entitiesWithTag lists for each tag
	m.TagEntity(e, r.Tags...)
	// apply the unique tag if provided
	if uniqueTag != "" {
		m.createEntitiesWithTagListIfNeeded(uniqueTag)
	}
	// set entity active and notify entity is active
	m.setActiveState(e, true)
	// return Entity
	return e, nil
}
