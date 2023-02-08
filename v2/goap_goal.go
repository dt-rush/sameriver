package sameriver

import (
	"math"
	"strings"
)

type NumericInterval struct {
	a, b float64
}

// given an open interval [a, b], return the least amount needed
// to modify x such that it will be in bounds (0 if in bounds, + if below,
// - if above)
//
// NOTE: we have to use float64 for +,- Inf, which is float64 type... annoying
func (i *NumericInterval) diff(x float64) int {
	if x >= i.a && x <= i.b {
		return 0
	} else if x < i.a {
		return int(i.a - x)
	} else {
		return int(x - i.b)
	}
}

func (i *NumericInterval) merged(j *NumericInterval) (result *NumericInterval, ok bool) {
	e := math.Max(i.a, j.a)
	f := math.Min(i.b, j.b)
	if e <= f {
		return &NumericInterval{e, f}, true
	} else {
		return nil, false
	}
}

type GOAPGoal struct {
	goals     map[string]*NumericInterval
	fulfilled map[string]bool
}

func MakeNumericInterval(op string, val int) *NumericInterval {
	switch op {
	case "<":
		return &NumericInterval{math.Inf(-1), float64(val - 1)}
	case "<=":
		return &NumericInterval{math.Inf(-1), float64(val)}
	case "=":
		return &NumericInterval{float64(val), float64(val)}
	case ">=":
		return &NumericInterval{float64(val), math.Inf(+1)}
	case ">":
		return &NumericInterval{float64(val + 1), math.Inf(+1)}
		/*
			case ">;<":
				// TODO
			case ">=;<":
				// TODO
			case ">;<=":
				// TODO
			case ">=;<=":
				// TODO
		*/
	default:
		panic("Got undefined op in GOAPGoalFunc() [valid: >=,>,=,<,<=]")
	}
}

func NewGOAPGoal(def map[string]int) *GOAPGoal {
	g := &GOAPGoal{
		goals:     make(map[string]*NumericInterval),
		fulfilled: make(map[string]bool),
	}
	for spec, val := range def {
		split := strings.Split(spec, ",")
		varName, op := split[0], split[1]
		g.goals[varName] = MakeNumericInterval(op, val)
	}
	return g
}

func (g *GOAPGoal) remaining(ws *GOAPWorldState) (remaining *GOAPGoal, diffs map[string]int) {
	remaining = NewGOAPGoal(nil)
	diffs = make(map[string]int)
	for varName, interval := range g.goals {
		var diff int
		// NOTE: we assume a var not present in the state as having val = 0
		if stateVal, ok := ws.vals[varName]; ok {
			diff = interval.diff(float64(stateVal))
		} else {
			diff = interval.diff(0)
		}
		diffs[varName] = diff
		// goal fulfilled
		if diff == 0 {
			// NOTE: this will *entirely remove the goal related to this var*,
			// so if, in backtracking, we add an action that fulfills another
			// goal but modifies this var as well after it was already fulfilled,
			// we may end up - playing the whole chain - with an invalid solution.
			// Thus we should be sure to - upon finding a "satisfied chain",
			// check it by playing it forward and seeing if
			// the resulting GOAPWorldState's remaining with the true
			// original goal is actually all zero'd that is, this should be
			// done when `deepen` returns a path in `solutions`
			remaining.fulfilled[varName] = true
		} else {
			remaining.goals[varName] = interval
		}
	}
	return remaining, diffs
}

func (g *GOAPGoal) stateCloserInSomeVar(after, before *GOAPWorldState) (closer bool, afterRemaining *GOAPGoal) {
	afterRemaining, afterDiffs := g.remaining(after)
	_, beforeDiffs := g.remaining(before)
	debugGOAPPrintf("*** stateCloserInSomeVar()")
	for varName, _ := range beforeDiffs {
		debugGOAPPrintf("*** %s,before: %d", varName, intAbs(beforeDiffs[varName]))
		debugGOAPPrintf("*** %s,after:  %d", varName, intAbs(afterDiffs[varName]))
		if intAbs(afterDiffs[varName]) < intAbs(beforeDiffs[varName]) {
			return true, afterRemaining
		}
	}
	return false, afterRemaining
}

func (g *GOAPGoal) prependingMerge(newPathGoal *GOAPGoal) (result *GOAPGoal, valid bool) {
	result = NewGOAPGoal(nil)
	// first, copy this goal
	for varName, interval := range g.goals {
		result.goals[varName] = interval
	}
	for varName, _ := range g.fulfilled {
		result.fulfilled[varName] = true
	}
	// now........ the monstrous task
	for varNameB, intervalB := range newPathGoal.goals {
		nameAlreadyInGoals := false
		for varNameA, intervalA := range result.goals {
			if varNameB == varNameA {
				nameAlreadyInGoals = true
				newInterval, ok := intervalB.merged(intervalA)
				if !ok {
					debugGOAPPrintf("interval conflict! favouring new:")
					debugGOAPPrintf("%s: old:[%.0f, %.0f], new:[%.0f, %.0f]", varNameA, intervalA.a, intervalA.b, intervalB.a, intervalB.b)
					// NOTE: we favour the new interval in the case of a conflict,
					// since this is required to allow solutions where we have contradictory
					// goals at different parts of the path
					// (see TestGOAPPlannerPurifyOneself, where at one point, for purifyOneself,
					// we want booze = 0, but for drink, earlier in the chain, we want booze > 1
					result.goals[varNameA] = intervalB
				} else {
					result.goals[varNameA] = newInterval
				}
			}
		}
		if !nameAlreadyInGoals {
			result.goals[varNameB] = intervalB
			if _, ok := result.fulfilled[varNameB]; ok {
				delete(result.fulfilled, varNameB)
			}
		}
	}
	for varName, _ := range newPathGoal.fulfilled {
		result.fulfilled[varName] = true
		delete(result.goals, varName)
	}
	return result, true
}
