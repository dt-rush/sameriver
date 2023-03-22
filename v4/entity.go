package sameriver

import (
	"bytes"
	"fmt"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type Entity struct {
	ID                int
	World             *World
	Active            bool
	Despawned         bool
	ComponentBitArray bitarray.BitArray
	Lists             []*UpdatedEntityList
	Logics            map[string]*LogicUnit
	funcs             *FuncSet
	mind              map[string]any
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

func (e *Entity) AddLogic(name string, F func(e *Entity, dt_ms float64)) *LogicUnit {
	closureF := func(dt_ms float64) {
		F(e, dt_ms)
	}
	l := e.makeLogicUnit(name, closureF)
	e.Logics[name] = l
	e.World.addEntityLogic(e, l)
	return l
}

func (e *Entity) AddLogicWithSchedule(name string, F func(e *Entity, dt_ms float64), period float64) *LogicUnit {
	l := e.AddLogic(name, F)
	runSchedule := NewTimeAccumulator(period)
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
		e.World.RuntimeSharer.RunnerMap["entities"].Remove(l)
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

func (e *Entity) AddFuncs(funcs map[string](func(e *Entity, params any) any)) {
	for name, f := range funcs {
		e.AddFunc(name, f)

	}
}

func (e *Entity) AddFunc(name string, f func(e *Entity, params any) any) {
	closureF := func(params any) any {
		return f(e, params)
	}
	e.funcs.Add(name, closureF)
}

func (e *Entity) RemoveFunc(name string) {
	e.funcs.Remove(name)
}

func (e *Entity) GetFunc(name string) func(any) any {
	return e.funcs.funcs[name]
}

func (e *Entity) HasFunc(name string) bool {
	return e.funcs.Has(name)
}

func (e *Entity) GetMind(name string) any {
	if v, ok := e.mind[name]; ok {
		return v
	}
	return nil
}

func (e *Entity) SetMind(name string, val any) {
	e.mind[name] = val
}

func (e *Entity) String() string {
	return fmt.Sprintf("{id:%d, tags:%s, components:%s}",
		e.ID,
		e.GetTagList("GenericTags").AsSlice(),
		e.World.em.components.BitArrayToString(e.ComponentBitArray),
	)
}

func EntitySliceToString(entities []*Entity) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, e := range entities {
		buf.WriteString(e.String())
		if i != len(entities)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}
