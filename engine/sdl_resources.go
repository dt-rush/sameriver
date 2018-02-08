package engine


import (
    "github.com/veandco/go-sdl2/sdl"

    "github.com/dt-rush/donkeys-qquest/utils"
)


/*
 * Builds and returns an SDL window and renderer object
 * for the game to use
 */
func BuildWindowAndRenderer (window_title string, width int32, height int32) (*sdl.Window, *sdl.Renderer) {

    window, err := sdl.CreateWindow (window_title,
        sdl.WINDOWPOS_UNDEFINED,
        sdl.WINDOWPOS_UNDEFINED,
        0, 0,
        sdl.WINDOW_SHOWN | sdl.WINDOW_FULLSCREEN_DESKTOP)
    if err != nil {
        panic(err)
    }

    renderer, err := sdl.CreateRenderer (window,
        -1,
        sdl.RENDERER_SOFTWARE)
    if err != nil {
        panic (err)
    }

    window_w, window_h := window.GetSize() 
    utils.DebugPrintf ("window.GetSize() (w x h): %d x %d\n",
                        window_w, window_h)
    // renderer.SetLogicalSize (width, height)
    sdl.SetHint (sdl.HINT_RENDER_SCALE_QUALITY, "linear")
    scale_w := float32 (window_w) / float32 (width)
    scale_h := float32 (window_h) / float32 (height)
    renderer.SetScale (scale_w, scale_h)

    return window, renderer
}
