package sameriver

type System interface {
	LinkWorld(w *World)
	Update(dt_ms float64)
	GetComponentDeps() map[ComponentID]ComponentKind
	Expand(n int)
}
