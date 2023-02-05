package sameriver

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

// Provides services related to entities
type EntityManager struct {
	// the world this EntityManager is inside
	w *World
	// list of entities currently spawned (whether active or not)
	entities [MAX_ENTITIES]*Entity
	// Component data for entities
	components *ComponentTable
	// EntityTable stores: a list of allocated Entitys and a
	// list of available IDs from previous deallocations
	entityTable *EntityTable
	// updated entity lists created by the user according to provided filters
	lists map[string]*UpdatedEntityList
	// updated entity lists of entities with given tags
	entitiesWithTag map[string]*UpdatedEntityList
	// entities which have been tagged uniquely
	uniqueEntities map[string]*Entity
	// Channel for spawn entity requests (processed as a batch each Update())
	spawnSubscription *EventChannel
	// Channel for despawn entity requests (processed as a batch each Update())
	despawnSubscription *EventChannel
}

// Construct a new entity manager
func NewEntityManager(w *World) *EntityManager {
	em := &EntityManager{
		w:               w,
		components:      NewComponentTable(),
		entityTable:     NewEntityTable(w.IdGen),
		lists:           make(map[string]*UpdatedEntityList),
		entitiesWithTag: make(map[string]*UpdatedEntityList),
		uniqueEntities:  make(map[string]*Entity),
		spawnSubscription: w.Events.Subscribe("EntityManager::SpawnRequest",
			SimpleEventFilter("spawn-request")),
		despawnSubscription: w.Events.Subscribe("EntityManager::DespawnRequest",
			SimpleEventFilter("despawn-request")),
	}
	return em
}

// called once per scene Update() for scenes holding an entity manager
func (m *EntityManager) Update(allowance_ms float64) float64 {
	// TODO: base spawning off allowance. Spawn enough and do no more.
	t0 := time.Now()
	m.processDespawnChannel()
	m.processSpawnChannel()
	return float64(time.Since(t0).Nanoseconds()) / 1.0e6
}

// set an entity Active and notify all active entity lists
func (m *EntityManager) Activate(e *Entity) {
	m.setActiveState(e, true)
}

// set an entity inactive and notify all active entity lists
func (m *EntityManager) Deactivate(e *Entity) {
	m.setActiveState(e, false)
}

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) setActiveState(e *Entity, state bool) {
	// only act if the state is different to that which exists
	if e.Active != state {
		if state {
			m.entityTable.active++
		} else {
			m.entityTable.active--
		}
		// start / stop all logics of this entity accordingly
		for _, l := range e.Logics {
			l.active = state
		}
		// set active state
		e.Active = state
		// notify any listening lists
		m.notifyActiveState(e, state)
	}
}

// Gets the first entity with the given tag. Warns to console if the entity is
// not unique. Returns an error if the entity doesn't exist
func (m *EntityManager) UniqueTaggedEntity(tag string) (*Entity, error) {
	if e, ok := m.uniqueEntities[tag]; ok {
		return e, nil
	} else {
		errorMsg := fmt.Sprintf("tried to fetch unique entity %s, but did "+
			"not exist", tag)
		return nil, errors.New(errorMsg)
	}
}

func (m *EntityManager) UpdatedEntitiesWithTag(tag string) *UpdatedEntityList {
	m.createEntitiesWithTagListIfNeeded(tag)
	return m.entitiesWithTag[tag]
}

func (m *EntityManager) createEntitiesWithTagListIfNeeded(tag string) {
	if _, exists := m.entitiesWithTag[tag]; !exists {
		m.entitiesWithTag[tag] =
			m.GetUpdatedEntityList(m.w.entityFilterFromTag(tag))
	}
}

func (m *EntityManager) EntityHasComponent(e *Entity, name string) bool {
	b, _ := e.ComponentBitArray.GetBit(uint64(m.components.ixs[name]))
	return b
}

func (m *EntityManager) EntityHasTag(e *Entity, tag string) bool {
	return m.EntityHasComponent(e, "GenericTags") &&
		e.GetTagList("GenericTags").Has(tag)
}

// apply the given tags to the given entity
func (m *EntityManager) TagEntity(e *Entity, tags ...string) {
	if !m.EntityHasComponent(e, "GenericTags") {
		e.ComponentBitArray.SetBit(uint64(m.components.ixs["GenericTags"]))
		*e.GetTagList("GenericTags") = NewTagList()
	}
	for _, tag := range tags {
		e.GetTagList("GenericTags").Add(tag)
		if e.Active {
			m.createEntitiesWithTagListIfNeeded(tag)
		}
	}
	if e.Active {
		m.checkActiveEntity(e)
	}
}

// Tag each of the entities in the provided list
func (m *EntityManager) TagEntities(entities []*Entity, tag string) {
	for _, e := range entities {
		m.TagEntity(e, tag)
	}
}

// Remove a tag from an entity
func (m *EntityManager) UntagEntity(e *Entity, tag string) {
	e.GetTagList("GenericTags").Remove(tag)
	m.checkActiveEntity(e)
}

// Remove a tag from each of the entities in the provided list
func (m *EntityManager) UntagEntities(entities []*Entity, tag string) {
	for _, e := range entities {
		m.UntagEntity(e, tag)
	}
}

// Get the number of allocated entities (not number of active, mind you)
func (m *EntityManager) NumEntities() (total int, active int) {
	return len(m.entityTable.currentEntities), m.entityTable.active
}

// return a map where the keys are the current entities (aka an idiomatic
// go "set")
func (m *EntityManager) GetCurrentEntitiesSet() map[*Entity]bool {
	return m.entityTable.currentEntities
}

// return a copy of the current entities map (allowing you to spawn/despawn
// while iterating over it)
func (m *EntityManager) GetCurrentEntitiesSetCopy() map[*Entity]bool {
	setCopy := make(map[*Entity]bool, len(m.entityTable.currentEntities))
	for e, _ := range m.entityTable.currentEntities {
		setCopy[e] = true
	}
	return setCopy
}

func (m *EntityManager) String() string {
	return fmt.Sprintf("EntityManager[ %d / %d active ]\n",
		len(m.entityTable.currentEntities), m.entityTable.active)
}

// dump entities with tags
func (m *EntityManager) DumpEntities() string {
	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for _, e := range m.entities {
		if e == nil {
			continue
		}
		tags := e.GetTagList("GenericTags")
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			e.ID, tags)
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
