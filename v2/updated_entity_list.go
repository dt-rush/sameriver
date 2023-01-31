// a list of entities which have active = true, which receives updates to its
// contents via the Signal() method
package sameriver

import (
	"errors"
	"math/rand"
)

// used to communicate insert / remove events
type EntitySignalType int

const (
	ENTITY_ADD    = iota
	ENTITY_REMOVE = iota
)

type EntitySignal struct {
	SignalType EntitySignalType
	Entity     *Entity
}

type UpdatedEntityList struct {
	// the entities in the list (tagged with gen)
	entities []*Entity
	// possibly nil Filter defining the list
	Filter *EntityFilter
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
		entities: make([]*Entity, 0),
		sorted:   false,
	}
	return &l
}

// create a new SORTED UpdatedEntityList by giving it a channel on which it will
// receive entity IDs
func NewSortedUpdatedEntityList() *UpdatedEntityList {
	l := UpdatedEntityList{
		entities: make([]*Entity, 0),
		sorted:   true,
	}
	return &l
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
func (l *UpdatedEntityList) add(e *Entity) {
	// NOTE: idempotent
	lenBefore := len(l.entities)
	if l.sorted {
		SortedEntitySliceInsertIfNotPresent(&l.entities, e)
	} else {
		if indexOfEntityInSlice(&l.entities, e) == -1 {
			l.entities = append(l.entities, e)
		}
	}
	if len(l.entities) == lenBefore+1 {
		e.Lists = append(e.Lists, l)
	}
}

// removes an entity from the list (private so only called by the update loop)
func (l *UpdatedEntityList) remove(e *Entity) {
	// NOTE: both idempotent
	if l.sorted {
		SortedEntitySliceRemove(&l.entities, e)
	} else {
		removeEntityFromSlice(&l.entities, e)
	}
	removeUpdatedEntityListFromSlice(&e.Lists, l)
}

// add a callback to the callbacks slice
func (l *UpdatedEntityList) AddCallback(
	callback func(EntitySignal)) {

	l.callbacks = append(l.callbacks, callback)
}

// get the length of the list
func (l *UpdatedEntityList) Length() int {
	return len(l.entities)
}

// get the contents via copy
func (l *UpdatedEntityList) GetEntities() []*Entity {
	copyOfentities := make([]*Entity, len(l.entities))
	copy(copyOfentities, l.entities)
	return copyOfentities
}

// get the first element of the list
func (l *UpdatedEntityList) FirstEntity() (*Entity, error) {
	if len(l.entities) == 0 {
		return nil, errors.New("list is empty, no first element")
	}
	return l.entities[0], nil
}

// get a random element of the list
func (l *UpdatedEntityList) RandomEntity() (*Entity, error) {
	if len(l.entities) == 0 {
		return nil, errors.New("list is empty, can't get random element")
	}
	return l.entities[rand.Intn(len(l.entities))], nil
}

// For printing the list
func (l *UpdatedEntityList) String() string {
	return EntitySliceToString(l.entities)
}
