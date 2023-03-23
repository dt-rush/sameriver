package sameriver

type EntityManagerInterface interface {
	Update(allowance_ms float64) float64

	MaxEntities() int
	Components() *ComponentTable

	Spawn(spec map[string]any) *Entity
	QueueSpawn(spec map[string]any)
	Despawn(e *Entity)
	DespawnAll()
	Activate(e *Entity)
	Deactivate(e *Entity)

	TagEntity(e *Entity, tags ...string)
	TagEntities(entities []*Entity, tag string)
	UntagEntity(e *Entity, tag string)
	UntagEntities(entities []*Entity, tag string)

	NumEntities() (total int, active int)
	GetActiveEntitiesSet() map[*Entity]bool
	GetCurrentEntitiesSet() map[*Entity]bool
	GetCurrentEntitiesSetCopy() map[*Entity]bool

	UniqueTaggedEntity(tag string) (*Entity, error)
	UpdatedEntitiesWithTag(tag string) *UpdatedEntityList
	EntityHasComponent(e *Entity, name ComponentID) bool
	EntityHasTag(e *Entity, tag string) bool

	GetUpdatedEntityList(q EntityFilter) *UpdatedEntityList
	GetSortedUpdatedEntityList(q EntityFilter) *UpdatedEntityList
	GetUpdatedEntityListByName(name string) *UpdatedEntityList
	GetUpdatedEntityListByComponents(names []ComponentID) *UpdatedEntityList

	ApplyComponentSet(e *Entity, spec map[ComponentID]any)

	String() string
	DumpEntities() string
}
