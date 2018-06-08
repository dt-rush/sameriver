package engine

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"

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
	// entityLogicTable contains references to the LogicUnits of the entities, if
	// they supplied one
	entityLogicTable EntityLogicTable
	// EntityClassTable stores references to entity classes, which can be
	// retrieved by string ("crow", "turtle", "bear") in GetEntityClass()
	entityClasses EntityClassTable
	// ActiveEntityListCollection is used by GetUpdatedEntityList to
	// store EntityQueryWatchers and references to UpdatedEntityLists used
	// to implement GetUpdatedEntityList
	activeEntityLists ActiveEntityListCollection
	// Channel for spawn entity requests which we don't need to get the
	// entity returned from (processed as a batch each Update())
	spawnChannel chan EntitySpawnRequest
	// spawnMutex prevents despawn / spawn events from occurring at the same time
	spawnMutex sync.Mutex
}

func (m *EntityManager) Init() {
	// allocate space for the spawn buffer
	m.spawnChannel = make(chan EntitySpawnRequest,
		MAX_ENTITIES)
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
	spawnDebug("applying component set for %v...", entity)
	m.ApplyComponentSetAtomic(r.Components)(entity)
	// allocate ComponentValueLock's for each component for this entity
	for i := 0; i < N_COMPONENTS; i++ {
		m.Components[i].locks[entity.ID] = NewComponentValueLock()
	}
	// apply the tags
	spawnDebug("applying tags for %v...", entity)
	for _, tag := range r.Tags {
		m.TagEntityAtomic(tag)(entity)
	}
	// apply the unique tag if provided
	if r.UniqueTag != "" {
		m.TagEntityAtomic(r.UniqueTag)(entity)
	}
	// start the logic goroutine if supplied
	if r.Logic != nil {
		spawnDebug("Setting and starting logic for %d...", entity.ID)
		logicUnit := r.entityLogicTable.setLogic(entity, r.Logic)
	}
	// set entity active and notify entity is active
	m.setActiveState(entity, true)
	// add the entity to the list of current entities
	spawnDebug("adding %v to current entities...", entity)
	m.entityTable.addToCurrentEntities(entity)
	// set despawnFlag to 0 for the entity
	m.entityTable.despawnFlags[entity.ID].Store(0)
	// create the Array-Based Read-Write Queuing Lock for activeModificationLocks
	// (allocating on heap on spawn this way is good since it means the
	// locks will be less likely to false-share with other locks, causing
	// interconnect traffic when the queue moves, since we expect that in any
	// given game situation, the pattern of locks is a uniform distribution,
	// hence there's little likelyhood that two locks "next to" each other
	// will both be updated within the same cache-line lifetime)
	m.entityTable.activeModificationLocks[entity.ID] = NewActiveModificationLock()
	// return EntityToken
	return entity, nil
}

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) setActiveState(entity EntityToken, state bool) {
	// only act if the state is different to that which exists
	if m.entityTable.activeState[entity.ID] != state {
		// start / stop logic accordingly
		logic := m.logicTable.getLogic(entity)
		if state == true {
			go logicUnit.f(entity, logicUnit.stopChannel, m)
		} else {
			go func() {
				logicUnit.stopChannel <- true
			}()
		}
		// set active state
		m.entitytable.activeState[entity.ID] = state
		// notify any listening lists
		go m.activeEntityLists.notifyActiveState(entity, state)
	}
}

// returns the entity's active state if gen matches, else an error
// (using an option-type pattern)
func (m *EntityManager) GetActiveState(entity EntityToken) (bool, error) {
	if !m.entityTable.genValidate(entity) {
		return false, errors.New("tried to get active state of despawned entity")
	}
	return m.entityTable.activeStates[entity.ID], nil
}

// User facing function which is used to drain the state of the
// entity manager, and will also kill any pending spawn requests
func (m *EntityManager) DespawnAll() {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()
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

	// Deactivate and notify
	despawnDebug("setting %v inactive...", entity)
	m.setActiveState(entity, false)
	// delete the entity's logic (we have to do this *after* stopping it)
	despawnDebug("deleting %v logic...", entity)
	m.entityLogicTable.deleteLogic(entity)
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedEntityList(
	q EntityQuery) *UpdatedEntityList {

	return m.activeEntityLists.GetUpdatedEntityList(q)
}

// hold a single entity for modification, invoking a function which will be
// clear to access the entity components directly, releasing the entity on
// return. Return value is whether the entity was locked and f was run
// (remember lockEntity will wait as long as it needs to access the entity,
// but will fail if, upon locking, gen does not match)
func (m *EntityManager) AtomicEntityModify(
	req EntityModificationRequest,
	f func()) bool {

	// keep track of locks acquired
	locksAcquired := make([]EntityComponent, 0)
	defer m.releaseEntityComponents(locksAcquired)
	// acquire locks on components in sorted order
	sort.Slice(req.components, func(i int, j int) bool {
		return req.components[i] < req.components[j]
	})
	// lock the activemodification lock as a "reader" for the duration of this
	// function (that is, other copies of AtomicEntityModify can also lock)
	m.entityTable.activeModificationLocks[req.entity.ID].RLock()
	defer m.entityTable.activeModificationLocks[req.entity.ID].RUnlock()
	// attempt to lock each component in order
	for _, component := range req.components {
		// if lock acquired, add it to the list of locks acquired
		if m.lockEntityComponent(req.entity, component) {
			locksAcquired = append(locksAcquired,
				EntityComponent{entity, component})
		} else {
			// if here, the lock on the entity's component was failed, because
			// the gen changed (was despawned in between when the caller got
			// their entity token and when they called this function ), so we
			// should return false to notify caller (this triggers the deferred
			// release of all locks on entity components)
			return false
		}
	}

	// if we're here, all components were acquired for the entity
	f()
	return true
}

// hold several entities on several components for modification, invoking a
// function which will be clear to access the entity components directly,
// releasing them on return. Return value is whether the entities were locked
// successfully (and hence, f ran)
//
// this is an extension of AtomicEntityModify and the comments for that function
// should be seen for reference in understanding this code
func (m *EntityManager) AtomicEntitiesModify(
	reqs []EntityModificationRequest,
	f func()) bool {

	// acquire entities in sorted order of ID's
	sort.Slice(reqs, func(i int, j int) bool {
		return reqs[i].entity.ID < reqs[j].entity.ID
	})

	locksAcquired := make([]EntityComponent, 0)
	defer m.releaseEntityComponents(locksAcquired)
	for _, req := range reqs {
		sort.Slice(req.components, func(i int, j int) bool {
			return req.components[i] < req.components[j]
		})
		m.entityTable.activeModificationLocks[req.entity.ID].RLock()
		defer m.entityTable.activeModificationLocks[req.entity.ID].RUnlock()
		for _, component := range req.components {
			if m.lockEntityComponent(req.entity, component) {
				locksAcquired = append(locksAcquired,
					EntityComponent{entity, component})
			} else {
				return false
			}
		}
	}

	f()
	return true
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
	m.tags.mutex.RLock()
	defer m.tags.mutex.RUnlock()

	list := m.EntitiesWithTag(tag)
	if list.Length() == 0 {
		errorMsg := fmt.Sprintf("tried to fetch unique entity %s, but did "+
			"not exist", tag)
		tagsDebug(errorMsg)
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

// Somewhat expensive conversion of entire entity list to string, locking
// spawn/despawn from occurring while we read the entities (best to use for
// debugging, very ocassional diagnostic output)
func (m *EntityManager) String() string {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for _, entity := range m.entityTable.currentEntities {
		tags, err := m.Components.TagList.SafeGet(entity)
		if err != nil {
			continue // entity was despawned
		}
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			entity.ID, tags)
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
