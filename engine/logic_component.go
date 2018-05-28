/**
  *
  * Component for a piece of running logic attached to an entity
  *
  *
**/

package engine

import (
	"errors"
	"fmt"
)

type LogicComponent struct {
	Data [MAX_ENTITIES]LogicUnit
	em   *EntityManager
}

func (c *LogicComponent) SafeGet(e EntityToken) (LogicUnit, error) {
	if !c.em.lockEntity(e) {
		return LogicUnit{}, errors.New(fmt.Sprintf("%v no longer exists", e))
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, nil
}
