package sameriver

import (
	"fmt"
)

// state values are either bool or GOAPCtxState (bool resolved by get())
type GOAPState interface{}

var EmptyGOAPState = map[string]GOAPState{}

// a state val whose value must be read from the entity/world
// and which can *set* values in the worldstate (and modal) as the action
// chain runs forward
type GOAPCtxStateVal struct {
	name string
	get  func() bool
	set  func(ws *GOAPWorldState)
}

type GOAPWorldState struct {
	Vals  map[string]GOAPState
	modal map[string]interface{}
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
		Logger.Printf("%s .get()...", name)
		return ctxStateVal.get()
	} else {
		return ws.Vals[name]
	}
}

func (ws GOAPWorldState) applyAction(action GOAPAction) GOAPWorldState {
	for name, val := range action.effs {
		if ctxStateVal, ok := val.(GOAPCtxStateVal); ok {
			ws.Vals[name] = ctxStateVal
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
	for name, _ := range other.Vals {
		if ws.get(name) == other.get(name) {
			return true
		}
	}
	return false
}

type GOAPAction struct {
	name string
	// values are either bool or GOAPCtxState (bool resolved by get())
	pres map[string]GOAPState
	effs map[string]GOAPState
}

func (a *GOAPAction) presFulfilled(ws GOAPWorldState) bool {
	state := NewGOAPWorldState(nil)
	for name, val := range a.pres {
		state.Vals[name] = val
	}
	return ws.fulfills(state)
}

type GOAPActionSet struct {
	set map[string]GOAPAction
}

func NewGOAPActionSet() *GOAPActionSet {
	return &GOAPActionSet{
		set: make(map[string]GOAPAction),
	}
}

func (as *GOAPActionSet) Add(actions ...GOAPAction) {
	for _, action := range actions {
		as.set[action.name] = action
	}
}

func (as *GOAPActionSet) thoseThatHelpFulfill(ws GOAPWorldState) *GOAPActionSet {
	helpers := NewGOAPActionSet()
	for _, action := range as.set {
		effState := NewGOAPWorldState(nil)
		effState.applyAction(action)
		if effState.isSubset(ws) {
			helpers.Add(action)
		}
	}
	return helpers
}
