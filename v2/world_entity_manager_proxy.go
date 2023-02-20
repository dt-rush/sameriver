package sameriver

func (w *World) MaxEntities() int {
	return w.em.MaxEntities()
}

func (w *World) Components() *ComponentTable {
	return w.em.components
}

func (w *World) Spawn(spec map[string]any) *Entity {
	return w.em.Spawn(spec)
}

func (w *World) QueueSpawn(spec map[string]any) {
	w.em.QueueSpawn(spec)
}

func (w *World) QueueDespawn(e *Entity) {
	w.em.QueueDespawn(e)
}

func (w *World) Despawn(e *Entity) {
	w.em.Despawn(e)
	w.RemoveAllEntityLogics(e)
}

func (w *World) DespawnAll() {
	for e, _ := range w.em.GetCurrentEntitiesSetCopy() {
		w.RemoveAllEntityLogics(e)
	}
	w.em.DespawnAll()
}

func (w *World) Activate(e *Entity) {
	w.em.Activate(e)
}

func (w *World) Deactivate(e *Entity) {
	w.em.Deactivate(e)
}

func (w *World) GetUpdatedEntityList(q EntityFilter) *UpdatedEntityList {
	return w.em.GetUpdatedEntityList(q)
}

func (w *World) GetSortedUpdatedEntityList(q EntityFilter) *UpdatedEntityList {
	return w.em.GetSortedUpdatedEntityList(q)
}

func (w *World) GetUpdatedEntityListByName(name string) *UpdatedEntityList {
	return w.em.GetUpdatedEntityListByName(name)
}

func (w *World) GetUpdatedEntityListByComponentNames(names []string) *UpdatedEntityList {
	return w.em.GetUpdatedEntityListByComponentNames(names)
}

func (w *World) UniqueTaggedEntity(tag string) (*Entity, error) {
	return w.em.UniqueTaggedEntity(tag)
}

func (w *World) UpdatedEntitiesWithTag(tag string) *UpdatedEntityList {
	return w.em.UpdatedEntitiesWithTag(tag)
}

func (w *World) EntityHasComponent(e *Entity, name string) bool {
	return w.em.EntityHasComponent(e, name)
}

func (w *World) EntityHasTag(e *Entity, tag string) bool {
	return w.em.EntityHasTag(e, tag)
}

func (w *World) TagEntity(e *Entity, tags ...string) {
	w.em.TagEntity(e, tags...)
}

func (w *World) TagEntities(entities []*Entity, tag string) {
	w.em.TagEntities(entities, tag)
}

func (w *World) UntagEntity(e *Entity, tag string) {
	w.em.UntagEntity(e, tag)
}

func (w *World) UntagEntities(entities []*Entity, tag string) {
	w.em.UntagEntities(entities, tag)
}

func (w *World) NumEntities() (total int, active int) {
	return w.em.NumEntities()
}

func (w *World) GetActiveEntitiesSet() map[*Entity]bool {
	return w.em.GetActiveEntitiesSet()
}

func (w *World) GetCurrentEntitiesSet() map[*Entity]bool {
	return w.em.GetCurrentEntitiesSet()
}

func (w *World) GetCurrentEntitiesSetCopy() map[*Entity]bool {
	return w.em.GetCurrentEntitiesSetCopy()
}

func (w *World) DumpEntities() string {
	return w.em.DumpEntities()
}
