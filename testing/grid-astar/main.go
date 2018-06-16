package main

import (
	"github.com/fatih/color"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	color.NoColor = false
}

func main() {

	g := MakeTerrain(W, H)
	pc := NewPathComputer(&g)
	N := 1024
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
		path := pc.Path(start, end)
		PrintTerrainWithPath(g, start, end, path)
		time.Sleep(100 * time.Millisecond)
	}
}
