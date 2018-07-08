package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type BaseComponentSet struct {
	TagList *TagList

	Logic *LogicUnit

	Box    *sdl.Rect
	Sprite *Sprite

	Position       *Vec2D
	Velocity       *Vec2D
	MovementTarget *Vec2D
	Steer          *float64
}
