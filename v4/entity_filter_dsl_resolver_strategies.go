package sameriver

type IdentifierResolver interface {
	Resolve(identifier string) any
}

type EntityResolver struct {
	entity *Entity
}

func (er *EntityResolver) Resolve(identifier string) any {
	// TODO
	return nil
}

type WorldResolver struct {
	world *World
}

func (wr *WorldResolver) Resolve(identifier string) any {
	// TODO
	return nil
}
