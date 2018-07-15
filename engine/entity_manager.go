package engine

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

// Provides services related to entities
type EntityManager struct {
	// list of entities currently spawned (whether active or not)
	Entities [MAX_ENTITIES]*EntityToken
	// Component data for entities
	Components *ComponentsTable
	// EntityTable stores: a list of allocated EntityTokens and a
	// list of available IDs from previous deallocations
	entityTable *EntityTable
	// updated entity lists of entities with given tags
	entitiesWithTag map[string]*UpdatedEntityList
	// entities which have been tagged uniquely
	uniqueEntities map[string]*EntityToken
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
		entityTable:     NewEntityTable(w.IDGen),
		entitiesWithTag: make(map[string]*UpdatedEntityList),
		uniqueEntities:  make(map[string]*EntityToken),
		eventBus:        w.Ev,
		spawnSubscription: w.Ev.Subscribe(
			"EntityManager::SpawnRequest",
			NewSimpleEventQuery(SPAWNREQUEST_EVENT)),
		despawnSubscription: w.Ev.Subscribe(
			"EntityManager::DespawnRequest",
			NewSimpleEventQuery(DESPAWNREQUEST_EVENT)),
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
func (m *EntityManager) Activate(entity *EntityToken) {
	m.setActiveState(entity, true)
}

// set an entity inactive and notify all active entity lists
func (m *EntityManager) Deactivate(entity *EntityToken) {
	m.setActiveState(entity, false)
}

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) setActiveState(entity *EntityToken, state bool) {
	// only act if the state is different to that which exists
	if entity.Active != state {
		if state {
			m.entityTable.active++
		} else {
			m.entityTable.active--
		}
		// start / stop logic accordingly
		m.Components.Logic[entity.ID].Active = state
		// set active state
		entity.Active = state
		// notify any listening lists
		m.activeEntityLists.notifyActiveState(entity, state)
	}
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedEntityList(
	q EntityQuery) *UpdatedEntityList {
	return m.activeEntityLists.GetUpdatedEntityList(q)
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetSortedUpdatedEntityList(
	q EntityQuery) *UpdatedEntityList {
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
func (m *EntityManager) UniqueTaggedEntity(tag string) (*EntityToken, error) {
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
		m.entitiesWithTag[tag] = m.GetUpdatedEntityList(EntityQueryFromTag(tag))
	}
}

func (m *EntityManager) EntityHasComponent(
	entity *EntityToken, COMPONENT int) bool {

	b, _ := entity.ComponentBitArray.GetBit(uint64(COMPONENT))
	return b
}

func (m *EntityManager) EntityHasTag(
	entity *EntityToken, tag string) bool {

	return m.EntityHasComponent(entity, TAGLIST_COMPONENT) &&
		m.Components.TagList[entity.ID].Has(tag)
}

// apply the given tags to the given entity
func (m *EntityManager) TagEntity(entity *EntityToken, tags ...string) {
	if !m.EntityHasComponent(entity, TAGLIST_COMPONENT) {
		entity.ComponentBitArray.SetBit(TAGLIST_COMPONENT)
	}
	for _, tag := range tags {
		m.Components.TagList[entity.ID].Add(tag)
		if entity.Active {
			m.createEntitiesWithTagListIfNeeded(tag)
		}
	}
	if entity.Active {
		m.activeEntityLists.checkActiveEntity(entity)
	}
}

// Tag each of the entities in the provided list
func (m *EntityManager) TagEntities(entities []*EntityToken, tag string) {
	for _, entity := range entities {
		m.TagEntity(entity, tag)
	}
}

// Remove a tag from an entity
func (m *EntityManager) UntagEntity(entity *EntityToken, tag string) {
	list := m.Components.TagList[entity.ID]
	list.Remove(tag)
	m.activeEntityLists.checkActiveEntity(entity)
}

// Remove a tag from each of the entities in the provided list
func (m *EntityManager) UntagEntities(entities []*EntityToken, tag string) {
	for _, entity := range entities {
		m.UntagEntity(entity, tag)
	}
}

// Get the number of allocated entities (not number of active, mind you)
func (m *EntityManager) NumEntities() (total int, active int) {
	return m.entityTable.n, m.entityTable.active
}

// Returns the Entities field, copied. Notice that this is of size MAX_ENTITIES
// and can have many nil elements, so the caller must checka and discard
// nil elements as they iterate
func (m *EntityManager) GetCurrentEntities() []*EntityToken {
	entities := make([]*EntityToken, 0, m.entityTable.n)
	for _, e := range m.Entities {
		if e != nil {
			entities = append(entities, e)
		}
	}
	return entities
}

func (m *EntityManager) String() string {
	return fmt.Sprintf("EntityManager[ %d / %d active ]\n",
		m.entityTable.n, m.entityTable.active)
}

// Somewhat expensive conversion of entire entity list to string, locking
// spawn/despawn from occurring while we read the entities (best to use for
// debugging, very ocassional diagnostic output)
func (m *EntityManager) Dump() string {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for _, entity := range m.Entities {
		if entity == nil {
			continue
		}
		tags := m.Components.TagList[entity.ID]
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			entity.ID, tags)
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
