package sameriver

type System interface {
	LinkWorld(w *World)
	Update(dt_ms float64)
	// return an array specifying the component dependencies
	// where every 3 elements groups together into a component spec
	// ComponentID, ComponentKind, string
	// eg.
	// POSITION, VEC2D, "POSITION"
	GetComponentDeps() []any
	Expand(n int)
}
