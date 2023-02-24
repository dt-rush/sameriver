package sameriver

import (
	"testing"
)

func TestTagListToSlice(t *testing.T) {
	t1 := &TagList{
		tags: map[string]bool{
			"hello": true,
			"world": true,
		},
	}
	slice := t1.AsSlice()
	if !(t1.Has(slice...) && len(slice) == len(t1.tags)) {
		t.Fatal("did not convert to slice properly")
	}
}

func TestTagListCopy(t *testing.T) {
	t1 := TagList{
		tags: map[string]bool{
			"hello": true,
			"world": true,
		},
	}
	t2 := t1.CopyOf()
	same := t2.Has(t1.AsSlice()...) && t1.Has(t2.AsSlice()...)
	if !same {
		t.Fatal("didn't Copy properly")
	}
	t1.Remove("hello")
	if !t2.Has("hello") {
		t.Fatal("taglists are still coupled to same underlying map")
	}
}
