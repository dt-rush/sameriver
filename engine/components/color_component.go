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



type ColorComponent struct {

	data map[int]uint32
	
}

func (c *ColorComponent) Init (capacity int, game *engine.Game) {
	c.data = make (map[int]uint32, capacity)
}

func (c *ColorComponent) Get (id int) uint32 {
	return c.data[id]
}

func (c *ColorComponent) Set (id int, val interface{}) {
	val_ := val.(uint32)
	c.data[id] = val_
}

func (c *ColorComponent) DefaultValue () interface{} {
	return uint32 (0xff888888)
}

func (c *ColorComponent) String() string {
	return fmt.Sprintf ("%v", c.data)
}

func (c *ColorComponent) Name() string {
	return "ColorComponent"
}



