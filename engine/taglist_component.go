/**
  *
  * Component for an entity's taglist
  *
  *
**/

package engine

type TagListComponent struct {
	Data [MAX_ENTITIES]TagList
	em   *EntityManager
}

func (c *TagListComponent) SafeGet(e EntityToken) (TagList, bool) {
	if !c.em.lockEntity(e) {
		return TagList{}, false
	}
	val := c.Data[e.ID]
	c.em.releaseEntity(e)
	return val, true
}
