/**
  *
  * Component for whether an entity is active (if inactive, no system
  * should operate on its components)
  *
  *
**/

package engine

type ActiveComponent struct {
	Data [MAX_ENTITIES]bool
	em   *EntityManager
}

func (c *ActiveComponent) SafeGet(e EntityToken) (bool, bool) {
	if !c.em.lockEntity(e) {
		return false, false
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, true
}
