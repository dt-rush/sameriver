package sameriver

import (
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func MainMediaThread(f func()) {
	sdl.Main(f)
}

func InitMediaLayer() {
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
}

func CreateWindowAndRenderer(spec WindowSpec) (*sdl.Window, *sdl.Renderer) {
	// create the window
	flags := uint32(sdl.WINDOW_SHOWN)
	if spec.Fullscreen {
		flags |= sdl.WINDOW_FULLSCREEN_DESKTOP
	}
	window, err := sdl.CreateWindow(spec.Title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(spec.Width), int32(spec.Height),
		flags)
	if err != nil {
		panic(err)
	}

	// create the renderer
	renderer, err := sdl.CreateRenderer(window,
		-1,
		sdl.RENDERER_SOFTWARE)
	if err != nil {
		panic(err)
	}

	// set renderer scale
	if spec.Fullscreen {
		window_w, window_h := window.GetSize()
		Logger.Printf("window.GetSize() (w x h): %d x %d",
			window_w, window_h)
		sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")
		scale_w := float32(window_w) / float32(spec.Width)
		scale_h := float32(window_h) / float32(spec.Height)
		renderer.SetScale(scale_w, scale_h)
	}

	// set renderer alpha
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	return window, renderer
}
