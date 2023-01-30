package engine

import (
	"fmt"
	"strings"
)

type ComponentSet struct {
	// names of all components given values in this set
	names map[string]bool
	// data storage
	vec2DMap     map[string]Vec2D
	logicUnitMap map[string]*LogicUnit
	boolMap      map[string]bool
	intMap       map[string]int
	float64Map   map[string]float64
	stringMap    map[string]string
	spriteMap    map[string]Sprite
	tagListMap   map[string]TagList
	genericMap   map[string]interface{}
}

// takes as input a map whose keys are components specified by {kind},{name}
// and whose values are interface{} for the value
func MakeComponentSet(input map[string]interface{}) ComponentSet {
	cs := ComponentSet{
		names:        make(map[string]bool),
		vec2DMap:     make(map[string]Vec2D),
		logicUnitMap: make(map[string]*LogicUnit),
		boolMap:      make(map[string]bool),
		intMap:       make(map[string]int),
		float64Map:   make(map[string]float64),
		stringMap:    make(map[string]string),
		spriteMap:    make(map[string]Sprite),
		tagListMap:   make(map[string]TagList),
		genericMap:   make(map[string]interface{}),
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
		case "*LogicUnit":
			if l, ok := value.(*LogicUnit); ok {
				cs.logicUnitMap[name] = l
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
