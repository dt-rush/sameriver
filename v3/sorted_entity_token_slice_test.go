package sameriver

import (
	"testing"
)

func TestSortedEntitySliceSearch(t *testing.T) {
	e1 := &Entity{ID: 1}
	e2 := &Entity{ID: 2}
	e3 := &Entity{ID: 3}
	e4 := &Entity{ID: 4}
	e5 := &Entity{ID: 5}
	e6 := &Entity{ID: 6}
	e7 := &Entity{ID: 7}
	slice := []*Entity{e1, e2, e3, e4, e5, e6, e7}
	if SortedEntitySliceSearch(slice, e1) != 0 {
		t.Fatal("failed to find element e1")
	}
	if SortedEntitySliceSearch(slice, e4) != 3 {
		t.Fatal("failed to find element e4")
	}
	if SortedEntitySliceSearch(slice, e7) != 6 {
		t.Fatal("failed to find element e7")
	}
}

func TestSortedEntitySliceInsertIfNotPresent(t *testing.T) {
	e1 := &Entity{ID: 1}
	e2 := &Entity{ID: 2}
	e3 := &Entity{ID: 3}
	slice := []*Entity{e1, e2, e3}
	SortedEntitySliceInsertIfNotPresent(&slice, e1)
	SortedEntitySliceInsertIfNotPresent(&slice, e2)
	SortedEntitySliceInsertIfNotPresent(&slice, e3)
	if len(slice) != 3 {
		t.Fatal("inserted when already present")
	}
}

func TestSortedEntitySliceRemove(t *testing.T) {
	e1 := &Entity{ID: 1}
	e2 := &Entity{ID: 2}
	e3 := &Entity{ID: 3}
	e4 := &Entity{ID: 4}
	slice := []*Entity{e1, e2, e3, e4}
	SortedEntitySliceRemove(&slice, e1)
	if !(slice[0].ID == 2 &&
		slice[1].ID == 3 &&
		slice[2].ID == 4) {
		t.Fatal("did not remove 0 properly")
	}
	slice = []*Entity{e1, e2, e3, e4}
	SortedEntitySliceRemove(&slice, e2)
	if !(slice[0].ID == 1 &&
		slice[1].ID == 3 &&
		slice[2].ID == 4) {
		t.Fatal("did not remove 1 properly")
	}
	slice = []*Entity{e1, e2, e3, e4}
	SortedEntitySliceRemove(&slice, e3)
	if !(slice[0].ID == 1 &&
		slice[1].ID == 2 &&
		slice[2].ID == 4) {
		t.Fatal("did not remove 2 properly")
	}
	slice = []*Entity{e1, e2, e3, e4}
	SortedEntitySliceRemove(&slice, e4)
	if !(slice[0].ID == 1 &&
		slice[1].ID == 2 &&
		slice[2].ID == 3) {
		t.Fatal("did not remove 3 properly")
	}
}
