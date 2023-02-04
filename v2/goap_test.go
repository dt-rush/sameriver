package sameriver

import (
	"testing"
)

func TestGOAPWorldStateIsSubset(t *testing.T) {
	wsA := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   false,
		},
	)
	wsB := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe": true,
		},
	)
	if !wsB.isSubset(wsA) {
		t.Fatal("wsB should be subset of wsA")
	}
	wsC := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe": false,
		},
	)
	if wsC.isSubset(wsA) {
		t.Fatal("wsC should not be subset of wsA")
	}
}

func TestGOAPWorldStateFulfillsSimple(t *testing.T) {
	wsA := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   false,
		},
	)
	wsB := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe": true,
		},
	)
	if !wsA.fulfills(wsB) {
		t.Fatal("wsA should fulfill wsB")
	}
	wsC := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   true,
		},
	)
	if wsA.fulfills(wsC) {
		t.Fatal("wsA should not fulfill wsC")
	}
}

func TestGOAPWorldStateApplyActionSimple(t *testing.T) {
	ws := NewGOAPWorldState(nil)
	goToAxe := GOAPAction{
		name: "goToAxe",
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atAxe": true,
		},
	}
	ws.applyAction(goToAxe)
	goal := NewGOAPWorldState(
		map[string]GOAPState{
			"atAxe": true,
		},
	)
	if !ws.fulfills(goal) {
		t.Fatal("ws should fulfill goal after applyAction")
	}
}

func TestGOAPWorldStateFulfillsCtx(t *testing.T) {
	axeDistance := 2
	ctxAtAxe := GOAPCtxStateVal{
		name: "atAxe",
		get: func() bool {
			return axeDistance < 5
		},
	}
	ws := NewGOAPWorldState(
		map[string]GOAPState{
			"atAxe": ctxAtAxe,
		},
	)
	goal := NewGOAPWorldState(
		map[string]GOAPState{
			"atAxe": true,
		},
	)
	if !ws.fulfills(goal) {
		t.Fatal("GOAPWorldState with value in map of type GOAPCtxStateVal should have worked")
	}
}

func TestGOAPWorldStateApplyActionCtx(t *testing.T) {
	ws := NewGOAPWorldState(nil)
	ctxAtAxe := GOAPCtxStateVal{
		name: "atAxe",
		get: func() bool {
			return true
		},
		set: func(ws *GOAPWorldState) {
			ws.Vals["atAxe"] = true
		},
	}
	goToAxe := GOAPAction{
		name: "goToAxe",
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atAxe": ctxAtAxe,
		},
	}
	ws.applyAction(goToAxe)
	goal := NewGOAPWorldState(
		map[string]GOAPState{
			"atAxe": true,
		},
	)
	if !ws.fulfills(goal) {
		t.Fatal("ws should fulfill goal after applyAction with ctx val set()")
	}
}

func TestGOAPWorldStateSetModal(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)
	axePos := Vec2D{11, 11}

	ctxAtAxe := GOAPCtxStateVal{
		name: "atAxe",
		get: func() bool {
			Logger.Println("in get...")
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(axePos)
			return d < 2
		},
		set: func(ws *GOAPWorldState) {
			Logger.Println("in set...")
			nearAxe := axePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearAxe)
		},
	}

	goToAxe := GOAPAction{
		name: "goToAxe",
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atAxe": ctxAtAxe,
		},
	}

	Logger.Println("applying action goToAxe...")
	ws.applyAction(goToAxe)
	Logger.Println(ws.modal)

	goal := NewGOAPWorldState(
		map[string]GOAPState{
			"atAxe": true,
		},
	)
	Logger.Println("testing if ws.fulfills(goal)")
	if !ws.fulfills(goal) {
		t.Fatal("ws should fulfill goal after applyAction with ctx val set()")
	}
}

func TestGOAPActionPresFulfilled(t *testing.T) {
	ws := NewGOAPWorldState(map[string]GOAPState{
		"atTree":   true,
		"hasAxe":   true,
		"hasGlove": true,
	})

	chopTree := GOAPAction{
		name: "chopTree",
		pres: map[string]GOAPState{
			"atTree":   true,
			"hasAxe":   true,
			"hasGlove": true,
		},
		effs: map[string]GOAPState{
			"woodChopped": true,
		},
	}

	if !chopTree.presFulfilled(ws) {
		t.Fatal("ws should have fulfilled the pres of chopTree")
	}
}

func TestGOAPActionPresFulfilledCtx(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)
	axePos := Vec2D{11, 11}

	ctxAtAxe := GOAPCtxStateVal{
		name: "atAxe",
		get: func() bool {
			Logger.Println("in get...")
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(axePos)
			return d < 2
		},
		set: func(ws *GOAPWorldState) {
			Logger.Println("in set...")
			nearAxe := axePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearAxe)
		},
	}

	goToAxe := GOAPAction{
		name: "goToAxe",
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atAxe": ctxAtAxe,
		},
	}

	Logger.Println("applying action goToAxe...")
	ws.applyAction(goToAxe)

	getAxe := GOAPAction{
		name: "getAxe",
		pres: map[string]GOAPState{
			"atAxe": true,
		},
		effs: map[string]GOAPState{
			"hasAxe": true,
		},
	}

	if !getAxe.presFulfilled(ws) {
		t.Fatal("ws should have fulfilled the pres of getAxe after goToAxe")
	}
}

func TestGOAPThoseThatHelpFulfill(t *testing.T) {
	getAxe := GOAPAction{
		name: "getAxe",
		pres: map[string]GOAPState{
			"atAxe": true,
		},
		effs: map[string]GOAPState{
			"hasAxe": true,
		},
	}
	goToAxe := GOAPAction{
		name: "goToAxe",
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atAxe": true,
		},
	}
	getGlove := GOAPAction{
		name: "getGlove",
		pres: map[string]GOAPState{
			"atAxe": true,
		},
		effs: map[string]GOAPState{
			"hasGlove": true,
		},
	}
	goToTree := GOAPAction{
		name: "goToTree",
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atTree": true,
		},
	}
	candidates := NewGOAPActionSet()
	candidates.Add(getAxe, goToAxe, getGlove, goToTree)
	want := NewGOAPWorldState(map[string]GOAPState{
		"hasAxe":   true,
		"hasGlove": true,
		"atTree":   true,
	})
	helpers := candidates.thoseThatHelpFulfill(want)

	helpersMatchExpected := func(expected []string) bool {
		valid := true
		for _, name := range expected {
			found := false
			for _, helper := range helpers.set {
				if helper.name == name {
					found = true
					break
				}
			}
			valid = valid && found
		}
		return valid
	}

	if !helpersMatchExpected([]string{"getAxe", "getGlove", "goToTree"}) {
		t.Fatal("Couldn't find expected fulfilling action in result of thoseThatFulfill")
	}

	if helpersMatchExpected([]string{"goToAxe"}) {
		t.Fatal("Should not have considered goToAxe as a helper to fulfill the goal")
	}

}
