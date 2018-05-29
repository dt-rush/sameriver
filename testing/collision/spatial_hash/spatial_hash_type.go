package main

import (
	"bytes"
	"fmt"
	"github.com/dt-rush/donkeys-qquest/engine"
)

type EntityPosition struct {
	entity   engine.EntityToken
	position [2]int16
}

type SpatialHash [GRID][GRID][]EntityPosition

func NewSpatialHash() *SpatialHash {
	var buckets = SpatialHash{}
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			buckets[y][x] = make([]EntityPosition,
				UPPER_ESTIMATE_ENTITIES_PER_SQUARE)
		}
	}
	return &buckets
}

func (h *SpatialHash) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for y := 0; y < GRID; y++ {
		for x := 0; x < GRID; x++ {
			buffer.WriteString(fmt.Sprintf(
				"CELL(%d, %d): %v", x, y, h[y][x]))
			if !(y == GRID-1 && x == GRID-1) {
				buffer.WriteString(",")
			} else if !(y == GRID-1) {
				buffer.WriteString("\n")
			}
		}
	}
	buffer.WriteString("]")
	return buffer.String()
}
