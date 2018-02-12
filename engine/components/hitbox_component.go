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
)

// stores a (length, width) ordered pair
type HitboxComponent struct {
    data map[int]([2]float64)
}

func (c *HitboxComponent) Init (capacity int, game *engine.Game) {
    // init data storage
    c.data = make (map[int]([2]float64), capacity)
}

// connected to gamescene.go:
// TODO factor out the "get all active components 
// with hitbox and position" into a tag
func (c *HitboxComponent) Has (id int) bool {
    _, ok := c.data [id]
    return ok
}

func (c *HitboxComponent) Set (id int, val interface{}) {
    val_ := val.([2]float64)
    c.data[id] = val_
}

func (c *HitboxComponent) Get (id int) [2]float64 {
    return c.data [id]
}

func (c *HitboxComponent) DefaultValue () interface{} {
    r := [2]float64{0, 0}
    return r
}

func (c *HitboxComponent) String() string {
    return fmt.Sprintf ("%v", c.data)
}

func (c *HitboxComponent) Name() string {
    return "HitboxComponent"
}

// TODO implement
// func (c *AudioComponent) destroy() {
    
// }



