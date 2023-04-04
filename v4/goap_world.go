package sameriver

import (
	"fmt"
	//	"math"
)

type GOAPWorldState struct {
	w *World
	// TODO: export vals
	vals map[string]int
	// TODO: change this to a map[int](map[string]any) [ID][component]
	modal map[string]any
	//
	ModalEntities map[string]*Entity
}

func (ws *GOAPWorldState) CopyOf() *GOAPWorldState {
	copyvals := make(map[string]int)
	for k, v := range ws.vals {
		copyvals[k] = v
	}
	copyModal := make(map[string]any)
	for k, v := range ws.modal {
		copyModal[k] = v
	}
	copyModalEntities := make(map[string]*Entity)
	for k, v := range ws.ModalEntities {
		copyModalEntities[k] = v
	}
	copyWS := &GOAPWorldState{
		ws.w,
		copyvals,
		copyModal,
		copyModalEntities,
	}
	return copyWS
}

func NewGOAPWorldState(vals map[string]int) *GOAPWorldState {
	ws := &GOAPWorldState{
		vals:  vals,
		modal: make(map[string]any),
	}
	if vals == nil {
		ws.vals = make(map[string]int)
	}
	return ws
}

func (ws *GOAPWorldState) ecKey(e *Entity, name ComponentID) string {
	return fmt.Sprintf("%d-%s", e.ID, ws.w.em.components.strings[name])
}

func (ws *GOAPWorldState) GetModal(e *Entity, name ComponentID) any {
	if val, ok := ws.modal[ws.ecKey(e, name)]; ok {
		return val
	} else {
		return e.GetVal(name)
	}
}

func (ws *GOAPWorldState) SetModal(e *Entity, name ComponentID, val any) {
	ws.modal[ws.ecKey(e, name)] = val
}
