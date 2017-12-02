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
	
	id_system IDSystem
	
	tag_system TagSystem

	component_registry ComponentRegistry

	// TODO: implement, along with selector functions and logical operators...
	// entity_tag_bitarray // TODO: type?
	// tag_entity_bitarray // TODO: type?
	// entity_component_bitarray 
	// component_entity_bitarray 
	
	entities []int
}

func (m *EntityManager) Init (capacity int, components []Component) {
	// init ID subsystem
	m.id_system.Init (capacity)
	// init tag subsystem
	m.tag_system.Init (capacity)
	// init component registry subsystem
	m.component_registry.Init (components)
}


func (m *EntityManager) Entities () []int {
	return m.entities
}



// ECS maxim: entities are just IDs! These two numbers, of entities and IDs,
// are always in sync, along with creation and destruction of entities
// and the freeing of their IDs (currently we just make entities inactive,
// but at a certain point, deletion code needs to exist) (TODO)
func (m *EntityManager) NumberOfEntities () int {
	return len (m.entities)
}


// given a list of components, spawn an entity with the default values
// returns the ID
func (m *EntityManager) SpawnEntity (components []Component) int {

	// LOG message
	fmt.Printf ("spawning entity with components [")
	for _, component := range components {
		fmt.Printf ("%s,", component.Name())
	}
	fmt.Printf ("]")

	// generate an id
	id := m.id_system.Gen()
	fmt.Printf (" #%d\n", id)
	// allocate component data storage
	for _, component := range components {
		component.Set (id, component.DefaultValue())
	}
	// register the entity and its components with the component_registry
	m.component_registry.RegisterEntity (id, components)
	// append ID to this entity manager's internal list of entities
	m.entities = append (m.entities, id)
	// return the created ID to the caller
	return id
}

// given a list of lists of components, return the result of
// mapping spawn_entity() over the list of lists of components
// returns a list of IDs, therefore
func (m *EntityManager) SpawnEntities (args [][]Component) []int {
	ids := make ([]int, len (args))
	for _, arg := range args {
		ids = append (ids, m.SpawnEntity (arg))
	}
	return ids
}

// seems the array shift makes this inefficient.
// TODO this is probably thread-unsafe
func (m *EntityManager) DespawnEntity (id int) {
	for i := 0; i < len (m.entities); i++ {
		if i == id {
			// delete the element at i
			// (put it at the end and return a truncated by 1 list
			_i := len (m.entities) - 1 
			m.entities [_i], m.entities [i] = m.entities [i], m.entities [_i]
			m.entities = m.entities [:_i]

		}
	}
}



// TODO: refactor somehow
// TAG SUPPORTING FUNCTIONS


func (m *EntityManager) TagEntity (id int, tag string) {
	m.tag_system.TagEntity (id, tag)
}

func (m *EntityManager) TagEntities (ids []int, tag string) {
	for _, id := range (ids) {
		m.tag_system.TagEntity (id, tag)
	}
}

func (m *EntityManager) GetEntityTags (id int) []string {
	return m.tag_system.entity_tags [id]
}

func (m *EntityManager) GetTagEntities (tag string) []int {
	return m.tag_system.tag_entities [tag]
}




// TODO: refactor somehow
// COMPONENT REGISTRY FUNCTIONS

func (m *EntityManager) EntityHas (id int, component Component) bool {
	return m.component_registry.EntityHas (id, component)
}
