package engine

// Used to represent an entity with an ID at a point in time. Despawning the
// entity at a given ID will increment gen (gen ("generation") data is stored
// in EntityTable). The token storing *gen* prevents goroutines from
// requesting modifications on an entity after it has been despawened, or
// once a new entity has been spawned with its ID, which may happen quite
// readily otherwise
type EntityToken struct {
	ID  int
	gen uint32
}

var ENTITY_TOKEN_NIL = EntityToken{-1, 0}

func RemovalToken(token EntityToken) EntityToken {
	return EntityToken{
		-(token.ID + 1),
		token.gen}
}
