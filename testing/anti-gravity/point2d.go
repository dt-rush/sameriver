package main

import (
	"math"
)

type Point2D struct {
	X float64
	Y float64
}

type Point2DPair struct {
	p1 Point2D
	p2 Point2D
}

func Distance(p1 Point2D, p2 Point2D) (dx, dy, d float64) {
	dx = float64(p2.X - p1.X)
	dy = float64(p2.Y - p1.Y)
	d = math.Sqrt(dx*dx + dy*dy)
	return dx, dy, d
}

func PointDelta(p Point2D, dx float64, dy float64) Point2D {
	return Point2D{p.X + dx, p.Y + dy}
}
