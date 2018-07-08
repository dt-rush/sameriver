// a list of entities which have active = true, which receives updates to its
// contents by an EntityQueryWatcher
package engine

import (
	"errors"
	"fmt"
	"math/rand"
)

type UpdatedEntityList struct {
	// the entities in the list (tagged with gen)
	Entities []*EntityToken
	// possibly nil query defining the list
	Query *EntityQuery
	// a channel used to receive the entity add / remove signals
	Channel chan EntitySignal
	// used to stop the update loop's goroutine in
	// the case that they're done with the list (by calling Stop())
	stopUpdateLoopChannel chan bool
	// whether the entities slice should be sorted
	sorted bool
	// a slice of funcs who want to be called *before* the entity gets
	// added/removed
	callbacks []func(EntitySignal)
}

// create a new UpdatedEntityList by giving it a channel on which it will
// receive entity IDs
func NewUpdatedEntityList(channel chan EntitySignal) *UpdatedEntityList {
	l := UpdatedEntityList{
		Entities:              make([]*EntityToken, 0),
		Channel:               channel,
		stopUpdateLoopChannel: make(chan bool),
		sorted:                false,
	}
	return &l
}

// create a new SORTED UpdatedEntityList by giving it a channel on which it will
// receive entity IDs
func NewSortedUpdatedEntityList(channel chan EntitySignal) *UpdatedEntityList {
	l := UpdatedEntityList{
		Entities:              make([]*EntityToken, 0),
		Channel:               channel,
		stopUpdateLoopChannel: make(chan bool),
		sorted:                true,
	}
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
func (l *UpdatedEntityList) Start() {
	go func() {
	updateloop:
		for {
			select {
			case _ = <-l.stopUpdateLoopChannel:
				break updateloop
			case signal := <-l.Channel:
				l.actOnSignal(signal)
			}
		}
	}()
}

func (l *UpdatedEntityList) Stop() {
	l.stopUpdateLoopChannel <- true
}

func (l *UpdatedEntityList) actOnSignal(signal EntitySignal) {
	// callbacks list want to be notified of each signal we get
	for _, callback := range l.callbacks {
		go callback(signal)
	}
	// act on signal
	switch signal.SignalType {
	case ENTITY_ADD:
		l.add(signal.Entity)
	case ENTITY_REMOVE:
		l.remove(signal.Entity)
	}
}

// adds an entity into the list (private so only called by the update loop)
func (l *UpdatedEntityList) add(e *EntityToken) {
	// note: both sorted and regular list add will not double-add an entity
	if l.sorted {
		SortedEntityTokenSliceInsertIfNotPresent(&l.Entities, e)
	} else if indexOfEntityTokenInSlice(&l.Entities, e) == -1 {
		l.Entities = append(l.Entities, e)
	}
	e.ListsMutex.Lock()
	e.Lists = append(e.Lists, l)
	e.ListsMutex.Unlock()
}

// removes an entity from the list (private so only called by the update loop)
func (l *UpdatedEntityList) remove(e *EntityToken) {
	if l.sorted {
		SortedEntityTokenSliceRemove(&l.Entities, e)
	} else {
		removeEntityTokenFromSlice(&l.Entities, e)
	}
	e.ListsMutex.Lock()
	removeUpdatedEntityListFromSlice(&e.Lists, l)
	e.ListsMutex.Unlock()
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
