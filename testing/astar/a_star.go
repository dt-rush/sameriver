package main

import (
	"github.com/beefsack/go-astar"
	"math"
)

func (c *WorldMapCell) PathNeighbors() []astar.Pather {
	neighbors := make([]astar.Pather, 0)
	for dy := -1; dy <= 1; dy++ {
		if c.pos.y+dy < 0 ||
			c.pos.y+dy > WORLD_CELLHEIGHT-1 {
			continue
		}
		for dx := -1; dx <= 1; dx++ {
			if c.pos.x+dx < 0 ||
				c.pos.x+dx > WORLD_CELLWIDTH-1 {
				continue
			}
			neighbors = append(neighbors,
				&c.m.cells[c.pos.y+dy][c.pos.x+dx])
		}
	}
	return neighbors
}

func (c *WorldMapCell) PathNeighborCost(to astar.Pather) float64 {
	dx := math.Abs(float64(c.pos.x - to.(*WorldMapCell).pos.x))
	dy := math.Abs(float64(c.pos.y - to.(*WorldMapCell).pos.y))
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance * c.CostToTransitionTo(to.(*WorldMapCell))
}

func (c *WorldMapCell) PathEstimatedCost(to astar.Pather) float64 {
	dx := math.Abs(float64(c.pos.x - to.(*WorldMapCell).pos.x))
	dy := math.Abs(float64(c.pos.y - to.(*WorldMapCell).pos.y))
	return dx + dy
}
