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

// get the ID for a new entity. Only called by SpawnEntity, which locks
// the entityTable, so it's safe that this method operates on that data.
// Returns int32 so that we can return -1 in case we have run out of space
// to spawn entities
func (t *EntityTable) allocateID() (EntityToken, error) {
	t.IDMutex.Lock()
	defer t.IDMutex.Unlock()
	// if maximum entity count reached, fail with message
	if t.numEntities == MAX_ENTITIES {
		msg := fmt.Sprintf("Reached max entity count: %d. "+
			"Will not allocate ID.\n", MAX_ENTITIES)
		Logger.Println(msg)
		return ENTITY_TOKEN_NIL, errors.New(msg)
	}
	// Increment the entity count
	t.numEntities++
	// if there is a deallocated entity somewhere in the table before the
	// highest ID, return that ID to the caller
	n_avail := len(t.availableIDs)
	var id int
	if n_avail > 0 {
		// there is an ID available for a previously deallocated entity.
		// pop it from the list and continue with that as the ID
		id = t.availableIDs[n_avail-1]
		t.availableIDs = t.availableIDs[:n_avail-1]
	} else {
		// every slot in the table before the highest ID is filled
		id = t.numEntities - 1
	}
	// add the ID to the list of allocated IDs
	entity := EntityToken{id, t.gens[id]}
	t.currentEntities = append(t.currentEntities, entity)
	return entity, nil
}

// lock the ID table after waiting on spawn mutex to be unlocked,
// and grab a copy of the currently allocated IDs
func (t *EntityTable) snapshotAllocatedEntities() []EntityToken {
	t.IDMutex.RLock()
	updatedEntityListDebug("got IDMutex in snapshot")
	defer t.IDMutex.RUnlock()

	snapshot := make([]EntityToken, len(t.currentEntities))
	copy(snapshot, t.currentEntities)
	return snapshot
}


