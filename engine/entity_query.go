package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityQuery struct {
	Name     string
	TestFunc func(entity *EntityToken, em *EntityManager) bool
}

func (q EntityQuery) Test(entity *EntityToken, em *EntityManager) bool {
	result := q.TestFunc(entity, em)
	return result
}

func EntityQueryFromTag(tag string) EntityQuery {

	return EntityQuery{
		Name: tag,
		TestFunc: func(entity *EntityToken, em *EntityManager) bool {
			return em.Components.TagList[entity.ID].Has(tag)
		}}
}

func EntityQueryFromComponentBitArray(
	name string,
	q bitarray.BitArray) EntityQuery {

	return EntityQuery{
		Name: name,
		TestFunc: func(entity *EntityToken, em *EntityManager) bool {
			// determine if q = q&b
			// that is, if every set bit of q is set in b
			b := em.entityComponentBitArray(entity.ID)
			return q.Equals(q.And(b))
		}}
}
