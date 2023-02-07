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

func NewGOAPAction(
	name string,
	cost int,
	pres map[string]int,
	effs map[string]int) *GOAPAction {

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

func NewGOAPActionModal(
	name string,
	cost int,
	pres map[string]int,
	preModalChecks map[string]GOAPModalVal,
	effs map[string]GOAPStateVal) *GOAPAction {

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
