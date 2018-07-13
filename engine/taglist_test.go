package engine

import (
	"fmt"
	"testing"
)

func TestTagListCopy(t *testing.T) {
	t1 := &TagList{Tags: []string{"hello", "world"}}
	t2 := t1.Copy()
	if t1 == t2 {
		t.Fatal("Copy returned pointer to same taglist")
	}
	if fmt.Sprintf("%v", t1.Tags) != fmt.Sprintf("%v", t2.Tags) {
		t.Fatal("Copy returned pointer to same taglist")
	}
}
