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

type EntityManager struct {

	// bitarray used to keep track of which entities have which components
	// (indexes are IDs, bitarrays have bit set if entity has the
	// component corresponding to that index)
	EntityComponentBitArrays [MAX_ENTITIES]bitarray.BitArray
	// how many entities there are
	NumEntities int
	// list of IDs which have been allocated
	allocatedIDs []uint16
	// list of available entity ID's which have previously been deallocated
	availableIDs []uint16
	// to protedct modification of the above data
	entityTableMutex sync.Mutex

	// Component data
	Components ComponentsTable

	// used to allow systems to keep an updated list of entities which have
	// components they're interested in operating on (eg. physics watches
	// for entities with position, velocity, and hitbox)
	activeWatchers []EntityQueryWatcher
	// to protect modifying the above slice
	activeWatchersMutex sync.Mutex

	// data members to support the entity tagging system, which allows us to
	// associate a set of strings with an entity
	// tag -> []IDs
	EntitiesWithTag map[string]([]uint16)
	// ID -> []tag
	TagsOfEntity map[uint16]([]string)
	// to avoid race conditions
	tagSystemMutex sync.Mutex
}

func (m *EntityManager) Init() {
	// allocate component data
	m.Components = AllocateComponentsMemoryBlock()
	// allocate tag system data members
	m.EntitiesWithTag = make(map[string]([]uint16))
	m.TagsOfEntity = make(map[uint16]([]string))
}

// get the ID for a new entity
func (m *EntityManager) AllocateID() uint16 {
	// lock the entity table while we operate on it
	m.entityTableMutex.Lock()
	defer m.entityTableMutex.Unlock()
	// if there is a deallocated entity somewhere in the table before the
	// highest ID, return that ID to the caller
	n_avail := len(m.availableIDs)
	var id uint16
	if n_avail > 0 {
		// there is an ID available for a previously deallocated entity.
		// pop it from the list and continue with that as the ID
		id = m.availableIDs[n_avail-1]
		m.availableIDs = m.availableIDs[:n_avail-1]
	} else {
		// every slot in the table before the highest ID is filled
		m.NumEntities++
		id = uint16(m.NumEntities - 1)
	}
	// add the ID to the list of allocated IDs
	m.allocatedIDs = append(m.allocatedIDs, id)
	return id
}

func (m *EntityManager) Despawn(id uint16) {
	// Deactivate the entity to ensure that all updated entity lists are
	// notified
	m.Deactivate(id)
	// Lock the entity table while we operate on it
	m.entityTableMutex.Lock()
	defer m.entityTableMutex.Unlock()
	// add the ID to the list of available IDs
	m.availableIDs = append(m.availableIDs, id)
	// remove the ID from the list of allocated IDs
	last_ix := len(m.allocatedIDs) - 1
	for i := 0; i <= last_ix; i++ {
		if m.allocatedIDs[i] == id {
			m.allocatedIDs[i] = m.allocatedIDs[last_ix]
			m.allocatedIDs = m.allocatedIDs[:last_ix]
			break
		}
	}
	// clear the tags for entity
	tags_to_clear := m.TagsOfEntity[id]
	for _, tag_to_clear := range tags_to_clear {
		m.UntagEntity(id, tag_to_clear)
	}
	delete(m.TagsOfEntity, id)
	// unset the bitarray for the entity
	m.EntityComponentBitArrays[id].Reset()
}

func (m *EntityManager) DespawnAll() {
	Logger.Println("[Entity manager] Deallocating all...")
	Logger.Printf("[Entity manager] Currently allocated: %v\n", m.allocatedIDs)
	// iterate all IDs which could have been allocated
	for i := 0; i < m.NumEntities; i++ {
		// we can call this safely on each ID, even those unallocated,
		// since it's idempotent - it will exit early if the entity is already
		// deactivated
		m.Despawn(uint16(i))
	}
}

// given a list of components, spawn an entity with the default values
// returns the ID
func (m *EntityManager) SpawnEntity(
	id uint16,
	component_set ComponentSet) {

	Logger.Printf("[Entity manager] Spawning: %d\n", id)

	m.entityTableMutex.Lock()
	defer m.entityTableMutex.Unlock()

	// set the bitarray for this entity
	m.EntityComponentBitArrays[id] = component_set.ToBitArray()

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
}

// NOTE: do not call this if you've locked Components.Active for reading, haha
func (m *EntityManager) Activate(id uint16) {
	Logger.Printf("[Entity manager] Activating: %d\n", id)
	// Activate is idempotent - only enter if not active
	if !m.Components.Active.SafeGet(id) {
		// Set active = true
		m.Components.Active.SafeSet(id, true)
		// check if anybody has set a query watch on this entity.
		// If so, notify them of activate by sending this id
		// through the channel
		for _, watcher := range m.activeWatchers {
			if watcher.Query.Test(id, m) {
				// warn if the channel is full (we will block here if so)
				// NOTE: this can be very bad indeed, since now whatever
				// called Activate is blocking
				if len(watcher.Channel) == ACTIVE_ENTITY_WATCHER_CHANNEL_CAPACITY {
					Logger.Printf("[Entity manager] ⚠⚠⚠  active watcher channel %s is full, causing block in Activate()\n", watcher.Name)
				}
				watcher.Channel <- int16(id)
			}
		}
	}
}

func (m *EntityManager) Deactivate(id uint16) {
	Logger.Printf("[Entity manager] Deactivating: %d\n", id)
	// Deactivate is idempotent - only enter if active
	if m.Components.Active.SafeGet(id) {
		// Set active = false
		m.Components.Active.SafeSet(id, false)
		// check if anybody has set a query watch on this entity.
		// If so, notify them of activate by sending -(id + 1)
		// through the channel
		for _, watcher := range m.activeWatchers {
			if watcher.Query.Test(id, m) {
				if len(watcher.Channel) == ACTIVE_ENTITY_WATCHER_CHANNEL_CAPACITY {
					Logger.Printf("[Entity manager] ⚠  active watcher channel %s is full - blocking in Deactivate()\n", watcher.Name)
				}
				watcher.Channel <- -(int16(id + 1))
			}
		}
	}
}

func (m *EntityManager) GetUpdatedActiveList(
	q EntityQuery, name string) *UpdatedEntityList {

	return NewUpdatedEntityList(m.SetActiveWatcher(q, name), name)
}

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
	// find the index of the EntityQueryWatcher in the list and splice it out
	last_ix := len(m.activeWatchers) - 1
	for i := uint16(0); i <= uint16(last_ix); i++ {
		// channel equality is watcher equality, since watchers are created
		// at the same time as their channel (1:1 mapping)
		if m.activeWatchers[i].Channel == qw.Channel {
			m.activeWatchersMutex.Lock()
			m.activeWatchers[i] = m.activeWatchers[last_ix]
			m.activeWatchers = m.activeWatchers[:last_ix]
			m.activeWatchersMutex.Unlock()
			break
		}
	}
}

// apply the given tag to the given entity
func (m *EntityManager) TagEntity(id uint16, tag string) {
	Logger.Printf("[Entity manager] Tagging %d with: %s\n", id, tag)
	m.tagSystemMutex.Lock()
	defer m.tagSystemMutex.Unlock()
	_, t_of_e_exists := m.TagsOfEntity[id]
	_, e_with_t_exists := m.EntitiesWithTag[tag]
	if !t_of_e_exists {
		m.TagsOfEntity[id] = make([]string, 0)
	}
	if !e_with_t_exists {
		m.EntitiesWithTag[tag] = make([]uint16, 0)
	}
	m.TagsOfEntity[id] = append(m.TagsOfEntity[id], tag)
	m.EntitiesWithTag[tag] = append(m.EntitiesWithTag[tag], id)
}

// remove a tag from an entity
func (m *EntityManager) UntagEntity(id uint16, tag string) {
	m.tagSystemMutex.Lock()
	defer m.tagSystemMutex.Unlock()
	// remove the id from the list of entities with the tag
	last_ix := len(m.EntitiesWithTag[tag]) - 1
	for i := 0; i <= last_ix; i++ {
		if m.EntitiesWithTag[tag][i] == id {
			// thanks to https://stackoverflow.com/a/37359662 for this nice
			// little splice idiom when we don't care about slice order (saves
			// reallocating the whole dang thing)
			m.EntitiesWithTag[tag][i] = m.EntitiesWithTag[tag][last_ix]
			m.EntitiesWithTag[tag] = m.EntitiesWithTag[tag][:last_ix]
			break
		}
	}
	// remove the tag from the list of tags for the entity
	last_ix = len(m.TagsOfEntity[id]) - 1
	for i := 0; i <= last_ix; i++ {
		if m.TagsOfEntity[id][i] == tag {
			m.TagsOfEntity[id][i] = m.TagsOfEntity[id][last_ix]
			m.TagsOfEntity[id] = m.TagsOfEntity[id][:last_ix]
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
	for _, entity_tag := range m.TagsOfEntity[id] {
		if entity_tag == tag {
			return true
		}
	}
	return false
}

// Boolean check of whether a given entity has a given component
func (m *EntityManager) EntityHasComponent(id uint16, COMPONENT int) bool {
	b, _ := m.EntityComponentBitArrays[id].GetBit(uint64(COMPONENT))
	return b
}
