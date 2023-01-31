package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
	"testing"
	"time"
)

func TestLayeredRendererAddRemove(t *testing.T) {
	lr := NewLayeredRenderer()
	// add l8
	l8 := NewRenderLayer("l8", 8, testingNullRenderF)
	lr.AddLayer(l8)
	if lr.NumLayers() != 1 {
		t.Fatal("did not insert first layer")
	}
	// add l2
	l2 := NewRenderLayer("l2", 2, testingNullRenderF)
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
	lA := NewRenderLayer("la", 0, func(w *sdl.Window, r *sdl.Renderer) {
		time.Sleep(1 * time.Millisecond)
		a++
	})
	lr.AddLayer(lA)
	// add layer B
	b := 0
	lB := NewRenderLayer("lb", 0, func(w *sdl.Window, r *sdl.Renderer) {
		time.Sleep(1 * time.Millisecond)
		b++
	})
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

func TestLayeredRendererZIndexing(t *testing.T) {
	lr := NewLayeredRenderer()
	// add layer A
	x := 0
	lA := NewRenderLayer("la", 0, func(w *sdl.Window, r *sdl.Renderer) {
		time.Sleep(1 * time.Millisecond)
		x++
	})
	// add layer B
	lB := NewRenderLayer("lb", 4, func(w *sdl.Window, r *sdl.Renderer) {
		time.Sleep(1 * time.Millisecond)
		x *= 3
	})
	// add B first, even though A is a lower Zindex and should render first
	lr.AddLayer(lB)
	lr.AddLayer(lA)
	// render and get elapsed time
	elapsed := lr.Render(nil, nil)
	if x != 3 {
		t.Fatal("Render() did not render in order")
	}
	if elapsed < 2 {
		t.Fatal("did not return elapsed time properly")
	}
}

func TestLayeredRendererByName(t *testing.T) {
	lr := NewLayeredRenderer()
	// add layer A
	a := 0
	lA := NewRenderLayer("la", 0, func(w *sdl.Window, r *sdl.Renderer) {
		a++
	})
	lr.AddLayer(lA)
	// add layer B
	b := 0
	lB := NewRenderLayer("lb", 0, func(w *sdl.Window, r *sdl.Renderer) {
		b++
	})
	lr.AddLayer(lB)
	// get layer A by name
	l, err := lr.GetLayerByName("la")
	if err != nil {
		t.Fatal(err)
	}
	if l != lA {
		t.Fatal("GetLayerByName() returned wrong layer")
	}
	// remove layer B by name
	lr.RemoveLayerByName("lb")
	if lr.NumLayers() != 1 {
		t.Fatal("did not remove layer")
	}
	// render and get elapsed time
	lr.Render(nil, nil)
	if !(a == 1 && b == 0) {
		t.Fatal("did not remove right layer by name")
	}
	// check failure for non-existent layer
	l, err = lr.GetLayerByName("!!!!!!!!!!")
	if err == nil {
		t.Fatal("should have thrown error for GetLayerByName when name doesn't exist")
	}
}
