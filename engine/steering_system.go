package engine

type SteeringSystem struct {
	w                *World
	movementEntities *UpdatedEntityList
}

func NewSteeringSystem(w *World) *SteeringSystem {
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
		w:                w,
		movementEntities: w.em.GetUpdatedEntityList(query),
	}
}

func (s *SteeringSystem) Seek(e *EntityToken) {
	p0 := s.w.em.Components.Position[e.ID]
	p1 := s.w.em.Components.MovementTarget[e.ID]
	v := &s.w.em.Components.Velocity[e.ID]
	maxV := s.w.em.Components.MaxVelocity[e.ID]
	st := &s.w.em.Components.Steer[e.ID]
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
	v := &s.w.em.Components.Velocity[e.ID]
	maxV := s.w.em.Components.MaxVelocity[e.ID]
	st := &s.w.em.Components.Steer[e.ID]
	mass := s.w.em.Components.Mass[e.ID]
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
