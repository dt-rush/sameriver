package engine

import (
	"sync"

	"github.com/golang-collections/go-datastructures/bitarray"

	"github.com/dt-rush/sameriver/engine/utils"
)

type EntityToken struct {
	ID                int
	WorldID           utils.ID
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
