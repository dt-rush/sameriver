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
	PrimaryLogic      *LogicUnit
	Logics            map[string]*LogicUnit
}

func (e *Entity) LogicUnitName(name string) string {
	return fmt.Sprintf("entity-logic-%d-%s", e.ID, name)
}

func (e *Entity) MakeLogicUnit(name string, F func()) *LogicUnit {
	return &LogicUnit{
		name:    e.LogicUnitName(name),
		f:       F,
		active:  true,
		worldID: e.World.IdGen.Next(),
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
