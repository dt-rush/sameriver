package sameriver

import (
	"testing"
)

func TestEntityIDAllocatorAllocateID(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, NewIDGenerator())
	et.allocateID()
}

func TestEntityIDAllocatorDeallocateID(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, NewIDGenerator())
	e := et.allocateID()
	et.deallocate(e)
	if len(et.currentEntities) != 0 {
		t.Fatal("didn't update allocated count")
	}
	if !(len(et.availableIDs) == 1 && et.availableIDs[0] == e.ID) {
		t.Fatal("didn't add deallocated ID to list of available IDs")
	}
}

func TestEntityIDAllocatorAllocateMaxIDs(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, NewIDGenerator())
	for i := 0; i < MAX_ENTITIES; i++ {
		et.allocateID()
	}
	et.allocateID()
	et.expand(1)
	et.allocateID()
}

func TestEntityIDAllocatorReallocateDeallocatedID(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, NewIDGenerator())
	var e *Entity
	for i := 0; i < MAX_ENTITIES; i++ {
		allocated := et.allocateID()
		if i == MAX_ENTITIES/2 {
			e = allocated
		}
	}
	et.deallocate(e)
	e = et.allocateID()
	if e.ID != MAX_ENTITIES/2 {
		t.Fatal("should have used deallocated ID to serve new allocate request")
	}
}
