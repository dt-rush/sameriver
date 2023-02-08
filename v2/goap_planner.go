package sameriver

import (
	"bytes"
	"os"
	"strings"
	"time"
)

func debugGOAPPrintf(s string, args ...any) {
	if val, ok := os.LookupEnv("DEBUG_GOAP"); ok && val == "true" {
		Logger.Printf(s, args...)
	}
}

func GOAPPathToString(plan []*GOAPAction) string {
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

func debugGOAPPrintGoal(g *GOAPGoal) {
	if g == nil || (len(g.goals) == 0 && len(g.fulfilled) == 0) {
		debugGOAPPrintf("    nil")
		return
	}
	for spec, interval := range g.goals {
		split := strings.Split(spec, ",")
		varName := split[0]
		debugGOAPPrintf("    %s: [%.0f, %.0f]", varName, interval.a, interval.b)
	}
	for varName, _ := range g.fulfilled {
		debugGOAPPrintf("    fulfilled: %s", varName)
	}
}

func debugGOAPPrintWorldState(ws *GOAPWorldState) {
	if ws == nil || len(ws.vals) == 0 {
		debugGOAPPrintf("    nil")
		return
	}
	for name, val := range ws.vals {
		debugGOAPPrintf("    %s: %d", name, val)
	}
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

func (p *GOAPPlanner) deepen(
	start *GOAPWorldState,
	here *GOAPPQueueItem) (frontier []*GOAPPQueueItem) {

	debugGOAPPrintf("deepen-----------------")
	frontier = make([]*GOAPPQueueItem, 0)
	for _, action := range p.eval.actions.set {
		debugGOAPPrintf("    ------------------------------ considering prepending action %s", action.name)
		extended, useful := p.eval.prepend(start, action, here)
		debugGOAPPrintf("    --- useful? %t", useful)
		if !useful {
			debugGOAPPrintf("    --- was not a useful prepending action")
			continue
		}
		debugGOAPPrintf("    --- OK!")
		debugGOAPPrintf("    ---=== frontier expanded to: %s", GOAPPathToString(extended.path))
		frontier = append(frontier, extended)
	}
	debugGOAPPrintf("-----------------/deepen")
	return frontier
}

func (p *GOAPPlanner) traverseFulfillers(
	pq *GOAPPriorityQueue,
	start *GOAPWorldState,
	here *GOAPPQueueItem) {

	debugGOAPPrintf("traverse--------------------------")
	debugGOAPPrintf("backtrack path so far: ")
	debugGOAPPrintf(GOAPPathToString(here.path))

	frontier := p.deepen(start, here)
	debugGOAPPrintf("newPaths:")
	for _, x := range frontier {
		debugGOAPPrintf("---")
		debugGOAPPrintf("    %s", GOAPPathToString(x.path))
		pq.Push(x)
	}
	debugGOAPPrintf("--------------------------/traverse")
}

func (p *GOAPPlanner) Plan(
	start *GOAPWorldState,
	goal *GOAPGoal,
	maxIter int) (solution []*GOAPAction, ok bool) {

	// populate start state with any modal vals at start
	p.eval.populateModalStartState(start)

	resultPq := &GOAPPriorityQueue{}

	pq := &GOAPPriorityQueue{}

	backtrackRoot := &GOAPPQueueItem{
		path:          []*GOAPAction{},
		presRemaining: make(map[string]*GOAPGoal),
		remaining:     goal,
		nUnfulfilled:  len(goal.goals),
		endState:      NewGOAPWorldState(nil),
		cost:          0,
		index:         -1, // going to be set by Push()
	}
	p.traverseFulfillers(pq, start, backtrackRoot)

	iter := 0
	// TODO: should we just pop out the *very first result*?
	// why wait for 2 or exhausting the pq?
	t0 := time.Now()
	for iter < maxIter && pq.Len() > 0 && resultPq.Len() < 2 {
		debugGOAPPrintf("=== iter ===")
		here := pq.Pop().(*GOAPPQueueItem)
		debugGOAPPrintf("here:")
		debugGOAPPrintf(GOAPPathToString(here.path))
		debugGOAPPrintf("(%d unfulfilled)", here.nUnfulfilled)

		if here.nUnfulfilled == 0 {
			debugGOAPPrintf(">>>>>>>>>>>>>>>>>>>>>>")
			debugGOAPPrintf(">>>>>>>>>>>>>>>>>>>>>>")
			debugGOAPPrintf(">>>>>>>>>>>>>>>>>>>>>>")
			debugGOAPPrintf("    SOLUTION: %s", GOAPPathToString(here.path))
			debugGOAPPrintf(">>>>>>>>>>>>>>>>>>>>>>")
			debugGOAPPrintf(">>>>>>>>>>>>>>>>>>>>>>")
			debugGOAPPrintf(">>>>>>>>>>>>>>>>>>>>>>")
			resultPq.Push(here)
		} else {
			p.traverseFulfillers(pq, start, here)
			iter++
		}
	}

	dt := time.Since(t0).Milliseconds()
	if iter >= maxIter {
		debugGOAPPrintf("Took %d ms to reach max iter", dt)
		debugGOAPPrintf("================================ REACHED MAX ITER !!!")
	}
	if pq.Len() == 0 && resultPq.Len() == 0 {
		debugGOAPPrintf("Took %d ms to exhaust pq without solution", dt)
	}
	if resultPq.Len() > 0 {
		debugGOAPPrintf("Took %d ms to find solution", dt)
		if pq.Len() == 0 {
			debugGOAPPrintf("Even though >0 solutions were found, exhausted pq")
		}
		return resultPq.Pop().(*GOAPPQueueItem).path, true
	} else {
		return nil, false
	}
}

func (p *GOAPPlanner) validateForward(
	start *GOAPWorldState,
	path []*GOAPAction,
	goal *GOAPGoal) bool {

	world := start
	for _, action := range path {
		if !p.eval.presFulfilled(action, world) {
			return false
		}
		world = p.eval.applyAction(action, world)
	}

	remaining, _ := goal.remaining(world)
	return len(remaining.goals) == 0
}
