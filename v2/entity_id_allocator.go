package sameriver

import (
	"errors"
	"fmt"

	"github.com/dt-rush/sameriver/v2/utils"
)

// used by the EntityManager to hold info about the allocated entities
type EntityIDAllocator struct {
	// the ID Generator given by the world the entity manager is in
	IdGen *utils.IDGenerator
	// list of available entity ID's which have previously been deallocated
	availableIDs []int
	// list of Entities which have been allocated
	currentEntities map[*Entity]bool
	// how many entities are active
	active int
	// capacity of how many ID's we can allocate without expanding
	capacity int
}

func NewEntityIDAllocator(capacity int, IDGen *utils.IDGenerator) *EntityIDAllocator {
	return &EntityIDAllocator{
		IdGen:           IDGen,
		currentEntities: make(map[*Entity]bool),
		capacity:        capacity,
	}
}

func (a *EntityIDAllocator) expand(n int) {
	a.capacity += n
}

// get the ID for a new e. Only called by SpawnEntity, which locks
// the entityTable, so it's safe that this method operates on that data.
// Returns int32 so that we can return -1 in case we have run out of space
// to spawn entities
func (a *EntityIDAllocator) allocateID() (*Entity, error) {
	// if maximum entity count reached, fail with message
	if len(a.currentEntities) == a.capacity {
		msg := fmt.Sprintf("[WARNING] Reached entity capacity: %d. "+
			"Will not allocate ID (entity manager should expand and retry).", a.capacity)
		Logger.Println(msg)
		return nil, errors.New(msg)
	}
	// if there is a deallocated entity somewhere in the table before the
	// highest ID, return that ID to the caller
	var ID int
	n_avail := len(a.availableIDs)
	if n_avail > 0 {
		// there is an ID available for a previously deallocated e.
		// pop it from the list and continue with that as the ID
		ID = a.availableIDs[n_avail-1]
		a.availableIDs = a.availableIDs[:n_avail-1]
	} else {
		// every slot in the table before the highest ID is filled
		ID = len(a.currentEntities)
	}
	entity := &Entity{
		ID:      ID,
		WorldID: a.IdGen.Next(),
	}
	a.currentEntities[entity] = true
	return entity, nil
}

func (a *EntityIDAllocator) deallocate(e *Entity) {
	// guards against false deallocation (edge case, but hey)
	if _, ok := a.currentEntities[e]; ok {
		a.availableIDs = append(a.availableIDs, e.ID)
		a.IdGen.Free(e.WorldID)
		delete(a.currentEntities, e)
	}
}
