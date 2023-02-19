package sameriver

import (
	"testing"
)

func TestInventoryDebitCredit(t *testing.T) {
	w := testingWorld()
	i := NewItemSystem()
	w.RegisterSystems(i)
	w.RegisterComponents([]string{"Generic,Inventory"})
	e := w.Spawn(map[string]any{
		"components": map[string]any{"Generic,Inventory": NewInventory()},
	})

	inv := e.GetGeneric("Inventory").(*Inventory)

	i.CreateArchetype(map[string]any{
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

	i.CreateArchetype(map[string]any{
		"name":        "bottle_booze",
		"displayName": "bottle of booze",
		"flavourText": "a fiery brew, potent!",
		"properties": map[string]int{
			"drunkness": 2,
			"value":     5,
		},
		"tags": []string{"consumable"},
	})

	sword := i.CreateItem(map[string]any{
		"archetype": "sword_iron",
	})
	bottleOfBooze := i.CreateItem(map[string]any{
		"archetype": "bottle_booze",
	})

	inv.Credit(sword)
	boozeStack := bottleOfBooze.StackOf(5)
	inv.Credit(boozeStack)

	Logger.Printf("Inventory: %s", inv.String())

	if len(inv.Items) != 2 {
		t.Fatal("Credit should have added items!")
	}

	swordInInv := inv.NameFilter("sword_iron")[0]
	retrieved := inv.Debit(swordInInv)

	if retrieved.GetArchetype().Name != "sword_iron" {
		t.Fatal("Did not retrieve debited item!")
	}

	if len(inv.Items) != 1 {
		t.Fatal("Should have one item left")
	}

	boozeInInv := inv.NameFilter("bottle_booze")[0]
	retrieved = inv.DebitN(boozeInInv, 2)

	if retrieved.Count != 2 {
		t.Fatal("Did not retrieve the right number of booze!")
	}

	if boozeInInv.Count != 3 {
		Logger.Printf("%d remaining", boozeInInv.Count)
		t.Fatal("Did not subtract when debiting!")
	}

	inv.DebitN(boozeInInv, 3)

	if len(inv.Items) != 0 {
		Logger.Println(inv)
		t.Fatal("Should've taken every last bottle of booze!")
	}

	boozeStack = bottleOfBooze.StackOf(10)
	inv.Credit(boozeStack)
	inv.DebitN(boozeStack, 10)

	if len(inv.Items) != 0 {
		t.Fatal("Should've taken every last bottle of booze!")
	}

	swordStack := sword.StackOf(10)
	inv.Credit(swordStack)
	s := inv.Debit(swordStack)
	s.SetProperty("degradation", 5)
	inv.Credit(s)
	Logger.Println(inv)

	if len(inv.Items) != 2 {
		t.Fatal("Should've put in a modified sword as a separate item from the stack")
	}
}

func TestInventoryFromListing(t *testing.T) {
	w := testingWorld()
	items := NewItemSystem()
	inventories := NewInventorySystem()
	w.RegisterSystems(items, inventories)

	items.CreateArchetype(map[string]any{
		"name":        "sword_iron",
		"displayName": "iron sword",
		"flavourText": "a good irons word, decently sharp",
		"properties": map[string]int{
			"damage":      3,
			"value":       50,
			"degradation": 0,
			"durability":  5,
		},
		"tags": []string{"weapon", "common"},
	})

	items.CreateArchetype(map[string]any{
		"name":        "copper_coin",
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
				"copper_coin": 100,
				"heart_sutra": 1,
			}),
		},
	})

	inv := e.GetGeneric("Inventory").(*Inventory)

	coin := inv.NameFilter("copper_coin")[0]
	inv.DebitN(coin, coin.Count/2)
}
