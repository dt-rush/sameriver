package engine

// a collision has occured between two entities
type CollisionData struct {
	This  *Entity
	Other *Entity
}

// the EntityManager is requested to spawn an entity
type SpawnRequestData struct {
	UniqueTag  string
	Tags       []string
	Components ComponentSet
}

// the EntityManager is requested to despawn an entity
type DespawnRequestData struct {
	Entity *Entity
}
