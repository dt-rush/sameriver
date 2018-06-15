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
