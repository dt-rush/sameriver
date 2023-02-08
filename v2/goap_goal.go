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
func (i *NumericInterval) diff(x float64) float64 {
	if x >= i.a && x <= i.b {
		return 0
	} else if x < i.a {
		return i.a - x
	} else {
		return i.b - x
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

func (g *GOAPGoal) remaining(ws *GOAPWorldState) (remaining *GOAPGoal, diffs map[string]float64) {
	remaining = NewGOAPGoal(nil)
	//    map[drunk:0xc0001080b0 templeEntered:0xc0001080c0]
	diffs = make(map[string]float64)
	debugGOAPPrintf("            checking remaining for goal: %v", g.goals)
	debugGOAPPrintf("            ws: %v", ws.vals)
	for varName, interval := range g.goals {

		if stateVal, ok := ws.vals[varName]; ok {
			diff := interval.diff(float64(stateVal))
			diffs[varName] = diff
			if diff == 0 {
				remaining.fulfilled[varName] = true
			} else {
				remaining.goals[varName] = interval
			}
		} else {
			// varName not in worldstate - diff is infinite and goal is unchanged for this var
			diffs[varName] = math.Inf(+1)
			remaining.goals[varName] = g.goals[varName]
		}
	}
	return remaining, diffs
}

func (g *GOAPGoal) stateCloserInSomeVar(after, before *GOAPWorldState) (closer bool, afterRemaining *GOAPGoal) {
	debugGOAPPrintf("*** stateCloserInSomeVar()")
	debugGOAPPrintf("*** goal: %v", g.goals)
	debugGOAPPrintf("***    before")
	debugGOAPPrintWorldState(before)
	debugGOAPPrintf("***    after")
	debugGOAPPrintWorldState(after)
	afterRemaining, afterDiffs := g.remaining(after)
	debugGOAPPrintf("***    afterRemaining:")
	debugGOAPPrintGoal(afterRemaining)
	_, beforeDiffs := g.remaining(before)
	debugGOAPPrintf("*** beforeDiffs: %v", beforeDiffs)
	debugGOAPPrintf("*** afterDiffs: %v", afterDiffs)
	for varName, _ := range afterDiffs {
		if math.Abs(afterDiffs[varName]) < math.Abs(beforeDiffs[varName]) {
			debugGOAPPrintf("****************** closer!")
			return true, afterRemaining
		}
	}
	debugGOAPPrintf("****************** not closer")
	return false, afterRemaining
}

func (g *GOAPGoal) stateAssuresInSomeVar(state *GOAPWorldState) (assures bool) {
	debugGOAPPrintf("*** stateAssuresInSomeVar()")
	_, diffs := g.remaining(state)
	for varName, _ := range diffs {
		if diffs[varName] == 0 {
			debugGOAPPrintf("****************** assures!")
			return true
		}
	}
	debugGOAPPrintf("****************** doesn't assure")
	return false
}
