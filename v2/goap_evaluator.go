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

func (e *GOAPEvaluator) populateModalStartState(ws *GOAPWorldState) {
	for varName, val := range e.modalVals {
		ws.vals[varName] = val.check(ws)
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
		x := ws.vals[varName]
		debugGOAPPrintf("            applying %s::%s%d(%d) ; = %d", action.name, spec, eff.val, x, eff.f(x))
		newWS.vals[varName] = eff.f(x)
		// do modal set
		if setter, ok := action.effModalSetters[varName]; ok {
			setter(newWS, op, eff.val)
		}
	}
	debugGOAPPrintf("            ws after action: %v", newWS.vals)
	// TODO: do we need this?
	// re-check any modal vals
	/*
		for varName, _ := range newWS.vals {
			if modalVal, ok := e.modalVals[varName]; ok {
				newWS.vals[varName] = modalVal.check(newWS)
			}
		}
		debugGOAPPrintf("            ws after re-checking modal vals: %v", newWS.vals)
	*/
	return newWS
}

func (e *GOAPEvaluator) applyPath(path *GOAPPath, ws *GOAPWorldState) (result *GOAPWorldState) {
	result = ws.copyOf()
	for _, action := range path.path {
		result = e.applyAction(action, result)
	}
	return result
}

func (e *GOAPEvaluator) remainingsOfPath(path *GOAPPath, start *GOAPWorldState, main *GOAPGoal) (remainings *GOAPGoalRemainingSurface) {
	ws := start.copyOf()
	remainings = NewGOAPGoalRemainingSurface()
	remainings.path = path
	for _, action := range path.path {
		preRemaining := action.pres.remaining(ws)
		remainings.nUnfulfilled += len(preRemaining.goal.goals)
		remainings.pres = append(remainings.pres, preRemaining)
		ws = e.applyAction(action, ws)
	}
	debugGOAPPrintf("  --- remainingsOfPath: end state: %v", ws.vals)
	mainRemaining := main.remaining(ws)
	remainings.nUnfulfilled += len(mainRemaining.goal.goals)
	remainings.main = mainRemaining
	return remainings
}

func (e *GOAPEvaluator) presFulfilled(a *GOAPAction, ws *GOAPWorldState) bool {
	modifiedWS := ws.copyOf()
	for varName, checkF := range a.preModalChecks {
		modifiedWS.vals[varName] = checkF(ws)
	}
	remaining := a.pres.remaining(modifiedWS)
	return len(remaining.goal.goals) == 0
}

func (e *GOAPEvaluator) validateForward(path []*GOAPAction, start *GOAPWorldState, main *GOAPGoal) bool {
	ws := start.copyOf()
	for _, action := range path {
		if !e.presFulfilled(action, ws) {
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

func (e *GOAPEvaluator) prepend(
	start *GOAPWorldState,
	action *GOAPAction,
	here *GOAPPQueueItem) (extended *GOAPPQueueItem, useful bool) {
	/*

		prependSlice := func(a *GOAPAction, path []*GOAPAction) []*GOAPAction {
			prepended := make([]*GOAPAction, len(path))
			copy(prepended, path)
			prepended = append([]*GOAPAction{a}, path...)
			return prepended
		}
		// TODO: presRemaining should not be a name map, it should be actually a list...
		// each action has to have its pres fulfilled, and there may be multiple actions
		// with the same name -- doesn't necessarily mean their pres are valid...
		// say we have start with 0 booze and do:
		//
		// getBooze -> drinkBooze -> drinkBooze.
		//
		// the second one doesn't have its pres fulfilled yet
		if !action.affectsAnUnfulfilledVar(here.remaining, here.presRemaining) {
			debugGOAPPrintf("Action %s doesn't help with any goal -- skipping", action.name)
			return nil, false
		}
		debugGOAPPrintf("    starting state:")
		debugGOAPPrintWorldState(start)
		debugGOAPPrintf("    here.remaining:")
		debugGOAPPrintGoal(here.remaining)
		debugGOAPPrintf("    here.presRemaining:")
		for _, pre := range here.presRemaining {
			debugGOAPPrintGoal(pre)
		}
		newPath := prependSlice(action, here.path)
		debugGOAPPrintf("    --- evaluating path: %s", GOAPPathToString(newPath))
		// note: we start from nothing because, say the end goal wants fridgeClosed = 0
		// and we start with it closed, this will prematurely *fulfill* the fridgeClosed goal
		// of {hasFood: 1, fridgeClosed: 1}
		// instead, the start state is only used when checking pres for an otherwise fulfilled goal
		nothing := NewGOAPWorldState(nil)
		newResult := e.applyPath(newPath, nothing)

		debugGOAPPrintf("    --- result of chain with %s prepended:", action.name)
		debugGOAPPrintWorldState(newResult)

		beforeResult := here.endState
		debugGOAPPrintf("    --- result of chain without:")
		debugGOAPPrintWorldState(beforeResult)

		singleResult := e.applyAction(action, start)
		debugGOAPPrintf("    --- single result of %s", action.name)
		debugGOAPPrintWorldState(singleResult)

		// does newResult fulfill goal better than beforeResult?
		debugGOAPPrintf("    --- remaining of goal?")
		closerToGoal, remaining := here.remaining.stateCloserInSomeVar(newResult, beforeResult)
		// does newResult fulfill any goal?
		remainingNothing, _ := here.remaining.remaining(nothing)
		betterThanNothing := len(remainingNothing.fulfilled) > 0
		// calculate pres remaining
		// does singleResult fulfill any of presRemaining better than nil?
		helpsWithAPre := false
		afterPresRemaining := make(map[string]*GOAPGoal)
		for preActionName, pre := range here.presRemaining {
			if len(pre.goals) == 0 {
				continue // pre already fulfilled
			}
			betterThanNothingForPre, preRemaining := pre.stateCloserInSomeVar(singleResult, nothing)
			debugGOAPPrintf("    action %s better than nothing for pre of %s? %t", action.name, preActionName, betterThanNothingForPre)
			helpsWithAPre = helpsWithAPre || betterThanNothingForPre
			afterPresRemaining[preActionName] = preRemaining
		}
		afterPresRemaining[action.name] = action.pres

		debugGOAPPrintf("    === remaining:")
		debugGOAPPrintGoal(remaining)
		debugGOAPPrintf("    === afterPresRemaining:")
		for _, pre := range afterPresRemaining {
			debugGOAPPrintGoal(pre)
		}

		// if remaining.goals == 0, see if we can fulfill all remaining pres with start state
		if len(remaining.goals) == 0 {
			presAfterStart := make(map[string]*GOAPGoal)
			nPresFulfilled := 0
			for preActionName, pre := range afterPresRemaining {
				preRemaining, _ := pre.remaining(start)
				if len(pre.goals) == 0 || len(preRemaining.goals) == 0 {
					debugGOAPPrintf("        pre:%s is fulfilled.", preActionName)
					nPresFulfilled++
					presAfterStart[preActionName] = preRemaining
				}
			}
			if nPresFulfilled == len(afterPresRemaining) {
				debugGOAPPrintf("!!!!!!! start state fulfills all remaining pres!")
				afterPresRemaining = presAfterStart
			}
		}

		// if this action is good for something,
		if closerToGoal || betterThanNothing || helpsWithAPre {
			debugGOAPPrintf("    remaining for %s:", GOAPPathToString(newPath))
			debugGOAPPrintGoal(remaining)
			debugGOAPPrintf("    presRemaining for %s:", GOAPPathToString(newPath))
			if len(afterPresRemaining) == 0 {
				debugGOAPPrintf("    none")
			}
			for actionName, g := range afterPresRemaining {
				debugGOAPPrintf("%s now wants:", actionName)
				debugGOAPPrintGoal(g)
			}
			// add path cost
			cost := 0
			for _, action := range newPath {
				switch action.cost.(type) {
				case int:
					cost += action.cost.(int)
				case func() int:
					cost += action.cost.(func() int)()
				}
			}
			// calculate n goals unfulfilled of main goal and pres
			nUnfulfilledMain := len(remaining.goals)
			nUnfulfilledPres := 0
			for _, pre := range afterPresRemaining {
				nUnfulfilledPres += len(pre.goals)
			}
			debugGOAPPrintf("    nUnfulfilledMain: %d", nUnfulfilledMain)
			debugGOAPPrintf("    nUnfulfilledPres: %d", nUnfulfilledPres)

			// add heuristic
			cost += nUnfulfilledMain + nUnfulfilledPres

			extended = &GOAPPQueueItem{
				path:          newPath,
				presRemaining: afterPresRemaining,
				remaining:     remaining,
				endState:      newResult,
				nUnfulfilled:  nUnfulfilledMain + nUnfulfilledPres,
				cost:          cost,
				index:         -1, // going to be set by Push()
			}
			return extended, true
		}
		return nil, false
	*/
	return nil, false
}
