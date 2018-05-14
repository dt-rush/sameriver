/**
  *
  * Component for the hitbox of an entity
  *
  *
**/

package engine

import (
	"fmt"
	"sync"
)

type VelocityComponent struct {
	Data       [MAX_ENTITIES][2]uint16
	WriteMutex sync.Mutex
}

func (c *VelocityComponent) SafeSet(id int, val [2]uint16) {
	c.WriteMutex.Lock()
	c.Data[id] = val
	c.WriteMutex.Unlock()
}
