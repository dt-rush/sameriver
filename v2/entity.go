package sameriver

import (
	"bytes"
	"fmt"

	"github.com/dt-rush/sameriver/v2/utils"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type Entity struct {
	ID                int
	World             *World
	WorldID           int
	Active            bool
	Despawned         bool
	ComponentBitArray bitarray.BitArray
	Lists             []*UpdatedEntityList
	Logics            map[string]*LogicUnit
	funcs             *FuncSet
}

func (e *Entity) LogicUnitName(name string) string {
	return fmt.Sprintf("entity-logic-%d-%s", e.ID, name)
}

func (e *Entity) makeLogicUnit(name string, F func(dt_ms float64)) *LogicUnit {
	return &LogicUnit{
		name:        e.LogicUnitName(name),
		f:           F,
		active:      true,
		worldID:     e.World.IdGen.Next(),
		runSchedule: nil,
	}
}

func (e *Entity) AddLogic(name string, F func(dt_ms float64)) *LogicUnit {
	l := e.makeLogicUnit(name, F)
	e.Logics[name] = l
	e.World.addEntityLogic(e, l)
	return l
}

func (e *Entity) AddLogicWithSchedule(name string, F func(dt_ms float64), period float64) *LogicUnit {
	l := e.AddLogic(name, F)
	runSchedule := utils.NewTimeAccumulator(period)
	l.runSchedule = &runSchedule
	return l
}

func (e *Entity) RemoveLogic(name string) {
	if _, ok := e.Logics[name]; !ok {
		panic(fmt.Sprintf("Trying to remove logic %s - but entity doesn't have it!", name))
	}
	e.World.removeEntityLogic(e, e.Logics[name])
	delete(e.Logics, name)
}

func (e *Entity) RemoveAllLogics() {
	for _, l := range e.Logics {
		e.World.removeEntityLogic(e, l)
	}
}

func (e *Entity) ActivateLogics() {
	for _, logic := range e.Logics {
		logic.Activate()
	}
}

func (e *Entity) DeactivateLogics() {
	for _, logic := range e.Logics {
		logic.Deactivate()
	}
}

func (e *Entity) AddFuncs(funcs map[string](func(interface{}) interface{})) {
	for name, f := range funcs {
		e.funcs.Add(name, f)
	}
}

func (e *Entity) AddFunc(name string, f func(interface{}) interface{}) {
	e.funcs.Add(name, f)
}

func (e *Entity) RemoveFunc(name string) {
	e.funcs.Remove(name)
}

func (e *Entity) HasFunc(name string) bool {
	return e.funcs.Has(name)
}

func (e *Entity) GetFunc(name string) func(interface{}) interface{} {
	return e.funcs.funcs[name]
}

func EntitySliceToString(entities []*Entity) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, e := range entities {
		buf.WriteString(fmt.Sprintf("%d", e.ID))
		if i != len(entities)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}
