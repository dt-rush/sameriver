/**
  *
  *
  *
  *
**/

package engine

type Component interface {
	Init(capacity int, game *Game)
	// val will be type-asserted by implementers
	Set(id int, val interface{})
	// to be type-asserted by receivers
	DefaultValue() interface{}
	String() string
	// TODO be sure to document this somewhere other than right
	// in this comment:
	// MUST be unique (for engine/component_registry.go)
	// TODO: implement via reflection?
	Name() string
}
