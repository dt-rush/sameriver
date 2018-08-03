package engine

import (
	"testing"
)

func TestRenderLayerActivateDeactivate(t *testing.T) {
	l := NewRenderLayer("null", 0, testingNullRenderF)
	if !l.IsActive() {
		t.Fatal("new layers should be active")
	}
	l.Deactivate()
	if l.IsActive() {
		t.Fatal("Deactivate() did not deactivate")
	}
	l.Activate()
	if !l.IsActive() {
		t.Fatal("Activate() did not activate")
	}
}

func TestRenderLayerName(t *testing.T) {
	l := NewRenderLayer("hello", 0, nil)
	if l.Name() != "hello" {
		t.Fatal("Name() doesn't return name given to constructor")
	}
}
