package components

import (
    "github.com/dt-rush/donkeys-qquest/engine"
)

type LogicUnit struct {
    Name string
    Logic func (float64)
}



type LogicComponent struct {
    data map[int](LogicUnit)
}

func (c *LogicComponent) Init (capacity int, game *engine.Game) {
    c.data = make (map[int](LogicUnit), capacity)
}

func (c *LogicComponent) Set (id int, val interface{}) {
    val_ := val.(LogicUnit)
    c.data [id] = val_
}

func (c *LogicComponent) Get (id int) LogicUnit {
    return c.data [id]
}

func (c *LogicComponent) DefaultValue() interface{} {
    return LogicUnit {"empty function", func (float64) {}}
}

func (c *LogicComponent) String() string {
    return "logic component print implementation is TODO" 
}

func (c *LogicComponent) Name() string {
    return "LogicComponent"
}
