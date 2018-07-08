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
	// whether the entities slice should be sorted
	sorted bool
	// a slice of funcs who want to be called *before* the entity gets
	// added/removed
	callbacks []func(EntitySignal)
}

// create a new UpdatedEntityList by giving it a channel on which it will
// receive entity IDs
func NewUpdatedEntityList() *UpdatedEntityList {
	l := UpdatedEntityList{
		Entities: make([]*EntityToken, 0),
		sorted:   false,
	}
	return &l
}

// create a new SORTED UpdatedEntityList by giving it a channel on which it will
// receive entity IDs
func NewSortedUpdatedEntityList() *UpdatedEntityList {
	l := UpdatedEntityList{
		Entities: make([]*EntityToken, 0),
		sorted:   true,
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

func (l *UpdatedEntityList) Signal(signal EntitySignal) {
	// callbacks list want to be notified of each signal we get
	for _, callback := range l.callbacks {
		callback(signal)
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
	// NOTE: idempotent
	lenBefore := len(l.Entities)
	if l.sorted {
		SortedEntityTokenSliceInsertIfNotPresent(&l.Entities, e)
	} else {
		if indexOfEntityTokenInSlice(&l.Entities, e) == -1 {
			l.Entities = append(l.Entities, e)
		}
	}
	if len(l.Entities) == lenBefore+1 {
		e.Lists = append(e.Lists, l)
	}
}

// removes an entity from the list (private so only called by the update loop)
func (l *UpdatedEntityList) remove(e *EntityToken) {
	// NOTE: both idempotent
	if l.sorted {
		SortedEntityTokenSliceRemove(&l.Entities, e)
	} else {
		removeEntityTokenFromSlice(&l.Entities, e)
	}
	removeUpdatedEntityListFromSlice(&e.Lists, l)
}

// add a callback to the callbacks slice
func (l *UpdatedEntityList) addCallback(
	callback func(EntitySignal)) {

	l.callbacks = append(l.callbacks, callback)
}

// For printing the list
func (l *UpdatedEntityList) String() string {
	return fmt.Sprintf("%v", l.Entities)
}
