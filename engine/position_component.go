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
	Data       [MAX_ENTITIES][2]uint16
	WriteMutex sync.Mutex
}

func (c *PositionComponent) SafeSet(id int, val [2]uint16) {
	c.WriteMutex.Lock()
	c.Data[id] = val
	c.WriteMutex.Unlock()
}
