package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func ScreenSpaceToWorldSpace(p Point2D) (x int, y int) {
	return int(WORLD_WIDTH * (p.X / float64(WINDOW_WIDTH))),
		int(WORLD_HEIGHT * (1.0 - p.Y/float64(WINDOW_HEIGHT)))
}

func WorldSpaceToScreenSpace(p Point2D) (x int, y int) {
	return int(WINDOW_WIDTH * (p.X / float64(WORLD_WIDTH))),
		int(WINDOW_HEIGHT * (1.0 - p.Y/float64(WORLD_HEIGHT)))
}

func MouseButtonEventToPoint2D(me *sdl.MouseButtonEvent) Point2D {
	x, y := ScreenSpaceToWorldSpace(Point2D{float64(me.X), float64(me.Y)})
	return Point2D{float64(x), float64(y)}
}

func MouseMotionEventToPoint2D(me *sdl.MouseMotionEvent) Point2D {
	x, y := ScreenSpaceToWorldSpace(Point2D{float64(me.X), float64(me.Y)})
	return Point2D{float64(x), float64(y)}
}
