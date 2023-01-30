package engine

import (
	"reflect"
	"testing"
)

func TestValidateComponentSetFields(t *testing.T) {
	ty := reflect.TypeOf(ComponentSet{})
	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)
		if f.Type.Kind() != reflect.Ptr {
			t.Errorf("field %s was not of pointer type", f.Name)
		}
	}
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
