/**
  *
  *
  *
  *
**/

package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Scene interface {
	Init(game *Game) chan bool
	Run()
	Update(dt_ms uint16)
	Draw(window *sdl.Window, renderer *sdl.Renderer)
	HandleKeyboardState(keyboard_state []uint8)
	HandleKeyboardEvent(keyboard_event *sdl.KeyboardEvent)
	IsRunning() bool
	Stop()
	Name() string
	Destroy()
}
