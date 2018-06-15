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

	start := Position{0, 0}
	end := Position{W - 1, H - 1}
	t := MakeTerrain(start, end)

	fmt.Println("Finding path...")

	path := FindPath(start, end, t)
	for _, p := range path {
		t[p.X][p.Y] = PATH
	}

	PrintTerrain(t)
}
