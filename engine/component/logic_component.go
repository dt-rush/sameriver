package component

import (
    "sync"
    "github.com/dt-rush/donkeys-qquest/engine"
)

type LogicUnit struct {
    Name string
    Logic func (int)
}



type LogicComponent struct {
    data map[int](LogicUnit)
    write_mutex sync.Mutex
}

func (c *LogicComponent) Init (capacity int, game *engine.Game) {
    c.data = make (map[int](LogicUnit), capacity)
}

func (c *LogicComponent) Set (id int, val interface{}) {
    c.write_mutex.Lock()
    val_ := val.(LogicUnit)
    c.data [id] = val_
    c.write_mutex.Unlock()
}

func (c *LogicComponent) Get (id int) LogicUnit {
    return c.data [id]
}

func (c *LogicComponent) DefaultValue() interface{} {
    return LogicUnit {"empty function", func (int) {}}
}

func (c *LogicComponent) String() string {
    return "logic component print implementation is TODO"
}

func (c *LogicComponent) Name() string {
    return "LogicComponent"
}
