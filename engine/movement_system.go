package engine

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
		mv := math.Sqrt(v.x*v.x + v.y*v.y) // magnitude of velocity vector
		// our steer factor
		s := s.em.Components.Steer[e.ID]

		// compute vec with same magnitude as v pointing toward p1
		d := Vec2D{p1.x - p0.x, p1.y - p0.y}
		md := math.Sqrt(d.x*d.x + d.y*d.y)
		d.x *= mv / md
		d.y *= mv / md

		*v.x = v.x + d.x*s
		*v.y = v.y + d.y*s
	}
}
