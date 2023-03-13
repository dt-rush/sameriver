package sameriver

import (
	"container/heap"
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
	goal *GOAPTemporalGoal,
	pathsSeen map[string]bool) {

	if DEBUG_GOAP() {
		logGOAPDebug("traverse--------------------------")
		logGOAPDebug(color.InRedOverGray("remaining:"))
		debugGOAPPrintGoalRemainingSurface(here.path.remainings)
		logGOAPDebug("%d possible actions", len(p.eval.actions.set))
		logGOAPDebug("regionOffsets: %v", here.path.regionOffsets)
	}

	// determine if action is good to insert anywhere
	// consider
	// path: [A, C]
	// A.pre = [q], A fulfills s
	// C.pre = [s t], C fulfills u
	// remainings surface: [[q] [s t] [u]]

	// iterate path with index i
	// and iterate temporal of surface (ex [s t]) with index regionIx

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
				for action := range p.eval.varActions[varName] {
					logGOAPDebug("       ...")
					logGOAPDebug("        |")
					logGOAPDebug("        └>varName: %s", varName)

					insertionIx := i + here.path.regionOffsets[i][regionIx]
					logGOAPDebug("    insertionIx: %d", insertionIx)

					if DEBUG_GOAP() {
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
					scale, helpful := p.eval.actionHelpsToInsert(
						start,
						here.path,
						insertionIx,
						tg,
						action)
					if helpful {
						if DEBUG_GOAP() {
							logGOAPDebug("[X] %s helpful!", action.DisplayName())
						}
						var toInsert *GOAPAction
						if scale > 1 {
							toInsert = action.Parametrized(scale) // yields a copy
						} else {
							toInsert = action.CopyOf()
						}
						toInsert = toInsert.ChildOf(parent)
						newPath := here.path.inserted(toInsert, insertionIx, regionIx)
						pathStr := newPath.String()
						if _, ok := pathsSeen[pathStr]; ok {
							continue
						} else {
							pathsSeen[pathStr] = true
						}
						p.eval.computeRemainingsOfPath(newPath, start, goal)
						if DEBUG_GOAP() {
							msg := fmt.Sprintf("{} - {} - {}    new path: %s     (cost %d)",
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

	logGOAPDebug("Planning...")

	// convert goal spec into GOAPTemporalGoal
	goal := NewGOAPTemporalGoal(goalSpec)

	// populate start state with any modal vals at start
	p.eval.PopulateModalStartState(start)

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
	p.eval.computeRemainingsOfPath(rootPath, start, goal)
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
	for iter < maxIter && pq.Len() > 0 && resultPq.Len() < 2 {
		logGOAPDebug("=== iter ===")
		here := heap.Pop(pq).(*GOAPPQueueItem)
		if DEBUG_GOAP() {
			logGOAPDebug(color.InRedOverGray("here:"))
			logGOAPDebug(color.InWhiteOverBlue(color.InBold(GOAPPathToString(here.path))))
			logGOAPDebug(color.InRedOverGray(fmt.Sprintf("(%d unfulfilled)",
				here.path.remainings.NUnfulfilled())))
		}

		if here.path.remainings.NUnfulfilled() == 0 {
			ok := p.eval.validateForward(here.path, start, goal)
			if !ok {
				logGOAPDebug(">>>>>>> potential solution rejected")
				continue
			}

			if DEBUG_GOAP() {
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
			logGOAPDebug("solution (cost %d): %s", pop.cost, color.InWhiteOverBlue(GOAPPathToString(pop)))
		}
		return results[0], true
	} else {
		return nil, false
	}
}
