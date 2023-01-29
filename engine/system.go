package engine

type System interface {
	LinkWorld(w *World)
	GetComponentDeps()
	Update()
}
