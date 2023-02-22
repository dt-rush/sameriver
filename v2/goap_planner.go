package sameriver

import (
	"fmt"
	"time"

	"github.com/TwiN/go-color"
)

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

func (p *GOAPPlanner) traverseFulfillers(
	pq *GOAPPriorityQueue,
	start *GOAPWorldState,
	here *GOAPPQueueItem,
	goal *GOAPGoal) {

	debugGOAPPrintf("traverse--------------------------")
	debugGOAPPrintf(color.InRedOverGray("remaining:"))
	debugGOAPPrintGoalRemainingSurface(here.path.remainings)

	debugGOAPPrintf("%d possible actions", len(p.eval.actions.set))

	for _, action := range p.eval.actions.set {
		debugGOAPPrintf("[ ] Considering action %s", action.name)
		frontier := make([]*GOAPPQueueItem, 0)
		if p.eval.actionMightHelp(start, action, here.path, GOAP_PATH_PREPEND) {
			helpfulItem := p.eval.tryPrepend(start, action, here.path, goal)
			if helpfulItem != nil {
				frontier = append(frontier, helpfulItem)
			}
		}
		if p.eval.actionMightHelp(start, action, here.path, GOAP_PATH_APPEND) {
			helpfulItem := p.eval.tryAppend(start, action, here.path, goal)
			if helpfulItem != nil {
				frontier = append(frontier, helpfulItem)
			}
		}
		if len(frontier) == 0 {
			debugGOAPPrintf("[_] %s not helpful", action.name)
		} else {
			debugGOAPPrintf("[X] %s helpful!", action.name)
			for _, item := range frontier {
				debugGOAPPrintf(
					color.Ize(color.Cyan,
						fmt.Sprintf("}-}-}-}-}-}-}-} new path: %s",
							GOAPPathToString(item.path))))
				pq.Push(item)
			}
		}
	}
	debugGOAPPrintf("--------------------------/traverse")
}

func (p *GOAPPlanner) Plan(
	start *GOAPWorldState,
	goal *GOAPGoal,
	maxIter int) (solution *GOAPPath, ok bool) {

	// populate start state with any modal vals at start
	p.eval.PopulateModalStartState(start)

	// used to return the solution with lowest cost among solutions found
	resultPq := &GOAPPriorityQueue{}

	// used for the search
	pq := &GOAPPriorityQueue{}

	rootPath := NewGOAPPath([]*GOAPAction{}, 0)
	p.eval.remainingsOfPath(rootPath, start, goal)
	backtrackRoot := &GOAPPQueueItem{
		path:  rootPath,
		index: -1, // going to be set by Push()
	}
	pq.Push(backtrackRoot)

	iter := 0
	// TODO: should we just pop out the *very first result*?
	// why wait for 2 or exhausting the pq?
	t0 := time.Now()
	for iter < maxIter && pq.Len() > 0 && resultPq.Len() < 2 {
		debugGOAPPrintf("=== iter ===")
		here := pq.Pop().(*GOAPPQueueItem)
		debugGOAPPrintf(color.InRedOverGray("here:"))
		debugGOAPPrintf(color.InWhiteOverBlack(color.InBold(GOAPPathToString(here.path))))
		debugGOAPPrintf(color.InRedOverGray(fmt.Sprintf("(%d unfulfilled)", here.path.remainings.nUnfulfilled)))

		if here.path.remainings.nUnfulfilled == 0 {
			ok := p.eval.validateForward(here.path, start, goal)
			if !ok {
				debugGOAPPrintf(">>>>>>> potential solution rejected")
			}

			debugGOAPPrintf(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
			debugGOAPPrintf(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
			debugGOAPPrintf(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
			debugGOAPPrintf(color.InGreenOverWhite(color.InBold(fmt.Sprintf("    SOLUTION: %s", GOAPPathToString(here.path)))))
			debugGOAPPrintf(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
			debugGOAPPrintf(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
			debugGOAPPrintf(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
			debugGOAPPrintf(color.InGreenOverWhite(color.InBold(GOAPPathToString(here.path))))
			resultPq.Push(here)
			debugGOAPPrintf(color.InGreenOverWhite(color.InBold(fmt.Sprintf("%d solutions found so far", resultPq.Len()))))
		} else {
			p.traverseFulfillers(pq, start, here, goal)
			iter++
		}
	}

	dt := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	if iter >= maxIter {
		debugGOAPPrintf("Took %f ms to reach max iter (%d)", dt, iter)
		debugGOAPPrintf("================================ REACHED MAX ITER !!!")
	}
	if pq.Len() == 0 && resultPq.Len() == 0 {
		debugGOAPPrintf("Took %f ms to exhaust pq without solution (%d iterations)", dt, iter)
	}
	if resultPq.Len() > 0 {
		debugGOAPPrintf("Took %f ms to find solution (%d iterations)", dt, iter)
		return resultPq.Pop().(*GOAPPQueueItem).path, true
	} else {
		return nil, false
	}
}
