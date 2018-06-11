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
	"math/rand"
	"time"
)

// A list of entities which is can be regularly updated by one goroutine
// while another reads and uses it
type UpdatedEntityList struct {
	// the query watcher this list is attached to
	qw EntityQueryWatcher
	// the entities in the list (tagged with gen)
	Entities []*EntityToken
	// a list of ID's the channel has yet to check in being created
	backlog []*EntityToken
	// a function used to test if an entity belongs in the list
	// (supplied by EntityManager, it will have a reference to the EntityManager
	// and will run a query's Test() function against the entity)
	backlogTester func(entity *EntityToken) bool
	// set to true while we're processing the backlog
	processingBacklog bool
	// used to stop the update loop's goroutine in
	// the case that they're done with the list (by calling Stop())
	stopUpdateLoopChannel chan bool
	// whether the entities slice should be sorted
	Sorted bool
	// a slice of funcs who want to be called *before* the entity gets
	// added/removed (that is, before the mutex unlocks)
	callbacks []func(EntitySignal)
}

// create a new UpdatedEntityList by giving it a channel on which it will
// receive entity IDs
func NewUpdatedEntityList(
	qw EntityQueryWatcher,
	backlog []*EntityToken,
	backlogTester func(entity *EntityToken) bool) *UpdatedEntityList {

	l := UpdatedEntityList{}
	l.Entities = make([]*EntityToken, 0)
	l.qw = qw
	l.backlog = backlog
	l.backlogTester = backlogTester
	l.processingBacklog = len(backlog) > 0
	l.stopUpdateLoopChannel = make(chan (bool))
	return &l
}

// get the length of the list
func (l *UpdatedEntityList) Length() int {
	return len(l.Entities)
}

// get the first element of the list
func (l *UpdatedEntityList) FirstEntity() (*EntityToken, error) {
	if len(l.Entities) == 0 {
		return nil, errors.New("list is empty, no first element")
	}
	return l.Entities[0], nil
}

// get a random element of the list
func (l *UpdatedEntityList) RandomEntity() (*EntityToken, error) {
	if len(l.Entities) == 0 {
		return nil, errors.New("list is empty, can't get random element")
	}
	return l.Entities[rand.Intn(len(l.Entities))], nil
}

// called during the creation of a list. Starts a goroutine which listens
// on the channel and either adds or removes entities as appropriate
func (l *UpdatedEntityList) start() {
	go func() {
	updateloop:
		for {
			select {
			case _ = <-l.stopUpdateLoopChannel:
				break updateloop
			case signal := <-l.qw.Channel:
				l.actOnSignal(signal)
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

	if len(l.backlog) == 0 {
		l.processingBacklog = false
		l.backlogTester = nil
		return
	}

	last_ix := len(l.backlog) - 1
	entity := l.backlog[last_ix]
	l.backlog = l.backlog[:last_ix]
	if l.backlogTester(entity) {
		l.actOnSignal(EntitySignal{ENTITY_ADD, entity})
	}
}

func (l *UpdatedEntityList) stop() {
	l.stopUpdateLoopChannel <- true
}

func (l *UpdatedEntityList) actOnSignal(signal EntitySignal) {

	// callbacks list want to be notified of each signal we get
	for _, callback := range l.callbacks {
		go callback(signal)
	}
	// act on signal
	switch signal.signalType {
	case ENTITY_REMOVE:
		removeEntityTokenFromSlice(&l.backlog, signal.entity)
		l.remove(signal.entity)
	case ENTITY_ADD:
		removeEntityTokenFromSlice(&l.backlog, signal.entity)
		l.add(signal.entity)
	}
}

// adds an entity into the list (private so only called by the update loop)
func (l *UpdatedEntityList) add(e *EntityToken) {
	// note: both sorted and regular list add will not double-add an entity
	if l.Sorted {
		SortedEntityTokenSliceInsertIfNotPresent(&l.Entities, e)
	} else if indexOfEntityTokenInSlice(&l.Entities, e) == -1 {
		l.Entities = append(l.Entities, e)
	}
}

// removes an entity from the list (private so only called by the update loop)
func (l *UpdatedEntityList) remove(e *EntityToken) {
	if l.Sorted {
		SortedEntityTokenSliceRemove(&l.Entities, e)
	} else {
		removeEntityTokenFromSlice(&l.Entities, e)
	}
}

// add a callback to the callbacks slice
func (l *UpdatedEntityList) addCallback(
	callback func(EntitySignal)) {

	l.callbacks = append(l.callbacks, callback)
}

// For printing the list
func (l *UpdatedEntityList) String() string {
	return fmt.Sprintf("%s", l.Entities)
}
