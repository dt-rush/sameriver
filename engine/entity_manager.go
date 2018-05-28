/**
  * Manages the spawning and querying of entities
  *
**/

package engine

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/atomic"
	"sync"
	"time"

	"github.com/golang-collections/go-datastructures/bitarray"
)

// created by game scene as a singleton, containing the component, entity,
// and tag data
type EntityManager struct {
    // Component data (can be accessed by users (but only safely inside an
    // AtomicEntit(y|ies)Modify callback)
	Components ComponentsTable
	// EntityTable stores component bitarrays, a list of allocated EntityTokens,
	// and a list of available IDs from previous deallocations
	entityTable EntityTable
	// TagTable stores data for the entity tagging system
	tagTable TagTable
	// EntityClassTable stores references to entity classes, which can be
	// retrieved by string ("crow", "turtle", "bear") in GetEntityClass()
	entityClasses entityClassTable
    // ActiveEntityListCollection is used by GetUpdatedActiveEntityList to
    // store EntityQueryWatchers and references to UpdatedEntityLists used
    // to implement GetUpdatedActiveEntityList
    activeEntityLists ActiveEntityListCollection
	// Channel for spawn entity requests which we don't need to get the
	// entity returned from (processed as a batch each Update())
	spawnChannel chan EntitySpawnRequest
	// updateMutex is used so that String() can grab the whole entity table
	// (massively interrupting Update()) and stringify it, safely at any rate
	// (this is not used often, or ever, unless you call String() - really this
    // should only be used in developent, in debugging, or in printing state
    // during a crash
	updateMutex sync.Mutex
}

func (m *EntityManager) Init() {
	// allocate component data
	m.Components = AllocateComponentsMemoryBlock()
	m.Components.LinkEntityManager(m)
	// allocate space for the spawn buffer
	m.spawnChannel = make(chan EntitySpawnRequest,
		MAX_ENTITIES)
	// init entity class table
	m.entityClasses.Init()
	// init tag table
	m.tagTable.Init()
}

// called once per scene Update() for scenes holding an entity manager
func (m *EntityManager) Update() {

	// proces any requests to spawn new entities queued in the
	// buffered channel
	var t0 time.Time
	m.processSpawnChannel()
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		fmt.Printf("spawn: %d ms\n", time.Since(t0).Nanoseconds()/1e6)
	}
}

// process the spawn requests in the channel buffer
func (m *EntityManager) processSpawnChannel() {
	// get the current number of requests in the channel and only process
	// them. More may continue to pile up. They'll get processed next time.
    m.spawnMutex.Lock()
    defer m.spawnMutex.Unlock()
    n := len(m.spawnChannel)
    for i := 0; i < n; i++ {
        // get the request from the channel
        r := <-m.spawnChannel
        spawnDebug("processing request: %v", r)
        m.Spawn(r)
    }
}

// used by goroutines to request the spawning of an entity
func (m *EntityManager) RequestSpawn(r EntitySpawnRequest) {
	m.spawnChannel <- r
}

// given a list of components, spawn an entity with the default values
// returns the EntityToken (used to spawn an entity for which we *want* the
// token back)
func (m *EntityManager) Spawn(r EntitySpawnRequest) (EntityToken, error) {
	defer functionEndDebug("spawn %v", r)

	// used if spawn is impossible for various reasons
	var fail = func(msg string) (EntityToken, error) {
		spawnDebug(msg)
		return ENTITY_TOKEN_NIL, errors.New(msg)
	}

	// if the spawn request has a unique tag, return error if tag already
	// has an entity
	if r.UniqueTag != "" &&
		m.EntitiesWithTag(r.UniqueTag).Length() != 0 {
		return fail(fmt.Sprintf("requested to spawn unique entity for %s, "+
			"but %s already exists", r.UniqueTag))
	}

	// get an ID for the entity
	spawnDebug("trying to allocate ID for spawn request with tags %v", r.Tags)
	entity, err := m.entityTable.allocateID()
	if err != nil {
		errorMsg := fmt.Sprintf("⚠ Error in allocateID(): %s. Will not spawn "+
			"entity with tags: %v\n", err, r.Tags)
		spawnDebug(errorMsg)
		return fail("ran out of entity space")
	}
	// lock the entity (this prevents GetUpdatedActiveEntityList's list-building
	// goroutine from querying the entity if its now-allocated entity is included
	// in the snapshot of allocated entities)
	m.lockEntity(entity)
	defer m.releaseEntity(entity)
	// print a debug message
	spawnDebug("Spawning: %v\n", entity)
	// set the bitarray for this entity
	m.entityTable.componentBitArrays[entity.ID] = r.Components.ToBitArray()
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
	m.Components.Active.Data[entity.ID] = true
	m.Components.ApplyComponentSet(entity.ID, r.Components)
	// apply the tags
	for _, tag := range r.Tags {
		m.TagEntityAtomic(tag)(entity)
	}
	// apply the unique tag if provided
	if r.UniqueTag != "" {
		m.TagEntityAtomic(r.UniqueTag)(entity)
	}
	// start the logic goroutine if supplied
	if r.Components.Logic != nil {
		entityLogicDebug("Starting logic for %d...", entity.ID)
		go r.Components.Logic.f(
			entity,
			r.Components.Logic.StopChannel,
			m)
	}
	// notify entity is active
	go m.activeEntityLists.notifyActiveState(entity, true)
	// return EntityToken
	return entity, nil
}

// User facing function which is used to drain the state of the
// entity manager, and will also kill any pending spawn requests
func (m *EntityManager) DespawnAll() {
	despawnDebug("setting despawningAll flag")
	m.despawningAll.Store(1)
	// iterate all IDs which could have been allocated and despawn them
	// (each time a despawn goes through, the entityTable.currentEntities list
	// will shrink)
	for len(m.entityTable.currentEntities) > 0 {
		// for each allocated ID, build a token based on the current gen
		// lockEntity here will always return true because
		// the gen will never mismatch, since only a despawnInternal() call
		// could change that, and that occurs only in two places: here, where
		// we iterate one despawn for each entity while the entityTable is
		// locked, and in a call to Despawn() (in entity_modifications.go)
		// from a user which would only be able to proceed after we
		// released the entityTable lock, and would then exit since gen
		// had changed
		entity := m.entityTable.currentEntities[0]
		m.lockEntity(entity)
		m.despawnInternal(entity)
		m.releaseEntity(entity)
	}
	// drain the spawn channel
	for len(m.spawnChannel) > 0 {
		// we're draining the channel, so do nothing
		_ = <-m.spawnChannel
	}
	m.despawningAll.Store(0)
}

// internal despawn function which assumes the EntityTable is locked
func (m *EntityManager) despawnInternal(entity EntityToken) {
	m.entityTable.IDMutex.Lock()
	defer m.entityTable.IDMutex.Unlock()

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

	// Deactivate the entity to ensure that all updated entity lists are
	// notified
	despawnDebug("about to setActiveState(%v, false)...", entity)
	m.setActiveState(entity, false)
	despawnDebug("finished setActiveState(%v, false)", entity)
	// remove each tag from this ID in the tag table
	// (also sends removal signals to the tag lists)
	t0 := time.Now()
	tags_to_clear := m.tagTable.tagsOfEntity[entity.ID]
	despawnDebug("about to remove tags for %v...", entity)
	for _, tag_to_clear := range tags_to_clear {
		m.UntagEntityAtomic(tag_to_clear)(entity)
	}
	despawnDebug("removed tags for %v...", entity)
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		fmt.Printf("removing tags took: %d ms\n",
			time.Since(t0).Nanoseconds()/1e6)
	}
	// stop the entity's logic
	// NOTE: we don't need to worry about reading the component value
	// directly since this is called exclusively from AtomicEntityModify
	go func() {
		m.Components.Logic.Data[entity.ID].StopChannel <- true
	}()
}

// sets an entity active and notifies all watchers
func (m *EntityManager) activate(entity EntityToken) {
	entityManagerDebug("Activating: %d\n", entity.ID)
	m.setActiveState(entity, true)
}

// sets an entity inactive and notifies all watchers
func (m *EntityManager) deactivate(entity EntityToken) {
	entityManagerDebug("Deactivating: %d\n", entity.ID)
	m.setActiveState(entity, false)
}

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) setActiveState(entity EntityToken, state bool) {
	// NOTE: we can access the active value directly since this is called
	// exclusively when the entityLock is set (will be reset at the end of
	// the loop iteration in processStateModificationChannel which called
	// this function via one of activate, deactivate, or despawn)
	if m.Components.Active.Data[entity.ID] != state {
		// setActiveState is only called when the entity is locked, so we're
		// good to write directly to the component
		m.Components.Active.Data[entity.ID] = state
		m.activeEntityLists.notifyActiveState(entity, state)
	}
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedEntityList(q EntityQuery) *UpdatedEntityList {
	return m.activeEntityLists.GetUpdatedEntityList(q)
}

// hold a single entity for modification, invoking a function which will be
// clear to access the entity components directly, releasing the entity on
// return. Return value is whether the entity was locked and f was run
// (remember lockEntity will wait as long as it needs to access the entity,
// but will fail if, upon locking, gen does not match)
func (m *EntityManager) AtomicEntityModify(
	entity EntityToken,
	f func()) bool {

	if !m.lockEntity(entity) {
		return false
	}
	f()
	m.releaseEntity(entity)
	return true
}

// hold several entities for modification, invoking a function which will be
// clear to access the entity components directly, releasing the entities on
// return. Return value is whether the entities were locked and f was run
func (m *EntityManager) AtomicEntitiesModify(
	entities []EntityToken,
	f func()) bool {

	var time = time.Now().UnixNano()

	atomicEntityModifyDebug("[%d] trying to lock %v", time, entities)
	if !m.lockEntities(entities) {
		return false
	}
	atomicEntityModifyDebug("[%d] lock succeeded, trying to run f()", time)
	f()
	atomicEntityModifyDebug("[%d] f() completed, relesing entities", time)
	m.releaseEntities(entities)
	return true
}

// Boolean check of whether a given entity has a given tag
func (m *EntityManager) EntityHasTag(entity EntityToken, tag string) bool {
	tagsDebug("in EntityHasTag(%d, %s), trying to acquire tagTable mutex",
		entity.ID, tag)
	m.tagTable.mutex.RLock()
	tagsDebug("in EntityHasTag(%d, %s) got tagTable mutex",
		entity.ID, tag)
	defer tagsDebug("EntityHasTag(%d, %s) released tagTable mutex",
		entity.ID, tag)
	defer m.tagTable.mutex.RUnlock()

	for _, entity_tag := range m.tagTable.tagsOfEntity[entity.ID] {
		if entity_tag == tag {
			return true
		}
	}
	return false
}

// Register an entity class (subsequently retrievable)
func (m *EntityManager) AddEntityClass(class EntityClass) {
	// add the class to the EntityClassTable and return
	m.entityClasses.addClass(class)
}

// Get an entity class by name
func (m *EntityManager) GetEntityClass(name string) EntityClass {
	return m.entityClasses.getClass(name)
}

// Gets the first entity with the given tag. Warns to console if the entity is
// not unique. Returns an error if the entity doesn't exist
func (m *EntityManager) UniqueTaggedEntity(tag string) (EntityToken, error) {
	m.tagTable.mutex.RLock()
	defer m.tagTable.mutex.RUnlock()

	list := m.EntitiesWithTag(tag)
	if list.Length() == 0 {
		tagsDebug("tried to fetch unique entity %s, but did not exist", tag)
		return ENTITY_TOKEN_NIL, errors.New("no such entity")
	}
	if list.Length() > 1 {
		tagsDebug("⚠ more than one entity tagged with %s, but "+
			"GetUniqueTaggedEntity was called. This is a logic error. "+
			"Returning the first entity.", tag)
	}
	return list.First()
}

func (m *EntityManager) EntitiesWithTag(
	tag string) *UpdatedEntityList {

	return m.tagTable.entitiesWithTag(tag)
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

// Somewhat expensive conversion of entire entity list to string
func (m *EntityManager) String() string {
	m.updateMutex.Lock()
	m.tagTable.mutex.RLock()
	defer m.tagTable.mutex.RUnlock()
	defer m.updateMutex.Unlock()

	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for _, entity := range m.entityTable.currentEntities {
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			entity.ID, m.tagTable.tagsOfEntity[entity.ID])
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
