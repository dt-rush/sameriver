package engine

// Used to represent an entity with an ID at a point in time. Despawning the
// entity at a given ID will increment gen (gen ("generation") data is stored
// in EntityTable). The token storing *gen* prevents goroutines from
// requesting modifications on an entity after it has been despawened, or
// once a new entity has been spawned with its ID, which may happen quite
// readily otherwise
type EntityToken struct {
	ID  int32
	gen uint32
}
