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
	Data       [MAX_ENTITIES]bool
	WriteMutex sync.Mutex
}

func (c *ActiveComponent) SafeSet(id int, val bool) {
	c.WriteMutex.Lock()
	c.Data[id] = val
	c.WriteMutex.Unlock()
}
