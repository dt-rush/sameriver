package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"time"
)

type Game struct {
	w  *World
	c  *Controls
	ui *UI

	r *sdl.Renderer
	f *ttf.Font
}

func NewGame(r *sdl.Renderer, f *ttf.Font) *Game {
	g := &Game{r: r, f: f}

	g.w = NewWorld(g)
	g.c = NewControls()
	g.ui = NewUI(r, f)

	g.ui.UpdateMsg(0, MODENAMES[g.c.mode])
	return g
}

func (g *Game) HandleKeyEvents(e sdl.Event) {
	switch e.(type) {
	case *sdl.KeyboardEvent:
		ke := e.(*sdl.KeyboardEvent)
		if ke.Type == sdl.KEYDOWN {
			if ke.Keysym.Sym == sdl.K_m {
				g.c.ToggleMode()
				g.ui.UpdateMsg(0, MODENAMES[g.c.mode])
			}
			if ke.Keysym.Sym == sdl.K_c {
				g.w.ClearObstacles()
			}
		}
	}
}

func (g *Game) HandleMouseMotionEvents(me *sdl.MouseMotionEvent) {
	p := MouseMotionEventToPoint2D(me)
	if me.State&sdl.ButtonLMask() != 0 {
		if g.c.mode == MODE_PLACING_WAYPOINT {
			g.HandleWayPointInput(sdl.BUTTON_LEFT, p)
		}
	}
	if me.State&sdl.ButtonRMask() != 0 {
		if g.c.mode == MODE_PLACING_WAYPOINT {
			g.HandleWayPointInput(sdl.BUTTON_RIGHT, p)
		} else if g.c.mode == MODE_PLACING_OBSTACLE {
			g.HandleObstacleInput(sdl.BUTTON_LEFT, p)
		}
	}
}

func (g *Game) HandleMouseButtonEvents(me *sdl.MouseButtonEvent) {
	p := MouseButtonEventToPoint2D(me)
	if me.Type == sdl.MOUSEBUTTONDOWN {
		if g.c.mode == MODE_PLACING_WAYPOINT {
			g.HandleWayPointInput(me.Button, p)
		} else if g.c.mode == MODE_PLACING_OBSTACLE {
			g.HandleObstacleInput(me.Button, p)
		}
	}
}

func (g *Game) HandleWayPointInput(button uint8, pos Point2D) {
	if button == sdl.BUTTON_LEFT {
		g.w.e = NewEntity(pos, g.w)
	}
	if button == sdl.BUTTON_RIGHT {
		if g.w.e != nil {
			g.w.e.moveTarget = &pos
		}
	}
}

func (g *Game) HandleObstacleInput(button uint8, pos Point2D) {
	if button == sdl.BUTTON_LEFT {
		g.w.obstacles = append(g.w.obstacles, CenteredSquare(pos, OBSTACLESZ))
	}
	if button == sdl.BUTTON_RIGHT {

	}
}

func (g *Game) HandleQuit(e sdl.Event) bool {
	switch e.(type) {
	case *sdl.QuitEvent:
		return false
	case *sdl.KeyboardEvent:
		ke := e.(*sdl.KeyboardEvent)
		if ke.Keysym.Sym == sdl.K_ESCAPE ||
			ke.Keysym.Sym == sdl.K_q {
			return true
		}
	}
	return false
}

func (g *Game) HandleEvents() bool {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		if g.HandleQuit(e) {
			return false
		}
		switch e.(type) {
		case *sdl.KeyboardEvent:
			g.HandleKeyEvents(e)
		case *sdl.MouseMotionEvent:
			g.HandleMouseMotionEvents(e.(*sdl.MouseMotionEvent))
		case *sdl.MouseButtonEvent:
			g.HandleMouseButtonEvents(e.(*sdl.MouseButtonEvent))
		}
	}
	return true
}

func (g *Game) gameloop() int {

	fpsTicker := time.NewTicker(time.Millisecond * (1000 / FPS))

gameloop:
	for {
		// try to draw
		select {
		case _ = <-fpsTicker.C:
			sdl.Do(func() {
				g.r.Clear()
				g.DrawWorld()
				g.DrawUI()
				g.r.Present()
			})
		default:
		}

		// handle input
		if !g.HandleEvents() {
			break gameloop
		}

		// update world
		g.w.Update()

		sdl.Delay(1000 / (2 * FPS))
	}
	return 0
}
