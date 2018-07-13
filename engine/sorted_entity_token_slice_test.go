package engine

import (
	"testing"
)

func TestSortedEntityTokenSliceSearch(t *testing.T) {
	e1 := &EntityToken{ID: 1}
	e2 := &EntityToken{ID: 2}
	e3 := &EntityToken{ID: 3}
	e4 := &EntityToken{ID: 4}
	e5 := &EntityToken{ID: 5}
	e6 := &EntityToken{ID: 6}
	e7 := &EntityToken{ID: 7}
	slice := []*EntityToken{e1, e2, e3, e4, e5, e6, e7}
	if SortedEntityTokenSliceSearch(slice, e1) != 0 {
		t.Fatal("failed to find element e1")
	}
	if SortedEntityTokenSliceSearch(slice, e4) != 3 {
		t.Fatal("failed to find element e4")
	}
	if SortedEntityTokenSliceSearch(slice, e7) != 6 {
		t.Fatal("failed to find element e7")
	}
}

func TestSortedEntityTokenSliceInsertIfNotPresent(t *testing.T) {
	e1 := &EntityToken{ID: 1}
	e2 := &EntityToken{ID: 2}
	e3 := &EntityToken{ID: 3}
	slice := []*EntityToken{e1, e2, e3}
	SortedEntityTokenSliceInsertIfNotPresent(&slice, e2)
	if len(slice) != 3 {
		t.Fatal("inserted when already present")
	}
}

func TestSortedEntityTokenSliceRemove(t *testing.T) {
	e1 := &EntityToken{ID: 1}
	e2 := &EntityToken{ID: 2}
	e3 := &EntityToken{ID: 3}
	e4 := &EntityToken{ID: 4}
	slice := []*EntityToken{e1, e2, e3, e4}
	SortedEntityTokenSliceRemove(&slice, e2)
	if !(slice[0].ID == 1 &&
		slice[1].ID == 3 &&
		slice[2].ID == 4) {
		t.Fatal("did not remove the right element")
	}
}
