/*
 *  You are now dead, I'm sorry.
 *  You managed to catch %d donkeys!
 *
 *  Play again? [y/n]
 */

package scenes

import (

    "github.com/dt-rush/donkeys-qquest/engine"
    "github.com/dt-rush/donkeys-qquest/constants"

    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
    // "github.com/veandco/go-sdl2/gfx"
)

type GameOverScene struct {
    // whether the scene is running
    running bool
    // used to make destroy() idempotent
    destroyed bool
    // the game
    game *engine.Game

    // needed to write strings to graphics
    font *ttf.Font
    // texture used to write multiple text textures to. saved statically
    // and drawn to the screen
    screen_texture *sdl.Texture
}

func (s *GameOverScene) Init (game *engine.Game) chan bool {

    s.game = game
    init_done_signal_chan := make (chan bool)

    go func () {
        var err error

        if s.font, err = ttf.OpenFont("assets/test.ttf", 16); err != nil {
            panic(err)
        }

        // create a texture to render to 
        s.screen_texture, err = s.game.Renderer.CreateTexture (
            sdl.PIXELFORMAT_RGBA8888,
            sdl.TEXTUREACCESS_TARGET,
            constants.WINDOW_WIDTH,
            constants.WINDOW_HEIGHT)
        if err != nil { panic (err) }

        // set the renderer's texture to screen_texture
        s.game.Renderer.SetRenderTarget (s.screen_texture)

        // write a black background to the screen texture 
        s.game.Renderer.SetDrawColor (0, 0, 0, 255)
        s.game.Renderer.Clear()

        // render the "game over" message
        s.render_message_to_texture (
            "GAME OVER",
            sdl.Color{255,255,255,255},
            &sdl.Rect{
                constants.WINDOW_WIDTH / 8,
                (constants.WINDOW_HEIGHT * 3) / 8,
                (constants.WINDOW_WIDTH * 6) / 8,
                constants.WINDOW_HEIGHT / 8})

        // render the player score message
        s.render_message_to_texture (
            // TODO determine how to pass info between scenes
            // fmt.Sprintf ("You managed to capture %d donkeys", score),
            "Play again? y/n",
            sdl.Color{230,230,230,255},
            &sdl.Rect{
                constants.WINDOW_WIDTH / 8,
                (constants.WINDOW_HEIGHT * 5) / 8,
                (constants.WINDOW_WIDTH * 6) / 8,
                constants.WINDOW_HEIGHT / 12})

        // write the "replay?" message
        // TODO: implement dialogue-selection struct to enable the left/right
        // arrow selection of a continue choice, and to draw the rectangle
        // using sdl_gfx ThickLineRGBA

        // restore original render target
        s.game.Renderer.SetRenderTarget (nil)

        init_done_signal_chan <- true

    }()

    return init_done_signal_chan
}

func (s *GameOverScene) Run() {
    s.running = true
}

func (s *GameOverScene) Name() string {
    return "game over scene"
}

func (s *GameOverScene) Stop() {
    s.running = false
}

func (s *GameOverScene) IsRunning () bool {
    return s.running
}

func (s *GameOverScene) Destroy () {
    if ! s.destroyed {
        s.destroyed = true
        s.screen_texture.Destroy()
        s.font.Close()
    }
}

func (s *GameOverScene) Update (dt_ms int) {
    // TODO: update menu selection element (it will
    // use dt to blink the selection)
}

func (s *GameOverScene) Draw (window *sdl.Window, renderer *sdl.Renderer) {
    renderer.Copy (
        s.screen_texture,
        nil,
        &sdl.Rect{
            0,
            0,
            constants.WINDOW_WIDTH,
            constants.WINDOW_HEIGHT})
}

func (s *GameOverScene) HandleKeyboardState (keyboard_state []uint8) {
    // null implementation
}

func (s *GameOverScene) HandleKeyboardEvent (keyboard_event *sdl.KeyboardEvent) {
    // TODO: use left / right arrow to pass state-change to
    // menu selection element. For now, we just take the y or n keys
    switch keyboard_event.Keysym.Sym {
    case sdl.K_y:
        game_scene := GameScene{}
        s.game.NextScene = &game_scene
        s.Stop()
    case sdl.K_n:
        s.game.NextScene = nil
        s.Stop()
    }
}

func (s *GameOverScene) render_message_to_texture (
    msg string,
    color sdl.Color,
    dst *sdl.Rect) {

    // surface & texture to be used in writing the message to the texture
    var surface *sdl.Surface
    var texture *sdl.Texture
    var err error
        surface, err = s.font.RenderUTF8Solid (msg, color)
    if err != nil { panic (err) }
    texture, err = s.game.Renderer.CreateTextureFromSurface (surface)
    if err != nil { panic (err) }
    // this copies our texture to the *target* texture (screen_texture)
    s.game.Renderer.SetDrawColor (
        color.R,
        color.G,
        color.B,
        color.A)
    s.game.Renderer.Copy (texture, nil, dst)
    // free the resources allocated above
    surface.Free()
    texture.Destroy()
}
