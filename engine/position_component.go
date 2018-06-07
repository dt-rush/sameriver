/**
  *
  * Component for the position of an entity
  *
  *
**/

package engine

import (
	"errors"
	"fmt"
)

type PositionComponent struct {
	Data  [MAX_ENTITIES][2]int16
	locks [MAX_ENTITIES]*ABRWPQL
	em    *EntityManager
}

func (c *PositionComponent) SafeGet(e EntityToken) ([2]int16, error) {
	if !c.em.lockEntityComponent(e, POSITION_COMPONENT) {
		return [2]int16{}, errors.New(fmt.Sprintf("%+v no longer exists", e))
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, nil
}
