package sameriver

import (
	"testing"
)

func TestInventoryDebitCredit(t *testing.T) {
	w := testingWorld()
	inventories := NewInventorySystem()
	items := NewItemSystem(nil)
	w.RegisterSystems(items, inventories)
	w.RegisterComponents([]string{"Generic,Inventory"})
	e := w.Spawn(map[string]any{
		"components": map[string]any{"Generic,Inventory": NewInventory()},
	})

	inv := e.GetGeneric("Inventory").(*Inventory)

	items.CreateArchetype(map[string]any{
		"name":        "sword_iron",
		"displayName": "iron sword",
		"flavourText": "a good irons word, decently sharp",
		"properties": map[string]int{
			"damage":      3,
			"value":       20,
			"degradation": 0,
			"durability":  5,
		},
		"tags": []string{"weapon"},
	})

	items.CreateArchetype(map[string]any{
		"name":        "bottle_booze",
		"displayName": "bottle of booze",
		"flavourText": "a fiery brew, potent!",
		"properties": map[string]int{
			"drunkness": 2,
			"value":     5,
		},
		"tags": []string{"consumable"},
	})

	sword := items.CreateItem(map[string]any{
		"archetype": "sword_iron",
	})
	Logger.Println(sword.Degradations)
	boozeStack := items.CreateStack(5, map[string]any{
		"archetype": "bottle_booze",
	})

	inv.Credit(sword)
	Logger.Println(sword.Degradations)
	inv.Credit(boozeStack)

	Logger.Printf("Inventory: %s", inv.String())

	if len(inv.Stacks) != 2 {
		t.Fatal("Credit should have added items!")
	}

	swordInInv := inv.NameFilter("sword_iron")[0]
	Logger.Println(swordInInv.Degradations)
	Logger.Printf("swordInInv.Count = %d", swordInInv.Count)
	retrieved := inv.Debit(swordInInv)

	if retrieved.GetArchetype().Name != "sword_iron" {
		t.Fatal("Did not retrieve debited item!")
	}

	if len(inv.Stacks) != 1 {
		t.Fatal("Should have one item left")
	}

	boozeInInv := inv.NameFilter("bottle_booze")[0]
	retrieved = inv.DebitN(boozeInInv, 2)
	Logger.Println(retrieved)

	if retrieved.Count != 2 {
		t.Fatal("Did not retrieve the right number of booze!")
	}

	if boozeInInv.Count != 3 {
		Logger.Printf("%d remaining", boozeInInv.Count)
		t.Fatal("Did not subtract when debiting!")
	}

	inv.DebitN(boozeInInv, 3)

	if len(inv.Stacks) != 0 {
		Logger.Println(inv)
		t.Fatal("Should've taken every last bottle of booze!")
	}

	boozeStack = items.CreateStackSimple(10, "bottle_booze")
	inv.Credit(boozeStack)
	inv.DebitN(boozeStack, 10)

	if len(inv.Stacks) != 0 {
		t.Fatal("Should've taken every last bottle of booze!")
	}

	swordStack := items.CreateStackSimple(10, "sword_iron")
	inv.Credit(swordStack)
	s := inv.Debit(swordStack)
	s.SetProperty("damage", 5)
	inv.Credit(s)

	if len(inv.Stacks) != 2 {
		t.Fatal("Should've put in a modified sword as a separate item from the stack")
	}

	newSword := items.CreateItemSimple("sword_iron")
	inv.Credit(newSword)

	if len(inv.Stacks) != 2 || swordStack.Count != 10 {
		t.Fatal("Stacks should absorb if all props match")
	}
}

func TestInventoryFromListing(t *testing.T) {
	w := testingWorld()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(items, inventories)

	items.CreateArchetype(map[string]any{
		"name":        "sword_iron",
		"displayName": "iron sword",
		"flavourText": "a good irons word, decently sharp",
		"properties": map[string]int{
			"damage":     3,
			"value":      50,
			"durability": 5,
		},
		"tags": []string{"weapon", "common", "degrades"},
	})

	items.CreateArchetype(map[string]any{
		"name":        "coin_copper",
		"displayName": "copper coin",
		"flavourText": "copper die-cast coin with an elephant on it",
		"properties": map[string]int{
			"value": 2,
		},
		"tags": []string{"currency"},
	})

	items.CreateArchetype(map[string]any{
		"name":        "heart_sutra",
		"displayName": "heart sutra",
		"flavourText": "small copy of the heart sutra",
		"properties": map[string]int{
			"value": 10,
		},
		"tags": []string{"book"},
	})

	e := w.Spawn(map[string]any{
		"components": map[string]any{
			"Generic,Inventory": inventories.Create(map[string]int{
				"sword_iron":  1,
				"coin_copper": 100,
				"heart_sutra": 1,
			}),
		},
	})

	inv := e.GetGeneric("Inventory").(*Inventory)

	coin := inv.NameFilter("coin_copper")[0]
	Logger.Println(coin)
	inv.DebitN(coin, coin.Count/2)
}

func TestInventoryStacksForDisplay(t *testing.T) {
	w := testingWorld()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(items, inventories)

	items.LoadArchetypesFile("test_data/basic_archetypes.json")

	e := w.Spawn(map[string]any{
		"components": map[string]any{
			"Generic,Inventory": inventories.Create(map[string]int{
				"sword_iron":  1,
				"coin_copper": 100,
				"heart_sutra": 1,
			}),
		},
	})

	eInv := e.GetGeneric("Inventory").(*Inventory)
	Logger.Println(eInv.StacksForDisplay())
}
