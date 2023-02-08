package sameriver

import (
	"strings"
)

type GOAPAction struct {
	name            string
	cost            IntOrFunc // (interface{})
	pres            *GOAPGoal
	preModalChecks  map[string]func(ws *GOAPWorldState) int
	effs            map[string]*GOAPEff
	effModalSetters map[string]func(ws *GOAPWorldState, op string, x int)
	varsAffected    map[string]bool
}

type GOAPEff struct {
	val int
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
		varsAffected:    make(map[string]bool),
	}
	a.effs = make(map[string]*GOAPEff)
	for spec, val := range effs {
		split := strings.Split(spec, ",")
		varName, op := split[0], split[1]
		a.effs[spec] = &GOAPEff{
			val,
			GOAPEffFunc(op, val),
		}
		a.varsAffected[varName] = true
	}
	return a
}

func (a *GOAPAction) affectsAnUnfulfilledVar(goal *GOAPGoal, preGoals map[string]*GOAPGoal) bool {
	for varName, _ := range goal.goals {
		if _, ok := a.varsAffected[varName]; ok {
			return true
		}
	}
	for _, pre := range preGoals {
		for varName, _ := range pre.goals {
			if _, ok := a.varsAffected[varName]; ok {
				return true
			}
		}
	}
	return false
}
