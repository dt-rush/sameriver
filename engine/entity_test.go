package engine

import (
	"testing"
)

func TestEntityMakeLogicUnit(t *testing.T) {
	ID := 1000
	worldID := 2000
	e := Entity{ID: ID, WorldID: worldID}
	lu := e.MakeLogicUnit(func() {})
	if lu.name != e.LogicUnitName() {
		t.Fatal("did not set logic unit name")
	}
	if lu.worldID != e.WorldID {
		t.Fatal("did not set logic unit world ID")
	}
}
