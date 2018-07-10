package engine

type SteeringSystem struct {
	em               *EntityManager
	movementEntities *UpdatedEntityList
}

func NewSteeringSystem(em *EntityManager) *SteeringSystem {
	query := EntityQueryFromComponentBitArray(
		"steering",
		MakeComponentBitArray([]ComponentType{
			POSITION_COMPONENT,
			VELOCITY_COMPONENT,
			MOVEMENTTARGET_COMPONENT,
			STEER_COMPONENT,
			MASS_COMPONENT,
		}))
	return &SteeringSystem{
		em:               em,
		movementEntities: em.GetUpdatedEntityList(query),
	}
}

func (s *SteeringSystem) Seek(e *EntityToken) {
	p0 := s.em.Components.Position[e.ID]
	p1 := s.em.Components.MovementTarget[e.ID]
	v := &s.em.Components.Velocity[e.ID]
	maxV := s.em.Components.MaxVelocity[e.ID]
	st := &s.em.Components.Steer[e.ID]
	desired := p1.Sub(p0)
	distance := desired.Magnitude()
	desired = desired.Unit()
	// do slowing for arrival behavior
	// TODO: define this properly
	slowingRadius := 30.0
	if distance <= slowingRadius {
		desired = desired.Scale(maxV * distance / slowingRadius)
	} else {
		desired = desired.Scale(maxV)
	}
	force := desired.Sub(*v)
	st.Inc(force)
}

func (s *SteeringSystem) Apply(e *EntityToken) {
	v := &s.em.Components.Velocity[e.ID]
	maxV := s.em.Components.MaxVelocity[e.ID]
	st := &s.em.Components.Steer[e.ID]
	mass := s.em.Components.Mass[e.ID]
	// TODO: define this properly
	maxSteerForce := 3.0
	*st = st.Truncate(maxSteerForce)
	*st = st.Scale(1 / mass)
	*v = v.Add(*st).Truncate(maxV)
}

func (s *SteeringSystem) Update(dt_ms int64) {
	for _, e := range s.movementEntities.Entities {
		s.Seek(e)
		s.Apply(e)
	}
}
