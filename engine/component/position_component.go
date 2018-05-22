/**
  *
  * Component for the position of an entity
  *
  *
**/

package component

import (
	"sync"
)

type PositionComponent struct {
	Data  [MAX_ENTITIES][2]int16
	Mutex sync.RWMutex
}

func (c *PositionComponent) SafeSet(id uint16, val [2]int16) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[id] = val
}

func (c *PositionComponent) SafeGet(id uint16) [2]int16 {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	val := c.Data[id]
	return val
}
