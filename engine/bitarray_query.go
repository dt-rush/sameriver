package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type BitArraySubsetQuery struct {
	Match   bitarray.BitArray
}

func NewBitArraySubsetQuery (match bitarray.BitArray) BitArraySubsetQuery {
	return BitArraySubsetQuery{match}
}

func (bq BitArraySubsetQuery) Test (
	id uint16, entity_manager *EntityManager) bool {

	Logger.Printf ("Testing q%s against e%s...\n",
		ComponentBitArrayToString(bq.Match),
		ComponentBitArrayToString(entity_manager.EntityComponentBitArrays[id]))
	return bq.Match.Equals (
		bq.Match.And (entity_manager.EntityComponentBitArrays[id]))
}
