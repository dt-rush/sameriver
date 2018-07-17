package engine

import (
	"testing"
)

func TestEntityTokenMakeLogicUnit(t *testing.T) {
	ID := 1000
	worldID := 2000
	e := EntityToken{ID: ID, WorldID: worldID}
	lu := e.MakeLogicUnit(func() {})
	if lu.Name != e.LogicUnitName() {
		t.Fatal("did not set logic unit name")
	}
	if lu.WorldID != e.WorldID {
		t.Fatal("did not set logic unit world ID")
	}
}
