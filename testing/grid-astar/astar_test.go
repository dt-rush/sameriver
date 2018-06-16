package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkAstarTest(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	g := MakeTerrain(W, H)
	pc := NewPathComputer(&g)
	points := make([][2]Position, 0)
	N := 4096
	for i := 0; i < N; i++ {
		var start Position
		startValid := false
		var end Position
		endValid := false
		for !startValid {
			start = Position{rand.Intn(W), rand.Intn(H)}
			startValid = g.Cells[start.X][start.Y] != OBSTACLE
		}
		for !endValid {
			end = Position{rand.Intn(W), rand.Intn(H)}
			endValid = g.Cells[end.X][end.Y] != OBSTACLE && end != start
		}
		points = append(points, [2]Position{start, end})
	}
	b.ResetTimer()
	for i := 0; i < N; i++ {
		pc.Path(points[i][0], points[i][1])
	}
}
