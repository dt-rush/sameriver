/**
  *
  *
  *
  *
**/



package component



import (
    "sync"
    "fmt"
    "github.com/dt-rush/donkeys-qquest/engine"
    "github.com/veandco/go-sdl2/sdl"
)



type ColorComponent struct {
    data map[int]sdl.Color
    write_mutex sync.Mutex
}

func (c *ColorComponent) Init (capacity int, game *engine.Game) {
    c.data = make (map[int]sdl.Color, capacity)
}

func (c *ColorComponent) Get (id int) sdl.Color {
    return c.data[id]
}

func (c *ColorComponent) Set (id int, val interface{}) {
    c.write_mutex.Lock()
    val_ := val.(sdl.Color)
    c.data[id] = val_
    c.write_mutex.Unlock()
}

func (c *ColorComponent) DefaultValue () interface{} {
    return sdl.Color{0,0,0,255}
}

func (c *ColorComponent) String() string {
    return fmt.Sprintf ("%v", c.data)
}

func (c *ColorComponent) Name() string {
    return "ColorComponent"
}



