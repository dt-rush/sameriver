package sameriver

import "testing"

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

func TestItemSystemLoadArchetypes(t *testing.T) {
	i := NewItemSystem()
	i.LoadArchetypesFile("test_data/basic_archetypes.json")
	Logger.Println(i.Archetypes)
	if len(i.Archetypes) != 3 {
		t.Fatal("Did not load from JSON file!")
	}
}
