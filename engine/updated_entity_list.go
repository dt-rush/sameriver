/*
 * a list of entities which will be updated by another goroutine maybe,
 * which has a mutex that the user can lock when they wish to look at the
 * current contents. Can be sorted (needed by the data structure / algorithm
 * used in CollisionSystem)
 *
 */

package engine

import (
	"sync"
	"time"
)

// A list of entities which is can be regularly updated by one goroutine
// while another reads and uses it
type UpdatedEntityList struct {
	// the entities in the list (tagged with gen)
	Entities []uint16
	// used to protect the Entities slice when adding or removing an entity
	Mutex sync.RWMutex
	// the channel along which updates to the list will come
	EntityChannel chan int32
	// used to stop the update loop's goroutine in
	// the case that they're done with the list (by calling Stop())
	stopUpdateLoopChannel chan bool
	// the name of the list (for debugging)
	Name string
	// the ID of the list (used for memory management)
	ID uint16
	// whether the entities slice should be sorted
	Sorted bool
	// a slice of funcs who want to be called *before* the entity gets
	// added/removed (that is, before the mutex unlocks)
	callbacks []func(int32)
}

// create a new UpdatedEntityList by giving it a channel on which it will
// receive entity IDs
func NewUpdatedEntityList(
	EntityChannel chan int32, ID uint16, Name string) *UpdatedEntityList {

	l := UpdatedEntityList{}
	l.EntityChannel = EntityChannel
	l.Entities = make([]uint16, 0)
	l.stopUpdateLoopChannel = make(chan (bool))
	l.Name = Name
	l.start()
	return &l
}

// called during the creation of a list. Starts a goroutine which listens
// on the channel and either inserts or deletes entities as appropriate
func (l *UpdatedEntityList) start() {
	go func() {
	updateloop:
		for {
			select {
			case _ = <-l.stopUpdateLoopChannel:
				break updateloop
			case id := <-l.EntityChannel:
				l.actOnIDSignal(id)
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (l *UpdatedEntityList) Stop() {
	l.stopUpdateLoopChannel <- true
}

// acts on an ID signal, which is either an ID to insert or -(ID + 1) to remove
func (l *UpdatedEntityList) actOnIDSignal(idEncoded int32) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	for _, callback := range l.callbacks {
		go callback(idEncoded)
	}

	if idEncoded >= 0 {
		id := uint16(idEncoded)
		updatedEntityListDebug("%s got insert:%d", l.Name, idEncoded)
		l.insert(id)
	} else {
		id := uint16(-(idEncoded + 1))
		updatedEntityListDebug("%s got remove:%d", l.Name, id)
		l.remove(id)
	}
}

// inserts an entity into the list (private so only called by the update loop)
func (l *UpdatedEntityList) insert(id uint16) {
	if l.Sorted {
		SortedUint16SliceInsert(&l.Entities, id)
	} else {
		l.Entities = append(l.Entities, id)
	}
}

// removes an entity from the list (private so only called by the update loop)
func (l *UpdatedEntityList) remove(id uint16) {
	if l.Sorted {
		SortedUint16SliceRemove(&l.Entities, id)
	} else {
		removeUint16FromSlice(id, &l.Entities)
	}
}

// add a callback to the callbacks slice
func (l *UpdatedEntityList) addCallback(callback func(idEncoded int32)) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	l.callbacks = append(l.callbacks, callback)
}
