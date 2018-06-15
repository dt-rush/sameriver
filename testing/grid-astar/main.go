package main

import (
	"fmt"
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

	start := Position{0, 0}
	end := Position{W - 1, H - 1}
	g.Cells[start.X][start.Y] = START
	g.Cells[end.X][end.Y] = END

	pc := NewPathComputer(&g)

	fmt.Println("Finding path...")

	path := pc.Path(start, end)
	for _, p := range path {
		g.Cells[p.X][p.Y] = PATH
	}
	if path == nil {
		fmt.Println("NO PATH")
	}

	PrintTerrain(g)
}
