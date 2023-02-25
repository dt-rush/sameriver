// Package sameriver is a game engine, ya underdig?
package sameriver

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Window     *sdl.Window
	Renderer   *sdl.Renderer
	WindowSpec WindowSpec
	Screen     GameScreen

	running      bool
	loadingScene Scene
	currentScene Scene
	endScene     chan bool
}

type GameInitSpec struct {
	WindowSpec   WindowSpec
	LoadingScene Scene
	FirstScene   Scene
}

func RunGame(spec GameInitSpec) {
	MainMediaThread(func() {
		InitMediaLayer()
		g := &Game{
			WindowSpec: spec.WindowSpec,
			Screen: GameScreen{
				W: spec.WindowSpec.Width,
				H: spec.WindowSpec.Height,
			},
			loadingScene: spec.LoadingScene,
			currentScene: spec.FirstScene,
			endScene:     make(chan bool),
		}
		g.Window, g.Renderer = CreateWindowAndRenderer(spec.WindowSpec)
		g.run()
	})
}

func (g *Game) SetLoadingScene(scene Scene) {
	g.loadingScene = scene
}

func (g *Game) run() {
	g.running = true
	stopLoading := make(chan (bool))
	for g.running {
		if g.currentScene == nil {
			Logger.Println("next scene is nil, ending game")
			break
		}
		go func() {
			g.loadingScene.Init(g, nil)
			g.RunScene(g.loadingScene, stopLoading)
			stopLoading <- true
		}()
		g.currentScene.Init(g, nil)
		stopLoading <- true
		<-stopLoading
		g.currentScene = g.RunScene(g.currentScene, g.endScene)
	}
	g.Destroy()
}

func (g *Game) RunScene(scene Scene, endScene chan bool) Scene {
	Logger.Printf("started: %s ▷", scene.Name())
	fpsTicker := time.NewTicker(FRAME_DURATION)
	lastUpdate := time.Now()
	overrun_ms := 0.0
gameloop:
	for {
		loopStart := time.Now()
		// break the game loop when the end game loop channel gets a signal
		select {
		case <-endScene:
			break gameloop
		default:
			if scene.IsDone() {
				break gameloop
			}
			sdl.Do(func() {
				g.handleKeyboard(scene)
			})
			dt_ms := float64(time.Since(lastUpdate).Nanoseconds()) / 1e6
			// if we overran last loop, we get proportionally less time this loop
			// (this keeps frame-rate steady while we try to run scene.Update() as
			// often as possible)
			allowance_ms := float64(FRAME_DURATION_INT)
			if overrun_ms > 0 {
				allowance_ms -= overrun_ms
			}
			scene.Update(dt_ms, allowance_ms)
			lastUpdate = time.Now()
			select {
			case <-fpsTicker.C:
				sdl.Do(func() {
					g.blankScreen()
					scene.Draw(g.Window, g.Renderer)
					g.Renderer.Present()
				})
			default:
			}
		}
		// (World.runtimeSharer shares an allowance of engine.FRAME_DURATION
		// among all {systems,world,entities} logics. This is the max it can
		// share per Update(). If it goes over, elapsed will be > FRAME_DURATION
		// and we'll skip any sleeping. Note that scene.Draw() only occurs
		// every fpsTicker tick anyway, so
		elapsed_ms := float64(time.Since(loopStart) / 1e6)
		overrun_ms = elapsed_ms - FRAME_DURATION_INT
		if overrun_ms < 0 {
			time.Sleep(time.Duration(-overrun_ms * float64(time.Millisecond)))
		}
	}
	// once gameloop ends, get next scene and destroy scene if transient
	nextScene := scene.NextScene()
	scene.End()
	if scene.IsTransient() {
		Logger.Printf("destroying scene: %s\n", scene.Name())
		go scene.Destroy()
	}
	Logger.Printf("ended: %s ■", scene.Name())
	return nextScene
}

func (g *Game) blankScreen() {
	g.Renderer.SetDrawColor(0, 0, 0, 255)
	g.Renderer.Clear()
}

func (g *Game) handleKeyboard(scene Scene) {
	// poll for events
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			Logger.Printf("sdl.QuitEvent received: %v", t)
			// notice we use a nonblocking goroutine
			g.GoEndGame()
			return
		case *sdl.KeyboardEvent:
			keyboard_event := event.(*sdl.KeyboardEvent)
			// if escape, exit immediately, else pass to the scene
			if keyboard_event.Keysym.Sym == sdl.K_ESCAPE {
				g.GoEndGame()
				return
			} else {
				scene.HandleKeyboardEvent(keyboard_event)
			}
		}
	}
	// pass keyboard state to scene
	keyboard_state := sdl.GetKeyboardState()
	scene.HandleKeyboardState(keyboard_state)
}

func (g *Game) GoEndGame() {
	Logger.Println("in Game.End()")
	if g.running {
		go func() {
			g.running = false
			g.endScene <- true
		}()
	}
}

func (g *Game) Destroy() {
	// TODO: make sure this is actually a proper and complete destroy method
	g.Renderer.Destroy()
	g.Window.Destroy()
}
