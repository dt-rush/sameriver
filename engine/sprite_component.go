/**
  *
  *
  *
  *
**/

package engine

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	Texture *sdl.Texture
	Frame   uint8
	Visible bool
	Flip    sdl.RendererFlip
}

type SpriteComponent struct {
	Data  [MAX_ENTITIES]Sprite
	Mutex sync.RWMutex
}

func (c *SpriteComponent) SafeSet(id uint16, val Sprite) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[id] = val
}

func (c *SpriteComponent) SafeGet(id uint16) Sprite {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	val := c.Data[id]
	return val
}
