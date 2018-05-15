/**
  *
  * Component for the position of an entity
  *
  *
**/

package engine

import (
	"sync"
)

type PositionComponent struct {
	Data  [MAX_ENTITIES][2]int16
	Mutex sync.Mutex
}

func (c *PositionComponent) SafeSet(id uint16, val [2]int16) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}

func (c *PositionComponent) SafeGet(id uint16) [2]int16 {
	c.Mutex.Lock()
	val := c.Data[id]
	c.Mutex.Unlock()
	return val
}
