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
	// the highest index a registered entity resides at
	HighestID uint16
	// stack of available entity ID's < n_entities (when Deallocate is called
	// for an ID, we add to this stack)
	availableIDs []uint16
	// to avoid race conditions involving the modification of the above
	// when multiple goroutines may want to spawn or despawn entities
	entityTableMutex sync.Mutex

	// Component data
	Components ComponentsTable

	// used to allow systems to keep an updated list of entities which have
	// components they're interested in operating on (eg. physics watches
	// for entities with position, velocity, and hitbox)
	activeWatchers []QueryWatcher
	// to generate ID's for the active watchers
	watcherIDGen IDGenerator
	// to avoid race conditions
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
	m.entityTableMutex.Lock()
	defer m.entityTableMutex.Unlock()
	// if there is a deallocated entity somewhere in the table before the
	// highest ID, return that ID to the caller
	n_avail := len(m.availableIDs)
	if n_avail > 0 {
		id := m.availableIDs[n_avail-1]
		m.availableIDs = m.availableIDs[:n_avail-1]
		return id
	} else {
		// every slot in the table before the highest ID is filled.
		// Increment the highest ID (by setting it = to the number of entities,
		// which will be, given that the table is full up to this point,
		// highest_id + 1) and return it
		m.HighestID = uint16(m.NumEntities)
		m.NumEntities += 1
		return m.HighestID
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

// forget an entity existed (don't clear the data stored in any component
// tables though (these will simply be overwritten later)
func (m *EntityManager) DespawnEntity(id uint16) {
	m.entityTableMutex.Lock()
	defer m.entityTableMutex.Unlock()
	// set entity inactive
	m.Deactivate(id)
	// add the ID to the pool of now-available ID's
	m.availableIDs = append(m.availableIDs, id)
	// remove tag metadata for this entity
	tags_to_clear := m.TagsOfEntity[id]
	for _, tag_to_clear := range tags_to_clear {
		m.UntagEntity(id, tag_to_clear)
	}
	delete(m.TagsOfEntity, id)

	// unset the bitarray for the entity
	m.EntityComponentBitArrays[id].Reset()
}

// NOTE: do not call this if you've locked Components.Active for reading, haha
func (m *EntityManager) Activate(id uint16) {
	Logger.Printf("[Entity manager] Activating: %d\n", id)
	m.Components.Active.SafeSet(id, true)
	// check if anybody has set a query watch on the specific component mix
	// of this entity. If so, notify them of activate by sending this id
	// through the channel
	for _, watcher := range m.activeWatchers {
		if watcher.Query.Test(id, m) {
			watcher.Channel <- int16(id)
		}
	}
}

func (m *EntityManager) Deactivate(id uint16) {
	Logger.Printf("[Entity manager] Deactivating: %d\n", id)
	m.Components.Active.SafeSet(id, false)
	// check if anybody has set a query watch on the specific component mix
	// of this entity. If so, notify them of deactivate by sending -(id + 1)
	// through the channel
	for _, watcher := range m.activeWatchers {
		if watcher.Query.Test(id, m) {
			watcher.Channel <- -(int16(id + 1))
		}
	}
}

func (m *EntityManager) GetUpdatedActiveList(
	q Query, name string) *UpdatedEntityList {

	return NewUpdatedEntityList(m.SetActiveWatcher(q), name)
}

func (m *EntityManager) StopUpdatedActiveList(l UpdatedEntityList) {
	m.UnsetActiveWatcher(l.Watcher)
	l.StopUpdateChannel <- true
}

// Return a channel which will receive the id of an entity whenever an entity
/// becomes active with a component set matching the query bitarray, and which
// will receive -(id + 1) whenever an entity is *despawned* with a component
// set matching the query bitarray
func (m *EntityManager) SetActiveWatcher(q Query) QueryWatcher {

	// TODO: this seems as if we're just hoping that the capacity won't exceed
	// 8 for any reason, which it could if we spawn a lot of entities and the
	// channel readers aren't fast enough. At the least we should have
	// a system of wrapping channels with a given capacity such that if they
	// near their capacity at any point, we print a clear warning to the
	// console.
	// Rationale: if we try to push to a channel whose buffer is full, we'll
	// block. And that can be quite bad in some circumstances. Without a nice
	// and thorough analysis of all the dependent flows of channel
	// sends/reads, in fact, we could end up in a deadlock through some
	// obscure condition we hadn't forseen

	c := make(chan (int16), 8)
	watcherID := m.watcherIDGen.Gen()
	qw := QueryWatcher{q, c, watcherID}
	m.activeWatchersMutex.Lock()
	m.activeWatchers = append(m.activeWatchers, qw)
	m.activeWatchersMutex.Unlock()
	return qw
}

func (m *EntityManager) UnsetActiveWatcher(qw QueryWatcher) {
	// find the index of the QueryWatcher in the list and splice it out
	last_ix := len(m.activeWatchers) - 1
	for i := uint16(0); i <= uint16(last_ix); i++ {
		if m.activeWatchers[i].ID == qw.ID {
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
