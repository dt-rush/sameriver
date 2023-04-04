package sameriver

import (
	"fmt"
	"strings"
	"time"

	"testing"

	"github.com/TwiN/go-color"
)

func printWorldState(ws *GOAPWorldState) {
	if ws == nil || len(ws.vals) == 0 {
		Logger.Println("    nil")
		return
	}
	for name, val := range ws.vals {
		Logger.Printf("    %s: %d", name, val)
	}
}

func printGoal(g *GOAPGoal) {
	if g == nil || len(g.vars) == 0 {
		Logger.Println("    nil")
		return
	}
	for varName, interval := range g.vars {
		Logger.Printf("    want %s: [%.0f, %.0f]", varName, interval.A, interval.B)
	}
}

func printGoalRemaining(g *GOAPGoalRemaining) {
	if g.nUnfulfilled == 0 {
		msg := "    satisfied    "
		Logger.Printf(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		Logger.Printf(color.InBlackOverGreen(msg))
		Logger.Printf(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		return
	}
	for varName, interval := range g.goalLeft {
		msg := fmt.Sprintf("    %s: [%.0f, %.0f]    ", varName, interval.A, interval.B)

		Logger.Printf(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))
		Logger.Printf(color.InBold(color.InRedOverBlack(msg)))
		Logger.Printf(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))

	}
}

func printGoalRemainingSurface(s *GOAPGoalRemainingSurface) {
	if s.NUnfulfilled() == 0 {
		Logger.Println("    nil")
	} else {
		for i, tgs := range s.surface {
			if i == len(s.surface)-1 {
				Logger.Printf(color.InBold(color.InRedOverGray("main:")))

			}
			for _, tg := range tgs {
				printGoalRemaining(tg)
			}
		}
	}
}

func printDiffs(diffs map[string]float64) {
	for name, diff := range diffs {
		Logger.Printf("    %s: %.0f", name, diff)
	}
}

func TestGOAPGoalRemaining(t *testing.T) {
	doTest := func(
		g *GOAPGoal,
		ws *GOAPWorldState,
		nRemaining int,
		expectedRemaining []string,
	) {

		remaining := g.remaining(ws)

		Logger.Printf("goal:")
		printGoal(g)
		Logger.Printf("state:")
		printWorldState(ws)
		Logger.Printf("remaining:")
		printGoal(remaining.goal)
		Logger.Printf("diffs:")
		printDiffs(remaining.diffs)
		Logger.Println("-------------------")

		if len(remaining.goalLeft) != nRemaining {
			t.Fatalf("Should have had %d goals remaining, had %d", nRemaining, len(remaining.goalLeft))
		}
		for _, name := range expectedRemaining {
			if diffVal, ok := remaining.diffs[name]; !ok || diffVal == 0 {
				t.Fatalf("Should have had %s in diffs with value != 0", name)
			}
		}
	}

	doTest(
		newGOAPGoal(map[string]int{
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
		newGOAPGoal(map[string]int{
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
		newGOAPGoal(map[string]int{
			"drunk,>=": 3,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 1,
		}),
		1,
		[]string{"drunk"},
	)
}

func TestGOAPGoalRemainingsOfPath(t *testing.T) {
	w := testingWorld()
	const (
		BOOZEAMOUNT = GENERICTAGS + 1 + iota
	)
	w.RegisterComponents([]any{
		BOOZEAMOUNT, INT, "BOOZEAMOUNT",
	})

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			BOOZEAMOUNT: 0,
		},
	})
	Logger.Println(e)

	p := NewGOAPPlanner(e)

	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			amount := ws.GetModal(e, BOOZEAMOUNT).(*int)
			return *amount
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			amount := ws.GetModal(e, BOOZEAMOUNT).(*int)
			if op == "-" {
				newVal := *amount - x
				ws.SetModal(e, BOOZEAMOUNT, &newVal)
			}
		},
	}
	drink := NewGOAPAction(map[string]any{
		"name": "drink",
		"cost": 1,
		"pres": map[string]int{
			"EACH:hasBooze,>=": 1,
		},
		"effs": map[string]int{
			"drunk,+":    1,
			"hasBooze,-": 1,
		},
	})

	p.AddModalVals(hasBoozeModal)
	p.AddActions(drink)

	start := NewGOAPWorldState(nil)
	start.w = w // this would be done automatically in Plan()
	p.checkModalInto("hasBooze", start)

	goal := map[string]int{
		"drunk,>=": 3,
	}

	path := NewGOAPPath([]*GOAPAction{drink.Parametrized(2)})

	Logger.Printf("-------------------------------------------- 1")

	p.computeCostAndRemainingsOfPath(path, start, NewGOAPTemporalGoal(goal))

	Logger.Printf("%d unfulfilled", path.remainings.NUnfulfilled())
	printGoalRemainingSurface(path.remainings)
	mainGoalRemaining := path.remainings.surface[len(path.remainings.surface)-1][0]
	if path.remainings.NUnfulfilled() != 2 || len(mainGoalRemaining.goalLeft) != 1 {
		t.Fatal("Remaining was not calculated properly")
	}

	Logger.Printf("-------------------------------------------- 2")

	path = NewGOAPPath([]*GOAPAction{drink.Parametrized(3)})

	p.computeCostAndRemainingsOfPath(path, start, NewGOAPTemporalGoal(goal))

	Logger.Printf("%d unfulfilled", path.remainings.NUnfulfilled())
	printGoalRemainingSurface(path.remainings)
	mainGoalRemaining = path.remainings.surface[len(path.remainings.surface)-1][0]
	if path.remainings.NUnfulfilled() != 1 || len(mainGoalRemaining.goalLeft) != 0 {
		t.Fatal("Remaining was not calculated properly")
	}

	Logger.Printf("-------------------------------------------- 3")

	booze := e.GetInt(BOOZEAMOUNT)
	*booze = 3

	p.checkModalInto("hasBooze", start)

	Logger.Printf("start: %v", start.vals)

	p.computeCostAndRemainingsOfPath(path, start, NewGOAPTemporalGoal(goal))

	Logger.Printf("%d unfulfilled", path.remainings.NUnfulfilled())
	printGoalRemainingSurface(path.remainings)
	if path.remainings.NUnfulfilled() != 0 || len(mainGoalRemaining.goalLeft) != 0 {
		t.Fatal("Remaining was not calculated properly")
	}
}

func TestGOAPActionPresFulfilled(t *testing.T) {

	w := testingWorld()
	e := w.Spawn(nil)
	p := NewGOAPPlanner(e)

	doTest := func(ws *GOAPWorldState, a *GOAPAction, expected bool) {
		if p.presFulfilled(a, ws) != expected {
			Logger.Println("world state:")
			printWorldState(ws)
			Logger.Println("action.pres:")
			for _, tg := range a.pres.temporalGoals {
				printGoal(tg)
			}
			t.Fatal("Did not get expected value for action presfulfilled")
		}
	}

	// NOTE: both of these in reality should be modal
	goToAxe := NewGOAPAction(map[string]any{
		"name": "goToAxe",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atAxe,=": 1,
		},
	})
	drink := NewGOAPAction(map[string]any{
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
	chopTree := NewGOAPAction(map[string]any{
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

	p.AddActions(goToAxe, drink, chopTree)

	doDrinkTest(0, false)
	doDrinkTest(1, true)
	doDrinkTest(2, true)

	if !p.presFulfilled(
		chopTree,
		NewGOAPWorldState(map[string]int{
			"hasGlove": 1,
			"hasAxe":   1,
			"atTree":   1,
		})) {
		t.Fatal("chopTree pres should have been fulfilled")
	}

	if p.presFulfilled(
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
	e := testingSpawnPhysics(w)

	nilWS := NewGOAPWorldState(nil)
	nilWS.w = w // this would be done automatically in Plan()

	ws := nilWS.CopyOf()
	treePos := &Vec2D{11, 11}

	p := NewGOAPPlanner(e)

	atTreeModal := GOAPModalVal{
		name: "atTree",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, POSITION).(*Vec2D)
			_, _, d := ourPos.Distance(*treePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearTree := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, POSITION, &nearTree)
		},
	}
	oceanPos := &Vec2D{500, 0}
	atOceanModal := GOAPModalVal{
		name: "atOcean",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, POSITION).(*Vec2D)
			_, _, d := ourPos.Distance(*oceanPos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearOcean := oceanPos.Add(Vec2D{1, 0})
			ws.SetModal(e, POSITION, &nearOcean)
		},
	}
	goToTree := NewGOAPAction(map[string]any{
		"name": "goToTree",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atTree,=": 1,
		},
	})
	chopTree := NewGOAPAction(map[string]any{
		"name": "chopTree",
		"cost": 1,
		"pres": map[string]int{
			"atTree,=": 1,
			"hasAxe,>": 0,
		},
		"effs": map[string]int{
			"woodChopped,+": 1,
		},
	})
	hugTree := NewGOAPAction(map[string]any{
		"name": "hugTree",
		"cost": 1,
		"pres": map[string]int{
			"atTree,=": 1,
		},
		"effs": map[string]int{
			"connectionToNature,+": 2,
		},
	})
	goToOcean := NewGOAPAction(map[string]any{
		"name": "goToOcean",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atOcean,=": 1,
		},
	})

	p.AddModalVals(atTreeModal, atOceanModal)
	p.AddActions(goToTree, hugTree, chopTree, goToOcean)

	//
	// test presfulfilled
	//
	*e.GetVec2D(POSITION) = *treePos

	if !p.presFulfilled(hugTree, ws) {
		t.Fatal("check result of atTreeModal should have returned 1, satisfying atTree,=: 1")
	}

	*e.GetVec2D(POSITION) = Vec2D{-100, -100}

	if p.presFulfilled(hugTree, ws) {
		t.Fatal("check result of atTreeModal should have returned 0, failing to satisfy atTree,=: 1")
	}

	badWS := NewGOAPWorldState(map[string]int{
		"atTree": 0,
	})
	badWS.w = w

	*e.GetVec2D(POSITION) = *treePos

	if !p.presFulfilled(hugTree, badWS) {
		t.Fatal("regardless of what worldstate says, modal pre should decide and should've been true based on entity position = tree position")
	}

	axeWS := NewGOAPWorldState(map[string]int{
		"hasAxe": 1,
	})
	axeWS.w = w
	if !p.presFulfilled(chopTree, axeWS) {
		t.Fatal("mix of modal and basic world state vals should fulfill pre")
	}

	//
	// test applyAction
	//

	g := newGOAPGoal(map[string]int{
		"atTree,=": 1,
	})
	appliedState := p.applyActionBasic(goToTree, nilWS, true)
	remaining := g.remaining(appliedState)
	Logger.Println("goal:")
	printGoal(g)
	Logger.Println("state after applying goToTree:")
	printWorldState(appliedState)
	if appliedState.vals["atTree"] != 1 {
		t.Fatal("atTree should've been 1 after goToTree")
	}
	Logger.Println("goal remaining:")
	printGoal(remaining.goal)
	if len(remaining.goalLeft) != 0 {
		t.Fatal("Goal should have been satisfied")
	}
	Logger.Println("diffs:")
	printDiffs(remaining.diffs)

	g2 := newGOAPGoal(map[string]int{
		"atTree,=": 1,
		"drunk,>=": 10,
	})
	remaining = g2.remaining(appliedState)
	if len(remaining.goalLeft) != 1 {
		t.Fatal("drunk goal should be unfulfilled by atTree state")
	}

	//
	// test modal effect of applyAction
	//

	*e.GetVec2D(POSITION) = Vec2D{-100, -100}

	atTreeApplied, _ := p.applyActionModal(goToTree, nilWS)
	Logger.Println("state after applying modal action eff of atTree:")
	printWorldState(atTreeApplied)
	if val, ok := atTreeApplied.vals["atTree"]; !ok || val != 1 {
		t.Fatal("Modal action eff should've set atTree=1")
	}
	Logger.Println("modal position of entity after modal action eff of atTree:")
	posAfter := atTreeApplied.GetModal(e, POSITION).(*Vec2D)
	Logger.Printf("[%f, %f]", posAfter.X, posAfter.Y)

	//
	// test modal pre after modal set
	//

	*e.GetVec2D(POSITION) = Vec2D{-100, -100}
	atOceanApplied, _ := p.applyActionModal(goToOcean, nilWS)
	Logger.Println("state after applying modal action eff of atOcean:")
	printWorldState(atOceanApplied)

	if p.presFulfilled(hugTree, atOceanApplied) {
		t.Fatal("atTree modal pre of hugTree should fail when modal position is set at ocean")
	}

	nowGoToTreeApplied, _ := p.applyActionModal(goToTree, atOceanApplied)
	Logger.Println("state after goToOcean->goToTree:")
	printWorldState(nowGoToTreeApplied)
	if nowGoToTreeApplied.vals["atOcean"] != 0 {
		t.Fatal("Should've had atOcean=0 after goToTree")
	}
	if nowGoToTreeApplied.vals["atTree"] != 1 {
		t.Fatal("Should've had atTree=1 after goToTree")
	}

}

func TestGOAPPlanSimple(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)

	treePos := &Vec2D{19, 19}

	atTreeModal := GOAPModalVal{
		name: "atTree",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, POSITION).(*Vec2D)
			_, _, d := ourPos.Distance(*treePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearTree := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, POSITION, &nearTree)
		},
	}
	goToTree := NewGOAPAction(map[string]any{
		"name": "goToTree",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atTree,=": 1,
		},
	})

	goal := map[string]int{
		"atTree,=": 1,
	}

	Logger.Println(*e.GetVec2D(POSITION))

	ws := NewGOAPWorldState(nil)

	p := NewGOAPPlanner(e)
	p.AddModalVals(atTreeModal)
	p.AddActions(goToTree)

	Logger.Println(p.Plan(ws, goal, 50))

}

func TestGOAPPlanSimpleIota(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)

	const (
		STATE = GENERICTAGS + 1 + iota
	)

	w.RegisterComponents([]any{
		STATE, INTMAP, "STATE",
	})

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			STATE: map[string]int{
				"drunk": 0,
			},
		},
	})

	drunkModal := GOAPModalVal{
		name: "drunk",
		check: func(ws *GOAPWorldState) int {
			state := ws.GetModal(e, STATE).(*IntMap)
			return state.m["drunk"]
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			state := ws.GetModal(e, STATE).(*IntMap).CopyOf()
			if op == "+" {
				state.m["drunk"] += x
			}
			ws.SetModal(e, STATE, &state)
		},
	}
	drink := NewGOAPAction(map[string]any{
		"name": "drink",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"drunk,+": 1,
		},
	})

	goal := newGOAPGoal(map[string]int{
		"drunk,=": 1,
	})

	ws := NewGOAPWorldState(nil)

	p := NewGOAPPlanner(e)
	p.AddModalVals(drunkModal)
	p.AddActions(drink)

	Logger.Println(p.Plan(ws, goal, 50))

	goal = newGOAPGoal(map[string]int{
		"drunk,=": 3,
	})
	Logger.Println(p.Plan(ws, goal, 50))

}

func TestGOAPPlanSimpleEnough(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)

	const (
		STATE = GENERICTAGS + 1 + iota
	)

	w.RegisterComponents([]any{
		STATE, INTMAP, "STATE",
	})

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			STATE: map[string]int{
				"drunk": 0,
			},
		},
	})

	drunkModal := GOAPModalVal{
		name: "drunk",
		check: func(ws *GOAPWorldState) int {
			state := ws.GetModal(e, STATE).(*IntMap)
			return state.m["drunk"]
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			state := ws.GetModal(e, STATE).(*IntMap).CopyOf()
			if op == "+" {
				state.m["drunk"] += x
			}
			ws.SetModal(e, STATE, &state)
		},
	}
	drink := NewGOAPAction(map[string]any{
		"name": "drink",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"drunk,+": 1,
		},
	})
	purifyOneself := NewGOAPAction(map[string]any{
		"name": "purifyOneself",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"rituallyPure,=": 1,
		},
	})

	ws := NewGOAPWorldState(nil)

	p := NewGOAPPlanner(e)
	p.AddModalVals(drunkModal)
	p.AddActions(drink, purifyOneself)

	goal := newGOAPGoal(map[string]int{
		"drunk,=":        10,
		"rituallyPure,=": 1,
	})
	Logger.Println(p.Plan(ws, goal, 50))
}

func TestGOAPPlanAlanWatts(t *testing.T) {
	w := testingWorld()

	ps := NewPhysicsSystem()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(ps, items, inventories)

	const (
		STATE = GENERICTAGS + 1 + iota
	)

	w.RegisterComponents([]any{
		STATE, INTMAP, "STATE",
	})

	items.CreateArchetype(map[string]any{
		"name":        "bottle_booze",
		"displayName": "bottle of booze",
		"flavourText": "a potent brew!",
		"properties": map[string]int{
			"value":     10,
			"drunkness": 2,
		},
		"tags": []string{"booze"},
	})

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION: Vec2D{0, 0},
			STATE: map[string]int{
				"drunk": 0,
			},
			INVENTORY: inventories.Create(nil),
		},
	})

	boozePos := &Vec2D{19, 19}
	templePos := &Vec2D{-19, 19}

	inTempleModal := GOAPModalVal{
		name: "inTemple",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, POSITION).(*Vec2D)
			_, _, d := ourPos.Distance(*templePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearTemple := templePos.Add(Vec2D{1, 0})
			ws.SetModal(e, POSITION, &nearTemple)
		},
	}
	atBoozeModal := GOAPModalVal{
		name: "atBooze",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, POSITION).(*Vec2D)
			_, _, d := ourPos.Distance(*boozePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearBooze := boozePos.Add(Vec2D{1, 0})
			ws.SetModal(e, POSITION, &nearBooze)
		},
	}
	drunkModal := GOAPModalVal{
		name: "drunk",
		check: func(ws *GOAPWorldState) int {
			state := ws.GetModal(e, STATE).(*IntMap)
			return state.m["drunk"]
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			state := ws.GetModal(e, STATE).(*IntMap).CopyOf()
			if op == "+" {
				state.m["drunk"] += x
			}
			ws.SetModal(e, STATE, &state)
		},
	}
	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			inv := ws.GetModal(e, INVENTORY).(*Inventory)
			count := inv.CountTags("booze")
			return count
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			inv := ws.GetModal(e, INVENTORY).(*Inventory).CopyOf()
			if op == "-" {
				inv.DebitNTags(x, "booze")
			}
			if op == "=" {
				if x == 0 {
					inv.DebitAllTags("booze")
				}
				count := inv.CountTags("booze")
				if count == 0 {
					inv.Credit(items.CreateStackSimple(1, "bottle_booze"))
				}
				inv.SetCountTags(x, "booze")
			}
			if op == "+" {
				count := inv.CountTags("booze")
				if count == 0 {
					inv.Credit(items.CreateStackSimple(x, "bottle_booze"))
				} else {
					inv.SetCountTags(count+x, "booze")
				}
			}
			ws.SetModal(e, INVENTORY, inv)
		},
	}
	goToBooze := NewGOAPAction(map[string]any{
		"name": "goToBooze",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atBooze,=": 1,
		},
	})
	getBooze := NewGOAPAction(map[string]any{
		"name": "getBooze",
		"cost": 1,
		"pres": map[string]int{
			"atBooze,=": 1,
		},
		"effs": map[string]int{
			"hasBooze,+": 1,
		},
	})
	drink := NewGOAPAction(map[string]any{
		"name": "drink",
		"cost": 1,
		"pres": map[string]int{
			"EACH:hasBooze,>=": 1,
		},
		"effs": map[string]int{
			"drunk,+":    2,
			"hasBooze,-": 1,
		},
	})
	dropAllBooze := NewGOAPAction(map[string]any{
		"name": "dropAllBooze",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"hasBooze,=": 0,
		},
	})
	purifyOneself := NewGOAPAction(map[string]any{
		"name": "purifyOneself",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,<": 1,
		},
		"effs": map[string]int{
			"rituallyPure,=": 1,
		},
	})
	enterTemple := NewGOAPAction(map[string]any{
		"name": "enterTemple",
		"cost": 1,
		"pres": map[string]int{
			"rituallyPure,=": 1,
		},
		"effs": map[string]int{
			"inTemple,=": 1,
		},
	})

	p := NewGOAPPlanner(e)

	p.AddModalVals(drunkModal, hasBoozeModal, atBoozeModal, inTempleModal)
	p.AddActions(drink, dropAllBooze, purifyOneself, enterTemple, goToBooze, getBooze)

	ws := NewGOAPWorldState(map[string]int{
		"rituallyPure": 0,
	})

	goal := []any{
		map[string]int{
			"drunk,>=": 3,
		},
		map[string]int{
			"inTemple,=": 1,
		},
	}
	t0 := time.Now()
	plan, ok := p.Plan(ws, goal, 500)
	if !ok {
		t.Fatal("Should've found a solution")
	}
	Logger.Println(color.InGreenOverWhite(GOAPPathToString(plan)))
	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	Logger.Printf("Took %f ms to find solution", dt_ms)

	e.SetGeneric(INVENTORY, inventories.Create(map[string]int{
		"bottle_booze": 10,
	}))
	t0 = time.Now()
	plan, ok = p.Plan(ws, goal, 500)
	if !ok {
		t.Fatal("Should've found a solution")
	}
	Logger.Println(color.InGreenOverWhite(GOAPPathToString(plan)))
	dt_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
	Logger.Printf("Took %f ms to find solution", dt_ms)
}

func TestGOAPPlanClassic(t *testing.T) {
	w := testingWorld()

	ps := NewPhysicsSystem()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(ps, items, inventories)

	items.CreateArchetype(map[string]any{
		"name":        "axe",
		"displayName": "axe",
		"flavourText": "a nice axe for chopping trees",
		"properties": map[string]int{
			"value":     10,
			"sharpness": 2,
		},
		"tags": []string{"tool"},
	})
	items.CreateArchetype(map[string]any{
		"name":        "glove",
		"displayName": "glove",
		"flavourText": "good hand protection",
		"properties": map[string]int{
			"value": 2,
		},
		"tags": []string{"wearable"},
	})

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION:  Vec2D{0, 0},
			INVENTORY: inventories.Create(nil),
		},
	})

	hasModal := func(name string, archetype string, tags ...string) GOAPModalVal {
		return GOAPModalVal{
			name: fmt.Sprintf("has%s", name),
			check: func(ws *GOAPWorldState) int {
				inv := ws.GetModal(e, INVENTORY).(*Inventory)
				return inv.CountName(archetype)
			},
			effModalSet: func(ws *GOAPWorldState, op string, x int) {
				inv := ws.GetModal(e, INVENTORY).(*Inventory).CopyOf()
				if op == "-" {
					inv.DebitNTags(x, archetype)
				}
				if op == "=" {
					count := inv.CountTags(tags...)
					if count == 0 {
						inv.Credit(items.CreateStackSimple(x, archetype))
					} else {
						inv.SetCountName(x, archetype)
					}
				}
				if op == "+" {
					count := inv.CountName(archetype)
					if count == 0 {
						inv.Credit(items.CreateStackSimple(x, archetype))
					} else {
						inv.SetCountName(count+x, archetype)
					}
				}
				ws.SetModal(e, INVENTORY, inv)
			},
		}
	}

	hasAxeModal := hasModal("Axe", "axe")
	hasGloveModal := hasModal("Glove", "glove")

	get := func(name string) *GOAPAction {
		return NewGOAPAction(map[string]any{
			"name": fmt.Sprintf("get%s", name),
			"cost": 1,
			"pres": map[string]int{
				fmt.Sprintf("at%s,=", name): 1,
			},
			"effs": map[string]int{
				fmt.Sprintf("has%s,+", name): 1,
			},
		})
	}

	getAxe := get("Axe")
	getGlove := get("Glove")

	axePos := Vec2D{7, 7}
	glovePos := Vec2D{-7, 7}
	treePos := Vec2D{0, 19}

	atModal := func(name string, pos Vec2D) GOAPModalVal {
		return GOAPModalVal{
			name: fmt.Sprintf("at%s", name),
			check: func(ws *GOAPWorldState) int {
				ourPos := ws.GetModal(e, POSITION).(*Vec2D)
				_, _, d := ourPos.Distance(pos)
				if d < 2 {
					return 1
				} else {
					return 0
				}
			},
			effModalSet: func(ws *GOAPWorldState, op string, x int) {
				near := pos.Add(Vec2D{1, 0})
				ws.SetModal(e, POSITION, &near)
			},
		}
	}

	atAxeModal := atModal("Axe", axePos)
	atGloveModal := atModal("Glove", glovePos)
	atTreeModal := atModal("Tree", treePos)

	goTo := func(name string) *GOAPAction {
		return NewGOAPAction(map[string]any{
			"name": fmt.Sprintf("goTo%s", name),
			"cost": 1,
			"pres": nil,
			"effs": map[string]int{
				fmt.Sprintf("at%s,=", name): 1,
			},
		})
	}

	goToAxe := goTo("Axe")
	goToGlove := goTo("Glove")
	goToTree := goTo("Tree")

	chopTree := NewGOAPAction(map[string]any{
		"name": "chopTree",
		"cost": 1,
		"pres": []any{
			map[string]int{
				"hasGlove,=": 1,
				"hasAxe,=":   1,
			},
			map[string]int{
				"atTree,=": 1,
			},
		},
		"effs": map[string]int{
			"woodChopped,+": 1,
		},
	})

	p := NewGOAPPlanner(e)

	p.AddModalVals(hasGloveModal, hasAxeModal, atAxeModal, atGloveModal, atTreeModal)
	p.AddActions(getAxe, getGlove, goToAxe, goToGlove, goToTree, chopTree)

	ws := NewGOAPWorldState(nil)

	goal := map[string]int{
		"woodChopped,=": 3,
	}
	t0 := time.Now()
	plan, ok := p.Plan(ws, goal, 500)
	if !ok {
		t.Fatal("Should've found a solution")
	}
	Logger.Println(color.InGreenOverWhite(GOAPPathToString(plan)))
	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	Logger.Printf("Took %f ms to find solution", dt_ms)
}

func TestGOAPPlanResponsibleFridgeUsage(t *testing.T) {
	w := testingWorld()

	e := w.Spawn(nil)

	openFridge := NewGOAPAction(map[string]any{
		"name": "openFridge",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"fridgeOpen,=": 1,
		},
	})
	closeFridge := NewGOAPAction(map[string]any{
		"name": "closeFridge",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"fridgeOpen,=": 0,
		},
	})
	getFoodFromFridge := NewGOAPAction(map[string]any{
		"name": "getFoodFromFridge",
		"cost": 1,
		"pres": map[string]int{
			"fridgeOpen,=": 1,
		},
		"effs": map[string]int{
			"food,+": 1,
		},
	})

	p := NewGOAPPlanner(e)

	p.AddActions(openFridge, getFoodFromFridge, closeFridge)

	ws := NewGOAPWorldState(map[string]int{
		"fridgeOpen": 0,
	})

	goal := map[string]int{
		"fridgeOpen,=": 0,
		"food,=":       1,
	}
	t0 := time.Now()
	plan, ok := p.Plan(ws, goal, 500)
	if !ok {
		t.Fatal("Should've found a solution")
	}
	Logger.Println(color.InGreenOverWhite(GOAPPathToString(plan)))
	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	Logger.Printf("Took %f ms to find solution", dt_ms)

}

func TestGOAPPlanFarmer2000(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(ps, items, inventories)

	items.CreateArchetype(map[string]any{
		"name":        "yoke",
		"displayName": "a yoke for cattle",
		"flavourText": "one of mankind's greatest inventions... an ancestral gift!",
		"properties": map[string]int{
			"value": 25,
		},
		"tags": []string{"item.agricultural"},
		"entity": map[string]any{
			"sprite": "yoke",
			"box":    [2]float64{0.2, 0.2},
		},
	})

	// NOTE: all spawns are on x = 0
	// villager
	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION:  Vec2D{0, 0},
			BOX:       Vec2D{1, 1},
			INVENTORY: inventories.Create(nil),
		},
	})
	// yoke
	yoke := items.CreateItemSimple("yoke")
	items.SpawnItemEntity(Vec2D{0, 5}, yoke)
	// oxen
	spawnOxen := func(positions []Vec2D) {
		for i := 0; i < len(positions); i++ {
			w.Spawn(map[string]any{
				"components": map[ComponentID]any{
					POSITION: positions[i],
					BOX:      Vec2D{3, 2},
				},
				"tags": []string{"ox"},
			})
		}
	}
	spawnOxen([]Vec2D{})
	// field
	field := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION: Vec2D{0, 100},
			BOX:      Vec2D{100, 100},
		},
	})

	oxInFieldModal := GOAPModalVal{
		name:  "oxInField",
		nodes: []string{"ox", "field"},
		check: func(ws *GOAPWorldState) int {
			ox := ws.ModalEntities["ox"]
			oxPos := ws.GetModal(ox, POSITION).(*Vec2D)
			if RectIntersectsRect(
				*oxPos, *ox.GetVec2D(BOX),
				*field.GetVec2D(POSITION), *field.GetVec2D(BOX)) {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			ox := ws.ModalEntities["ox"]
			field := ws.ModalEntities["field"]
			if op == "=" {
				switch x {
				case 0:
					// TODO: this should really be a call to some kind of sophisticated
					// relocation function that avoids obstacles and makes sure there's a path
					// to be able to get there via navmesh/grid
					awayFromField := field.GetVec2D(POSITION).Add(Vec2D{200, 200})
					ws.SetModal(ox, POSITION, &awayFromField)
				case 1:
					fieldCenter := *field.GetVec2D(POSITION)
					ws.SetModal(ox, POSITION, &fieldCenter)
				}
			}
		},
	}
	hasYokeModal := GOAPModalVal{
		name: "hasYoke",
		check: func(ws *GOAPWorldState) int {
			inv := ws.GetModal(e, INVENTORY).(*Inventory)
			count := inv.CountName("yoke")
			return count
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			inv := ws.GetModal(e, INVENTORY).(*Inventory).CopyOf()
			if op == "-" {
				inv.DebitNName(x, "yoke")
			}
			if op == "=" {
				if x == 0 {
					inv.DebitAllName("yoke")
				}
				count := inv.CountName("yoke")
				if count == 0 {
					inv.Credit(items.CreateStackSimple(1, "yoke"))
				}
				inv.SetCountName(x, "yoke")
			}
			if op == "+" {
				count := inv.CountName("yoke")
				if count == 0 {
					inv.Credit(items.CreateStackSimple(x, "yoke"))
				} else {
					inv.SetCountName(count+x, "yoke")
				}
			}
			ws.SetModal(e, INVENTORY, inv)
		},
	}

	leadOxToField := NewGOAPAction(map[string]any{
		"name":           "leadOxToField",
		"node":           "ox",
		"travelWithNode": true,
		"cost":           1,
		"pres":           nil,
		"effs": map[string]int{
			"oxInField,=": 1,
		},
	})
	getYoke := NewGOAPAction(map[string]any{
		"name": "getYoke",
		"node": "yoke",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"hasYoke,=": 1,
		},
	})
	yokeOxplow := NewGOAPAction(map[string]any{
		"name":       "yokeOxplow",
		"node":       "ox",
		"otherNodes": []string{"field"},
		"cost":       1,
		"pres": []any{
			map[string]int{
				"hasYoke,=": 1,
			},
		},
		"effs": map[string]int{
			"oxYoked,=": 1,
		},
	})
	oxplow := NewGOAPAction(map[string]any{
		"name": "oxplow",
		"node": "ox",
		"cost": 1,
		"pres": []any{
			map[string]int{
				"oxInField,=": 1,
			},
			map[string]int{
				"oxYoked,=": 1,
			},
		},
		"effs": map[string]int{
			"fieldTilled,=": 1,
		},
	})

	p := NewGOAPPlanner(e)

	p.AddModalVals(oxInFieldModal, hasYokeModal)
	p.AddActions(leadOxToField, getYoke, yokeOxplow, oxplow)

	mockMakeTillPlan := func() {
		e.SetMind("plan.field", field)
		planField := e.GetMind("plan.field").(*Entity)

		// this would really be a filtering not of all entities but of perception
		closestOxToField := e.World.ClosestEntityFilter(
			*planField.GetVec2D(POSITION),
			*planField.GetVec2D(BOX),
			func(e *Entity) bool {
				return e.GetTagList(GENERICTAGS).Has("ox")
			})
		e.SetMind("plan.ox", closestOxToField)
	}

	p.BindEntitySelectors(map[string]func(*Entity) bool{
		// ox from blackboard plan - the closest to the field
		"ox": func(candidate *Entity) bool {
			return candidate == e.GetMind("plan.ox")
		},
		// any yoke
		"yoke": func(candidate *Entity) bool {
			if candidate.HasComponent(INVENTORY) {
				inv := candidate.GetGeneric(INVENTORY).(*Inventory)
				return inv.CountName("yoke") > 0
			}
			return false
		},
		// the field from the blackboard plan
		"field": func(candidate *Entity) bool {
			return candidate == e.GetMind("plan.field")
		},
	})

	ws := NewGOAPWorldState(map[string]int{
		"fieldTilled": 0,
	})

	goal := map[string]int{
		"fieldTilled,=": 1,
	}

	// first run with no oxen
	t0 := time.Now()
	mockMakeTillPlan() // include the blackboard context setting in the time
	_, ok := p.Plan(ws, goal, 500)
	if ok {
		t.Fatal("No oxen spawned, how did you find a solution?")
	}
	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	Logger.Printf("Took %f ms to fail", dt_ms)

	// spawn them
	spawnOxen([]Vec2D{Vec2D{0, 100}, Vec2D{0, 20}, Vec2D{0, -100}})

	// second run with oxen
	t0 = time.Now()
	mockMakeTillPlan() // include the blackboard context setting in the time
	plan, ok := p.Plan(ws, goal, 500)
	if !ok {
		t.Fatal("Should've found a solution")
	}
	Logger.Println(color.InGreenOverWhite(GOAPPathToString(plan)))
	dt_ms = float64(time.Since(t0).Nanoseconds()) / 1.0e6
	Logger.Printf("Took %f ms to find solution", dt_ms)

}
