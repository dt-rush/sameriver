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
