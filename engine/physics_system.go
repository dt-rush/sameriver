// moves entities according to their velocity
package engine

import (
	"time"
)

type PhysicsSystem struct {
	w               *World
	physicsEntities *UpdatedEntityList
	lastUpdate      *time.Time
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (s *PhysicsSystem) LinkWorld(w *World) {
	s.w = w
	s.physicsEntities = w.Em.GetUpdatedEntityList(
		EntityQueryFromComponentBitArray(
			"physical",
			MakeComponentBitArray([]ComponentType{
				POSITION_COMPONENT,
				VELOCITY_COMPONENT,
				BOX_COMPONENT,
				MASS_COMPONENT, // TODO: make some use of mass?
			})))
}

func (s *PhysicsSystem) Update() {
	now := time.Now()
	defer func() {
		s.lastUpdate = &now
	}()
	if s.lastUpdate == nil {
		return
	} else {
		dt_ms := float64(time.Since(*s.lastUpdate).Nanoseconds()) / 1.0e6
		// note: there are no function calls in the below, so we won't
		// be preempted while computin physics (this is very good, get it over with)
		for _, e := range s.physicsEntities.entities {
			pos := &s.w.Em.Components.Position[e.ID]
			box := s.w.Em.Components.Box[e.ID]
			vel := s.w.Em.Components.Velocity[e.ID]
			dx := vel.X * dt_ms
			dy := vel.Y * dt_ms
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
			if pos.Y+dy < 0 {
				// max out on the bottom
				pos.Y = 0
			} else if pos.Y+box.Y+dy > float64(s.w.Height) {
				// max out on the top
				pos.Y = float64(s.w.Height) - box.Y
			} else {
				// otherwise move in y freely
				pos.Y += dy
			}
		}
	}
}
