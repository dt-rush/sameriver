package engine

type EntityManagerInterface interface {
	Update(allowance float64) float64

	Spawn(tags []string, components ComponentSet) (*Entity, error)
	QueueSpawn(tags []string, components ComponentSet)
	SpawnUnique(
		uniqueTag string, tags []string, components ComponentSet) (*Entity, error)
	QueueSpawnUnique(uniqueTag string, tags []string, components ComponentSet)
	Despawn(e *Entity)
	QueueDespawn(e *Entity)
	DespawnAll()
	Activate(e *Entity)
	Deactivate(e *Entity)
	ApplyComponentSet(cs ComponentSet) func(*Entity)

	TagEntity(e *Entity, tags ...string)
	TagEntities(entities []*Entity, tag string)
	UntagEntity(e *Entity, tag string)
	UntagEntities(entities []*Entity, tag string)

	NumEntities() (total int, active int)
	GetCurrentEntitiesSet() map[*Entity]bool
	GetCurrentEntitiesSetCopy() map[*Entity]bool

	UniqueTaggedEntity(tag string) (*Entity, error)
	EntitiesWithTag(tag string) *UpdatedEntityList
	EntityHasComponent(e *Entity, COMPONENT int) bool
	EntityHasTag(e *Entity, tag string) bool

	GetUpdatedEntityList(q EntityFilter) *UpdatedEntityList
	GetSortedUpdatedEntityList(q EntityFilter) *UpdatedEntityList
	GetUpdatedEntityListByName(name string) *UpdatedEntityList

	String() string
	DumpEntities() string
}
