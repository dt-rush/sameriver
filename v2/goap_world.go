package sameriver

import (
	"fmt"
)

type GOAPWorldState struct {
	Vals  map[string]GOAPState
	modal map[string]interface{}
}

func (ws GOAPWorldState) copyOf() GOAPWorldState {
	copyVals := make(map[string]GOAPState)
	for k, v := range ws.Vals {
		copyVals[k] = v
	}
	copyModal := make(map[string]interface{})
	for k, v := range ws.modal {
		copyModal[k] = v
	}
	copyWS := GOAPWorldState{
		copyVals,
		copyModal,
	}
	return copyWS
}

func NewGOAPWorldState(vals map[string]GOAPState) GOAPWorldState {
	ws := GOAPWorldState{
		Vals:  vals,
		modal: make(map[string]interface{}),
	}
	if vals == nil {
		ws.Vals = make(map[string]GOAPState)
	}
	return ws
}

func (ws GOAPWorldState) ecKey(e *Entity, name string) string {
	return fmt.Sprintf("%d-%s", e.ID, name)
}

func (ws GOAPWorldState) GetModal(e *Entity, name string) interface{} {
	if val, ok := ws.modal[ws.ecKey(e, name)]; ok {
		return val
	} else {
		return e.GetVal(name)
	}
}

func (ws GOAPWorldState) SetModal(e *Entity, name string, val interface{}) {
	ws.modal[ws.ecKey(e, name)] = val
}

func (ws GOAPWorldState) get(name string) interface{} {
	if ctxStateVal, ok := ws.Vals[name].(GOAPCtxStateVal); ok {
		return ctxStateVal.val
	} else {
		return ws.Vals[name]
	}
}

func (ws GOAPWorldState) applyAction(action GOAPAction) GOAPWorldState {
	ws = ws.copyOf()
	for name, val := range action.effs {
		if ctxStateVal, ok := val.(GOAPCtxStateVal); ok {
			ctxStateVal.set(&ws)
			ws.Vals[name] = ctxStateVal.val
		} else {
			ws.Vals[name] = val
		}
	}
	return ws
}

func (ws GOAPWorldState) fulfills(other GOAPWorldState) bool {
	for name, _ := range other.Vals {
		if ws.get(name) != other.get(name) {
			return false
		}
	}
	return true
}

func (ws GOAPWorldState) isSubset(other GOAPWorldState) bool {
	// Logger.Println("        isSubset")
	// Logger.Printf("        ws: %v", ws)
	// Logger.Printf("        other: %v", other)
	for name, _ := range other.Vals {
		// Logger.Printf("        %s?", name)
		if ws.get(name) == other.get(name) {
			// Logger.Printf("        true")
			return true
		}
	}
	// Logger.Printf("        false")
	return false
}

func (ws GOAPWorldState) unfulfilledBy(action GOAPAction) GOAPWorldState {
	ws = ws.copyOf()
	for name, val := range action.effs {
		if ctxStateVal, ok := val.(GOAPCtxStateVal); ok {
			val = ctxStateVal.val
		}
		if _, ok := ws.Vals[name]; ok {
			if ws.Vals[name] == val {
				delete(ws.Vals, name)
			}
		}
	}
	return ws
}

func (ws GOAPWorldState) mergeActionPres(action GOAPAction) GOAPWorldState {
	ws = ws.copyOf()
	for name, val := range action.pres {
		if ctxStateVal, ok := val.(GOAPCtxStateVal); ok {
			val = ctxStateVal.val
		}
		ws.Vals[name] = val
	}
	return ws
}
