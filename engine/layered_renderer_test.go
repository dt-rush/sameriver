package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"testing"
	"time"
)

func TestLayeredRendererAddRemove(t *testing.T) {
	lr := NewLayeredRenderer()
	// add l8
	l8 := NewRenderLayer(testingNullRenderF, 8)
	lr.AddLayer(l8)
	if lr.NumLayers() != 1 {
		t.Fatal("did not insert first layer")
	}
	// add l2
	l2 := NewRenderLayer(testingNullRenderF, 2)
	lr.AddLayer(l2)
	if lr.layers[0] != l2 {
		t.Fatal("did not insert in z-order")
	}
	if lr.NumLayers() != 2 {
		t.Fatal("did not insert second layer properly")
	}
	// remove l2
	lr.RemoveLayer(l2)
	if lr.NumLayers() != 1 {
		t.Fatal("did not remove layer")
	}
	if lr.layers[0] != l8 {
		t.Fatal("did not remove layer properly")
	}
	// remove l8
	lr.RemoveLayer(l8)
	if lr.NumLayers() != 0 {
		t.Fatal("did not remove layer")
	}
}

func TestLayeredRendererRender(t *testing.T) {
	lr := NewLayeredRenderer()
	// add layer A
	a := 0
	lA := NewRenderLayer(func(w *sdl.Window, r *sdl.Renderer) {
		time.Sleep(1 * time.Millisecond)
		a++
	}, 0)
	lr.AddLayer(lA)
	// add layer B
	b := 0
	lB := NewRenderLayer(func(w *sdl.Window, r *sdl.Renderer) {
		time.Sleep(1 * time.Millisecond)
		b++
	}, 0)
	lr.AddLayer(lB)
	// render and get elapsed time
	elapsed := lr.Render(nil, nil)
	if !(a == 1 && b == 1) {
		t.Fatal("Render() did not render each layer")
	}
	if elapsed < 2 {
		t.Fatal("did not return elapsed time properly")
	}
}
