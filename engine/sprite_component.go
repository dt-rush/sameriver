/**
  *
  *
  *
  *
**/

package engine

type SpriteComponent struct {
	Data [MAX_ENTITIES]Sprite
	em   *EntityManager
}

func (c *SpriteComponent) SafeGet(e EntityToken) (Sprite, bool) {
	if !c.em.lockEntity(e) {
		return Sprite{}, false
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, true
}
