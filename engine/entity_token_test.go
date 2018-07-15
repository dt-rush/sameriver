package engine

import (
	"testing"
)

func TestEntityTokenMakeLogicUnit(t *testing.T) {
	worldID := 2000
	name := "mylogic"
	e := EntityToken{WorldID: worldID}
	lu := e.MakeLogicUnit(name, func() {})
	if lu.Name != name {
		t.Fatal("did not set logic unit name")
	}
	if lu.WorldID != e.WorldID {
		t.Fatal("did not set logic unit world ID")
	}
}
