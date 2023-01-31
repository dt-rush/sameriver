package sameriver

import (
	"errors"
	"fmt"
)

type SpawnRequestData struct {
	Components ComponentSet
	Tags       []string
	UniqueTag  string
}

// get the current number of requests in the channel and only process
// them. More may continue to pile up. They'll get processed next time.
func (m *EntityManager) processSpawnChannel() {
	n := len(m.spawnSubscription.C)
	for i := 0; i < n; i++ {
		// get the request from the channel
		e := <-m.spawnSubscription.C
		req := e.Data.(SpawnRequestData)
		_, err := m.Spawn(req.Tags, req.Components)
		if err != nil {
			Logger.Println(err)
		}
	}
}

func (m *EntityManager) Spawn(tags []string,
	components ComponentSet) (*Entity, error) {
	return m.doSpawn("", tags, components)
}

func (m *EntityManager) queueSpawn(req SpawnRequestData) {
	if len(m.spawnSubscription.C) >= EVENT_SUBSCRIBER_CHANNEL_CAPACITY {
		go func() {
			m.spawnSubscription.C <- Event{"spawn-request", req}
		}()
	} else {
		m.spawnSubscription.C <- Event{"spawn-request", req}
	}
}

func (m *EntityManager) QueueSpawn(tags []string, components ComponentSet) {
	m.queueSpawn(SpawnRequestData{
		Tags:       tags,
		Components: components,
	})
}

func (m *EntityManager) QueueSpawnUnique(
	uniqueTag string, tags []string, components ComponentSet) {
	m.queueSpawn(SpawnRequestData{
		UniqueTag:  uniqueTag,
		Tags:       tags,
		Components: components,
	})
}

func (m *EntityManager) SpawnUnique(
	tag string, tags []string, components ComponentSet) (*Entity, error) {

	if _, ok := m.uniqueEntities[tag]; ok {
		return nil, errors.New(fmt.Sprintf("requested to spawn unique "+
			"entity for %s, but %s already exists", tag, tag))
	}
	e, err := m.doSpawn(tag, tags, components)
	if err == nil {
		m.uniqueEntities[tag] = e
	}
	return e, err
}

// given a list of components, spawn an entity with the default values
// returns the Entity (used to spawn an entity for which we *want* the
// token back)

func (m *EntityManager) doSpawn(
	uniqueTag string, tags []string, components ComponentSet) (
	*Entity, error) {

	// used if spawn is impossible for various reasons
	var fail = func(msg string) (*Entity, error) {
		return nil, errors.New(msg)
	}

	// get an ID for the entity
	e, err := m.entityTable.allocateID()
	if err != nil {
		errorMsg := fmt.Sprintf("⚠ Error in allocateID(): %s. Will not spawn "+
			"entity with tags: %v\n", err, tags)
		return fail(errorMsg)
	}
	e.World = m.w
	// add the entity to the list of current entities
	m.entities[e.ID] = e
	// set the bitarray for this entity
	e.ComponentBitArray = m.components.BitArrayFromComponentSet(components)
	// copy the data inNto the component storage for each component
	m.components.ApplyComponentSet(e, components)
	// create (if doesn't exist) entitiesWithTag lists for each tag
	m.TagEntity(e, tags...)
	// apply the unique tag if provided
	if uniqueTag != "" {
		m.TagEntity(e, uniqueTag)
	}
	// set entity active and notify entity is active
	m.setActiveState(e, true)
	// return Entity
	return e, nil
}
