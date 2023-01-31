package engine

import (
	"testing"
)

func TestEntityMakeLogicUnit(t *testing.T) {
	w := testingWorld()
	e, _ := w.Spawn([]string{}, ComponentSet{})
	lu := e.makeLogicUnit("loggyboi", func() {})
	if lu.name != e.LogicUnitName("loggyboi") {
		t.Fatal("did not set logic unit name")
	}
}
