/**
  *
  * Interface defining a component, essentially a map from int ID's to data
  *
  *
**/

package engine

type Component interface {
	// initialize the internal data structure(s) of the component
	Init(capacity int, game *Game)
	// set the value for an ID
	Set(id int, val interface{})
	// returns the default value for the component
	// (ie. the point (0,0) for position)
	DefaultValue() interface{}
	// convert the component to string for printing a display / summary
	String() string
	// print the name of the component
	// TODO be sure to document this somewhere other than right
	// in this comment:
	// MUST be unique (for engine/component_registry.go)
	// TODO: implement via reflection?
	Name() string
}
