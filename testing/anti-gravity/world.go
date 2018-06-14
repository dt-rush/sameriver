package main

import (
	"sync"
)

type World struct {
	e         *Entity
	obstacles []Point2D
	mutex     sync.Mutex
	param     int
}

func NewWorld() *World {
	w := World{}
	w.param = 0
	return &w
}

func (w *World) ClearObstacles() {
	w.obstacles = w.obstacles[:0]
}
