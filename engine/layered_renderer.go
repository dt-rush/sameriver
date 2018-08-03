package engine

import (
	"errors"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"sort"
	"time"
)

type LayeredRenderer struct {
	layers []*RenderLayer
	names  map[string]*RenderLayer
}

func NewLayeredRenderer() *LayeredRenderer {
	return &LayeredRenderer{
		names: make(map[string]*RenderLayer),
	}
}

// sorted insert to layers array according to layer of provided RenderLayer
func (lr *LayeredRenderer) AddLayer(l *RenderLayer) {
	// sorted insert into array
	index := sort.Search(len(lr.layers),
		func(i int) bool { return lr.layers[i].z > l.z })
	lr.layers = append(lr.layers, nil)
	copy(lr.layers[index+1:], lr.layers[index:])
	lr.layers[index] = l
	// add to name map
	lr.names[l.name] = l
}

// remove a given RenderLayer
func (lr *LayeredRenderer) RemoveLayer(l *RenderLayer) {
	// NOTE: we subtract 1 since Search() will return the index after the
	// element's index if already found
	index := sort.Search(len(lr.layers),
		func(i int) bool { return lr.layers[i].z > l.z }) - 1
	// there may be multiple layers with the same z, so move backward
	// til we no longer have any layers left with this z in the array, til we
	// exhaust the array, or we find our layer
	for i := index; i >= 0 && lr.layers[i].z == l.z; i-- {
		if lr.layers[i] == l {
			// if found:
			// remove from name map
			delete(lr.names, l.name)
			// remove from array
			copy(lr.layers[i:], lr.layers[i+1:])
			// NOTE: we set the last element to nil so that we don't memory leak
			// by leaving pointers around in the slice capacity tail
			lr.layers[len(lr.layers)-1] = nil
			lr.layers = lr.layers[:len(lr.layers)-1]
			break
		}
	}
}

// remove a layer by name
func (lr *LayeredRenderer) RemoveLayerByName(name string) {
	if l, ok := lr.names[name]; ok {
		lr.RemoveLayer(l)
	}
}

// get a layer by name
func (lr *LayeredRenderer) GetLayerByName(name string) (*RenderLayer, error) {
	Logger.Println(lr.names)
	if l, ok := lr.names[name]; ok {
		return l, nil
	} else {
		return nil, errors.New(fmt.Sprintf("layer %s not found", name))
	}
}

// render each layer, if active, returning elapsed time in ms
// (since lr.layers is sorted by layer.z, iterating is automatically from
// lowest to highest z)
func (lr *LayeredRenderer) Render(
	w *sdl.Window, r *sdl.Renderer) (elapsed float64) {

	t0 := time.Now()
	for _, l := range lr.layers {
		if l.IsActive() {
			l.Render(w, r)
		}
	}
	return float64(time.Since(t0).Nanoseconds()) / 1.0e6
}

// get the number of layers currently active
func (lr *LayeredRenderer) NumLayers() int {
	return len(lr.layers)
}
