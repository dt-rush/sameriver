package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

// Query for whether the bitarray (Match) is a subset of the target
// BitArray
type EntityComponentBitArrayQuery struct {
	q bitarray.BitArray
}

func NewEntityComponentBitArrayQuery(
	q bitarray.BitArray) EntityComponentBitArrayQuery {

	return EntityComponentBitArrayQuery{q}
}

// determine if q = q&b
// that is, if every set bit of q is set in b
func (bq EntityComponentBitArrayQuery) Test(
	id uint16, entity_manager *EntityManager) bool {

	b := entity_manager.EntityComponentBitArray(id)
	return bq.q.Equals(bq.q.And(b))
}
