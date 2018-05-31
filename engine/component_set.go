/*
 *
 * A set of component values which may be provided or not
 * Also functions as the place where components are defined so all the
 * supporting source files can be generated
 *
 */

package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ComponentSet struct {
	Active   *bool
	Color    *sdl.Color
	Health   *uint8
	HitBox   *[2]uint16
	Logic    *LogicUnit
	Mind     *Mind
	Position *[2]int16
	Sprite   *Sprite
	TagList  *TagList
	Velocity *[2]float32
}
