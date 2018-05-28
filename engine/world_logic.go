package engine

import (
	"sync"
	"sync/atomic"
	"time"
)

type WorldLogicFunc func(
	em *EntityManager,
	ev *EventBus,
	wl *WorldLogicManager)

type WorldLogic struct {
	Name   string
	sleep  time.Duration
	f      WorldLogicFunc
	active uint32
}

type WorldLogicManager struct {
	em     *EntityManager
	ev     *EventBus
	Logics map[string]*WorldLogic
	lists  map[string]*UpdatedEntityList
	mutex  sync.RWMutex
}

func (wl *WorldLogicManager) Init(
	em *EntityManager,
	ev *EventBus) {

	wl.em = em
	wl.ev = ev
	wl.Logics = make(map[string]*WorldLogic)
	wl.lists = make(map[string]*UpdatedEntityList)
}

func (wl *WorldLogicManager) AddList(
	name string, query GenericEntityQuery) {

	wl.mutex.Lock()
	defer wl.mutex.Unlock()
	wl.lists[name] = wl.em.GetUpdatedActiveEntityList(
		name, query)
}

func (wl *WorldLogicManager) GetEntitiesFromList(name string) []EntityToken {
	wl.mutex.RLock()
	wl.lists[name].Mutex.RLock()
	defer wl.lists[name].Mutex.RUnlock()
	defer wl.mutex.RUnlock()

	entities := wl.lists[name].Entities
	copyOfEntities := make([]EntityToken, len(entities))
	copy(copyOfEntities, entities)
	return copyOfEntities
}

// NOTE: we make no check of whether the key exists in the map,
// if used improperly, this will cause panics
func (wl *WorldLogicManager) ActivateLogic(name string) {
	wl.mutex.RLock()
	defer wl.mutex.RUnlock()

	Logic := wl.Logics[name]
	atomic.StoreUint32(&Logic.active, 1)
	go wl.run(name)
}

// NOTE: we make no check of whether the key exists in the map,
// if used improperly, this will cause panics
func (wl *WorldLogicManager) DeactivateLogic(name string) {
	wl.mutex.RLock()
	defer wl.mutex.RUnlock()

	Logic := wl.Logics[name]
	atomic.StoreUint32(&Logic.active, 0)
}

func (wl *WorldLogicManager) IsActive(name string) bool {
	wl.mutex.RLock()
	defer wl.mutex.RUnlock()

	Logic := wl.Logics[name]
	return atomic.LoadUint32(&Logic.active) == 1
}

func (wl *WorldLogicManager) AddLogic(
	name string, sleep time.Duration, f WorldLogicFunc) {

	wl.mutex.Lock()
	defer wl.mutex.Unlock()

	wl.Logics[name] = &WorldLogic{
		Name:  name,
		sleep: sleep,
		f:     f}

	go wl.run(name)
}

func (wl *WorldLogicManager) run(name string) {

	// set the logic active
	wl.mutex.RLock()
	Logic := wl.Logics[name]
	wl.mutex.RUnlock()
	atomic.StoreUint32(&Logic.active, 1)

	worldLogicDebug("running %s...", name)

	// while it's active, invoke the function and sleep in a loop
	for atomic.LoadUint32(&Logic.active) == 1 {
		// the WorldLogicFunc wants to be invoked with
		// (*EntityManager, *EventBus, *WorldLogicManager)
		Logic.f(wl.em, wl.ev, wl)
		time.Sleep(Logic.sleep)
	}
}
