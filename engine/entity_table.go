package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"sync"
	"sync/atomic"
)

// used by the EntityManager to hold info about the allocated entities
type EntityTable struct {
	// guards all changes to this table as atomic
	mutex sync.RWMutex
	// how many entities there are
	numEntities int
	// list of IDs which have been allocated
	allocatedIDs []uint16
	// list of available entity ID's which have previously been deallocated
	availableIDs []uint16
	// bitarray used to keep track of which entities have which components
	// (indexes are IDs, bitarrays have bit set if entity has the
	// component corresponding to that index)
	componentBitArrays [MAX_ENTITIES]bitarray.BitArray
	// the gen of an ID is how many times an entity has been
	// spawned on that ID
	gens [MAX_ENTITIES]uint32
	// locks so that goroutines can operate atomically on individual entities
	// (eg. imagine two squirrels coming upon a nut and trying to eat it. One
	// must win!). Also used by systems like PhysicsSystem to avoid modifying
	// those entities while they're held for modification (hence the
	// importance of not holding entities for modification longer than, say,
	// one update cycle (at 60fps, around 16 ms). In fact, one update
	// cycle is a hell of a long time. It should be like 4 milliseconds at
	// most (thinking here of an inventory modification after comparing
	// inventory contents to entity desires)
	locks [MAX_ENTITIES]uint32
}

func (t *EntityTable) getGen(id uint16) uint32 {
	return atomic.LoadUint32(&t.gens[id])
}

func (t *EntityTable) incrementGen(id uint16) uint32 {
	return atomic.AddUint32(&t.gens[id], 1)
}

func (t *EntityTable) getEntityToken(id uint16) EntityToken {
	return EntityToken{int32(id), t.getGen(id)}
}

func (t *EntityTable) genValidate(entity EntityToken) bool {
	return t.getGen(uint16(entity.ID)) == entity.gen
}
