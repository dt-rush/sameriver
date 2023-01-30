package engine

type System interface {
	LinkWorld(w *World)
	Update()
	GetComponentDeps() []string
}
