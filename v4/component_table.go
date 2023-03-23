package sameriver

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type ComponentKind int
type ComponentID int

const (
	VEC2D ComponentKind = iota
	BOOL
	INT
	FLOAT64
	TIME
	TIMEACCUMULATOR
	STRING
	SPRITE
	TAGLIST
	INTMAP
	FLOATMAP
	GENERIC
	CUSTOM
)

var componentKindStrings = map[ComponentKind]string{
	VEC2D:           "VEC2D",
	BOOL:            "BOOL",
	INT:             "INT",
	FLOAT64:         "FLOAT64",
	TIME:            "TIME",
	TIMEACCUMULATOR: "TIMEACCUMULATOR",
	STRING:          "STRING",
	SPRITE:          "SPRITE",
	TAGLIST:         "TAGLIST",
	INTMAP:          "INTMAP",
	FLOATMAP:        "FLOATMAP",
	GENERIC:         "GENERIC",
	CUSTOM:          "CUSTOM",
}

type ComponentTable struct {
	// the size of the tables
	capacity int

	nextIx  int
	ixs     map[ComponentID]int
	ixsRev  map[int]ComponentID
	strings map[ComponentID]string
	kinds   map[ComponentID]ComponentKind

	// data storage
	vec2DMap           map[ComponentID][]Vec2D
	boolMap            map[ComponentID][]bool
	intMap             map[ComponentID][]int
	float64Map         map[ComponentID][]float64
	timeMap            map[ComponentID][]time.Time
	timeAccumulatorMap map[ComponentID][]TimeAccumulator
	stringMap          map[ComponentID][]string
	spriteMap          map[ComponentID][]Sprite
	tagListMap         map[ComponentID][]TagList
	intMapMap          map[ComponentID][]IntMap
	floatMapMap        map[ComponentID][]FloatMap
	genericMap         map[ComponentID][]any
	cccMap             map[ComponentID]CustomContiguousComponent
}

func NewComponentTable(capacity int) *ComponentTable {
	return &ComponentTable{
		capacity: capacity,

		ixs:     make(map[ComponentID]int),
		ixsRev:  make(map[int]ComponentID),
		strings: make(map[ComponentID]string),
		kinds:   make(map[ComponentID]ComponentKind),

		vec2DMap:           make(map[ComponentID][]Vec2D),
		boolMap:            make(map[ComponentID][]bool),
		intMap:             make(map[ComponentID][]int),
		float64Map:         make(map[ComponentID][]float64),
		timeMap:            make(map[ComponentID][]time.Time),
		timeAccumulatorMap: make(map[ComponentID][]TimeAccumulator),
		stringMap:          make(map[ComponentID][]string),
		spriteMap:          make(map[ComponentID][]Sprite),
		tagListMap:         make(map[ComponentID][]TagList),
		intMapMap:          make(map[ComponentID][]IntMap),
		floatMapMap:        make(map[ComponentID][]FloatMap),
		genericMap:         make(map[ComponentID][]any),
		cccMap:             make(map[ComponentID]CustomContiguousComponent),
	}
}

// this is likely to be an expensive operation
func (ct *ComponentTable) expand(n int) {
	Logger.Printf("Expanding component tables from %d to %d", ct.capacity, ct.capacity+n)
	for name, slice := range ct.vec2DMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]Vec2D, n)
		ct.vec2DMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.boolMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]bool, n)
		ct.boolMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.intMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]int, n)
		ct.intMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.float64Map {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]float64, n)
		ct.float64Map[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.timeMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]time.Time, n)
		ct.timeMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.timeAccumulatorMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]TimeAccumulator, n)
		ct.timeAccumulatorMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.stringMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]string, n)
		ct.stringMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.spriteMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]Sprite, n)
		ct.spriteMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.tagListMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]TagList, n)
		ct.tagListMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.intMapMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]IntMap, n)
		ct.intMapMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.floatMapMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]FloatMap, n)
		ct.floatMapMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.genericMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.kinds[name]], ct.strings[name])
		extraSpace := make([]any, n)
		ct.genericMap[name] = append(slice, extraSpace...)
	}
	for name, ccc := range ct.cccMap {
		Logger.Printf("Requesting expanding of internal storage of CustomContiguousComponent,%s", ct.strings[name])
		ccc.ExpandTable(n)
	}
	ct.capacity += n
}

func (ct *ComponentTable) RegisterComponentStrings(strings map[ComponentID]string) {
	for name, str := range strings {
		ct.strings[name] = str
	}
}

func (ct *ComponentTable) index(name ComponentID) {
	// increment index and store (used for bitarray generation)
	ct.ixs[name] = ct.nextIx
	ct.ixsRev[ct.nextIx] = name
	ct.nextIx++
}

func (ct *ComponentTable) ComponentExists(name ComponentID) bool {
	if _, ok := ct.ixs[name]; ok {
		return true
	}
	return false
}

func (ct *ComponentTable) addComponent(kind ComponentKind, name ComponentID, str string) {
	// create table in appropriate map
	// (note we allocate with capacity 2* so that if we reach max entities the
	// first time expanding the tables won't necessarily be expensive; but
	// then again, if we do reach the NEW capacity, the slices will have to
	// be reallocated to new memory locations as they'll have totally
	// eaten up the capacity)
	switch kind {
	case VEC2D:
		ct.vec2DMap[name] = make([]Vec2D, ct.capacity, 2*ct.capacity)
	case BOOL:
		ct.boolMap[name] = make([]bool, ct.capacity, 2*ct.capacity)
	case INT:
		ct.intMap[name] = make([]int, ct.capacity, 2*ct.capacity)
	case FLOAT64:
		ct.float64Map[name] = make([]float64, ct.capacity, 2*ct.capacity)
	case TIME:
		ct.timeMap[name] = make([]time.Time, ct.capacity, 2*ct.capacity)
	case TIMEACCUMULATOR:
		ct.timeAccumulatorMap[name] = make([]TimeAccumulator, ct.capacity, 2*ct.capacity)
	case STRING:
		ct.stringMap[name] = make([]string, ct.capacity, 2*ct.capacity)
	case SPRITE:
		ct.spriteMap[name] = make([]Sprite, ct.capacity, 2*ct.capacity)
	case TAGLIST:
		ct.tagListMap[name] = make([]TagList, ct.capacity, 2*ct.capacity)
	case INTMAP:
		ct.intMapMap[name] = make([]IntMap, ct.capacity, 2*ct.capacity)
	case FLOATMAP:
		ct.floatMapMap[name] = make([]FloatMap, ct.capacity, 2*ct.capacity)
	case GENERIC:
		ct.genericMap[name] = make([]any, ct.capacity, 2*ct.capacity)
	default:
		panic(fmt.Sprintf("added component of kind %s has no case in component_table.go", componentKindStrings[kind]))
	}

	// note name and kind
	ct.index(name)
	ct.kinds[name] = kind

	// note string
	ct.strings[name] = str
}

func (ct *ComponentTable) addCCC(name ComponentID, custom CustomContiguousComponent) {
	// guard against double insertion (many say it's a great time, but not here)
	if _, already := ct.ixs[name]; already {
		logWarning("trying to add CCC but component with id %d already exists. Skipping.", name)
		return
	}
	ct.cccMap[name] = custom
	ct.kinds[name] = CUSTOM
	ct.index(name)
	ct.strings[name] = custom.Name()
	custom.AllocateTable(MAX_ENTITIES)
}

func (ct *ComponentTable) AssertValidComponentSet(cs ComponentSet) {
	for name := range cs.vec2DMap {
		if _, ok := ct.vec2DMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in vec2DMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.boolMap {
		if _, ok := ct.boolMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in boolMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.intMap {
		if _, ok := ct.intMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in intMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.float64Map {
		if _, ok := ct.float64Map[name]; !ok {
			panic(fmt.Sprintf("%s not found in float64Map - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.timeMap {
		if _, ok := ct.timeMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in timeMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.stringMap {
		if _, ok := ct.stringMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in stringMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.spriteMap {
		if _, ok := ct.spriteMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in spriteMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.tagListMap {
		if _, ok := ct.tagListMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in tagListMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.intMapMap {
		if _, ok := ct.intMapMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in intMapMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.floatMapMap {
		if _, ok := ct.floatMapMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in floatMapMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.genericMap {
		if _, ok := ct.genericMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in genericMap - maybe not registered yet?", ct.strings[name]))
		}
	}
	for name := range cs.customComponentsMap {
		if _, ok := ct.cccMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in cccMap - maybe not registered yet?", ct.strings[name]))
		}
	}
}

func (ct *ComponentTable) ApplyComponentSet(e *Entity, spec map[ComponentID]any) {
	ct.applyComponentSet(e, ct.makeComponentSet(spec))
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

func (ct *ComponentTable) BitArrayFromIDs(IDs []ComponentID) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(len(ct.ixs)))
	for _, name := range IDs {
		b.SetBit(uint64(ct.ixs[name]))
	}
	return b
}

func (ct *ComponentTable) BitArrayFromComponentSet(spec map[ComponentID]any) bitarray.BitArray {
	return ct.bitArrayFromComponentSet(ct.makeComponentSet(spec))
}

func (ct *ComponentTable) bitArrayFromComponentSet(cs ComponentSet) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(len(ct.ixs)))
	for name := range cs.names {
		b.SetBit(uint64(ct.ixs[name]))
	}
	return b
}

// BitArrayToString prints a string representation of a component bitarray as a set with
// string representations of each component type whose bit is set
func (ct *ComponentTable) BitArrayToString(b bitarray.BitArray) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	names := make([]string, 0)
	for name, ix := range ct.ixs {
		bit, _ := b.GetBit(uint64(ix))
		// the index into the array is the component type int from the
		// iota const block in component_enum.go
		if bit {
			names = append(names, ct.strings[name])
		}
	}
	for i, name := range names {
		buf.WriteString(name)
		if i != len(names)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (ct *ComponentTable) guardInvalidComponentGet(e *Entity, name ComponentID) {
	var ix int
	var ok bool
	if ix, ok = ct.ixs[name]; !ok {
		msg := fmt.Sprintf("Tried to access %s component; but there is no component with that name", ct.strings[name])
		panic(msg)
	}
	bit, _ := e.ComponentBitArray.GetBit(uint64(ix))
	if !bit {
		msg := fmt.Sprintf("Tried to get %s component of entity without: %s", ct.strings[name], e)
		panic(msg)
	}
}

func (e *Entity) GetVec2D(name ComponentID) *Vec2D {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.vec2DMap[name][e.ID]
}
func (e *Entity) GetBool(name ComponentID) *bool {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.boolMap[name][e.ID]
}
func (e *Entity) GetInt(name ComponentID) *int {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.intMap[name][e.ID]
}
func (e *Entity) GetFloat64(name ComponentID) *float64 {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.float64Map[name][e.ID]
}
func (e *Entity) GetTime(name ComponentID) *time.Time {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.timeMap[name][e.ID]
}
func (e *Entity) GetTimeAccumulator(name ComponentID) *TimeAccumulator {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.timeAccumulatorMap[name][e.ID]
}
func (e *Entity) GetString(name ComponentID) *string {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.stringMap[name][e.ID]
}
func (e *Entity) GetSprite(name ComponentID) *Sprite {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.spriteMap[name][e.ID]
}
func (e *Entity) GetTagList(name ComponentID) *TagList {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.tagListMap[name][e.ID]
}
func (e *Entity) GetIntMap(name ComponentID) *IntMap {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.intMapMap[name][e.ID]
}
func (e *Entity) GetFloatMap(name ComponentID) *FloatMap {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.floatMapMap[name][e.ID]
}
func (e *Entity) GetGeneric(name ComponentID) any {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return e.World.em.components.genericMap[name][e.ID]
}
func (e *Entity) SetGeneric(name ComponentID, val any) {
	e.World.em.components.guardInvalidComponentGet(e, name)
	e.World.em.components.genericMap[name][e.ID] = val
}
func (e *Entity) GetVal(name ComponentID) any {
	e.World.em.components.guardInvalidComponentGet(e, name)
	kind := e.World.em.components.kinds[name]
	switch kind {
	case VEC2D:
		return &e.World.em.components.vec2DMap[name][e.ID]
	case BOOL:
		return &e.World.em.components.boolMap[name][e.ID]
	case INT:
		return &e.World.em.components.intMap[name][e.ID]
	case FLOAT64:
		return &e.World.em.components.float64Map[name][e.ID]
	case STRING:
		return &e.World.em.components.stringMap[name][e.ID]
	case SPRITE:
		return &e.World.em.components.spriteMap[name][e.ID]
	case TAGLIST:
		return &e.World.em.components.tagListMap[name][e.ID]
	case INTMAP:
		return &e.World.em.components.intMapMap[name][e.ID]
	case FLOATMAP:
		return &e.World.em.components.floatMapMap[name][e.ID]
	case GENERIC:
		return e.World.em.components.genericMap[name][e.ID]
	case CUSTOM:
		return e.World.em.components.cccMap[name].Get(e)
	default:
		panic(fmt.Sprintf("Can't get component with ID %d - it doesn't seem to exist", name))
	}
}

// GetCustom returns the custom component data for the entity
// NOTE: we have to provide a get and set method since we can't
// return a pointer to any
func (e *Entity) GetCustom(name ComponentID) any {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return e.World.em.components.cccMap[name].Get(e)
}
func (e *Entity) SetCustom(name ComponentID, x any) {
	e.World.em.components.guardInvalidComponentGet(e, name)
	e.World.em.components.cccMap[name].Set(e, x)
}
