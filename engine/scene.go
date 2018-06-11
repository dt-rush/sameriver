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
	Name() string

	Init(game *Game, config map[string]string, stopChannel chan bool)
	StartLogic()
	StopLogic()

	Update(dt_ms int64)
	Draw(window *sdl.Window, renderer *sdl.Renderer)
	HandleKeyboardState(keyboard_state []uint8)
	HandleKeyboardEvent(keyboard_event *sdl.KeyboardEvent)

	IsTransient() bool
	Destroy()
}
