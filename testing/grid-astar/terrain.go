package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
)

const (
	OPEN     = 0
	OBSTACLE = iota
	START    = iota
	PATH     = iota
	END      = iota
)

var bgRed = color.New(color.BgRed).SprintFunc()
var bgBlack = color.New(color.BgBlack).SprintFunc()
var bgGreen = color.New(color.BgGreen).SprintFunc()
var bgWhite = color.New(color.BgWhite).SprintFunc()
var bgCyan = color.New(color.BgCyan).SprintFunc()

func CellString(val int) (rep string) {
	switch val {
	case OPEN:
		rep = bgBlack("  ")
	case OBSTACLE:
		rep = bgRed("  ")
	case START:
		rep = bgGreen("  ")
	case PATH:
		rep = bgWhite("  ")
	case END:
		rep = bgCyan("  ")
	}
	return rep
}

func MakeTerrain(w int, h int) Grid {
	t := make([][]int, w)
	for x := 0; x < w; x++ {
		t[x] = make([]int, h)
		for y := 0; y < h; y++ {
			if rand.Float64() < DENSITY {
				t[x][y] = OBSTACLE
			} else {
				t[x][y] = OPEN
			}
		}
	}
	return Grid{w, h, t}
}

func PrintTerrain(g Grid) {
	for y := g.H - 1; y >= 0; y-- {
		for x := 0; x < g.W; x++ {
			fmt.Printf("%s", CellString(g.Cells[y][x]))
		}
		fmt.Println()
	}
	fmt.Println()
}

func PrintTerrainWithPath(
	g Grid, start Position, end Position, path []Position) {

	inPath := make(map[Position]bool)
	for _, p := range path {
		inPath[p] = true
	}
	for y := g.H - 1; y >= 0; y-- {
		for x := 0; x < g.W; x++ {
			p := Position{x, y}
			if p == start {
				fmt.Printf("%s", CellString(START))
			} else if p == end {
				fmt.Printf("%s", CellString(END))
			} else if _, ok := inPath[Position{x, y}]; ok {
				fmt.Printf("%s", CellString(PATH))
			} else {
				fmt.Printf("%s", CellString(g.Cells[y][x]))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
