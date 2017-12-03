package engine


import (
	"github.com/veandco/go-sdl2/sdl"
)


/*
 * Builds and returns an SDL window and renderer object 
 * for the game to use
 */
func BuildWindowAndRenderer (window_title string, width int32, height int32) (*sdl.Window, *sdl.Renderer) {

	window, err := sdl.CreateWindow (window_title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		width,
		height,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	
	renderer, err := sdl.CreateRenderer (window,
		-1,
		sdl.RENDERER_SOFTWARE)
	if err != nil {
		panic (err)
	}

	return window, renderer
}
