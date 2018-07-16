package engine

import (
	"testing"

	"github.com/dt-rush/sameriver/engine/utils"
)

func TestEntityTableAllocateID(t *testing.T) {
	et := NewEntityTable(utils.NewIDGenerator())
	_, err := et.allocateID()
	if !(err == nil &&
		len(et.currentEntities) == 1) {
		t.Fatal("didn't allocate ID properly")
	}
}

func TestEntityTableDeallocateID(t *testing.T) {
	et := NewEntityTable(utils.NewIDGenerator())
	e, _ := et.allocateID()
	et.deallocate(e)
	if len(et.currentEntities) != 0 {
		t.Fatal("didn't update allocated count")
	}
	if !(len(et.availableIDs) == 1 && et.availableIDs[0] == e.ID) {
		t.Fatal("didn't add deallocated ID to list of available IDs")
	}
}

func TestEntityTableAllocateMaxIDs(t *testing.T) {
	et := NewEntityTable(utils.NewIDGenerator())
	for i := 0; i < MAX_ENTITIES; i++ {
		et.allocateID()
	}
	_, err := et.allocateID()
	if err == nil {
		t.Fatal("should have returned error on allocating > MAX_ENTITIES")
	}
}

func TestEntityTableReallocateDeallocatedID(t *testing.T) {
	et := NewEntityTable(utils.NewIDGenerator())
	var e *EntityToken
	for i := 0; i < MAX_ENTITIES; i++ {
		allocated, _ := et.allocateID()
		if i == MAX_ENTITIES/2 {
			e = allocated
		}
	}
	et.deallocate(e)
	e, err := et.allocateID()
	if err != nil {
		t.Fatal("should have had space after deallocate to allocate")
	}
	if e.ID != MAX_ENTITIES/2 {
		t.Fatal("should have used deallocated ID to serve new allocate request")
	}
}
