package engine

// a collision has occured between two entities
type CollisionData struct {
	EntityA *Entity
	EntityB *Entity
}

// the EntityManager is requested to spawn an entity
type SpawnRequestData struct {
	Components ComponentSet
	Tags       []string
}

// the EntityManager is requested to despawn an entity
type DespawnRequestData struct {
	Entity *Entity
}
