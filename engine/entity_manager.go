/**
  * Manages the spawning and querying of entities
  *
**/

package engine

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-collections/go-datastructures/bitarray"
)

// created by game scene as a singleton, containing the component, entity,
// and tag data
type EntityManager struct {
	// EntityTable stores component bitarrays, a list of allocated IDs,
	// and a list of available IDs from previous deallocations
	entityTable EntityTable
	// TagTable stores data for entity tagging system
	tagTable TagTable
	// EntityClassTable stores references to entity classes, which can be
	// retrieved by string ("crow", "turtle", "bear") in GetEntityClass()
	entityClassTable entityClassTable
	// Component data
	Components ComponentsTable

	// Channel for spawn entity requests which we don't need to get the
	// entity returned from (processed as a batch each Update())
	spawnChannel chan EntitySpawnRequest

	// a fag signifying that we are currently despawning all entities
	despawningAll uint32
	// updateMutex is used so that String() can grab the whole entity table
	// (massively interrupting Update()) and stringify it, safely at any rate
	// (this is not used often, or ever, unless you call String())
	updateMutex sync.Mutex

	// used to allow systems to keep an updated list of entities which have
	// components they're interested in operating on (eg. physics watches
	// for entities with position, velocity, and hitbox)
	activeEntityWatchers []EntityQueryWatcher
	// to protect modifying the above slice
	activeEntityWatchersMutex sync.RWMutex
}

func (m *EntityManager) Init() {
	// allocate component data
	m.Components = AllocateComponentsMemoryBlock()
	m.Components.LinkEntityManager(m)
	// allocate space for the spawn buffer
	m.spawnChannel = make(chan EntitySpawnRequest,
		MAX_ENTITIES)
	// init entity class table
	m.entityClassTable.Init()
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

// get the ID for a new entity. Only called by SpawnEntity, which locks
// the entityTable, so it's safe that this method operates on that data.
// Returns int32 so that we can return -1 in case we have run out of space
// to spawn entities
func (m *EntityManager) allocateID() (EntityToken, error) {
	m.entityTable.IDMutex.Lock()
	defer m.entityTable.IDMutex.Unlock()
	// if maximum entity count reached, fail with message
	if m.entityTable.numEntities == MAX_ENTITIES {
		msg := fmt.Sprintf("Reached max entity count: %d. "+
			"Will not allocate ID.\n", MAX_ENTITIES)
		Logger.Println(msg)
		return ENTITY_TOKEN_NIL, errors.New(msg)
	}
	// Increment the entity count
	m.entityTable.numEntities++
	// if there is a deallocated entity somewhere in the table before the
	// highest ID, return that ID to the caller
	n_avail := len(m.entityTable.availableIDs)
	var id int
	if n_avail > 0 {
		// there is an ID available for a previously deallocated entity.
		// pop it from the list and continue with that as the ID
		id = m.entityTable.availableIDs[n_avail-1]
		m.entityTable.availableIDs = m.entityTable.availableIDs[:n_avail-1]
	} else {
		// every slot in the table before the highest ID is filled
		id = m.entityTable.numEntities - 1
	}
	// add the ID to the list of allocated IDs
	entity := EntityToken{id, m.entityTable.gens[id]}
	m.entityTable.currentEntities = append(m.entityTable.currentEntities, entity)
	return entity, nil
}

// lock the ID table after waiting on spawn mutex to be unlocked,
// and grab a copy of the currently allocated IDs
func (m *EntityManager) snapshotAllocatedEntities() []EntityToken {
	m.entityTable.IDMutex.RLock()
	updatedEntityListDebug("got IDMutex in snapshot")
	defer m.entityTable.IDMutex.RUnlock()

	snapshot := make([]EntityToken, len(m.entityTable.currentEntities))
	copy(snapshot, m.entityTable.currentEntities)
	return snapshot
}

// process the spawn requests in the channel buffer
func (m *EntityManager) processSpawnChannel() {
	// get the current number of requests in the channel and only process
	// them. More may continue to pile up. They'll get processed next time.
	if atomic.LoadUint32(&m.despawningAll) != 1 {
		n := len(m.spawnChannel)
		for i := 0; i < n; i++ {
			// get the request from the channel
			r := <-m.spawnChannel
			spawnDebug("processing request: %v", r)
			m.Spawn(r)
		}
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
		m.tagTable.EntitiesWithTag(r.UniqueTag).Length() != 0 {
		return fail(fmt.Sprintf("requested to spawn unique entity for %s, "+
			"but %s already exists", r.UniqueTag))
	}

	// get an ID for the entity
	spawnDebug("trying to allocate ID for spawn request with tags %v", r.Tags)
	entity, err := m.allocateID()
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
	go m.notifyActiveState(entity, true)
	// return EntityToken
	return entity, nil
}

// User facing function which is used to drain the state of the
// entity manager, and will also kill any pending spawn requests
func (m *EntityManager) DespawnAll() {
	despawnDebug("setting despawningAll flag")
	atomic.StoreUint32(&m.despawningAll, 1)
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
	atomic.StoreUint32(&m.despawningAll, 0)
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
		m.notifyActiveState(entity, state)
	}
}

// Send a signal to all registered watchers that an entity has a certain
// active state, either true or false
func (m *EntityManager) notifyActiveState(entity EntityToken, active bool) {

	var time = time.Now().UnixNano()
	updatedEntityListDebug("notifyActiveState[%d] (%d, %v)",
		time, entity.ID, active)

	m.activeEntityWatchersMutex.Lock()
	defer m.activeEntityWatchersMutex.Unlock()
	defer updatedEntityListDebug("notifyActiveState[%d] unlocked "+
		"activeEntityWatchersMutex", entity.ID)
	for _, watcher := range m.activeEntityWatchers {
		updatedEntityListDebug("notifyActiveState[%d] testing Query %s...",
			time, watcher.Name)
		if watcher.Query.Test(entity, m) {
			updatedEntityListDebug("notifyActiveState[%d] Query %s matched %v",
				time, watcher.Name, entity)
			// warn if the channel is full (we will block here if so)
			// NOTE: this can be very bad indeed, since now whatever
			// called Activate is blocking
			if len(watcher.Channel) == ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY {
				entityManagerDebug("⚠  active entity "+
					" watcher channel %s is full, causing block in "+
					" for NotifyActiveState(%d, %v)\n",
					watcher.Name, entity.ID, active)
			}
			// send the ID signal, or -(ID + 1), if active == false
			signal := entity
			if !active {
				signal.ID = -(signal.ID + 1)
				updatedEntityListDebug("notifyActiveState[%d] sending "+
					"remove:%v to %s", time, entity, watcher.Name)
			} else {
				updatedEntityListDebug("notifyActiveState[%d] sending "+
					"insert:%v to %s", time, entity, watcher.Name)
			}
			watcher.Channel <- signal
		}
	}
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedActiveEntityList(
	name string, q EntityQuery) *UpdatedEntityList {

	defer functionEndDebug("GetUpdatedActiveEntityList() returning %s list "+
		"(as yet unbuilt)",
		name)

	// register a query watcher for the query given
	queryWatcher := m.GetActiveEntityQueryWatcher(name, q)
	// build the list (as yet empty), provide it with a backlog, and start it
	backlog := m.snapshotAllocatedEntities()
	backlogTester := func(entity EntityToken) bool {
		return q.Test(entity, m)
	}
	list := NewUpdatedEntityList(
		name,
		queryWatcher.ID,
		queryWatcher.Channel,
		backlog,
		backlogTester)
	list.start()
	return list
}

// Stops the channel-watching update-loop goroutine for the entity list
// and deletes the active watcher created for it
func (m *EntityManager) DeleteUpdatedActiveEntityList(l UpdatedEntityList) {
	m.DeleteActiveEntityQueryWatcher(l.ID)
	l.stopUpdateLoopChannel <- true
}

// Return a channel which will receive the id of an entity whenever an entity
/// becomes active with a component set matching the query bitarray, and which
// will receive -(id + 1) whenever an entity is *despawned* with a component
// set matching the query bitarray
func (m *EntityManager) GetActiveEntityQueryWatcher(
	name string, q EntityQuery) EntityQueryWatcher {

	entityManagerDebug("about to make active-EntityQueryWatcher for %s", name)
	defer functionEndDebug("made active-EntityQueryWatcher for %s", name)

	// create the query watcher
	qw := NewEntityQueryWatcher(q, name, IDGEN())
	// add it to the list of activeEntity watchers
	go func() {
		m.activeEntityWatchersMutex.Lock()
		m.activeEntityWatchers = append(m.activeEntityWatchers, qw)
		m.activeEntityWatchersMutex.Unlock()
	}()
	// return to the caller
	return qw
}

func (m *EntityManager) DeleteActiveEntityQueryWatcher(ID int) {
	m.activeEntityWatchersMutex.Lock()
	defer m.activeEntityWatchersMutex.Unlock()
	// remove the EntityQueryWatcher from the list of active watchers
	removeEntityQueryWatcherFromSliceByID(&m.activeEntityWatchers, ID)
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

// wait until we can lock an entity and check gen-match after lock. If
// gen mismatch, release and return false. else return true
func (m *EntityManager) lockEntity(entity EntityToken) bool {
	// wait until we can lock the entity
	for !atomic.CompareAndSwapUint32(&m.entityTable.locks[entity.ID], 0, 1) {
		// if we can't lock the entity, sleep half a frame
		time.Sleep(FRAME_SLEEP / 2)
	}
	// return value is whether we locked the entity with the same gen
	// (in between when the caller acquired the EntityToken and now, the
	// entity may have despawned)
	if !m.entityTable.genValidate(entity) {
		m.releaseEntity(entity)
		return false
	}
	entityManagerDebug("LOCKED entity %v", entity)
	return true
}

// used by physics system to release an entity locked for modification
func (m *EntityManager) releaseEntity(entity EntityToken) {
	atomic.StoreUint32(&m.entityTable.locks[entity.ID], 0)
	entityManagerDebug("RELEASED entity %v", entity)
}

// lock multiple entities (with return value true only if gen matches for
// all entities locked)
func (m *EntityManager) lockEntities(entities []EntityToken) bool {
	// attempt to lock all entities, keeping track of which ones we have
	var allValid = true
	var locked = make([]EntityToken, 0)
	var time = time.Now().UnixNano()
	entityLocksDebug("[%d] attempting to lock %d entities: %v",
		time, len(entities), entities)
	for _, entity := range entities {
		if !m.lockEntity(entity) {
			entityLocksDebug("[%d] locking failed for %d", time, entity.ID)
			allValid = false
			break
		} else {
			entityLocksDebug("[%d] locking succeeded for %d", time, entity.ID)
			locked = append(locked, entity)
		}
	}
	// if one was invalid, the locking of this group no longer makes sense
	// (one was despawned since the []EntityToken was formulated by the caller)
	// so, release all those entities we've already locked and return false
	if !allValid {
		m.releaseEntities(locked)
		return false
	}
	// else return true
	return true
}

// release multiple entities
func (m *EntityManager) releaseEntities(entities []EntityToken) {
	for _, entity := range entities {
		m.releaseEntity(entity)
	}
}

// self explanatory
func (m *EntityManager) releaseTwoEntities(
	entityA EntityToken, entityB EntityToken) {
	atomic.StoreUint32(&m.entityTable.locks[entityA.ID], 0)
	atomic.StoreUint32(&m.entityTable.locks[entityB.ID], 0)
}

// used by physics system to attempt to lock an entity for modification, but
// will not sleep and retry if the lock fails, simply returns false if we
// didn't lock it, or if it had a new gen. Releases the entity if we locked it
// and gen doesn't match.
func (m *EntityManager) attemptLockEntityOnce(entity EntityToken) bool {
	// do a single attempt to lock
	locked := atomic.CompareAndSwapUint32(
		&m.entityTable.locks[entity.ID], 0, 1)
	// if we locked the entity but gen mistmatches, release it
	// and return false
	if locked &&
		!m.entityTable.genValidate(entity) {
		m.releaseEntity(entity)
		return false
	}
	// else, either we locked it and gen matched (pass), or we didn't lock
	// (fail), so return `locked`
	return locked
}

// used by collision system to attempt to lock two entities for modification
// (if both can't be acquired, we just back off and try again another cycle;
// collision between those two entities won't occur this cycle (there are many
// per second, so it's not noticeable to the user)
func (m *EntityManager) attemptLockTwoEntitiesOnce(
	entityA EntityToken, entityB EntityToken) bool {

	// attempt to lock entity A
	if !m.attemptLockEntityOnce(entityA) {
		// NOTE: we don't need to release the entity since if it failed
		// due to gen mismatch, attemptLockEntityOnce will itself release it,
		// and if it failed due to not locking, there's nothing to release
		return false
	}
	// attempt to lock entity B
	if !m.attemptLockEntityOnce(entityB) {
		// NOTE: if we're here, we *did* acquire A
		m.releaseEntity(entityA)
		return false
	}
	// if we're here, we locked both
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
func (m *EntityManager) RegisterEntityClass(class EntityClass) {
	// add the class to the EntityClassTable and return
	m.entityClassTable.addEntityClass(class)
}

// Get an entity class by name
func (m *EntityManager) EntityClass(name string) EntityClass {
	return m.entityClassTable.getClass(name)
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

	m.tagTable.mutex.Lock()
	defer m.tagTable.mutex.Unlock()

	m.createTagListIfNeeded(tag)
	return m.tagTable.entitiesWithTag[tag]
}

func (m *EntityManager) createTagListIfNeeded(tag string) {
	if _, exists := m.tagTable.entitiesWithTag[tag]; !exists {
		m.tagTable.entitiesWithTag[tag] = m.GetUpdatedActiveEntityList(
			tag, GenericEntityQueryFromTag(tag))
	}
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
