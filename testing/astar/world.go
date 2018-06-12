package main

import (
	"fmt"
)

type World struct {
	m *WorldMap
	e *Entity
}

func (w *World) RegenMap() {
	w.m = GenerateWorldMap()
	fmt.Printf("seed: %d\n", w.m.seed)
}
