package sameriver

import (
	"fmt"
	"sort"
)

// notice the sortf returned by Evaluate() is a closure that wants the result string so it can actually use i, j int

func DSLEval(expr string, resolver IdentifierResolver) (func(*Entity) bool, func(xs []*Entity) func(i, j int) bool, error) {
	parser := &EntityFilterDSLParser{}
	evaluator := NewEntityFilterDSLEvaluator(EntityFilterDSLPredicates, EntityFilterDSLSorts)

	ast, err := parser.Parse(expr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse expr: %s", err)
	}

	filter, sort := evaluator.Evaluate(ast, resolver)

	return filter, sort, nil
}

func DSLFilter(expr string, resolver IdentifierResolver, world *World) ([]*Entity, error) {
	filter, _, err := DSLEval(expr, resolver)
	if err != nil {
		return nil, err
	}
	result := world.FilterAllEntities(filter)
	return result, nil
}

func DSLFilterSort(expr string, resolver IdentifierResolver, world *World) ([]*Entity, error) {
	filterf, sortf, err := DSLEval(expr, resolver)
	if err != nil {
		return nil, err
	}
	result := world.FilterAllEntities(filterf)
	if sortf != nil {
		sort.Slice(result, sortf(result))
	}
	return result, nil
}

func (e *Entity) DSLFilter(expr string) ([]*Entity, error) {
	resolver := &EntityResolver{e: e}
	return DSLFilter(expr, resolver, e.World)
}

func (w *World) DSLFilter(expr string) ([]*Entity, error) {
	resolver := &WorldResolver{w: w}
	return DSLFilter(expr, resolver, w)
}

func (e *Entity) DSLFilterSort(expr string) ([]*Entity, error) {
	resolver := &EntityResolver{e: e}
	return DSLFilterSort(expr, resolver, e.World)
}

func (w *World) DSLFilterSort(expr string) ([]*Entity, error) {
	resolver := &WorldResolver{w: w}
	return DSLFilterSort(expr, resolver, w)
}
