/*
  *
  *
  *
  *
**/

package engine

import (
	"time"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Game struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer

	currentScene                Scene
	currentSceneEndGameLoopChan chan bool

	loadingScene                Scene
	loadingSceneEndGameLoopChan chan bool
}

func NewGame() *Game {
	return &Game{
		currentSceneEndGameLoopChan: make(chan bool),
		loadingSceneEndGameLoopChan: make(chan bool),
	}
}

func (g *Game) Init(windowSpec WindowSpec, firstScene Scene) {
	g.currentScene = firstScene
	g.InitSDL(windowSpec)
}

func (g *Game) InitSDL(windowSpec WindowSpec) {
	Logger.Println("Starting to init SDL")
	defer func() {
		Logger.Println("Finished init of SDL")
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
	g.Window, g.Renderer = BuildWindowAndRenderer(windowSpec)
}

func (g *Game) Destroy() {
	// free anything else that needs to be destroyed (happens once)
	// do we even need to do this, since the whole program exits
	// when the game does? anyways...
	// TODO: make sure this is actually a proper and complete destroy method
	g.Renderer.Destroy()
}

func (g *Game) AsyncEnd() {
	Logger.Println("in Game.End()")
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
			Logger.Printf("sdl.QuitEvent received: %v", t)
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

func (g *Game) AsyncRunLoadingScene() chan bool {
	g.loadingScene.Init(g, nil)
	loading_scene_stopped_signal_chan := make(chan (bool))
	go func() {
		g.RunScene(g.loadingScene, g.loadingSceneEndGameLoopChan)
		loading_scene_stopped_signal_chan <- true
	}()
	return loading_scene_stopped_signal_chan
}

func (g *Game) RunScene(scene Scene, endGameLoopChan chan bool) Scene {
	Logger.Printf("started: %s ▷", scene.Name())
	defer Logger.Printf("ended: %s ■", scene.Name())
	fpsTicker := time.NewTicker(FRAME_SLEEP)
	lastUpdate := time.Now()
gameloop:
	for {
		loopStart := time.Now()
		// break the game loop when the end game loop channel gets a signal
		select {
		case _ = <-endGameLoopChan:
			break gameloop
		default:
			if scene.IsDone() {
				break gameloop
			}
			sdl.Do(func() {
				g.handleKeyboard(scene)
			})
			scene.Update(float64(time.Since(lastUpdate).Nanoseconds()) / 1e6)
			lastUpdate = time.Now()
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
	// once gameloop ends, get next scene and destroy scene if transient
	nextScene := scene.NextScene()
	if scene.IsTransient() {
		Logger.Printf("destroying scene: %s\n", scene.Name())
		go scene.Destroy()
	}
	return nextScene
}

func (g *Game) Run() {
	for g.currentScene != nil {
		loadingSceneStoppedChan := g.AsyncRunLoadingScene()
		g.currentScene.Init(g, nil)
		g.loadingSceneEndGameLoopChan <- true
		<-loadingSceneStoppedChan
		g.currentScene = g.RunScene(g.currentScene, g.currentSceneEndGameLoopChan)
	}
}
