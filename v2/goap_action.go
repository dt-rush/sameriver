package sameriver

import (
	"strings"
)

type GOAPAction struct {
	name            string
	cost            IntOrFunc // (interface{})
	pres            *GOAPGoal
	preModalChecks  map[string]func(ws *GOAPWorldState) int
	effs            map[string]func(int) int
	effModalSetters map[string]func(ws *GOAPWorldState)
}

func GOAPEffFunc(op string, val int) func(int) int {
	switch op {
	case "+":
		return func(x int) int { return x + val }
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
		effModalSetters: make(map[string]func(ws *GOAPWorldState)),
	}
	a.effs = make(map[string]func(int) int)
	for spec, val := range effs {
		split := strings.Split(spec, ",")
		varName := split[0]
		op := split[1]
		a.effs[varName] = GOAPEffFunc(op, val)
	}
	return a
}

func NewGOAPActionModal(def map[string]interface{}) *GOAPAction {

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
		effModalSetters: make(map[string]func(ws *GOAPWorldState)),
	}
	a.effs = make(map[string]func(int) int)
	for spec, val := range effs {
		split := strings.Split(spec, ",")
		varName := split[0]
		op := split[1]
		a.effs[varName] = GOAPEffFunc(op, val)
	}
	return a
}
