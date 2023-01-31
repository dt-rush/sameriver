package sameriver

import (
	"fmt"
	"strings"
)

type ComponentSet struct {
	// names of all components given values in this set
	names map[string]bool
	// data storage
	vec2DMap   map[string]Vec2D
	boolMap    map[string]bool
	intMap     map[string]int
	float64Map map[string]float64
	stringMap  map[string]string
	spriteMap  map[string]Sprite
	tagListMap map[string]TagList
	genericMap map[string]interface{}
	customMap  map[string]interface{}
	cccs       map[string]CustomContiguousComponent
}

func MakeCustomComponentSet(
	input map[string]interface{},
	customs map[string]interface{},
	cccs map[string]CustomContiguousComponent) ComponentSet {

	baseCS := MakeComponentSet(input)
	for spec, value := range customs {
		// decode spec string
		split := strings.Split(spec, ",")
		kind := split[0]
		name := split[1]
		if kind != "Custom" {
			panic(fmt.Sprintf("custom component spec should have type Custom, got: %s", kind))
		}
		// take note in names map that this component name occurs
		baseCS.names[name] = true
		baseCS.customMap[name] = value
		// store the interface object itself so ComponentTable.ApplyComponentSet()
		// can call its ApplyToEntity() function to set the value
		baseCS.cccs[name] = cccs[name]
	}
	return baseCS
}

// takes as input a map whose keys are components specified by {kind},{name}
// and whose values are interface{} for the value
func MakeComponentSet(input map[string]interface{}) ComponentSet {
	cs := ComponentSet{
		names:      make(map[string]bool),
		vec2DMap:   make(map[string]Vec2D),
		boolMap:    make(map[string]bool),
		intMap:     make(map[string]int),
		float64Map: make(map[string]float64),
		stringMap:  make(map[string]string),
		spriteMap:  make(map[string]Sprite),
		tagListMap: make(map[string]TagList),
		genericMap: make(map[string]interface{}),
		customMap:  make(map[string]interface{}),
		cccs:       make(map[string]CustomContiguousComponent),
	}
	for spec, value := range input {
		// decode spec string
		split := strings.Split(spec, ",")
		kind := split[0]
		name := split[1]
		// take note in names map that this component name occurs
		cs.names[name] = true
		// assign values into appropriate maps
		switch kind {
		case "Vec2D":
			if v, ok := value.(Vec2D); ok {
				cs.vec2DMap[name] = v
			}
		case "Bool":
			if b, ok := value.(bool); ok {
				cs.boolMap[name] = b
			}
		case "Int":
			if i, ok := value.(int); ok {
				cs.intMap[name] = i
			}
		case "Float64":
			if f, ok := value.(float64); ok {
				cs.float64Map[name] = f
			}
		case "String":
			if s, ok := value.(string); ok {
				cs.stringMap[name] = s
			}
		case "Sprite":
			if s, ok := value.(Sprite); ok {
				cs.spriteMap[name] = s
			}
		case "TagList":
			if t, ok := value.(TagList); ok {
				cs.tagListMap[name] = t
			}
		case "Generic":
			cs.genericMap[name] = value
		default:
			panic(fmt.Sprintf("unknown component kind %s", kind))
		}
	}
	return cs
}