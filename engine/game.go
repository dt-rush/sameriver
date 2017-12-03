/**
  * 
  * 
  * 
  * 
**/



package engine

import (
	"fmt"
	"time"

	"github.com/dt-rush/donkeys-qquest/utils"
	
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/veandco/go-sdl2/mix"
)

type Game struct {
	scene Scene
	next_scene_chan chan Scene
	scene_end_sig_chan chan bool
	window *sdl.Window
	renderer *sdl.Renderer
	running bool
	loading_scene Scene
	func_profiler utils.FuncProfiler
}

func (g *Game) Init (WINDOW_TITLE string,
	WINDOW_WIDTH int32,
	WINDOW_HEIGHT int32) {

	// init systems
	g.InitSystems()

	// set state
	g.running = true
	g.next_scene_chan = make (chan Scene, 1)
	g.scene_end_sig_chan = make (chan bool)
	g.func_profiler = utils.FuncProfiler{}
	g.window, g.renderer = BuildWindowAndRenderer (
		WINDOW_TITLE,
		WINDOW_WIDTH,
		WINDOW_HEIGHT)
}

func (g *Game) InitSystems() {

	var err error
	
	// init SDL
	sdl.Init (sdl.INIT_EVERYTHING)
	
	// init SDL TTF
	err = ttf.Init()
	if err != nil {
		panic (err) 
	}
	
	// init SDL Audio
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

func (g *Game) IsRunning() bool {
	return g.running
}

func (g *Game) Destroy() {
	// free anything else that needs to be destroyed (happens once)
	// do we even need to do this, since the whole program exits
	// when the game does? anyways...
	g.renderer.Destroy()
}

func (g *Game) EndScene() {
	g.scene.Stop()
}

func (g *Game) NextSceneChan() chan Scene {
	return g.next_scene_chan
}

func (g *Game) End() {
	fmt.Println ("in game.end()")
	g.scene.Stop()
	fmt.Println ("g.scene.stop() finished")
	<- g.scene_end_sig_chan // wait for scene to end (TODO investigate, can this be used for saves?)
	fmt.Println ("got scene_end_sig_chan")
	g.running = false
	g.next_scene_chan <- nil
	fmt.Println ("g.destroy()")
	g.Destroy()
}

func (g *Game) HandleEvents () {
	// if sdl.QuitEvent occurred, exit immediately
	sdl.PumpEvents()
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			fmt.Printf ("sdl.QuitEvent received: %v\n", t)
			go g.End()
			return
		}
	}
	// if escape, exit immediately
	keyboard_state := sdl.GetKeyboardState()
	if keyboard_state [sdl.SCANCODE_ESCAPE] == 1 {
		fmt.Printf ("sdl.SCANCODE_ESCAPE was 1 in keyboard_state array\n")
		go g.End()
		return
	}
	g.scene.HandleKeyboardState (keyboard_state)
}


func (g *Game) RunLoadingScene () {
	g.RunScene (g.loading_scene)
}

// game loop logic
func (g *Game) RunScene (scene Scene) {
	// stop existing scene
	old_scene := g.scene
	if old_scene != nil {
		old_scene.Stop()
		// wait for scene end sig sent from outside of
		// old_scene's game loop goroutine soon to terminate its looping
		<- g.scene_end_sig_chan
	}
	// we're officially okay to enter another game loop!
	// set the scene to run
	fmt.Printf ("///  \\\\\\ %s starting to run\n", scene.Name())
	scene.Run()
	g.scene = scene
	
	// start a game loop on the game's scene
	ticker := time.NewTicker (16 * time.Millisecond)
	t0 := <- ticker.C
	accum := 0
	gameloop_counter := 0
	gameloop_ms_accum := 0.0
	// loop
	for g.scene.IsRunning() {
		// profiling wrapper TODO find a cleaner way to do this
		gameloop_counter++
		gameloop_ms_accum += g.func_profiler.Time (func () {
			// update ticker, calculate loop dt
			t1 := <- ticker.C
			dt_ms := float64 (t1.UnixNano() - t0.UnixNano()) / 1e6
			
			// draw active scene at framerate (60 fps)
			accum += int (dt_ms)
			if accum > 1000 / 60 {
				// kill any backlog
				for accum > 1000 / 60 {
					accum = accum % (1000 / 60)
				}
				sdl.Do (func () {
					g.renderer.Present()
					g.renderer.Clear()
					g.scene.Draw (g.window, g.renderer)
				})
			}

			// sdl sleep 16 ms
			sdl.Delay (16)

			// handle events
			sdl.Do (g.HandleEvents)

			// update active scene
			g.scene.Update (dt_ms)

			
			// shift t0 to t1 for next loop
			t0 = t1
			
		})
	}
	fmt.Printf ("\\\\\\\\//// %s stopped running [gameloop_ms_avg = %.3f]\n\n", scene.Name(), float64 (gameloop_ms_accum) / float64 (gameloop_counter))
	if scene != g.loading_scene {
	   	go scene.Destroy()
	}
	g.scene_end_sig_chan <- true
}


func (g *Game) CreateTextureFromSurface (surface *sdl.Surface) (*sdl.Texture, error) {
	return g.renderer.CreateTextureFromSurface (surface)
}



func (g *Game) SetLoadingScene (scene Scene) {
	g.loading_scene = scene
	g.loading_scene.Init (g)
}

func (g *Game) PushScene (scene Scene) {
	go func () {
		g.next_scene_chan <- scene
	}()
}
