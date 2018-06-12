package main

import (
	"fmt"
	"time"
)

type WorldMap struct {
	seed  int64
	cells [WORLD_CELLWIDTH][WORLD_CELLHEIGHT]WorldMapCell
}

func GenerateWorldMap() *WorldMap {
	m := WorldMap{}
	m.seed = time.Now().UnixNano()
	terrain := PerlinNoiseInt2D(
		WORLD_CELLWIDTH, WORLD_CELLHEIGHT, 16,
		2.0, 2.0, 3,
		m.seed)
	water := PerlinNoiseInt2D(
		WORLD_CELLWIDTH, WORLD_CELLHEIGHT, 32,
		4.0, 2.0, 3,
		m.seed)
	world := OpPerlins(terrain, water, func(a float64, b float64) float64 {
		x := (a + (a + 0.3) - b) / 2
		if x < 0 {
			return 0
		} else if x > 1 {
			return 1
		} else {
			return x
		}
	})
	for y := 0; y < WORLD_CELLHEIGHT; y++ {
		for x := 0; x < WORLD_CELLWIDTH; x++ {
			var v = world[y][x]
			var c WorldMapCell
			if v < 0.4 {
				depth := int(v / 0.1)
				c = m.WaterCell(depth)
			} else if v < 0.45 {
				c = m.SandCell()
			} else if v < 0.55 {
				c = m.GrassCell()
			} else {
				density := int((1.0 - v) / 0.1)
				c = m.ForestCell(density)
			}
			c.pos = Position{x, y}
			m.cells[y][x] = c
		}
	}
	return &m
}

func (m *WorldMap) CellAt(pos Position) *WorldMapCell {
	return &m.cells[pos.y][pos.x]
}

func (m *WorldMap) Print() {
	for y := 0; y < WORLD_CELLHEIGHT; y++ {
		for x := 0; x < WORLD_CELLWIDTH; x++ {
			fmt.Printf("%s", m.cells[y][x].rep)
		}
		fmt.Println()
	}
}
