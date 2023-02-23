package sameriver

import (
	"fmt"
	"strings"
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
	goal *GOAPGoal,
	pathsSeen map[string]bool) {

	debugGOAPPrintf("traverse--------------------------")
	debugGOAPPrintf(color.InRedOverGray("remaining:"))
	debugGOAPPrintGoalRemainingSurface(here.path.remainings)

	debugGOAPPrintf("%d possible actions", len(p.eval.actions.set))

	for _, action := range p.eval.actions.set {
		debugGOAPPrintf("[ ] Considering action %s", action.DisplayName())
		// determine if action is good to insert anywhere
		// consider, surface: [Apre, Bpre, Main]
		// consider inserting at 0 means fulfilling Apre
		for i, g := range here.path.remainings.surface {
			if g.nUnfulfilled == 0 {
				continue
			}
			if DEBUG_GOAP {
				debugGOAPPrintf(color.InGreenOverGray(
					fmt.Sprintf("checking if %s can be inserted at %d to satisfy %v",
						action.DisplayName(), i, g.goalLeft)))
			}
			scale, helpful := p.eval.actionHelpsToInsert(start, here.path, i, action)
			if helpful {
				debugGOAPPrintf("[X] %s helpful!", action.DisplayName())
				var toInsert *GOAPAction
				if scale > 1 {
					toInsert = action.Parametrized(scale)
				} else {
					toInsert = action
				}
				newPath := here.path.inserted(toInsert, i)
				pathStr := newPath.String()
				if _, ok := pathsSeen[pathStr]; ok {
					continue
				} else {
					pathsSeen[pathStr] = true
				}
				p.eval.computeRemainingsOfPath(newPath, start, goal)
				if DEBUG_GOAP {
					msg := fmt.Sprintf("{} - {} - {}    new path: %s     (cost %d)",
						GOAPPathToString(newPath), newPath.cost)
					debugGOAPPrintf(color.InWhiteOverCyan(strings.Repeat(" ", len(msg))))
					debugGOAPPrintf(color.InWhiteOverCyan(msg))
					debugGOAPPrintf(color.InWhiteOverCyan(strings.Repeat(" ", len(msg))))
				}
				pq.Push(&GOAPPQueueItem{path: newPath})
			} else {
				debugGOAPPrintf("[_] %s not helpful", action.DisplayName())
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

	// used to keep track of which paths we've already seen since there's multiple ways to
	// reach a path in the insertion-based logic we use
	pathsSeen := make(map[string]bool)

	// used for the search
	pq := &GOAPPriorityQueue{}

	rootPath := NewGOAPPath(nil)
	p.eval.computeRemainingsOfPath(rootPath, start, goal)
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
		debugGOAPPrintf(color.InWhiteOverBlue(color.InBold(GOAPPathToString(here.path))))
		debugGOAPPrintf(color.InRedOverGray(fmt.Sprintf("(%d unfulfilled)", here.path.remainings.NUnfulfilled())))

		if here.path.remainings.NUnfulfilled() == 0 {
			ok := p.eval.validateForward(here.path, start, goal)
			if !ok {
				debugGOAPPrintf(">>>>>>> potential solution rejected")
				continue
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
			p.traverseFulfillers(pq, start, here, goal, pathsSeen)
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
