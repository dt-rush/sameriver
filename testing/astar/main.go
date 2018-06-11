package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func handleQuit(e sdl.Event) bool {
	switch e.(type) {
	case *sdl.QuitEvent:
		return false
	case *sdl.KeyboardEvent:
		ke := e.(*sdl.KeyboardEvent)
		if ke.Keysym.Sym == sdl.K_ESCAPE ||
			ke.Keysym.Sym == sdl.K_q {
			return false
		}
	}
	return true
}

func handleKeyEvents(w *World, e sdl.Event) {
	switch e.(type) {
	case *sdl.KeyboardEvent:
		ke := e.(*sdl.KeyboardEvent)
		if ke.Keysym.Sym == sdl.K_g && ke.Type == sdl.KEYDOWN {
			w.NewWorldMap()
			w.ComputeEntityPath()
		}
	}
}

func handleMouseEvents(w *World, e sdl.Event) {
	switch e.(type) {
	case *sdl.MouseButtonEvent:
		me := e.(*sdl.MouseButtonEvent)
		if me.Type == sdl.MOUSEBUTTONDOWN {
			pos := Position{
				int(float64(me.X) / WORLD_CELL_PIXEL_WIDTH),
				int(float64(WINDOW_HEIGHT-me.Y) / WORLD_CELL_PIXEL_HEIGHT)}
			if me.Button == sdl.BUTTON_LEFT {
				e := Entity{pos: pos}
				w.e = &e
			}
			if me.Button == sdl.BUTTON_RIGHT {
				if w.e != nil {
					w.e.moveTarget = &pos
					w.ComputeEntityPath()
				}
			}
		}
	}
}

func gameloop() int {

	w := World{GenerateWorldMap(), nil}

	sdl.Init(sdl.INIT_EVERYTHING)
	r, exitcode := GetRenderer()
	if exitcode != 0 {
		return exitcode
	}

	moveTicker := time.NewTicker(50 * time.Millisecond)

	running := true
	for running {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			running = handleQuit(e)
			handleKeyEvents(&w, e)
			handleMouseEvents(&w, e)
		}
		r.SetDrawColor(0, 0, 0, 255)
		r.Clear()

		w.DrawWorldMap(r)
		w.DrawEntityAndPath(r)

		select {
		case _ = <-moveTicker.C:
			w.MoveEntity()
		default:
		}

		r.Present()
		sdl.Delay(1000 / FPS)
	}
	fmt.Println("Done")
	return 0
}

func main() {

	var exitcode int
	sdl.Main(func() {
		runtime.LockOSThread()
		exitcode = gameloop()
	})
	os.Exit(exitcode)

}
