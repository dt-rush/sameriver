package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type World struct {
	e         *Entity
	obstacles []Rect2D

	dif *DiffusionMap

	param int
}

func NewWorld(r *sdl.Renderer) *World {
	w := World{}
	w.param = 0
	w.dif = NewDiffusionMap(r, &w.obstacles, 100*time.Millisecond)

	return &w
}

func (w *World) Update() {
	if w.e != nil {
		w.e.Update()
	}
	select {
	case _ = <-w.dif.tick.C:
		if w.e != nil {
			w.dif.Diffuse(w.e.pos)
		}
	default:
	}
}

func (w *World) ClearObstacles() {
	w.obstacles = w.obstacles[:0]
}
