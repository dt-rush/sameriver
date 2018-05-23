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

// Used to represent an entity with an ID at a point in time. Despawning the
// entity at a given ID will increment gen (gen ("generation") data is stored
// in EntityTable). The token storing *gen* prevents goroutines from
// requesting modifications on a new entity in edge cases which may occur
// quite readily depending on timing of event processing
type EntityToken struct {
	id  uint16
	gen uint32
}

// Used by goroutines which have requested to modify an entity to communicate
// their desired modification, whether an entity state change or a change
// to component values
type EntityModificationType int

const (
	ENTITY_STATE_MODIFICATION     = EntityModificationType(iota)
	ENTITY_COMPONENT_MODIFICATION = EntityModificationType(iota)
)

type EntityModification struct {
	// the EntityToken is filled by the Manager, not visible to the
	// caller
	entity EntityToken
	// Type is used to allow the type-assertion of Modification to be either
	// an instance of EntityState or ComponentSet
	Type         EntityModificationType
	Modification interface{}
}

// Used for goroutines to request that an entity be activated, deactivated, or
// despawned
type EntityState int

const (
	ENTITY_ACTIVATE   = EntityState(iota)
	ENTITY_DEACTIVATE = EntityState(iota)
	ENTITY_DESPAWN    = EntityState(iota)
)

type entityStateModification struct {
	// the entity to apply the state change to
	entity EntityToken
	// the requested state
	State EntityState
}

// Used for goroutines to request modifications to entity components
type entityComponentModification struct {
	// the entity to apply the component change to
	entity EntityToken
	// the component values for the entity
	Components ComponentSet
}

// Used to spawn entities
type EntitySpawnRequest struct {
	Components ComponentSet
	Tags       []string
}

// used by the EntityManager to hold info about the allocated entities
type EntityTable struct {
	// guards all changes to this table as atomic
	mutex sync.RWMutex
	// how many entities there are
	numEntities int
	// list of IDs which have been allocated
	allocatedIDs []uint16
	// list of available entity ID's which have previously been deallocated
	availableIDs []uint16
	// bitarray used to keep track of which entities have which components
	// (indexes are IDs, bitarrays have bit set if entity has the
	// component corresponding to that index)
	componentBitArrays [MAX_ENTITIES]bitarray.BitArray
	// the gen of an ID is how many times an entity has been spawned on that ID
	gens [MAX_ENTITIES]uint32
	// locks so that goroutines can operate atomically on individual entities
	// (eg. imagine two squirrels coming upon a nut and trying to eat it. One
	// must win!). Also used by systems like PhysicsSystem to avoid modifying
	// those entities while they're held for modification (hence the importance
	// of not holding entities for modification longer than, say, one update
	// cycle (at 60fps, around 16 ms). In fact, one update cycle is a hell of
	// a long time. It should be less than a millisecond or two.
	heldForModificationLocks [MAX_ENTITIES]uint32
}

// used by the EntityManager to tag entities
type TagTable struct {
	// guards all state
	mutex sync.RWMutex
	// data members to support the entity tagging system, which allows us to
	// associate a set of strings with an entity
	// tag -> []IDs
	entitiesWithTag map[string]([]uint16)
	// ID -> []tag
	tagsOfEntity map[uint16]([]string)
}

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

	// Channel for logic goroutines to send requests to modify entity
	// active/spawn state (internal only, accessed implicitly via the
	// AtomicEntityModify() method)
	stateModificationChannel chan entityStateModification
	// Channel for logic goroutines to send requests to modify entity
	// components (internal only, accessed implicitly via the
	// AtomicEntityModify() method)
	componentModificationChannel chan entityComponentModification
	// set atomically when DespawnAll() is called, to short-circuit any
	// processing of the prior two channels. We simply despawn all entities
	// if this is 1
	despawnAllFlag uint32
	// Channel for spawn entity requests
	spawnChannel chan EntitySpawnRequest

	// locked during the entirety of Update()
	UpdateMutex sync.Mutex

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
	// allocate the state and component modification channels with the buffer
	// size defined in constants
	m.stateModificationChannel = make(chan entityStateModification,
		ENTITY_MODIFICATION_CHANNEL_CAPACITY)
	m.componentModificationChannel = make(chan entityComponentModification,
		ENTITY_MODIFICATION_CHANNEL_CAPACITY)
	m.spawnChannel = make(chan EntitySpawnRequest,
		MAX_ENTITIES)
	// allocate tag system data members
	m.tagTable.entitiesWithTag = make(map[string]([]uint16))
	m.tagTable.tagsOfEntity = make(map[uint16]([]string))
}

// called once per scene Update() for scenes holding an entity manager
func (m *EntityManager) Update() {
	m.UpdateMutex.Lock()
	defer m.UpdateMutex.Unlock()
	// First, act on the despawn all flag, despawning all entities if
	// it's set.
	var t0 time.Time
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		t0 = time.Now()
	}
	m.actOnDespawnAllFlag()
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		fmt.Printf("despawnall: %d ms\n", time.Since(t0).Nanoseconds()/1e6)
	}
	// Second, process the spawn/despawn, activate/deactivate requests queued in
	// the buffered channel
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		t0 = time.Now()
	}
	m.processStateModificationChannel()
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		fmt.Printf("statemod: %d ms\n", time.Since(t0).Nanoseconds()/1e6)
	}
	// Third, process the component modifications queued in the buffered channel
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		t0 = time.Now()
	}
	m.processComponentModificationChannel()
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		fmt.Printf("componentmod: %d ms\n", time.Since(t0).Nanoseconds()/1e6)
	}
	// Finally, process any requests to spawn new entities queued in the
	// buffered channel
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		t0 = time.Now()
	}
	m.processSpawnChannel()
	if DEBUG_ENTITY_MANAGER_UPDATE_TIMING {
		fmt.Printf("spawn: %d ms\n", time.Since(t0).Nanoseconds()/1e6)
	}

}

// react to the despawnall flag
func (m *EntityManager) actOnDespawnAllFlag() {
	// if the flag is 1, set to 0 and proceed to despawn all
	if atomic.CompareAndSwapUint32(&m.despawnAllFlag, 1, 0) {
		Logger.Println("[Entity manager] Despawning all...")
		entityManagerDebug("Currently spawned: %v\n",
			m.entityTable.allocatedIDs)
		// iterate all IDs which could have been allocated and despawn them
		for i := 0; i < m.entityTable.numEntities; i++ {
			// we can call this safely on each ID, even those unallocated,
			// since it's idempotent with respect to unspawned entities
			// (it will exit early if the entity is already despawned)
			m.despawn(uint16(i))
		}
		// drain the modification and spawn channels (NOTE: by this point,
		// all logic goroutines should have been terminated, so nothing new
		// should be coming to these channels)
		for len(m.stateModificationChannel) > 0 {
			// we're draining the channel, so do nothing
			_ = <-m.stateModificationChannel
		}
		for len(m.componentModificationChannel) > 0 {
			// we're draining the channel, so do nothing
			_ = <-m.componentModificationChannel
		}
		for len(m.spawnChannel) > 0 {
			// we're draining the channel, so do nothing
			_ = <-m.spawnChannel
		}
	}
}

// Process the entityStateModifications on stateModificationChannel
func (m *EntityManager) processStateModificationChannel() {
	// get the current number of requests in the channel and only process
	// them. More may continue to pile up. They'll get processed next time.
	n := len(m.stateModificationChannel)
	entityManagerDebug("%d state modifications to process...", n)
	for i := 0; i < n; i++ {
		// get the request from the channel
		r := <-m.stateModificationChannel
		// if gen doesn't match, this request is invalid
		// (TODO: We don't really need this check. If AtomicEntityModify is
		// being respected, we won't ever receive a modification for an entity
		// unless its gen matches, since a despawn event is itself a
		// modification which will change the gen, resulting in a failure to
		// acquire the entity in AtomicEntityModify)
		if r.entity.gen != m.entityTable.gens[r.entity.id] {
			continue
		}
		// process the event
		switch r.State {
		case ENTITY_ACTIVATE:
			entityManagerDebug("processing activate()")
			m.activate(r.entity.id)
		case ENTITY_DEACTIVATE:
			entityManagerDebug("processing deactivate()")
			m.deactivate(r.entity.id)
		case ENTITY_DESPAWN:
			entityManagerDebug("processing despawn()")
			m.despawn(r.entity.id)
		}
		// now that we've applied the modification, release the lock on the
		// ID
		atomic.StoreUint32(
			&m.entityTable.heldForModificationLocks[r.entity.id], 0)
	}
}

// process the requests to set entity components
func (m *EntityManager) processComponentModificationChannel() {
	// get the current number of requests in the channel and only process
	// them. More may continue to pile up. They'll get processed next time.
	n := len(m.componentModificationChannel)
	for i := 0; i < n; i++ {
		// get the request from the channel
		r := <-m.componentModificationChannel
		// if gen doesn't match, this request is invalid
		// (TODO: We don't really need this check. If AtomicEntityModify is
		// being respected, we won't ever receive a modification for an entity
		// unless its gen matches, since a despawn event is itself a
		// modification which will change the gen, resulting in a failure to
		// acquire the entity in AtomicEntityModify)
		if r.entity.gen != m.entityTable.gens[r.entity.id] {
			continue
		}
		// take the component data and set it, on the ID
		m.Components.ApplyComponentSet(r.entity.id, r.Components)
		// now that we've applied the modification, release the lock on the
		// ID
		atomic.StoreUint32(
			&m.entityTable.heldForModificationLocks[r.entity.id], 0)
	}
}

// process the spawn requests in the channel buffer
func (m *EntityManager) processSpawnChannel() {
	// get the current number of requests in the channel and only process
	// them. More may continue to pile up. They'll get processed next time.
	n := len(m.spawnChannel)
	for i := 0; i < n; i++ {
		// get the request from the channel
		r := <-m.spawnChannel
		m.spawn(r)
	}
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
	m.Components.Active.SafeSet(id, true)
	m.Components.ApplyComponentSet(id, r.Components)
	// apply the tags
	for _, tag := range r.Tags {
		m.TagEntity(id, tag)
	}
	// start the logic goroutine if supplied
	if r.Components.Logic != nil {
		logicDebug("Starting logic for %d...", id)
		go r.Components.Logic.LogicFunc(id,
			r.Components.Logic.StopChannel,
			m)
	}
	// notify entity is active
	go m.notifyActiveState(id, true)
	// return ID
	return id
}

func (m *EntityManager) despawn(id uint16) {
	m.entityTable.mutex.Lock()
	defer m.entityTable.mutex.Unlock()

	// decrement the entity count
	m.entityTable.numEntities--
	// add the ID to the list of available IDs
	m.entityTable.availableIDs = append(m.entityTable.availableIDs, id)
	// remove the ID from the list of allocated IDs
	removeUint16FromSlice(id, &m.entityTable.allocatedIDs)
	// Deactivate the entity to ensure that all updated entity lists are
	// notified
	m.setActiveState(id, false)
	// clear the entity from lists of tagged entities it's in
	tags_to_clear := m.tagTable.tagsOfEntity[id]
	for _, tag_to_clear := range tags_to_clear {
		m.UntagEntity(id, tag_to_clear)
	}
	// remove the taglist for this entity
	m.tagTable.mutex.Lock()
	delete(m.tagTable.tagsOfEntity, id)
	m.tagTable.mutex.Unlock()
	// stop the entity's logic
	m.Components.Logic.SafeGet(id).StopChannel <- true
	// Increment the gen for the ID
	// NOTE: it's important that we increment gen before resetting the
	// heldForModificationLocks, since any goroutines waiting for the
	// lock to be 0 so they can claim it in AtomicEntityModify() will then
	// immediately want to check if the gen of the entity still matches.
	atomic.AddUint32(&m.entityTable.gens[id], 1)
	// Clear the modificationLock for the entity (any goroutine either trying
	// to set it to 1 with an old gen or holding it for modification with
	// an old gen will fail
	atomic.StoreUint32(&m.entityTable.heldForModificationLocks[id], 0)
}

// setting the flag will cause the entities to all get despawned next time
// processEntityModificationRequests() is called
func (m *EntityManager) DespawnAll() {
	atomic.StoreUint32(&m.despawnAllFlag, 1)
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
	// TODO: this pattern seems weird. Active should not be a component,
	// but metadata. Refactor it
	// Only set (and notify) if not already in given state
	if m.Components.Active.SafeGet(id) != state {
		m.Components.Active.SafeSet(id, state)
		go m.notifyActiveState(id, state)
	}
}

// Send a signal to all registered watchers that an entity has a certain
// active state, either true or false
func (m *EntityManager) notifyActiveState(id uint16, active bool) {

	m.activeEntityWatchersMutex.Lock()
	defer m.activeEntityWatchersMutex.Unlock()
	for _, watcher := range m.activeEntityWatchers {
		if watcher.Query.Test(id, m) {
			// warn if the channel is full (we will block here if so)
			// NOTE: this can be very bad indeed, since now whatever
			// called Activate is blocking
			if len(watcher.Channel) == ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY {
				entityManagerDebug("⚠  active entity "+
					" watcher channel %s is full, causing block in goroutine "+
					" for NotifyActiveState(%d, %v)\n",
					watcher.Name, id, active)
			}
			// send the ID signal, or -(ID + 1), if active == false
			idSignal := int32(id)
			if !active {
				idSignal = -(idSignal + 1)
			}
			watcher.Channel <- idSignal
		}
	}
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedActiveEntityList(
	q EntityQuery, name string) *UpdatedEntityList {

	// get a new UpdatedEntityList
	queryWatcher := m.GetActiveEntityQueryWatcher(q, name)
	list := NewUpdatedEntityList(
		queryWatcher.Channel,
		queryWatcher.ID,
		name)
	// lock the entity table (to prevent any new activate/deactive events)
	// and iterate all entities, adding those which currently match the query
	// to the new UpdatedEntityList's Entities list.
	m.entityTable.mutex.Lock()
	for _, id := range m.entityTable.allocatedIDs {
		if q.Test(id, m) {
			updatedEntityListDebug("sending signal %d in "+
				"GetUpdatedActiveEntityList", id)
			list.EntityChannel <- int32(id)
		}
	}
	// unlock the entity table and return our list
	m.entityTable.mutex.Unlock()
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
	removeEntityQueryWatcherFromSliceByID(ID, &m.activeEntityWatchers)
}

// hold an entity for modification, queueing the modification when
// the function
// returns
func (m *EntityManager) AtomicEntityModify(
	id uint16,
	f func(e *EntityModification)) bool {

	genRequested := atomic.LoadUint32(&m.entityTable.gens[id])

	// wait to obtain heldForModificationLocks
	for !atomic.CompareAndSwapUint32(
		&m.entityTable.heldForModificationLocks[id], 0, 1) {
		// if we didn't manage to grab the lock, sleep for a good 8 frames
		// (128 ms, hardly a problem) so that even if there are several
		// goroutines currently sleeping in this loop, they won't starve the
		// physics/collision update which needs to sometimes hold this flag
		// as well. If we only wanted goroutines to atomically access
		// entities, it wouldn't matter how long we sleep, but since we want
		// to really atomically modify the state of the entity *across* the
		// goroutines *including* the priveleged goroutine of the GameScene
		// Update(), we need to make sure not to starve physics and collision
		time.Sleep(4 * FRAME_SLEEP)
	}
	// if the gen has changed by the time the lock was released, return false
	if atomic.LoadUint32(&m.entityTable.gens[id]) != genRequested {
		return false
	}
	// create a modification object whose address we will pass to f.
	// we add an entity token to it (which the caller cannot remove or write)
	mod := EntityModification{entity: EntityToken{id, genRequested}}
	// invoke the function, passing it a pointer to the modification object
	f(&mod)
	// determine the type of the modification and pass it to the appropriate
	// channel after type asserting its contained modification
	switch mod.Type {
	case ENTITY_STATE_MODIFICATION:
		m.stateModificationChannel <- entityStateModification{
			entity: mod.entity,
			State:  mod.Modification.(EntityState)}
	case ENTITY_COMPONENT_MODIFICATION:
		m.componentModificationChannel <- entityComponentModification{
			entity:     mod.entity,
			Components: mod.Modification.(ComponentSet)}
	}
	// let the caller know their modification went through
	return true
}

// used by collision system to hold two entities for modification
func (m *EntityManager) holdTwoEntities(i uint16, j uint16) bool {
	// attempt to hold i
	if atomic.CompareAndSwapUint32(
		&m.entityTable.heldForModificationLocks[i], 0, 1) {
		// attempt to hold j
		if atomic.CompareAndSwapUint32(
			&m.entityTable.heldForModificationLocks[j], 0, 1) {
			// if we're here, we have held both i and j. return true
			return true
		}
		// if we're here, we have held i but failed to hold j. let go of i
		// and return false
		atomic.StoreUint32(&m.entityTable.heldForModificationLocks[i], 0)
		return false
	}
	// if we're here, we failed to hold i
	return false
}

// used by collision system to release two entities for modification
func (m *EntityManager) releaseTwoEntities(i uint16, j uint16) {
	atomic.StoreUint32(&m.entityTable.heldForModificationLocks[i], 0)
	atomic.StoreUint32(&m.entityTable.heldForModificationLocks[j], 0)
}

// used by physics system to hold an entity for modification
func (m *EntityManager) holdEntity(i uint16) bool {
	return atomic.CompareAndSwapUint32(
		&m.entityTable.heldForModificationLocks[i], 0, 1)
}

// used by physics system to release an entity held for modification
func (m *EntityManager) releaseEntity(i uint16) {
	atomic.StoreUint32(&m.entityTable.heldForModificationLocks[i], 0)
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
func (m *EntityManager) EntityComponentBitArray(id uint16) bitarray.BitArray {
	m.entityTable.mutex.RLock()
	defer m.entityTable.mutex.RUnlock()

	return m.entityTable.componentBitArrays[id]
}

// Somewhat expensive conversion of entire entity list to string
func (m *EntityManager) String() string {
	m.UpdateMutex.Lock()
	m.tagTable.mutex.RLock()
	defer m.tagTable.mutex.RUnlock()
	defer m.UpdateMutex.Unlock()

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
