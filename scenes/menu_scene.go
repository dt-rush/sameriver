/**
  * 
  * 
  * 
  * 
**/



package scenes

import (
	"fmt"
	"math"

	"github.com/dt-rush/donkeys-qquest/engine"
	"github.com/dt-rush/donkeys-qquest/constants"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type MenuScene struct {

	// TODO separate
	// Scene "abstract class members"

	// whether the scene is running
	running bool
	// used to make destroy() idempotent
	destroyed bool
	// the game
	game *engine.Game

	// TODO preserve
	// data specific to this scene

	// TODO figure out why go-sdl2 is failing when
	// closing two fonts from the same file	
	// title_font *ttf.Font
	small_font *ttf.Font

	rainbow_colors []sdl.Color // for prerendering the below
	rainbow_surfaces []*sdl.Surface // rainbow of title texts, prerendered
	rainbow_textures []*sdl.Texture // textures of the above, for Renderer.Copy() in draw()
	rainbow_index int // for iterating thrugh the above
	rainbow_dt_accum float64 // for accumulating dt's to shift which rainbow_surface is displayed
	message_surface *sdl.Surface // message = "press space"
	message_texture *sdl.Texture // texture of the above, for Renderer.Copy() in draw()
	
	five_second_dt_accum float64 // for general timing, resets at 5000 ms to 0 ms
}



func (s MenuScene) CreateScene () engine.Scene {
	return engine.Scene (&s)
}



func (s *MenuScene) Init (game *engine.Game) chan bool {

	s.game = game
	init_done_sig_chan := make (chan bool)

	go func () {
		
		s.destroyed = false
		
		var err error

		// load fonts
		
		// TODO figure out why go-sdl2 is failing when
		// closing two fonts from the same file
		// if s.title_font, err = ttf.OpenFont("assets/test.ttf", 24); err != nil {
		// 	panic(err)
		// }
		
		if s.small_font, err = ttf.OpenFont("assets/test.ttf", 16); err != nil {
			panic(err)
		}

		// create rainbow of title surfaces
		s.rainbow_colors = []sdl.Color{sdl.Color{255, 128, 237, 19},
			sdl.Color{255, 176, 205, 3},
			sdl.Color{255, 217, 160, 6},
			sdl.Color{255, 245, 112, 28},
			sdl.Color{255, 255, 65, 65},
			sdl.Color{255, 245, 28, 112},
			sdl.Color{255, 217, 6, 160},
			sdl.Color{255, 176, 3, 205},
			sdl.Color{255, 128, 19, 237},
			sdl.Color{255, 80, 51, 253},
			sdl.Color{255, 39, 96, 250},
			sdl.Color{255, 11, 144, 228},
			sdl.Color{255, 1, 191, 191},
			sdl.Color{255, 11, 228, 144},
			sdl.Color{255, 39, 250, 96},
			sdl.Color{255, 80, 253, 51}}

		s.rainbow_surfaces = make ([]*sdl.Surface, len (s.rainbow_colors))
		s.rainbow_textures = make ([]*sdl.Texture, len (s.rainbow_colors))

		// render rainbow of titles

		// iterate the rainbow colors
		// and prerender the text at the given rainbow color
		for i, _ := range s.rainbow_surfaces {

			color := s.rainbow_colors [i]

			// TODO figure out why go-sdl2 is failing when
			// closing two fonts from the same file
			// s.rainbow_surfaces[i], err = s.title_font.RenderUTF8Solid ("Donkeys QQuest", color)
			// create the surface
			s.rainbow_surfaces[i], err = s.small_font.RenderUTF8Solid ("Donkeys QQuest", color)
			if err != nil {
				panic (err)
			}
			// create the texture
			s.rainbow_textures[i], err = s.game.CreateTextureFromSurface (s.rainbow_surfaces [i])
			if err != nil {
				panic (err)
			}
		}

		s.rainbow_index = 0

		s.rainbow_dt_accum = 0
		s.five_second_dt_accum = 0

		// render message ("press space") surface
		
		s.message_surface, err = s.small_font.RenderUTF8Solid ("Press Space",
			sdl.Color{255, 255, 255, 255})
		if err != nil {
			panic (err)
		}
		// create the texture
		s.message_texture, err = s.game.CreateTextureFromSurface (s.message_surface)
		if err != nil {
			panic (err)
		}
		init_done_sig_chan <- true
	}()

	return init_done_sig_chan

}

func (s *MenuScene) Stop () {
	s.running = false
}

func (s *MenuScene) IsRunning () bool {
	return s.running
}





func (s *MenuScene) Update (dt_ms float64) {
	s.rainbow_dt_accum += dt_ms
	for s.rainbow_dt_accum > 1000 {
		s.rainbow_index++
		s.rainbow_dt_accum -= 1000
		s.rainbow_index %= len (s.rainbow_colors)
	}

	s.five_second_dt_accum += dt_ms
	for s.five_second_dt_accum > 5000 {
		s.five_second_dt_accum -= 5000
	}
}





func (s *MenuScene) Draw (window *sdl.Window, renderer *sdl.Renderer) {
	// window_surface, err := window.GetSurface()
	// if err != nil {
	// 	panic(err)
	// }
	
	windowRect := sdl.Rect{0,
		0,
		constants.WINDOW_WIDTH,
		constants.WINDOW_HEIGHT}

	renderer.SetDrawColor (0, 0, 0, 255)
	renderer.FillRect (&windowRect)
	
//	window_surface.FillRect(&rect, 0xFF000000)


	// write title
	dst := sdl.Rect{20, 100, 355, 64}
	renderer.Copy (s.rainbow_textures [s.rainbow_index], nil, &dst)
	// err = s.rainbow_surfaces [s.rainbow_index].Blit (nil,
	// 	window_surface,
	// 	&sdl.Rect{65, 100, 512, 512})
	// if err != nil {
	// 	panic (err)
	// }

	// write message ("press space")
	dst = sdl.Rect{40,
		int32 (165 + 20 * math.Sin (3 * 2 * math.Pi * s.five_second_dt_accum / 5000.0)),
		constants.WINDOW_WIDTH - 80,
		24}
	renderer.Copy (s.message_texture, nil, &dst)
	// err = s.message_surface.Blit (nil, window_surface, &sdl.Rect{160, 150, 512, 512})
	// if err != nil {
	// 	panic (err)
	// }
	
	// window.UpdateSurface()
}



func (s *MenuScene) transition () {
	s.game.PushScene (GameScene{}.CreateScene())
	s.Stop()
}

func (s *MenuScene) HandleKeyboardState (keyboard_state []uint8) {
	k := keyboard_state
	// if space, transition (push game scene)
	if k [sdl.SCANCODE_SPACE] == 1 {
		s.transition()
	}
}



func (s *MenuScene) Destroy() {

	fmt.Println ("MenuScene.Destroy() called")
	
	if ! s.destroyed {
		
		fmt.Println ("inside MenuScene.Destroy(), ! s.destroyed")

		// TODO figure out why go-sdl2 is failing when
		// closing two fonts from the same file
		// [note, in line with the below, this might not even be the root cause]
		// s.title_font.Close()
		

		// issue was still occurring even when using one font:

		/*

goroutine 73 [syscall, locked to thread]:
runtime.cgocall(0x4e9f70, 0xc8200aa730, 0x100000000000000)
	/usr/lib/go-1.6/src/runtime/cgocall.go:123 +0x11b fp=0xc8200aa6f8 sp=0xc8200aa6c8
github.com/veandco/go-sdl2/sdl._Cfunc_SDL_FreeSurface(0x7f0f3c01b250)
	??:0 +0x36 fp=0xc8200aa730 sp=0xc8200aa6f8
github.com/veandco/go-sdl2/sdl.(*Surface).Free(0x7f0f3c01b250)
	/home/anon/gocode/src/github.com/veandco/go-sdl2/sdl/surface.go:80 +0x75 fp=0xc8200aa780 sp=0xc8200aa730
main.(*MenuScene).destroy(0xc82007e000)
	/home/anon/gocode/src/github.com/dt-rush/donkeys-qquest/main/menuscene.go:214 +0x8e fp=0xc8200aa7b8 sp=0xc8200aa780
runtime.goexit()
	/usr/lib/go-1.6/src/runtime/asm_amd64.s:1998 +0x1 fp=0xc8200aa7c0 sp=0xc8200aa7b8
created by main.(*MenuScene).stop
	/home/anon/gocode/src/github.com/dt-rush/donkeys-qquest/main/menuscene.go:120 +0x3a


*/
		// * possibly related to https://github.com/veandco/go-sdl2/issues/187
		//
		// * changing the below to use sdl.Do in the meantime (above issue
		//   ended with discussion of sdl.CallQueue, which is no longer
		//   directly exported, but interfaced-with via func Do()
		//
		// * gameloop is a goroutine per scene, which then calls functions
		//   such as update, draw, and handle_keyboard_state which could
		//   cause function calls to end the scene / game which will call
		//   the destroy() method of the scene outside the main OS thread?
		//
		// * since it only manifested intermittently in the past (race condition
		//   fits this), the true test is over time, to see whether it ever occurs
		//   after this change
		//
		//   -- (27 of May, 2017)
		fmt.Println ("about to call sdl.Do in menuscene.destroy()")
		sdl.Do (func () {
			fmt.Println ("in sdl.Do() func for menuscene.destroy()")
			s.small_font.Close()
			for i, _ := range s.rainbow_surfaces {
				s.rainbow_surfaces [i].Free()
			}
			s.message_surface.Free()
		})
		fmt.Println ("finished sdl.Do() call")
	}
	s.destroyed = true
}



func (s *MenuScene) Run () {

	// any scene-specific routines can be spawned in here

	s.running = true
	
}



func (s *MenuScene) Name () string {
	return "menu scene"
}
