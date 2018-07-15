package engine

import (
	"sync"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type EntityToken struct {
	ID                int
	WorldID           int
	Active            bool
	Despawned         bool
	ComponentBitArray bitarray.BitArray
	ListsMutex        sync.RWMutex
	Lists             []*UpdatedEntityList
}

func (e *EntityToken) MakeLogicUnit(Name string, F func()) *LogicUnit {
	return &LogicUnit{
		Name:    Name,
		F:       F,
		Active:  false,
		WorldID: e.WorldID,
	}
}
