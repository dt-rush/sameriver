/**
  *
  * Component for whether an entity is active (if inactive, no system
  * should operate on its components)
  *
  *
**/

package component

import (
	"sync"
)

type ActiveComponent struct {
	Data  [MAX_ENTITIES]bool
	Mutex sync.RWMutex
}

func (c *ActiveComponent) SafeSet(id uint16, val bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[id] = val
}

func (c *ActiveComponent) SafeGet(id uint16) bool {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	val := c.Data[id]
	return val
}
