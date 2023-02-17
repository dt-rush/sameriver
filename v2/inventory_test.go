package sameriver

import "testing"

func testingSpawnInventory(em EntityManagerInterface) (*Entity, error) {
	return em.Spawn([]string{}, MakeComponentSet(map[string]interface{}{
		"Generic,Inventory": NewInventory(),
	}))
}

func TestInventoryDebitCredit(t *testing.T) {
	w := testingWorld()
	i := NewInventorySystem()
	w.RegisterSystems(i)
	e, _ := testingSpawnInventory(w)

	inv := e.GetGeneric("Inventory").(Inventory)

	sword := NewItem(
		"sword_iron",
		"iron sword",
		"a good iron sword, decently sharp",
		map[string]int{
			"damage": 3,
			"value":  20,
		},
		[]string{"weapon"},
	)

	bottleOfBooze := NewItem(
		"bottle_booze",
		"bottle of booze",
		"a fiery brew, potent!",
		map[string]int{
			"drunkness": 2,
			"value":     5,
		},
		[]string{"consumable"},
	)

	inv.Credit(sword)
	inv.CreditN(bottleOfBooze, 5)

	Logger.Printf("Inventory: %s", inv.String())

	if len(inv.Items) != 2 {
		t.Fatal("Credit should have added items!")
	}

	swordInInv := inv.NameFilter("sword_iron")[0]
	retrieved := inv.Debit(swordInInv)

	if retrieved.Name != "sword_iron" {
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
}
