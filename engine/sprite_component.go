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
	texture *sdl.Texture
	frame   uint8
	visible bool
	flip    sdl.RendererFlip
}

type SpriteComponent struct {
	Data       [MAX_ENTITIES]Sprite
	Mutex sync.Mutex
}

func (c *SpriteComponent) SafeSet(id int, val Sprite) {
	c.Mutex.Lock()
	c.Data[id] = val
	c.Mutex.Unlock()
}
