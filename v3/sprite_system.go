package sameriver

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type SpriteSystem struct {
	w              *World
	SpriteEntities *UpdatedEntityList

	Textures   map[string]*sdl.Texture
	NilTexture *sdl.Texture
}

func NewSpriteSystem(renderer *sdl.Renderer) *SpriteSystem {
	s := &SpriteSystem{}
	s.Textures = make(map[string]*sdl.Texture, 256)
	s.LoadFiles(renderer)
	s.generateNilTexture(renderer)
	return s
}

func (s *SpriteSystem) GetSprite(name string) Sprite {
	texture, ok := s.Textures[name]
	if !ok {
		texture = s.NilTexture
	}
	return Sprite{
		texture,       // texture
		0,             // frame
		true,          // visible
		sdl.FLIP_NONE, // flip
	}
}

func (s *SpriteSystem) generateNilTexture(renderer *sdl.Renderer) {
	surface, err := sdl.CreateRGBSurface(
		0,          // flags
		8,          // width
		8,          //height
		int32(32),  // depth
		0xff000000, // rgba masks
		0x00ff0000,
		0x0000ff00,
		0x000000ff)
	if err != nil {
		panic(err)
	}
	rect := sdl.Rect{0, 0, 8, 8}
	color := uint32(0x9fddbcff) // feijoa
	surface.FillRect(&rect, color)
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	s.NilTexture = texture
}

func (s *SpriteSystem) LoadFiles(renderer *sdl.Renderer) {
	files, err := ioutil.ReadDir("assets/images/sprites")
	if err != nil {
		Logger.Println(err)
		logWarning("could not open assets/images/sprites; skipping SpriteSystem.LoadFiles()")
		return
	}
	for _, f := range files {
		var err error
		log_err := func(err error) {
			Logger.Printf("[Sprite manager] failed to load %s", f.Name())
			panic(err)
		}
		// get image, convert to texture, and store
		// image to texture
		surface, err := img.Load(fmt.Sprintf("assets/images/sprites/%s", f.Name()))
		if err != nil {
			log_err(err)
			continue
		}
		mapkey := strings.Split(f.Name(), ".png")[0]
		s.Textures[mapkey], err = renderer.CreateTextureFromSurface(surface)
		if err != nil {
			log_err(err)
			continue
		}
		surface.Free()
	}
}

// System funcs

func (s *SpriteSystem) GetComponentDeps() []string {
	return []string{"Sprite,Sprite"}
}

func (s *SpriteSystem) LinkWorld(w *World) {
	s.w = w

	s.SpriteEntities = w.GetUpdatedEntityListByComponentNames([]string{"Sprite"})
}

func (s *SpriteSystem) Update(dt_ms float64) {
	// nil?
}

func (s *SpriteSystem) Expand(n int) {
	// nil?
}
