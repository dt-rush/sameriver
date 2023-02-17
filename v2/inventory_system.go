package sameriver

type InventorySystem struct {
	w                 *World
	inventoryEntities *UpdatedEntityList
	archetypes        map[string]ItemSpec
}

func NewInventorySystem() *InventorySystem {
	return &InventorySystem{}
}

func (i *InventorySystem) RegisterArchetype(arch ItemSpec) {
	i.archetypes[arch.Name] = arch
}

func (i *InventorySystem) GetArchetype(name string) ItemSpec {
	return i.archetypes[name]
}

// System funcs

func (i *InventorySystem) GetComponentDeps() []string {
	return []string{"Generic,Inventory"}
}

func (i *InventorySystem) LinkWorld(w *World) {
	i.w = w
	i.inventoryEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"inventory",
			w.em.components.BitArrayFromNames([]string{"Inventory"})))
}

func (i *InventorySystem) Update(dt_ms float64) {
	// nil?
}
