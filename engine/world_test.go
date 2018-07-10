package engine

import (
	"testing"
)

func TestCanConstructWorld(t *testing.T) {
	w := NewWorld(1024, 1024)
	if w == nil {
		t.Fatal("NewWorld() was nil")
	}
}
