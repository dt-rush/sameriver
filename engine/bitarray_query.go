package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type BitArraySubsetQuery struct {
	Match bitarray.BitArray
}

func NewBitArraySubsetQuery(match bitarray.BitArray) BitArraySubsetQuery {
	return BitArraySubsetQuery{match}
}

func (bq BitArraySubsetQuery) Test(
	id uint16, entity_manager *EntityManager) bool {

	return bq.Match.Equals(
		bq.Match.And(entity_manager.EntityComponentBitArrays[id]))
}
