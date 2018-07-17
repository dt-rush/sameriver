package engine

func (w *World) QueueSpawn(req SpawnRequestData) {
	go func() {
		w.em.spawnSubscription.C <- Event{SPAWNREQUEST_EVENT, req}
	}()
}

func (w *World) Spawn(req SpawnRequestData) (*Entity, error) {
	return w.em.spawn(req)
}

func (w *World) SpawnUnique(
	tag string, req SpawnRequestData) (*Entity, error) {

	return w.em.spawnUnique(tag, req)
}

func (w *World) QueueDespawn(req DespawnRequestData) {
	go func() {
		w.em.despawnSubscription.C <- Event{DESPAWNREQUEST_EVENT, req}
	}()
}

func (w *World) Despawn(e *Entity) {
	w.em.doDespawn(e)
	w.RemoveEntityLogic(e)
}

func (w *World) Activate(entity *Entity) {
	w.em.Activate(entity)
}

func (w *World) Deactivate(entity *Entity) {
	w.em.Deactivate(entity)
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

func (w *World) UniqueTaggedEntity(tag string) (*Entity, error) {
	return w.em.UniqueTaggedEntity(tag)
}

func (w *World) EntitiesWithTag(tag string) *UpdatedEntityList {
	return w.em.EntitiesWithTag(tag)
}

func (w *World) EntityHasComponent(entity *Entity, COMPONENT int) bool {
	return w.em.EntityHasComponent(entity, COMPONENT)
}

func (w *World) EntityHasTag(entity *Entity, tag string) bool {
	return w.em.EntityHasTag(entity, tag)
}

func (w *World) TagEntity(entity *Entity, tags ...string) {
	w.em.TagEntity(entity, tags...)
}

func (w *World) TagEntities(entities []*Entity, tag string) {
	w.em.TagEntities(entities, tag)
}

func (w *World) UntagEntity(entity *Entity, tag string) {
	w.em.UntagEntity(entity, tag)
}

func (w *World) UntagEntities(entities []*Entity, tag string) {
	w.em.UntagEntities(entities, tag)
}

func (w *World) NumEntities() (total int, active int) {
	return w.em.NumEntities()
}

func (w *World) GetCurrentEntities() []*Entity {
	return w.em.GetCurrentEntities()
}

func (w *World) EntityManagerString() string {
	return w.em.String()
}

func (w *World) DumpEntities() string {
	return w.em.Dump()
}
