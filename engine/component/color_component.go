/**
  *
  * Component for the color of an entity
  *
  *
**/

package component

import (
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type ColorComponent struct {
	Data  [MAX_ENTITIES]sdl.Color
	Mutex sync.RWMutex
}

func (c *ColorComponent) SafeSet(id uint16, val sdl.Color) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[id] = val
}

func (c *ColorComponent) SafeGet(id uint16) sdl.Color {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	val := c.Data[id]
	return val
}
