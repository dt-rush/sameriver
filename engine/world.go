package engine

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"time"
	"unsafe"

	"github.com/dt-rush/sameriver/engine/utils"
)

type World struct {
	Width  float64
	Height float64

	IdGen  *utils.IDGenerator
	Events *EventBus
	em     *EntityManager

	systems       map[string]System
	systemsRunner *RuntimeLimiter
	// this is needed to associate ID's with Systems, since System is an
	// interface, not a struct type like LogicUnit that can have a name field
	systemsIDs map[System]int

	worldLogics       map[string]*LogicUnit
	worldLogicsRunner *RuntimeLimiter

	// for primary entity logics
	primaryEntityLogics       map[int]*LogicUnit
	primaryEntityLogicsRunner *RuntimeLimiter

	// for additional specialised entity logics
	logicUnitComponentRunner *RuntimeLimiter

	totalRuntime *float64
}

func NewWorld(width int, height int) *World {
	w := &World{
		Width:                     float64(width),
		Height:                    float64(height),
		Events:                    NewEventBus(),
		IdGen:                     utils.NewIDGenerator(),
		systems:                   make(map[string]System),
		systemsIDs:                make(map[System]int),
		systemsRunner:             NewRuntimeLimiter(),
		worldLogics:               make(map[string]*LogicUnit),
		worldLogicsRunner:         NewRuntimeLimiter(),
		primaryEntityLogics:       make(map[int]*LogicUnit),
		primaryEntityLogicsRunner: NewRuntimeLimiter(),
		logicUnitComponentRunner:  NewRuntimeLimiter(),
	}
	// init entitymanager
	w.em = NewEntityManager(w)
	// register generic taglist
	w.em.components.AddComponent("TagList,GenericTags")
	// register generic logic
	w.em.components.AddComponent("*LogicUnit,GenericLogic")
	return w
}

func (w *World) Update(allowance float64) (overrun_ms float64) {
	t0 := time.Now()
	w.em.Update(FRAME_SLEEP_MS / 2)
	// systems update functions, world logic, and entity logic can use
	// whatever time is left over after entity manager update
	overunder := RuntimeLimitShare(
		allowance-float64(time.Since(t0).Nanoseconds())/1.0e6,
		w.systemsRunner,
		w.worldLogicsRunner,
		w.primaryEntityLogicsRunner,
		w.logicUnitComponentRunner)
	total := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	if w.totalRuntime == nil {
		w.totalRuntime = &total
	} else {
		*w.totalRuntime = (*w.totalRuntime + total) / 2.0
	}
	return overunder
}

func (w *World) RegisterComponents(specs []string) {
	// register given specs
	for _, spec := range specs {
		w.em.components.AddComponent(spec)
	}
}

func (w *World) AddSystems(systems ...System) {
	// add all systems
	for _, s := range systems {
		w.RegisterComponents(s.GetComponentDeps())
		w.addSystem(s)
	}
	// link up all systems' dependencies
	for _, s := range systems {
		w.linkSystemDependencies(s)
	}
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
	logicName := fmt.Sprintf("%s-update", name)
	w.systemsRunner.Add(
		&LogicUnit{
			name:    logicName,
			worldID: w.systemsIDs[s],
			f:       s.Update,
			active:  true})
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
		if f.Tag.Get("sameriver-system-dependency") != "" {
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
				panic(fmt.Sprintf("%s %v of %s dependency could not be "+
					"resolved. No system found of type %v.",
					f.Name, f.Type, sType.Elem().Name(), f.Type))
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

func (w *World) AddWorldLogic(Name string, F func()) *LogicUnit {
	if _, ok := w.worldLogics[Name]; ok {
		panic(fmt.Sprintf("double-add of world logic %s", Name))
	}
	l := &LogicUnit{
		name:    Name,
		f:       F,
		active:  false,
		worldID: w.IdGen.Next(),
	}
	w.worldLogics[Name] = l
	w.worldLogicsRunner.Add(l)
	return l
}

func (w *World) RemoveWorldLogic(Name string) {
	if logic, ok := w.worldLogics[Name]; ok {
		w.worldLogicsRunner.Remove(logic.worldID)
		delete(w.worldLogics, Name)
	}
}

func (w *World) ActivateAllWorldLogics() {
	w.worldLogicsRunner.ActivateAll()
}

func (w *World) DeactivateAllWorldLogics() {
	w.worldLogicsRunner.DeactivateAll()
}

func (w *World) ActivateWorldLogic(name string) {
	w.SetWorldLogicActiveState(name, true)
}

func (w *World) DeactivateWorldLogic(name string) {
	w.SetWorldLogicActiveState(name, false)
}

func (w *World) SetWorldLogicActiveState(name string, state bool) {
	if logic, ok := w.worldLogics[name]; ok {
		logic.active = state
	}
}

func (w *World) SetPrimaryEntityLogic(e *Entity, F func()) *LogicUnit {
	// idempotent remove
	w.RemovePrimaryEntityLogic(e)
	l := e.MakeLogicUnit("primary", F)
	w.primaryEntityLogics[e.ID] = l
	w.primaryEntityLogicsRunner.Add(l)
	return l
}

func (w *World) RemovePrimaryEntityLogic(e *Entity) {
	if logic, ok := w.primaryEntityLogics[e.ID]; ok {
		w.primaryEntityLogicsRunner.Remove(logic.worldID)
		delete(w.primaryEntityLogics, e.ID)
	}
}

func (w *World) RemoveAllLogics(e *Entity) {
	w.RemovePrimaryEntityLogic(e)
	for _, logic := range e.Logics {
		w.logicUnitComponentRunner.Remove(logic.worldID)
	}
}

func (w *World) ActivateAllEntityLogics() {
	w.primaryEntityLogicsRunner.ActivateAll()
	w.logicUnitComponentRunner.ActivateAll()
}

func (w *World) DeactivateAllEntityLogics() {
	w.primaryEntityLogicsRunner.DeactivateAll()
	w.logicUnitComponentRunner.DeactivateAll()
}

func (w *World) ActivateEntityLogic(e *Entity) {
	w.setEntityPrimaryLogicActiveState(e, true)
	for _, logic := range e.Logics {
		logic.active = true
	}
}

func (w *World) DeactivateEntityLogic(e *Entity) {
	w.setEntityPrimaryLogicActiveState(e, false)
	for _, logic := range e.Logics {
		logic.active = false
	}
}

func (w *World) setEntityPrimaryLogicActiveState(e *Entity, state bool) {
	if logic, ok := w.primaryEntityLogics[e.ID]; ok {
		logic.active = state
	}
}

func (w *World) String() string {
	// TODO: implement
	return "TODO"
}

func (w *World) DumpStats() (stats map[string](map[string]float64)) {
	stats = make(map[string](map[string]float64))
	systemStats, systemTotal := w.systemsRunner.DumpStats()
	worldStats, worldTotal := w.worldLogicsRunner.DumpStats()
	entityStats, entityTotal := w.primaryEntityLogicsRunner.DumpStats()
	stats["system"] = systemStats
	stats["world"] = worldStats
	stats["entity"] = entityStats
	stats["totals"] = make(map[string]float64)
	stats["totals"]["system"] = systemTotal
	stats["totals"]["world"] = worldTotal
	stats["totals"]["entity"] = entityTotal
	if w.totalRuntime != nil {
		stats["totals"]["total"] = *w.totalRuntime
	} else {
		stats["totals"]["total"] = 0.0
	}
	return
}

func (w *World) DumpStatsString() string {
	stats := w.DumpStats()
	b, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
