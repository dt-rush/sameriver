package sameriver

import (
	"fmt"

	"strings"
)

type GOAPAction struct {
	// the action for one of whose pre's this action is a satisfier
	// (nil if it's a satisfier for the main goal)
	parent *GOAPAction

	// the index in the list at which this action is inserted
	insertionIx int
	// the region index of the temporal region this was inserted into, satisfying
	regionIx int

	// the object used to construct this (used in Parametrized() to reconstruct)
	spec map[string]any

	// something like "chopWood"
	Name string

	// the entity from planner.BindEntititySelectors() that will be used for modal
	// considerations
	Node string

	// whether, during the action, we move with the node
	travelWithNode bool

	// other nodes we want to look at modally - named by string a la BindEntitySelectors()
	otherNodes []string

	// how many times we'll do the action
	Count int

	// the inherent cost of the action as an int or function
	cost IntOrFunc

	// preconditions for the action to be performed
	pres *GOAPTemporalGoal

	// functions that will check the numeric value of a varName (the string key) modally
	preModalChecks map[string]func(ws *GOAPWorldState) int
	// functions that will set the varName (the string key) to a certain value modally
	effModalSetters map[string]func(ws *GOAPWorldState, op string, x int)

	// effects of this action on ws numerically (non-modal)
	effs map[string]*GOAPEff

	// k: varName, v: op
	// eg, "health": "+"
	ops map[string]string
}

type GOAPEff struct {
	val int
	op  string
	f   func(count int, x int) int
}

func GOAPEffFunc(op string, val int) func(count int, x int) int {
	switch op {
	case "+":
		return func(count int, x int) int { return count*val + x }
	case "-":
		return func(count int, x int) int { return x - count*val }
	case "=":
		return func(count int, x int) int { return val }
	default:
		panic("Got an unspecified op in GOAPEffFunc() [valid: +,-,=]")
	}
}

func NewGOAPAction(spec map[string]any) *GOAPAction {
	name := spec["name"].(string)
	node := spec["node"].(string)
	travelWithNode, ok := spec["travelWithNode"].(bool)
	if !ok {
		travelWithNode = false
	}
	otherNodes, ok := spec["otherNodes"].([]string)
	if !ok {
		otherNodes = []string{}
	}
	cost := spec["cost"].(int)
	pres := spec["pres"]
	effs := spec["effs"].(map[string]int)

	a := &GOAPAction{
		spec:            spec,
		Name:            name,
		Node:            node,
		travelWithNode:  travelWithNode,
		otherNodes:      otherNodes,
		Count:           1,
		cost:            cost,
		pres:            NewGOAPTemporalGoal(pres),
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
			f:   GOAPEffFunc(op, val),
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
	result := *a
	return &result
}

func (a *GOAPAction) Parametrized(n int) *GOAPAction {
	logGOAPDebug("    parametrizing %s x %d", a.Name, n)
	result := a.CopyOf()
	result.Count = n
	result.pres = result.pres.Parametrized(n)
	return result
}

func (a *GOAPAction) ChildOf(p *GOAPAction) *GOAPAction {
	result := a.CopyOf()
	result.parent = p
	return result
}
