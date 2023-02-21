package sameriver

import (
	"testing"

	"github.com/dt-rush/sameriver/v2/utils"
)

func TestEntityIDAllocatorAllocateID(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, utils.NewIDGenerator())
	_, err := et.allocateID()
	if !(err == nil &&
		len(et.currentEntities) == 1) {
		t.Fatal("didn't allocate ID properly")
	}
}

func TestEntityIDAllocatorDeallocateID(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, utils.NewIDGenerator())
	e, _ := et.allocateID()
	et.deallocate(e)
	if len(et.currentEntities) != 0 {
		t.Fatal("didn't update allocated count")
	}
	if !(len(et.availableIDs) == 1 && et.availableIDs[0] == e.ID) {
		t.Fatal("didn't add deallocated ID to list of available IDs")
	}
}

func TestEntityIDAllocatorAllocateMaxIDs(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, utils.NewIDGenerator())
	for i := 0; i < MAX_ENTITIES; i++ {
		et.allocateID()
	}
	_, err := et.allocateID()
	if err == nil {
		t.Fatal("should have returned error on allocating > MAX_ENTITIES")
	}
	et.expand(1)
	_, err = et.allocateID()
	if err != nil {
		t.Fatal("should have been able to allocate after expand")
	}
}

func TestEntityIDAllocatorReallocateDeallocatedID(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, utils.NewIDGenerator())
	var e *Entity
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
