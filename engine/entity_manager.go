package engine

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

// Provides services related to entities
type EntityManager struct {
	// the world this EntityManager is inside
	w *World
	// list of entities currently spawned (whether active or not)
	Entities [MAX_ENTITIES]*Entity
	// Component data for entities
	Components *ComponentsTable
	// EntityTable stores: a list of allocated Entitys and a
	// list of available IDs from previous deallocations
	entityTable *EntityTable
	// updated entity lists of entities with given tags
	entitiesWithTag map[string]*UpdatedEntityList
	// entities which have been tagged uniquely
	uniqueEntities map[string]*Entity
	// ActiveEntityListCollection is used by GetUpdatedEntityList to
	// store references to existing UpdatedEntityLists (by name)
	activeEntityLists *ActiveEntityListCollection
	// used to communicate with other systems
	eventBus *EventBus
	// Channel for spawn entity requests (processed as a batch each Update())
	spawnSubscription *EventChannel
	// Channel for despawn entity requests (processed as a batch each Update())
	despawnSubscription *EventChannel
	// spawnMutex prevents despawn / spawn events from occurring while we
	// convert the entire EntityManager to string (expensive!)
	spawnMutex sync.Mutex
}

// Construct a new entity manager
func NewEntityManager(w *World) *EntityManager {
	em := &EntityManager{
		w:               w,
		entityTable:     NewEntityTable(w.IDGen),
		entitiesWithTag: make(map[string]*UpdatedEntityList),
		uniqueEntities:  make(map[string]*Entity),
		eventBus:        w.Ev,
		spawnSubscription: w.Ev.Subscribe(
			"EntityManager::SpawnRequest",
			SimpleEventFilter(SPAWNREQUEST_EVENT)),
		despawnSubscription: w.Ev.Subscribe(
			"EntityManager::DespawnRequest",
			SimpleEventFilter(DESPAWNREQUEST_EVENT)),
	}
	em.Components = NewComponentsTable(em)
	em.activeEntityLists = NewActiveEntityListCollection(em)
	return em
}

// called once per scene Update() for scenes holding an entity manager
func (m *EntityManager) Update() {
	m.processDespawnChannel()
	m.processSpawnChannel()
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
		// start / stop logic accordingly
		m.Components.Logic[e.ID].Active = state
		// set active state
		e.Active = state
		// notify any listening lists
		m.activeEntityLists.notifyActiveState(e, state)
	}
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedEntityList(q EntityFilter) *UpdatedEntityList {
	return m.activeEntityLists.GetUpdatedEntityList(q)
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetSortedUpdatedEntityList(
	q EntityFilter) *UpdatedEntityList {
	return m.activeEntityLists.GetSortedUpdatedEntityList(q)
}

// get a previously-created UpdatedEntityList by name, or nil if does not exist
func (m *EntityManager) GetUpdatedEntityListByName(
	name string) *UpdatedEntityList {

	if list, ok := m.activeEntityLists.lists[name]; ok {
		return list
	} else {
		return nil
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

func (m *EntityManager) EntitiesWithTag(tag string) *UpdatedEntityList {
	m.createEntitiesWithTagListIfNeeded(tag)
	return m.entitiesWithTag[tag]
}

func (m *EntityManager) createEntitiesWithTagListIfNeeded(tag string) {
	if _, exists := m.entitiesWithTag[tag]; !exists {
		m.entitiesWithTag[tag] =
			m.GetUpdatedEntityList(m.w.entityFilterFromTag(tag))
	}
}

func (m *EntityManager) EntityHasComponent(
	e *Entity, COMPONENT int) bool {

	b, _ := e.ComponentBitArray.GetBit(uint64(COMPONENT))
	return b
}

func (m *EntityManager) EntityHasTag(
	e *Entity, tag string) bool {

	return m.EntityHasComponent(e, TAGLIST_COMPONENT) &&
		m.Components.TagList[e.ID].Has(tag)
}

// apply the given tags to the given entity
func (m *EntityManager) TagEntity(e *Entity, tags ...string) {
	if !m.EntityHasComponent(e, TAGLIST_COMPONENT) {
		e.ComponentBitArray.SetBit(TAGLIST_COMPONENT)
	}
	for _, tag := range tags {
		m.Components.TagList[e.ID].Add(tag)
		if e.Active {
			m.createEntitiesWithTagListIfNeeded(tag)
		}
	}
	if e.Active {
		m.activeEntityLists.checkActiveEntity(e)
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
	list := m.Components.TagList[e.ID]
	list.Remove(tag)
	m.activeEntityLists.checkActiveEntity(e)
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

// Returns the Entities field, copied. Notice that this is of size MAX_ENTITIES
// and can have many nil elements, so the caller must checka and discard
// nil elements as they iterate
func (m *EntityManager) GetCurrentEntities() []*Entity {
	entities := make([]*Entity, 0, len(m.entityTable.currentEntities))
	for _, e := range m.Entities {
		if e != nil {
			entities = append(entities, e)
		}
	}
	return entities
}

func (m *EntityManager) String() string {
	return fmt.Sprintf("EntityManager[ %d / %d active ]\n",
		len(m.entityTable.currentEntities), m.entityTable.active)
}

// Somewhat expensive conversion of entire entity list to string, locking
// spawn/despawn from occurring while we read the entities (best to use for
// debugging, very ocassional diagnostic output)
func (m *EntityManager) Dump() string {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for _, e := range m.Entities {
		if e == nil {
			continue
		}
		tags := m.Components.TagList[e.ID]
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			e.ID, tags)
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
