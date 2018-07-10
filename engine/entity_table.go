package engine

import (
	"errors"
	"fmt"
)

// used by the EntityManager to hold info about the allocated entities
type EntityTable struct {
	// how many entities there are
	numEntities int
	// list of Entities which have been allocated
	currentEntities []*EntityToken
	// list of available entity ID's which have previously been deallocated
	availableIDs []int
}

// get the ID for a new entity. Only called by SpawnEntity, which locks
// the entityTable, so it's safe that this method operates on that data.
// Returns int32 so that we can return -1 in case we have run out of space
// to spawn entities
func (t *EntityTable) allocateID() (*EntityToken, error) {
	// if maximum entity count reached, fail with message
	if t.numEntities == MAX_ENTITIES {
		msg := fmt.Sprintf("Reached max entity count: %d. "+
			"Will not allocate ID.\n", MAX_ENTITIES)
		Logger.Println(msg)
		return nil, errors.New(msg)
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
	// return the token
	entity := EntityToken{ID: id, Active: false, Despawned: false}
	return &entity, nil
}

func (t *EntityTable) addToCurrentEntities(entity *EntityToken) {
	t.currentEntities = append(t.currentEntities, entity)
}

func (t *EntityTable) snapshotAllocatedEntities() []*EntityToken {
	snapshot := make([]*EntityToken, len(t.currentEntities))
	copy(snapshot, t.currentEntities)
	return snapshot
}
