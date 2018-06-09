package engine

import (
	"bytes"
	"github.com/golang-collections/go-datastructures/bitarray"
)

func BitArrayToString(arr bitarray.BitArray) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i := uint64(0); i < N_COMPONENT_TYPES; i++ {
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
