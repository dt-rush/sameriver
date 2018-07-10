package engine

type WorldLogicManager struct {
	w      *World
	Logics map[string]*LogicUnit
}

func NewWorldLogicManager(w *World) *WorldLogicManager {
	return &WorldLogicManager{
		w:      w,
		Logics: make(map[string]*LogicUnit),
	}
}

func (wl *WorldLogicManager) GetEntitiesFromList(name string) []*EntityToken {
	var entities []*EntityToken
	list := wl.w.em.GetUpdatedEntityListByName(name)
	if list != nil {
		entities = list.Entities
	}
	copyOfEntities := make([]*EntityToken, len(entities))
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
		Logic.F()
	}
}
