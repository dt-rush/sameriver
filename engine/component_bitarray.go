package engine

import (
	"bytes"
	"fmt"

	"github.com/golang-collections/go-datastructures/bitarray"
)

// builds a bitarray where the bit corresponding to each component in a list
// of components is set
func MakeComponentBitArray(components []ComponentType) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(N_COMPONENT_TYPES))
	for _, COMPONENT := range components {
		b.SetBit(uint64(COMPONENT))
	}
	return b
}

// prints a string representation of a component bitarray as a set with
// string representations of each component type whose bit is set
func ComponentBitArrayToString(b bitarray.BitArray) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i := uint64(0); i < N_COMPONENT_TYPES; i++ {
		bit, _ := b.GetBit(i)
		// the index into the array is the component type int from the
		// iota const block in component_enum.go
		if bit {
			buf.WriteString(fmt.Sprintf("%s", COMPONENT_NAMES[ComponentType(i)]))
			if i != N_COMPONENT_TYPES-1 {
				buf.WriteString(", ")
			}
		}
	}
	buf.WriteString("]")
	return buf.String()
}
