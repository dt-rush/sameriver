package sameriver

import (
	"fmt"

	"testing"
)

/*
benchmark output comment data in reference to machine:

goos: linux
goarch: amd64
pkg: github.com/dt-rush/sameriver/v3
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
*/

/*
BenchmarkGOAPClassic
BenchmarkGOAPClassic-8   	    2443	    549238 ns/op	  211534 B/op	    4367 allocs/op
BenchmarkGOAPClassic-8   	    2934	    391883 ns/op	  211543 B/op	    4367 allocs/op
BenchmarkGOAPClassic-8   	    3348	    407299 ns/op	  211537 B/op	    4367 allocs/op
BenchmarkGOAPClassic-8   	    2844	    426119 ns/op	  211546 B/op	    4367 allocs/op
BenchmarkGOAPClassic-8   	    2755	    404956 ns/op	  211545 B/op	    4367 allocs/op
BenchmarkGOAPClassic-8   	    3117	    397921 ns/op	  211538 B/op	    4367 allocs/op
BenchmarkGOAPClassic-8   	    2510	    405034 ns/op	  211543 B/op	    4367 allocs/op
PASS
ok  	github.com/dt-rush/sameriver/v3	10.639s

about 2380 per second or 238 if planning gets 10% of runtime per frame
*/
func BenchmarkGOAPClassic(b *testing.B) {
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

	ws := NewGOAPWorldState(map[string]int{
		"woodChopped": 0,
	})

	goal := map[string]int{
		"woodChopped,=": 3,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Plan(ws, goal, 500)
	}

}

/*
BenchmarkGOAPAlanWatts
BenchmarkGOAPAlanWatts-8   	    6974	    149900 ns/op	   81253 B/op	    1763 allocs/op
BenchmarkGOAPAlanWatts-8   	    7536	    136743 ns/op	   81252 B/op	    1763 allocs/op
BenchmarkGOAPAlanWatts-8   	    8928	    131337 ns/op	   81250 B/op	    1763 allocs/op
BenchmarkGOAPAlanWatts-8   	    9349	    137459 ns/op	   81251 B/op	    1763 allocs/op
BenchmarkGOAPAlanWatts-8   	    7573	    143881 ns/op	   81250 B/op	    1763 allocs/op
BenchmarkGOAPAlanWatts-8   	    8674	    146659 ns/op	   81252 B/op	    1763 allocs/op
BenchmarkGOAPAlanWatts-8   	    8470	    136132 ns/op	   81252 B/op	    1763 allocs/op
PASS
ok  	github.com/dt-rush/sameriver/v3	9.050s

avg about 7000 per second or 700 if we give planning 10% of time per frame
*/
func BenchmarkGOAPAlanWatts(b *testing.B) {
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
			STATE: map[string]int{
				"drunk": 0,
			},
			INVENTORY: inventories.Create(nil),
			POSITION:  Vec2D{10, 10},
			VELOCITY:  Vec2D{0, 0},
			BOX:       Vec2D{1, 1},
			MASS:      3.0,
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Plan(ws, goal, 500)
	}
}

/*
BenchmarkGOAPSimple
BenchmarkGOAPSimple-8   	   34522	     39616 ns/op	   22424 B/op	     490 allocs/op
BenchmarkGOAPSimple-8   	   33502	     35879 ns/op	   22424 B/op	     489 allocs/op
BenchmarkGOAPSimple-8   	   37714	     31034 ns/op	   22424 B/op	     489 allocs/op
BenchmarkGOAPSimple-8   	   37502	     32084 ns/op	   22424 B/op	     490 allocs/op
BenchmarkGOAPSimple-8   	   39062	     32116 ns/op	   22424 B/op	     489 allocs/op
BenchmarkGOAPSimple-8   	   36994	     32680 ns/op	   22424 B/op	     490 allocs/op
BenchmarkGOAPSimple-8   	   37473	     32743 ns/op	   22424 B/op	     490 allocs/op
PASS
ok  	github.com/dt-rush/sameriver/v3	10.998s

about 28000 per second, or 2800 if planning gets 10% of runtime per frame
*/
func BenchmarkGOAPSimple(b *testing.B) {
	w := testingWorld()

	ps := NewPhysicsSystem()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(ps, items, inventories)

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION:  Vec2D{0, 0},
			INVENTORY: inventories.Create(nil),
		},
	})

	equipBow := NewGOAPAction(map[string]any{
		"name": "equipBow",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"bowEquipped,=": 1,
		},
	})

	moveToTarget := NewGOAPAction(map[string]any{
		"name": "moveToTarget",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"inRangeOfTarget,=": 1,
		},
	})

	rangedCombat := NewGOAPAction(map[string]any{
		"name": "rangedCombat",
		"cost": 1,
		"pres": map[string]int{
			"bowEquipped,=":     1,
			"inRangeOfTarget,=": 1,
		},
		"effs": map[string]int{
			"combat,=": 1,
		},
	})

	p := NewGOAPPlanner(e)

	p.AddActions(equipBow, moveToTarget, rangedCombat)

	ws := NewGOAPWorldState(nil)

	goal := map[string]int{
		"combat,=": 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Plan(ws, goal, 500)
	}
}

/*
	first result of the new algorithm/pattern

BenchmarkGOAPFarmer2000
BenchmarkGOAPFarmer2000-8   	    6187	    183668 ns/op
BenchmarkGOAPFarmer2000-8   	    6724	    175238 ns/op
BenchmarkGOAPFarmer2000-8   	    7591	    150992 ns/op
BenchmarkGOAPFarmer2000-8   	    7333	    148206 ns/op
BenchmarkGOAPFarmer2000-8   	    7449	    150329 ns/op      <---- that's about 1108 entities per second doing a farmer2000 plan,
BenchmarkGOAPFarmer2000-8   	    7596	    148335 ns/op            assuming planning gets 1/6th (0.167) of the time per frame
BenchmarkGOAPFarmer2000-8   	    7628	    150825 ns/op
BenchmarkGOAPFarmer2000-8   	    7988	    149477 ns/op
BenchmarkGOAPFarmer2000-8   	    7665	    153718 ns/op
BenchmarkGOAPFarmer2000-8   	    7519	    153100 ns/op
ok  	github.com/dt-rush/sameriver/v4	13.348s
*/
func BenchmarkGOAPFarmer2000(b *testing.B) {
	//
	// world init
	//
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

	//
	// item system init
	//
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

	//
	// spawn entities
	//

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
			STATE: map[string]int{
				"tilled": 0,
			},
		},
		"tags": []string{"field"},
	})

	//
	// modal world state vars
	//
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
	fieldTilledModal := GOAPModalVal{
		name:  "fieldTilled",
		nodes: []string{"field"},
		check: func(ws *GOAPWorldState) int {
			field := ws.ModalEntities["field"]
			return ws.GetModal(field, STATE).(*IntMap).m["tilled"]
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			if op == "=" {
				field := ws.ModalEntities["field"]
				state := field.GetIntMap(STATE).CopyOf()
				state.m["tilled"] = x
				ws.SetModal(field, STATE, &state)
			}
		},
	}

	// GOAP actions
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

	//
	// PLANNER INIT
	//

	// in this test, only the yoke selector gets used as a generic fallback since
	// we don't bind any more specific selector for "yoke" before any Plan() call.
	// *all* of these would be used for more general planning than this specific constrained
	// example where we have a field in mind allowing us to choose a closer ox.
	// so really, this would happen not before the planning as the BindEntitySelectors() call below,
	// but this RegisterGenericEntitySelectors() call would happen on setup of the planner itself
	p := NewGOAPPlanner(e)
	p.RegisterGenericEntitySelectors(map[string]func(*Entity) bool{
		// any ox
		"ox": func(candidate *Entity) bool {
			return candidate.GetTagList(GENERICTAGS).Has("ox")
		},
		// any yoke
		"yoke": func(candidate *Entity) bool {
			if candidate.HasComponent(INVENTORY) {
				inv := candidate.GetGeneric(INVENTORY).(*Inventory)
				return inv.CountName("yoke") > 0
			}
			return false
		},
		// any field
		"field": func(candidate *Entity) bool {
			return candidate.GetTagList(GENERICTAGS).Has("field")
		},
	})
	p.AddModalVals(oxInFieldModal, hasYokeModal, fieldTilledModal)
	p.AddActions(leadOxToField, getYoke, yokeOxplow, oxplow)

	//
	// bb workplan
	//

	// NOTE: we'd *get* the currently active bb work plan for the field rather than
	// generate it if someone was already doing plant
	tillPlanBB := func() {
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
	tillPlanBindEntities := func() {
		p.BindEntitySelectors(map[string]any{
			// ox from blackboard plan - the closest to the field
			"ox": "mind.plan.ox",
			// the field from the blackboard plan
			"field": "mind.plan.field",
		})
	}
	mockMakeTillPlan := func() {
		tillPlanBB()
		tillPlanBindEntities()
	}

	//
	// initial world state
	//
	// TODO: this would be a perception system thing - we would get the current
	// state of the world from perception/memory
	ws := NewGOAPWorldState(nil)

	// TODO: this would derive from a utility, not be hardcoded
	goal := map[string]int{
		"fieldTilled,=": 1,
	}

	// spawn them
	// spawnOxen([]Vec2D{Vec2D{0, 100}, Vec2D{0, 20}, Vec2D{0, -100}}) // one in field already
	spawnOxen([]Vec2D{Vec2D{0, 20}, Vec2D{0, -100}})
	b.ResetTimer()
	mockMakeTillPlan()
	for i := 0; i < b.N; i++ {
		p.Plan(ws, goal, 500)
	}
}
