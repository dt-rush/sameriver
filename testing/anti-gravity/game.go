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
	g := &Game{
		NewWorld(),
		NewControls(),
		NewUI(r, f),
		r, f}
	g.ui.UpdateMsg(MODENAMES[g.c.mode])
	return g
}

func (g *Game) HandleKeyEvents(e sdl.Event) {
	switch e.(type) {
	case *sdl.KeyboardEvent:
		ke := e.(*sdl.KeyboardEvent)
		if ke.Type == sdl.KEYDOWN {
			if ke.Keysym.Sym == sdl.K_m {
				g.c.ToggleMode()
				g.ui.UpdateMsg(MODENAMES[g.c.mode])
			}
			if ke.Keysym.Sym == sdl.K_c {
				g.w.ClearObstacles()
			}
		}
	}
}

func (g *Game) HandleMouseEvents(me *sdl.MouseButtonEvent) {
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
		e := Entity{pos: pos}
		g.w.e = &e
	}
	if button == sdl.BUTTON_RIGHT {
		if g.w.e != nil {
			g.w.e.moveTarget = &pos
		}
	}
}

func (g *Game) HandleObstacleInput(button uint8, pos Point2D) {
	if button == sdl.BUTTON_LEFT {
		g.w.obstacles = append(g.w.obstacles, pos)
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
		case *sdl.MouseButtonEvent:
			g.HandleMouseEvents(e.(*sdl.MouseButtonEvent))
		}
	}
	return true
}

func (g *Game) gameloop() int {

	fpsTicker := time.NewTicker(time.Millisecond * (1000 / FPS))
	stopRenderChan := make(chan bool)

	go func() {
	renderloop:
		for {
			select {
			case _ = <-stopRenderChan:
				stopRenderChan <- true
				break renderloop
			case _ = <-fpsTicker.C:
				sdl.Do(func() {
					g.r.Clear()
					g.w.mutex.Lock()
					g.DrawWorld()
					g.DrawUI()
					g.r.Present()
					g.w.mutex.Unlock()
				})
			}
		}
	}()

	moveTicker := time.NewTicker(16 * time.Millisecond)
	velTicker := time.NewTicker(500 * time.Millisecond)

gameloop:
	for {
		if !g.HandleEvents() {
			break gameloop
		}

		g.w.mutex.Lock()
		select {
		case _ = <-moveTicker.C:
			g.w.MoveEntity()
		case _ = <-velTicker.C:
			g.w.UpdateEntityVel()
		default:
		}
		g.w.mutex.Unlock()
		sdl.Delay(1000 / (2 * FPS))
	}
	stopRenderChan <- true
	<-stopRenderChan
	return 0
}
