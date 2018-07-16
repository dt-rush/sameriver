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

	Init(game *Game, config map[string]string)

	Update(dt_ms float64)
	Draw(window *sdl.Window, renderer *sdl.Renderer)
	HandleKeyboardState(keyboard_state []uint8)
	HandleKeyboardEvent(keyboard_event *sdl.KeyboardEvent)

	IsDone() bool
	NextScene() Scene
	End()
	IsTransient() bool
	Destroy()
}
