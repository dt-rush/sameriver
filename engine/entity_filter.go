package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityFilter struct {
	Name      string
	Predicate func(entity *EntityToken) bool
}

func NewEntityFilter(
	name string, f func(e *EntityToken) bool) EntityFilter {
	return EntityFilter{Name: name, Predicate: f}
}

func (q EntityFilter) Test(entity *EntityToken) bool {
	return q.Predicate(entity)
}

func (w *World) entityFilterFromTag(tag string) EntityFilter {
	return EntityFilter{
		Name: tag,
		Predicate: func(entity *EntityToken) bool {
			return w.Components.TagList[entity.ID].Has(tag)
		}}
}

func EntityFilterFromComponentBitArray(
	name string, q bitarray.BitArray) EntityFilter {
	return EntityFilter{
		Name: name,
		Predicate: func(entity *EntityToken) bool {
			// determine if q = q&b
			// that is, if every set bit of q is set in b
			b := entity.ComponentBitArray
			return q.Equals(q.And(b))
		}}
}
