package engine

import (
	"bytes"
	"github.com/golang-collections/go-datastructures/bitarray"
)

// Query for whether the bitarray (Match) is a subset of the target
// BitArray
func NewEntityComponentBitArrayQuery(
	q bitarray.BitArray) GenericEntityQuery {

	return GenericEntityQuery{
		Name: BitArrayToString(q),
		TestFunc: func(id uint16, em *EntityManager) bool {
			// determine if q = q&b
			// that is, if every set bit of q is set in b
			b := em.EntityComponentBitArray(id)
			return q.Equals(q.And(b))
		}}
}

func BitArrayToString(arr bitarray.BitArray) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i := uint64(0); i < N_COMPONENTS; i++ {
		bit, _ := arr.GetBit(i)
		var val int
		if bit {
			val = 1
		} else {
			val = 0
		}
		buf.WriteString(string(val))
	}
	buf.WriteString("]")
	return buf.String()
}
