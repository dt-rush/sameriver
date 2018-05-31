package engine

type CollisionEvent struct {
	EntityA EntityToken
	EntityB EntityToken
}

type SpawnRequest struct {
	EntityType int
	Position   [2]int16
	Active     bool
}
