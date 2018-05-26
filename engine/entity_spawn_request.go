package engine

// Used to spawn entities
type EntitySpawnRequest struct {
	Components ComponentSet
	Tags       []string
}
