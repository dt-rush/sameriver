package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (s *GameScreen) ScreenSpaceRect(pos *Vec2D, box *Vec2D) *sdl.Rect {
	return &sdl.Rect{
		int32(pos.X),
		int32(s.H - pos.Y - box.Y),
		int32(box.X),
		int32(box.Y),
	}
}
