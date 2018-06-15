package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	CELL_WATER  = iota
	CELL_SAND   = iota
	CELL_GRASS  = iota
	CELL_FOREST = iota
)

var WorldMapCellTransitionCostFuncs = []func(otherKind int) float64{
	// WATER
	func(otherKind int) float64 {
		var cost float64
		switch otherKind {
		case CELL_WATER:
			cost = 10
		case CELL_SAND:
			cost = 1
		case CELL_GRASS:
			cost = 1
		case CELL_FOREST:
			cost = 1.8
		}
		return cost
	},
	// SAND
	func(otherKind int) float64 {
		var cost float64
		switch otherKind {
		case CELL_WATER:
			cost = 3
		case CELL_SAND:
			cost = 1
		case CELL_GRASS:
			cost = 1
		case CELL_FOREST:
			cost = 1.5
		}
		return cost
	},
	// GRASS
	func(otherKind int) float64 {
		var cost float64
		switch otherKind {
		case CELL_WATER:
			cost = 3
		case CELL_SAND:
			cost = 1
		case CELL_GRASS:
			cost = 1
		case CELL_FOREST:
			cost = 1.5
		}
		return cost
	},
	// FOREST
	func(otherKind int) float64 {
		var cost float64
		switch otherKind {
		case CELL_WATER:
			cost = 3
		case CELL_SAND:
			cost = 0.5
		case CELL_GRASS:
			cost = 0.5
		case CELL_FOREST:
			cost = 1
		}
		return cost
	},
}

type WorldMapCell struct {
	m     *WorldMap
	rep   string
	kind  int
	pos   Position
	color sdl.Color
	data  interface{}
}

func NewWorldMapCell(m *WorldMap, rep string, kind int, color sdl.Color) WorldMapCell {
	return WorldMapCell{
		m:     m,
		rep:   rep,
		kind:  kind,
		color: color}
}

func (m *WorldMap) WaterCell(depth int) WorldMapCell {
	return NewWorldMapCell(m, "o", CELL_WATER,
		sdl.Color{R: 0, G: 0, B: uint8(48 + 16*depth)})
}

func (m *WorldMap) SandCell() WorldMapCell {
	return NewWorldMapCell(m, ".", CELL_SAND,
		sdl.Color{R: 182, G: 182, B: 0})
}

func (m *WorldMap) GrassCell() WorldMapCell {
	return NewWorldMapCell(m, ".", CELL_GRASS,
		sdl.Color{R: 0, G: 182, B: 0})
}

type ForestCellData struct {
	density int
}

func (m *WorldMap) ForestCell(density int) WorldMapCell {
	c := NewWorldMapCell(m, "#", CELL_FOREST,
		sdl.Color{R: 0, G: uint8(48 + density*16), B: 0})
	c.data = ForestCellData{density}
	return c
}

func (c1 *WorldMapCell) CostToTransitionTo(c2 *WorldMapCell) float64 {
	return WorldMapCellTransitionCostFuncs[c1.kind](c2.kind)
}