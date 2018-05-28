/*
 * Component used by an entity to share info across
 * its various behavior functions
 *
 */

package engine

import (
	"errors"
	"fmt"
)

// allows arbitrary storage of data

type Mind map[string]interface{}

func NewMind() Mind {
	return make(map[string]interface{})
}

type MindComponent struct {
	Data [MAX_ENTITIES]Mind
	em   *EntityManager
}

func (c *MindComponent) SafeGet(e EntityToken) (Mind, error) {
	if !c.em.lockEntity(e) {
		return Mind{}, errors.New(fmt.Sprintf("%v no longer exists", e))
	}
	returnCopy := NewMind()
	for k, v := range c.Data[e.ID] {
		returnCopy[k] = v
	}
	c.em.releaseEntity(e)
	return returnCopy, nil
}
