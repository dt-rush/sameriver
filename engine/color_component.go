/**
  *
  * Component for the color of an entity
  *
  *
**/

package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"sync"
	"sync/atomic"
	"time"
)

type ColorComponent struct {
	Data [MAX_ENTITIES]sdl.Color
	// component-wide mutex is write-locked by any system which operates on
	// this component in bulk, read-locked by calls to SafeGet()
	mutex sync.RWMutex
	// entityLocks is a pointer to the entityLocks array in the EntityManager
	// which holds these components through ComponentsTable
	entityLocks *[MAX_ENTITIES]uint32
}

// get the value for an entity from the component in a safe manner
func (c *ColorComponent) SafeGet(id uint16) sdl.Color {
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
