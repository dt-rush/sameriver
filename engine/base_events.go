package engine

// a collision has occured between two entities
type CollisionData struct {
	EntityA *EntityToken
	EntityB *EntityToken
}

// the EntityManager is requested to spawn an entity
type SpawnRequestData struct {
	Components ComponentSet
	Logic      func()
	Tags       []string
	UniqueTag  string
}

// the EntityManager is requested to despawn an entity
type DespawnRequestData struct {
	Entity *EntityToken
}
