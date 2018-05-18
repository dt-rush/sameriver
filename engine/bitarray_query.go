package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

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
func (bq BitArraySubsetQuery) Test(b bitarray.BitArray) bool {
	q := bq.Match
	return q.Equals(q.And(b))
}
