/**
  *
  * Component for an entity's taglist
  *
  *
**/

package engine

import (
	"errors"
	"fmt"
)

type TagListComponent struct {
	Data [MAX_ENTITIES]TagList
	em   *EntityManager
}

func (c *TagListComponent) SafeGet(e EntityToken) (TagList, error) {
	if !c.em.lockEntity(e) {
		return TagList{}, errors.New(fmt.Sprintf("%v no longer exists", e))
	}
	returnCopy := c.Data[e.ID].Copy()
	c.em.releaseEntity(e)
	return returnCopy, nil
}
