/**
  *
  * Component for the velocity of an entity
  *
  *
**/

package engine

type VelocityComponent struct {
	Data [MAX_ENTITIES][2]float32
	em   *EntityManager
}

func (c *VelocityComponent) SafeGet(e EntityToken) ([2]float32, bool) {
	if !c.em.lockEntity(e) {
		return [2]float32{}, false
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, true
}
