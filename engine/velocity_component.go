/**
  *
  * Component for the velocity of an entity
  *
  *
**/

package engine

import (
	"errors"
	"fmt"
)

type VelocityComponent struct {
	Data [MAX_ENTITIES][2]float32
	em   *EntityManager
}

func (c *VelocityComponent) SafeGet(e EntityToken) ([2]float32, error) {
	if !c.em.lockEntity(e) {
		return [2]float32{}, errors.New(fmt.Sprintf("%v no longer exists", e))
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, nil
}
