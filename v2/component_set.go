package sameriver

import (
	"fmt"
	"strings"
	"time"

	"github.com/dt-rush/sameriver/v2/utils"
)

type ComponentSet struct {
	// names of all components given values in this set
	names map[string]bool
	// data storage
	vec2DMap             map[string]Vec2D
	boolMap              map[string]bool
	intMap               map[string]int
	float64Map           map[string]float64
	timeMap              map[string]time.Time
	timeAccumulatorMap   map[string]utils.TimeAccumulator
	stringMap            map[string]string
	spriteMap            map[string]Sprite
	tagListMap           map[string]TagList
	intMapMap            map[string]IntMap
	floatMapMap          map[string]FloatMap
	genericMap           map[string]any
	customComponentsMap  map[string]any
	customComponentsImpl map[string]CustomContiguousComponent
}

func makeCustomComponentSet(
	componentSpecs map[string]any,
	customComponentSpecs map[string]any,
	customComponentsImpl map[string]CustomContiguousComponent) ComponentSet {

	baseCS := makeComponentSet(componentSpecs)
	for spec, value := range customComponentSpecs {
		// decode spec string
		split := strings.Split(spec, ",")
		kind := split[0]
		name := split[1]
		if kind != "Custom" {
			panic(fmt.Sprintf("custom component spec should have type Custom, got: %s", kind))
		}
		// take note in names map that this component name occurs
		baseCS.names[name] = true
		baseCS.customComponentsMap[name] = value
		// store the ccc implementation interface object itself so
		// ComponentTable.applyComponentSet() can call its Set() function to
		// set the value
		baseCS.customComponentsImpl[name] = customComponentsImpl[name]
	}
	return baseCS
}

// takes as componentSpecs a map whose keys are components specified by {kind},{name}
// and whose values are any for the value
func makeComponentSet(componentSpecs map[string]any) ComponentSet {
	cs := ComponentSet{
		names: make(map[string]bool),
	}
	for spec, value := range componentSpecs {
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
				if cs.vec2DMap == nil {
					cs.vec2DMap = make(map[string]Vec2D)
				}
				cs.vec2DMap[name] = v
			}
		case "Bool":
			if b, ok := value.(bool); ok {
				if cs.boolMap == nil {
					cs.boolMap = make(map[string]bool)
				}
				cs.boolMap[name] = b
			}
		case "Int":
			if i, ok := value.(int); ok {
				if cs.intMap == nil {
					cs.intMap = make(map[string]int)
				}
				cs.intMap[name] = i
			}
		case "Float64":
			if f, ok := value.(float64); ok {
				if cs.float64Map == nil {
					cs.float64Map = make(map[string]float64)
				}
				cs.float64Map[name] = f
			}
		case "Time":
			if t, ok := value.(time.Time); ok {
				if cs.timeMap == nil {
					cs.timeMap = make(map[string]time.Time)
				}
				cs.timeMap[name] = t
			}
		case "TimeAccumulator":
			if t, ok := value.(utils.TimeAccumulator); ok {
				if cs.timeAccumulatorMap == nil {
					cs.timeAccumulatorMap = make(map[string]utils.TimeAccumulator)
				}
				cs.timeAccumulatorMap[name] = t
			}
		case "String":
			if s, ok := value.(string); ok {
				if cs.stringMap == nil {
					cs.stringMap = make(map[string]string)
				}
				cs.stringMap[name] = s
			}
		case "Sprite":
			if s, ok := value.(Sprite); ok {
				if cs.spriteMap == nil {
					cs.spriteMap = make(map[string]Sprite)
				}
				cs.spriteMap[name] = s
			}
		case "TagList":
			if t, ok := value.(TagList); ok {
				if cs.tagListMap == nil {
					cs.tagListMap = make(map[string]TagList)
				}
				cs.tagListMap[name] = t
			}
		case "IntMap":
			if m, ok := value.(map[string]int); ok {
				if cs.intMapMap == nil {
					cs.intMapMap = make(map[string]IntMap)
				}
				cs.intMapMap[name] = NewIntMap(m)
			}
		case "FloatMap":
			if m, ok := value.(map[string]float64); ok {
				if cs.floatMapMap == nil {
					cs.floatMapMap = make(map[string]FloatMap)
				}
				cs.floatMapMap[name] = NewFloatMap(m)
			}
		case "Generic":
			if cs.genericMap == nil {
				cs.genericMap = make(map[string]any)
			}
			cs.genericMap[name] = value
		default:
			panic(fmt.Sprintf("unknown component kind %s", kind))
		}
	}
	return cs
}
