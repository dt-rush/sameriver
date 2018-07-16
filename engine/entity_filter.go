package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityFilter struct {
	Name     string
	TestFunc func(entity *EntityToken, em *EntityManager) bool
}

func (q EntityFilter) Test(entity *EntityToken, em *EntityManager) bool {
	result := q.TestFunc(entity, em)
	return result
}

func EntityFilterFromTag(tag string) EntityFilter {
	return EntityFilter{
		Name: tag,
		TestFunc: func(entity *EntityToken, em *EntityManager) bool {
			return em.Components.TagList[entity.ID].Has(tag)
		}}
}

func EntityFilterFromComponentBitArray(
	name string, q bitarray.BitArray) EntityFilter {
	return EntityFilter{
		Name: name,
		TestFunc: func(entity *EntityToken, em *EntityManager) bool {
			// determine if q = q&b
			// that is, if every set bit of q is set in b
			b := entity.ComponentBitArray
			return q.Equals(q.And(b))
		}}
}
