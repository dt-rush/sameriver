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
		p0 := s.em.Components.Box[e.ID]
		// our target
		p1 := s.em.Components.MovementTarget[e.ID]
		v := &s.em.Components.Velocity[e.ID]
		mv := math.Sqrt(float64(v.X*v.X + v.Y*v.Y)) // magnitude of velocity vector
		// our steer factor
		s := s.em.Components.Steer[e.ID]

		// compute vec with same magnitude as v pointing toward p1
		d := Vec2D{p1.X - float32(p0.X), p1.Y - float32(p0.X)}
		md := math.Sqrt(float64(d.X*d.X + d.Y*d.Y))
		d.X *= float32(mv / md)
		d.Y *= float32(mv / md)

		v.X = v.X + d.X*s
		v.Y = v.Y + d.Y*s
	}
}
