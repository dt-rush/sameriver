package sameriver

import (
	"fmt"
	"sort"
)

var EFDSL = NewEFDSLEvaluator()

// notice the sortf returned by Evaluate() is a closure that wants the result string so it can actually use i, j int

func EFDSLEval(expr string, resolver IdentifierResolver) (func(*Entity) bool, func(xs []*Entity) func(i, j int) bool, error) {
	parser := &EFDSLParser{}

	ast, err := parser.Parse(expr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse expr: %s", err)
	}

	filter, sort := EFDSL.Evaluate(ast, resolver)

	return filter, sort, nil
}

func EFDSLFilter(expr string, resolver IdentifierResolver, world *World) ([]*Entity, error) {
	filter, _, err := EFDSLEval(expr, resolver)
	if err != nil {
		return nil, err
	}
	result := world.FilterAllEntities(filter)
	return result, nil
}

func EFDSLFilterSort(expr string, resolver IdentifierResolver, world *World) ([]*Entity, error) {
	filterf, sortf, err := EFDSLEval(expr, resolver)
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
	return EFDSLFilter(expr, resolver, e.World)
}

func (w *World) DSLFilter(expr string) ([]*Entity, error) {
	resolver := &WorldResolver{w: w}
	return EFDSLFilter(expr, resolver, w)
}

func (e *Entity) DSLFilterSort(expr string) ([]*Entity, error) {
	resolver := &EntityResolver{e: e}
	return EFDSLFilterSort(expr, resolver, e.World)
}

func (w *World) DSLFilterSort(expr string) ([]*Entity, error) {
	resolver := &WorldResolver{w: w}
	return EFDSLFilterSort(expr, resolver, w)
}
