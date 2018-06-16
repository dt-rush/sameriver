package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func ScreenSpaceToWorldSpace(p Vec2D) (x int, y int) {
	return int(WORLD_WIDTH * (p.X / float64(WINDOW_WIDTH))),
		int(WORLD_HEIGHT * (1.0 - p.Y/float64(WINDOW_HEIGHT)))
}

func WorldSpaceToScreenSpace(p Vec2D) (x int, y int) {
	return int(WINDOW_WIDTH * (p.X / float64(WORLD_WIDTH))),
		int(WINDOW_HEIGHT * (1.0 - p.Y/float64(WORLD_HEIGHT)))
}

func MouseButtonEventToVec2D(me *sdl.MouseButtonEvent) Vec2D {
	x, y := ScreenSpaceToWorldSpace(Vec2D{float64(me.X), float64(me.Y)})
	return Vec2D{float64(x), float64(y)}
}

func MouseMotionEventToVec2D(me *sdl.MouseMotionEvent) Vec2D {
	x, y := ScreenSpaceToWorldSpace(Vec2D{float64(me.X), float64(me.Y)})
	return Vec2D{float64(x), float64(y)}
}
