package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ScreenSpaceConverter struct {
	W float64
	H float64
}

func (c *ScreenSpaceConverter) DrawRect(r *sdl.Renderer, pos *Vec2D, box *Vec2D) {
	r.DrawRect(&sdl.Rect{
		int32(pos.X),
		int32(c.H - pos.Y - box.Y),
		int32(box.X),
		int32(box.Y),
	})
}

func (c *ScreenSpaceConverter) FillRect(r *sdl.Renderer, pos *Vec2D, box *Vec2D) {
	r.FillRect(&sdl.Rect{
		int32(pos.X),
		int32(c.H - pos.Y - box.Y),
		int32(box.X),
		int32(box.Y),
	})
}
