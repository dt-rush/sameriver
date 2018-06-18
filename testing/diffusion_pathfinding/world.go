package main

import (
	"fmt"
	"math/rand"
	"time"
)

type World struct {
	g         *Game
	e         *Entity
	obstacles []Rect2D

	dm *DiffusionMap
	pc *PathComputer

	param int
}

func NewWorld(g *Game) *World {
	w := World{g: g}
	w.param = 0
	w.dm = NewDiffusionMap(g.r, &w.obstacles, 16*time.Millisecond)
	w.pc = NewPathComputer(w.dm)

	return &w
}

func (w *World) Update() {
	if w.e != nil {
		w.e.Update()
	}
	select {
	case _ = <-w.dm.tick.C:
		if w.e != nil {
			t0 := time.Now()
			w.dm.Diffuse(w.e.pos)
			msg := fmt.Sprintf("diffusion took %.1f ms",
				float64(time.Since(t0).Nanoseconds()/1e6)/1000.0)
			w.g.ui.UpdateMsg(4, msg)

		}
	default:
	}
}

func (w *World) AddObstacle(pos Vec2D) {
	o := CenteredSquare(pos, OBSTACLESZ)
	w.obstacles = append(w.obstacles, o)
	w.dm.AddObstacle(o)
}

func (w *World) ClearObstacles() {
	w.obstacles = w.obstacles[:0]
	w.dm.ClearObstacles()
}

func (w *World) RandomObstacles() {
	for i := 0; i < 20; i++ {
		var x, y = -999, -999
		for !w.dm.InGrid(
			x/GRIDCELL_WORLD_W,
			y/GRIDCELL_WORLD_H) {
			x = rand.Intn(GRID_WORLD_DIMENSION)

			y = rand.Intn(GRID_WORLD_DIMENSION)
		}
		o := CenteredSquare(Vec2D{float64(x), float64(y)}, OBSTACLESZ)
		w.obstacles = append(w.obstacles, o)
		w.dm.AddObstacle(o)
	}
}
