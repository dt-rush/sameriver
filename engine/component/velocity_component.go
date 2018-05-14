/**
  *
  * Component for the velocity of an entity
  *
  *
**/

package component

import (
	"fmt"
	"sync"

	"github.com/dt-rush/donkeys-qquest/engine"
)

type VelocityComponent struct {
	Data       [engine.MAX_ENTITIES][2]uint16
	WriteMutex sync.Mutex
}

func (c *VelocityComponent) SafeSet(id int, val [2]uint16) {
	c.WriteMutex.Lock()
	c.Data[id] = val
	c.WriteMutex.Unlock()
}
