package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

// TODO: rename to clearly explain that this is an EntityQuery

// Query for whether the bitarray (Match) is a subset of the target
// BitArray
type BitArraySubsetQuery struct {
	q bitarray.BitArray
}

func NewBitArraySubsetQuery(q bitarray.BitArray) BitArraySubsetQuery {
	return BitArraySubsetQuery{q}
}

// determine if q = q&b
// that is, if every set bit of q is set in b
func (bq BitArraySubsetQuery) Test(
	id uint16, entity_manager *EntityManager) bool {

	b := entity_manager.EntityComponentBitArray(id)
	return bq.q.Equals(bq.q.And(b))
}
