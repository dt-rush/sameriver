package engine

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
	"unsafe"
)

type World struct {
	Width         int
	Height        int
	Ev            *EventBus
	Em            *EntityManager
	systems       []System
	logics        []LogicUnit
	logicRunIndex int
}

func NewWorld(width int, height int) *World {
	ev := NewEventBus()
	em := NewEntityManager(ev)
	w := World{
		Width:   width,
		Height:  height,
		Ev:      ev,
		Em:      em,
		systems: make([]System, 0),
		logics:  make([]LogicUnit, 0),
	}
	return &w
}

func (w *World) AddSystems(systems ...System) {
	// add all systems
	for _, s := range systems {
		w.addSystem(s)
	}
	// link up all systems' dependencies
	for _, s := range systems {
		w.linkSystemDependencies(s)
	}
}

func (w *World) addSystem(s System) {
	w.assertSystemTypeValid(reflect.TypeOf(s))
	w.systems = append(w.systems, s)
	s.LinkWorld(w)
}

func (w *World) assertSystemTypeValid(t reflect.Type) {
	if t.Kind() != reflect.Ptr {
		panic("Implementers of engine.System must be pointer-receivers")
	}
	typeName := t.Elem().Name()
	validName, _ := regexp.MatchString(".+System$", typeName)
	if !validName {
		panic(fmt.Sprintf("implementers of System must have a name "+
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

func (w *World) Update(dt_ms float64) {
	w.Em.Update()
	for _, s := range w.systems {
		s.Update(dt_ms)
	}
}

func (w *World) AddLogic(l LogicUnit) {
	w.logics = append(w.logics, l)
}

func (w *World) ActivateAllLogic() {
	for i, _ := range w.logics {
		w.logics[i].Active = true
	}
}

func (w *World) DeactivateAllLogic() {
	for i, _ := range w.logics {
		w.logics[i].Active = false
	}
}

func (w *World) ActivateLogic(name string) {
	w.SetLogicActiveState(name, true)
}

func (w *World) DeactivateLogic(name string) {
	w.SetLogicActiveState(name, false)
}

func (w *World) SetLogicActiveState(name string, state bool) {
	for i, _ := range w.logics {
		if w.logics[i].Name == name {
			w.logics[i].Active = state
		}
	}
}

// run as many logics as we can in the time limit, picking up
// where we left off next time (and returning the amount we overrun)
func (w *World) RunLogic(limit_ms int64) (overrun_ms int64) {
	startLogicRunIndex := w.logicRunIndex
	for limit_ms > 0 {
		t0 := time.Now()
		w.logics[w.logicRunIndex].F()
		elapsed_ms := time.Since(t0).Nanoseconds() / 1e6
		limit_ms -= elapsed_ms
		w.logicRunIndex = (w.logicRunIndex + 1) % len(w.logics)
		if w.logicRunIndex == startLogicRunIndex {
			break
		}
	}
	return limit_ms * -1
}
