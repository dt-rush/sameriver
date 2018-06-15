package main

/*==========================================================================
module:      Path.cpp
programmer:  Marco Pinter
             Copyright (c) Marco Pinter, 2001.

A* algorithm for path-finding, written from 8/20/2000.

   This code is "open source".
   You may incorporate it royalty-free in any project,
   without restriction.  If possible, please maintain a reference
   to the author, and this message, within the source.

=========================================================================-*/

const MAXTILEDIST = 40
const MAXTILE_XLEN = (MAXTILEDIST * 3 / 2)
const MAXTILE_YLEN = (MAXTILEDIST * 3 / 2)
const TILE_CENTER = (MAXTILE_XLEN / 2)

type PointSmall struct {
	x   uint8
	y   uint8
	dir uint8
}

func NewPointSmall(x int, y int, dir int) PointSmall {
	return PointSmall{uint8(x), uint8(y), uint8(dir)}
}

type LineSegment struct {
}

type PathTileNodeType struct {
	bOpen              bool
	direction          int
	costFromStart      int
	totalCost          int
	parentDir          int
	parentX            uint8
	parentY            uint8
	prevOpen, nextOpen PointSmall
}

type PathCalculator struct {
}
