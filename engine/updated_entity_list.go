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
	"time"
)

type UpdatedEntityList struct {
	Watcher QueryWatcher
	Mutex sync.Mutex
	Entities []uint16
	StopUpdateChannel chan(bool)
	Name string
}

func NewUpdatedEntityList (watcher QueryWatcher, name string) *UpdatedEntityList {
	l := UpdatedEntityList{}
	l.Watcher = watcher
	l.Entities = make([]uint16,0)
	l.StopUpdateChannel = make(chan(bool))
	l.Name = name
	l.start()
	return &l
}

func (l *UpdatedEntityList) start() {
	go func () {
		for {
			select {
			case _ = <-l.StopUpdateChannel:
				return
			case id := <-l.Watcher.Channel:
				l.Mutex.Lock()
				Logger.Printf ("inserting #%d to %s\n", id, l.Name)
				if id >= 0 {
					l.insert (id)
				} else {
					l.remove (-(id+1))
				}
				l.Mutex.Unlock()
				Logger.Println (l.Entities)
			default:
				time.Sleep (100 * time.Millisecond)
			}
		}
	}()
}

func (l *UpdatedEntityList) insert(id uint16) {
	l.Entities = append (l.Entities, id)
}

func (l *UpdatedEntityList) remove(id uint16) {
	last_ix := len(l.Entities) - 1
	for i := uint16(0); i <= uint16(last_ix); i++ {
		if i == id {
			l.Entities[i] = l.Entities[last_ix]
			l.Entities = l.Entities[:last_ix]
		}
	}
}
