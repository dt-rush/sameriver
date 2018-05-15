/*
 *
 * a list of entities which will be updated by another goroutine maybe,
 * which has a mutex that the user can lock when they wish to look at the 
 * current contents
 *
*/


package engine

import (
	"sync"
)

type UpdatedEntityList struct {
	Watcher QueryWatcher
	Mutex sync.Mutex
	Entities []uint16
	StopUpdateChannel chan(bool)
}

func NewUpdatedEntityList (watcher QueryWatcher) UpdatedEntityList {
	l := UpdatedEntityList{}
	l.Watcher = watcher
	l.Entities = make([]uint16,0)
	l.StopUpdateChannel = make(chan(bool))
	l.start()
}

func (l *UpdatedEntityList) start() {
	go func () {
		for {
			select {
			case signal := <-l.StopUpdateChannel:
				return
			case id := <-l.Watcher.Channel:
				l.Mutex.Lock()
				if id >= 0 {
					l.insert (id)
				} else {
					l.remove (-(id+1))
				}
			}
		}
	}()
}

func (l *UpdatedEntityList) insert(id uint16) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	l.Entities = append (l.Entities, id)
}

func (l *UpdatedEntityList) remove(id uint16) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	last_ix := len(l.Entities) - 1
	for i := 0; i <= last_ix; i++ {
		if i == id {
			l.Entities[i] = l.Entities[last_ix]
			l.Entities = l.Entities[:last_ix]
		}
	}
}
