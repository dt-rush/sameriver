package engine

type World struct {
	Width  int
	Height int
	ev     *EventBus
	em     *EntityManager
	wl     *WorldLogicManager
	sh     *SpatialHash
	ps     *PhysicsSystem
	ss     *SteeringSystem
	cs     *CollisionSystem
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
	w.sh = NewSpatialHash(&w, GRID)
	w.ps = NewPhysicsSystem(&w)
	w.ss = NewSteeringSystem(&w)
	w.cs = NewCollisionSystem(&w)
	return &w
}
