package sameriver

import "testing"

func testingSpawnInventory(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[string]any{"Generic,Inventory": NewInventory()},
	})
}

func TestInventoryDebitCredit(t *testing.T) {
	w := testingWorld()
	i := NewItemSystem()
	w.RegisterSystems(i)
	e := testingSpawnInventory(w)

	inv := e.GetGeneric("Inventory").(Inventory)

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
	inv.CreditN(bottleOfBooze, 5)

	Logger.Printf("Inventory: %s", inv.String())

	if len(inv.Items) != 2 {
		t.Fatal("Credit should have added items!")
	}

	swordInInv := inv.NameFilter("sword_iron")[0]
	retrieved := inv.Debit(swordInInv)

	if retrieved.Archetype.Name != "sword_iron" {
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

	inv.DebitNName("bottle_booze", 3)

	if len(inv.Items) != 0 {
		t.Fatal("Should've taken every last bottle of booze!")
	}

	inv.CreditN(bottleOfBooze, 10)
	inv.DebitAllName("bottle_booze")

	if len(inv.Items) != 0 {
		t.Fatal("Should've taken every last bottle of booze!")
	}

	inv.CreditN(sword, 10)
	s := inv.DebitName("sword_iron")
	s.SetProperty("degradation", 5)
	inv.Credit(s)
	Logger.Println(inv)

	if len(inv.Items) != 2 {
		t.Fatal("Should've put in a modified sword as a separate item from the stack")
	}
}

func TestItemFromArchetype(t *testing.T) {
	w := testingWorld()
	i := NewItemSystem()
	w.RegisterSystems(i)

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

	i.CreateSubArchetype(map[string]any{
		"parent":      "sword_iron",
		"name":        "sword_iron_manjushris",
		"displayName": "manjushri's sword",
		"flavourText": "the sword of the legendary bodhisattva Manjushri; it can cut illusion itself",
		"tagDiff":     []string{"+legendary"},
	})

	manjushrisSword := i.CreateItem(map[string]any{
		"archetype": "sword_iron_manjushris",
	})

	Logger.Printf("Created: %s", manjushrisSword.String())

	if !manjushrisSword.Tags.Has("weapon") {
		t.Fatal("did not inherit tags!")
	}

	if manjushrisSword.GetProperty("damage") != 3 {
		t.Fatal("did not inherit properties!")
	}

	manjushrisSword.SetProperty("damage", 108)

	if manjushrisSword.GetProperty("damage") != 108 {
		t.Fatal("Did not set property!")
	}
}
