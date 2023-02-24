package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
	"testing"
)

func TestScreenCoords(t *testing.T) {
	s := GameScreen{W: 100, H: 100}
	pos := Vec2D{10, 10}
	box := Vec2D{5, 5}
	expected := sdl.Rect{10, 85, 5, 5}
	rect := s.ScreenSpaceRect(&pos, &box)
	if !(rect.X == expected.X &&
		rect.Y == expected.Y &&
		rect.W == expected.W &&
		rect.H == expected.H) {
		t.Fatal("screenspace conversion done wrong")
	}
}
