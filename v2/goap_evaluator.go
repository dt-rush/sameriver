package sameriver

import (
	"strings"
	"time"

	"github.com/dt-rush/sameriver/v2/utils"
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
		e.actions.Add(action)
		// link up modal setters for effs matching modal varnames
		for spec, _ := range action.effs {
			split := strings.Split(spec, ",")
			varName := split[0]
			if modal, ok := e.modalVals[varName]; ok {
				debugGOAPPrintf("[][][] adding modal setter for %s", varName)
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
	t0 := time.Now()
	newWS = ws.copyOf()
	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	debugGOAPPrintf("============================================================================== took %f ms to copy ws", dt_ms)

	t0 = time.Now()
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
	dt_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
	debugGOAPPrintf("================================================================================ took %f ms to apply effs", dt_ms)

	// re-check any modal vals
	t0 = time.Now()
	for varName, _ := range newWS.vals {
		if modalVal, ok := e.modalVals[varName]; ok {
			debugGOAPPrintf("              re-checking modal val %s", varName)
			newWS.vals[varName] = modalVal.check(newWS)
		}
	}
	dt_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
	debugGOAPPrintf("================================================================================== took %f ms to re-check modal vals", dt_ms)

	debugGOAPPrintf("            ws after re-checking modal vals: %v", newWS.vals)
	return newWS
}

func (e *GOAPEvaluator) remainingsOfPath(path *GOAPPath, start *GOAPWorldState, main *GOAPGoal) (remainings *GOAPGoalRemainingSurface) {
	t0 := time.Now()
	ws := start.copyOf()
	remainings = NewGOAPGoalRemainingSurface()
	remainings.path = path
	for _, action := range path.path {
		t1 := time.Now()
		preRemaining := action.pres.remaining(ws)
		remainings.nUnfulfilled += len(preRemaining.goal.goals)
		remainings.pres = append(remainings.pres, preRemaining)
		dt_ms := float64(time.Since(t1).Nanoseconds()) / 1.0e6
		debugGOAPPrintf(" ============================================================================================= took %f ms to check presRemaining for action", dt_ms)

		t2 := time.Now()
		ws = e.applyAction(action, ws)
		dt_ms_2 := float64(time.Since(t2).Nanoseconds()) / 1.0e6
		debugGOAPPrintf(" ==================================================================================================== took %f ms to apply action", dt_ms_2)

	}
	debugGOAPPrintf("  --- ws after path: %v", ws.vals)
	mainRemaining := main.remaining(ws)
	remainings.nUnfulfilled += len(mainRemaining.goal.goals)
	remainings.main = mainRemaining
	path.remainings = remainings
	path.endState = ws

	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	debugGOAPPrintf("================================================================================================== Took %f ms to check remainings of path by applying it", dt_ms)

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
		ws = e.applyAction(action, ws)
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

	var appendedPrependedMsg string
	if prependAppendFlag == GOAP_PATH_PREPEND {
		appendedPrependedMsg = "prepended"
	}
	if prependAppendFlag == GOAP_PATH_APPEND {
		appendedPrependedMsg = "appended"
	}
	debugGOAPPrintf("checking if %s can be %s", action.name, appendedPrependedMsg)

	actionChangesVarWell := func(spec string, interval *utils.NumericInterval, action *GOAPAction) bool {
		debugGOAPPrintf("    Considering effs of %s: %v", action.name, action.effs)
		split := strings.Split(spec, ",")
		varName := split[0]
		for effSpec, eff := range action.effs {
			split = strings.Split(effSpec, ",")
			effVarName := split[0]
			if varName == effVarName {
				debugGOAPPrintf("      [ ] eff affects var: %v", effSpec)
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
					debugGOAPPrintf("      [X] eff is good for var")
					return true
				} else {
					debugGOAPPrintf("      [_] eff doesn't help var")
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
