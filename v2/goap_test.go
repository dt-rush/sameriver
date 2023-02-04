package sameriver

import (
	"fmt"

	"testing"
)

func TestGOAPWorldStateUnfulfilledBySimple(t *testing.T) {
	ws := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   true,
		},
	)
	getAxe := GOAPAction{
		name: "getAxe",
		pres: map[string]GOAPState{
			"atAxe": true,
		},
		effs: map[string]GOAPState{
			"hasAxe": true,
		},
	}
	unfulfilled := ws.unfulfilledBy(getAxe)
	expected := map[string]bool{
		"hasGlove": true,
		"atTree":   true,
	}
	if len(unfulfilled.Vals) != 2 {
		t.Fatal("unfulfilled should have length 2")
	}
	for name, val := range expected {
		if _, ok := unfulfilled.Vals[name]; !ok {
			t.Fatal(fmt.Sprintf("%s not found in unfulfilled", name))
		}
		if unfulfilled.Vals[name] != val {
			t.Fatal(fmt.Sprintf("%s should have been %t in unfulfilled", name, val))
		}
	}
}

func TestGOAPWorldStateUnfulfilledByCtx(t *testing.T) {
	ws := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   true,
		},
	)
	hasAxeCtx := GOAPCtxStateVal{
		name: "hasAxe",
		val:  true,
		get: func(ws GOAPWorldState) bool {
			return true
		},
	}
	getAxe := GOAPAction{
		name: "getAxe",
		pres: map[string]GOAPState{
			"atAxe": true,
		},
		effs: map[string]GOAPState{
			"hasAxe": hasAxeCtx,
		},
	}
	unfulfilled := ws.unfulfilledBy(getAxe)
	expected := map[string]bool{
		"hasGlove": true,
		"atTree":   true,
	}
	if len(unfulfilled.Vals) != 2 {
		t.Fatal("unfulfilled should have length 2")
	}
	for name, val := range expected {
		if _, ok := unfulfilled.Vals[name]; !ok {
			t.Fatal(fmt.Sprintf("%s not found in unfulfilled", name))
		}
		if unfulfilled.Vals[name] != val {
			t.Fatal(fmt.Sprintf("%s should have been %t in unfulfilled", name, val))
		}
	}
}

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
	ws = ws.applyAction(goToAxe)
	goal := NewGOAPWorldState(
		map[string]GOAPState{
			"atAxe": true,
		},
	)
	Logger.Println(ws)
	if !ws.fulfills(goal) {
		t.Fatal("ws should fulfill goal after applyAction")
	}
}

func TestGOAPWorldStateFulfillsCtx(t *testing.T) {
	axeDistance := 2
	ctxAtAxe := GOAPCtxStateVal{
		name: "atAxe",
		val:  true,
		get: func(ws GOAPWorldState) bool {
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
		val:  true,
		get: func(ws GOAPWorldState) bool {
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
	ws = ws.applyAction(goToAxe)
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
		val:  true,
		get: func(ws GOAPWorldState) bool {
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
	ws = ws.applyAction(goToAxe)
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
		val:  true,
		get: func(ws GOAPWorldState) bool {
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
	ws = ws.applyAction(goToAxe)

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

func TestGOAPPlannerBasic(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)

	axePos := Vec2D{11, 11}

	ctxAtAxe := GOAPCtxStateVal{
		name: "atAxe",
		val:  true,
		get: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(axePos)
			return d < 2

		},
		set: func(ws *GOAPWorldState) {
			nearAxe := axePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearAxe)
		},
	}
	goToAxe := GOAPAction{
		name: "goToAxe",
		cost: 1,
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atAxe": ctxAtAxe,
		},
	}
	getAxe := GOAPAction{
		name: "getAxe",
		cost: 1,
		pres: map[string]GOAPState{
			"atAxe": true,
		},
		effs: map[string]GOAPState{
			"hasAxe": true,
		},
	}

	p := NewGOAPPlanner(e)
	p.AddActions(goToAxe, getAxe)

	want := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe": true,
		},
	)

	plans := p.Plans(ws, want)
	Logger.Printf("Found %d plans.", len(plans))

	expected := "[goToAxe,getAxe]"
	if GOAPPlanToString(plans[0]) != expected {
		t.Fatal("Did not find correct plan.")
	}
}

func TestGOAPPlannerHarder(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)

	axePos := Vec2D{11, 11}
	glovePos := Vec2D{2, 2}

	ctxAtAxe := GOAPCtxStateVal{
		name: "atAxe",
		val:  true,
		get: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(Vec2D)
			_, _, d := ourPos.Distance(axePos)
			return d < 2
		},
		set: func(ws *GOAPWorldState) {
			nearAxe := axePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearAxe)
		},
	}
	goToAxe := GOAPAction{
		name: "goToAxe",
		cost: 1,
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atAxe": ctxAtAxe,
		},
	}
	getAxe := GOAPAction{
		name: "getAxe",
		cost: 1,
		pres: map[string]GOAPState{
			"atAxe": ctxAtAxe,
		},
		effs: map[string]GOAPState{
			"hasAxe": true,
		},
	}

	ctxAtGlove := GOAPCtxStateVal{
		name: "atGlove",
		val:  true,
		get: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(Vec2D)
			_, _, d := ourPos.Distance(glovePos)
			return d < 2
		},
		set: func(ws *GOAPWorldState) {
			nearGlove := glovePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearGlove)
		},
	}
	goToGlove := GOAPAction{
		name: "goToGlove",
		cost: 1,
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atGlove": ctxAtGlove,
		},
	}
	getGlove := GOAPAction{
		name: "getGlove",
		cost: 1,
		pres: map[string]GOAPState{
			"atGlove": ctxAtGlove,
		},
		effs: map[string]GOAPState{
			"hasGlove": true,
		},
	}

	p := NewGOAPPlanner(e)
	p.AddActions(goToAxe, getAxe, goToGlove, getGlove)

	want := NewGOAPWorldState(
		map[string]GOAPState{
			"hasAxe":   true,
			"hasGlove": true,
		},
	)

	plans := p.Plans(ws, want)
	Logger.Printf("Found %d plans.", len(plans))

}

func TestGOAPPlannerHardest(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)

	axePos := Vec2D{11, 11}
	glovePos := Vec2D{2, 2}
	treePos := Vec2D{-7, -7}

	ctxAtAxe := GOAPCtxStateVal{
		name: "atAxe",
		val:  true,
		get: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(Vec2D)
			_, _, d := ourPos.Distance(axePos)
			return d < 2
		},
		set: func(ws *GOAPWorldState) {
			nearAxe := axePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearAxe)
		},
	}
	goToAxe := GOAPAction{
		name: "goToAxe",
		cost: 1,
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atAxe": ctxAtAxe,
		},
	}
	getAxe := GOAPAction{
		name: "getAxe",
		cost: 1,
		pres: map[string]GOAPState{
			"atAxe": ctxAtAxe,
		},
		effs: map[string]GOAPState{
			"hasAxe": true,
		},
	}

	ctxAtGlove := GOAPCtxStateVal{
		name: "atGlove",
		val:  true,
		get: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(Vec2D)
			_, _, d := ourPos.Distance(glovePos)
			return d < 2
		},
		set: func(ws *GOAPWorldState) {
			nearGlove := glovePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearGlove)
		},
	}
	goToGlove := GOAPAction{
		name: "goToGlove",
		cost: 1,
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atGlove": ctxAtGlove,
		},
	}
	getGlove := GOAPAction{
		name: "getGlove",
		cost: 1,
		pres: map[string]GOAPState{
			"atGlove": ctxAtGlove,
		},
		effs: map[string]GOAPState{
			"hasGlove": true,
		},
	}

	ctxAtTree := GOAPCtxStateVal{
		name: "atTree",
		val:  true,
		get: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(Vec2D)
			_, _, d := ourPos.Distance(treePos)
			return d < 2
		},
		set: func(ws *GOAPWorldState) {
			nearTree := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearTree)
		},
	}
	goToTree := GOAPAction{
		name: "goToTree",
		cost: 1,
		pres: EmptyGOAPState,
		effs: map[string]GOAPState{
			"atTree": ctxAtTree,
		},
	}

	chopWood := GOAPAction{
		name: "chopWood",
		cost: 1,
		pres: map[string]GOAPState{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   ctxAtTree,
		},
		effs: map[string]GOAPState{
			"woodChopped": true,
		},
	}

	p := NewGOAPPlanner(e)
	p.AddActions(goToAxe, getAxe, goToGlove, getGlove, goToTree, chopWood)

	want := NewGOAPWorldState(
		map[string]GOAPState{
			"woodChopped": true,
		},
	)

	plans := p.Plans(ws, want)
	Logger.Printf("Found %d plans.", len(plans))

}
