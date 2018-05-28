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
	// the query watcher this list is attached to
	qw EntityQueryWatcher
	// the entities in the list (tagged with gen)
	Entities []EntityToken
	// used to protect the Entities slice when adding or removing an entity
	Mutex sync.RWMutex
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
	// whether the entities slice should be sorted
	Sorted bool
	// a slice of funcs who want to be called *before* the entity gets
	// added/removed (that is, before the mutex unlocks)
	callbacks []func(EntityToken)
}

// create a new UpdatedEntityList by giving it a channel on which it will
// receive entity IDs
func NewUpdatedEntityList(
	qw EntityQueryWatcher,
	backlog []EntityToken,
	backlogTester func(entity EntityToken) bool) *UpdatedEntityList {

	l := UpdatedEntityList{}
	l.Entities = make([]EntityToken, 0)
	l.qw = qw
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
// on the channel and either adds or removes entities as appropriate
func (l *UpdatedEntityList) start() {
	updatedEntityListDebug("starting UpdatedEntityList %s...", l.qw.Name)
	go func() {
	updateloop:
		for {
			select {
			case _ = <-l.stopUpdateLoopChannel:
				break updateloop
			case id := <-l.qw.Channel:
				updatedEntityListDebug("%s received signal", l.qw.Name)
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
		return
	}

	last_ix := len(l.backlog) - 1
	entity := l.backlog[last_ix]
	l.backlog = l.backlog[:last_ix]
	updatedEntityListDebug("popped %v from backlog for list %s",
		entity, l.qw.Name)
	if l.backlogTester(entity) {
		l.actOnSignal(EntitySignal{ENTITY_ADD, entity})
	}
}

func (l *UpdatedEntityList) stop() {
	l.stopUpdateLoopChannel <- true
}

func (l *UpdatedEntityList) actOnEntitySignal(signal EntitySignal) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	defer updatedEntityListDebug("%s now: %v", l.qw.Name, l.Entities)

	// callbacks list want to be notified of each signal we get
	for _, callback := range l.callbacks {
		go callback(entitySignal)
	}
	// act on signal
	switch signal.signalType {
	case ENTITY_REMOVE:
		updatedEntityListDebug("%s got remove:%d", l.qw.Name, entitySignal.ID)
		removeEntityTokenFromSlice(&l.backlog, entitySignal)
		l.remove(entitySignal)
	case ENTITY_ADD:
		updatedEntityListDebug("%s got add:%d", l.qw.Name, entitySignal.ID)
		removeEntityTokenFromSlice(&l.backlog, entitySignal)
		l.add(entitySignal)
	}
}

// adds an entity into the list (private so only called by the update loop)
func (l *UpdatedEntityList) add(e EntityToken) {
	// note: both sorted and regular list add will not double-add an entity
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
