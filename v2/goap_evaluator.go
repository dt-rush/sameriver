package sameriver

import (
	"strings"
)

type GOAPEvaluator struct {
	modalVals map[string]GOAPModalVal
	actions   *GOAPActionSet
}

func NewGOAPEvaluator() *GOAPEvaluator {
	return &GOAPEvaluator{
		modalVals: make(map[string]GOAPModalVal),
		actions:   NewGOAPActionSet(),
	}
}

func (e *GOAPEvaluator) addModalVals(vals ...GOAPModalVal) {
	for _, val := range vals {
		e.modalVals[val.name] = val
	}
}

func (e *GOAPEvaluator) addActions(actions ...*GOAPAction) {
	e.actions.Add(actions...)
}

func (e *GOAPEvaluator) applyAction(action *GOAPAction, ws *GOAPWorldState) (newWS *GOAPWorldState) {
	newWS = ws.copyOf()
	for spec, f := range action.effs {
		split := strings.Split(spec, ",")
		varName := split[0]
		if x, ok := ws.vals[varName]; ok {
			newWS.vals[varName] = f(x)
		} else {
			newWS.vals[varName] = f(0)
		}
		// do modal set
		if setter, ok := action.effModalSetters[varName]; ok {
			setter(newWS)
		}
	}
	// re-check any modal vals
	for varName, _ := range newWS.vals {
		if modalVal, ok := e.modalVals[varName]; ok {
			newWS.vals[varName] = modalVal.check(newWS)
		}
	}
	return newWS
}

func (e *GOAPEvaluator) applyPath(path []*GOAPAction, ws *GOAPWorldState) *GOAPWorldState {
	ws = ws.copyOf()
	for _, action := range path {
		ws = e.applyAction(action, ws)
	}
	return ws
}

func (e *GOAPEvaluator) presFulfilled(a *GOAPAction, ws *GOAPWorldState) bool {
	modifiedWS := ws.copyOf()
	for varName, modalVal := range a.preModalChecks {
		modifiedWS.vals[varName] = modalVal.check(ws)
	}
	remaining, _ := a.pres.goalRemaining(modifiedWS)
	return len(remaining.goals) == 0
}

/*
func (ws GOAPWorldState) mergeActionPres(action GOAPAction) GOAPWorldState {
	ws = ws.copyOf()
	for name, val := range action.pres {
		ws.vals[name] = resolveGOAPStateVal(val)
	}
	return ws
}
*/

// TODO: have the GOAPPlanner contain a GOAPEvaluator,
// and move the AddModalVal and AddAction to GOAPPlanner
/*
func (e *GOAPEvaluator) thoseThatCanPrepend(
	start GOAPWorldState,
	path []GOAPAction,
	goal GOAPGoal) (prependers *GOAPActionSet, solutions [][]GOAPAction) {

	helpers := NewGOAPActionSet()
	for _, action := range e.actions.set {
		if action.movesToward(ws, path) {
			helpers.Add(action)
		}
	}
	return helpers
}
*/
