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

func (c *ActiveComponent) SafeSet(id int, val bool) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}
