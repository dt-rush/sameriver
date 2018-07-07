package engine

import (
	"time"
)

type WorldLogicManager struct {
	em     *EntityManager
	ev     *EventBus
	Logics map[string]*LogicUnit
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

func (wl *WorldLogicManager) ActivateLogic(name string) {
	if l, ok := wl.Logics[name]; ok {
		l.Active = true
	}
}

func (wl *WorldLogicManager) DeactivateLogic(name string) {
	if l, ok := wl.Logics[name]; ok {
		l.Active = false
	}
}

func (wl *WorldLogicManager) IsActive(name string) bool {
	if l, ok := wl.Logics[name]; ok {
		return l.Active
	} else {
		return false
	}
}

func (wl *WorldLogicManager) AddLogic(Logic *LogicUnit) {
	wl.Logics[Logic.Name] = Logic
}

func (wl *WorldLogicManager) run(name string) {
	Logic := wl.Logics[name]
	if Logic.Active {
		Logic.f()
	}
}
