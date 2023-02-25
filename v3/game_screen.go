package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
)

type GameScreen struct {
	W int
	H int
}

func (s *GameScreen) ScreenSpaceY(y int) int {
	return s.H - y
}

// expects pos shifted to bottom-left corner
func (s *GameScreen) ScreenSpaceRect(pos *Vec2D, box *Vec2D) *sdl.Rect {
	return &sdl.Rect{
		int32(pos.X),
		int32(float64(s.H) - pos.Y - box.Y),
		int32(box.X),
		int32(box.Y),
	}
}

func (s *GameScreen) DrawRect(r *sdl.Renderer, pos *Vec2D, box *Vec2D) {
	r.DrawRect(s.ScreenSpaceRect(pos, box))
}

func (s *GameScreen) FillRect(r *sdl.Renderer, pos *Vec2D, box *Vec2D) {
	r.FillRect(s.ScreenSpaceRect(pos, box))
}
