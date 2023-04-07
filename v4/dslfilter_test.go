package sameriver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDSLBasic(t *testing.T) {
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

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION:  Vec2D{0, 0},
			BOX:       Vec2D{1, 1},
			INVENTORY: inventories.Create(nil),
		},
	})

	yoke := items.CreateItemSimple("yoke")
	items.SpawnItemEntity(Vec2D{0, 5}, yoke)

	positions := []Vec2D{
		Vec2D{0, 10}, // close ox
		Vec2D{0, 30}, // far ox
	}
	tags := []string{
		"a",
		"b",
	}

	oxen := make([]*Entity, len(positions))
	for i := 0; i < len(positions); i++ {
		oxen[i] = w.Spawn(map[string]any{
			"components": map[ComponentID]any{
				POSITION: positions[i],
				BOX:      Vec2D{3, 2},
				STATE: map[string]int{
					"yoked": 0,
				},
			},
			"tags": []string{"ox", tags[i]},
		})
	}
	field := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION: Vec2D{0, 100},
			BOX:      Vec2D{30, 30},
			STATE: map[string]int{
				"tilled": 0,
			},
		},
		"tags": []string{"field"},
	})
	w.Blackboard("somebb").Set("field", field)

	//
	// entity
	//

	// Test Entity.DSLFilter
	Logger.Println("1")
	entities, err := e.DSLFilter("HasTag(ox)")
	assert.NoError(t, err)
	assert.ElementsMatch(t, oxen, entities)

	Logger.Println("2")
	entities, err = e.DSLFilter("HasComponent(position)")
	assert.NoError(t, err)
	assert.Equal(t, 5, len(entities)) // e, item.yoke, 2 oxen, field

	Logger.Println("3")
	entities, err = e.DSLFilter("WithinDistance(self, 15)")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(entities)) // e, item.yoke, close ox

	Logger.Println("4")
	oxen[0].GetIntMap(STATE).Set("yoked", 1)
	entities, err = e.DSLFilter("State(yoked, 1)")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entities))
	assert.Equal(t, oxen[0], entities[0])

	// Test Entity.DSLFilterSort
	Logger.Println("5")
	filterSortEntities, err := e.DSLFilterSort("HasTag(ox); Closest(self)")
	assert.NoError(t, err)
	assert.Equal(t, oxen[0], filterSortEntities[0]) // a close
	assert.Equal(t, oxen[1], filterSortEntities[1]) // b far

	//
	// world
	//

	// Test World.DSLFilter
	Logger.Println("6")
	worldEntities, err := w.DSLFilter("HasTag(ox)")
	assert.NoError(t, err)
	assert.ElementsMatch(t, oxen, worldEntities)
	// Test World.DSLFilterSort
	Logger.Println("7")
	worldFilterSortEntities, err := w.DSLFilterSort("HasTag(ox); Closest(bb.somebb.field)")
	assert.NoError(t, err)
	assert.Equal(t, oxen[1], worldFilterSortEntities[0]) // b close
	assert.Equal(t, oxen[0], worldFilterSortEntities[1]) // a far
}
