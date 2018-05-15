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

func (c *HitboxComponent) SafeSet(id int, val [2]uint16) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}
