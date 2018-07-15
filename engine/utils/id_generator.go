package utils

import (
	"go.uber.org/atomic"
	"math"
)

type ID = int

type IDGenerator struct {
	universe map[ID]bool
	freed    map[ID]bool
	x        atomic.Uint32
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{
		universe: make(map[ID]bool),
		freed:    make(map[ID]bool),
	}
}

func (g *IDGenerator) Next() ID {
	var ID ID
	// try to get ID from already-available freed IDs
	if len(g.freed) > 0 {
		for freeID, _ := range g.freed {
			ID = freeID
			break
		}
		delete(g.freed, ID)
	}
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
	g.universe[ID] = true
	return ID
}

func (g *IDGenerator) Free(ID ID) {
	delete(g.universe, ID)
	g.freed[ID] = true
}
