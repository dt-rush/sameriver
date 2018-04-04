/**
  *
  *
  *
  *
**/

package engine

import (
	"fmt"
)

type EntityManager struct {
	id_generator IDGenerator

	component_registry ComponentRegistry

	// TODO: implement, along with selector functions and logical operators...
	// entity_tag_bitarray // TODO: type?
	// tag_entity_bitarray // TODO: type?
	// entity_component_bitarray
	// component_entity_bitarray

	entities []int
	// two one-way maps support a many-to-many relationship
	// tag -> []IDs
	tag_entities map[string]([]int)
	// ID -> []tag
	entity_tags map[int]([]string)
}

func (m *EntityManager) Init(components []Component) {
	// 4 is arbitrary (could be tuned?). this should be expected to grow anyway
	m.id_generator.Init()
	// init component registry subsystem
	m.component_registry.Init(components)

	// init slices, maps
	m.entities = make([]int, 0)
	m.tag_entities = make(map[string]([]int))
	m.entity_tags = make(map[int]([]string))
}

func (m *EntityManager) Entities() []int {
	return m.entities
}

// ECS maxim: entities are just IDs! These two numbers, of entities and IDs,
// are always in sync, along with creation and destruction of entities
// and the freeing of their IDs (currently we just make entities inactive,
// but at a certain point, deletion code needs to exist) (TODO)
func (m *EntityManager) NumberOfEntities() int {
	return len(m.entities)
}

// given a list of components, spawn an entity with the default values
// returns the ID
func (m *EntityManager) SpawnEntity(components []Component) int {

	// LOG message
	Logger.Printf("spawning entity with components [")
	for _, component := range components {
		Logger.Printf("%s,", component.Name())
	}
	Logger.Printf("]")

	// generate an id
	id := m.id_generator.Gen()
	Logger.Printf(" #%d", id)
	// allocate component data storage
	for _, component := range components {
		component.Set(id, component.DefaultValue())
	}
	// register the entity and its components with the component_registry
	m.component_registry.RegisterEntity(id, components)
	// append ID to this entity manager's internal list of entities
	m.entities = append(m.entities, id)
	// return the created ID to the caller
	return id
}

// seems the array shift makes this inefficient.
// TODO this is probably thread-unsafe
func (m *EntityManager) DespawnEntity(id int) {
	for i := 0; i < len(m.entities); i++ {
		if i == id {
			// delete the element at i
			// (put it at the end and return a truncated by 1 list
			_i := len(m.entities) - 1
			m.entities[_i], m.entities[i] = m.entities[i], m.entities[_i]
			m.entities = m.entities[:_i]

		}
	}
}

// TAG SUPPORTING FUNCTIONS

func (m *EntityManager) TagEntity(id int, tag string) {
	_, et_ok := m.entity_tags[id]
	_, te_ok := m.tag_entities[tag]
	if !et_ok {
		m.entity_tags[id] = make([]string, 0)
	}
	if !te_ok {
		m.tag_entities[tag] = make([]int, 0)
	}
	m.entity_tags[id] = append(m.entity_tags[id], tag)
	m.tag_entities[tag] = append(m.tag_entities[tag], id)
}

func (m *EntityManager) TagEntityUnique(id int, tag string) {
	if len(m.tag_entities[tag]) != 0 {
		panic(fmt.Sprintf("trying to tag unique %s more than once", tag))
	}
	m.TagEntity(id, tag)
}

func (m *EntityManager) TagEntities(ids []int, tag string) {
	for _, id := range ids {
		m.TagEntity(id, tag)
	}
}

func (m *EntityManager) GetEntityTags(id int) []string {
	return m.entity_tags[id]
}

func (m *EntityManager) GetTagEntities(tag string) []int {
	return m.tag_entities[tag]
}

func (m *EntityManager) GetTagEntityUnique(tag string) int {
	entity_list := m.tag_entities[tag]
	if len(entity_list) == 0 {
		return -1
	} else {
		return entity_list[0]
	}
}

func (m *EntityManager) EntityHasTag(id int, tag string) bool {
	for _, entity_tag := range m.entity_tags[id] {
		if entity_tag == tag {
			return true
		}
	}
	return false
}

func (m *EntityManager) EntityHasComponent(id int, component Component) bool {
	return m.component_registry.EntityHas(id, component)
}
