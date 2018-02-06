/**
  *
  *
  *
  *
**/



package scenes

import (
    "math"

    "github.com/dt-rush/donkeys-qquest/engine"
    "github.com/dt-rush/donkeys-qquest/constants"
    "github.com/dt-rush/donkeys-qquest/utils"

    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
)

type LoadingScene struct {

    // TODO separate
    // Scene "abstract class members"

    // whether the scene is running
    running bool
    // used to make destroy() idempotent
    destroyed bool
    // used to prevent double-initialization
    initialized bool
    // the game
    game *engine.Game

    // TODO preserve
    // data specific to this scene

    message_font *ttf.Font

    message_surface *sdl.Surface // message = "loading"
    message_texture *sdl.Texture // texture of the above, for Renderer.Copy() in draw()

    five_second_dt_accum float64 // for general timing, resets at 5000 ms to 0 ms
}






func (s *LoadingScene) Init (game *engine.Game) chan bool {

    s.game = game
    init_done_signal_chan := make (chan bool)

    go func () {
        if ! s.initialized {
            s.destroyed = false
            var err error
            // load font
            if s.message_font, err = ttf.OpenFont ("./assets/test.ttf", 8); err != nil {
                panic(err)
            }
            s.five_second_dt_accum = 0
            // render message ("press space") surface
            s.message_surface, err = s.message_font.RenderUTF8Solid ("Loading",
                sdl.Color{255, 255, 255, 255})
            if err != nil {
                panic (err)
            }
            // create the texture
            s.message_texture, err = s.game.CreateTextureFromSurface (s.message_surface)
            if err != nil {
                panic (err)
            }
            s.initialized = true
        }
        init_done_signal_chan <- true
    }()
    return init_done_signal_chan
}

func (s *LoadingScene) Stop () {
    s.running = false
}

func (s *LoadingScene) IsRunning () bool {
    return s.running
}





func (s *LoadingScene) Update (dt_ms float64) {

    s.five_second_dt_accum += dt_ms
    for s.five_second_dt_accum > 5000 {
        s.five_second_dt_accum -= 5000
    }
}





func (s *LoadingScene) Draw (window *sdl.Window, renderer *sdl.Renderer) {

    windowRect := sdl.Rect{0,
        0,
        constants.WINDOW_WIDTH,
        constants.WINDOW_HEIGHT}

    renderer.SetDrawColor (0, 0, 0, 255)
    renderer.FillRect (&windowRect)

    // write message ("loading")
    dst := sdl.Rect{40,
        int32 (135 + 20 * math.Sin (5 * 2 * math.Pi * s.five_second_dt_accum / 5000.0)),
        constants.WINDOW_WIDTH - 80,
        24}
    renderer.Copy (s.message_texture, nil, &dst)
}



func (s *LoadingScene) HandleKeyboardState (keyboard_state []uint8) {}



func (s *LoadingScene) Destroy() {
    utils.DebugPrintln ("loadingscene.destroy()")
    if ! s.destroyed {

        s.message_font.Close()
        s.message_surface.Free()

    }
    s.destroyed = true
}



func (s *LoadingScene) Run () {

    // any scene-specific routines can be spawned in here

    s.running = true

}



func (s *LoadingScene) Name () string {
    return "loading scene"
}
