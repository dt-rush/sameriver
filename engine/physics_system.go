/**
  *
  *
  *
  *
**/

package engine

type PhysicsSystem struct {
	w               *World
	physicsEntities *UpdatedEntityList
}

func NewPhysicsSystem(w *World) *PhysicsSystem {
	// get a regularly updated list of the entities which have physics
	query := EntityQueryFromComponentBitArray(
		"physical",
		MakeComponentBitArray([]ComponentType{
			POSITION_COMPONENT,
			VELOCITY_COMPONENT,
			BOX_COMPONENT,
			// TODO: make use of mass?
			MASS_COMPONENT,
		}))
	return &PhysicsSystem{
		w:               w,
		physicsEntities: w.em.GetUpdatedEntityList(query),
	}
}

func (s *PhysicsSystem) Update(dt_ms int64) {

	// note: there are no function calls in the below, so we won't
	// be preempted while computin physics (this is very good, get it over with)
	for _, e := range s.physicsEntities.Entities {
		// read the position and velocity, using dt to compute dx, dy
		pos := &s.w.em.Components.Position[e.ID]
		box := s.w.em.Components.Box[e.ID]
		vel := s.w.em.Components.Velocity[e.ID]
		dx := vel.X * float64(dt_ms)
		dy := vel.Y * float64(dt_ms)
		// motion in x
		if pos.X+dx < 0 {
			// max out on the left
			pos.X = 0
		} else if pos.X+box.X+dx > float64(s.w.Width) {
			// max out on the right
			pos.X = float64(s.w.Width) - box.X
		} else {
			// otherwise move in x freely
			pos.X += dx
		}
		// motion in y
		if pos.Y+dy < box.Y {
			// max out on the bottom
			pos.Y = box.Y
		} else if pos.Y+dy > float64(s.w.Height) {
			// max out on the top
			pos.Y = float64(s.w.Height)
		} else {
			// otherwise move in y freely
			pos.Y += dy
		}
	}
}
