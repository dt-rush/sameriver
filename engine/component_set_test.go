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
	cs := fullZeroedComponentSet()
	cs.ToBitArray()
}

func TestComponentSetApply(t *testing.T) {
	w := testingWorld()
	e, _ := testingSpawnSimple(w)
	cs := fullZeroedComponentSet()
	w.ApplyComponentSet(cs)(e)
	if !e.ComponentBitArray.Equals(cs.ToBitArray()) {
		t.Fatal("failed to modify bitarray")
	}
}
