package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityToken struct {
	ID                int
	active            bool
	despawned         bool
	ComponentBitArray bitarray.BitArray
}
