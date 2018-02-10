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
    "github.com/dt-rush/donkeys-qquest/engine/utils"

    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
)

const COLOR_CHANGE_MS = 50

type MenuScene struct {

    // whether the scene is running
    running bool
    // used to make destroy() idempotent
    destroyed bool
    // the game
    game *engine.Game

    // TODO figure out why go-sdl2 is failing when
    // closing two fonts from the same file
    title_font *ttf.Font
    small_font *ttf.Font

    // for prerendering the below
    rainbow_colors []sdl.Color

    // rainbow of title texts, prerendered
    rainbow_surfaces []*sdl.Surface
    // textures of the above, for Renderer.Copy() in draw()
    rainbow_textures []*sdl.Texture
    // for iterating thrugh the above
    rainbow_index int
    // for accumulating dt's to manage color shifting in time
    accum_rainbow engine.TimeAccumulator
    // message = "press space"
    message_surface *sdl.Surface
    // texture of the above, for Renderer.Copy() in draw()
    message_texture *sdl.Texture
    // for general timing, resets at 5000 ms to 0 ms
    accum_5000 engine.TimeAccumulator
}






func (s *MenuScene) Init (game *engine.Game) chan bool {

    s.game = game
    init_done_signal_chan := make (chan bool)
    s.destroyed = false

    go func () {

        var err error
        // used to render text to textures
        var surface *sdl.Surface

        // load fonts

        // TODO figure out why go-sdl2 is failing when
        // closing two fonts from the same file
        s.title_font, err = ttf.OpenFont("assets/test.ttf", 16)
        if err != nil { panic(err) }
        s.small_font, err = ttf.OpenFont("assets/test.ttf", 10)
        if err != nil { panic(err) }

        // create rainbow of title surfaces
        s.rainbow_colors = []sdl.Color{sdl.Color{255, 128, 237, 19},
            sdl.Color{255, 176, 205, 3},
            sdl.Color{255, 217, 160, 6},
            sdl.Color{255, 245, 112, 28},
            sdl.Color{255, 255, 65, 65},
            sdl.Color{255, 245, 28, 112},
            sdl.Color{255, 217, 6, 160},
            sdl.Color{255, 176, 3, 205},
            sdl.Color{255, 128, 19, 237},
            sdl.Color{255, 80, 51, 253},
            sdl.Color{255, 39, 96, 250},
            sdl.Color{255, 11, 144, 228},
            sdl.Color{255, 1, 191, 191},
            sdl.Color{255, 11, 228, 144},
            sdl.Color{255, 39, 250, 96},
            sdl.Color{255, 80, 253, 51}}

        s.rainbow_textures = make ([]*sdl.Texture, len (s.rainbow_colors))

        // render rainbow of titles

        // iterate the rainbow colors
        // and prerender the text at the given rainbow color
        for i := 0; i < len (s.rainbow_colors); i++ {

            color := s.rainbow_colors [i]

            // create the surface
            surface, err = s.title_font.RenderUTF8Solid (
                "Donkeys QQuest", color)
            if err != nil { panic (err) }
            // create the texture
            s.rainbow_textures[i], err = s.game.Renderer.CreateTextureFromSurface (surface)
            if err != nil { panic (err) }
            // free the surface we used
            surface.Free()
        }

        s.rainbow_index = 0

        s.accum_rainbow = engine.CreateTimeAccumulator (COLOR_CHANGE_MS)
        s.accum_5000 = engine.CreateTimeAccumulator (5000)

        // render message ("press space") surface
        surface, err = s.small_font.RenderUTF8Solid ("Press Space",
            sdl.Color{255, 255, 255, 255})
        if err != nil { panic (err) }
        // create the texture
        s.message_texture, err = s.game.Renderer.CreateTextureFromSurface (surface)
        if err != nil { panic (err) }
        // free the surface we used
        surface.Free()

        init_done_signal_chan <- true
    }()

    return init_done_signal_chan

}

func (s *MenuScene) Stop() {
    s.running = false
}

func (s *MenuScene) IsRunning () bool {
    return s.running
}





func (s *MenuScene) Update (dt_ms int) {

    if s.accum_rainbow.Tick (dt_ms) {
        s.rainbow_index++
        s.rainbow_index %= len (s.rainbow_colors)
    }

    s.accum_5000.Tick (dt_ms)

}





func (s *MenuScene) Draw (window *sdl.Window, renderer *sdl.Renderer) {

    // write title
    title_dst := sdl.Rect{
        constants.WINDOW_WIDTH / 8,
        (constants.WINDOW_HEIGHT * 3) / 8,
        (constants.WINDOW_WIDTH * 6) / 8,
        constants.WINDOW_HEIGHT / 8}
    renderer.Copy (s.rainbow_textures [s.rainbow_index], nil, &title_dst)

    // write message ("press space")
    x_offset := constants.WINDOW_WIDTH / 3
    msg_dst := sdl.Rect{
        x_offset,
        int32 (
            float64 (constants.WINDOW_HEIGHT * 2 / 5) +
            float64 (constants.WINDOW_HEIGHT / 10) * math.Sin (5 * 2 * math.Pi * s.accum_5000.Completion())),
        constants.WINDOW_WIDTH - x_offset * 2,
        20}
    renderer.Copy (s.message_texture, nil, &msg_dst)
}

func (s *MenuScene) transition () {
    game_scene := GameScene{}
    s.game.NextScene = &game_scene
    s.Stop()
}

func (s *MenuScene) HandleKeyboardState (keyboard_state []uint8) {
    k := keyboard_state
    // if space, transition (push game scene)
    if k [sdl.SCANCODE_SPACE] == 1 {
        s.transition()
    }
}

func (s *MenuScene) HandleKeyboardEvent (keyboard_event *sdl.KeyboardEvent) {
    // null implementation
}

func (s *MenuScene) Destroy() {

    utils.DebugPrintln ("MenuScene.Destroy() called")

    if ! s.destroyed {

        utils.DebugPrintln ("inside MenuScene.Destroy(), ! s.destroyed")

        // TODO figure out why go-sdl2 is failing when
        // closing two fonts from the same file
        // [note, in line with the below, this might not even be the root cause]
        // s.title_font.Close()


        // issue was still occurring even when using one font:

        /*

goroutine 73 [syscall, locked to thread]:
runtime.cgocall(0x4e9f70, 0xc8200aa730, 0x100000000000000)
    /usr/lib/go-1.6/src/runtime/cgocall.go:123 +0x11b fp=0xc8200aa6f8 sp=0xc8200aa6c8
github.com/veandco/go-sdl2/sdl._Cfunc_SDL_FreeSurface(0x7f0f3c01b250)
    ??:0 +0x36 fp=0xc8200aa730 sp=0xc8200aa6f8
github.com/veandco/go-sdl2/sdl.(*Surface).Free(0x7f0f3c01b250)
    /home/anon/gocode/src/github.com/veandco/go-sdl2/sdl/surface.go:80 +0x75 fp=0xc8200aa780 sp=0xc8200aa730
main.(*MenuScene).destroy(0xc82007e000)
    /home/anon/gocode/src/github.com/dt-rush/donkeys-qquest/main/menuscene.go:214 +0x8e fp=0xc8200aa7b8 sp=0xc8200aa780
runtime.goexit()
    /usr/lib/go-1.6/src/runtime/asm_amd64.s:1998 +0x1 fp=0xc8200aa7c0 sp=0xc8200aa7b8
created by main.(*MenuScene).stop
    /home/anon/gocode/src/github.com/dt-rush/donkeys-qquest/main/menuscene.go:120 +0x3a


*/
        // * possibly related to https://github.com/veandco/go-sdl2/issues/187
        //
        // * changing the below to use sdl.Do in the meantime (above issue
        //   ended with discussion of sdl.CallQueue, which is no longer
        //   directly exported, but interfaced-with via func Do()
        //
        // * gameloop is a goroutine per scene, which then calls functions
        //   such as update, draw, and handle_keyboard_state which could
        //   cause function calls to end the scene / game which will call
        //   the destroy() method of the scene outside the main OS thread?
        //
        // * since it only manifested intermittently in the past (race condition
        //   fits this), the true test is over time, to see whether it ever occurs
        //   after this change
        //
        //   -- (27 of May, 2017)
        utils.DebugPrintln ("about to call sdl.Do in menuscene.destroy()")
        sdl.Do (func () {
            utils.DebugPrintln ("in sdl.Do() func for menuscene.destroy()")
            s.small_font.Close()
            for i, _ := range s.rainbow_surfaces {
                s.rainbow_surfaces [i].Free()
            }
            s.message_surface.Free()
        })
        utils.DebugPrintln ("finished sdl.Do() call in menu_scene.Destroy() which called .Free() on some surfaces")
    }
    s.destroyed = true
}



func (s *MenuScene) Run () {

    // any scene-specific routines can be spawned in here

    s.running = true

}



func (s *MenuScene) Name () string {
    return "menu scene"
}
