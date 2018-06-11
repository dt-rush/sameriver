package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type BaseComponentSet struct {
	Box            *sdl.Rect
	Sprite         *Sprite
	TagList        *TagList
	Velocity       *Vec2D
	MovementTarget *Vec2D
	Steer          *float32
}
