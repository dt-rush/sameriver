package engine

import (
	"bytes"
	"fmt"
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

func EntityTokenSliceToString(entities []*EntityToken) string {
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
