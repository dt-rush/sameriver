package main

import (
	"fmt"
	"github.com/beefsack/go-astar"
	"time"
)

func (w *World) ComputeEntityPath() {
	if w.e != nil && w.e.moveTarget != nil {
		t0 := time.Now()
		path, distance, found := astar.Path(
			w.m.CellAt(w.e.pos),
			w.m.CellAt(*w.e.moveTarget))
		fmt.Printf("path calculation took %.3f ms\n",
			float64(time.Since(t0).Nanoseconds())/float64(1e6))
		if found {
			w.e.distance = distance
			cellsPath := make([]Position, len(path))
			for i, pather := range path {
				cellsPath[i] = pather.(*WorldMapCell).pos
			}
			w.e.path = cellsPath
		}
	}
}

func (w *World) MoveEntity() {
	if w.e != nil && w.e.moveTarget != nil {
		last_ix := len(w.e.path) - 1
		w.e.pos = w.e.path[last_ix]
		w.e.path = w.e.path[:last_ix]
		if len(w.e.path) == 0 {
			w.e.moveTarget = nil
			w.e.path = nil
		}
	}

}
