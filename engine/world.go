package engine

import (
	"time"
)

type World struct {
	Width         int
	Height        int
	ev            *EventBus
	em            *EntityManager
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
		ev:      ev,
		em:      em,
		systems: make([]System, 0),
		logics:  make([]LogicUnit, 0),
	}
	return &w
}

func (w *World) AddSystem(s System) {
	w.systems = append(w.systems, s)
	s.LinkWorld(w)
}

func (w *World) Update(dt_ms float64) {
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
