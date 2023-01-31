package sameriver

import (
	"testing"
)

func TestTagListToSlice(t *testing.T) {
	t1 := &TagList{
		Tags: map[string]bool{
			"hello": true,
			"world": true,
		},
	}
	slice := t1.ToSlice()
	if !(t1.Has(slice...) && len(slice) == len(t1.Tags)) {
		t.Fatal("did not convert to slice properly")
	}
}

func TestTagListCopy(t *testing.T) {
	t1 := &TagList{
		Tags: map[string]bool{
			"hello": true,
			"world": true,
		},
	}
	t2 := t1.Copy()
	if t1 == t2 {
		t.Fatal("Copy returned pointer to same taglist")
	}
	if !(t2.Has(t1.ToSlice()...) &&
		t1.Has(t2.ToSlice()...)) {
		t.Fatal("didn't Copy properly")
	}
}
