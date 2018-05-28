/**
  *
  * Component for the color of an entity
  *
  *
**/

package engine

import (
	"errors"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type ColorComponent struct {
	Data [MAX_ENTITIES]sdl.Color
	em   *EntityManager
}

func (c *ColorComponent) SafeGet(e EntityToken) (sdl.Color, error) {
	if !c.em.lockEntity(e) {
		return sdl.Color{}, errors.New(fmt.Sprintf("%v no longer exists", e))
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, nil
}
