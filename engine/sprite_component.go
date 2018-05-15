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
	Data       [MAX_ENTITIES]Sprite
	Mutex sync.Mutex
}

func (c *SpriteComponent) SafeSet(id uint16, val Sprite) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}
