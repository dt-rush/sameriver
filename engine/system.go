package engine

type System interface {
	LinkWorld(w *World)
	Update(dt_ms float64)
}
