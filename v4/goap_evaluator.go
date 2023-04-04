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

func (e *GOAPEvaluator) checkModalInto(varName string, ws *GOAPWorldState) {
	if _, ok := e.modalVals[varName]; ok {
		ws.vals[varName] = e.modalVals[varName].check(ws)
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
		for _, tg := range action.pres.temporalGoals {
			for varName := range tg.vars {
				if modal, ok := e.modalVals[varName]; ok {
					action.preModalChecks[varName] = modal.check
				}
			}
		}
	}
}

func (e *GOAPEvaluator) applyActionBasic(
	action *GOAPAction, ws *GOAPWorldState, makeCopy bool) *GOAPWorldState {

	if makeCopy {
		ws = ws.CopyOf()
	}

	logGOAPDebug("     %s   applying %s",
		color.InWhiteOverYellow(">>>"),
		action.DisplayName())
	for varName, eff := range action.effs {
		op := action.ops[varName]
		x := ws.vals[varName]
		if DEBUG_GOAP {
			logGOAPDebug("     %s       %d x %s%s%d(%d) ; = %d",
				color.InWhiteOverYellow(">>>"),
				action.Count, varName, op, eff.val, x,
				eff.f(action.Count, x))
		}
		ws.vals[varName] = eff.f(action.Count, x)
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
			eff.f(action.Count, x))
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

func (e *GOAPEvaluator) computeRemainingsOfPath(path *GOAPPath, start *GOAPWorldState, main *GOAPTemporalGoal) {
	ws := start.CopyOf()
	// one []*GOAPGoalRemaining for each action pre + 1 for main
	surfaceLen := len(path.path) + 1
	surface := newGOAPGoalRemainingSurface(surfaceLen)
	// create the storage space for statesAlong
	// consider, a path [A B C] will have 4 states: [start, postA, postB, postC (end)]
	path.statesAlong = make([]*GOAPWorldState, len(path.path)+1)
	path.statesAlong[0] = ws
	for i, action := range path.path {
		for _, tg := range action.pres.temporalGoals {
			surface.surface[i] = append(
				surface.surface[i],
				tg.remaining(ws))
		}
		ws = e.applyActionBasic(action, ws, true)
		path.statesAlong[i+1] = ws
	}
	for _, tg := range main.temporalGoals {
		surface.surface[surfaceLen-1] = append(
			surface.surface[surfaceLen-1],
			tg.remaining(ws))
	}
	path.remainings = surface
	// TODO: is this on path?
	logGOAPDebug("  --- ws after path: %v", ws.vals)
}

func (e *GOAPEvaluator) presFulfilled(a *GOAPAction, ws *GOAPWorldState) bool {
	logGOAPDebug("Checking presFulfilled")
	modifiedWS := ws.CopyOf()
	for varName, checkF := range a.preModalChecks {
		modifiedWS.vals[varName] = checkF(ws)
	}
	goalLeftCount := 0
	for _, tg := range a.pres.temporalGoals {
		goalLeftCount += len(tg.remaining(modifiedWS).goalLeft)
	}
	return goalLeftCount == 0
}

func (e *GOAPEvaluator) validateForward(path *GOAPPath, start *GOAPWorldState, main *GOAPTemporalGoal) bool {

	ws := start.CopyOf()
	for _, action := range path.path {
		if len(action.pres.temporalGoals) > 0 && !e.presFulfilled(action, ws) {
			logGOAPDebug(">>>>>>> in validateForward, %s was not fulfilled", action.DisplayName())
			return false
		}
		ws = e.applyActionModal(action, ws)
	}
	endRemainingCount := 0
	for _, tg := range main.temporalGoals {
		endRemainingCount += len(tg.remaining(ws).goalLeft)
	}
	if endRemainingCount != 0 {
		logGOAPDebug(">>>>>>> in validateForward, main goal was not fulfilled at end of path")
		return false
	}
	return true
}

func (e *GOAPEvaluator) actionHelpsToInsert(
	start *GOAPWorldState,
	path *GOAPPath,
	insertionIx int,
	goalToHelp *GOAPGoalRemaining,
	action *GOAPAction) (scale int, helpful bool) {

	actionChangesVarWell := func(
		varName string,
		interval *NumericInterval,
		action *GOAPAction) (scale int, helpful bool) {

		if DEBUG_GOAP {
			logGOAPDebug("    Considering effs of %s for var %s. effs: %v", color.InYellowOverWhite(action.DisplayName()), varName, action.effs)
		}
		for effVarName, eff := range action.effs {
			if varName != effVarName {
				logGOAPDebug("      [_] eff for %s doesn't affect var %s", effVarName, varName)
			} else {
				logGOAPDebug("      [ ] eff affects var: %s; is it satisfactory/closer?", effVarName)
				// all other cases
				stateAtPoint := path.statesAlong[insertionIx].vals[varName]
				needToBeat := interval.Diff(float64(stateAtPoint))
				actionDiff := interval.Diff(float64(eff.f(action.Count, stateAtPoint)))
				if DEBUG_GOAP {
					logGOAPDebug(path.String())
					logGOAPDebug("            ws[%s] = %d (before)", varName, stateAtPoint)
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

	return helpsGoal(goalToHelp.goalLeft)
}

func (e *GOAPEvaluator) setVarInStartIfNotDefined(start *GOAPWorldState, varName string) {
	if _, already := start.vals[varName]; !already {
		if _, isModal := e.modalVals[varName]; isModal {
			e.checkModalInto(varName, start)
			logGOAPDebug(color.InPurple(fmt.Sprintf("[ ] start.modal[%s] %d", varName, start.vals[varName])))
		} else {
			// NOTE: vars that don't have modal check default to 0
			logGOAPDebug(color.InYellow(fmt.Sprintf("[ ] %s not defined in GOAP start state, and no modal check exists. Defaulting to 0.", varName)))
			start.vals[varName] = 0
		}
	}
}
