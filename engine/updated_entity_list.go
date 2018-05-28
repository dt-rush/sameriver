/*
 * a list of entities which will be updated by another goroutine maybe,
 * which has a mutex that the user can lock when they wish to look at the
 * current contents. Can be sorted (needed by the data structure / algorithm
 * used in CollisionSystem)
 *
 */

package engine

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// A list of entities which is can be regularly updated by one goroutine
// while another reads and uses it
type UpdatedEntityList struct {
	// the entities in the list (tagged with gen)
	Entities []EntityToken
	// used to protect the Entities slice when adding or removing an entity
	Mutex sync.RWMutex
	// the channel along which updates to the list will come
	InputChannel chan EntityToken
	// a list of ID's the channel has yet to check in being created
	backlog []EntityToken
	// a function used to test if an entity belongs in the list
	// (supplied by EntityManager, it will have a reference to the EntityManager
	// and will run a query's Test() function against the entity)
	backlogTester func(entity EntityToken) bool
	// set to true while we're processing the backlog
	processingBacklog bool
	// used to stop the update loop's goroutine in
	// the case that they're done with the list (by calling Stop())
	stopUpdateLoopChannel chan bool
	// the name of the list (for debugging)
	Name string
	// the ID of the list (used for memory management)
	ID int
	// whether the entities slice should be sorted
	Sorted bool
	// a slice of funcs who want to be called *before* the entity gets
	// added/removed (that is, before the mutex unlocks)
	callbacks []func(EntityToken)
}

// create a new UpdatedEntityList by giving it a channel on which it will
// receive entity IDs
func NewUpdatedEntityList(
	Name string,
	ID int,
	InputChannel chan EntityToken,
	backlog []EntityToken,
	backlogTester func(entity EntityToken) bool) *UpdatedEntityList {

	l := UpdatedEntityList{}
	l.Name = Name
	l.ID = ID
	l.InputChannel = InputChannel
	l.Entities = make([]EntityToken, 0)

	l.backlog = backlog
	l.backlogTester = backlogTester
	l.processingBacklog = len(backlog) > 0

	l.stopUpdateLoopChannel = make(chan (bool))
	return &l
}

// get the length of the list
func (l *UpdatedEntityList) Length() int {
	l.Mutex.RLock()
	l.Mutex.RUnlock()
	return len(l.Entities)
}

// get the first element of the list
func (l *UpdatedEntityList) First() (EntityToken, error) {
	l.Mutex.RLock()
	defer l.Mutex.RUnlock()
	if len(l.Entities) > 0 {
		return l.Entities[0], nil
	} else {
		return EntityToken{ID: -1}, errors.New("no first entity found")
	}
}

// called during the creation of a list. Starts a goroutine which listens
// on the channel and either inserts or deletes entities as appropriate
func (l *UpdatedEntityList) start() {
	updatedEntityListDebug("starting UpdatedEntityList %s...", l.Name)
	go func() {
	updateloop:
		for {
			select {
			case _ = <-l.stopUpdateLoopChannel:
				break updateloop
			case id := <-l.InputChannel:
				updatedEntityListDebug("%s received signal", l.Name)
				l.actOnEntitySignal(id)
			default:
				if l.processingBacklog {
					l.popBacklog()
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (l *UpdatedEntityList) popBacklog() {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if len(l.backlog) == 0 {
		l.processingBacklog = false
		l.backlogTester = nil
	}

	last_ix := len(l.backlog) - 1
	entity := l.backlog[last_ix]
	l.backlog = l.backlog[:last_ix]
	updatedEntityListDebug("popped %v from backlog for list %s", entity, l.Name)
	if l.backlogTester(entity) {
		l.actOnEntitySignal(entity)
	}
}

func (l *UpdatedEntityList) stop() {
	l.stopUpdateLoopChannel <- true
}

// acts on an ID signal, which is either an ID to insert or -(ID + 1) to remove
func (l *UpdatedEntityList) actOnEntitySignal(entitySignal EntityToken) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	// callbacks list want to be notified of each signal we get
	for _, callback := range l.callbacks {
		go callback(entitySignal)
	}
	if entitySignal.ID < 0 {
		entitySignal = RemovalToken(entitySignal)
		updatedEntityListDebug("%s got remove:%d", l.Name, entitySignal.ID)
		removeEntityTokenFromSlice(&l.backlog, entitySignal)
		l.remove(entitySignal)
	} else {
		updatedEntityListDebug("%s got insert:%d", l.Name, entitySignal.ID)
		removeEntityTokenFromSlice(&l.backlog, entitySignal)
		l.insert(entitySignal)
	}
	updatedEntityListDebug("%s now: %v", l.Name, l.Entities)
}

// inserts an entity into the list (private so only called by the update loop)
func (l *UpdatedEntityList) insert(e EntityToken) {
	// note: both sorted and regular list insert will not double-insert an entity
	// (this deals with certain cases when lists are created at the same time
	// as tags are being added and entities spawned, some entities being added
	// to tags as the first entity with that tag. This will most often occur
	if l.Sorted {
		SortedEntityTokenSliceInsertIfNotPresent(&l.Entities, e)
	} else if indexOfEntityTokenInSlice(&l.Entities, e) == -1 {
		l.Entities = append(l.Entities, e)
	}
}

// removes an entity from the list (private so only called by the update loop)
func (l *UpdatedEntityList) remove(e EntityToken) {
	if l.Sorted {
		SortedEntityTokenSliceRemove(&l.Entities, e)
	} else {
		removeEntityTokenFromSlice(&l.Entities, e)
	}
}

// add a callback to the callbacks slice
func (l *UpdatedEntityList) addCallback(callback func(e EntityToken)) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	l.callbacks = append(l.callbacks, callback)
}

// For printing the list
func (l *UpdatedEntityList) String() string {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	return fmt.Sprintf("%s", l.Entities)
}
