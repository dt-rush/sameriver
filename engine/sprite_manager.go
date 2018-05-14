package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"strings"
)

type SpriteManager struct {
	SpriteComponent *component.SpriteComponent
	Textures        map[string]*sdl.Texture
}

func (c *SpriteManager) Init(
	sprite_component *component.SpriteComponent,
	renderer *sdl.Renderer) {

	c.sprite_component = sprite_component
	c.Textures = make(map[string]*sdl.Texture, 256)
	c.LoadFiles()
}

func (c *SpriteManager) NewSprite(name string) component.Sprite {
	return component.Sprite{
		c.Textures[name], // texture
		0,                // frame
		true,             // visible
		sdl.FLIP_NONE,    // flip
	}
}

func (c *SpriteManager) LoadFiles() {
	files, err := ioutil.ReadDir("assets/images/sprites")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		var err error
		log_err := func(err error) {
			engine.Logger.Printf("failed to load %s", f.Name())
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
		c.Textures[mapkey], err = renderer.CreateTextureFromSurface(surface)
		if err != nil {
			log_err(err)
			continue
		}
		surface.Free()
	}
}