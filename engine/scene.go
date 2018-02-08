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
    Init (game *Game) chan bool
    Run()
    Update (dt_ms int)
    Draw (window *sdl.Window, renderer *sdl.Renderer)
    HandleKeyboardState (keyboard_state []uint8)
    IsRunning() bool
    Stop()
    Name() string
    Destroy()
}


