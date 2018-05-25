/**
  *
  * Component for a piece of running logic attached to an entity
  *
  *
**/

package engine

import (
	"sync"
	"sync/atomic"
	"time"
)

// Each LogicFunc will started as a goroutine, supplied with the EntityToken
// of the entity it's attached to, a channel on which a stop signal may
// arrive, and a reference to the EntityManager
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
type EntityLogicFunc func(
	entity EntityToken,
	StopChannel chan bool,
	em *EntityManager)

type LogicUnit struct {
	f           EntityLogicFunc
	Name        string
	StopChannel chan bool
}

// Create a new LogicUnit instance
func NewLogicUnit(Name string, f EntityLogicFunc) LogicUnit {
	return LogicUnit{
		f,
		Name,
		make(chan bool)}
}

type LogicComponent struct {
	Data [MAX_ENTITIES]LogicUnit
	// component-wide mutex is write-locked by any system which operates on
	// this component in bulk, read-locked by calls to SafeGet()
	mutex sync.RWMutex
	// entityLocks is a pointer to the entityLocks array in the EntityManager
	// which holds these components through ComponentsTable
	entityLocks *[MAX_ENTITIES]uint32
}

// get the value for an entity from the component in a safe manner
func (c *LogicComponent) SafeGet(id uint16) LogicUnit {
	// NOTE: holding the entity lock and the component mutex is not really
	// that bad since the total duration is the time it takes to read a single
	// value

	// wait for the entity to not be held for modification
	for !atomic.CompareAndSwapUint32(&c.entityLocks[id], 0, 1) {
		time.Sleep(FRAME_SLEEP)
	}
	// read-lock the component
	c.mutex.RLock()
	val := c.Data[id]
	// release the mutex and entity lock
	c.mutex.RUnlock()
	atomic.StoreUint32(&c.entityLocks[id], 0)
	return val
}
