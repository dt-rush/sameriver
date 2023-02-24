package sameriver

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"github.com/TwiN/go-color"
)

type World struct {

	// rand.Seed for this world's run
	Seed int

	Width  float64
	Height float64

	IdGen  *IDGenerator
	Events *EventBus
	em     *EntityManager

	// systems registered
	systems map[string]System
	// this is needed to associate ID's with Systems, since System is an
	// interface, not a struct type like LogicUnit that can have a name field
	systemsIDs map[System]int

	// logics invoked regularly by runtimeSharer
	worldLogics map[string]*LogicUnit

	// funcs that can be called by name with data and get a result,
	// or to produce an effect
	funcs *FuncSet

	// blackboards that entity's can join to share events and state
	blackboards map[string]*Blackboard

	// for sharing runtime among the various runtimelimiter kinds
	// and contains the RuntimeLimiters to which we Add() LogicUnits
	runtimeSharer *RuntimeLimitSharer

	// for statistics tracking - the avg ms used to run World.Update()
	totalRuntimeAvg_ms *float64

	// used for entity distance queries
	SpatialHasher *SpatialHasher
}

type WorldSpec struct {
	Width               int
	Height              int
	DistanceHasherGridX int
	DistanceHasherGridY int
}

func destructureWorldSpec(spec map[string]any) WorldSpec {
	var width, height int
	var distanceHasherGridX, distanceHasherGridY int
	if _, ok := spec["width"].(int); ok {
		width = spec["width"].(int)
	} else {
		width = 100
	}
	if _, ok := spec["height"].(int); ok {
		height = spec["height"].(int)
	} else {
		height = 100
	}
	if _, ok := spec["distanceHasherGridX"].(int); ok {
		distanceHasherGridX = spec["distanceHasherGridX"].(int)
	} else {
		distanceHasherGridX = 10
	}
	if _, ok := spec["distanceHasherGridY"].(int); ok {
		distanceHasherGridY = spec["distanceHasherGridY"].(int)
	} else {
		distanceHasherGridY = 10
	}

	return WorldSpec{
		Width:               width,
		Height:              height,
		DistanceHasherGridX: distanceHasherGridX,
		DistanceHasherGridY: distanceHasherGridY,
	}
}

func NewWorld(spec map[string]any) *World {
	// seed a random number from [1,108]
	seed := rand.Intn(108) + 1
	rand.Seed(int64(seed))
	Logger.Println(color.InBold(color.InWhiteOverCyan(fmt.Sprintf("[world seed: %d]", seed))))
	destructured := destructureWorldSpec(spec)
	w := &World{
		Seed:          seed,
		Width:         float64(destructured.Width),
		Height:        float64(destructured.Height),
		Events:        NewEventBus(),
		IdGen:         NewIDGenerator(),
		systems:       make(map[string]System),
		systemsIDs:    make(map[System]int),
		worldLogics:   make(map[string]*LogicUnit),
		funcs:         NewFuncSet(nil),
		blackboards:   make(map[string]*Blackboard),
		runtimeSharer: NewRuntimeLimitSharer(),
	}

	// set up runtimesharer
	w.runtimeSharer.RegisterRunner("entity-manager")
	w.runtimeSharer.RegisterRunner("systems")
	w.runtimeSharer.RegisterRunner("world")
	w.runtimeSharer.RegisterRunner("entities")
	w.runtimeSharer.RegisterRunner("distance-query-spatial-hasher")
	// init entitymanager
	w.em = NewEntityManager(w)
	w.runtimeSharer.AddLogic("entity-manager",
		&LogicUnit{
			name:    "entity-manager",
			worldID: w.IdGen.Next(),
			f: func(dt_ms float64) {
				w.em.Update(FRAME_DURATION_INT / 2)
			},
			active:      true,
			runSchedule: nil,
		})
	// set up distance spatial hasher
	w.SpatialHasher = NewSpatialHasher(
		destructured.DistanceHasherGridX,
		destructured.DistanceHasherGridY,
		w,
	)
	w.runtimeSharer.AddLogic("distance-query-spatial-hasher",
		&LogicUnit{
			name:    "distance-query-spatial-hasher",
			worldID: w.IdGen.Next(),
			f: func(dt_ms float64) {
				w.SpatialHasher.ClearTable()
				w.SpatialHasher.ScanAndInsertEntities()
			},
			active:      true,
			runSchedule: nil,
		})
	// register basic components
	w.RegisterComponents(
		"TagList,GenericTags",
		"Vec2D,Position",
		"Vec2D,Box",
	)
	return w
}

func (w *World) Update(allowance_ms float64) (overunder_ms float64) {
	t0 := time.Now()
	w.em.Update(FRAME_DURATION_INT / 2)
	overunder_ms, starved := w.runtimeSharer.Share(allowance_ms)
	if starved > 0 {
		Logger.Println("Starvation of RuntimeLimiters occuring in World.Update(); Logic Units will be getting run less frequently.")
	}
	// maintain total runtime moving average
	total := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	if w.totalRuntimeAvg_ms == nil {
		w.totalRuntimeAvg_ms = &total
	} else {
		*w.totalRuntimeAvg_ms = (*w.totalRuntimeAvg_ms + total) / 2.0
	}
	return overunder_ms
}

func (w *World) RegisterComponents(specs ...string) {
	// register given specs
	for _, spec := range specs {
		// guard against double insertion (many say it's a great time, but not here)
		split := strings.Split(spec, ",")
		kind := split[0]
		name := split[1]
		if w.em.components.ComponentExists(spec) {
			Logger.Println(fmt.Sprintf("[%s,%s already exists. Skipping...]", kind, name))
			continue
		} else {
			Logger.Printf("%s%s%s", color.InGreen("[registering component: "), fmt.Sprintf("%s,%s", color.InBlue(kind), name), color.InGreen("]"))
			w.em.components.addComponent(kind, name)
		}
	}
}

func (w *World) RegisterCCCs(customs []CustomContiguousComponent) {
	// register custom contiguous components
	for _, custom := range customs {
		w.em.components.AddCCC(custom)
	}
}

func (w *World) RegisterSystems(systems ...System) {
	// add all systems
	for _, s := range systems {
		for _, spec := range s.GetComponentDeps() {
			if !w.em.components.ComponentExists(spec) {
				w.RegisterComponents(spec)
			}
		}
		w.addSystem(s)
	}
	// link up all systems' dependencies
	for _, s := range systems {
		w.linkSystemDependencies(s)
	}
}

func (w *World) SetSystemSchedule(systemName string, period_ms float64) {
	system := w.systems[systemName]
	systemLogicWorldID := w.systemsIDs[system]
	w.runtimeSharer.SetSchedule("systems", systemLogicWorldID, period_ms)
}

func (w *World) addSystem(s System) {
	w.assertSystemValid(s)
	name := reflect.TypeOf(s).Elem().Name()
	if _, ok := w.systems[name]; ok {
		panic(fmt.Sprintf("double-add of system %s", name))
	}
	w.systems[name] = s
	ID := w.IdGen.Next()
	w.systemsIDs[s] = ID
	s.LinkWorld(w)
	logicName := fmt.Sprintf("%s", name)
	// add logic immediately rather than wait for runtimeSharer.Share() to
	// process the add/remove logic channel so that if we call SetSystemSchedule()
	// immediately after RegisterSystems(), the LogicUnit will be in the runner
	// to set the runSchedule on
	w.runtimeSharer.addLogicImmediately("systems",
		&LogicUnit{
			name:        logicName,
			worldID:     w.systemsIDs[s],
			f:           s.Update,
			active:      true,
			runSchedule: nil,
		})
}

func (w *World) assertSystemValid(s System) {
	t := reflect.TypeOf(s)
	typeName := t.Elem().Name()
	if _, ok := s.(System); !ok {
		panic(fmt.Sprintf("Can't add object of type %s - doesn't implement System interface", typeName))
	}
	w.assertSystemTypeValid(t)
}

func (w *World) assertSystemTypeValid(t reflect.Type) {
	if t.Kind() != reflect.Ptr {
		panic("Implementers of engine.System must be pointer-receivers")
	}
	typeName := t.Elem().Name()
	validName, _ := regexp.MatchString(".+System$", typeName)
	if !validName {
		panic(fmt.Sprintf("Implementers of System must have a name "+
			"matching regexp .+System$. %s did not", typeName))
	}
}

func (w *World) linkSystemDependencies(s System) {
	// foreach field of the underlying struct,
	// check if it has the tag `sameriver-system-dependency`
	// if it does, search for the system with the same type as that
	// field and assign it as a pointer, cast to the expected type,
	// to that field
	//
	// sType is going to be something like *CollisionSystem
	sType := reflect.TypeOf(s).Elem()
	// get a type to represent the System interface (to ensure dependencies
	// are to implementers of System)
	systemInterface := reflect.TypeOf((*System)(nil)).Elem()
	for i := 0; i < sType.NumField(); i++ {
		// for each field of the struct
		// f would be something like sh *SpatialHashSystem, possibly with a tag
		f := sType.Field(i)
		tagVal := f.Tag.Get("sameriver-system-dependency")
		if tagVal != "" {
			// check that tagged field implements System and is a valid System
			// implemented
			isSystem := f.Type.Implements(systemInterface)
			if !isSystem {
				panic(fmt.Sprintf("fields tagged sameriver-system-dependency "+
					"must implement engine.System "+
					"(field %s %v of %s did not pass this requirement",
					f.Name, f.Type, sType.Name()))
			}
			w.assertSystemTypeValid(f.Type)
			// iterate through the other systems and find one whose type matches
			// the field's type
			var foundSystem System
			for _, otherSystem := range w.systems {
				if otherSystem == s {
					continue
				}
				if reflect.TypeOf(otherSystem) == f.Type {
					foundSystem = otherSystem
					break
				}
			}
			if foundSystem == nil {
				if tagVal == "optional" {
					continue
				} else {
					panic(fmt.Sprintf("%s %v of %s dependency could not be "+
						"resolved. No system found of type %v.",
						f.Name, f.Type, sType.Name(), f.Type))
				}
			}
			// now that we have found the system which corresponds to the
			// dependency, we will assign it to the place it should be
			//
			// thank you to feilengcui008 from golang-nuts for this method of
			// assigning to an unexported pointer field whose value is nil
			//
			// since vf is nil value, vf.Elem() will be the zero value, and
			// since the zero value is not addressable or settable, we
			// need to allocate a new settable value at the same address
			v := reflect.Indirect(reflect.ValueOf(s))
			vf := v.Field(i)
			vf = reflect.NewAt(vf.Type(), unsafe.Pointer(vf.UnsafeAddr())).Elem()
			vf.Set(reflect.ValueOf(foundSystem))
		}
	}
}

func (w *World) AddWorldLogic(Name string, F func(dt_ms float64)) *LogicUnit {
	if _, ok := w.worldLogics[Name]; ok {
		panic(fmt.Sprintf("double-add of world logic %s", Name))
	}
	l := &LogicUnit{
		name:        Name,
		f:           F,
		active:      true,
		worldID:     w.IdGen.Next(),
		runSchedule: nil,
	}
	w.worldLogics[Name] = l
	w.runtimeSharer.AddLogic("world", l)
	return l
}

func (w *World) AddWorldLogicWithSchedule(Name string, F func(dt_ms float64), period_ms float64) *LogicUnit {
	l := w.AddWorldLogic(Name, F)
	runSchedule := NewTimeAccumulator(period_ms)
	l.runSchedule = &runSchedule
	return l
}

func (w *World) RemoveWorldLogic(Name string) {
	if logic, ok := w.worldLogics[Name]; ok {
		w.runtimeSharer.RemoveLogic("world", logic)
		delete(w.worldLogics, Name)
		w.IdGen.Free(logic.worldID)
	}
}

func (w *World) ActivateAllWorldLogics() {
	w.runtimeSharer.ActivateAll("world")
}

func (w *World) DeactivateAllWorldLogics() {
	w.runtimeSharer.DeactivateAll("world")
}

func (w *World) ActivateWorldLogic(name string) {
	if logic, ok := w.worldLogics[name]; ok {
		logic.Activate()
	}
}

func (w *World) DeactivateWorldLogic(name string) {
	if logic, ok := w.worldLogics[name]; ok {
		logic.Deactivate()
	}
}

func (w *World) addEntityLogic(e *Entity, l *LogicUnit) *LogicUnit {
	w.runtimeSharer.AddLogic("entities", l)
	return l
}

func (w *World) removeEntityLogic(e *Entity, l *LogicUnit) {
	w.runtimeSharer.RemoveLogic("entities", l)
	w.IdGen.Free(l.worldID)
}

func (w *World) RemoveAllEntityLogics(e *Entity) {
	for _, l := range e.Logics {
		w.runtimeSharer.RemoveLogic("entities", l)
		w.IdGen.Free(l.worldID)
	}
}

func (w *World) ActivateAllEntityLogics() {
	w.runtimeSharer.ActivateAll("entities")
}

func (w *World) DeactivateAllEntityLogics() {
	w.runtimeSharer.DeactivateAll("entities")
}

func (w *World) ActivateEntityLogics(e *Entity) {
	for _, logic := range e.Logics {
		logic.Activate()
	}
}

func (w *World) DeactivateEntityLogics(e *Entity) {
	for _, logic := range e.Logics {
		logic.Deactivate()
	}
}

func (w *World) AddFuncs(funcs map[string](func(interface{}) interface{})) {
	for name, f := range funcs {
		w.funcs.Add(name, f)
	}
}

func (w *World) AddFunc(name string, f func(interface{}) interface{}) {
	w.funcs.Add(name, f)
}

func (w *World) RemoveFunc(name string) {
	w.funcs.Remove(name)
}

func (w *World) GetFunc(name string) func(interface{}) interface{} {
	return w.funcs.funcs[name]
}

func (w *World) HasFunc(name string) bool {
	return w.funcs.Has(name)
}

func (w *World) Blackboard(name string) *Blackboard {
	if _, ok := w.blackboards[name]; !ok {
		w.blackboards[name] = NewBlackboard(name)
	}
	return w.blackboards[name]
}

func (w *World) String() string {
	// TODO: implement
	return "TODO"
}

func (w *World) DumpStats() map[string](map[string]float64) {
	stats := w.runtimeSharer.DumpStats()
	// add total Update() runtime avg
	if w.totalRuntimeAvg_ms != nil {
		stats["totals"]["total"] = *w.totalRuntimeAvg_ms
	} else {
		stats["totals"]["total"] = 0.0
	}
	return stats
}

func (w *World) DumpStatsString() string {
	stats := w.DumpStats()
	b, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
