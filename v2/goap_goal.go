package sameriver

import (
	"math"
	"strings"

	"github.com/dt-rush/sameriver/v2/utils"
)

type GOAPGoal struct {
	spec map[string]int
	vars map[string]*utils.NumericInterval
}

func NewGOAPGoal(spec map[string]int) *GOAPGoal {
	g := &GOAPGoal{
		spec: spec,
		vars: make(map[string]*utils.NumericInterval),
	}
	g.Parametrize(1)

	return g
}

func (g *GOAPGoal) Parametrize(n int) *GOAPGoal {
	for spec, val := range g.spec {
		var split []string
		macroSplit := strings.Split(spec, ":")
		// if there is a macro ("EACH")
		if macroSplit[0] == "EACH" {
			val *= n
			split = strings.Split(macroSplit[1], ",")
		} else {
			split = strings.Split(spec, ",")
		}
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
		goal:         g,
		goalLeft:     make(map[string]*utils.NumericInterval),
		diffs:        make(map[string]float64),
		nUnfulfilled: 0,
	}
	if DEBUG_GOAP {
		debugGOAPPrintf("      -+- checking remaining for goal: %s", debugGOAPGoalToString(g))
		debugGOAPPrintf("      -+-     ws: %v", ws.vals)
	}
	for varName, interval := range g.vars {
		if stateVal, ok := ws.vals[varName]; ok {
			diff := interval.Diff(float64(stateVal))
			debugGOAPPrintf("                diff for %s: %.0f", varName, diff)
			result.diffs[varName] = diff
			if diff != 0 {
				result.nUnfulfilled++
				result.goalLeft[varName] = interval
			}
		} else {
			// varName not in worldstate - diff is infinite and goal is unchanged for this var
			result.nUnfulfilled++
			result.diffs[varName] = math.Inf(+1)
			result.goalLeft[varName] = interval
		}
	}
	return result
}
