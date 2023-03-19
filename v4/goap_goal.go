package sameriver

import (
	"math"
	"strings"
)

type GOAPGoal struct {
	spec map[string]int
	vars map[string]*NumericInterval
}

func newGOAPGoal(spec map[string]int) *GOAPGoal {
	g := &GOAPGoal{
		spec: spec,
		vars: make(map[string]*NumericInterval),
	}
	return g.Parametrized(1)
}

func (g *GOAPGoal) Parametrized(n int) *GOAPGoal {
	result := &GOAPGoal{
		spec: g.spec,
		vars: make(map[string]*NumericInterval),
	}
	for spec, val := range g.spec {
		logGOAPDebug("        parametrizing %s:%d by %d", spec, val, n)
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
		interval := MakeNumericInterval(op, val)
		result.vars[varName] = interval
	}
	return result
}

func (g *GOAPGoal) remaining(ws *GOAPWorldState) (result *GOAPGoalRemaining) {
	result = &GOAPGoalRemaining{
		goal:         g,
		goalLeft:     make(map[string]*NumericInterval),
		diffs:        make(map[string]float64),
		nUnfulfilled: 0,
	}
	if DEBUG_GOAP {
		logGOAPDebug("      -+- checking remaining for goal: %s", debugGOAPGoalToString(g))
		logGOAPDebug("      -+-     ws: %v", ws.vals)
	}
	for varName, interval := range g.vars {
		if stateVal, ok := ws.vals[varName]; ok {
			diff := interval.Diff(float64(stateVal))
			logGOAPDebug("                diff for %s: %.0f", varName, diff)
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
