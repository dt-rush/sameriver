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
	Data       [MAX_ENTITIES][2]float32
	Mutex sync.Mutex
}

func (c *VelocityComponent) SafeSet(id uint16, val [2]float32) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}

func (c *VelocityComponent) SafeGet(id uint16) [2]float32 {
	c.Mutex.Lock()
	val := c.Data[id]
	c.Mutex.Unlock()
	return val
}
