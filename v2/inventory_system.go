package sameriver

type InventorySystem struct {
	InventoryEntities *UpdatedEntityList
	itemSystem        *ItemSystem `sameriver-system-dependency:"-"`
}

func NewInventorySystem() *InventorySystem {
	return &InventorySystem{}
}

func (i *InventorySystem) Create(listing map[string]int) *Inventory {
	result := NewInventory()
	for arch, count := range listing {
		item := i.itemSystem.CreateStackSimple(count, arch)
		result.Credit(item)
	}
	return result
}

// System funcs

func (i *InventorySystem) GetComponentDeps() []string {
	return []string{"Generic,Inventory"}
}

func (i *InventorySystem) LinkWorld(w *World) {

	i.InventoryEntities = w.em.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"inventory",
			w.em.components.BitArrayFromNames([]string{"Inventory"})))
}

func (i *InventorySystem) Update(dt_ms float64) {
	// nil?
}

func (i *InventorySystem) Expand(n int) {
	// nil?
}
