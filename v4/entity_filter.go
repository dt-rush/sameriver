package sameriver

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityFilter struct {
	Name      string
	Predicate func(e *Entity) bool
}

func NewEntityFilter(
	name string, f func(e *Entity) bool) EntityFilter {
	return EntityFilter{Name: name, Predicate: f}
}

func (q EntityFilter) Test(e *Entity) bool {
	return q.Predicate(e)
}

func EntityFilterFromTag(tag string) EntityFilter {
	return EntityFilter{
		Name: tag,
		Predicate: func(e *Entity) bool {
			return e.GetTagList(GENERICTAGS).Has(tag)
		}}
}

func EntityFilterFromComponentBitArray(
	name string, q bitarray.BitArray) EntityFilter {
	return EntityFilter{
		Name: name,
		Predicate: func(e *Entity) bool {
			// determine if q = q&b
			// that is, if every set bit of q is set in b
			b := e.ComponentBitArray
			return q.Equals(q.And(b))
		}}
}
