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



type PositionComponent struct {

	data map[int]([2]float64)

}

func (c *PositionComponent) Init (capacity int, game *engine.Game) {
	c.data = make (map[int]([2]float64), capacity)
}


// connected to gamescene.go update():
// TODO factor out the "get all active components with hitbox and position"
// ...
func (c *PositionComponent) Has (id int) bool {
	_, ok := c.data[id]
	return ok
}

func (c *PositionComponent) Get (id int) [2]float64 {
	return c.data [id]
}

func (c *PositionComponent) Set (id int, val interface{}) {
	// type assert
	val_ := val.([2]float64)
	c.data[id] = val_
}

func (c *PositionComponent) DefaultValue () interface{} {
	return [2]float64{0, 0}
}

func (c *PositionComponent) String() string {
	return fmt.Sprintf ("%v", c.data)
}

func (c *PositionComponent) Name() string {
	return "PositionComponent"
}



