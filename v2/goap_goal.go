package sameriver

import (
	"math"
	"strings"

	"github.com/dt-rush/sameriver/v2/utils"
)

type GOAPGoal struct {
	vars map[string]*utils.NumericInterval
}

func NewGOAPGoal(def map[string]int) *GOAPGoal {
	g := &GOAPGoal{
		vars: make(map[string]*utils.NumericInterval),
	}
	for spec, val := range def {
		split := strings.Split(spec, ",")
		varName, op := split[0], split[1]
		interval := utils.MakeNumericInterval(op, val)
		g.vars[varName] = interval
	}
	return g
}

func (g *GOAPGoal) copyOf() *GOAPGoal {
	result := &GOAPGoal{
		vars: make(map[string]*utils.NumericInterval),
	}
	for varName, interval := range g.vars {
		result.vars[varName] = interval
	}
	return result
}

func (g *GOAPGoal) remaining(ws *GOAPWorldState) (result *GOAPGoalRemaining) {
	result = &GOAPGoalRemaining{
		goal:  NewGOAPGoal(nil),
		diffs: make(map[string]float64),
	}
	debugGOAPPrintf("      -+- checking remaining for goal: %v", g.vars)
	debugGOAPPrintf("      -+-     ws: %v", ws.vals)
	for varName, interval := range g.vars {
		if stateVal, ok := ws.vals[varName]; ok {
			diff := interval.Diff(float64(stateVal))
			result.diffs[varName] = diff
			if diff != 0 {
				result.goal.vars[varName] = interval
			}
		} else {
			// varName not in worldstate - diff is infinite and goal is unchanged for this var
			result.diffs[varName] = math.Inf(+1)
			result.goal.vars[varName] = interval
		}
	}
	return result
}

/*
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
*/
