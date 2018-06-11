package engine

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityManager struct {
	// Component data (can be accessed by users (but only safely inside an
	// AtomicEntit(y|ies)Modify callback)
	Components ComponentsTable
	// EntityTable stores component bitarrays, a list of allocated EntityTokens,
	// active states, and a list of available IDs from previous deallocations
	entityTable EntityTable
	// TagTable contains UpdatedEntityLists for tagged entities
	tags TagTable
	// Logics contains references to logic funcs of the entities, if
	// they supplied one
	Logics map[EntityToken]func()
	// entityClasses stores references to entity classes, which can be
	// retrieved by string ("crow", "turtle", "bear") in GetEntityClass()
	entityClasses map[string]EntityClass
	// ActiveEntityListCollection is used by GetUpdatedEntityList to
	// store EntityQueryWatchers and references to UpdatedEntityLists used
	// to implement GetUpdatedEntityList
	activeEntityLists ActiveEntityListCollection
	// used to communicate with other systems
	ev *EventBus
	// Channel for spawn entity requests (processed as a batch each Update())
	spawnSubscription EventChannel
	// Channel for despawn entity requests (processed as a batch each Update())
	despawnSubscription EventChannel
	// spawnMutex prevents despawn / spawn events from occurring while we
	// convert the entire EntityManager to string (expensive!)
	spawnMutex sync.Mutex
}

func (m *EntityManager) Init(ev *EventBus) {
	// take down a reference to the event bus
	m.ev = ev
	// set up spawn / despawn channels as listeners on the appropriate
	// events
	m.spawnSubscription = ev.Subscribe(
		"EntityManager::SpawnRequest",
		NewSimpleEventQuery(SPAWNREQUEST_EVENT))
	m.despawnSubscription = ev.Subscribe(
		"EntityManager::DespawnRequest",
		NewSimpleEventQuery(DESPAWNREQUEST_EVENT))
	// init ComponentsTable
	m.Components.Init(m)
	// init EntityClassTable
	m.entityClasses.Init()
	// init TagTable
	m.tags.Init(m)
	// init ActiveEntityListCollection
	m.activeEntityLists.Init(m)
	// init EntityLogicTable
	m.EntityLogicTable.Init()
}

// called once per scene Update() for scenes holding an entity manager
func (m *EntityManager) Update() {

	// proces any requests to spawn new entities queued in the
	// buffered channel
	m.processDespawnChannel()
	m.processSpawnChannel()
}

// set an entity Active and notify all active entity lists
func (m *EntityManager) Activate(entity EntityToken) {
	m.setActiveState(entity, true)
}

// set an entity inactive and notify all active entity lists
func (m *EntityManager) Deactivate(entity EntityToken) {
	m.setActiveState(entity, false)
}

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) setActiveState(entity EntityToken, state bool) {
	// only act if the state is different to that which exists
	if m.entityTable.activeStates[entity.ID] != state {
		// start / stop logic accordingly
		if state == true {
			m.EntityLogicTable.ActivateLogic(entity)
		} else {
			m.EntityLogicTable.DeactivateLogic(entity)
		}
		// set active state
		m.entityTable.activeStates[entity.ID] = state
		// notify any listening lists
		m.activeEntityLists.notifyActiveState(entity, state)
	}
}

// returns the entity's active state if gen matches, else an error
// (using an option-type pattern)
func (m *EntityManager) getActiveState(entity EntityToken) (bool, error) {
	if !m.entityTable.genValidate(entity) {
		return false, errors.New("tried to get active state of despawned entity")
	}
	return m.entityTable.activeStates[entity.ID], nil
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedEntityList(
	q EntityQuery) *UpdatedEntityList {

	return m.activeEntityLists.GetUpdatedEntityList(q)
}

// Register an entity class (subsequently retrievable)
func (m *EntityManager) AddEntityClass(class EntityClass) {
	m.classes[ec.Name()] = ec
}

// Get an entity class by name
func (m *EntityManager) GetEntityClass(name string) EntityClass {
	return m.classes[name]
}

// Gets the first entity with the given tag. Warns to console if the entity is
// not unique. Returns an error if the entity doesn't exist
func (m *EntityManager) UniqueTaggedEntity(tag string) (EntityToken, error) {
	list := m.EntitiesWithTag(tag)
	if list.Length() == 0 {
		errorMsg := fmt.Sprintf("tried to fetch unique entity %s, but did "+
			"not exist", tag)
		return ENTITY_TOKEN_NIL, errors.New(errorMsg)
	}
	if list.Length() > 1 {
		tagsDebug("⚠ more than one entity tagged with %s, but "+
			"GetUniqueTaggedEntity was called. This is a logic error. "+
			"Returning the first entity.", tag)
	}
	return list.FirstEntity()
}

func (m *EntityManager) EntitiesWithTag(
	tag string) *UpdatedEntityList {

	return m.tags.GetEntitiesWithTag(tag)
}

// Boolean check of whether a given entity has a given component
func (m *EntityManager) EntityHasComponent(id int, COMPONENT int) bool {
	b, _ := m.entityTable.componentBitArrays[id].GetBit(uint64(COMPONENT))
	return b
}

// Returns the component bit array for an entity
func (m *EntityManager) entityComponentBitArray(id int) bitarray.BitArray {
	return m.entityTable.componentBitArrays[id]
}

// apply the given tag to the given entity
func (m *EntityManager) TagEntity(tag string) func(EntityToken) {
	return func(entity EntityToken) {

		// add the tag to the taglist component
		m.Components.TagList[entity.ID].Add(tag)
		// if the entity is active, it has already been checked by all lists,
		// thus generate a new signal to add it to the list of the tag
		if m.entityTable.activeStates[entity.ID] {
			m.tags.createEntitiesWithTagListIfNeeded(tag)
			m.activeEntityLists.checkActiveEntity(entity)
		}
	}
}

// remove a tag from an entity
func (m *EntityManager) UntagEntity(tag string) func(EntityToken) {
	return func(entity EntityToken) {
		m.Components.TagList[entity.ID].Remove(tag)
		m.activeEntityLists.checkActiveEntity(entity)
	}
}

// Tag each of the entities in the provided array of ID's with the given tag
func (m *EntityManager) TagEntities(tag string) func([]EntityToken) {
	return func(entities []EntityToken) {

		for _, entity := range entities {
			m.TagEntity(tag)(entity)
		}
	}
}

// Somewhat expensive conversion of entire entity list to string, locking
// spawn/despawn from occurring while we read the entities (best to use for
// debugging, very ocassional diagnostic output)
func (m *EntityManager) String() string {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for _, entity := range m.entityTable.currentEntities {
		tags := m.Components.TagList[entity.ID]
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			entity.ID, tags)
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
