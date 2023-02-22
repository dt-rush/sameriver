package sameriver

import (
	"fmt"
	"math"

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
		debugGOAPPrintf("[][][] adding action %s", action.DisplayName())
		e.actions.Add(action)
		// link up modal setters for effs matching modal varnames
		for varName, _ := range action.effs {
			if modal, ok := e.modalVals[varName]; ok {
				debugGOAPPrintf("[][][]     adding modal setter for %s", varName)
				action.effModalSetters[varName] = modal.effModalSet
			}
		}
		// link up modal checks for pres matching modal varnames
		for varName, _ := range action.pres.vars {
			if modal, ok := e.modalVals[varName]; ok {
				action.preModalChecks[varName] = modal.check
			}
		}
	}
}

func (e *GOAPEvaluator) applyActionBasic(action *GOAPAction, ws *GOAPWorldState) (newWS *GOAPWorldState) {
	newWS = ws.copyOf()

	for varName, eff := range action.effs {
		op := action.ops[varName]
		x := ws.vals[varName]
		debugGOAPPrintf("     %s       applying %s::%d x %s%s%d(%d) ; = %d",
			color.InWhiteOverYellow(">>>"),
			action.DisplayName(), action.Count, varName, op, eff.val, x,
			action.Count*eff.f(x))
		newWS.vals[varName] = action.Count * eff.f(x)
	}
	debugGOAPPrintf("            ws after action: %v", newWS.vals)

	return newWS
}

func (e *GOAPEvaluator) applyActionModal(action *GOAPAction, ws *GOAPWorldState) (newWS *GOAPWorldState) {

	newWS = e.applyActionBasic(action, ws)
	for varName, eff := range action.effs {
		op := action.ops[varName]
		x := ws.vals[varName]
		debugGOAPPrintf("    %s        applying %s::%d x %s%s%d(%d) ; = %d",
			color.InPurpleOverWhite(" >>>modal "),
			action.DisplayName(), action.Count, varName, op, eff.val, x,
			action.Count*eff.f(x))
		// do modal set
		if setter, ok := action.effModalSetters[varName]; ok {
			setter(newWS, op, action.Count*eff.val)
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
		remainings.nUnfulfilled += len(preRemaining.goal.vars)
		remainings.pres = append(remainings.pres, preRemaining)

		ws = e.applyActionBasic(action, ws)
	}
	debugGOAPPrintf("  --- ws after path: %v", ws.vals)
	mainRemaining := main.remaining(ws)
	remainings.nUnfulfilled += len(mainRemaining.goal.vars)
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
	return len(remaining.goal.vars) == 0
}

func (e *GOAPEvaluator) validateForward(path *GOAPPath, start *GOAPWorldState, main *GOAPGoal) bool {

	ws := start.copyOf()
	for _, action := range path.path {
		if len(action.pres.vars) > 0 && !e.presFulfilled(action, ws) {
			debugGOAPPrintf(">>>>>>> in validateForward, %s was not fulfilled", action.DisplayName())
			return false
		}
		ws = e.applyActionModal(action, ws)
	}
	endRemaining := main.remaining(ws)
	if len(endRemaining.goal.vars) != 0 {
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

func (e *GOAPEvaluator) actionHelps(
	start *GOAPWorldState,
	action *GOAPAction,
	path *GOAPPath,
	prependAppendFlag int) (scale int, helpful bool) {

	if DEBUG_GOAP {
		if prependAppendFlag == GOAP_PATH_PREPEND {
			debugGOAPPrintf(color.InBlueOverGray(fmt.Sprintf("checking if %s can be prepended", action.DisplayName())))
		}
		if prependAppendFlag == GOAP_PATH_APPEND {
			debugGOAPPrintf(color.InGreenOverGray(fmt.Sprintf("checking if %s can be appended", action.DisplayName())))
		}
	}

	actionChangesVarWell := func(
		varName string,
		interval *utils.NumericInterval,
		action *GOAPAction) (scale int, helpful bool) {

		debugGOAPPrintf("    Considering effs of %s for var %s. effs: %v", action.DisplayName(), varName, action.effs)
		for effVarName, eff := range action.effs {
			if varName == effVarName {
				debugGOAPPrintf("      [ ] eff affects var: %v; is it satisfactory/closer?", effVarName)
				// special handler for =
				if eff.op == "=" {
					switch prependAppendFlag {
					case GOAP_PATH_PREPEND:
						if interval.Diff(float64(eff.f(start.vals[varName]))) == 0 {
							debugGOAPPrintf("      [x] eff satisfactory")
							return 1, true
						} else {
							debugGOAPPrintf("      [_] eff not satisfactory")
							return -1, false
						}
					case GOAP_PATH_APPEND:
						if interval.Diff(float64(eff.f(path.endState.vals[varName]))) == 0 {
							debugGOAPPrintf("      [x] eff satisfactory")
							return 1, true
						} else {
							debugGOAPPrintf("      [_] eff not satisfactory")
							return -1, false
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
					// compute how many of this action we need
					// if we had diff 0, we just need one
					if actionDiff == 0 {
						return 1, true
					}
					// but if diff is nonzero, we need some scale
					// (note that diff is missing 1 val since we computed
					// the diff after applying 1 of the action, so we do
					// some tricky stuff to reconstruct the original diff)
					if actionDiff != 0 {
						var diffMagnitude float64
						if eff.op == "-" {
							diffMagnitude = -needToBeat
						} else if eff.op == "+" {
							diffMagnitude = needToBeat
						}
						scale := int(math.Ceil(diffMagnitude / float64(eff.val)))
						return scale, true
					}
				} else {
					debugGOAPPrintf("      [_] eff not closer")
					return -1, false
				}
			}
		}
		return -1, false
	}

	helpsGoal := func(goal *GOAPGoal) (scale int, helpful bool) {
		for varName, interval := range goal.vars {
			scale, helpful := actionChangesVarWell(varName, interval, action)
			if helpful {
				return scale, true
			}
		}
		return -1, false
	}

	if scale, helpful := helpsGoal(path.remainings.main.goal); helpful {
		return scale, true
	}
	for _, pre := range path.remainings.pres {
		if scale, helpful := helpsGoal(pre.goal); helpful {
			return scale, true
		}
	}
	return -1, false
}
