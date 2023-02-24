package sameriver

import (
	"testing"
)

func TestEntityMakeLogicUnit(t *testing.T) {
	w := testingWorld()
	e := w.Spawn(nil)
	lu := e.makeLogicUnit("loggyboi", func(dt_ms float64) {})
	if lu.name != e.LogicUnitName("loggyboi") {
		t.Fatal("did not set logic unit name")
	}
}
