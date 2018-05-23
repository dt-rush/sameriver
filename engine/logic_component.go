/**
  *
  * Component for a piece of running logic attached to an entity
  *
  *
**/

package engine

import (
	"sync"
)

// Each LogicFunc will started as a goroutine, supplied with the ID of the
// entity it's attached to, a channel on which a stop signal may arrive, and
// a reference to the EntityManager
//
// Through the EntityManager, the goroutine will be able to:
//
// - request an UpdatedEntityList with an arbitrary EntityQuery
// - get entity component data via SafeGet() on the component in em.Components
// - send EntitySpawnRequest messages to the SpawnChannel
// - send EntityStateModification messages to the StateSetChannel
// - send EntityComponentModification messages to the ComponentSetChannel
// - atomically modify an entity using em.AtomicEntityModify()
//
// The StopChannel is not buffered, since we need to be sure when the
// logic has ended.
type LogicFunc func(
	entityID uint16,
	StopChannel chan bool,
	em *EntityManager)

type LogicUnit struct {
	LogicFunc   LogicFunc
	Name        string
	StopChannel chan bool
}

// Create a new LogicUnit instance
func NewLogicUnit(LogicFunc LogicFunc, Name string) LogicUnit {
	return LogicUnit{
		LogicFunc,
		Name,
		make(chan bool)}
}

type LogicComponent struct {
	Data  [MAX_ENTITIES]LogicUnit
	Mutex sync.RWMutex
}

func (c *LogicComponent) SafeSet(id uint16, val LogicUnit) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[id] = val
}

func (c *LogicComponent) SafeGet(id uint16) LogicUnit {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	val := c.Data[id]
	return val
}
