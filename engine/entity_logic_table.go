package engine

import (
	"sync"
)

type EntityLogicUnit struct {
	stopChannel chan bool
	f           EntityLogicFunc
}

type EntityLogicTable struct {
	entityLogicUnits map[EntityToken]*EntityLogicUnit
	mutex            sync.RWMutex
}

func (t *EntityLogicTable) Init() {
	t.entityLogicUnits = make(map[EntityToken]*EntityLogicUnit)
}

func (t *EntityLogicTable) setLogic(
	entity EntityToken, f EntityLogicFunc) EntityLogicUnit {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if unit, exists := t.entityLogicUnits[entity]; exists {
		unit.StopChannel <- true
	}
	unit := EntityLogicUnit{f, make(chan bool)}
	t.entityLogicUnits[id] = unit
	return unit
}

func (t *EntityLogicTable) getLogic(entity EntityToken) EntityLogicUnit {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if unit, exists := t.entityLogicUnits[entity]; exists {
		return unit
	}
}

func (t *EntityLogicTable) deleteLogic(entity EntityToken) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if unit, exists := t.entityLogicUnits[entity]; exists {
		delete(t.entityLogicUnits, entity)
	}
}
