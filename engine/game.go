/*
  *
  *
  *
  *
**/

package engine

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Game struct {

	// the scene which the game will be running presently
	Scene Scene
	// the next scene (exported field, set by scenes
	// before they stop themselves)
	NextScene Scene
	// the scene to play while the next scene is running Init()
	LoadingScene Scene
	// Map of scenes by ints (constants) so scenes can identify each other
	// without import cycles
	SceneMap SceneMap
	// to allow scenes to store data somewhere that other scenes
	// can access it (TODO: currently unused - refactor?)
	GameState map[string]string

	Window    *sdl.Window
	Renderer  *sdl.Renderer
	accum_fps TimeAccumulator

	func_profiler        FuncProfiler
	gameloop_profiler_id uint16
}

func (g *Game) Init(WINDOW_TITLE string,
	WINDOW_WIDTH int32,
	WINDOW_HEIGHT int32,
	FPS int) {
	// seed rand
	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)
	Logger.Printf("rand seeded with %d", seed)

	// init systems
	Logger.Println("Starting to init SDL systems")
	g.InitSystems()
	Logger.Println("Finished init of SDL systems")

	// set up func profiler
	g.func_profiler = NewFuncProfiler(FUNC_PROFILER_SIMPLE)
	g.gameloop_profiler_id = g.func_profiler.RegisterFunc("gameloop")

	// build the window and renderer
	g.Window, g.Renderer = BuildWindowAndRenderer(
		WINDOW_TITLE,
		WINDOW_WIDTH,
		WINDOW_HEIGHT)

	// hide the cursor
	sdl.ShowCursor(0)

	// set the FPS rate
	g.accum_fps = NewTimeAccumulator(uint16(1000 / FPS))

	g.SceneMap.Map = make(map[int]Scene)

	// set up game state
	g.GameState = make(map[string]string)
}

func (g *Game) InitSystems() {

	var err error

	// init SDL
	sdl.Init(sdl.INIT_EVERYTHING)

	// init SDL TTF
	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	// init SDL Audio
	if AUDIO_ON {
		err = sdl.Init(sdl.INIT_AUDIO)
		if err != nil {
			panic(err)
		}
		err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
		if err != nil {
			panic(err)
		}
	}
}

func (g *Game) Destroy() {
	// free anything else that needs to be destroyed (happens once)
	// do we even need to do this, since the whole program exits
	// when the game does? anyways...
	g.Renderer.Destroy()
}

func (g *Game) EndCurrentScene() {
	if g.Scene != nil {
		Logger.Printf("in Game.EndCurrentScene(), g.scene is %s",
			g.Scene.Name())
		g.Scene.Stop()
	} else {
		Logger.Println("in Game.EndCurrentScene(), g.scene is nil")
	}
}

func (g *Game) End() {
	Logger.Println("in Game.End()")
	g.Scene.Stop()
}

func (g *Game) handleKeyboard() {

	// poll for events
	sdl.PumpEvents()
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			Logger.Printf("sdl.QuitEvent received: %v", t)
			// notice we use a nonblocking goroutine
			g.End()
			return
		case *sdl.KeyboardEvent:
			keyboard_event := event.(*sdl.KeyboardEvent)
			// if escape, exit immediately, else pass to the scene
			if keyboard_event.Keysym.Sym == sdl.K_ESCAPE {
				g.End()
			} else {
				g.Scene.HandleKeyboardEvent(keyboard_event)
			}
		}
	}

	// pass keyboard state to scene
	keyboard_state := sdl.GetKeyboardState()
	g.Scene.HandleKeyboardState(keyboard_state)
}

func (g *Game) RunLoadingScene() chan bool {
	loading_scene_stopped_signal_chan := make(chan bool)
	go func() {
		g.LoadingScene.Run()
		g.runGameLoopOnScene(g.LoadingScene)
		loading_scene_stopped_signal_chan <- true
	}()
	return loading_scene_stopped_signal_chan
}

func (g *Game) blankScreen() {
	g.Renderer.SetDrawColor(0, 0, 0, 255)
	g.Renderer.Clear()
}

func (g *Game) RunScene() {
	Logger.Printf("wanting to run %s", g.Scene.Name())
	g.Scene.Run()
	Logger.Printf("about to run game loop on %s", g.Scene.Name())
	g.runGameLoopOnScene(g.Scene)
}

func (g *Game) runGameLoopOnScene(scene Scene) {
	Logger.Printf("\\\\\\  /// %s starting to run",
		scene.Name())
	// set profiler name
	g.func_profiler.SetName(g.gameloop_profiler_id,
		fmt.Sprintf("%s gameloop", scene.Name()))
	// set up ticker to get millisecond dt's to send to scene Update()
	ticker := time.NewTicker(16 * time.Millisecond)
	t0 := <-ticker.C
	// gameloop
	for scene.IsRunning() {
		// start the profiling timer for the gameloop
		g.func_profiler.StartTimer(g.gameloop_profiler_id)
		// update ticker, calculate loop dt
		t1 := <-ticker.C
		dt_ms := uint16((t1.UnixNano() - t0.UnixNano()) / 1e6)
		// draw active scene at framerate
		if g.accum_fps.Tick(dt_ms) {
			sdl.Do(func() {
				g.blankScreen()
				scene.Draw(g.Window, g.Renderer)
				g.Renderer.Present()
			})
		}
		// handle events and update scene
		sdl.Do(g.handleKeyboard)
		scene.Update(dt_ms)
		// prepare the ticker for the next go-round
		t0 = t1
		// end the profiling timer
		g.func_profiler.EndTimer(g.gameloop_profiler_id)
		// everyone deserves some rest now and then
		sdl.Delay(16)
	}

	// Scene has ended. Print a summary
	Logger.Printf("//// \\\\\\\\ %s stopped running.",
		scene.Name())
	Logger.Print(g.func_profiler.GetSummaryString(g.gameloop_profiler_id))
	// clear the timer for the new scene to start its timing
	g.func_profiler.ClearData(g.gameloop_profiler_id)

	// destroy resources used by the scene
	// (but don't trash the laoding scene which is reused)
	if scene != g.LoadingScene {
		Logger.Printf("Destroying resources used by %s", scene.Name())
		go scene.Destroy()
	}

}

func (g *Game) Run() {
	for true {
		if g.NextScene == nil {
			Logger.Println("NextScene is nil. Game ending.")
			g.End()
			break
		} else {

			loading_scene_stopped_signal_chan := g.RunLoadingScene()
			g.Scene = g.NextScene
			g.NextScene = nil
			<-g.Scene.Init(g)
			Logger.Println("<-g.Scene.Init (g)")
			g.LoadingScene.Stop()
			<-loading_scene_stopped_signal_chan
			Logger.Println("<-loading_scene_stopped_signal_chan")
			if DEBUG_GOROUTINES {
				Logger.Printf("Before running %s, NumGoroutine = %d",
					g.Scene.Name(),
					runtime.NumGoroutine())
			}
			g.RunScene()
			if DEBUG_GOROUTINES {
				Logger.Printf("After running %s, NumGoroutine = %d",
					g.Scene.Name(),
					runtime.NumGoroutine())
			}
		}
	}
}
