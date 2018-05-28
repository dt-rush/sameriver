/**
  *
  * Component for health of an entity
  *
  *
**/

package engine

import (
	"errors"
	"fmt"
)

type HealthComponent struct {
	Data [MAX_ENTITIES]uint8
	em   *EntityManager
}

func (c *HealthComponent) SafeGet(e EntityToken) (uint8, error) {
	if !c.em.lockEntity(e) {
		return 0, errors.New(fmt.Sprintf("%v no longer exists", e))
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, nil
}
