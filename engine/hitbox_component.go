/**
  *
  * Component for the hitbox of an entity
  *
  *
**/

package engine

import (
	"errors"
	"fmt"
)

type HitBoxComponent struct {
	// TODO: consider making hitbox [2]uint8 - nothing really needs to be
	// bigger than 255...
	Data [MAX_ENTITIES][2]uint16
	em   *EntityManager
}

func (c *HitBoxComponent) SafeGet(e EntityToken) ([2]uint16, error) {
	if !c.em.lockEntity(e) {
		return [2]uint16{}, errors.New(fmt.Sprintf("%v no longer exists", e))
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, nil
}
