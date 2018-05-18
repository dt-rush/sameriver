/**
  *
  * Component for the hitbox of an entity
  *
  *
**/

package engine

import (
	"sync"
)

type HitboxComponent struct {
	Data  [MAX_ENTITIES][2]uint16
	Mutex sync.RWMutex
}

func (c *HitboxComponent) SafeSet(id uint16, val [2]uint16) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[id] = val
}

func (c *HitboxComponent) SafeGet(id uint16) [2]uint16 {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	val := c.Data[id]
	return val
}
