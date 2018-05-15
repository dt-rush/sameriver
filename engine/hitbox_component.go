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
	Data       [MAX_ENTITIES][2]uint16
	Mutex sync.Mutex
}

func (c *HitboxComponent) SafeSet(id uint16, val [2]uint16) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}

func (c *HitboxComponent) SafeGet(id uint16) [2]uint16 {
	c.Mutex.Lock()
	val := c.Data[id]
	c.Mutex.Unlock()
	return val
}
