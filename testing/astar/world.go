package main

type World struct {
	m *WorldMap
	e *Entity
}

func (w *World) NewWorldMap() {
	w.m = GenerateWorldMap()
}
