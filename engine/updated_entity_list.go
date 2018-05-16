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

const VERBOSE = false

type UpdatedEntityList struct {
	Watcher           QueryWatcher
	Mutex             sync.Mutex
	Entities          []uint16
	StopUpdateChannel chan (bool)
	Name              string
}

func NewUpdatedEntityList(watcher QueryWatcher, name string) *UpdatedEntityList {
	l := UpdatedEntityList{}
	l.Watcher = watcher
	l.Entities = make([]uint16, 0)
	l.StopUpdateChannel = make(chan (bool))
	l.Name = name
	l.start()
	return &l
}

func (l *UpdatedEntityList) start() {
	go func() {
	updateloop:
		for {
			select {
			case _ = <-l.StopUpdateChannel:
				break updateloop
			case id := <-l.Watcher.Channel:
				l.Mutex.Lock()
				if id >= 0 {
					if VERBOSE {
						Logger.Printf("[Updated entity list] %s got insert:%d\n", l.Name, id)
					}
					l.insert(uint16(id))
				} else {
					id = -(id + 1)
					if VERBOSE {
						Logger.Printf("[Updated entity list] %s got remove:%d\n", l.Name, id)
					}
					l.remove(uint16(id))
				}
				l.Mutex.Unlock()
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (l *UpdatedEntityList) insert(id uint16) {
	l.Entities = append(l.Entities, id)
}

func (l *UpdatedEntityList) remove(id uint16) {
	last_ix := len(l.Entities) - 1
	for i := uint16(0); i <= uint16(last_ix); i++ {
		if l.Entities[i] == id {
			l.Entities[i] = l.Entities[last_ix]
			l.Entities = l.Entities[:last_ix]
			break
		}
	}
}
