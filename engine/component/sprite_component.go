/**
  *
  *
  *
  *
**/

package component

import (
	"fmt"
	"sync"

	"github.com/dt-rush/donkeys-qquest/engine"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	texture *sdl.Texture
	frame   uint8
	visible bool
	flip    sdl.RendererFlip
}

type SpriteComponent struct {
	Data       [engine.MAX_ENTITIES]Sprite
	WriteMutex sync.Mutex
}

func (c *SpriteComponent) SafeSet(id int, val Sprite) {
	c.WriteMutex.Lock()
	c.Data[id] = val
	c.WriteMutex.Unlock()
}
