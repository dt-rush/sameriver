/**
  *
  * Component for a piece of running logic attached to an entity
  *
  *
**/

package engine

type LogicComponent struct {
	Data [MAX_ENTITIES]LogicUnit
	em   *EntityManager
}

func (c *LogicComponent) SafeGet(e EntityToken) (LogicUnit, bool) {
	if !c.em.lockEntity(e) {
		return LogicUnit{}, false
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, true
}
