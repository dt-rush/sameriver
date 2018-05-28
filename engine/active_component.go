/**
  *
  * Component for whether an entity is active (if inactive, no system
  * should operate on its components)
  *
  *
**/

package engine

import (
	"errors"
	"fmt"
)

type ActiveComponent struct {
	Data [MAX_ENTITIES]bool
	em   *EntityManager
}

func (c *ActiveComponent) SafeGet(e EntityToken) (bool, error) {
	if !c.em.lockEntity(e) {
		return false, errors.New(fmt.Sprintf("%v no longer exists", e))
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, nil
}
