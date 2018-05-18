/**
  *
  * Manages the spawning and querying of entities
  *
**/

package engine

import (
	"sync"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityTable struct {
	// guards all state
	mutex sync.RWMutex
	// bitarray used to keep track of which entities have which components
	// (indexes are IDs, bitarrays have bit set if entity has the
	// component corresponding to that index)
	componentBitArrays [MAX_ENTITIES]bitarray.BitArray
	// how many entities there are
	numEntities int
	// list of IDs which have been allocated
	allocatedIDs []uint16
	// list of available entity ID's which have previously been deallocated
	availableIDs []uint16
}

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

type EntityManager struct {
	// Entity table stores component bitarrays, a list of allocated IDs,
	// and a list of available IDs from previous deallocations
	entityTable EntityTable
	// Component data
	Components ComponentsTable
	// Tag table stores data for entity tagging system
	tagTable TagTable

	// used to allow systems to keep an updated list of entities which have
	// components they're interested in operating on (eg. physics watches
	// for entities with position, velocity, and hitbox)
	activeWatchers []EntityQueryWatcher
	// to protect modifying the above slice
	activeWatchersMutex sync.RWMutex
}

func (m *EntityManager) Init() {
	// allocate component data
	m.Components = AllocateComponentsMemoryBlock()
	// allocate tag system data members
	m.tagTable.entitiesWithTag = make(map[string]([]uint16))
	m.tagTable.tagsOfEntity = make(map[uint16]([]string))
}

// get the ID for a new entity
func (m *EntityManager) allocateID() uint16 {
	m.entityTable.mutex.Lock()
	defer m.entityTable.mutex.Unlock()

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
	return id
}

func (m *EntityManager) Despawn(id uint16) {
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
	m.Deactivate(id)
	// clear the entity from lists of tagged entities it's in
	tags_to_clear := m.tagTable.tagsOfEntity[id]
	for _, tag_to_clear := range tags_to_clear {
		m.UntagEntity(id, tag_to_clear)
	}
	// remove the taglist for this entity
	delete(m.tagTable.tagsOfEntity, id)
}

func (m *EntityManager) DespawnAll() {
	Logger.Println("[Entity manager] Despawning all...")
	Logger.Printf("[Entity manager] Currently spawned: %v\n",
		m.entityTable.allocatedIDs)
	// iterate all IDs which could have been allocated
	for i := 0; i < m.entityTable.numEntities; i++ {
		// we can call this safely on each ID, even those unallocated,
		// since it's idempotent - it will exit early if the entity is already
		// deactivated
		m.Despawn(uint16(i))
	}
}

// given a list of components, spawn an entity with the default values
// returns the ID
func (m *EntityManager) SpawnEntity(
	component_set ComponentSet,
	tags []string) uint16 {

	// get an ID for the entity (locks and then unlocks the entityTable)
	id := m.allocateID()

	// lock the entityTable
	m.entityTable.mutex.Lock()
	defer m.entityTable.mutex.Unlock()

	Logger.Printf("[Entity manager] Spawning: %d\n", id)

	// set the bitarray for this entity
	m.entityTable.componentBitArrays[id] = component_set.ToBitArray()

	// copy the data into the component storage for each component
	// (note: we dereference the pointers, this is a real copy, so it's good
	// that component values are either small pieces of data like [2]uint16
	// or a pointer to a func, etc.).
	// We don't "zero" the values of components not in the entity's set,
	// because really if a system operating on the component data
	// expects to work on the data, it should be maintaining a list of
	// entities with the required components using an UpdatedEntityList

	m.Components.Active.SafeSet(id, false)

	if component_set.Color != nil {
		m.Components.Color.SafeSet(id, *(component_set.Color))
	}
	if component_set.Hitbox != nil {
		m.Components.Hitbox.SafeSet(id, *(component_set.Hitbox))
	}
	if component_set.Logic != nil {
		m.Components.Logic.SafeSet(id, *(component_set.Logic))
	}
	if component_set.Position != nil {
		m.Components.Position.SafeSet(id, *(component_set.Position))
	}
	if component_set.Sprite != nil {
		m.Components.Sprite.SafeSet(id, *(component_set.Sprite))
	}
	if component_set.Velocity != nil {
		m.Components.Velocity.SafeSet(id, *(component_set.Velocity))
	}

	return id
}

// Returns the component bit array for an entity
func (m *EntityManager) EntityComponentBitArray(id uint16) bitarray.BitArray {
	m.entityTable.mutex.RLock()
	defer m.entityTable.mutex.RUnlock()
	return m.entityTable.componentBitArrays[id]
}

// sets an entity active and notifies all watchers
func (m *EntityManager) Activate(id uint16) {
	Logger.Printf("[Entity manager] Activating: %d\n", id)
	m.setActiveState(id, true)
}

// sets an entity inactive and notifies all watchers
func (m *EntityManager) Deactivate(id uint16) {
	Logger.Printf("[Entity manager] Deactivating: %d\n", id)
	m.setActiveState(id, false)
}

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) setActiveState(id uint16, state bool) {
	// Only set (and notify) if not already in given state
	if m.Components.Active.SafeGet(id) != state {
		m.Components.Active.SafeSet(id, state)
		go m.notifyActiveState(id, state)
	}
}

// Send a signal to all registered watchers that an entity has a certain
// active state, either true or false
func (m *EntityManager) notifyActiveState(id uint16, active bool) {
	if !active {
		id = -(id + 1)
	}
	for _, watcher := range m.activeWatchers {
		if !watcher.Query.Test(id, m) {
			// warn if the channel is full (we will block here if so)
			// NOTE: this can be very bad indeed, since now whatever
			// called Activate is blocking
			if len(watcher.Channel) == ACTIVE_ENTITY_WATCHER_CHANNEL_CAPACITY {
				Logger.Printf("[Entity manager] ⚠  active watcher channel "+
					"%s is full, causing block in goroutine for "+
					"NotifyActiveState(%d, %v)\n",
					watcher.Name, id, active)
			}
			watcher.Channel <- int16(id)
		}
	}
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedActiveList(
	q EntityQuery, name string) *UpdatedEntityList {

	return NewUpdatedEntityList(m.SetActiveWatcher(q, name), name)
}

// Stops
func (m *EntityManager) StopUpdatedActiveList(l UpdatedEntityList) {
	m.UnsetActiveWatcher(l.Watcher)
	l.StopUpdateChannel <- true
}

// Return a channel which will receive the id of an entity whenever an entity
/// becomes active with a component set matching the query bitarray, and which
// will receive -(id + 1) whenever an entity is *despawned* with a component
// set matching the query bitarray
func (m *EntityManager) SetActiveWatcher(
	q EntityQuery, name string) EntityQueryWatcher {

	c := make(chan (int16), ACTIVE_ENTITY_WATCHER_CHANNEL_CAPACITY)
	qw := EntityQueryWatcher{q, c, name}
	m.activeWatchersMutex.Lock()
	m.activeWatchers = append(m.activeWatchers, qw)
	m.activeWatchersMutex.Unlock()
	return qw
}

func (m *EntityManager) UnsetActiveWatcher(qw EntityQueryWatcher) {
	m.activeWatchersMutex.Lock()
	defer m.activeWatchersMutex.Unlock()
	// remove the EntityQueryWatcher from the list of active watchers
	removeEntityQueryWatcherFromSlice(qw, &m.activeWatchers)
}

// apply the given tag to the given entity
func (m *EntityManager) TagEntity(id uint16, tag string) {
	m.tagTable.mutex.Lock()
	defer m.tagTable.mutex.Unlock()

	Logger.Printf("[Entity manager] Tagging %d with: %s\n", id, tag)
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

	Logger.Printf("[Entity manager] Removing tag %s from %d\n", tag, id)
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
	for _, id := range ids {
		m.TagEntity(id, tag)
	}
}

// Boolean check of whether a given entity has a given tag
func (m *EntityManager) EntityHasTag(id uint16, tag string) bool {
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
