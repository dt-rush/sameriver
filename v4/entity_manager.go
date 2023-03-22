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
	// Component data for entities
	components *ComponentTable
	// EntityIDALlocator stores: a list of allocated Entitys and a
	// list of available IDs from previous deallocations
	entityIDAllocator *EntityIDAllocator
	// updated entity lists created by the user according to provided filters
	lists map[string]*UpdatedEntityList
	// updated entity lists of entities with given tags
	entitiesWithTag map[string]*UpdatedEntityList
	// entities which have been tagged uniquely
	uniqueEntities map[string]*Entity
	// entities that are active
	activeEntities map[*Entity]bool
	// Channel for spawn entity requests (processed as a batch each Update())
	spawnSubscription *EventChannel
	// Channel for despawn entity requests (processed as a batch each Update())
	despawnSubscription *EventChannel
}

// Construct a new entity manager
func NewEntityManager(w *World) *EntityManager {
	em := &EntityManager{
		w:                   w,
		components:          NewComponentTable(MAX_ENTITIES),
		entityIDAllocator:   NewEntityIDAllocator(MAX_ENTITIES, w.IdGen),
		lists:               make(map[string]*UpdatedEntityList),
		entitiesWithTag:     make(map[string]*UpdatedEntityList),
		uniqueEntities:      make(map[string]*Entity),
		activeEntities:      make(map[*Entity]bool),
		spawnSubscription:   w.Events.Subscribe(SimpleEventFilter("spawn-request")),
		despawnSubscription: w.Events.Subscribe(SimpleEventFilter("despawn-request")),
	}
	return em
}

func (m *EntityManager) Components() *ComponentTable {
	return m.components
}

// called once per scene Update() for scenes holding an entity manager
func (m *EntityManager) Update(allowance_ms float64) float64 {
	t0 := time.Now()
	// TODO: base spawning off allowance. Spawn enough and do no more.
	m.processSpawnChannel()
	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1e6
	return allowance_ms - dt_ms
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
			m.entityIDAllocator.active++
			m.activeEntities[e] = true
		} else {
			m.entityIDAllocator.active--
			delete(m.activeEntities, e)
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
			m.GetUpdatedEntityList(EntityFilterFromTag(tag))
	}
}

func (m *EntityManager) EntityHasComponent(e *Entity, name ComponentID) bool {
	b, _ := e.ComponentBitArray.GetBit(uint64(m.components.ixs[name]))
	return b
}

func (m *EntityManager) EntityHasTag(e *Entity, tag string) bool {
	return e.GetTagList(GENERICTAGS).Has(tag)
}

// apply the given tags to the given entity
func (m *EntityManager) TagEntity(e *Entity, tags ...string) {
	for _, tag := range tags {
		e.GetTagList(GENERICTAGS).Add(tag)
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
	e.GetTagList(GENERICTAGS).Remove(tag)
	m.checkActiveEntity(e)
}

// Remove a tag from each of the entities in the provided list
func (m *EntityManager) UntagEntities(entities []*Entity, tag string) {
	for _, e := range entities {
		m.UntagEntity(e, tag)
	}
}

// get the maximum number of entities without a resizing and reallocating of
// components and system data (if Expand() is not a nil function for that system)
func (m *EntityManager) MaxEntities() int {
	return m.entityIDAllocator.capacity
}

// Get the number of allocated entities (not number of active, mind you)
func (m *EntityManager) NumEntities() (total int, active int) {
	return len(m.entityIDAllocator.currentEntities), m.entityIDAllocator.active
}

// returns a map of all active entities
func (m *EntityManager) GetActiveEntitiesSet() map[*Entity]bool {
	return m.activeEntities
}

// return a map where the keys are the current entities (aka an idiomatic
// go "set")
func (m *EntityManager) GetCurrentEntitiesSet() map[*Entity]bool {
	return m.entityIDAllocator.currentEntities
}

// return a copy of the current entities map (allowing you to spawn/despawn
// while iterating over it)
func (m *EntityManager) GetCurrentEntitiesSetCopy() map[*Entity]bool {
	setCopy := make(map[*Entity]bool, len(m.entityIDAllocator.currentEntities))
	for e := range m.entityIDAllocator.currentEntities {
		setCopy[e] = true
	}
	return setCopy
}

func (m *EntityManager) String() string {
	return fmt.Sprintf("EntityManager[ %d / %d active ]\n",
		len(m.entityIDAllocator.currentEntities), m.entityIDAllocator.active)
}

// dump entities with tags
func (m *EntityManager) DumpEntities() string {
	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for e := range m.entityIDAllocator.currentEntities {
		if e == nil {
			continue
		}
		tags := e.GetTagList(GENERICTAGS)
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			e.ID, tags)
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
