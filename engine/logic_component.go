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
	Data  [MAX_ENTITIES]LogicUnit
	Mutex sync.RWMutex
}

func (c *LogicComponent) SafeSet(id uint16, val LogicUnit) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[id] = val
}

func (c *LogicComponent) SafeGet(id uint16) LogicUnit {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	val := c.Data[id]
	return val
}
