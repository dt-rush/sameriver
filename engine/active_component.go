/**
  *
  * Component for whether an entity is active (if inactive, no system
  * should operate on its components)
  *
  *
**/

package engine

import (
	"sync"
)

type ActiveComponent struct {
	Data       [MAX_ENTITIES]bool
	Mutex sync.Mutex
}

func (c *ActiveComponent) SafeSet(id uint16, val bool) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}

func (c *ActiveComponent) SafeGet(id uint16) bool {
	c.Mutex.Lock()
	val := c.Data[id]
	c.Mutex.Unlock()
	return val
}
