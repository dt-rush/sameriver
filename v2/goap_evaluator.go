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

func (e *GOAPEvaluator) applyActionBasic(
	action *GOAPAction, ws *GOAPWorldState, makeCopy bool) *GOAPWorldState {

	if makeCopy {
		ws = ws.CopyOf()
	}

	for varName, eff := range action.effs {
		op := action.ops[varName]
		x := ws.vals[varName]
		if DEBUG_GOAP {
			debugGOAPPrintf("     %s       applying %s::%d x %s%s%d(%d) ; = %d",
				color.InWhiteOverYellow(">>>"),
				action.DisplayName(), action.Count, varName, op, eff.val, x,
				eff.f(x))
		}
		ws.vals[varName] = eff.f(x)
	}
	if DEBUG_GOAP {
		debugGOAPPrintf(color.InBlueOverWhite(fmt.Sprintf("            ws after action: %v", ws.vals)))
	}

	return ws
}

func (e *GOAPEvaluator) applyActionModal(action *GOAPAction, ws *GOAPWorldState) (newWS *GOAPWorldState) {

	newWS = e.applyActionBasic(action, ws, false)
	for varName, eff := range action.effs {
		op := action.ops[varName]
		x := ws.vals[varName]
		debugGOAPPrintf("    %s        applying %s::%d x %s%s%d(%d) ; = %d",
			color.InPurpleOverWhite(" >>>modal "),
			action.DisplayName(), action.Count, varName, op, eff.val, x,
			eff.f(x))
		// do modal set
		if setter, ok := action.effModalSetters[varName]; ok {
			setter(newWS, op, action.Count*eff.val)
		}
	}
	if DEBUG_GOAP {
		debugGOAPPrintf(color.InBlueOverWhite(fmt.Sprintf("            ws after action: %v", newWS.vals)))
	}

	// re-check any modal vals
	for varName, _ := range newWS.vals {
		if modalVal, ok := e.modalVals[varName]; ok {
			debugGOAPPrintf("              re-checking modal val %s", varName)
			newWS.vals[varName] = modalVal.check(newWS)
		}
	}

	if DEBUG_GOAP {
		debugGOAPPrintf(color.InPurpleOverWhite(fmt.Sprintf("            ws after re-checking modal vals: %v", newWS.vals)))
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
	debugGOAPPrintf("  --- ws after path: %v", ws.vals)
}

// action is the action to insert
// insertionIx is where to insert it
// before is the path *without* that action inserted
func (e *GOAPEvaluator) isBetter(
	insertionIx int, action *GOAPAction, path *GOAPPath, start *GOAPWorldState) (better bool) {

	// consider:
	// before      [A B C D E]
	// before.remainings: [Apre Bpre Cpre Dpre Epre Main]
	// after       [A B X C D E]
	// after.remainings: [Apre Bpre Xpre Cpre Dpre Epre Main]
	// insertionIx 2
	// we only need to check if [A B X] fulfills C's pre better than [A B]
	//
	// consider edge case:
	//
	// before []
	// before.remainings: [Main]
	// after [X]
	// after.remainings: [Xpre, Main]
	// insertionIx 0
	beforeRemaining := path.remainings.surface[insertionIx] // [2] = Cpre
	upToInsertion := make([]*GOAPAction, insertionIx+1)
	copy(upToInsertion[:insertionIx], path.path[:insertionIx]) // [:2] [A B]
	upToInsertion[insertionIx] = action                        // [A B X]
	relevantGoal := beforeRemaining.goal                       // Cpre
	ws := start.CopyOf()
	for _, action := range upToInsertion {
		ws = e.applyActionBasic(action, start, false)
	}
	insertionRemaining := relevantGoal.remaining(ws)
	return insertionRemaining.nUnfulfilled < beforeRemaining.nUnfulfilled
}

func (e *GOAPEvaluator) presFulfilled(a *GOAPAction, ws *GOAPWorldState) bool {
	debugGOAPPrintf("Checking presFulfilled")
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
			debugGOAPPrintf(">>>>>>> in validateForward, %s was not fulfilled", action.DisplayName())
			return false
		}
		ws = e.applyActionModal(action, ws)
	}
	endRemaining := main.remaining(ws)
	if len(endRemaining.goalLeft) != 0 {
		debugGOAPPrintf(">>>>>>> in validateForward, main goal was not fulfilled at end of path")
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
		interval *utils.NumericInterval,
		action *GOAPAction) (scale int, helpful bool) {

		if DEBUG_GOAP {
			debugGOAPPrintf("    Considering effs of %s for var %s. effs: %v", action.DisplayName(), varName, action.effs)
		}
		for effVarName, eff := range action.effs {
			if varName != effVarName {
				debugGOAPPrintf("      [_] eff for %s doesn't affect var %s", effVarName, varName)
			} else {
				debugGOAPPrintf("      [ ] eff affects var: %s; is it satisfactory/closer?", effVarName)
				// special handler for =
				if eff.op == "=" {
					if interval.Diff(float64(eff.f(start.vals[varName]))) == 0 {
						debugGOAPPrintf("      [x] eff satisfactory")
						return 1, true
					} else {
						debugGOAPPrintf("      [_] eff not satisfactory")
						return -1, false
					}
				}
				// all other cases
				stateAtPoint := path.statesAlong[insertionIx].vals[varName]
				needToBeat := interval.Diff(float64(stateAtPoint))
				actionDiff := interval.Diff(float64(eff.f(stateAtPoint)))
				if DEBUG_GOAP {
					debugGOAPPrintf(path.String())
					debugGOAPPrintf("            ws[%s] before insertion: %d", varName, stateAtPoint)
					debugGOAPPrintf("              needToBeat diff: %d", int(needToBeat))
					debugGOAPPrintf("              actionDiff: %d", int(actionDiff))
				}
				if math.Abs(actionDiff) < math.Abs(needToBeat) {
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

	helpsGoal := func(goalLeft map[string]*utils.NumericInterval) (scale int, helpful bool) {
		for varName, interval := range goalLeft {
			debugGOAPPrintf("    - considering effect on %s", varName)
			scale, helpful := actionChangesVarWell(varName, interval, action)
			if helpful {
				return scale, true
			}
		}
		return -1, false
	}

	return helpsGoal(path.remainings.surface[insertionIx].goalLeft)
}
