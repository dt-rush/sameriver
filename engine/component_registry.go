/**
  *
  * Used to implement efficient querying of components registered for each
  * entity ID. Each entity will be associated with a set of bits / bools
  * which will correspond, by position, to each of a set of registered
  * components. Thus, it becomes possible to query whether an entity has a
  * component, or get the full set of components an entity has.
  *
  *
**/

package engine

import (
// TODO: add bitarray implementation
)

// TODO:
// this type alias is to cover for the fact that bitarray aint' a type til
// we import an implementation. []bool is shite compared to a proper
// bitarray lib
type bitarray []bool

// The registry itself, containing the map of int -> bitarray
type ComponentRegistry struct {
	// the real data, a bitarray for the existence of a component
	// on a given entity
	// length of bitarrays determined by N components
	data map[int]bitarray
	// supporting map of the components to their
	// indexes in the bitarrays uses the fact that Name()
	// is unique (has to be!!!)
	component_indexes map[string]int
	// current bitarray size (number of components)
	bitarray_sz int
}

// Init takes an array of Components (that is, we can't add more dynamically
// later - the set of Components we register with Init is the set of all
// Components which entities can have)
func (r *ComponentRegistry) Init(components []Component) {
	// init id->bitarray storage
	r.bitarray_sz = len(components)
	r.data = make(map[int]bitarray)
	// init component -> bitarray-index supporting map
	r.component_indexes = make(map[string]int)
	for i, component := range components {
		r.component_indexes[component.Name()] = i
	}
}

// register an entity and its list of components
func (r *ComponentRegistry) RegisterEntity(id int, components []Component) {
	// mark this, horatio: the bitarray IS the entity, or its
	// representation, for this registry, hence the variable name
	// convention (seen again in other methods)
	// it be madness, yet there's method in't
	entity := make(bitarray, r.bitarray_sz)
	for _, component := range components {
		// note we assume no index failures because
		// we assume all components submitted are registered
		// already with Init() or the forthcoming (TODO: implement)
		// RegisterComponent for dynamic insertion of components
		entity[r.component_indexes[component.Name()]] = true
	}
	r.data[id] = entity
}

// func (r *ComponentRegistry) RegisterComponent (component Component) {
//  // register a component even though there are no entities which have it
//  // as yet
//  r.bitarray_sz += 1
//  // TODO implement the actual growing of the existing bitarrays for damn sake
// }

func (r *ComponentRegistry) EntityHas(id int, component Component) bool {
	return r.data[id][r.component_indexes[component.Name()]]
}

func (r *ComponentRegistry) GetEntity(id int) bitarray {
	entity, ok := r.data[id]
	if !ok {
		Logger.Printf("[ComponentRegistry] attempt to get #%d, not in data", id)
		Logger.Printf("[ComponentRegistry] returning empty bitarray")
		return make(bitarray, r.bitarray_sz)
	}
	return entity
}
