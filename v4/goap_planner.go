package sameriver

import (
	"container/heap"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/TwiN/go-color"
)

type GOAPPlanner struct {
	e *Entity

	// the functions selecting an entity
	boundSelectors map[string]func(*Entity) bool
	modalVals      map[string]GOAPModalVal
	actions        *GOAPActionSet
	// map of [varName](Set (in the sense of a map->bool) of actions that affect that var)
	varActions map[string](map[*GOAPAction]bool)
}

func NewGOAPPlanner(e *Entity) *GOAPPlanner {
	return &GOAPPlanner{
		e:          e,
		modalVals:  make(map[string]GOAPModalVal),
		actions:    NewGOAPActionSet(),
		varActions: make(map[string](map[*GOAPAction]bool)),
	}
}

func (p *GOAPPlanner) bindEntities(nodes []string, ws *GOAPWorldState, start bool) {
	var pos *Vec2D
	if start {
		pos = p.e.GetVec2D(POSITION)
	} else {
		pos = ws.GetModal(p.e, POSITION).(*Vec2D)
	}
	box := p.e.GetVec2D(BOX)
	world := p.e.World
	// don't overwrite one that we already have (inherit bindings from the earliest point
	// in the chain that they are set)
	for _, node := range nodes {
		if _, ok := ws.ModalEntities[node]; !ok {
			ws.ModalEntities[node] = world.ClosestEntityFilter(*pos, *box, p.boundSelectors[node])
		}
	}
}

func (p *GOAPPlanner) BindEntitySelectors(selectors map[string]func(*Entity) bool) {
	p.boundSelectors = make(map[string]func(*Entity) bool)
	for k, v := range selectors {
		p.boundSelectors[k] = v
	}
}

func (p *GOAPPlanner) AddModalVals(vals ...GOAPModalVal) {
	for _, val := range vals {
		p.modalVals[val.name] = val
	}
}

func (p *GOAPPlanner) checkModalInto(varName string, ws *GOAPWorldState) {
	if _, ok := p.modalVals[varName]; ok {
		ws.vals[varName] = p.modalVals[varName].check(ws)
	}
}

func (p *GOAPPlanner) actionAffectsVar(action *GOAPAction, varName string) {
	if _, mapExists := p.varActions[varName]; !mapExists {
		p.varActions[varName] = make(map[*GOAPAction]bool)
	}
	p.varActions[varName][action] = true
}

func (p *GOAPPlanner) AddActions(actions ...*GOAPAction) {
	for _, action := range actions {
		logGOAPDebug("[][][] adding action %s", action.DisplayName())
		p.actions.Add(action)
		// link up modal setters for effs matching modal varnames
		for varName := range action.effs {
			p.actionAffectsVar(action, varName)
			if modal, ok := p.modalVals[varName]; ok {
				logGOAPDebug("[][][]     adding modal setter for %s", varName)
				action.effModalSetters[varName] = modal.effModalSet
			}
		}
		// link up modal checks for pres matching modal varnames
		for _, tg := range action.pres.temporalGoals {
			for varName := range tg.vars {
				if modal, ok := p.modalVals[varName]; ok {
					action.preModalChecks[varName] = modal.check
				}
			}
		}
	}
}

func (p *GOAPPlanner) applyActionBasic(
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

func (p *GOAPPlanner) applyActionModal(a *GOAPAction, ws *GOAPWorldState) (newWS *GOAPWorldState, cost float64) {

	// calculate cost of this action as cost to modally move position here + action.cost
	beforePos := ws.GetModal(p.e, POSITION).(*Vec2D)
	// find nearest matching entity
	p.bindEntities(append(a.otherNodes, a.Node), ws, false)
	node := ws.ModalEntities[a.Node]
	nodePos := ws.GetModal(node, POSITION).(*Vec2D)
	distToGetHere := nodePos.Sub(*beforePos).Magnitude()
	// now we are at it
	nowPos := *nodePos
	ws.SetModal(p.e, POSITION, &nowPos)
	logGOAPDebug("        distance to get to node for action %s: %f", a.Name, distToGetHere)
	cost = distToGetHere
	switch a.cost.(type) {
	case int:
		cost += float64(a.Count * a.cost.(int))
	case func() int:
		cost += float64(a.Count * a.cost.(func() int)())
	}

	// apply the modal state changes involved in this actions effs
	newWS = p.applyActionBasic(a, ws, false)
	for varName, eff := range a.effs {
		op := a.ops[varName]
		x := ws.vals[varName]
		logGOAPDebug("    %s        applying %s::%d x %s%s%d(%d) ; = %d",
			color.InPurpleOverWhite(" >>>modal "),
			a.DisplayName(), a.Count, varName, op, eff.val, x,
			eff.f(a.Count, x))
		// do modal set
		if setter, ok := a.effModalSetters[varName]; ok {
			setter(newWS, op, a.Count*eff.val)
		}
	}
	if DEBUG_GOAP {
		logGOAPDebug(color.InBlueOverWhite(fmt.Sprintf("            ws after action: %v", newWS.vals)))
	}

	// re-check any modal vals
	for varName := range newWS.vals {
		if modalVal, ok := p.modalVals[varName]; ok {
			logGOAPDebug("              re-checking modal val %s", varName)
			newWS.vals[varName] = modalVal.check(newWS)
		}
	}

	// we are still at the node after the action is done
	if a.travelWithNode {
		nowPos := *(newWS.GetModal(node, POSITION).(*Vec2D))
		newWS.SetModal(p.e, POSITION, &nowPos)
	}

	afterPos := newWS.GetModal(p.e, POSITION).(*Vec2D)
	distTravelled := nowPos.Sub(*afterPos).Magnitude()
	logGOAPDebug("        distance travelled during action %s: %f", a.Name, distTravelled)
	cost += distTravelled

	if DEBUG_GOAP {
		logGOAPDebug(color.InPurpleOverWhite(fmt.Sprintf("            ws after re-checking modal vals: %v", newWS.vals)))
	}

	return newWS, cost
}

func (p *GOAPPlanner) computeCostAndRemainingsOfPath(path *GOAPPath, start *GOAPWorldState, main *GOAPTemporalGoal) {
	ws := start.CopyOf()
	// one []*GOAPGoalRemaining for each action pre + 1 for main
	surfaceLen := len(path.path) + 1
	surface := newGOAPGoalRemainingSurface(surfaceLen)
	// create the storage space for statesAlong
	// consider, a path [A B C] will have 4 states: [start, postA, postB, postC (end)]
	path.statesAlong = make([]*GOAPWorldState, len(path.path)+1)
	path.statesAlong[0] = ws
	totalCost := 0.0
	for i, action := range path.path {
		for _, tg := range action.pres.temporalGoals {
			surface.surface[i] = append(
				surface.surface[i],
				tg.remaining(ws))
		}
		var cost float64
		ws, cost = p.applyActionModal(action, ws)
		totalCost += cost
		path.statesAlong[i+1] = ws
	}
	for _, tg := range main.temporalGoals {
		surface.surface[surfaceLen-1] = append(
			surface.surface[surfaceLen-1],
			tg.remaining(ws))
	}
	path.remainings = surface
	// total distance travelled and inherent effort of actions modally
	// we offset by +1 so that when NUnfulfilled == 0, we have a solution,
	// and the path cost is just the distance + inherent effort cost computed modally
	path.cost = totalCost * float64(path.remainings.NUnfulfilled()+1)
	logGOAPDebug("  --- ws after path: %v", ws.vals)
}

func (p *GOAPPlanner) presFulfilled(a *GOAPAction, ws *GOAPWorldState) bool {
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

func (p *GOAPPlanner) validateForward(path *GOAPPath, start *GOAPWorldState, main *GOAPTemporalGoal) bool {

	ws := start.CopyOf()
	for _, action := range path.path {
		if len(action.pres.temporalGoals) > 0 && !p.presFulfilled(action, ws) {
			logGOAPDebug(">>>>>>> in validateForward, %s was not fulfilled", action.DisplayName())
			return false
		}
		ws, _ = p.applyActionModal(action, ws)
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

func (p *GOAPPlanner) actionHelpsToInsert(
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
			affectors := p.varActions[varName]
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

func (p *GOAPPlanner) setPositionInStartModalIfNotDefined(start *GOAPWorldState) {
	start.SetModal(p.e, POSITION, p.e.GetVec2D(POSITION))
}

func (p *GOAPPlanner) setVarInStartIfNotDefined(start *GOAPWorldState, varName string) {
	logGOAPDebug("[ ] setVarInStartIfNotDefined(%s)", varName)
	if _, already := start.vals[varName]; !already {
		if modal, isModal := p.modalVals[varName]; isModal {
			p.bindEntities(modal.nodes, start, true)
			p.checkModalInto(varName, start)
			logGOAPDebug(color.InPurple(fmt.Sprintf("[ ] start.modal[%s] %d", varName, start.vals[varName])))
		} else {
			// NOTE: vars that don't have modal check default to 0
			logGOAPDebug(color.InYellow(fmt.Sprintf("[ ] %s not defined in GOAP start state, and no modal check exists. Defaulting to 0.", varName)))
			start.vals[varName] = 0
		}
	}
}

/*
The part of the GOAP backtracking A* algorithm where we traverse through the
possible actions that can be inserted into the plan generated so far. The
purpose of this function is to explore different combinations of actions that
can satisfy the remaining goals or preconditions.

The algorithm iterates through the unfulfilled goals at each stage (indexed by
'i') and their respective regions (indexed by 'regionIx'). It then checks if
any action can be inserted at the relevant position to satisfy the unfulfilled
goal. We only iterate as candidates those actions which effect some var in a
goal.

For each action, the algorithm evaluates if the action is helpful to be
inserted at the current position in the plan. If the action is helpful, it's
inserted into the plan, and the resulting new path is pushed into the priority
queue (pq) for further exploration. The paths seen so far are stored in the '
pathsSeen' map to avoid revisiting already explored paths.
*/
func (p *GOAPPlanner) traverseFulfillers(
	pq *GOAPPriorityQueue,
	start *GOAPWorldState,
	here *GOAPPQueueItem,
	goal *GOAPTemporalGoal,
	pathsSeen map[string]bool) {

	if DEBUG_GOAP {
		logGOAPDebug("traverse--------------------------")
		logGOAPDebug(color.InRedOverGray("remaining:"))
		debugGOAPPrintGoalRemainingSurface(here.path.remainings)
		logGOAPDebug("%d possible actions", len(p.actions.set))
		logGOAPDebug("regionOffsets: %v", here.path.regionOffsets)
	}

	// determine if action is good to insert anywhere
	// consider
	// path: [A, C]
	// A.pre = [q], A fulfills s
	// C.pre = [s t], C fulfills u
	// remainings.surface: [[q] [s t] [u]]

	// iterate path with index i
	// and iterate temporal goal of surface (ex iterate inside [s t]) with index regionIx

	for i, tgs := range here.path.remainings.surface {
		if here.path.remainings.nUnfulfilledAtIx(i) == 0 {
			continue
		}
		var parent *GOAPAction
		if i == len(here.path.remainings.surface)-1 {
			logGOAPDebug("    nil parent (satisfying main goal)")
			parent = nil
		} else {
			logGOAPDebug("    parent: %s", here.path.path[i].Name)
			parent = here.path.path[i]
		}
		for regionIx, tg := range tgs {
			logGOAPDebug("  surface[i:%d][regionIx:%d]", i, regionIx)
			logGOAPDebug("        |")
			logGOAPDebug("        |")
			for varName := range tg.goalLeft {
				for action := range p.varActions[varName] {
					// can't self-append (guard against considering i == len(path) (that's the main goal)
					if i < len(here.path.path) && here.path.path[i].Name == action.Name {
						continue
					}
					logGOAPDebug("       ...")
					logGOAPDebug("        |")
					logGOAPDebug("        └>varName: %s", varName)

					insertionIx := i + here.path.regionOffsets[i][regionIx]
					logGOAPDebug("    insertionIx: %d", insertionIx)

					if DEBUG_GOAP {
						var toSatisfyMsg string
						if i == len(here.path.remainings.surface)-1 {
							toSatisfyMsg = "main goal"
						} else {
							toSatisfyMsg = fmt.Sprintf("pre of %s", here.path.path[i].Name)
						}
						logGOAPDebug(color.InGreenOverGray(
							fmt.Sprintf("checking if %s can be inserted at %d to satisfy %s",
								action.DisplayName(), i, toSatisfyMsg)))
					}
					scale, helpful := p.actionHelpsToInsert(
						start,
						here.path,
						insertionIx,
						tg,
						action)
					if helpful {
						if DEBUG_GOAP {
							logGOAPDebug("[X] %s helpful!", action.DisplayName())
						}
						// construct the path (with parametrised action)
						var toInsert *GOAPAction
						if scale > 1 {
							toInsert = action.Parametrized(scale) // yields a copy
						} else {
							toInsert = action.CopyOf()
						}
						toInsert = toInsert.ChildOf(parent)
						newPath := here.path.inserted(toInsert, insertionIx, regionIx)
						// guard against visit to already-seen path
						pathStr := newPath.String()
						if _, ok := pathsSeen[pathStr]; ok {
							logGOAPDebug(color.InBold(color.InWhiteOverCyan("path seen already")))
							continue
						} else {
							pathsSeen[pathStr] = true
						}
						// compute remainings of path
						p.computeCostAndRemainingsOfPath(newPath, start, goal)
						// check any modal vals in the pres of action that aren't already
						// in the start state
						for _, tg := range toInsert.pres.temporalGoals {
							for varName := range tg.vars {
								p.setVarInStartIfNotDefined(start, varName)
							}
						}
						if DEBUG_GOAP {
							msg := fmt.Sprintf("{} - {} - {}    new path: %s     (cost %.2f)",
								GOAPPathToString(newPath), newPath.cost)
							logGOAPDebug(color.InWhiteOverCyan(strings.Repeat(" ", len(msg))))
							logGOAPDebug(color.InWhiteOverCyan(msg))
							logGOAPDebug(color.InWhiteOverCyan(strings.Repeat(" ", len(msg))))
						}
						heap.Push(pq, &GOAPPQueueItem{path: newPath})
					} else {
						logGOAPDebug("[_] %s not helpful", action.DisplayName())
					}
				}
			}
		}
	}
	logGOAPDebug("--------------------------/traverse")
}

func (p *GOAPPlanner) Plan(
	start *GOAPWorldState,
	goalSpec any,
	maxIter int) (solution *GOAPPath, ok bool) {

	// we may be writing to this with modal vals as we explore and don't want
	// to pollute the caller's state object
	start = start.CopyOf()
	start.w = p.e.World
	p.setPositionInStartModalIfNotDefined(start)

	logGOAPDebug("Planning...")

	// convert goal spec into GOAPTemporalGoal
	goal := NewGOAPTemporalGoal(goalSpec)

	// populate start state with any modal vals at start
	for _, tg := range goal.temporalGoals {
		for varName := range tg.vars {
			p.setVarInStartIfNotDefined(start, varName)
		}
	}

	// used to return the solution with lowest cost among solutions found
	resultPq := &GOAPPriorityQueue{}
	heap.Init(resultPq)

	// used to keep track of which paths we've already seen since there's multiple ways to
	// reach a path in the insertion-based logic we use
	pathsSeen := make(map[string]bool)

	// used for the search
	pq := &GOAPPriorityQueue{}
	heap.Init(pq)

	rootPath := NewGOAPPath(nil)
	p.computeCostAndRemainingsOfPath(rootPath, start, goal)
	rootPath.regionOffsets[0] = make([]int, len(goal.temporalGoals))
	backtrackRoot := &GOAPPQueueItem{
		path:  rootPath,
		index: -1, // going to be set by Push()
	}
	heap.Push(pq, backtrackRoot)

	iter := 0
	// TODO: should we just pop out the *very first result*?
	// why wait for 2 or exhausting the pq?
	t0 := time.Now()
	// remove resultPq.Len() < 2 to exhaust space/max iter
	for iter < maxIter && pq.Len() > 0 {

		logGOAPDebug("=== iter ===")
		here := heap.Pop(pq).(*GOAPPQueueItem)
		if DEBUG_GOAP {
			logGOAPDebug(color.InRedOverGray("here:"))
			logGOAPDebug(color.InWhiteOverBlue(color.InBold(GOAPPathToString(here.path))))
			logGOAPDebug(color.InRedOverGray(fmt.Sprintf("(%d unfulfilled)",
				here.path.remainings.NUnfulfilled())))
		}

		if here.path.remainings.NUnfulfilled() == 0 {
			ok := p.validateForward(here.path, start, goal)
			if !ok {
				logGOAPDebug(">>>>>>> potential solution rejected")
				continue
			}

			if DEBUG_GOAP {
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(fmt.Sprintf("    SOLUTION: %s", GOAPPathToString(here.path)))))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(GOAPPathToString(here.path))))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(fmt.Sprintf("%d solutions found so far", resultPq.Len()+1))))
			}
			heap.Push(resultPq, here)
		} else {
			p.traverseFulfillers(pq, start, here, goal, pathsSeen)
			iter++
		}
	}

	dt := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	if iter >= maxIter {
		logGOAPDebug("Took %f ms to reach max iter (%d)", dt, iter)
		logGOAPDebug("================================ REACHED MAX ITER !!!")
	}
	if pq.Len() == 0 && resultPq.Len() == 0 {
		logGOAPDebug("Took %f ms to exhaust pq without solution (%d iterations)", dt, iter)
	}
	if resultPq.Len() > 0 {
		logGOAPDebug("Took %f ms to find %d solutions (%d iterations)", dt, resultPq.Len(), iter)
		if pq.Len() == 0 {
			logGOAPDebug("Exhausted pq")
		}
		results := make([]*GOAPPath, 0)
		for resultPq.Len() > 0 {
			pop := heap.Pop(resultPq).(*GOAPPQueueItem).path
			results = append(results, pop)
			logGOAPDebug("solution (cost %.3f): %s", pop.cost, color.InWhiteOverBlue(GOAPPathToString(pop)))
		}
		return results[0], true
	} else {
		return nil, false
	}
}
