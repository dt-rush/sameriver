/**
  *
  * Manages the spawning and querying of entities
  *
**/

package engine

import (
	"fmt"
	"sync"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityManager struct {

	// bitarray used to keep track of which entities have which components
	// (indexes are IDs, bitarrays have bit set if entity has the
	// component corresponding to that index)
	EntityComponentBitarrays [MAX_ENTITIES]bitarray.BitArray
	// how many entities there are
	NumEntities int
	// the highest index a registered entity resides at
	HighestID int
	// stack of available entity ID's < n_entities (when Deallocate is called
	// for an ID, we add to this stack)
	availableIDs []int
	// to avoid race conditions involving the modification of the above
	// when multiple goroutines may want to spawn or despawn entities
	entityTableMutex sync.Mutex

	// Component data
	Components ComponentsTable

	// used to allow systems to keep an updated list of entities which have
	// components they're interested in operating on (eg. physics watches
	// for entities with position, velocity, and hitbox)
	activeWatchers []QueryWatcher
	// to avoid race conditions
	activeWatchersMutex sync.Mutex

	// data members to support the entity tagging system, which allows us to
	// associate a set of strings with an entity
	// tag -> []IDs
	EntitiesWithTag map[string]([]int)
	// ID -> []tag
	TagsOfEntity map[int]([]string)
	// to avoid race conditions
	tagSystemMutex sync.Mutex
}

func (m *EntityManager) Init() {
	// allocate component data
	m.components = AllocateComponentsMemoryBlock()
	// allocate tag system data members
	m.EntitiesWithTag = make(map[string]([]int))
	m.TagsOfEntity = make(map[int]([]string))
}

// get the ID for a new entity
func (m *EntityManager) AllocateID() int {
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
		m.HighestID = m.NumEntities
		m.NumEntities += 1
		return m.HighestID
	}
}

// given a list of components, spawn an entity with the default values
// returns the ID
func (m *EntityManager) SpawnEntity(
	id int,
	component_set ComponentSet) int {

	m.entityTableMutex.Lock()
	defer m.entityTableMutex.Unlock()

	// set the bitarray for this entity
	m.EntityComponentBitarrays[id] = component_set.ToBitArray()

	// copy the data into the component storage for each component
	// (note: we dereference the pointers, this is a real copy, so it's good
	// that component values are either small pieces of data like [2]uint16
	// or a pointer to a func, etc.). We "zero" the values as best we can,
	// although really if a system operating on the component data
	// expects to work on the data, it should be maintaining a list of
	// entities with the required components according to the system of
	// activeWatchers

	m.components.Active.SafeSet(id, false)

	if component_set.Color == nil {
		component_set.Color = sdl.Color{}
	}
	m.components.Color.SafeSet(id, *(component_set.Color))

	if component_set.Hitbox == nil {
		component_set.Hitbox = [2]uint16{}
	}
	m.components.Hitbox.SafeSet(id, *(component_set.Hitbox))

	if component_set.Logic == nil {
		component_set.Logic = nil
	}
	m.components.Logic.SafeSet(id, *(component_set.Logic))

	if component_set.Position == nil {
		component_set.Position = [2]uint16{}
	}
	m.components.Position.SafeSet(id, *(component_set.Position))

	if component_set.Sprite == nil {
		component_set.Sprite = Sprite{}
	}
	m.components.Sprite.SafeSet(id, *(component_set.Sprite))

	if component_set.Velocity == nil {
		component_set.Velocity = [2]uint16{}
	}
	m.components.Velocity.SafeSet(id, *(component_set.Velocity))

	// return the created ID to the caller
	return id
}

// forget an entity existed (don't clear the data stored in any component
// tables though (these will simply be overwritten later)
func (m *EntityManager) DespawnEntity(id int) {
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
	m.EntityComponentBitarrays[id].Reset()
}

func Activate (id int) {
	m.Components.Active.SafeSet (id, true)
	// check if anybody has set a query watch on the specific component mix
	// of this entity. If so, notify them of activate by sending this id through
	// the channel
	for _, watcher := range activeWatchers {
		if watcher.Query.Equals(component_set.ToBitArray()) {
			watcher.Channel <- id
		}
	}
}

func Deactivate (id int) {
	m.Components.Active.SafeSet (id, false)
	// check if anybody has set a query watch on the specific component mix
	// of this entity. If so, notify them of deactivate by sending -(id + 1)
	// through the channel
	for _, watcher := range activeWatchers {
		if watcher.Query.Equals(component_set.ToBitArray()) {
			watcher.Channel <- -(id + 1)
		}
	}
}

func (m *EntityManager) GetUpdatedActiveList (query bitarray.BitArray) UpdatedEntityList {
	return NewUpdatedEntityList (m.SetActiveWatcher(query))
}

func (m *EntityManager) StopUpdatedActiveList (l UpdatedEntityList) {
	m.UnsetActiveWatcher (l.Watcher)
	m.StopUpdateChannel <- true
}

// Return a channel which will receive the id of an entity whenever an entity
// becomes active with a component set matching the query bitarray, and which
// will receive -(id + 1) whenever an entity is *despawned* with a component set
// matching the query bitarray
func (m *EntityManager) SetActiveWatcher(
	query bitarray.BitArray) QueryWatcher {

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

	c := make(chan (int), 8)
	watcherID := m.watcherIDGen.Gen()
	qw := QueryWatcher{query, c, watcherID}
	m.activeWatchersMutex.Lock()
	m.activeWatchers = append(m.activeWatchers, qw)
	m.activeWatchersMutex.Unlock()
	return qw
}

func (m *EntityManager) UnsetActiveWatcher(qw QueryWatcher) {
	// find the index of the QueryWatcher in the list and splice it out
	last_ix = len(m.activeWatchers) - 1
	for i := 0; i <= last_ix; i++ {
		if i == qw.ID {
			m.activeWatchersMutex.Lock()
			m.activeWatchers[i] = m.activeWatchers[last_ix]
			m.activeWatchers = m.activeWatchers[:last_ix]
			m.activeWatchersMutex.Unlock()
		}
	}
}

// apply the given tag to the given entity
func (m *EntityManager) TagEntity(id int, tag string) {
	m.tagSystemMutex.Lock()
	defer m.tagSystemMutex.Unlock()
	_, t_of_e_exists := m.TagsOfEntity[id]
	_, e_with_t_exists := m.EntitiesWithTag[tag]
	if !t_of_e_exists {
		m.TagsOfEntity[id] = make([]string, 0)
	}
	if !e_with_t_exists {
		m.EntitiesWithTag[tag] = make([]int, 0)
	}
	m.TagsOfEntity[id] = append(m.TagsOfEntity[id], tag)
	m.EntitiesWithTag[tag] = append(m.EntitiesWithTag[tag], id)
}

// remove a tag from an entity
func (m *EntityManager) UntagEntity(id int, tag string) {
	m.tagSystemMutex.Lock()
	defer m.tagSystemMutex.Unlock()
	// remove the id from the list of entities with the tag
	id_list := &m.EntitiesWithTag[tag]
	for i := 0; i < len(id_list); i++ {
		if id_list[i] == id {
			// thanks to https://stackoverflow.com/a/37359662 for this nice
			// little splice idiom when we don't care about slice order (saves
			// reallocating the whole dang thing)
			last_ix = len(id_list) - 1
			id_list[i] = id_list[last_ix]
			id_list = id_list[:last_ix]
		}
	}
	// remove the tag from the list of tags for the entity
	tag_list := &m.TagsOfEntity[id]
	for i := 0; i < len(tag_list); i++ {
		if tag_list[i] == tag {
			last_ix = len(tag_list) - 1
			tag_list[i] = tag_list[last_ix]
			tag_list = tag_list[:last_ix]
		}
	}
}

// Tag each of the entities in the provided array of ID's with the given tag
func (m *EntityManager) TagEntities(ids []int, tag string) {
	for _, id := range ids {
		m.TagEntity(id, tag)
	}
}

// Tag an entity uniquely. Panic if another entity is already tagged (this is
// probably not a good thing to do, TODO: find a better way to guard unique)
func (m *EntityManager) TagEntityUnique(id int, tag string) {
	if len(m.EntitiesWithTag[tag]) != 0 {
		panic(fmt.Sprintf("trying to TagEntityUnique for [%s] more than once",
			tag))
	}
	m.TagEntity(id, tag)
}

// Get the ID of the unique entity, returning -1 if no entity has that tag
func (m *EntityManager) GetTagEntityUnique(tag string) int {
	entity_list := m.EntitiesWithTag[tag]
	if len(entity_list) == 0 {
		return -1
	} else {
		return entity_list[0]
	}
}

// Boolean check of whether a given entity has a given tag
func (m *EntityManager) EntityHasTag(id int, tag string) bool {
	for _, entity_tag := range m.TagsOfEntity[id] {
		if entity_tag == tag {
			return true
		}
	}
	return false
}

// Boolean check of whether a given entity has a given component
func (m *EntityManager) EntityHasComponent(id int, COMPONENT int) bool {
	b, _ := m.EntityComponentBitArrays[id].GetBit(COMPONENT)
	return b
}
