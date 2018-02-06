/*
  *
  *
  *
  *
**/



package engine

import (
    "time"
    "math/rand"

    "github.com/dt-rush/donkeys-qquest/utils"

    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
    "github.com/veandco/go-sdl2/mix"
)

type Game struct {

    running bool
    scene Scene
    window *sdl.Window
    renderer *sdl.Renderer

    loading_scene Scene

    NextSceneChan chan Scene
    scene_end_signal_chan chan bool

    func_profiler utils.FuncProfiler
}

func (g *Game) Init (WINDOW_TITLE string,
                        WINDOW_WIDTH int32,
                        WINDOW_HEIGHT int32) {
    // seed rand
    seed := time.Now().UTC().UnixNano()
    rand.Seed (seed)
    if VERBOSE {
        utils.DebugPrintf ("rand seeded with %d\n", seed)
    }
    // init systems
    utils.DebugPrintln ("Starting to init SDL systems")
    g.InitSystems()
    utils.DebugPrintln ("Finished init of SDL systems")

    // set state
    g.running = true
    g.NextSceneChan = make (chan Scene, 1)
    g.scene_end_signal_chan = make (chan bool)
    g.func_profiler = utils.FuncProfiler{}
    g.window, g.renderer = BuildWindowAndRenderer (
        WINDOW_TITLE,
        WINDOW_WIDTH,
        WINDOW_HEIGHT)
}

func (g *Game) InitSystems() {

    var err error

    // init SDL
    // sdl.Init (sdl.INIT_VIDEO | sdl.INIT_TIMER)
    sdl.Init (sdl.INIT_EVERYTHING)

    // init SDL TTF
    err = ttf.Init()
    if err != nil {
        panic (err)
    }

    // init SDL Audio
    if (AUDIO_ON) {
        err = sdl.Init (sdl.INIT_AUDIO)
        if err != nil {
            panic (err)
        }
        err = mix.Init (mix.INIT_MP3)
        if err != nil {
            panic (err)
        }
        err = mix.OpenAudio (22050, mix.DEFAULT_FORMAT, 2, 4096)
        if err != nil {
            panic (err)
        }
    }
}

func (g *Game) IsRunning() bool {
    return g.running
}

func (g *Game) Destroy() {
    // free anything else that needs to be destroyed (happens once)
    // do we even need to do this, since the whole program exits
    // when the game does? anyways...
    g.renderer.Destroy()
}

func (g *Game) EndCurrentScene() {
    if g.scene != nil {
        g.scene.Stop()
        <-g.scene_end_signal_chan
    }
}

func (g *Game) End() {
    utils.DebugPrintln ("in game.end()")
    g.EndCurrentScene()
    utils.DebugPrintln ("g.scene.stop() finished")
    <-g.scene_end_signal_chan
    utils.DebugPrintln ("got scene_end_signal_chan")
    g.running = false
    g.NextSceneChan <- nil
    utils.DebugPrintln ("g.destroy()")
    g.Destroy()
}

func (g *Game) handleEvents () {
    // if sdl.QuitEvent occurred, exit immediately
    sdl.PumpEvents()
    var event sdl.Event
    for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
        switch t := event.(type) {
        case *sdl.QuitEvent:
            utils.DebugPrintf ("sdl.QuitEvent received: %v\n", t)
            go g.End()
            return
        }
    }
    // if escape, exit immediately
    keyboard_state := sdl.GetKeyboardState()
    if keyboard_state [sdl.SCANCODE_ESCAPE] == 1 {
        utils.DebugPrintf ("sdl.SCANCODE_ESCAPE was 1 in keyboard_state array\n")
        go g.End()
        return
    }
    g.scene.HandleKeyboardState (keyboard_state)
}

func (g *Game) RunLoadingScene () {
    g.RunScene (g.loading_scene)
}

func (g *Game) CreateTextureFromSurface (surface *sdl.Surface) (*sdl.Texture, error) {
    return g.renderer.CreateTextureFromSurface (surface)
}

func (g *Game) SetLoadingScene (scene Scene) {
    g.loading_scene = scene
    <-g.loading_scene.Init (g)
}

func (g *Game) PushScene (scene Scene) {
    go func () {
        g.NextSceneChan <- scene
    }()
}

func (g *Game) RunScene (scene Scene) {
    utils.DebugPrintf ("wanting to run %s\n", scene.Name())
    // will block until another copy of Game.RunScene()
    // which is running for the already-running scene
    // finishes Game.runGameLoopOnScene, at which point
    // (see the code for runGameLoopOnScene())
    // we send "true" into the channel
    g.EndCurrentScene()
    utils.DebugPrintln ("g.EndCurrentScene() completed")
    // we're okay to enter another game loop now
    g.scene = scene
    scene.Run()
    g.runGameLoopOnScene (scene)
    // destroy resources used by the scene
    // (but don't trash the laoding scene which is reused)
    if scene != g.loading_scene {
        go scene.Destroy()
    }
}

func (g *Game) runGameLoopOnScene (scene Scene) {
    ticker := time.NewTicker (16 * time.Millisecond)
    t0 := <-ticker.C
    accum := 0
    gameloop_counter := 0
    gameloop_ms_accum := 0.0
    utils.DebugPrintf ("///  \\\\\\ %s starting to run\n",
                        scene.Name())
    // loop
    for scene.IsRunning() {
        // profiling wrapper TODO find a cleaner way to do this
        gameloop_counter++
        gameloop_ms_accum += g.func_profiler.Time (func () {
            // update ticker, calculate loop dt
            t1 := <-ticker.C
            dt_ms := float64 (t1.UnixNano() - t0.UnixNano()) / 1e6
            // draw active scene at framerate (60 fps)
            accum += int (dt_ms)
            if accum > 1000 / 60 {
                // eat any backlog, just draw the current frame
                for accum > 1000 / 60 {
                    accum = accum % (1000 / 60)
                }
                sdl.Do (func () {
                    g.renderer.Clear()
                    scene.Draw (g.window, g.renderer)
                    g.renderer.Present()
                })
            }
            sdl.Delay (16)
            sdl.Do (g.handleEvents)
            scene.Update (dt_ms)
            t0 = t1
        })
    }
    utils.DebugPrintln ("\\\\\\\\//// %s stopped running.")
    utils.DebugPrintf ("[gameloop_ms_avg = %.3f]\n\n",
                        scene.Name(),
                        float64 (gameloop_ms_accum) /
                            float64 (gameloop_counter))
    g.scene_end_signal_chan <- true
}


