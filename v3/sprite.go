package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	Texture *sdl.Texture
	Frame   uint8
	Visible bool
	Flip    sdl.RendererFlip
}
