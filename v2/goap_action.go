package sameriver

import (
	"strings"
)

type GOAPStateVal interface{}

type GOAPAction struct {
	name            string
	cost            int
	pres            GOAPGoal
	preModalChecks  map[string]GOAPModalVal
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
		preModalChecks:  make(map[string]GOAPModalVal),
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
	var preModalChecks map[string]GOAPModalVal
	if def["checks"] == nil {
		preModalChecks = nil
	} else {
		preModalChecks = def["preModalChecks"].(map[string]GOAPModalVal)
	}
	effs := def["effs"].(map[string]GOAPStateVal)

	a := &GOAPAction{
		name:            name,
		cost:            cost,
		pres:            NewGOAPGoal(pres),
		preModalChecks:  preModalChecks,
		effModalSetters: make(map[string]func(ws *GOAPWorldState)),
	}
	a.effs = make(map[string]func(int) int)
	for spec, val := range effs {
		split := strings.Split(spec, ",")
		varName := split[0]
		op := split[1]
		if modalVal, ok := val.(GOAPModalVal); ok {
			a.effs[varName] = GOAPEffFunc(op, modalVal.valAsEff)
			a.effModalSetters[varName] = modalVal.effModalSet
		} else if intVal, ok := val.(int); ok {
			a.effs[varName] = GOAPEffFunc(op, intVal)
		} else {
			panic("value in effs map neither GOAPModalVal or int")
		}
	}
	return a
}
