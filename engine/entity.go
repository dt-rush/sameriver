package engine

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type Entity struct {
	ID                int
	World             *World
	WorldID           int
	Active            bool
	Despawned         bool
	ComponentBitArray bitarray.BitArray
	ListsMutex        sync.RWMutex
	Lists             []*UpdatedEntityList
}

func (e *Entity) LogicUnitName() string {
	return fmt.Sprintf("entity-logic-%d", e.ID)
}

func (e *Entity) MakeLogicUnit(F func()) *LogicUnit {
	return &LogicUnit{
		Name:    e.LogicUnitName(),
		F:       F,
		Active:  false,
		WorldID: e.WorldID,
	}
}

func EntitySliceToString(entities []*Entity) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, e := range entities {
		buf.WriteString(fmt.Sprintf("%d", e.ID))
		if i != len(entities)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}
