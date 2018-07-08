package engine

import (
	"math"
)

type MovementSystem struct {
	em               *EntityManager
	movementEntities *UpdatedEntityList
}

func (s *MovementSystem) Init(
	em *EntityManager) {

	s.em = em
	query := EntityQueryFromComponentBitArray(
		"moving",
		MakeComponentBitArray([]ComponentType{
			BOX_COMPONENT,
			VELOCITY_COMPONENT,
			MOVEMENTTARGET_COMPONENT}))
	s.movementEntities = s.em.GetUpdatedEntityList(query)
}

func (s *MovementSystem) Update() {
	for _, e := range s.movementEntities.Entities {
		// our position
		p0 := s.em.Components.Position[e.ID]
		// our target
		p1 := s.em.Components.MovementTarget[e.ID]
		v := &s.em.Components.Velocity[e.ID]
		// magnitude of velocity vector
		mv := math.Sqrt(v.X*v.X + v.Y*v.Y)
		// our steer factor
		s := s.em.Components.Steer[e.ID]

		// compute vec with same magnitude as v pointing toward p1
		d := Vec2D{p1.X - p0.X, p1.Y - p0.X}
		md := math.Sqrt(d.X*d.X + d.Y*d.Y)
		d.X *= (mv / md)
		d.Y *= (mv / md)

		v.X = v.X + d.X*s
		v.Y = v.Y + d.Y*s
	}
}
