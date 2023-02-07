package sameriver

import (
	"bytes"
	"os"
)

func debugGOAPPrintf(s string, args ...any) {
	if val, ok := os.LookupEnv("DEBUG_GOAP"); ok && val == "true" {
		Logger.Printf(s, args...)
	}
}

func GOAPPlanToString(plan []*GOAPAction) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, action := range plan {
		buf.WriteString(action.name)
		if i != len(plan)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

type GOAPPlanner struct {
	e    *Entity
	eval *GOAPEvaluator
}

type GOAPPathAndRemaining struct {
	path      []*GOAPAction
	remaining *GOAPGoal
}

func NewGOAPPlanner(e *Entity) *GOAPPlanner {
	return &GOAPPlanner{
		e:    e,
		eval: NewGOAPEvaluator(),
	}
}

func (p *GOAPPlanner) deepen(
	start *GOAPWorldState,
	path []*GOAPAction,
	goal *GOAPGoal) (newPaths []GOAPPathAndRemaining, solutions [][]*GOAPAction) {

	newPaths = make([]GOAPPathAndRemaining, 0)
	solutions = make([][]*GOAPAction, 0)

	prepend := func(a *GOAPAction, path []*GOAPAction) []*GOAPAction {
		prepended := make([]*GOAPAction, len(path))
		copy(prepended, path)
		prepended = append([]*GOAPAction{a}, path...)
		return prepended
	}

	pathResult := p.eval.applyPath(path, start)
	for _, action := range p.eval.actions.set {
		newPath := prepend(action, path)
		newResult := p.eval.applyPath(newPath, start)
		closer, remaining := goal.stateCloserInSomeVar(newResult, pathResult)
		if closer {
			if len(remaining.goals) == 0 {
				solutions = append(solutions, newPath)
			} else {
				newPaths = append(newPaths, GOAPPathAndRemaining{
					path:      newPath,
					remaining: remaining,
				})
			}
		}
	}
	return newPaths, solutions
}

/*
 TODO: in the result of deepen, we have the goalRemaining *without* the
 pres of the action. In merging two goals, if they don't coincide in varNames
 it's easy as pie, just union the maps. But if a varName coincides, we need
 the goal func to respond to the *intersection* of the spec-vals, and
 possibly return an error if their intersection is the null set.
 for example, if the goalRemaining is {drunk >= 3} and the pre is {drunk = 0},
 we can't make a proper goal. But if we want to merge {drunk >= 3} and {drunk >= 5,
 then the result is {drunk >= 5}. The logic depends on the operators

 we need to consider what will happen if our goal is itself incoherent from
 the start, and fail early if so rather than explore a bunch of partial paths.

 consider where we want the end goal

	 {drunk >= 3, admittedToTemple = 1}.

 the winning path of prependings is:

 [purifyOneself] (want: drunk >= 3, hasBooze = 0)
 [dropAllBooze, purifyOneself] (want: drunk >= 3)
 ...
 [drink, drink, drink, dropAllBooze, purifyOnself]

*/

/*
func (p *GOAPPlanner) Plans(
	world *GOAPWorldState,
	goal *GOAPWorldState,
	maxIter int) [][]GOAPAction {

	results := make([][]GOAPAction, 0)

	pq := GOAPPriorityQueue{}

	traverseFulfillers := func(path []GOAPAction, want *GOAPWorldState) {
		debugGOAPPrintf("--------------------------")
		Logger.Println("backtrack path: ")
		debugGOAPPrintf(GOAPPlanToString(path))
		debugGOAPPrintf("traversing fulfillers of want: %v", want)
		fulfillers := p.actions.thoseThatHelpFulfill(want, path)
		debugGOAPPrintf("fulfillers:")
		for _, fulfiller := range fulfillers.set {
			debugGOAPPrintf("    %v", fulfiller.name)
		}
		for _, action := range fulfillers.set {
			// TODO: what are we doing here?
			prependedPath := make([]GOAPAction, len(path))
			copy(prependedPath, path)
			prependedPath = append([]GOAPAction{action}, path...)
			debugGOAPPrintf("        Unfulfilled by %s:", action.name)
			unfulfilled := want.unfulfilledBy(path)
			debugGOAPPrintf("        %v", unfulfilled)
			want := unfulfilled.mergeActionPres(action)
			pq.Push(&GOAPPQueueItem{
				path: append([]GOAPAction{action}, path...),
				want: want,
			})
		}
		debugGOAPPrintf("--------------------------")
	}
	traverseFulfillers([]GOAPAction{}, goal)
	for pq.Len() > 0 {
		here := pq.Pop().(*GOAPPQueueItem)
		if len(here.want.Vals) == 0 || world.fulfills(here.want) {
			if p.validateForward(world, here.path, goal) {
				results = append(results, here.path)
				if len(results) == 2 {
					return results
				}
			}
		} else {
			traverseFulfillers(here.path, here.want)
		}
	}

	return results
}

func (p *GOAPPlanner) validateForward(
	world *GOAPWorldState,
	path []GOAPAction,
	goal *GOAPWorldState) bool {

	for _, action := range path {
		if !action.presFulfilled(world) {
			return false
		}
		world = world.applyAction(action)
	}

	return world.fulfills(goal)
}
*/
