package sameriver

import (
	"fmt"
	"time"
)

type ComponentSet struct {
	// names (ComponentID) of all components given values in this set
	names map[ComponentID]bool
	// data storage
	vec2DMap             map[ComponentID]Vec2D
	boolMap              map[ComponentID]bool
	intMap               map[ComponentID]int
	float64Map           map[ComponentID]float64
	timeMap              map[ComponentID]time.Time
	timeAccumulatorMap   map[ComponentID]TimeAccumulator
	stringMap            map[ComponentID]string
	spriteMap            map[ComponentID]Sprite
	tagListMap           map[ComponentID]TagList
	intMapMap            map[ComponentID]IntMap
	floatMapMap          map[ComponentID]FloatMap
	genericMap           map[ComponentID]any
	customComponentsMap  map[ComponentID]any
	customComponentsImpl map[ComponentID]CustomContiguousComponent
}

func (ct *ComponentTable) makeCustomComponentSet(
	componentSpecs map[ComponentID]any,
	customComponentSpecs map[ComponentID]any,
	customComponentsImpl map[ComponentID]CustomContiguousComponent) ComponentSet {

	baseCS := ct.makeComponentSet(componentSpecs)
	baseCS.customComponentsMap = make(map[ComponentID]any)
	baseCS.customComponentsImpl = make(map[ComponentID]CustomContiguousComponent)
	for name, value := range customComponentSpecs {
		kind := ct.kinds[name]
		if kind != CUSTOM {
			panic(fmt.Sprintf("custom component spec should have type Custom, it appears to be: %s", kindStrings[kind]))
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
func (ct *ComponentTable) makeComponentSet(componentSpecs map[ComponentID]any) ComponentSet {
	cs := ComponentSet{
		names: make(map[ComponentID]bool),
	}
	for name, value := range componentSpecs {
		kind := ct.kinds[name]
		// take note in names map that this component name occurs
		cs.names[name] = true
		// assign values into appropriate maps
		switch kind {
		case VEC2D:
			if v, ok := value.(Vec2D); ok {
				if cs.vec2DMap == nil {
					cs.vec2DMap = make(map[ComponentID]Vec2D)
				}
				cs.vec2DMap[name] = v
			}
		case BOOL:
			if b, ok := value.(bool); ok {
				if cs.boolMap == nil {
					cs.boolMap = make(map[ComponentID]bool)
				}
				cs.boolMap[name] = b
			}
		case INT:
			if i, ok := value.(int); ok {
				if cs.intMap == nil {
					cs.intMap = make(map[ComponentID]int)
				}
				cs.intMap[name] = i
			}
		case FLOAT64:
			if f, ok := value.(float64); ok {
				if cs.float64Map == nil {
					cs.float64Map = make(map[ComponentID]float64)
				}
				cs.float64Map[name] = f
			}
		case TIME:
			if t, ok := value.(time.Time); ok {
				if cs.timeMap == nil {
					cs.timeMap = make(map[ComponentID]time.Time)
				}
				cs.timeMap[name] = t
			}
		case TIMEACCUMULATOR:
			if t, ok := value.(TimeAccumulator); ok {
				if cs.timeAccumulatorMap == nil {
					cs.timeAccumulatorMap = make(map[ComponentID]TimeAccumulator)
				}
				cs.timeAccumulatorMap[name] = t
			}
		case STRING:
			if s, ok := value.(string); ok {
				if cs.stringMap == nil {
					cs.stringMap = make(map[ComponentID]string)
				}
				cs.stringMap[name] = s
			}
		case SPRITE:
			if s, ok := value.(Sprite); ok {
				if cs.spriteMap == nil {
					cs.spriteMap = make(map[ComponentID]Sprite)
				}
				cs.spriteMap[name] = s
			}
		case TAGLIST:
			if t, ok := value.(TagList); ok {
				if cs.tagListMap == nil {
					cs.tagListMap = make(map[ComponentID]TagList)
				}
				cs.tagListMap[name] = t
			}
		case INTMAP:
			if m, ok := value.(map[string]int); ok {
				if cs.intMapMap == nil {
					cs.intMapMap = make(map[ComponentID]IntMap)
				}
				cs.intMapMap[name] = NewIntMap(m)
			}
		case FLOATMAP:
			if m, ok := value.(map[string]float64); ok {
				if cs.floatMapMap == nil {
					cs.floatMapMap = make(map[ComponentID]FloatMap)
				}
				cs.floatMapMap[name] = NewFloatMap(m)
			}
		case GENERIC:
			if cs.genericMap == nil {
				cs.genericMap = make(map[ComponentID]any)
			}
			cs.genericMap[name] = value
		default:
			panic(fmt.Sprintf("unknown component kind %s", kind))
		}
	}
	return cs
}
