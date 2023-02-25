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

func TestEntityInvalidComponentAccess(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Should have paniced")
		}
	}()
	w := testingWorld()
	e := w.Spawn(nil)
	e.GetVec2D("Doesntexist")
}
