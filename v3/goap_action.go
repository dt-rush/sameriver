package sameriver

import (
	"fmt"

	"strings"
)

type GOAPAction struct {
	// the object used to construct this (used in Parametrized() to reconstruct)
	spec map[string]any

	Name            string
	Count           int
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

func GOAPEffFunc(a *GOAPAction, op string, val int) func(int) int {
	switch op {
	case "+":
		return func(x int) int { return a.Count*val + x }
	case "-":
		return func(x int) int { return x - a.Count*val }
	case "=":
		return func(x int) int { return val }
	default:
		panic("Got an unspecined op in GOAPEffFunc() [valid: +,-,=]")
	}
}

func NewGOAPAction(spec map[string]interface{}) *GOAPAction {
	name := spec["name"].(string)
	cost := spec["cost"].(int)
	var pres map[string]int
	if spec["pres"] == nil {
		pres = nil
	} else {
		pres = spec["pres"].(map[string]int)
	}
	effs := spec["effs"].(map[string]int)

	a := &GOAPAction{
		spec:            spec,
		Name:            name,
		Count:           1,
		cost:            cost,
		pres:            NewGOAPGoal(pres).Parametrize(1),
		preModalChecks:  make(map[string]func(ws *GOAPWorldState) int),               // set by GOAPEvaluator
		effModalSetters: make(map[string]func(ws *GOAPWorldState, op string, x int)), // set by GOAPEvaluator
		effs:            make(map[string]*GOAPEff),
		ops:             make(map[string]string),
	}
	for spec, val := range effs {
		split := strings.Split(spec, ",")
		varName, op := split[0], split[1]
		eff := &GOAPEff{
			val: val,
			op:  op,
			f:   GOAPEffFunc(a, op, val),
		}
		a.effs[varName] = eff
		a.ops[varName] = op
	}
	return a
}

func (a *GOAPAction) DisplayName() string {
	if a.Count == 1 {
		return a.Name
	} else {
		return fmt.Sprintf("%s(%d)", a.Name, a.Count)
	}
}

func (a *GOAPAction) CopyOf() *GOAPAction {
	result := NewGOAPAction(a.spec)
	result.Count = a.Count
	result.preModalChecks = a.preModalChecks
	result.effModalSetters = a.effModalSetters
	return result
}

func (a *GOAPAction) Parametrized(n int) *GOAPAction {
	result := a.CopyOf()
	result.Count = n
	result.pres = result.pres.Parametrize(n)
	return result
}
