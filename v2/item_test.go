package sameriver

import "testing"

func TestItemFromArchetype(t *testing.T) {
	ironSwordArchetype := NewItem(
		"sword_iron",
		"iron sword",
		"a good iron sword, decently sharp",
		map[string]int{
			"damage": 3,
			"value":  20,
		},
		[]string{"weapon"},
	)

	manjushrisSword := ItemFromArchetype(
		ironSwordArchetype,
		"manjushris",
		"manjushri's sword",
		"the sword of the legendary bodhisattva Manjushri; it can cut illusion",
	)
	manjushrisSword.AddTags("legendary")

	Logger.Printf("Created: %s", manjushrisSword.String())

	if !manjushrisSword.HasTag("weapon") {
		t.Fatal("did not inherit tags!")
	}

	if manjushrisSword.Properties["damage"] != 3 {
		t.Fatal("did not inherit properties!")
	}
}
