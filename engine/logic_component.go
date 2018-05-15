/**
  *
  * Component for a piece of running logic attached to an entity
  *
  *
**/

package engine

import (
	"sync"
)

type LogicUnit struct {
	Logic func(dt uint16)
	Name  string
}

type LogicComponent struct {
	Data       [MAX_ENTITIES]LogicUnit
	Mutex sync.Mutex
}

func (c *LogicComponent) SafeSet(id uint16, val LogicUnit) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}
