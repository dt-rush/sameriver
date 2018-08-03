package engine

import (
	"testing"
)

func TestRenderLayerActivateDeactivate(t *testing.T) {
	l := NewRenderLayer(testingNullRenderF, 0)
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
