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
	Data  [MAX_ENTITIES][2]float32
	Mutex sync.RWMutex
}

func (c *VelocityComponent) SafeSet(id uint16, val [2]float32) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[id] = val
}

func (c *VelocityComponent) SafeGet(id uint16) [2]float32 {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	val := c.Data[id]
	return val
}
