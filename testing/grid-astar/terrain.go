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

func MakeTerrain(start Position, end Position) [][]int {
	t := make([][]int, H)
	for y := 0; y < H; y++ {
		t[y] = make([]int, W)
		for x := 0; x < W; x++ {
			if rand.Float64() < DENSITY {
				t[y][x] = OBSTACLE
			} else {
				t[y][x] = OPEN
			}
		}
	}
	t[start.X][start.Y] = START
	t[end.X][end.Y] = END
	return t
}

func PrintTerrain(t [][]int) {
	for y := H - 1; y >= 0; y-- {
		for x := 0; x < W; x++ {
			fmt.Printf("%s", CellString(t[y][x]))
		}
		fmt.Println()
	}
	fmt.Println()
}
