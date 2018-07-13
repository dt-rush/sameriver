package engine

import (
	"testing"
)

func TestEntityTableAllocateID(t *testing.T) {
	et := EntityTable{}
	_, err := et.allocateID()
	if !(err == nil &&
		et.numEntities == 1) {
		t.Fatal("didn't allocate ID properly")
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
