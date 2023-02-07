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
	for _, action := range actions {
		e.actions.Add(action)
		// link up modal setters for effs matching modal varnames
		for spec, _ := range action.effs {
			split := strings.Split(spec, ",")
			varName := split[0]
			if modal, ok := e.modalVals[varName]; ok {
				action.effModalSetters[varName] = modal.effModalSet
			}
		}
		// link up modal checks for pres matching modal varnames
		for spec, _ := range action.pres.goals {
			split := strings.Split(spec, ",")
			varName := split[0]
			if modal, ok := e.modalVals[varName]; ok {
				action.preModalChecks[varName] = modal.check
			}
		}
	}
}

func (e *GOAPEvaluator) applyAction(action *GOAPAction, ws *GOAPWorldState) (newWS *GOAPWorldState) {
	newWS = ws.copyOf()
	for spec, eff := range action.effs {
		split := strings.Split(spec, ",")
		varName, op := split[0], split[1]
		var x int
		if val, ok := ws.vals[varName]; ok {
			x = val
		} else {
			x = 0
		}
		newWS.vals[varName] = eff.f(x)
		// do modal set
		if setter, ok := action.effModalSetters[varName]; ok {
			setter(newWS, op, eff.val)
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
	for varName, checkF := range a.preModalChecks {
		modifiedWS.vals[varName] = checkF(ws)
	}
	remaining, _ := a.pres.remaining(modifiedWS)
	return len(remaining.goals) == 0
}
