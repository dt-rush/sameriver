/**
  *
  * Component for the color of an entity
  *
  *
**/

package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ColorComponent struct {
	Data [MAX_ENTITIES]sdl.Color
	em   *EntityManager
}

func (c *ColorComponent) SafeGet(e EntityToken) (sdl.Color, bool) {
	if !c.em.lockEntity(e) {
		return sdl.Color{}, false
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, true
}
