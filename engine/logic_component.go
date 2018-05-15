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
	Logic func(dt int)
	Name  string
}

type LogicComponent struct {
	Data       [MAX_ENTITIES]LogicUnit
	Mutex sync.Mutex
}

func (c *LogicUnit) SafeSet(id int, val LogicUnit) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}
