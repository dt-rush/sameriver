package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"sync"
)

type EntityToken struct {
	ID                int
	Active            bool
	Despawned         bool
	ComponentBitArray bitarray.BitArray
	ListsMutex        sync.RWMutex
	Lists             []*UpdatedEntityList
}
