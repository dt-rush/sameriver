package engine

import (
	"testing"
)

func TestEntityTableAllocateID(t *testing.T) {
	et := EntityTable{}
	_, err := et.allocateID()
	if !(err == nil &&
		et.n == 1) {
		t.Fatal("didn't allocate ID properly")
	}
}

func TestEntityTableDeallocateID(t *testing.T) {
	et := EntityTable{}
	e, _ := et.allocateID()
	et.deallocateID(e.ID)
	if et.n != 0 {
		t.Fatal("didn't update allocated count")
	}
	if !(len(et.availableIDs) == 1 && et.availableIDs[0] == e.ID) {
		t.Fatal("didn't add deallocated ID to list of available IDs")
	}
}

func TestEntityTableAllocateMaxIDs(t *testing.T) {
	et := EntityTable{}
	for i := 0; i < MAX_ENTITIES; i++ {
		et.allocateID()
	}
	_, err := et.allocateID()
	if err == nil {
		t.Fatal("should have returned error on allocating > MAX_ENTITIES")
	}
}

func TestEntityTableReallocateDeallocatedID(t *testing.T) {
	et := EntityTable{}
	for i := 0; i < MAX_ENTITIES; i++ {
		et.allocateID()
	}
	et.deallocateID(MAX_ENTITIES / 2)
	e, err := et.allocateID()
	if err != nil {
		t.Fatal("should have had space after deallocate to allocate")
	}
	if e.ID != MAX_ENTITIES/2 {
		t.Fatal("should have used deallocated ID to serve new allocate request")
	}
}
