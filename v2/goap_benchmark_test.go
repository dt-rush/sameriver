package sameriver

import (
	"fmt"

	"testing"
)

func BenchmarkGOAPClassic(b *testing.B) {
	w := testingWorld()

	ps := NewPhysicsSystem()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(ps, items, inventories)

	w.RegisterComponents("IntMap,State", "Generic,Inventory")

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
		"components": map[string]any{
			"Vec2D,Position":    Vec2D{0, 0},
			"Generic,Inventory": inventories.Create(nil),
		},
	})

	hasModal := func(name string, archetype string, tags ...string) GOAPModalVal {
		return GOAPModalVal{
			name: fmt.Sprintf("has%s", name),
			check: func(ws *GOAPWorldState) int {
				inv := ws.GetModal(e, "Inventory").(*Inventory)
				return inv.CountName(archetype)
			},
			effModalSet: func(ws *GOAPWorldState, op string, x int) {
				inv := ws.GetModal(e, "Inventory").(*Inventory).CopyOf()
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
				ws.SetModal(e, "Inventory", inv)
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
				ourPos := ws.GetModal(e, "Position").(*Vec2D)
				_, _, d := ourPos.Distance(pos)
				if d < 2 {
					return 1
				} else {
					return 0
				}
			},
			effModalSet: func(ws *GOAPWorldState, op string, x int) {
				near := pos.Add(Vec2D{1, 0})
				ws.SetModal(e, "Position", &near)
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
		"pres": map[string]int{
			"hasGlove,=": 1,
			"hasAxe,=":   1,
			"atTree,=":   1,
		},
		"effs": map[string]int{
			"woodChopped,+": 1,
		},
	})

	p := NewGOAPPlanner(e)

	p.eval.AddModalVals(hasGloveModal, hasAxeModal, atAxeModal, atGloveModal, atTreeModal)
	p.eval.AddActions(getAxe, getGlove, goToAxe, goToGlove, goToTree, chopTree)

	ws := NewGOAPWorldState(map[string]int{
		"woodChopped": 0,
	})

	goal := NewGOAPGoal(map[string]int{
		"woodChopped,=": 3,
	})
	for i := 0; i < b.N; i++ {
		p.Plan(ws, goal, 500)
	}

}

func BenchmarkGOAPAlanWatts(b *testing.B) {
	w := testingWorld()

	ps := NewPhysicsSystem()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(ps, items, inventories)

	w.RegisterComponents("IntMap,State", "Generic,Inventory")

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
		"components": map[string]any{
			"IntMap,State": map[string]int{
				"drunk": 0,
			},
			"Generic,Inventory": inventories.Create(nil),
		},
	})

	boozePos := &Vec2D{19, 19}
	templePos := &Vec2D{-19, 19}

	inTempleModal := GOAPModalVal{
		name: "inTemple",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*templePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearTemple := templePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearTemple)
		},
	}
	atBoozeModal := GOAPModalVal{
		name: "atBooze",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*boozePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearBooze := boozePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearBooze)
		},
	}
	drunkModal := GOAPModalVal{
		name: "drunk",
		check: func(ws *GOAPWorldState) int {
			state := ws.GetModal(e, "State").(*IntMap)
			return state.m["drunk"]
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			state := ws.GetModal(e, "State").(*IntMap).CopyOf()
			if op == "+" {
				state.m["drunk"] += x
			}
			ws.SetModal(e, "State", &state)
		},
	}
	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			inv := ws.GetModal(e, "Inventory").(*Inventory)
			count := inv.CountTags("booze")
			return count
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			inv := ws.GetModal(e, "Inventory").(*Inventory).CopyOf()
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
			ws.SetModal(e, "Inventory", inv)
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

	p.eval.AddModalVals(drunkModal, hasBoozeModal, atBoozeModal, inTempleModal)
	p.eval.AddActions(drink, dropAllBooze, purifyOneself, enterTemple, goToBooze, getBooze)

	ws := NewGOAPWorldState(map[string]int{
		"rituallyPure": 0,
	})

	goal := NewGOAPGoal(map[string]int{
		"drunk,>=":   3,
		"inTemple,=": 1,
	})

	for i := 0; i < b.N; i++ {
		p.Plan(ws, goal, 500)
	}
}