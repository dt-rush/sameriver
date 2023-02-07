package sameriver

import (
	"fmt"
	//	"math"
)

type GOAPWorldState struct {
	// TODO: uncapitalise vals
	vals  map[string]int
	modal map[string]interface{}
}

func (ws *GOAPWorldState) copyOf() *GOAPWorldState {
	copyvals := make(map[string]int)
	for k, v := range ws.vals {
		copyvals[k] = v
	}
	copyModal := make(map[string]interface{})
	for k, v := range ws.modal {
		copyModal[k] = v
	}
	copyWS := &GOAPWorldState{
		copyvals,
		copyModal,
	}
	return copyWS
}

func NewGOAPWorldState(vals map[string]int) *GOAPWorldState {
	ws := &GOAPWorldState{
		vals:  vals,
		modal: make(map[string]interface{}),
	}
	if vals == nil {
		ws.vals = make(map[string]int)
	}
	return ws
}

func (ws *GOAPWorldState) ecKey(e *Entity, name string) string {
	return fmt.Sprintf("%d-%s", e.ID, name)
}

func (ws *GOAPWorldState) GetModal(e *Entity, name string) interface{} {
	if val, ok := ws.modal[ws.ecKey(e, name)]; ok {
		return val
	} else {
		return e.GetVal(name)
	}
}

func (ws *GOAPWorldState) SetModal(e *Entity, name string, val interface{}) {
	ws.modal[ws.ecKey(e, name)] = val
}
