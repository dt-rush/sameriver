/**
  *
  * Component for the velocity of an entity
  *
  *
**/

package engine

import (
	"sync"
)

type VelocityComponent struct {
	Data       [MAX_ENTITIES][2]int16
	Mutex sync.Mutex
}

func (c *VelocityComponent) SafeSet(id uint16, val [2]int16) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}
