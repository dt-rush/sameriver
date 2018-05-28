package engine

import (
	"github.com/golang-collections/go-datastructures/bitarray"
	"go.uber.org/atomic"
	"sync"
	goatomic "sync/atomic"
)

// used by the EntityManager to hold info about the allocated entities
type EntityTable struct {
	// mutex used to make the allocation / deallocation of an ID atomic
	IDMutex sync.RWMutex
	// how many entities there are
	numEntities int
	// list of Entities which have been allocated
	currentEntities []EntityToken
	// list of available entity ID's which have previously been deallocated
	availableIDs []int
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
	locks [MAX_ENTITIES]atomic.Uint32
}

func (t *EntityTable) incrementGen(id int) {
	goatomic.AddUint32(&t.gens[id], 1)
}

func (t *EntityTable) getGen(id int) uint32 {
	return goatomic.LoadUint32(&t.gens[id])
}

func (t *EntityTable) getEntityToken(id int) EntityToken {
	// we may want to get an entity token for -(id + 1) in the case
	// that this token represents a remove signal to an UpdatedEntityList.
	// handle this by getting the correct gen but leaving ID as negative
	var gen uint32
	if id < 0 {
		gen = t.getGen(-(id + 1))
	} else {
		gen = t.getGen(id)
	}
	return EntityToken{id, gen}
}

func (t *EntityTable) genValidate(entity EntityToken) bool {
	return t.getGen(entity.ID) == entity.gen
}
