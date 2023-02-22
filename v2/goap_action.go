package sameriver

import (
	"strings"
)

type GOAPAction struct {
	name            string
	cost            IntOrFunc
	pres            *GOAPGoal
	preModalChecks  map[string]func(ws *GOAPWorldState) int
	effModalSetters map[string]func(ws *GOAPWorldState, op string, x int)
	effs            map[string]*GOAPEff
	ops             map[string]string
}

type GOAPEff struct {
	val int
	op  string
	f   func(int) int
}

func GOAPEffFunc(op string, val int) func(int) int {
	switch op {
	case "+":
		return func(x int) int { return val + x }
	case "-":
		return func(x int) int { return x - val }
	case "=":
		return func(x int) int { return val }
	default:
		panic("Got an undefined op in GOAPEffFunc() [valid: +,-,=]")
	}
}

func NewGOAPAction(def map[string]interface{}) *GOAPAction {
	name := def["name"].(string)
	cost := def["cost"].(int)
	var pres map[string]int
	if def["pres"] == nil {
		pres = nil
	} else {
		pres = def["pres"].(map[string]int)
	}
	effs := def["effs"].(map[string]int)

	a := &GOAPAction{
		name:            name,
		cost:            cost,
		pres:            NewGOAPGoal(pres),
		preModalChecks:  make(map[string]func(ws *GOAPWorldState) int),
		effModalSetters: make(map[string]func(ws *GOAPWorldState, op string, x int)),
		effs:            make(map[string]*GOAPEff),
		ops:             make(map[string]string),
	}
	for spec, val := range effs {
		split := strings.Split(spec, ",")
		varName, op := split[0], split[1]
		eff := &GOAPEff{
			val: val,
			op:  op,
			f:   GOAPEffFunc(op, val),
		}
		a.effs[varName] = eff
		a.ops[varName] = op
	}
	return a
}

func (a *GOAPAction) affectsAnUnfulfilledVar(goal *GOAPGoal, preGoals map[string]*GOAPGoal) bool {
	for varName, _ := range goal.goals {
		if _, ok := a.effs[varName]; ok {
			return true
		}
	}
	for _, pre := range preGoals {
		for varName, _ := range pre.goals {
			if _, ok := a.effs[varName]; ok {
				return true
			}
		}
	}
	return false
}
