package sameriver

import (
	"testing"
)

func TestInventoryDebitCredit(t *testing.T) {
	w := testingWorld()
	inventories := NewInventorySystem()
	items := NewItemSystem(nil)
	w.RegisterSystems(items, inventories)

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			INVENTORY: NewInventory(),
		}})

	inv := e.GetGeneric(INVENTORY).(*Inventory)

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

	swordInInv := inv.FilterName("sword_iron")[0]
	Logger.Println(swordInInv.Degradations)
	Logger.Printf("swordInInv.Count = %d", swordInInv.Count)
	retrieved := inv.Debit(swordInInv)

	if retrieved.GetArchetype().Name != "sword_iron" {
		t.Fatal("Did not retrieve debited item!")
	}

	if len(inv.Stacks) != 1 {
		t.Fatal("Should have one item left")
	}

	boozeInInv := inv.FilterName("bottle_booze")[0]
	retrieved = inv.DebitN(2, boozeInInv)
	Logger.Println(retrieved)

	if retrieved.Count != 2 {
		t.Fatal("Did not retrieve the right number of booze!")
	}

	if boozeInInv.Count != 3 {
		Logger.Printf("%d remaining", boozeInInv.Count)
		t.Fatal("Did not subtract when debiting!")
	}

	inv.DebitN(3, boozeInInv)

	if len(inv.Stacks) != 0 {
		Logger.Println(inv)
		t.Fatal("Should've taken every last bottle of booze!")
	}

	boozeStack = items.CreateStackSimple(10, "bottle_booze")
	inv.Credit(boozeStack)
	inv.DebitN(10, boozeStack)

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
		"components": map[ComponentID]any{
			INVENTORY: inventories.Create(map[string]int{
				"sword_iron":  1,
				"coin_copper": 100,
				"heart_sutra": 1,
			}),
		},
	})

	inv := e.GetGeneric(INVENTORY).(*Inventory)

	coin := inv.FilterName("coin_copper")[0]
	Logger.Println(coin)
	inv.DebitN(coin.Count/2, coin)
}

func TestInventoryStacksForDisplay(t *testing.T) {
	w := testingWorld()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(items, inventories)

	items.LoadArchetypesFile("test_data/basic_archetypes.json")

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			INVENTORY: inventories.Create(map[string]int{
				"sword_iron":  1,
				"coin_copper": 100,
				"heart_sutra": 1,
			}),
		},
	})

	eInv := e.GetGeneric(INVENTORY).(*Inventory)
	Logger.Println(eInv.StacksForDisplay())
	for _, str := range eInv.StacksForDisplay() {
		if str.DisplayStr == "copper coin x 100" {
			return
		}
	}
	t.Fatal("Should've found copper coin x 100")
}

func TestInventoryDebitNByFilter(t *testing.T) {
	w := testingWorld()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(items, inventories)

	items.LoadArchetypesFile("test_data/basic_archetypes.json")

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			INVENTORY: inventories.Create(map[string]int{
				"sword_iron":  1,
				"coin_copper": 100,
			}),
		},
	})

	eInv := e.GetGeneric(INVENTORY).(*Inventory)
	perishableSutraSpec := map[string]any{
		"archetype":       "heart_sutra",
		"tags":            []string{"perishable"},
		"degradationRate": 0.01,
	}
	book := items.CreateItem(perishableSutraSpec)
	eInv.Credit(book)
	debited := eInv.DebitNByFilter(2, func(s *Item) bool {
		return s.Tags.Has("degrades")
	})
	Logger.Println(debited)
	if len(debited) != 2 {
		t.Fatal("should've grabbed a sword and a heart_sutra")
	}
}
func TestInventoryGetCountContains(t *testing.T) {
	w := testingWorld()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(items, inventories)

	items.LoadArchetypesFile("test_data/basic_archetypes.json")

	i := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			INVENTORY: inventories.Create(map[string]int{
				"coin_copper": 100,
			}),
		},
	})
	j := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			INVENTORY: inventories.Create(map[string]int{
				"sword_iron": 1,
			}),
		},
	})
	iInv := i.GetGeneric(INVENTORY).(*Inventory)
	jInv := j.GetGeneric(INVENTORY).(*Inventory)
	// i purchases a sword
	jInv.GetNByName(iInv, 50, "coin_copper")
	iInv.GetNByName(jInv, 1, "sword_iron")

	if jInv.CountName("coin_copper") != 50 {
		t.Fatal("Should've got 50 coins - (or is Count() wrong?)")
	}

	// j takes the sword back and takes the rest of i's money (evil!)
	jInv.GetNByFilter(iInv, 1, func(it *Item) bool {
		return it.GetArchetype().Name == "sword_iron"
	})
	jInv.GetAllByName(iInv, "coin_copper")

	if iInv.CountName("coin_copper") != 0 {
		t.Fatal("Should've lost all coins to robbery!")
	}

	Logger.Println(iInv)
	Logger.Println(jInv)

	if !jInv.ContainsName("sword_iron") {
		t.Fatal("Should've taken the sword back (or is ContainsName() wrong?)")
	}

	if !jInv.Contains(func(it *Item) bool {
		return it.Tags.Has("currency")
	}) {
		t.Fatal("Should've had some currency")
	}
}
