package engine

import (
	"go.uber.org/atomic"
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
	active atomic.Uint32
}

type WorldLogicManager struct {
	em     *EntityManager
	ev     *EventBus
	Logics map[string]*WorldLogic
	lists  map[string]*UpdatedEntityList
}

func (wl *WorldLogicManager) Init(
	em *EntityManager,
	ev *EventBus) {

	wl.em = em
	wl.ev = ev
	wl.Logics = make(map[string]*WorldLogic)
	wl.lists = make(map[string]*UpdatedEntityList)
}

func (wl *WorldLogicManager) AddList(query EntityQuery) {

	wl.lists[query.Name] = wl.em.GetUpdatedEntityList(query)
}

func (wl *WorldLogicManager) GetEntitiesFromList(name string) []EntityToken {

	entities := wl.lists[name].Entities
	copyOfEntities := make([]EntityToken, len(entities))
	copy(copyOfEntities, entities)
	return copyOfEntities
}

// NOTE: we make no check of whether the key exists in the map,
// if used improperly, this will cause panics
func (wl *WorldLogicManager) ActivateLogic(name string) {

	Logic := wl.Logics[name]
	Logic.active.Store(1)
	go wl.run(name)
}

// NOTE: we make no check of whether the key exists in the map,
// if used improperly, this will cause panics
func (wl *WorldLogicManager) DeactivateLogic(name string) {

	Logic := wl.Logics[name]
	Logic.active.Store(0)
}

func (wl *WorldLogicManager) IsActive(name string) bool {

	Logic := wl.Logics[name]
	return Logic.active.Load() == 1
}

func (wl *WorldLogicManager) AddLogic(
	name string, sleep time.Duration, f WorldLogicFunc) {

	wl.Logics[name] = &WorldLogic{
		Name:  name,
		sleep: sleep,
		f:     f}

	go wl.run(name)
}

// TODO: revisit how this works in light of the new sync gameloop we're working
// with
func (wl *WorldLogicManager) run(name string) {
	// set the logic active
	Logic := wl.Logics[name]
	Logic.active.Store(1)

	worldLogicDebug("running %s...", name)

	// while it's active, invoke the function and sleep in a loop
	for Logic.active.Load() == 1 {
		// the WorldLogicFunc wants to be invoked with
		// (*EntityManager, *EventBus, *WorldLogicManager)
		Logic.f(wl.em, wl.ev, wl)
		time.Sleep(Logic.sleep)
	}
}
