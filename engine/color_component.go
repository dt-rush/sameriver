/**
  *
  * Component for the color of an entity
  *
  *
**/

package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type ColorComponent struct {
	Data       [MAX_ENTITIES]sdl.Color
	Mutex sync.Mutex
}

func (c *ColorComponent) SafeSet(id uint16, val sdl.Color) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}

func (c *ColorComponent) SafeGet(id uint16) sdl.Color {
	c.Mutex.Lock()
	val := c.Data[id]
	c.Mutex.Unlock()
	return val
}
