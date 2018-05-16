/*
  *
  *
  *
  *
**/

package engine

import (
	"fmt"
	"runtime"
	"time"

	"github.com/dt-rush/go-func-profiler/func_profiler"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Game struct {
	// a channel (buffer size 1) through which appears (as if by magic)
	// the next scene
	NextSceneChan chan (Scene)
	// a channel (blocking) through which a signal can be sent (true or false
	// doesn't matter) to end the currently running scene.
	EndSceneChan chan (bool)
	// the scene to play while the next scene is running Init()
	LoadingScene Scene
	// Map of scenes by ints (constants) so scenes can identify each other
	// without import cycles
	SceneMap SceneMap
	// to allow scenes to store data somewhere that other scenes
	// can access it (TODO: currently unused - refactor?)
	GameState map[string]string
	// SDL resources to pass as references to each scene
	Window   *sdl.Window
	Renderer *sdl.Renderer
	// profiling members
	func_profiler        func_profiler.FuncProfiler
	gameloop_profiler_id uint16
}

func (g *Game) Init(WINDOW_TITLE string,
	WINDOW_WIDTH int32,
	WINDOW_HEIGHT int32) {

	// init systems
	g.InitSDL(WINDOW_TITLE, WINDOW_WIDTH, WINDOW_HEIGHT)
	// set up func profiler
	g.setupFuncProfiler()
	// initialize the scene map
	g.SceneMap.Map = make(map[int]Scene)
	// set up game state
	g.GameState = make(map[string]string)
}

func (g *Game) InitSDL(WINDOW_TITLE string,
	WINDOW_WIDTH int32,
	WINDOW_HEIGHT int32) {

	Logger.Println("[Game] Starting to init SDL")
	defer func() {
		Logger.Println("[Game] Finished init of SDL")
	}()
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
	sdl.ShowCursor(0)
	g.Window, g.Renderer = BuildWindowAndRenderer(
		WINDOW_TITLE,
		WINDOW_WIDTH,
		WINDOW_HEIGHT)
}

func (g *Game) setupFuncProfiler() {
	g.func_profiler = func_profiler.NewFuncProfiler(
		func_profiler.FUNC_PROFILER_SIMPLE)
	g.gameloop_profiler_id = g.func_profiler.RegisterFunc("gameloop")
}

func (g *Game) Destroy() {
	// free anything else that needs to be destroyed (happens once)
	// do we even need to do this, since the whole program exits
	// when the game does? anyways...
	// TODO: make sure this is actually a proper and complete destroy method
	g.Renderer.Destroy()
}

func (g *Game) End() {
	Logger.Println("[Game] in Game.End()")
	g.EndSceneChan <- true
}

func (g *Game) handleKeyboard(scene Scene) {

	// poll for events
	sdl.PumpEvents()
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			Logger.Printf("[Game] sdl.QuitEvent received: %v", t)
			// notice we use a nonblocking goroutine
			g.End()
			return
		case *sdl.KeyboardEvent:
			keyboard_event := event.(*sdl.KeyboardEvent)
			// if escape, exit immediately, else pass to the scene
			if keyboard_event.Keysym.Sym == sdl.K_ESCAPE {
				g.End()
			} else {
				scene.HandleKeyboardEvent(keyboard_event)
			}
		}
	}
	// pass keyboard state to scene
	keyboard_state := sdl.GetKeyboardState()
	scene.HandleKeyboardState(keyboard_state)
}

func (g *Game) RunLoadingScene() chan bool {
	loading_scene_stopped_signal_chan := make(chan bool)
	go func() {
		g.RunScene(g.LoadingScene)
		loading_scene_stopped_signal_chan <- true
	}()
	return loading_scene_stopped_signal_chan
}

func (g *Game) blankScreen() {
	g.Renderer.SetDrawColor(0, 0, 0, 255)
	g.Renderer.Clear()
}

func (g *Game) RunScene(scene Scene) {
	if DEBUG_GOROUTINES {
		Logger.Printf("[Game] Before running %s, NumGoroutine = %d",
			scene.Name(),
			runtime.NumGoroutine())
	}

	scene.StartLogic()
	g.runGameLoopOnScene(scene)
	scene.StopLogic()

	if DEBUG_GOROUTINES {
		Logger.Printf("[Game] After running %s, NumGoroutine = %d",
			scene.Name(),
			runtime.NumGoroutine())
	}
}

func (g *Game) logGameLoopStarted(scene Scene) {
	// print log message to notify scene starting
	Logger.Printf("[Game] \\\\\\  /// scene %s starting to run",
		scene.Name())
}

func (g *Game) initGameLoopProfiler(scene Scene) {
	// set profiler name
	g.func_profiler.SetName(g.gameloop_profiler_id,
		fmt.Sprintf("%s gameloop", scene.Name()))
}

func (g *Game) logGameLoopEnded(scene Scene) {
	// Scene has ended. Print a summary
	Logger.Printf("[Game] //// \\\\\\\\ %s stopped running.",
		scene.Name())
	Logger.Print(g.func_profiler.GetSummaryString(g.gameloop_profiler_id))
}

func (g *Game) clearGameLoopProfiler() {
	// clear the timer for the new scene to start its timing
	g.func_profiler.ClearData(g.gameloop_profiler_id)
}

func (g *Game) destroyScene(scene Scene) {
	// destroy resources used by the scene
	// (but don't trash the laoding scene which is reused)
	if scene != g.LoadingScene {
		Logger.Printf("[Game] Destroying resources used by %s", scene.Name())
		go scene.Destroy()
	}
}

func (g *Game) runGameLoopOnScene(scene Scene) {

	g.logGameLoopStarted(scene)
	g.initGameLoopProfiler(scene)
	defer g.logGameLoopEnded(scene)
	defer g.clearGameLoopProfiler()
	defer g.destroyScene(scene)

	// Actual gameloop code:
	fps_timer := NewPeriodicTimer(uint16(1000 / FPS))
	t0 := time.Now().UnixNano()
	// gameloop
	for {
		select {
		// if someone requested we end the scene, oblige them
		case <-g.EndSceneChan:
			break
		// else, run an iteration of the game loop
		default:
			// start the profiling timer for the gameloop
			g.func_profiler.StartTimer(g.gameloop_profiler_id)
			// calculate loop dt
			t1 := time.Now().UnixNano()
			dt_ms := uint16((t1 - t0) / 1e6)
			// draw active scene at framerate
			if fps_timer.Tick(dt_ms) {
				sdl.Do(func() {
					g.blankScreen()
					scene.Draw(g.Window, g.Renderer)
					g.Renderer.Present()
				})
			}
			// handle events and update scene
			sdl.Do(func() {
				g.handleKeyboard(scene)
			})
			scene.Update(dt_ms)
			// end the profiling timer
			g.func_profiler.EndTimer(g.gameloop_profiler_id)
			// set t0 so we can calculate dt next loop iteration
			t0 = t1
		}
		// everyone deserves some rest now and then
		sdl.Delay(16)
	}
}

func (g *Game) Run() {
	for {
		select {
		case scene := <-g.NextSceneChan:
			loadingSceneStoppedChan := g.RunLoadingScene()
			scene.Init(g)
			g.EndSceneChan <- true // end the loading scene
			<-loadingSceneStoppedChan
			g.RunScene(scene)
		default:
			Logger.Println("[Game] Last scene finished with no next scene. Game ending.")
		}
	}
}
