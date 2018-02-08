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



type ActiveComponent struct {

    data map[int]bool

}

func (c *ActiveComponent) Init (capacity int, game *engine.Game) {
    c.data = make (map[int]bool, capacity)
}



// connected to gamescene.go update():
// TODO factor out the "get all active components with hitbox and position"
// ...
func (c *ActiveComponent) Has (id int) bool {
    _, ok := c.data[id]
    return ok
}

func (c *ActiveComponent) Get (id int) bool {
    return c.data [id]
}

func (c *ActiveComponent) Set (id int, val interface{}) {
    val_ := val.(bool)
    c.data[id] = val_
}

func (c *ActiveComponent) DefaultValue () interface{} {
    return false
}

func (c *ActiveComponent) String() string {
    return fmt.Sprintf ("%v", c.data)
}

func (c *ActiveComponent) Name() string {
    return "ActiveComponent"
}



