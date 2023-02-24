package sameriver

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/dt-rush/sameriver/v2/utils"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type ComponentTable struct {
	em *EntityManager

	// the size of the tables
	capacity int

	next_ix int
	ixs     map[string]int
	ixs_rev map[int]string
	names   map[string]bool
	kinds   map[string]string

	// data storage
	vec2DMap           map[string][]Vec2D
	boolMap            map[string][]bool
	intMap             map[string][]int
	float64Map         map[string][]float64
	timeMap            map[string][]time.Time
	timeAccumulatorMap map[string][]utils.TimeAccumulator
	stringMap          map[string][]string
	spriteMap          map[string][]Sprite
	tagListMap         map[string][]TagList
	intMapMap          map[string][]IntMap
	floatMapMap        map[string][]FloatMap
	genericMap         map[string][]interface{}
	cccMap             map[string]CustomContiguousComponent
}

func NewComponentTable(capacity int) *ComponentTable {
	return &ComponentTable{
		capacity: capacity,

		ixs:     make(map[string]int),
		ixs_rev: make(map[int]string),
		names:   make(map[string]bool),
		kinds:   make(map[string]string),

		vec2DMap:           make(map[string][]Vec2D),
		boolMap:            make(map[string][]bool),
		intMap:             make(map[string][]int),
		float64Map:         make(map[string][]float64),
		timeMap:            make(map[string][]time.Time),
		timeAccumulatorMap: make(map[string][]utils.TimeAccumulator),
		stringMap:          make(map[string][]string),
		spriteMap:          make(map[string][]Sprite),
		tagListMap:         make(map[string][]TagList),
		intMapMap:          make(map[string][]IntMap),
		floatMapMap:        make(map[string][]FloatMap),
		genericMap:         make(map[string][]interface{}),
		cccMap:             make(map[string]CustomContiguousComponent),
	}
}

// this is likely to be an expensive operation
func (ct *ComponentTable) expand(n int) {
	Logger.Printf("Expanding component tables from %d to %d", ct.capacity, ct.capacity+n)
	for name, slice := range ct.vec2DMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]Vec2D, n)
		ct.vec2DMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.boolMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]bool, n)
		ct.boolMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.intMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]int, n)
		ct.intMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.float64Map {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]float64, n)
		ct.float64Map[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.timeMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]time.Time, n)
		ct.timeMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.timeAccumulatorMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]utils.TimeAccumulator, n)
		ct.timeAccumulatorMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.stringMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]string, n)
		ct.stringMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.spriteMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]Sprite, n)
		ct.spriteMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.tagListMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]TagList, n)
		ct.tagListMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.intMapMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]IntMap, n)
		ct.intMapMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.floatMapMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]FloatMap, n)
		ct.floatMapMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.genericMap {
		Logger.Printf("Expanding table of component %s,%s", ct.kinds[name], name)
		extraSpace := make([]interface{}, n)
		ct.genericMap[name] = append(slice, extraSpace...)
	}
	for name, ccc := range ct.cccMap {
		Logger.Printf("Requesting expanding of internal storage of CustomContiguousComponent,%s", name)
		ccc.ExpandTable(n)
	}
	ct.capacity += n
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

func (ct *ComponentTable) addComponent(kind, name string) {
	ct.nameAndIndex(name)
	// create table in appropriate map
	// (note we allocate with capacity 2* so that if we reach max entities the
	// first time expanding the tables won't necessarily be expensive; but
	// then again, if we do reach the NEW capacity, the slices will have to
	// be reallocated to new memory locations as they'll have totally
	// eaten up the capacity)
	switch kind {
	case "Vec2D":
		ct.vec2DMap[name] = make([]Vec2D, ct.capacity, 2*ct.capacity)
	case "Bool":
		ct.boolMap[name] = make([]bool, ct.capacity, 2*ct.capacity)
	case "Int":
		ct.intMap[name] = make([]int, ct.capacity, 2*ct.capacity)
	case "Float64":
		ct.float64Map[name] = make([]float64, ct.capacity, 2*ct.capacity)
	case "Time":
		ct.timeMap[name] = make([]time.Time, ct.capacity, 2*ct.capacity)
	case "TimeAccumulator":
		ct.timeAccumulatorMap[name] = make([]utils.TimeAccumulator, ct.capacity, 2*ct.capacity)
	case "String":
		ct.stringMap[name] = make([]string, ct.capacity, 2*ct.capacity)
	case "Sprite":
		ct.spriteMap[name] = make([]Sprite, ct.capacity, 2*ct.capacity)
	case "TagList":
		ct.tagListMap[name] = make([]TagList, ct.capacity, 2*ct.capacity)
	case "IntMap":
		ct.intMapMap[name] = make([]IntMap, ct.capacity, 2*ct.capacity)
	case "FloatMap":
		ct.floatMapMap[name] = make([]FloatMap, ct.capacity, 2*ct.capacity)
	case "Generic":
		ct.genericMap[name] = make([]interface{}, ct.capacity, 2*ct.capacity)
	default:
		panic(fmt.Sprintf("added component of kind %s has no case in component_table.go", kind))
	}

	// note name
	ct.kinds[name] = kind
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
			panic(fmt.Sprintf("%s not found in vec2DMap", name))
		}
	}
	for name, _ := range cs.boolMap {
		if _, ok := ct.boolMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in boolMap", name))
		}
	}
	for name, _ := range cs.intMap {
		if _, ok := ct.intMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in intMap", name))
		}
	}
	for name, _ := range cs.float64Map {
		if _, ok := ct.float64Map[name]; !ok {
			panic(fmt.Sprintf("%s not found in float64Map", name))
		}
	}
	for name, _ := range cs.timeMap {
		if _, ok := ct.timeMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in timeMap", name))
		}
	}
	for name, _ := range cs.stringMap {
		if _, ok := ct.stringMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in stringMap", name))
		}
	}
	for name, _ := range cs.spriteMap {
		if _, ok := ct.spriteMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in spriteMap", name))
		}
	}
	for name, _ := range cs.tagListMap {
		if _, ok := ct.tagListMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in tagListMap", name))
		}
	}
	for name, _ := range cs.intMapMap {
		if _, ok := ct.intMapMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in intMapMap", name))
		}
	}
	for name, _ := range cs.floatMapMap {
		if _, ok := ct.floatMapMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in floatMapMap", name))
		}
	}
	for name, _ := range cs.genericMap {
		if _, ok := ct.genericMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in genericMap", name))
		}
	}
	for name, _ := range cs.customComponentsMap {
		if _, ok := ct.cccMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in cccMap", name))
		}
	}
}

func (ct *ComponentTable) ApplyComponentSet(e *Entity, spec map[string]any) {
	ct.applyComponentSet(e, makeComponentSet(spec))
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
	for name, t := range cs.timeMap {
		ct.timeMap[name][e.ID] = t
	}
	for name, t := range cs.timeAccumulatorMap {
		ct.timeAccumulatorMap[name][e.ID] = t
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

func (ct *ComponentTable) BitArrayFromComponentSet(spec map[string]any) bitarray.BitArray {
	return ct.bitArrayFromComponentSet(makeComponentSet(spec))
}

func (ct *ComponentTable) bitArrayFromComponentSet(cs ComponentSet) bitarray.BitArray {
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
func (e *Entity) GetTime(name string) *time.Time {
	return &e.World.em.components.timeMap[name][e.ID]
}
func (e *Entity) GetTimeAccumulator(name string) *utils.TimeAccumulator {
	return &e.World.em.components.timeAccumulatorMap[name][e.ID]
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
		return &e.World.em.components.intMapMap[name][e.ID]
	case "FloatMap":
		return &e.World.em.components.floatMapMap[name][e.ID]
	case "Generic":
		return e.World.em.components.genericMap[name][e.ID]
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
