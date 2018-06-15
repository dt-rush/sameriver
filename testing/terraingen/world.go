package main

import (
	"fmt"
)

type World struct {
	m *WorldMap
	e *Entity
	c *PathCalculator
}

func NewWorld() *World {
	w := World{}
	w.RegenMap()
	w.c = NewPathCalculator(w.m)
	return &w
}

func (w *World) RegenMap() {
	w.m = GenerateWorldMap()
	fmt.Printf("seed: %d\n", w.m.seed)
}
