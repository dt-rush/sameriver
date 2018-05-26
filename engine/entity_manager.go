/**
  * Manages the spawning and querying of entities
  *
**/

package engine

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-collections/go-datastructures/bitarray"
)

// created by game scene as a singleton, containing the component, entity,
// and tag data
type EntityManager struct {
	// Entity table stores component bitarrays, a list of allocated IDs,
	// and a list of available IDs from previous deallocations
	entityTable EntityTable
	// Tag table stores data for entity tagging system
	tagTable TagTable
	// Component data
	Components ComponentsTable

	// Channel for spawn entity requests
	spawnChannel chan EntitySpawnRequest

	// a flag signifying that Update() is running, set atomically
	// used so that we can start work needing to lock the enitity table
	// as soon as Update finishes (Update runs in sync and can't be blocked,
	// this give us the best chance of being able to do some work while
	// holding the lock before Update wants to run again, being interrupted
	// by our holding onto the lock (even though this only takes a short
	// time, if you read the comments for GetUpdatedActiveEntityList)
	updateRunning uint32
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
	m.Components.LinkEntityLocks(&m.entityTable.locks)
	// allocate space for the spawn buffer
	m.spawnChannel = make(chan EntitySpawnRequest,
		MAX_ENTITIES)
	// allocate tag system data members
	m.tagTable.entitiesWithTag = make(map[string]([]uint16))
	m.tagTable.tagsOfEntity = make(map[uint16]([]string))
}

// called once per scene Update() for scenes holding an entity manager
func (m *EntityManager) Update() {

	atomic.StoreUint32(&m.updateRunning, 1)
	// proces any requests to spawn new entities queued in the
	// buffered channel
	var t0 time.Time
	m.processSpawnChannel()
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		fmt.Printf("spawn: %d ms\n", time.Since(t0).Nanoseconds()/1e6)
	}
	// Finally, set the updateRunning flag, so that
	// GetActiveUpdatedActiveEntityList won't conflict with Update() when
	// it wants to lock the entity table for reading
	atomic.StoreUint32(&m.updateRunning, 0)
}

// get the ID for a new entity. Only called by SpawnEntity, which locks
// the entityTable, so it's safe that this method operates on that data.
// Returns int32 so that we can return -1 in case we have run out of space
// to spawn entities
func (m *EntityManager) allocateID() int32 {
	// if maximum entity count reached, fail with message
	if m.entityTable.numEntities == MAX_ENTITIES {
		Logger.Printf("Reached max entity count: %d. Will not allocate ID.\n",
			MAX_ENTITIES)
		return -1
	}
	// Increment the entity count
	m.entityTable.numEntities++
	// if there is a deallocated entity somewhere in the table before the
	// highest ID, return that ID to the caller
	n_avail := len(m.entityTable.availableIDs)
	var id uint16
	if n_avail > 0 {
		// there is an ID available for a previously deallocated entity.
		// pop it from the list and continue with that as the ID
		id = m.entityTable.availableIDs[n_avail-1]
		m.entityTable.availableIDs = m.entityTable.availableIDs[:n_avail-1]
	} else {
		// every slot in the table before the highest ID is filled
		id = uint16(m.entityTable.numEntities - 1)
	}
	// add the ID to the list of allocated IDs
	m.entityTable.allocatedIDs = append(m.entityTable.allocatedIDs, id)
	return int32(id)
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
			m.spawn(r)
		}
	}
}

// used by goroutines to request the spawning of an entity
func (m *EntityManager) RequestSpawn(r EntitySpawnRequest) {
	m.spawnChannel <- r
}

// given a list of components, spawn an entity with the default values
// returns the ID
func (m *EntityManager) spawn(r EntitySpawnRequest) uint16 {

	// lock the entityTable
	m.entityTable.mutex.Lock()
	defer m.entityTable.mutex.Unlock()

	// get an ID for the entity
	allocateIDResponse := m.allocateID()
	if allocateIDResponse == -1 {
		Logger.Printf("Ran out of entity space. Will not spawn entity with "+
			"tags: %v\n", r.Tags)
	}
	id := uint16(allocateIDResponse)
	// print a debug message
	entityManagerDebug("[Entity manager] Spawning: %d\n", id)
	// set the bitarray for this entity
	m.entityTable.componentBitArrays[id] = r.Components.ToBitArray()
	// copy the data into the component storage for each component
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
	m.Components.Active.Data[id] = true
	m.Components.ApplyComponentSet(id, r.Components)
	// apply the tags
	for _, tag := range r.Tags {
		m.TagEntity(id, tag)
	}
	// start the logic goroutine if supplied
	if r.Components.Logic != nil {
		entityLogicDebug("Starting logic for %d...", id)
		go r.Components.Logic.f(
			m.entityTable.getEntityToken(int32(id)),
			r.Components.Logic.StopChannel,
			m)
	}
	// notify entity is active
	go m.notifyActiveState(id, true)
	// return ID
	return id
}

// User facing function which is used to drain the state of the
// entity manager, and will also kill any pending spawn requests
func (m *EntityManager) DespawnAll() {
	despawnDebug("setting despawningAll flag")
	atomic.StoreUint32(&m.despawningAll, 1)
	// iterate all IDs which could have been allocated and despawn them
	// (each time a despawn goes through, the entityTable.allocatedIDs list
	// will shrink)
	for len(m.entityTable.allocatedIDs) > 0 {
		// for each allocated ID, build a token based on the current gen
		ID := m.entityTable.allocatedIDs[0]
		// lockEntity here will always return true because
		// the gen will never mismatch, since only a despawnInternal() call
		// could change that, and that occurs only in two places: here, where
		// we iterate one despawn for each entity while the entityTable is
		// locked, and in a call to Despawn() (in entity_modifications.go)
		// from a user which would only be able to proceed after we
		// released the entityTable lock, and would then exit since gen
		// had changed
		e := m.entityTable.getEntityToken(int32(ID))
		m.lockEntity(e)
		m.despawnInternal(e)
		m.releaseEntity(e)
	}
	// drain the spawn channel
	for len(m.spawnChannel) > 0 {
		// we're draining the channel, so do nothing
		_ = <-m.spawnChannel
	}
	atomic.StoreUint32(&m.despawningAll, 0)
}

// internal despawn function which assumes the EntityTable is locked
func (m *EntityManager) despawnInternal(e EntityToken) {

	id := uint16(e.ID)

	// lock the Mutex on the EntityTable
	m.entityTable.mutex.Lock()
	// if the gen doesn't match, another despawn for this same entity
	// has already been through here (if a regular Despawn() call and
	// a DespawnAll() were racing)
	if !m.entityTable.genValidate(e) {
		return
	}
	// decrement the entity count
	m.entityTable.numEntities--
	// add the ID to the list of available IDs
	m.entityTable.availableIDs = append(m.entityTable.availableIDs, id)
	// remove the ID from the list of allocated IDs
	removeUint16FromSlice(&m.entityTable.allocatedIDs, id)
	// Increment the gen for the ID
	// NOTE: it's important that we increment gen before resetting the
	// locks, since any goroutines waiting for the
	// lock to be 0 so they can claim it in AtomicEntityModify() will then
	// immediately want to check if the gen of the entity still matches.
	m.entityTable.incrementGen(id)
	// release the mutex on the EntityTable
	m.entityTable.mutex.Unlock()

	// Deactivate the entity to ensure that all updated entity lists are
	// notified
	despawnDebug("about to setActiveState...")
	m.setActiveState(id, false)
	despawnDebug("finished setActiveState()")
	// clear the entity from lists of tagged entities it's in
	t0 := time.Now()
	tags_to_clear := m.tagTable.tagsOfEntity[id]
	for _, tag_to_clear := range tags_to_clear {
		m.UntagEntity(id, tag_to_clear)
	}
	// remove the taglist for this entity
	m.tagTable.mutex.Lock()
	delete(m.tagTable.tagsOfEntity, id)
	m.tagTable.mutex.Unlock()
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		fmt.Printf("removing tags took: %d ms\n",
			time.Since(t0).Nanoseconds()/1e6)
	}
	// stop the entity's logic
	// NOTE: we don't need to worry about reading the component value
	// directly since this is called exclusively from AtomicEntityModify
	go func() {
		m.Components.Logic.Data[id].StopChannel <- true
	}()
}

// sets an entity active and notifies all watchers
func (m *EntityManager) activate(id uint16) {
	m.entityTable.mutex.Lock()
	defer m.entityTable.mutex.Unlock()
	entityManagerDebug("Activating: %d\n", id)
	m.setActiveState(id, true)
}

// sets an entity inactive and notifies all watchers
func (m *EntityManager) deactivate(id uint16) {
	m.entityTable.mutex.Lock()
	defer m.entityTable.mutex.Unlock()
	entityManagerDebug("Deactivating: %d\n", id)
	m.setActiveState(id, false)
}

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) setActiveState(id uint16, state bool) {
	// NOTE: we can access the active value directly since this is called
	// exclusively when the entityLock is set (will be reset at the end of
	// the loop iteration in processStateModificationChannel which called
	// this function via one of activate, deactivate, or despawn)
	if m.Components.Active.Data[id] != state {
		// setActiveState is only called when the entity is locked, so we're
		// good to write directly to the component
		m.Components.Active.Data[id] = state
		m.notifyActiveState(id, state)
	}
}

// Send a signal to all registered watchers that an entity has a certain
// active state, either true or false
func (m *EntityManager) notifyActiveState(id uint16, active bool) {

	var time = time.Now().UnixNano()
	updatedEntityListDebug("[%d] in notifyActiveState for %d: %v",
		time, id, active)

	m.activeEntityWatchersMutex.Lock()
	defer m.activeEntityWatchersMutex.Unlock()
	for _, watcher := range m.activeEntityWatchers {
		updatedEntityListDebug("[%d] testing Query %s...",
			time, watcher.Name)
		if watcher.Query.Test(id, m) {
			updatedEntityListDebug("[%d] Query %s matched %d",
				time, watcher.Name, id)
			// warn if the channel is full (we will block here if so)
			// NOTE: this can be very bad indeed, since now whatever
			// called Activate is blocking
			if len(watcher.Channel) == ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY {
				entityManagerDebug("⚠  active entity "+
					" watcher channel %s is full, causing block in "+
					" for NotifyActiveState(%d, %v)\n",
					watcher.Name, id, active)
			}
			// send the ID signal, or -(ID + 1), if active == false
			idSignal := int32(id)
			if !active {
				idSignal = -(idSignal + 1)
			}
			watcher.Channel <- m.entityTable.getEntityToken(idSignal)
		}
	}
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedActiveEntityList(
	q EntityQuery, name string) *UpdatedEntityList {

	// The basic idea of how we're going to build this list is as follows:
	//
	// Wait until Update() has recently finished.
	// Lock the entityTable, but only long enough to grab a snapshot of
	// 		allocatedIDs
	// Register an active query watcher with the query of the list we want to
	// 		build.
	// Create a temporary EntityToken channel, tempChannel, which we will
	// 		attach to the UpdatedEntityList we're building.
	// Enter a loop while the snapshot of allocatedIDs still has IDs to process
	//		In the loop we select from the active query watcher's channel, and
	//			if we get a remove signal, we try to remove that ID from
	//				whatever remains of the snapshot we're processing. If the
	//				ID isn't in the snapshot, we must have already inserted it
	//				so we forward the remove signal to the list
	//			if we get an add signal, we send the signal to the tempChannel
	//				attached to the list
	//		if we didn't select an insert/remove signal from the active query
	//			watcher, we pop an element from the snapshot of ID's and test
	//			the query against that entity. If it matches, we send an insert
	//			event to the channel.
	// When we've processed all the snapshot ID's, we stop the list (which
	// stops its update loop, listening on its channel), connect the proper
	// channel, that of the active query watcher, and then start the list.
	// the list has now checked every entity against its query and is current
	// with all events. We return it.

	// sleep until Update() has just finished
	for atomic.LoadUint32(&m.updateRunning) != 0 {

		time.Sleep(FRAME_SLEEP / 6)
	}
	// register a query watcher for the query given
	queryWatcher := m.GetActiveEntityQueryWatcher(q, name)
	// make a channel to temporarily act as the input channel to the list
	tempChannel := make(chan EntityToken)
	// build the list with the tempChannel to which we'll send entities
	list := NewUpdatedEntityList(
		tempChannel,
		queryWatcher.ID,
		name)
	// lock the entity table as quick as possible and grab a snapshot of
	// the allocated IDs
	m.entityTable.mutex.RLock()
	allocatedIDsSnapshot := make([]uint16, len(m.entityTable.allocatedIDs))
	copy(allocatedIDsSnapshot, m.entityTable.allocatedIDs)
	m.entityTable.mutex.RUnlock()
	updatedEntityListDebug("ID snapshot in trying to build "+
		"UpdatedActiveEntityList %s: %v", name, allocatedIDsSnapshot)
	// our aim is to check each snapshotID while still keeping up with
	// activate/deactivate signals
	for len(allocatedIDsSnapshot) > 0 {
		select {
		// prioritize processing new events
		case entitySignal := <-queryWatcher.Channel:
			updatedEntityListDebug("got signal on queryWatcher Channel "+
				"while trying to build UpdatedActiveEntityList %s", name)
			if entitySignal.ID < 0 {
				idToRemove := uint16(-(entitySignal.ID + 1))
				updatedEntityListDebug("removing ID %d from snapshot list "+
					"in trying to build UpdatedActiveEntitylist %s",
					idToRemove, name)
				// remove from snapshot if yet to process. remove
				// from list otherwise
				indexOfIdToRemoveInSnapshot := indexOfUint16InSlice(
					&allocatedIDsSnapshot, idToRemove)
				if indexOfIdToRemoveInSnapshot != -1 {
					removeIndexFromUint16Slice(
						&allocatedIDsSnapshot, indexOfIdToRemoveInSnapshot)
				} else {
					// if the ID has already been removed from the snapshot,
					// we added it to the list (since getting a remove signal
					// on the query means it matches the query), so send the
					// remove signal to the list
					list.InputChannel <- entitySignal
				}
			} else {
				// signal was an activate event. send to list
				updatedEntityListDebug("sending signal %d in "+
					"GetUpdatedActiveEntityList for %s",
					entitySignal.ID, name)
				list.InputChannel <- entitySignal
			}
		default:
			// pop an allocatedID and test it
			last_ix := len(allocatedIDsSnapshot) - 1
			id := allocatedIDsSnapshot[last_ix]
			updatedEntityListDebug("popped id %d from allocated IDs snapshot "+
				"while trying to build UpdatedActiveEntityList %s", id, name)
			allocatedIDsSnapshot = allocatedIDsSnapshot[:last_ix]
			if q.Test(id, m) {
				updatedEntityListDebug("sending signal %d in "+
					"GetUpdatedActiveEntityList for %s", id, name)
				list.InputChannel <- m.entityTable.getEntityToken(int32(id))
			}
		}
	}
	// we've finished catching up the snapshot ID's with the current event
	// stream. Set the channel properly on the list and return it
	updatedEntityListDebug("finished reviewing existing entities in "+
		"building of list %s. About to Stop(), set channel, start()",
		list.Name)
	list.Stop()
	list.InputChannel = queryWatcher.Channel
	list.start()
	updatedEntityListDebug("finished building list %s", list.Name)
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
	q EntityQuery, name string) EntityQueryWatcher {

	// create the query watcher
	qw := NewEntityQueryWatcher(q, name, IDGEN())
	// add it to the list of activeEntity watchers
	m.activeEntityWatchersMutex.Lock()
	m.activeEntityWatchers = append(m.activeEntityWatchers, qw)
	m.activeEntityWatchersMutex.Unlock()
	// return to the caller
	return qw
}

func (m *EntityManager) DeleteActiveEntityQueryWatcher(ID uint16) {
	m.activeEntityWatchersMutex.Lock()
	defer m.activeEntityWatchersMutex.Unlock()
	// remove the EntityQueryWatcher from the list of active watchers
	removeEntityQueryWatcherFromSliceByID(&m.activeEntityWatchers, ID)
}

// hold a single entity for modification, invoking a function which will be
// clear to access the entity components directly, releasing the entity on
// return. Return value is whether the entity was locked and f was run
// (remember lockEntity will wait as long as it needs to access the entity,
// but will fail if, upon locking,
func (m *EntityManager) AtomicEntityModify(
	entity EntityToken,
	f func(EntityToken)) bool {

	if !m.lockEntity(entity) {
		return false
	}
	f(entity)
	m.releaseEntity(entity)
	return true
}

// hold several entities for modification, invoking a function which will be
// clear to access the entity components directly, releasing the entities on
// return. Return value is whether the entities were locked and f was run
func (m *EntityManager) AtomicEntitiesModify(
	entities []EntityToken,
	f func([]EntityToken)) bool {

	var time = time.Now().UnixNano()

	atomicEntityModifyDebug("[%d] trying to lock %v", time, entities)
	if !m.lockEntities(entities) {
		return false
	}
	atomicEntityModifyDebug("[%d] lock succeeded, trying to run f()", time)
	f(entities)
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
	return true
}

// used by physics system to release an entity locked for modification
func (m *EntityManager) releaseEntity(entity EntityToken) {
	atomic.StoreUint32(&m.entityTable.locks[entity.ID], 0)
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

// apply the given tag to the given entity
func (m *EntityManager) TagEntity(id uint16, tag string) {
	m.tagTable.mutex.Lock()
	defer m.tagTable.mutex.Unlock()

	entityManagerDebug("Tagging %d with: %s\n", id, tag)
	_, t_of_e_exists := m.tagTable.tagsOfEntity[id]
	_, e_with_t_exists := m.tagTable.entitiesWithTag[tag]
	if !t_of_e_exists {
		m.tagTable.tagsOfEntity[id] = make([]string, 0)
	}
	if !e_with_t_exists {
		m.tagTable.entitiesWithTag[tag] = make([]uint16, 0)
	}
	m.tagTable.tagsOfEntity[id] = append(m.tagTable.tagsOfEntity[id], tag)
	m.tagTable.entitiesWithTag[tag] = append(m.tagTable.entitiesWithTag[tag], id)
}

// remove a tag from an entity
func (m *EntityManager) UntagEntity(id uint16, tag string) {
	m.tagTable.mutex.Lock()
	defer m.tagTable.mutex.Unlock()

	// NOTE: I'm aware the below code looks like some gross PHP stuff and
	// might be hard to read. Basically we do the following:
	//
	// last_ix = len(L) - 1
	// when i == index of element to remove,
	// L[i] = L[last_ix]
	// L = L[:last_ix]
	// remove the id from the list of entities with the tag
	last_ix := len(m.tagTable.entitiesWithTag[tag]) - 1
	for i, idInList := range m.tagTable.entitiesWithTag[tag] {
		if idInList == id {
			m.tagTable.entitiesWithTag[tag][i] = m.tagTable.entitiesWithTag[tag][last_ix]
			m.tagTable.entitiesWithTag[tag] = m.tagTable.entitiesWithTag[tag][:last_ix]
			break
		}
	}
	// remove the tag from the list of tags for the entity
	last_ix = len(m.tagTable.tagsOfEntity[id]) - 1
	for i := 0; i <= last_ix; i++ {
		if m.tagTable.tagsOfEntity[id][i] == tag {
			m.tagTable.tagsOfEntity[id][i] = m.tagTable.tagsOfEntity[id][last_ix]
			m.tagTable.tagsOfEntity[id] = m.tagTable.tagsOfEntity[id][:last_ix]
			break
		}
	}
}

// Tag each of the entities in the provided array of ID's with the given tag
func (m *EntityManager) TagEntities(ids []uint16, tag string) {
	m.tagTable.mutex.Lock()
	defer m.tagTable.mutex.Unlock()

	for _, id := range ids {
		m.TagEntity(id, tag)
	}
}

// Boolean check of whether a given entity has a given tag
func (m *EntityManager) EntityHasTag(id uint16, tag string) bool {
	m.tagTable.mutex.RLock()
	defer m.tagTable.mutex.RUnlock()

	for _, entity_tag := range m.tagTable.tagsOfEntity[id] {
		if entity_tag == tag {
			return true
		}
	}
	return false
}

// Gets the first entity with the given tag. Warns to console if the entity is
// not unique, and returns -1 if there is no such tagged entity. The return
// type is int32 to make sure that we can properly handle the largest int16
// id without conflicting with the uint16 representation of -1
func (m *EntityManager) GetUniqueTaggedEntity(tag string) int32 {
	m.tagTable.mutex.RLock()
	defer m.tagTable.mutex.RUnlock()

	entities := m.tagTable.entitiesWithTag[tag]
	if len(entities) == 0 {
		return -1
	} else {
		if len(entities) > 1 {
			Logger.Printf("⚠ more than one entity tagged with %s, but "+
				"GetUniqueTaggedEntity was called. This is a logic error. "+
				"Returning entitiesWithTag[0].",
				tag)
		}
		return int32(entities[0])
	}
}

// Boolean check of whether a given entity has a given component
func (m *EntityManager) EntityHasComponent(id uint16, COMPONENT int) bool {
	b, _ := m.entityTable.componentBitArrays[id].GetBit(uint64(COMPONENT))
	return b
}

// Returns the component bit array for an entity
func (m *EntityManager) entityComponentBitArray(id uint16) bitarray.BitArray {
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
	for _, id := range m.entityTable.allocatedIDs {
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			id, m.tagTable.tagsOfEntity[id])
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
