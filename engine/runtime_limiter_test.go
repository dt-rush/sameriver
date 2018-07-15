package engine

import (
	"fmt"
	"testing"
)

func TestRuntimeLimiterAddLogic(t *testing.T) {
	r := NewRuntimeLimiter()
	for i := 0; i < 32; i++ {
		name := fmt.Sprintf("logic-%d", i)
		logic := &LogicUnit{
			Name:    name,
			WorldID: i,
			F:       func() {},
			Active:  true}
		r.Add(logic)
		if !(len(r.logicUnits) > 0 &&
			r.indexes[logic.WorldID] == len(r.logicUnits)-1) {
			t.Fatal("was not inserted properly")
		}
	}
}

func TestRuntimeLimiterRunLogic(t *testing.T) {
	r := NewRuntimeLimiter()
	x := 0
	name := "l1"
	r.Add(&LogicUnit{
		Name:    name,
		WorldID: 0,
		F:       func() { x += 1 },
		Active:  true})
	for i := 0; i < 32; i++ {
		r.Run(FRAME_SLEEP_MS)
	}
	if x != 32 {
		t.Fatal("didn't run logic")
	}
}

func TestRuntimeLimiterRemoveLogic(t *testing.T) {
	r := NewRuntimeLimiter()
	// test that we can remove a logic which doens't exist idempotently
	if r.Remove(0) != false {
		t.Fatal("somehow removed a logic which doesn't exist")
	}
	x := 0
	name := "l1"
	logic := &LogicUnit{
		Name:    name,
		WorldID: 0,
		F:       func() { x += 1 },
		Active:  true}
	r.Add(logic)
	// run logic a few times so that it has runtimeEstimate data
	for i := 0; i < 32; i++ {
		r.Run(FRAME_SLEEP_MS)
	}
	// remove it
	r.Remove(0)
	// test if removed
	if _, ok := r.runtimeEstimates[logic]; ok {
		t.Fatal("did not delete runtimeEstimates data")
	}
	if _, ok := r.indexes[logic.WorldID]; ok {
		t.Fatal("did not delete runtimeEstimates data")
	}
	if len(r.logicUnits) != 0 {
		t.Fatal("did not remove from logicUnits list")
	}
}
