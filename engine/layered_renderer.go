package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"sort"
	"time"
)

type LayeredRenderer struct {
	layers []*RenderLayer
}

func NewLayeredRenderer() *LayeredRenderer {
	return &LayeredRenderer{}
}

// sorted insert to layers array according to layer of provided RenderLayer
func (lr *LayeredRenderer) AddLayer(l *RenderLayer) {
	index := sort.Search(len(lr.layers),
		func(i int) bool { return lr.layers[i].z > l.z })
	lr.layers = append(lr.layers, nil)
	copy(lr.layers[index+1:], lr.layers[index:])
	lr.layers[index] = l
}

// remove a given RenderLayer
func (lr *LayeredRenderer) RemoveLayer(l *RenderLayer) {
	// NOTE: we subtract 1 since Search() will return the index after the
	// element's index if already found
	index := sort.Search(len(lr.layers),
		func(i int) bool { return lr.layers[i].z > l.z }) - 1
	// NOTE: we set the last element to nil so that we don't memory leak
	// by leaving pointers around in the slice capacity tail
	if lr.layers[index] == l {
		copy(lr.layers[index:], lr.layers[index+1:])
		lr.layers[len(lr.layers)-1] = nil
		lr.layers = lr.layers[:len(lr.layers)-1]
	}
}

// render each layer, if active, returning elapsed time in ms
// (since lr.layers is sorted by layer.z, iterating is automatically from
// lowest to highest z)
func (lr *LayeredRenderer) Render(
	w *sdl.Window, r *sdl.Renderer) (elapsed float64) {

	t0 := time.Now()
	for _, l := range lr.layers {
		if l.active {
			l.Render(w, r)
		}
	}
	return float64(time.Since(t0).Nanoseconds()) / 1.0e6
}

// get the number of layers currently active
func (lr *LayeredRenderer) NumLayers() int {
	return len(lr.layers)
}
