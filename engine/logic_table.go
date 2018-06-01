package engine

import (
	"sync"
)

type LogicUnit struct {
	f           EntityLogicFunc
	stopChannel chan bool
}

type LogicTable struct {
	logicUnits map[int]*LogicUnit
	mutex      sync.RWMutex
}

func (t *LogicTable) Init() {
	t.logicUnits = make(map[EntityToken]*LogicUnit)
}

func (t *LogicTable) setLogic(entity EntityToken, f EntityLogicFunc) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if unit, exists := t.logicUnits[entity]; exists {
		unit.StopChannel <- true
	}
	t.logicUnits[id] = LogicUnit{f, make(chan bool)}
}

func (t *LogicTable) getLogic(id int) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if unit, exists := t.logicUnits[entity]; exists {
		return unit
	}
}
