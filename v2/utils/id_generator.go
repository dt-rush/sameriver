package utils

import (
	"go.uber.org/atomic"
	"math"
)

type IDGenerator struct {
	universe map[int]bool
	freed    map[int]bool
	x        atomic.Uint32
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{
		universe: make(map[int]bool),
		freed:    make(map[int]bool),
	}
}

func (g *IDGenerator) Next() (ID int) {
	// try to get ID from already-available freed IDs
	if len(g.freed) > 0 {
		// get first of freed (break immediately)
		for freeID, _ := range g.freed {
			ID = freeID
			delete(g.freed, freeID)
			break
		}
	} else {
		unique := false
		for !unique {
			u32ID := g.x.Inc()
			if u32ID > math.MaxUint32/64 {
				panic("tried to generate more than (2^32 - 1) / 64 simultaneous " +
					"ID's without free. This is surely a logic error.")
			}
			ID = int(u32ID)
			_, already := g.universe[ID]
			unique = !already
		}
	}
	g.universe[ID] = true
	return ID
}

func (g *IDGenerator) Free(ID int) {
	delete(g.universe, ID)
	g.freed[ID] = true
}
