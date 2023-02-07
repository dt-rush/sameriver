package sameriver

import (
	"strings"
)

type GOAPGoalVal struct {
	val        int
	comparator func(int) int
}

type GOAPGoal struct {
	goals map[string]GOAPGoalVal
}

func GOAPGoalFunc(op string, val int) func(int) int {
	switch op {
	case "<":
		return func(x int) int { return -intMax(0, x-(val-1)) }
	case "<=":
		return func(x int) int { return -intMax(0, x-val) }
	case "=":
		return func(x int) int { return val - x }
	case ">":
		return func(x int) int { return intMax(0, (val+1)-x) }
	case ">=":
		return func(x int) int { return intMax(0, val-x) }
	default:
		panic("Got undefined op in GOAPGoalFunc() [valid: >=,>,=,<,<=]")
	}
}

func NewGOAPGoal(def map[string]int) *GOAPGoal {
	g := &GOAPGoal{
		goals: make(map[string]GOAPGoalVal),
	}
	for spec, val := range def {
		split := strings.Split(spec, ",")
		op := split[1]
		g.goals[spec] = GOAPGoalVal{
			val,
			GOAPGoalFunc(op, val),
		}
	}
	return g
}

func (g *GOAPGoal) goalRemaining(ws *GOAPWorldState) (goalRemaining *GOAPGoal, diffs map[string]int) {
	goalRemaining = NewGOAPGoal(nil)
	diffs = make(map[string]int)
	for spec, goalVal := range g.goals {
		split := strings.Split(spec, ",")
		varName := split[0]
		op := split[1]
		var diff int
		// NOTE: we assume a var not present in the state as having val = 0
		if stateVal, ok := ws.vals[varName]; ok {
			diff = goalVal.comparator(stateVal)
		} else {
			diff = goalVal.comparator(0)
		}
		diffs[varName] = diff
		if diff == 0 {
			// goal fulfilled
			// NOTE: this will *entirely remove any goal related to this var*,
			// so if, in backtracking, we add an action that fulfills another
			// goal but modifies this var as well, we may end up - playing
			// the whole chain - with an invalid solution. Thus we should be
			// sure to - upon finding a "satisfied chain", check it by playing
			// it forward and seeing if the resulting GOAPWorldState's
			// goalRemaining with the true original goal is actually all zero'd
			continue
		}
		switch op {
		case ">=":
			goalRemaining.goals[spec] = GOAPGoalVal{
				val:        diff,
				comparator: GOAPGoalFunc(op, diff),
			}
		case ">":
			goalRemaining.goals[spec] = GOAPGoalVal{
				val:        diff - 1,
				comparator: GOAPGoalFunc(op, diff-1),
			}
		case "=":
			goalRemaining.goals[spec] = GOAPGoalVal{
				val:        diff,
				comparator: GOAPGoalFunc(op, diff),
			}
		case "<=":
			goalRemaining.goals[spec] = GOAPGoalVal{
				val:        diff,
				comparator: GOAPGoalFunc(op, diff),
			}
		case "<":
			goalRemaining.goals[spec] = GOAPGoalVal{
				val:        diff + 1,
				comparator: GOAPGoalFunc(op, diff+1),
			}
		default:
			panic("Got an undefined op in goalRemaining() [valid: >=,>,=,<,<=]")
		}
	}
	return goalRemaining, diffs
}

func (g *GOAPGoal) stateCloserInSomeVar(after, before *GOAPWorldState) (closer bool, afterRemaining *GOAPGoal) {
	afterRemaining, afterDiffs := g.goalRemaining(after)
	_, beforeDiffs := g.goalRemaining(before)
	for varName, _ := range beforeDiffs {
		if intAbs(afterDiffs[varName]) < intAbs(beforeDiffs[varName]) {
			return true, afterRemaining
		}
	}
	return false, afterRemaining
}

func (g *GOAPGoal) merged(other *GOAPGoal) *GOAPGoal {
	result := &GOAPGoal{
		goals: make(map[string]GOAPGoalVal),
	}
	// first, copy this goal's vals in
	for spec, val := range g.goals {
		result.goals[spec] = val
	}
	// now........ the monstrous task
	// TODO
	return result
}
