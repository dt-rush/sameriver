package engine

import (
	"go.uber.org/atomic"
	"sync"
)

type EntityLogicUnit struct {
	stopChannel chan bool
	f           EntityLogicFunc
	running     *atomic.Uint32
}

type EntityLogicTable struct {
	em               *EntityManager
	entityLogicUnits map[EntityToken]*EntityLogicUnit
	mutex            sync.RWMutex
}

func (t *EntityLogicTable) Init(em *EntityManager) {
	t.em = em
	t.entityLogicUnits = make(map[EntityToken]*EntityLogicUnit)
}

func (t *EntityLogicTable) setLogic(
	entity EntityToken, f EntityLogicFunc) EntityLogicUnit {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if logicUnit, exists := t.entityLogicUnits[entity]; exists {
		logicUnit.stopChannel <- true
		logicUnit.running.Store(0)
	}
	logicUnit := EntityLogicUnit{make(chan bool), f, atomic.NewUint32(0)}
	t.entityLogicUnits[entity] = &logicUnit
	return logicUnit
}

func (t *EntityLogicTable) StartLogic(entity EntityToken) {
	t.mutex.RLock()
	t.mutex.RUnlock()
	if logicUnit, exists := t.entityLogicUnits[entity]; exists {
		if logicUnit.running.CAS(0, 1) {
			go logicUnit.f(entity, logicUnit.stopChannel, t.em)
		}
	}
}

func (t *EntityLogicTable) StopLogic(entity EntityToken) {
	t.mutex.RLock()
	t.mutex.RUnlock()
	if logicUnit, exists := t.entityLogicUnits[entity]; exists {
		if logicUnit.running.CAS(1, 0) {
			logicUnit.stopChannel <- true
		}
	}
}

func (t *EntityLogicTable) getLogic(entity EntityToken) *EntityLogicUnit {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if unit, exists := t.entityLogicUnits[entity]; exists {
		return unit
	}
	return nil
}

func (t *EntityLogicTable) deleteLogic(entity EntityToken) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exists := t.entityLogicUnits[entity]; exists {
		t.StopLogic(entity)
		delete(t.entityLogicUnits, entity)
	}
}
