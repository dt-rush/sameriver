package main

import (
	"sync"
)

// used to imitate the UpdatedEntityList from the engine
type UpdatedEntityList struct {
	Entities []EntityToken
	Mutex    sync.RWMutex
}

func (l *UpdatedEntityList) Length() int {
	l.Mutex.RLock()
	defer l.Mutex.RUnlock()
	return len(l.Entities)
}
