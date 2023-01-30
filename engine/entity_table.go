package engine

import (
	"errors"
	"fmt"

	"github.com/dt-rush/sameriver/engine/utils"
)

// used by the EntityManager to hold info about the allocated entities
type EntityTable struct {
	// the ID Generator given by the world the entity manager is in
	idGen *utils.IDGenerator
	// list of available entity ID's which have previously been deallocated
	availableIDs []int
	// list of Entities which have been allocated
	currentEntities map[*Entity]bool
	// how many entities are active
	active int
}

func NewEntityTable(IDGen *utils.IDGenerator) *EntityTable {
	return &EntityTable{
		idGen:           IDGen,
		currentEntities: make(map[*Entity]bool),
	}
}

// get the ID for a new e. Only called by SpawnEntity, which locks
// the entityTable, so it's safe that this method operates on that data.
// Returns int32 so that we can return -1 in case we have run out of space
// to spawn entities
func (t *EntityTable) allocateID() (*Entity, error) {
	// if maximum entity count reached, fail with message
	if len(t.currentEntities) == MAX_ENTITIES {
		msg := fmt.Sprintf("Reached max entity count: %d. "+
			"Will not allocate ID.", MAX_ENTITIES)
		Logger.Println(msg)
		return nil, errors.New(msg)
	}
	// if there is a deallocated entity somewhere in the table before the
	// highest ID, return that ID to the caller
	var ID int
	n_avail := len(t.availableIDs)
	if n_avail > 0 {
		// there is an ID available for a previously deallocated e.
		// pop it from the list and continue with that as the ID
		ID = t.availableIDs[n_avail-1]
		t.availableIDs = t.availableIDs[:n_avail-1]
	} else {
		// every slot in the table before the highest ID is filled
		ID = len(t.currentEntities)
	}
	entity := &Entity{
		ID:        ID,
		WorldID:   t.idGen.Next(),
		Active:    false,
		Despawned: false,
		Logics:    make(map[string]*LogicUnit),
	}
	t.currentEntities[entity] = true
	return entity, nil
}

func (t *EntityTable) deallocate(e *Entity) {
	// guards against false deallocation (edge case, but hey)
	if _, ok := t.currentEntities[e]; ok {
		t.availableIDs = append(t.availableIDs, e.ID)
		t.idGen.Free(e.WorldID)
		delete(t.currentEntities, e)
	}
}
