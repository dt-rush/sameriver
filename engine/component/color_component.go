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
	Data       [MAX_ENTITIES]sdl.Color
	WriteMutex sync.Mutex
}

func (c *ColorComponent) SafeSet(id int, val sdl.Color) {
	c.WriteMutex.Lock()
	c.Data[id] = val
	c.WriteMutex.Unlock()
}
