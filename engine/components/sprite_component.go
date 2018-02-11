/**
  *
  *
  *
  *
**/



package components

import (
    "fmt"

    "github.com/dt-rush/donkeys-qquest/engine"

    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/img"
)

type SpriteComponent struct {
    // entity ID -> sprite index
    data map[int](int)
    // sprite index -> *sdl.Texture
    sprites [](*sdl.Texture)
    // sprite index -> sdl.RendererFlip (see go-sdl2/sdl/render.go)
    flip []sdl.RendererFlip
    // for getting the index into the *sdl.Texture array via the original filename
    name_index_map map[string]int
}

func (c *SpriteComponent) IndexOf (s string) int {
    return c.name_index_map [s]
}

func (c *SpriteComponent) Get (id int) *sdl.Texture {
    return c.sprites [c.data [id]]
}

func (c *SpriteComponent) Set (id int, val interface{}) {
    val_ := val.(int)
    c.data [id] = val_
}


// TODO - separate "flip" into its own component... or just leave it in here,
// as a sort of side-car of get/set beside the main one
func (c *SpriteComponent) SetFlip (id int, val sdl.RendererFlip) {
    c.flip [c.data [id]] = val
}

func (c *SpriteComponent) GetFlip (id int) sdl.RendererFlip {
    return c.flip [c.data [id]]
}





func (c *SpriteComponent) DefaultValue () interface{} {
    return 0
}

func (c *SpriteComponent) String() string {
    return fmt.Sprintf ("%v", c.data)
}

func (c *SpriteComponent) Name() string {
    return "SpriteComponent"
}

// connected to gamescene.go draw():
// TODO refactor .has() to be recorded by a component using
// ...
func (c *SpriteComponent) Has (id int) bool {
    _, ok := c.data[id]
    return ok
}

func (c *SpriteComponent) Init (capacity int, game *engine.Game) {

    // init data storage
    c.data = make (map[int]int, capacity)
    c.sprites = make ([](*sdl.Texture), capacity)
    c.name_index_map = make (map[string]int, capacity)
    c.flip = make ([]sdl.RendererFlip, capacity)

    // set all c.flip default values just to be sure
    // (what _is_ the default value of a "const iota"
    // fake-enum? ... zero? Since it's just a type alias for int?
    // and what sdl constant would that be?
    for i, _ := range c.flip {
        c.flip [i] = sdl.FLIP_NONE
    }

    // image file enum for now, dynamic load later (TODO)
    const (
        FLAME = "flame.png"
        FLAME2 = "flame2.png"
        FLAME3 = "flame3.png"
        DONKEY1 = "donkey.png"
    )
    to_load := []string{FLAME, FLAME2, FLAME3, DONKEY1}
    for i, s := range to_load {
        var err error
        log_err := func (err error) {
            engine.Logger.Printf ("failed to load %s", s)
            engine.Logger.Printf ("%v", err)
        }
        // add s->i to name_index_map
        c.name_index_map [s] = i
        // get image, convert to texture, and store
        // image to texture
        surface, err := img.Load (fmt.Sprintf ("assets/%s", s))
        if err != nil {
            log_err (err)
            continue
        }
        c.sprites [i], err = game.Renderer.CreateTextureFromSurface (surface)
        if err != nil {
            log_err (err)
            continue
        }
        surface.Free()
    }

}

func (c *SpriteComponent) Destroy() {
    for _, sprite := range (c.sprites) {
        sprite.Destroy()
    }
}



