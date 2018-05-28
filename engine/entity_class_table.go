package engine

import (
	"sync"
)

type entityClassTable struct {
	classes map[string]EntityClass
	mutex   sync.RWMutex
}

func (ect *entityClassTable) Init() {
	ect.classes = make(map[string]EntityClass)
}

func (ect *entityClassTable) addEntityClass(ec EntityClass) {
	ect.mutex.Lock()
	defer ect.mutex.Unlock()
	ect.classes[ec.Name()] = ec
}

func (ect *entityClassTable) getClass(name string) EntityClass {
	ect.mutex.RLock()
	defer ect.mutex.RUnlock()
	return ect.classes[name]
}
