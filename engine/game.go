/*
  *
  *
  *
  *
**/

package engine

import (
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Game struct {
	// a channel (buffer size 1) through which appears (as if by magic)
	// the next scene
	NextSceneChan chan (Scene)
	// a channel through which a signal can be sent (true or false
	// doesn't matter) to end the currently running scene.
	currentSceneEndGameLoopChan chan (bool)
	// channel to end the loading scene
	endLoadingSceneGameLoopChan chan (bool)
	// the scene to play while the next scene is running Init()
	loadingScene Scene
	// Map of scenes by ints (constants) so scenes can identify each other
	// without import cycles
	SceneMap SceneMap
	// to allow scenes to store data somewhere that other scenes
	// can access it (TODO: currently unused - refactor?)
	GameState map[string]string
	// SDL resources to pass as references to each scene
	Window   *sdl.Window
	Renderer *sdl.Renderer
}

func (g *Game) Init(WINDOW_TITLE string,
	WINDOW_WIDTH int32,
	WINDOW_HEIGHT int32) {

	// init systems
	g.InitSDL(WINDOW_TITLE, WINDOW_WIDTH, WINDOW_HEIGHT)
	// initialize the scene map
	g.SceneMap.Map = make(map[int]Scene)
	// initialized the next scene channel
	g.NextSceneChan = make(chan (Scene), 1)
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

func (g *Game) Destroy() {
	// free anything else that needs to be destroyed (happens once)
	// do we even need to do this, since the whole program exits
	// when the game does? anyways...
	// TODO: make sure this is actually a proper and complete destroy method
	g.Renderer.Destroy()
}

func (g *Game) AsyncEnd() {
	Logger.Println("[Game] in Game.End()")
	go func() {
		g.currentSceneEndGameLoopChan <- true
	}()
}

func (g *Game) SetLoadingScene(scene Scene) {
	g.loadingScene = scene
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
			g.AsyncEnd()
			return
		case *sdl.KeyboardEvent:
			keyboard_event := event.(*sdl.KeyboardEvent)
			// if escape, exit immediately, else pass to the scene
			if keyboard_event.Keysym.Sym == sdl.K_ESCAPE {
				g.AsyncEnd()
			} else {
				scene.HandleKeyboardEvent(keyboard_event)
			}
		}
	}
	// pass keyboard state to scene
	keyboard_state := sdl.GetKeyboardState()
	scene.HandleKeyboardState(keyboard_state)
}

func (g *Game) blankScreen() {
	g.Renderer.SetDrawColor(0, 0, 0, 255)
	g.Renderer.Clear()
}

func (g *Game) destroyScene(scene Scene) {
	Logger.Printf("[Game] destroying scene: %s\n", scene.Name())
	go scene.Destroy()
}

func (g *Game) logGameLoopStarted(scene Scene) {
	// print log message to notify scene starting
	Logger.Printf("[Game] started: %s ▷",
		scene.Name())
}

func (g *Game) logGameLoopEnded(scene Scene) {
	// Scene has ended. Print a summary
	Logger.Printf("[Game] ended: %s ■",
		scene.Name())
}

func (g *Game) RunScene(scene Scene, endGameLoopChan chan (bool)) {
	if DEBUG_GOROUTINES {
		Logger.Printf("[Game] Before running %s, NumGoroutine = %d",
			scene.Name(),
			runtime.NumGoroutine())
	}

	scene.StartLogic()
	g.runGameLoopOnScene(scene, endGameLoopChan)
	scene.StopLogic()
	if scene.IsTransient() {
		g.destroyScene(scene)
	}

	if DEBUG_GOROUTINES {
		Logger.Printf("[Game] After running %s, NumGoroutine = %d",
			scene.Name(),
			runtime.NumGoroutine())
	}
}

func (g *Game) AsyncRunLoadingScene() chan bool {
	g.endLoadingSceneGameLoopChan = make(chan (bool), 1)
	g.loadingScene.Init(g, nil, g.endLoadingSceneGameLoopChan)
	loading_scene_stopped_signal_chan := make(chan (bool))
	go func() {
		g.RunScene(g.loadingScene, g.endLoadingSceneGameLoopChan)
		loading_scene_stopped_signal_chan <- true
	}()
	return loading_scene_stopped_signal_chan
}

func (g *Game) runGameLoopOnScene(scene Scene, endGameLoopChan chan (bool)) {

	g.logGameLoopStarted(scene)
	defer g.logGameLoopEnded(scene)

	// Actual gameloop code:
	fpsTicker := time.NewTicker(FRAME_SLEEP)
	lastUpdate := time.Now()
	// gameloop
	for {
		loopStart := time.Now()
		// break the game loop when the end game loop channel gets a signal
		if len(endGameLoopChan) > 0 {
			Logger.Printf("[Game] len (endGameLoopChan) > 0 for %s\n",
				scene.Name())
			break
		} else {
			// else, run an iteration of the game loop
			// handle input
			sdl.Do(func() {
				g.handleKeyboard(scene)
			})
			// update scene
			scene.Update(time.Since(lastUpdate).Nanoseconds() / 1e6)
			lastUpdate = time.Now()
			// draw scene at framerate
			select {
			case _ = <-fpsTicker.C:
				sdl.Do(func() {
					g.blankScreen()
					scene.Draw(g.Window, g.Renderer)
					g.Renderer.Present()
				})
			default:
			}
		}
		// sleep if we have time ("buffer time")
		elapsed := time.Since(loopStart) / 1e6 * time.Millisecond
		if elapsed > FRAME_SLEEP {
			continue
		} else {
			sdl.Delay(uint32((FRAME_SLEEP - elapsed) / time.Millisecond))
		}
	}
}

func (g *Game) Run() {

sceneloop:
	for {
		select {
		case scene := <-g.NextSceneChan:
			Logger.Printf("[Game] wanting to run %s\n", scene.Name())
			loadingSceneStoppedChan := g.AsyncRunLoadingScene()
			g.currentSceneEndGameLoopChan = make(chan (bool), 1)
			// TODO: nil below could be a config - how do we get the config
			// for the incoming scene? Is it attached to what comes through
			// the channel?
			scene.Init(g, nil, g.currentSceneEndGameLoopChan)
			Logger.Printf("[Game] %s.Init() finished\n", scene.Name())
			g.endLoadingSceneGameLoopChan <- true
			Logger.Printf("[Game] sent signal to stop loading scene\n")
			<-loadingSceneStoppedChan
			Logger.Printf("[Game] got loading scene stopped signal\n")
			Logger.Printf("[Game] calling g.RunScene (%s)\n",
				scene.Name())
			g.RunScene(scene, g.currentSceneEndGameLoopChan)
		default:
			Logger.Println("[Game] Last scene finished with no next scene. Game ending.")
			break sceneloop
		}
	}
}
