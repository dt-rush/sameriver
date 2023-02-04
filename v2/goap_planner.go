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
		// Logger.Println("------------------------")
		// Logger.Println("traversing to fulfill:")
		// Logger.Println(want.Vals)
		// Logger.Println("...")
		fulfillers := p.actions.thoseThatHelpFulfill(want)
		// Logger.Println("fulfillers:")
		// for _, action := range fulfillers.set {
		// 	Logger.Println(action.name)
		// }
		for _, action := range fulfillers.set {
			unfulfilled := want.unfulfilledBy(action)
			// Logger.Printf("unfulfilled by %s:", action.name)
			// Logger.Println(unfulfilled)
			want := unfulfilled.mergeActionPres(action)
			// Logger.Println("want:")
			// Logger.Println(want)
			// Logger.Printf("Pushing %s...", action.name)
			pq.Push(&GOAPPQueueItem{
				path: append([]GOAPAction{action}, path...),
				want: want,
			})
			// Logger.Println("---")
		}
	}
	traverseFulfillers([]GOAPAction{}, goal)
	for pq.Len() > 0 {
		here := pq.Pop().(*GOAPPQueueItem)
		// Logger.Println("------------")
		// Logger.Println("exploring path:")
		// Logger.Println(GOAPPlanToString(here.path))
		// Logger.Printf("%d wants.", len(here.want.Vals))
		// Logger.Println("wants:")
		// Logger.Println(here.want.Vals)
		if len(here.want.Vals) == 0 || world.fulfills(here.want) {
			// test chain forward
			// Logger.Println("found possible solution:")
			// Logger.Println(GOAPPlanToString(here.path))
			if p.validateForward(world, here.path, goal) {
				// Logger.Println("Valid solution!")
				results = append(results, here.path)
			}
		} else {
			traverseFulfillers(here.path, here.want)
		}
	}

	Logger.Println("==========")
	Logger.Println("VALID PLANS:")
	for _, plan := range results {
		Logger.Println(GOAPPlanToString(plan))
	}
	Logger.Println("==========")

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
