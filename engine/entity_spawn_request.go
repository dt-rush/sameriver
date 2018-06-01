package engine

// Used to spawn entities
type EntitySpawnRequest struct {
	Components ComponentSet
	Logic      EntityLogicFunc
	Tags       []string
	UniqueTag  string
}
