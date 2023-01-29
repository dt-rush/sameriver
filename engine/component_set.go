package engine

import (
	"strings"
)

type ComponentSet struct {
	// names of all components given values in this set
	names map[string]bool
	// data storage
	vec2DMap     map[string]Vec2D
	vec3DMap     map[string]Vec3D
	logicUnitMap map[string]LogicUnit
	intMap       map[string]int
	float64Map   map[string]float64
	spriteMap    map[string]Sprite
	tagListMap   map[string]TagList
	genericMap   map[string]interface{}
}

// takes as input a map whose keys are components specified by {kind},{name}
// and whose values are interface{} for the value
func MakeComponentSet(input map[string]interface{}) *ComponentSet {
	cs := &ComponentSet{
		vec2DMap:     make(map[string]Vec2D),
		vec3DMap:     make(map[string]Vec3D),
		logicUnitMap: make(map[string]LogicUnit),
		intMap:       make(map[string]int),
		float64Map:   make(map[string]float64),
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
			if vec2d, ok := value.(Vec2D); ok {
				cs.vec2DMap[name] = vec2d
			}
		case "Vec3D":
			if vec3d, ok := value.(Vec3D); ok {
				cs.vec3DMap[name] = vec3d
			}
		case "LogicUnit":
			if logicUnit, ok := value.(LogicUnit); ok {
				cs.logicUnitMap[name] = logicUnit
			}
		case "Int":
			if integer, ok := value.(int); ok {
				cs.intMap[name] = integer
			}
		case "Float64":
			if float, ok := value.(float64); ok {
				cs.float64Map[name] = float
			}
		case "Sprite":
			if sprite, ok := value.(Sprite); ok {
				cs.spriteMap[name] = sprite
			}
		case "TagList":
			if tagList, ok := value.(TagList); ok {
				cs.tagListMap[name] = tagList
			}
		case "Generic":
			cs.genericMap[name] = value
		default:
			panic(fmt.Sprintf("unknown component kind %s", kind))
		}
	}
	return cs
}
