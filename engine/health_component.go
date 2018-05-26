/**
  *
  * Component for health of an entity
  *
  *
**/

package engine

type HealthComponent struct {
	Data [MAX_ENTITIES]uint8
	em   *EntityManager
}

func (c *HealthComponent) SafeGet(e EntityToken) (uint8, bool) {
	if !c.em.lockEntity(e) {
		return 0, false
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, true
}
