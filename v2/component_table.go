package sameriver

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/golang-collections/go-datastructures/bitarray"
)

var COMPONENT_KINDS = []string{
	"Vec2D",
	"Bool",
	"Int",
	"Float64",
	"String",
	"Sprite",
	"TagList",
	"IntMap",
	"FloatMap",
	"Generic",
	"Custom",
}

type ComponentTable struct {
	em *EntityManager

	next_ix int
	ixs     map[string]int
	ixs_rev map[int]string
	names   map[string]bool
	kinds   map[string]string

	// data storage
	vec2DMap    map[string][]Vec2D
	boolMap     map[string][]bool
	intMap      map[string][]int
	float64Map  map[string][]float64
	stringMap   map[string][]string
	spriteMap   map[string][]Sprite
	tagListMap  map[string][]TagList
	intMapMap   map[string][]IntMap
	floatMapMap map[string][]FloatMap
	genericMap  map[string][]interface{}
	cccMap      map[string]CustomContiguousComponent
}

func NewComponentTable() *ComponentTable {
	ct := &ComponentTable{}
	ct.ixs = make(map[string]int)
	ct.ixs_rev = make(map[int]string)
	ct.names = make(map[string]bool)
	ct.kinds = make(map[string]string)

	ct.vec2DMap = make(map[string][]Vec2D)
	ct.boolMap = make(map[string][]bool)
	ct.intMap = make(map[string][]int)
	ct.float64Map = make(map[string][]float64)
	ct.stringMap = make(map[string][]string)
	ct.spriteMap = make(map[string][]Sprite)
	ct.tagListMap = make(map[string][]TagList)
	ct.intMapMap = make(map[string][]IntMap)
	ct.floatMapMap = make(map[string][]FloatMap)
	ct.genericMap = make(map[string][]interface{})
	ct.cccMap = make(map[string]CustomContiguousComponent)
	return ct
}

func (ct *ComponentTable) nameAndIndex(name string) bool {
	if _, ok := ct.names[name]; ok {
		return true
	} else {
		ct.names[name] = true
	}
	// increment index and store (used for bitarray generation)
	ct.ixs[name] = ct.next_ix
	ct.ixs_rev[ct.next_ix] = name
	ct.next_ix++
	return false
}

func (ct *ComponentTable) ComponentExists(spec string) bool {
	// decode spec string
	split := strings.Split(spec, ",")
	kind := split[0]
	name := split[1]
	if k, ok := ct.kinds[name]; ok {
		return kind == k
	}
	return false
}

func (ct *ComponentTable) AddComponent(spec string) {
	// decode spec string
	split := strings.Split(spec, ",")
	kind := split[0]
	name := split[1]

	// guard against double insertion (many say it's a great time, but not here)
	if already := ct.nameAndIndex(name); already {
		Logger.Println(fmt.Sprintf("Warning: trying to add component but component with name %s already exists. Skipping.", name))
		return
	}

	// create table in appropriate map
	switch kind {
	case "Vec2D":
		ct.vec2DMap[name] = make([]Vec2D, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "Vec2D"
	case "Bool":
		ct.boolMap[name] = make([]bool, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "Bool"
	case "Int":
		ct.intMap[name] = make([]int, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "Int"
	case "Float64":
		ct.float64Map[name] = make([]float64, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "Float64"
	case "String":
		ct.stringMap[name] = make([]string, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "String"
	case "Sprite":
		ct.spriteMap[name] = make([]Sprite, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "Sprite"
	case "TagList":
		ct.tagListMap[name] = make([]TagList, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "TagList"
	case "IntMap":
		ct.intMapMap[name] = make([]IntMap, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "IntMap"
	case "FloatMap":
		ct.floatMapMap[name] = make([]FloatMap, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "FloatMap"
	case "Generic":
		ct.genericMap[name] = make([]interface{}, MAX_ENTITIES, MAX_ENTITIES)
		ct.kinds[name] = "Generic"
	default:
		panic(fmt.Sprintf("added component of kind %s has no case in component_table.go", kind))
	}
}

func (ct *ComponentTable) AddCCC(custom CustomContiguousComponent) {
	// guard against double insertion (many say it's a great time, but not here)
	if already := ct.nameAndIndex(custom.Name()); already {
		Logger.Println(fmt.Sprintf("Warning: trying to add component but component with name %s already exists. Skipping.", custom.Name()))
		return
	}
	ct.cccMap[custom.Name()] = custom
	ct.kinds[custom.Name()] = "Custom"
	custom.AllocateTable(MAX_ENTITIES)
}

func (ct *ComponentTable) AssertValidComponentSet(cs ComponentSet) {
	for name, _ := range cs.vec2DMap {
		if _, ok := ct.vec2DMap[name]; !ok {
			panic("%s not found in vec2DMap")
		}
	}
	for name, _ := range cs.boolMap {
		if _, ok := ct.boolMap[name]; !ok {
			panic("%s not found in boolMap")
		}
	}
	for name, _ := range cs.intMap {
		if _, ok := ct.intMap[name]; !ok {
			panic("%s not found in intMap")
		}
	}
	for name, _ := range cs.float64Map {
		if _, ok := ct.float64Map[name]; !ok {
			panic("%s not found in float64Map")
		}
	}
	for name, _ := range cs.stringMap {
		if _, ok := ct.stringMap[name]; !ok {
			panic("%s not found in stringMap")
		}
	}
	for name, _ := range cs.spriteMap {
		if _, ok := ct.spriteMap[name]; !ok {
			panic("%s not found in spriteMap")
		}
	}
	for name, _ := range cs.tagListMap {
		if _, ok := ct.tagListMap[name]; !ok {
			panic("%s not found in tagListMap")
		}
	}
	for name, _ := range cs.intMapMap {
		if _, ok := ct.intMapMap[name]; !ok {
			panic("%s not found in intMapMap")
		}
	}
	for name, _ := range cs.floatMapMap {
		if _, ok := ct.floatMapMap[name]; !ok {
			panic("%s not found in floatMapMap")
		}
	}
	for name, _ := range cs.genericMap {
		if _, ok := ct.genericMap[name]; !ok {
			panic("%s not found in genericMap")
		}
	}
	for name, _ := range cs.customComponentsMap {
		if _, ok := ct.cccMap[name]; !ok {
			panic("%s not found in cccMap")
		}
	}
}

func (ct *ComponentTable) applyComponentSet(e *Entity, cs ComponentSet) {
	ct.AssertValidComponentSet(cs)
	for name, v := range cs.vec2DMap {
		ct.vec2DMap[name][e.ID] = v
	}
	for name, b := range cs.boolMap {
		ct.boolMap[name][e.ID] = b
	}
	for name, i := range cs.intMap {
		ct.intMap[name][e.ID] = i
	}
	for name, f := range cs.float64Map {
		ct.float64Map[name][e.ID] = f
	}
	for name, s := range cs.stringMap {
		ct.stringMap[name][e.ID] = s
	}
	for name, s := range cs.spriteMap {
		ct.spriteMap[name][e.ID] = s
	}
	for name, t := range cs.tagListMap {
		ct.tagListMap[name][e.ID] = t
	}
	for name, m := range cs.intMapMap {
		ct.intMapMap[name][e.ID] = m
	}
	for name, m := range cs.floatMapMap {
		ct.floatMapMap[name][e.ID] = m
	}
	for name, x := range cs.genericMap {
		ct.genericMap[name][e.ID] = x
	}
	for name, x := range cs.customComponentsMap {
		cs.customComponentsImpl[name].Set(e, x)
	}
}

func (ct *ComponentTable) BitArrayFromNames(names []string) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(len(ct.ixs)))
	for _, name := range names {
		b.SetBit(uint64(ct.ixs[name]))
	}
	return b
}

func (ct *ComponentTable) BitArrayFromComponentSet(cs ComponentSet) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(len(ct.ixs)))
	for name, _ := range cs.names {
		b.SetBit(uint64(ct.ixs[name]))
	}
	return b
}

// prints a string representation of a component bitarray as a set with
// string representations of each component type whose bit is set
func (ct *ComponentTable) BitArrayToString(b bitarray.BitArray) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i := uint64(0); i < uint64(len(ct.names)); i++ {
		bit, _ := b.GetBit(i)
		// the index into the array is the component type int from the
		// iota const block in component_enum.go
		if bit {
			buf.WriteString(fmt.Sprintf("%s", ct.ixs_rev[int(i)]))
			if i != uint64(len(ct.names)-1) {
				buf.WriteString(", ")
			}
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (e *Entity) GetVec2D(name string) *Vec2D {
	return &e.World.em.components.vec2DMap[name][e.ID]
}
func (e *Entity) GetBool(name string) *bool {
	return &e.World.em.components.boolMap[name][e.ID]
}
func (e *Entity) GetInt(name string) *int {
	return &e.World.em.components.intMap[name][e.ID]
}
func (e *Entity) GetFloat64(name string) *float64 {
	return &e.World.em.components.float64Map[name][e.ID]
}
func (e *Entity) GetString(name string) *string {
	return &e.World.em.components.stringMap[name][e.ID]
}
func (e *Entity) GetSprite(name string) *Sprite {
	return &e.World.em.components.spriteMap[name][e.ID]
}
func (e *Entity) GetTagList(name string) *TagList {
	return &e.World.em.components.tagListMap[name][e.ID]
}
func (e *Entity) GetIntMap(name string) *IntMap {
	return &e.World.em.components.intMapMap[name][e.ID]
}
func (e *Entity) GetFloatMap(name string) *FloatMap {
	return &e.World.em.components.floatMapMap[name][e.ID]
}
func (e *Entity) GetGeneric(name string) interface{} {
	return e.World.em.components.genericMap[name][e.ID]
}
func (e *Entity) SetGeneric(name string, val interface{}) {
	e.World.em.components.genericMap[name][e.ID] = val
}
func (e *Entity) GetVal(name string) interface{} {
	kind := e.World.em.components.kinds[name]
	switch kind {
	case "Vec2D":
		return &e.World.em.components.vec2DMap[name][e.ID]
	case "Bool":
		return &e.World.em.components.boolMap[name][e.ID]
	case "Int":
		return &e.World.em.components.intMap[name][e.ID]
	case "Float64":
		return &e.World.em.components.float64Map[name][e.ID]
	case "String":
		return &e.World.em.components.stringMap[name][e.ID]
	case "Sprite":
		return &e.World.em.components.spriteMap[name][e.ID]
	case "TagList":
		return &e.World.em.components.tagListMap[name][e.ID]
	case "IntMap":
		return e.World.em.components.intMapMap[name][e.ID]
	case "FloatMap":
		return e.World.em.components.floatMapMap[name][e.ID]
	case "Generic":
		return &e.World.em.components.genericMap[name][e.ID]
	case "Custom":
		return e.World.em.components.cccMap[name].Get(e)
	default:
		panic(fmt.Sprintf("Can't get component %s - it doesn't seem to exist", name))
	}
}

// NOTE: we have to provide a get and set method since we can't
// return a pointer to interface{}
func (e *Entity) GetCustom(name string) interface{} {
	return e.World.em.components.cccMap[name].Get(e)
}
func (e *Entity) SetCustom(name string, x interface{}) {
	e.World.em.components.cccMap[name].Set(e, x)
}
