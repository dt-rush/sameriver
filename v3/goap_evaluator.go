package sameriver

import (
	"fmt"
	"math"

	"github.com/TwiN/go-color"
)

type GOAPEvaluator struct {
	modalVals map[string]GOAPModalVal
	actions   *GOAPActionSet
	// map of [varName](Set (in the sense of a map->bool) of actions that affect that var)
	varActions map[string](map[*GOAPAction]bool)
}

func NewGOAPEvaluator() *GOAPEvaluator {
	return &GOAPEvaluator{
		modalVals:  make(map[string]GOAPModalVal),
		actions:    NewGOAPActionSet(),
		varActions: make(map[string](map[*GOAPAction]bool)),
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

func (e *GOAPEvaluator) actionAffectsVar(action *GOAPAction, varName string) {
	if _, mapExists := e.varActions[varName]; !mapExists {
		e.varActions[varName] = make(map[*GOAPAction]bool)
	}
	e.varActions[varName][action] = true
}

func (e *GOAPEvaluator) AddActions(actions ...*GOAPAction) {
	for _, action := range actions {
		logGOAPDebug("[][][] adding action %s", action.DisplayName())
		e.actions.Add(action)
		// link up modal setters for effs matching modal varnames
		for varName := range action.effs {
			e.actionAffectsVar(action, varName)
			if modal, ok := e.modalVals[varName]; ok {
				logGOAPDebug("[][][]     adding modal setter for %s", varName)
				action.effModalSetters[varName] = modal.effModalSet
			}
		}
		// link up modal checks for pres matching modal varnames
		for varName := range action.pres.vars {
			if modal, ok := e.modalVals[varName]; ok {
				action.preModalChecks[varName] = modal.check
			}
		}
	}
}

func (e *GOAPEvaluator) applyActionBasic(
	action *GOAPAction, ws *GOAPWorldState, makeCopy bool) *GOAPWorldState {

	if makeCopy {
		ws = ws.CopyOf()
	}

	for varName, eff := range action.effs {
		op := action.ops[varName]
		x := ws.vals[varName]
		if DEBUG_GOAP {
			logGOAPDebug("     %s       applying %s::%d x %s%s%d(%d) ; = %d",
				color.InWhiteOverYellow(">>>"),
				action.DisplayName(), action.Count, varName, op, eff.val, x,
				eff.f(x))
		}
		ws.vals[varName] = eff.f(x)
	}
	if DEBUG_GOAP {
		logGOAPDebug(color.InBlueOverWhite(fmt.Sprintf("            ws after action: %v", ws.vals)))
	}

	return ws
}

func (e *GOAPEvaluator) applyActionModal(action *GOAPAction, ws *GOAPWorldState) (newWS *GOAPWorldState) {

	newWS = e.applyActionBasic(action, ws, false)
	for varName, eff := range action.effs {
		op := action.ops[varName]
		x := ws.vals[varName]
		logGOAPDebug("    %s        applying %s::%d x %s%s%d(%d) ; = %d",
			color.InPurpleOverWhite(" >>>modal "),
			action.DisplayName(), action.Count, varName, op, eff.val, x,
			eff.f(x))
		// do modal set
		if setter, ok := action.effModalSetters[varName]; ok {
			setter(newWS, op, action.Count*eff.val)
		}
	}
	if DEBUG_GOAP {
		logGOAPDebug(color.InBlueOverWhite(fmt.Sprintf("            ws after action: %v", newWS.vals)))
	}

	// re-check any modal vals
	for varName := range newWS.vals {
		if modalVal, ok := e.modalVals[varName]; ok {
			logGOAPDebug("              re-checking modal val %s", varName)
			newWS.vals[varName] = modalVal.check(newWS)
		}
	}

	if DEBUG_GOAP {
		logGOAPDebug(color.InPurpleOverWhite(fmt.Sprintf("            ws after re-checking modal vals: %v", newWS.vals)))
	}
	return newWS
}

func (e *GOAPEvaluator) computeRemainingsOfPath(path *GOAPPath, start *GOAPWorldState, main *GOAPGoal) {
	ws := start.CopyOf()
	path.remainings = NewGOAPGoalRemainingSurface()
	// there is a goal in the surface for each action's pre + 1 for the end main goal
	path.remainings.surface = make([]*GOAPGoalRemaining, len(path.path)+1)
	// create the storage space for statesAlong
	// consider, a path [A B C] will have 4 states: [start, postA, postB, postC (end)]
	path.statesAlong = make([]*GOAPWorldState, len(path.path)+1)
	path.statesAlong[0] = ws
	for i, action := range path.path {
		path.remainings.surface[i] = action.pres.remaining(ws)
		ws = e.applyActionBasic(action, ws, true)
		path.statesAlong[i+1] = ws
	}
	path.remainings.surface[len(path.path)] = main.remaining(ws)
	logGOAPDebug("  --- ws after path: %v", ws.vals)
}

func (e *GOAPEvaluator) presFulfilled(a *GOAPAction, ws *GOAPWorldState) bool {
	logGOAPDebug("Checking presFulfilled")
	modifiedWS := ws.CopyOf()
	for varName, checkF := range a.preModalChecks {
		modifiedWS.vals[varName] = checkF(ws)
	}
	remaining := a.pres.remaining(modifiedWS)
	return len(remaining.goalLeft) == 0
}

func (e *GOAPEvaluator) validateForward(path *GOAPPath, start *GOAPWorldState, main *GOAPGoal) bool {

	ws := start.CopyOf()
	for _, action := range path.path {
		if len(action.pres.vars) > 0 && !e.presFulfilled(action, ws) {
			logGOAPDebug(">>>>>>> in validateForward, %s was not fulfilled", action.DisplayName())
			return false
		}
		ws = e.applyActionModal(action, ws)
	}
	endRemaining := main.remaining(ws)
	if len(endRemaining.goalLeft) != 0 {
		logGOAPDebug(">>>>>>> in validateForward, main goal was not fulfilled at end of path")
		return false
	}
	return true
}

func (e *GOAPEvaluator) actionHelpsToInsert(
	start *GOAPWorldState,
	path *GOAPPath,
	insertionIx int,
	action *GOAPAction) (scale int, helpful bool) {

	actionChangesVarWell := func(
		varName string,
		interval *NumericInterval,
		action *GOAPAction) (scale int, helpful bool) {

		if DEBUG_GOAP {
			logGOAPDebug("    Considering effs of %s for var %s. effs: %v", action.DisplayName(), varName, action.effs)
		}
		for effVarName, eff := range action.effs {
			if varName != effVarName {
				logGOAPDebug("      [_] eff for %s doesn't affect var %s", effVarName, varName)
			} else {
				logGOAPDebug("      [ ] eff affects var: %s; is it satisfactory/closer?", effVarName)
				// special handler for =
				if eff.op == "=" {
					if interval.Diff(float64(eff.f(start.vals[varName]))) == 0 {
						logGOAPDebug("      [x] eff satisfactory")
						return 1, true
					} else {
						logGOAPDebug("      [_] eff not satisfactory")
						return -1, false
					}
				}
				// all other cases
				stateAtPoint := path.statesAlong[insertionIx].vals[varName]
				needToBeat := interval.Diff(float64(stateAtPoint))
				actionDiff := interval.Diff(float64(eff.f(stateAtPoint)))
				if DEBUG_GOAP {
					logGOAPDebug(path.String())
					logGOAPDebug("            ws[%s] before insertion: %d", varName, stateAtPoint)
					logGOAPDebug("              needToBeat diff: %d", int(needToBeat))
					logGOAPDebug("              actionDiff: %d", int(actionDiff))
				}
				if math.Abs(actionDiff) < math.Abs(needToBeat) {
					logGOAPDebug("      [X] eff closer")
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
					logGOAPDebug("      [_] eff not closer")
					return -1, false
				}
			}
		}
		return -1, false
	}

	helpsGoal := func(goalLeft map[string]*NumericInterval) (scale int, helpful bool) {
		for varName, interval := range goalLeft {
			logGOAPDebug("    - considering effect on %s", varName)
			affectors := e.varActions[varName]
			if _, affects := affectors[action]; affects {
				scale, helpful := actionChangesVarWell(varName, interval, action)
				if helpful {
					return scale, true
				}
			}
		}
		return -1, false
	}

	return helpsGoal(path.remainings.surface[insertionIx].goalLeft)
}
