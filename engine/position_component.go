/**
  *
  * Component for the position of an entity
  *
  *
**/

package engine

type PositionComponent struct {
	Data [MAX_ENTITIES][2]int16
	em   *EntityManager
}

func (c *PositionComponent) SafeGet(e EntityToken) ([2]int16, bool) {
	if !c.em.lockEntity(e) {
		return [2]int16{}, false
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, true
}
