package engine

type CollisionData struct {
	EntityA EntityToken
	EntityB EntityToken
}

type SpawnRequestData struct {
	Components ComponentSet
	Logic      func()
	Tags       []string
	UniqueTag  string
}

type DespawnRequestData struct {
	Entity EntityToken
}

type LogicStartStopData struct {
	entity    EntityToken
	startStop bool
}
