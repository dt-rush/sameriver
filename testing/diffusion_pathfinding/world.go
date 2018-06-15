package main

import (
	"fmt"
	"time"
)

type World struct {
	g         *Game
	e         *Entity
	obstacles []Rect2D

	dif *DiffusionMap

	param int
}

func NewWorld(g *Game) *World {
	w := World{g: g}
	w.param = 0
	w.dif = NewDiffusionMap(g.r, &w.obstacles, 16*time.Millisecond)

	return &w
}

func (w *World) Update() {
	if w.e != nil {
		w.e.Update()
	}
	select {
	case _ = <-w.dif.tick.C:
		if w.e != nil {
			t0 := time.Now()
			w.dif.Diffuse(w.e.pos)

			msg := fmt.Sprintf("diffusion took %.1f ms",
				float64(time.Since(t0).Nanoseconds()/1e6)/1000.0)
			w.g.ui.UpdateMsg(1, msg)

		}
	default:
	}
}

func (w *World) ClearObstacles() {
	w.obstacles = w.obstacles[:0]
}
