package engine

type World struct {
	Width  int
	Height int
	ev     *EventBus
	em     *EntityManager
	wl     *WorldLogicManager
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
	w.ps = NewPhysicsSystem(&w)
	w.ss = NewSteeringSystem(&w)
	w.cs = NewCollisionSystem(&w)
	return &w
}

func (w *World) EntityIsWithinRect(
	e *EntityToken, topleft Vec2D, dimension Vec2D) bool {

	if !w.em.EntityHasComponent(e, POSITION_COMPONENT) ||
		!w.em.EntityHasComponent(e, BOX_COMPONENT) {
		return false
	}
	// NOTE: position is the top-left of the entity's box
	pos := w.em.Components.Position[e.ID]
	box := w.em.Components.Box[e.ID]
	return pos.X >= topleft.X &&
		pos.X+box.X <= topleft.X+dimension.X &&
		pos.Y >= topleft.Y+box.Y &&
		pos.Y <= topleft.Y+dimension.Y
}
