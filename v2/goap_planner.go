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

func NewGOAPPlanner(e *Entity) *GOAPPlanner {
	return &GOAPPlanner{
		e:    e,
		eval: NewGOAPEvaluator(),
	}
}

/*
func (p *GOAPPlanner) Plans(
	world *GOAPWorldState,
	goal *GOAPWorldState) [][]GOAPAction {

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
