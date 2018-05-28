/**
  *
  *
  *
  *
**/

package engine

import (
	"errors"
	"fmt"
)

type SpriteComponent struct {
	Data [MAX_ENTITIES]Sprite
	em   *EntityManager
}

func (c *SpriteComponent) SafeGet(e EntityToken) (Sprite, error) {
	if !c.em.lockEntity(e) {
		return Sprite{}, errors.New(fmt.Sprintf("%v no longer exists", e))
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, nil
}
