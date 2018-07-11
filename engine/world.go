package engine

type World struct {
	Width   int
	Height  int
	ev      *EventBus
	em      *EntityManager
	wl      *WorldLogicManager
	systems []System
}

func NewWorld(width int, height int) *World {
	ev := NewEventBus()
	em := NewEntityManager(ev)
	w := World{
		Width:  width,
		Height: height,
		ev:     ev,
		em:     em,
	}
	w.wl = NewWorldLogicManager(&w)
	return &w
}

func (w *World) AddSystem(s System) {
	w.systems = append(w.systems, s)
	s.LinkWorld(w)
}

func (w *World) Update(dt_ms float64) {
	for _, s := range w.systems {
		s.Update(dt_ms)
	}
}
