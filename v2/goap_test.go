package sameriver

import (
	"fmt"
	//	"time"

	"testing"
)

func printWorldState(ws *GOAPWorldState) {
	for name, val := range ws.vals {
		Logger.Printf("    %s: %d", name, val)
	}
}

func printGoal(g GOAPGoal) {
	for spec, goalVal := range g.goals {
		Logger.Printf("    %s: %d", spec, goalVal.val)
	}
}

func printDiffs(diffs map[string]int) {
	for name, diff := range diffs {
		Logger.Printf("    %s: %d", name, diff)
	}
}

func TestGOAPGoalRemaining(t *testing.T) {
	doTest := func(
		g GOAPGoal,
		ws *GOAPWorldState,
		nRemaining int,
		expectedRemaining []string,
	) {

		goalRemaining, diffs := g.goalRemaining(ws)

		Logger.Printf("goal:")
		printGoal(g)
		Logger.Printf("state:")
		printWorldState(ws)
		Logger.Printf("goalRemaining:")
		printGoal(goalRemaining)
		Logger.Printf("diffs:")
		printDiffs(diffs)
		Logger.Println("-------------------")

		if len(goalRemaining.goals) != nRemaining {
			t.Fatal(fmt.Sprintf("Should have had %d goals remaining, had %d", nRemaining, len(goalRemaining.goals)))
		}
		for _, name := range expectedRemaining {
			if diffVal, ok := diffs[name]; !ok || diffVal == 0 {
				t.Fatal(fmt.Sprintf("Should have had %s in diffs with value != 0", name))
			}
		}
	}

	doTest(
		NewGOAPGoal(map[string]int{
			"hasGlove,=": 1,
			"hasAxe,=":   1,
			"atTree,=":   1,
		}),
		NewGOAPWorldState(map[string]int{
			"hasGlove": 0,
			"hasAxe":   1,
			"atTree":   1,
		}),
		1,
		[]string{"hasGlove"},
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"hasGlove,=": 1,
			"hasAxe,=":   1,
			"atTree,=":   1,
		}),
		NewGOAPWorldState(map[string]int{
			"hasGlove": 1,
			"hasAxe":   1,
			"atTree":   1,
		}),
		0,
		[]string{},
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=": 3,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 1,
		}),
		1,
		[]string{"drunk"},
	)
}

func TestGOAPGoalStateCloserInSomeVar(t *testing.T) {
	doTest := func(g GOAPGoal, before, after *GOAPWorldState, expect bool) {
		if g.stateCloserInSomeVar(after, before) != expect {
			Logger.Println("Didn't get expected result for state B closer to goal than state A")
			Logger.Println("goal:")
			printGoal(g)
			Logger.Println("state A:")
			printWorldState(before)
			Logger.Println("state B:")
			printWorldState(after)
			t.Fatal("Didn't get expected result for stateCloserInSomeVar")
		}
	}

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=":       3,
			"rituallyPure,=": 1,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 1,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 1,
		}),
		false,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=":       3,
			"rituallyPure,=": 1,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 1,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 2,
		}),
		true,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=":       3,
			"rituallyPure,=": 1,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 1,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk":        1,
			"rituallyPure": 1,
		}),
		true,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=":       3,
			"rituallyPure,=": 1,
		}),
		NewGOAPWorldState(nil),
		NewGOAPWorldState(nil),
		false,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=":       3,
			"rituallyPure,=": 1,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 3,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 10,
		}),
		false,
	)

}

func TestGOAPActionPresFulfilled(t *testing.T) {

	eval := NewGOAPEvaluator()

	doTest := func(ws *GOAPWorldState, a *GOAPAction, expected bool) {
		if eval.presFulfilled(a, ws) != expected {
			Logger.Println("world state:")
			printWorldState(ws)
			Logger.Println("action.pres:")
			printGoal(a.pres)
			t.Fatal("Did not get expected value for action presfulfilled")
		}
	}

	// NOTE: both of these in reality should be modal
	goToAxe := NewGOAPAction(map[string]interface{}{
		"name": "goToAxe",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atAxe,=": 1,
		},
	})
	drink := NewGOAPAction(map[string]interface{}{
		"name": "drink",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,>": 0,
		},
		"effs": map[string]int{
			"hasBooze,-": 1,
		},
	})

	doDrinkTest := func(has int, expected bool) {
		doTest(
			NewGOAPWorldState(map[string]int{
				"hasBooze": has,
			}),
			drink,
			expected,
		)
	}
	chopTree := NewGOAPAction(map[string]interface{}{
		"name": "chopTree",
		"cost": 1,
		"pres": map[string]int{
			"hasGlove,>": 0,
			"hasAxe,>":   0,
			"atTree,=":   1,
		},
		"effs": map[string]int{
			"treeFelled,=": 1,
		},
	})

	eval.addActions(goToAxe, drink, chopTree)

	doDrinkTest(0, false)
	doDrinkTest(1, true)
	doDrinkTest(2, true)

	if !eval.presFulfilled(
		chopTree,
		NewGOAPWorldState(map[string]int{
			"hasGlove": 1,
			"hasAxe":   1,
			"atTree":   1,
		})) {
		t.Fatal("chopTree pres should have been fulfilled")
	}

	if eval.presFulfilled(
		chopTree,
		NewGOAPWorldState(map[string]int{
			"hasGlove": 1,
			"hasAxe":   1,
			"atTree":   0,
		})) {
		t.Fatal("chopTree pres shouldn't have been fulfilled")
	}
}

func TestGOAPActionModalVal(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)
	treePos := &Vec2D{11, 11}

	eval := NewGOAPEvaluator()

	atTreeModal := GOAPModalVal{
		name: "atTree",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*treePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		valAsEff: 1,
		effModalSet: func(ws *GOAPWorldState) {
			nearTree := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearTree)
		},
	}
	oceanPos := &Vec2D{500, 0}
	atOceanModal := GOAPModalVal{
		name: "atOcean",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*oceanPos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		valAsEff: 1,
		effModalSet: func(ws *GOAPWorldState) {
			nearOcean := oceanPos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearOcean)
		},
	}
	goToTree := NewGOAPActionModal(map[string]interface{}{
		"name":   "goToTree",
		"cost":   1,
		"pres":   nil,
		"checks": nil,
		"effs": map[string]GOAPStateVal{
			"atTree,=": atTreeModal,
		},
	})
	hugTree := NewGOAPActionModal(map[string]interface{}{
		"name": "hugTree",
		"cost": 1,
		// pres
		"pres": map[string]int{
			"atTree,=": 1,
		},
		// pre vars to be resolved by modalvals
		"checks": map[string]GOAPModalVal{
			"atTree": atTreeModal,
		},
		// effects
		"effs": map[string]GOAPStateVal{
			"connectionToNature,+": 2,
		},
	})
	goToOcean := NewGOAPActionModal(map[string]interface{}{
		"name":   "goToOcean",
		"cost":   1,
		"pres":   nil,
		"checks": nil,
		"effs": map[string]GOAPStateVal{
			"atOcean,=": atOceanModal,
		},
	})

	eval.addModalVals(atTreeModal, atOceanModal)
	eval.addActions(goToTree, hugTree, goToOcean)

	//
	// test presfulfilled
	//
	*e.GetVec2D("Position") = *treePos

	if !eval.presFulfilled(hugTree, ws) {
		t.Fatal("check result of atTreeModal should have returned 1, satisfying atTree,=: 1")
	}

	*e.GetVec2D("Position") = Vec2D{-100, -100}

	if eval.presFulfilled(hugTree, ws) {
		t.Fatal("check result of atTreeModal should have returned 0, failing to satisfy atTree,=: 1")
	}

	badWS := NewGOAPWorldState(map[string]int{
		"atTree": 0,
	})

	*e.GetVec2D("Position") = *treePos

	if !eval.presFulfilled(hugTree, badWS) {
		t.Fatal("regardless of what worldstate says, modal pre should decide and should've been true based on entity position = tree position")
	}

	//
	// test applyAction
	//

	g := NewGOAPGoal(map[string]int{
		"atTree,=": 1,
	})
	appliedState := eval.applyAction(goToTree, NewGOAPWorldState(nil))
	goalRemaining, diffs := g.goalRemaining(appliedState)
	Logger.Println("goal:")
	printGoal(g)
	Logger.Println("state after applying goToTree:")
	printWorldState(appliedState)
	if appliedState.vals["atTree"] != 1 {
		t.Fatal("atTree should've been 1 after goToTree")
	}
	Logger.Println("goal remaining:")
	printGoal(goalRemaining)
	Logger.Println("diffs:")
	printDiffs(diffs)

	//
	// test modal effect of applyAction
	//

	*e.GetVec2D("Position") = Vec2D{-100, -100}
	atTreeApplied := eval.applyAction(goToTree, NewGOAPWorldState(nil))
	Logger.Println("state after applying modal action eff of atTree:")
	printWorldState(atTreeApplied)
	if val, ok := atTreeApplied.vals["atTree"]; !ok || val != 1 {
		t.Fatal("Modal action eff should've set atTree=1")
	}
	Logger.Println("modal position of entity after modal action eff of atTree:")
	posAfter := atTreeApplied.GetModal(e, "Position").(*Vec2D)
	Logger.Printf("[%f, %f]", posAfter.X, posAfter.Y)

	//
	// test modal pre after modal set
	//

	*e.GetVec2D("Position") = Vec2D{-100, -100}
	atOceanApplied := eval.applyAction(goToOcean, NewGOAPWorldState(nil))
	Logger.Println("state after applying modal action eff of atOcean:")
	printWorldState(atOceanApplied)

	if eval.presFulfilled(hugTree, atOceanApplied) {
		t.Fatal("atTree modal pre of hugTree should fail when modal position is set at ocean")
	}

	nowGoToTreeApplied := eval.applyAction(goToTree, atOceanApplied)
	Logger.Println("state after goToOcean->goToTree:")
	printWorldState(nowGoToTreeApplied)
	if nowGoToTreeApplied.vals["atOcean"] != 0 {
		t.Fatal("Should've had atOcean=0 after goToTree")
	}
	if nowGoToTreeApplied.vals["atTree"] != 1 {
		t.Fatal("Should've had atTree=1 after goToTree")
	}

}

/*
	valAsEff: 1,
	effModalSet: func(ws *GOAPWorldState) {
		nearTree := treePos.Add(Vec2D{1, 0})
		ws.SetModal(e, "Position", nearTree)
	},

*/

/*


func TestGOAPWorldStateUnfulfilledByCtx(t *testing.T) {
	ws := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   true,
		},
	)
	hasAxeCtx := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
			return true
		},
	}
	getAxe := GOAPAction{
		name: "getAxe",
		pres: map[string]GOAPStateVal{
			"atAxe": true,
		},
		effs: map[string]GOAPStateVal{
			"hasAxe": hasAxeCtx,
		},
	}
	unfulfilled := ws.unfulfilledBy(getAxe)
	expected := map[string]bool{
		"hasGlove": true,
		"atTree":   true,
	}
	if len(unfulfilled.eals) != 2 {
		t.Fatal("unfulfilled should have length 2")
	}
	for name, val := range expected {
		if _, ok := unfulfilled.vals[name]; !ok {
			t.Fatal(fmt.Sprintf("%s not found in unfulfilled", name))
		}
		if unfulfilled.vals[name] != val {
			t.Fatal(fmt.Sprintf("%s should have been %t in unfulfilled", name, val))
		}
	}
}

func TestGOAPWorldStatePartlyCoversDoesntConflict(t *testing.T) {
	wsA := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   true,
		},
	)
	wsB := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"hasAxe": true,
		},
	)
	if !wsB.partlyCoversDoesntConflict(wsA) {
		t.Fatal("wsB should partly cover wsA")
	}
	wsC := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"hasAxe": false,
		},
	)
	if wsC.partlyCoversDoesntConflict(wsA) {
		t.Fatal("wsC should not partly cover wsA")
	}
	wsD := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"hasAxe": true,
			"atTree": false,
		},
	)
	if wsD.partlyCoversDoesntConflict(wsA) {
		t.Fatal("wsD should conflict with wsA")
	}
}

func TestGOAPWorldStateFulfillsSimple(t *testing.T) {
	wsA := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   false,
		},
	)
	wsB := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"hasAxe": true,
		},
	)
	if !wsA.fulfills(wsB) {
		t.Fatal("wsA should fulfill wsB")
	}
	wsC := NewGOAPWorldState(
		map[string]GOAPStateVal{
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
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atAxe": true,
		},
	}
	ws = ws.applyAction(goToAxe)
	goal := NewGOAPWorldState(
		map[string]GOAPStateVal{
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
		val: true,
		validate: func(ws GOAPWorldState) bool {
			return axeDistance < 5
		},
	}
	ws := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
	)
	goal := NewGOAPWorldState(
		map[string]GOAPStateVal{
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
		val: true,
		validate: func(ws GOAPWorldState) bool {
			return true
		},
		set: func(ws *GOAPWorldState) {
			ws.vals["atAxe"] = true
		},
	}
	goToAxe := GOAPAction{
		name: "goToAxe",
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
	}
	ws = ws.applyAction(goToAxe)
	goal := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
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
	e := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)
	axePos := Vec2D{11, 11}

	ctxAtAxe := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
			Logger.Println("in get...")
			ourPos := ws.GetModal(e, "Position").(Vec2D)
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
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
	}

	Logger.Println("applying action goToAxe...")
	ws = ws.applyAction(goToAxe)
	Logger.Println(ws.modal)

	goal := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
	)
	Logger.Println("testing if ws.fulfills(goal)")
	if !ws.fulfills(goal) {
		t.Fatal("ws should fulfill goal after applyAction with ctx val set()")
	}
}



func TestGOAPActionPresFulfilledCtx(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)
	axePos := Vec2D{11, 11}

	ctxAtAxe := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
			Logger.Println("in get...")
			ourPos := ws.GetModal(e, "Position").(Vec2D)
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
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
	}

	Logger.Println("applying action goToAxe...")
	ws = ws.applyAction(goToAxe)

	getAxe := GOAPAction{
		name: "getAxe",
		pres: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
		effs: map[string]GOAPStateVal{
			"hasAxe": true,
		},
	}

	if !getAxe.presFulfilled(ws) {
		t.Fatal("ws should have fulfilled the pres of getAxe after goToAxe")
	}
}

func TestGOAPPlannerBasic(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)

	axePos := Vec2D{11, 11}

	ctxAtAxe := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
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
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
	}
	getAxe := GOAPAction{
		name: "getAxe",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
		effs: map[string]GOAPStateVal{
			"hasAxe": true,
		},
	}

	p := NewGOAPPlanner(e)
	p.AddActions(goToAxe, getAxe)

	want := NewGOAPWorldState(
		map[string]GOAPStateVal{
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

func TestGOAPPlannerBasicMultiSuccess(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(map[string]GOAPStateVal{
		"bakeryHasBread":    true,
		"smokehouseHasFish": true,
	})

	bakeryPos := Vec2D{11, 11}
	smokehousePos := Vec2D{-11, 11}

	ctxAtBakery := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(Vec2D)
			_, _, d := ourPos.Distance(bakeryPos)
			return d < 2

		},
		set: func(ws *GOAPWorldState) {
			nearBakery := bakeryPos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearBakery)
		},
	}
	goToBakery := GOAPAction{
		name: "goToBakery",
		cost: 1,
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atBakery": ctxAtBakery,
		},
	}
	getBreadFromBakery := GOAPAction{
		name: "getBreadFromBakery",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"atBakery":       ctxAtBakery,
			"bakeryHasBread": true,
		},
		effs: map[string]GOAPStateVal{
			"hasBread": true,
			"hasFood":  true,
		},
	}

	ctxAtSmokehouse := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(Vec2D)
			_, _, d := ourPos.Distance(smokehousePos)
			return d < 2

		},
		set: func(ws *GOAPWorldState) {
			nearSmokehouse := smokehousePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearSmokehouse)
		},
	}
	goToSmokehouse := GOAPAction{
		name: "goToSmokehouse",
		cost: 1,
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atSmokehouse": ctxAtSmokehouse,
		},
	}
	getFishFromSmokehouse := GOAPAction{
		name: "getFishFromSmokehouse",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"atSmokehouse":      ctxAtSmokehouse,
			"smokehouseHasFish": true,
		},
		effs: map[string]GOAPStateVal{
			"hasFish": true,
			"hasFood": true,
		},
	}

	eat := GOAPAction{
		name: "eat",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"hasFood": true,
		},
		effs: map[string]GOAPStateVal{
			"sated": true,
		},
	}

	p := NewGOAPPlanner(e)
	p.AddActions(
		goToBakery,
		getBreadFromBakery,
		goToSmokehouse,
		getFishFromSmokehouse,
		eat,
	)

	want := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"sated": true,
		},
	)

	plans := p.Plans(ws, want)
	Logger.Printf("Found %d plans.", len(plans))

	Logger.Println("==========")
	Logger.Println("VALID PLANS:")
	for _, plan := range plans {
		Logger.Println(GOAPPlanToString(plan))
	}
	Logger.Println("==========")

	if len(plans) != 2 {
		t.Fatal("Should have found 2 plans (bakery, smokehouse)")
	}
}

func TestGOAPPlannerBasicMultiFailure(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)

	bakeryPos := Vec2D{11, 11}
	smokehousePos := Vec2D{-11, 11}

	ctxAtBakery := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(Vec2D)
			_, _, d := ourPos.Distance(bakeryPos)
			return d < 2

		},
		set: func(ws *GOAPWorldState) {
			nearBakery := bakeryPos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearBakery)
		},
	}
	goToBakery := GOAPAction{
		name: "goToBakery",
		cost: 1,
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atBakery": ctxAtBakery,
		},
	}
	getBreadFromBakery := GOAPAction{
		name: "getBreadFromBakery",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"atBakery":       ctxAtBakery,
			"bakeryHasBread": true,
		},
		effs: map[string]GOAPStateVal{
			"hasBread": true,
			"hasFood":  true,
		},
	}

	ctxAtSmokehouse := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
			ourPos := ws.GetModal(e, "Position").(Vec2D)
			_, _, d := ourPos.Distance(smokehousePos)
			return d < 2

		},
		set: func(ws *GOAPWorldState) {
			nearSmokehouse := smokehousePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", nearSmokehouse)
		},
	}
	goToSmokehouse := GOAPAction{
		name: "goToSmokehouse",
		cost: 1,
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atSmokehouse": ctxAtSmokehouse,
		},
	}
	getFishFromSmokehouse := GOAPAction{
		name: "getFishFromSmokehouse",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"atSmokehouse":      ctxAtSmokehouse,
			"smokehouseHasFish": true,
		},
		effs: map[string]GOAPStateVal{
			"hasFish": true,
			"hasFood": true,
		},
	}

	eat := GOAPAction{
		name: "eat",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"hasFood": true,
		},
		effs: map[string]GOAPStateVal{
			"sated": true,
		},
	}

	p := NewGOAPPlanner(e)
	p.AddActions(
		goToBakery,
		getBreadFromBakery,
		goToSmokehouse,
		getFishFromSmokehouse,
		eat,
	)

	want := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"sated": true,
		},
	)

	plans := p.Plans(ws, want)
	Logger.Printf("Found %d plans.", len(plans))

	Logger.Println("==========")
	Logger.Println("VALID PLANS:")
	for _, plan := range plans {
		Logger.Println(GOAPPlanToString(plan))
	}
	Logger.Println("==========")

	if len(plans) != 0 {
		t.Fatal("Should have found 0 plans (no bread in bakery or fish in smokehouse)")
	}
}

func TestGOAPPlannerHarder(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)

	axePos := Vec2D{11, 11}
	glovePos := Vec2D{2, 2}

	ctxAtAxe := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
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
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
	}
	getAxe := GOAPAction{
		name: "getAxe",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
		effs: map[string]GOAPStateVal{
			"hasAxe": true,
		},
	}

	ctxAtGlove := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
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
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atGlove": ctxAtGlove,
		},
	}
	getGlove := GOAPAction{
		name: "getGlove",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"atGlove": ctxAtGlove,
		},
		effs: map[string]GOAPStateVal{
			"hasGlove": true,
		},
	}

	p := NewGOAPPlanner(e)
	p.AddActions(goToAxe, getAxe, goToGlove, getGlove)

	want := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"hasAxe":   true,
			"hasGlove": true,
		},
	)

	plans := p.Plans(ws, want)
	Logger.Printf("Found %d plans.", len(plans))

	Logger.Println("==========")
	Logger.Println("VALID PLANS:")
	for _, plan := range plans {
		Logger.Println(GOAPPlanToString(plan))
	}
	Logger.Println("==========")

	if len(plans) != 2 {
		t.Fatal("Should have found 2 valid plans (glove,axe) or (axe,glove)")
	}

}

func TestGOAPPlannerHardest(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)

	axePos := Vec2D{11, 11}
	glovePos := Vec2D{2, 2}
	treePos := Vec2D{-7, -7}

	ctxAtAxe := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
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
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
	}
	getAxe := GOAPAction{
		name: "getAxe",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"atAxe": ctxAtAxe,
		},
		effs: map[string]GOAPStateVal{
			"hasAxe": true,
		},
	}

	ctxAtGlove := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
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
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atGlove": ctxAtGlove,
		},
	}
	getGlove := GOAPAction{
		name: "getGlove",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"atGlove": ctxAtGlove,
		},
		effs: map[string]GOAPStateVal{
			"hasGlove": true,
		},
	}

	ctxAtTree := GOAPCtxStateVal{
		val: true,
		validate: func(ws GOAPWorldState) bool {
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
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atTree": ctxAtTree,
		},
	}

	chopWood := GOAPAction{
		name: "chopWood",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"hasAxe":   true,
			"hasGlove": true,
			"atTree":   ctxAtTree,
		},
		effs: map[string]GOAPStateVal{
			"woodChopped": true,
		},
	}

	p := NewGOAPPlanner(e)
	p.AddActions(goToAxe, getAxe, goToGlove, getGlove, goToTree, chopWood)

	want := NewGOAPWorldState(
		map[string]GOAPStateVal{
			"woodChopped": true,
		},
	)

	t0 := time.Now()
	plans := p.Plans(ws, want)
	dt := time.Since(t0).Nanoseconds()
	Logger.Printf("Found %d plans from %d actions in %f ms", len(plans), len(p.actions.set), (float64(dt) / 1000000.0))

	Logger.Println("==========")
	Logger.Println("VALID PLANS:")
	for _, plan := range plans {
		Logger.Println(GOAPPlanToString(plan))
	}
	Logger.Println("==========")

	if len(plans) != 2 {
		t.Fatal("Should have found 2 valid plans (glove,axe,gototree) or (axe,glove,gototree)")
	}

}
*/

/*
func TestGOAPActionPresFulfilled(t *testing.T) {
	ws := NewGOAPWorldState(map[string]int{
		"atTree":   1,
		"hasAxe":   1,
		"hasGlove": 1,
	})

	chopTree := GOAPAction{
		name: "chopTree",
		pres: map[string]GOAPStateVal{
			"atTree":   1,
			"hasAxe":   1,
			"hasGlove": 1,
		},
		effs: map[string]GOAPStateVal{
			"woodChopped": 1,
		},
	}

	if !chopTree.presFulfilled(ws) {
		t.Fatal("ws should have fulfilled the pres of chopTree")
	}

	ws = NewGOAPWorldState(map[string]int{
		"atTree": 1,
	})

	if chopTree.presFulfilled(ws) {
		t.Fatal("ws should not have fulfilled pres of chopTree")
	}

	ws = NewGOAPWorldState(map[string]int{
		"atTree":   1,
		"hasAxe":   1,
		"hasGlove": 1,
		"drunk":    1,
	})

	if !chopTree.presFulfilled(ws) {
		t.Fatal("ws should have fulfilled the pres of chopTree")
	}
}

func TestGOAPWorldStateUnfulfilledBySimple(t *testing.T) {
	ws := NewGOAPWorldState(
		map[string]int{
			"hasAxe":   1,
			"hasGlove": 1,
			"atTree":   1,
		},
	)
	getAxe := GOAPAction{
		name: "getAxe",
		pres: map[string]GOAPStateVal{
			"atAxe": 1,
		},
		effs: map[string]GOAPStateVal{
			"hasAxe": 1,
		},
	}
	unfulfilled := ws.unfulfilledBy(getAxe)
	expected := map[string]int{
		"hasGlove": 1,
		"atTree":   1,
	}
	if len(unfulfilled.vals) != 2 {
		t.Fatal("unfulfilled should have length 2")
	}
	for name, val := range expected {
		if _, ok := unfulfilled.vals[name]; !ok {
			t.Fatal(fmt.Sprintf("%s not found in unfulfilled", name))
		}
		if unfulfilled.vals[name] != val {
			t.Fatal(fmt.Sprintf("%s should have been %d in unfulfilled", name, val))
		}
	}
}

func TestGOAPThoseThatHelpFulfill(t *testing.T) {
	getAxe := GOAPAction{
		name: "getAxe",
		pres: map[string]GOAPStateVal{
			"atAxe": 1,
		},
		effs: map[string]GOAPStateVal{
			"hasAxe": 1,
		},
	}
	goToAxe := GOAPAction{
		name: "goToAxe",
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atAxe": 1,
		},
	}
	getGlove := GOAPAction{
		name: "getGlove",
		pres: map[string]GOAPStateVal{
			"atAxe": 1,
		},
		effs: map[string]GOAPStateVal{
			"hasGlove": 1,
		},
	}
	goToTree := GOAPAction{
		name: "goToTree",
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"atTree": 1,
		},
	}
	candidates := NewGOAPActionSet()
	candidates.Add(getAxe, goToAxe, getGlove, goToTree)
	want := NewGOAPWorldState(map[string]int{
		"hasAxe":   1,
		"hasGlove": 1,
		"atTree":   1,
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
*/

/*
func TestGOAPPlannerPickUpDropNumeric(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	w.RegisterComponents([]string{"Generic,Inventory", "Generic,Disposition"})
	e, _ := w.Spawn([]string{}, MakeComponentSet(map[string]interface{}{
		"Generic,Inventory":   map[string]int{},
		"Generic,Disposition": map[string]int{},
	}))

	ws := NewGOAPWorldState(nil)

	getInventoryModal := func(ws GOAPWorldState) map[string]int {
		return (*ws.GetModal(e, "Inventory").(*interface{})).(map[string]int)
	}

	getDispositionModal := func(ws GOAPWorldState) map[string]int {
		return (*ws.GetModal(e, "Disposition").(*interface{})).(map[string]int)
	}

	ctxHasBooze := func(x int) GOAPCtxStateVal {
		return GOAPCtxStateVal{
			val: x,
			validate: func(ws GOAPWorldState) bool {
				inventory := getInventoryModal(ws)
				if n, ok := inventory["booze"]; ok {
					return n == x
				}
				return false
			},
			set: func(ws *GOAPWorldState) {
				inventory := getInventoryModal(*ws)
				inventory["booze"] += x
				ws.SetModal(e, "Inventory", inventory)
			},
		}
	}
	ctxDrunk := func(amount int) GOAPCtxStateVal {
		return GOAPCtxStateVal{
			// TODO: this seems to be the pattern that can distinguish a
			// setter from an incrementer
			val: func(val int) int {
				return val + 1
			},
			validate: func(ws GOAPWorldState) bool {
				disposition := getDispositionModal(ws)
				if _, ok := disposition["drunk"]; !ok {
					return false
				} else {
					return disposition["drunk"] >= amount
				}

			},
			set: func(ws *GOAPWorldState) {
				disposition := getDispositionModal(*ws)
				if _, ok := disposition["drunk"]; !ok {
					disposition["drunk"] = 1
				} else {
					disposition["drunk"] += 1
				}
				ws.SetModal(e, "Disposition", disposition)
			},
		}
	}

	getBooze := GOAPAction{
		name: "getBooze",
		cost: 1,
		pres: EmptyGOAPStateVal,
		effs: map[string]GOAPStateVal{
			"hasBooze": ctxHasBooze(1),
		},
	}
	drinkBooze := GOAPAction{
		name: "drinkBooze",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"hasBooze": 1,
		},
		effs: map[string]GOAPStateVal{
			"drunk":    ctxDrunk(1),
			"hasBooze": ctxHasBooze(-1),
		},
	}
	dropBooze := GOAPAction{
		name: "dropBooze",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"hasBooze": 1,
		},
		effs: map[string]GOAPStateVal{
			"hasBooze": ctxHasBooze(-1),
		},
	}
	purifyOneself := GOAPAction{
		name: "purifyOneself",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"hasBooze": 0,
		},
		effs: map[string]GOAPStateVal{
			"rituallyPure": 1,
		},
	}
	enterTemple := GOAPAction{
		name: "enterTemple",
		cost: 1,
		pres: map[string]GOAPStateVal{
			"rituallyPure": 1,
		},
		effs: map[string]GOAPStateVal{
			"admittedToTemple": 1,
		},
	}

	p := NewGOAPPlanner(e)
	p.AddActions(getBooze, drinkBooze, dropBooze, purifyOneself, enterTemple)

	want := NewGOAPWorldState(
		map[string]int{
			"drunk":            1,
			"admittedToTemple": 1,
		},
	)

	Logger.Println("Planning...")
	plans := p.Plans(ws, want)
	Logger.Printf("Found %d plans.", len(plans))

	expected := "[goToAxe,getAxe]"
	if GOAPPlanToString(plans[0]) != expected {
		t.Fatal("Did not find correct plan.")
	}
}
*/
