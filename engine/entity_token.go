package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityToken struct {
	ID                int
	Active            bool
	Despawned         bool
	ComponentBitArray bitarray.BitArray
}
