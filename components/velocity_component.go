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



type VelocityComponent struct {

	data map[int]([2]float64)

}

func (c *VelocityComponent) Init (capacity int, game *engine.Game) {
	c.data = make (map[int]([2]float64), capacity)
}

// connected to gamescene.go update():
// TODO factor out the "get all active components with hitbox and position"
// ...
func (c *VelocityComponent) Has (id int) bool {
	_, ok := c.data[id]
	return ok
}

func (c *VelocityComponent) Get (id int) [2]float64 {
	return c.data [id]
}

func (c *VelocityComponent) Set (id int, val interface{}) {
	// type assert
	val_ := val.([2]float64)
	c.data[id] = val_
}

func (c *VelocityComponent) DefaultValue () interface{} {
	return [2]float64{0, 0}
}

func (c *VelocityComponent) String() string {
	return fmt.Sprintf ("%v", c.data)
}

func (c *VelocityComponent) Name() string {
	return "VelocityComponent"
}



