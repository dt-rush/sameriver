package engine

type EntityLogicUnit struct {
	f      func()
	active bool
}

type EntityLogicTable struct {
	Logics map[EntityToken]EntityLogicUnit
}

func (t *EntityLogicTable) Init(em *EntityManager) {
	t.Logics = make(map[EntityToken]EntityLogicUnit)
}

func (t *EntityLogicTable) SetLogic(e EntityToken, f func()) {
	t.Logics[entity] = EntityLogicUnit{f, false}
}

func (t *EntityLogicTable) ActivateLogic(entity EntityToken) {
	if _, exists := t.Logics[entity]; exists {
		t.Logics[entity].Active = true
	}
}

func (t *EntityLogicTable) DeactivateLogic(entity EntityToken) {
	if _, exists := t.Logics[entity]; exists {
		t.Logics[entity].Active = true
	}
}

func (t *EntityLogicTable) DeleteLogic(entity EntityToken) {
	if _, exists := t.Logics[entity]; exists {
		delete(t.Logics, entity)
	}
}
