package main

type MapVertex struct {
	id        int
	pos       Point2D
	neighbors []*MapVertex
}
