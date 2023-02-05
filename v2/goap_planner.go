package sameriver

import (
	"bytes"
)

type GOAPPlan []GOAPAction

func GOAPPlanToString(plan GOAPPlan) string {
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
	e       *Entity
	actions *GOAPActionSet
}

func NewGOAPPlanner(e *Entity) *GOAPPlanner {
	return &GOAPPlanner{
		e:       e,
		actions: NewGOAPActionSet(),
	}
}

func (p *GOAPPlanner) AddActions(actions ...GOAPAction) {
	p.actions.Add(actions...)
}

func (p *GOAPPlanner) Plans(
	world GOAPWorldState,
	goal GOAPWorldState) []GOAPPlan {

	results := make([]GOAPPlan, 0)

	pq := GOAPPriorityQueue{}
	traverseFulfillers := func(path []GOAPAction, want GOAPWorldState) {
		fulfillers := p.actions.thoseThatHelpFulfill(want)
		for _, action := range fulfillers.set {
			unfulfilled := want.unfulfilledBy(action)
			want := unfulfilled.mergeActionPres(action)
			pq.Push(&GOAPPQueueItem{
				path: append([]GOAPAction{action}, path...),
				want: want,
			})
		}
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
	world GOAPWorldState,
	path []GOAPAction,
	goal GOAPWorldState) bool {

	for _, action := range path {
		if !action.presFulfilled(world) {
			return false
		}
		world = world.applyAction(action)
	}

	return world.fulfills(goal)

}
