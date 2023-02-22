package sameriver

import (
	"fmt"
	"strings"

	"github.com/dt-rush/sameriver/v2/utils"

	"github.com/TwiN/go-color"
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

func (e *GOAPEvaluator) AddModalVals(vals ...GOAPModalVal) {
	for _, val := range vals {
		e.modalVals[val.name] = val
	}
}

func (e *GOAPEvaluator) PopulateModalStartState(ws *GOAPWorldState) {
	for varName, val := range e.modalVals {
		ws.vals[varName] = val.check(ws)
	}
}

func (e *GOAPEvaluator) AddActions(actions ...*GOAPAction) {
	for _, action := range actions {
		debugGOAPPrintf("[][][] adding action %s", action.name)
		e.actions.Add(action)
		// link up modal setters for effs matching modal varnames
		for spec, _ := range action.effs {
			split := strings.Split(spec, ",")
			varName := split[0]
			if modal, ok := e.modalVals[varName]; ok {
				debugGOAPPrintf("[][][]     adding modal setter for %s", varName)
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

func (e *GOAPEvaluator) applyActionBasic(action *GOAPAction, ws *GOAPWorldState) (newWS *GOAPWorldState) {
	newWS = ws.copyOf()

	for spec, eff := range action.effs {
		split := strings.Split(spec, ",")
		varName, _ := split[0], split[1]
		x := ws.vals[varName]
		debugGOAPPrintf("            applying %s::%s%d(%d) ; = %d", action.name, spec, eff.val, x, eff.f(x))
		newWS.vals[varName] = eff.f(x)
	}
	debugGOAPPrintf("            ws after action: %v", newWS.vals)

	return newWS
}

func (e *GOAPEvaluator) applyActionModal(action *GOAPAction, ws *GOAPWorldState) (newWS *GOAPWorldState) {
	newWS = ws.copyOf()

	for spec, eff := range action.effs {
		split := strings.Split(spec, ",")
		varName, op := split[0], split[1]
		x := ws.vals[varName]
		debugGOAPPrintf("            applying %s::%s%d(%d) ; = %d", action.name, spec, eff.val, x, eff.f(x))
		newWS.vals[varName] = eff.f(x)
		// do modal set
		if setter, ok := action.effModalSetters[varName]; ok {
			setter(newWS, op, eff.val)
		}
	}
	debugGOAPPrintf("            ws after action: %v", newWS.vals)

	// re-check any modal vals
	for varName, _ := range newWS.vals {
		if modalVal, ok := e.modalVals[varName]; ok {
			debugGOAPPrintf("              re-checking modal val %s", varName)
			newWS.vals[varName] = modalVal.check(newWS)
		}
	}

	debugGOAPPrintf("            ws after re-checking modal vals: %v", newWS.vals)
	return newWS
}

func (e *GOAPEvaluator) remainingsOfPath(path *GOAPPath, start *GOAPWorldState, main *GOAPGoal) (remainings *GOAPGoalRemainingSurface) {
	ws := start.copyOf()
	remainings = NewGOAPGoalRemainingSurface()
	remainings.path = path
	for _, action := range path.path {
		preRemaining := action.pres.remaining(ws)
		remainings.nUnfulfilled += len(preRemaining.goal.goals)
		remainings.pres = append(remainings.pres, preRemaining)

		ws = e.applyActionBasic(action, ws)
	}
	debugGOAPPrintf("  --- ws after path: %v", ws.vals)
	mainRemaining := main.remaining(ws)
	remainings.nUnfulfilled += len(mainRemaining.goal.goals)
	remainings.main = mainRemaining
	path.remainings = remainings
	path.endState = ws

	return remainings
}

func (e *GOAPEvaluator) presFulfilled(a *GOAPAction, ws *GOAPWorldState) bool {
	debugGOAPPrintf("Checking presFulfilled")
	modifiedWS := ws.copyOf()
	for varName, checkF := range a.preModalChecks {
		modifiedWS.vals[varName] = checkF(ws)
	}
	remaining := a.pres.remaining(modifiedWS)
	return len(remaining.goal.goals) == 0
}

func (e *GOAPEvaluator) validateForward(path *GOAPPath, start *GOAPWorldState, main *GOAPGoal) bool {

	ws := start.copyOf()
	for _, action := range path.path {
		if len(action.pres.goals) > 0 && !e.presFulfilled(action, ws) {
			debugGOAPPrintf(">>>>>>> in validateForward, %s was not fulfilled", action.name)
			return false
		}
		ws = e.applyActionModal(action, ws)
	}
	endRemaining := main.remaining(ws)
	if len(endRemaining.goal.goals) != 0 {
		debugGOAPPrintf(">>>>>>> in validateForward, main goal was not fulfilled at end of path")
		return false
	}
	return true
}

func (e *GOAPEvaluator) tryPrepend(
	start *GOAPWorldState,
	action *GOAPAction,
	path *GOAPPath,
	goal *GOAPGoal) *GOAPPQueueItem {

	before := path.remainings
	prepended := path.prepended(action)
	if e.remainingsOfPath(prepended, start, goal).isCloser(before) {
		return &GOAPPQueueItem{path: prepended}
	} else {
		return nil
	}
}

func (e *GOAPEvaluator) tryAppend(
	start *GOAPWorldState,
	action *GOAPAction,
	path *GOAPPath,
	goal *GOAPGoal) *GOAPPQueueItem {

	before := path.remainings
	appended := path.appended(action)
	if e.remainingsOfPath(appended, start, goal).isCloser(before) {
		return &GOAPPQueueItem{path: appended}
	} else {
		return nil
	}
}

func (e *GOAPEvaluator) actionMightHelp(
	start *GOAPWorldState,
	action *GOAPAction,
	path *GOAPPath,
	prependAppendFlag int) bool {

	if prependAppendFlag == GOAP_PATH_PREPEND {
		debugGOAPPrintf(color.InBlueOverGray(fmt.Sprintf("checking if %s can be prepended", action.name)))
	}
	if prependAppendFlag == GOAP_PATH_APPEND {
		debugGOAPPrintf(color.InGreenOverGray(fmt.Sprintf("checking if %s can be appended", action.name)))
	}

	actionChangesVarWell := func(spec string, interval *utils.NumericInterval, action *GOAPAction) bool {
		split := strings.Split(spec, ",")
		varName := split[0]
		debugGOAPPrintf("    Considering effs of %s for var %s. effs: %v", action.name, varName, action.effs)
		for effSpec, eff := range action.effs {
			split = strings.Split(effSpec, ",")
			effVarName, op := split[0], split[1]
			if varName == effVarName {
				debugGOAPPrintf("      [ ] eff affects var: %v; is it satisfactory/closer?", effVarName)
				// special handler for =
				if op == "=" {
					switch prependAppendFlag {
					case GOAP_PATH_PREPEND:
						if interval.Diff(float64(eff.f(start.vals[varName]))) == 0 {
							debugGOAPPrintf("      [x] eff satisfactory")
							return true
						} else {
							debugGOAPPrintf("      [_] eff not satisfactory")
							return false
						}
					case GOAP_PATH_APPEND:
						if interval.Diff(float64(eff.f(path.endState.vals[varName]))) == 0 {
							debugGOAPPrintf("      [x] eff satisfactory")
							return true
						} else {
							debugGOAPPrintf("      [_] eff not satisfactory")
							return false
						}
					}
				}
				// all other cases
				var needToBeat, actionDiff float64
				switch prependAppendFlag {
				case GOAP_PATH_PREPEND:
					needToBeat = interval.Diff(float64(start.vals[varName]))
					actionDiff = interval.Diff(float64(eff.f(start.vals[varName])))
				case GOAP_PATH_APPEND:
					needToBeat = interval.Diff(float64(path.endState.vals[varName]))
					actionDiff = interval.Diff(float64(eff.f(path.endState.vals[varName])))
				}
				if actionDiff < needToBeat {
					debugGOAPPrintf("      [X] eff closer")
					return true
				} else {
					debugGOAPPrintf("      [_] eff not closer")
				}
			}
		}
		return false
	}

	mightHelpGoal := func(goal *GOAPGoal) bool {
		for spec, interval := range goal.goals {
			if actionChangesVarWell(spec, interval, action) {
				return true
			}
		}
		return false
	}

	if mightHelpGoal(path.remainings.main.goal) {
		return true
	}
	for _, pre := range path.remainings.pres {
		if mightHelpGoal(pre.goal) {
			return true
		}
	}
	return false
}
