package sameriver

import (
	"testing"
)

func TestInvalidComponentType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should panic if given component type Vec7D")
		}
	}()
	MakeComponentSet(map[string]interface{}{
		"Vec7D,Position": Vec2D{0, 0},
	})
}

func TestComponentSetToBitArray(t *testing.T) {
	w := testingWorld()
	l := NewTagList()
	cs := MakeComponentSet(map[string]interface{}{
		"TagList,GenericTag": &l,
	})
	w.em.components.BitArrayFromComponentSet(cs)
}

func TestComponentSetApply(t *testing.T) {
	w := testingWorld()
	e, _ := testingSpawnSimple(w)
	l := NewTagList()
	cs := MakeComponentSet(map[string]interface{}{
		"TagList,GenericTag": &l,
	})
	w.em.components.ApplyComponentSet(e, cs)
	if !e.ComponentBitArray.Equals(w.em.components.BitArrayFromComponentSet(cs)) {
		t.Fatal("failed to apply componentset according to bitarray")
	}
}
