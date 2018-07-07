package main

type EntityClassTable struct {
	// classes stores references to entity classes, which can be
	// retrieved by string ("crow", "turtle", "bear") in GetEntityClass()
	classes map[string]EntityClass
}

func NewEntityClassTable() *EntityClassTable {
	return &EntityClassTable{
		classes: make(map[string]EntityClass),
	}
}

// Register an entity class (subsequently retrievable)
func (t *EntityClassTable) AddEntityClass(c EntityClass) {
	t.classes[c.Name()] = c
}

// Get an entity class by name
func (t *EntityClassTable) GetEntityClass(name string) EntityClass {
	return t.classes[name]
}
