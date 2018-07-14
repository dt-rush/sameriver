package engine

type SteeringSystem struct {
	w                *World
	movementEntities *UpdatedEntityList
}

func NewSteeringSystem() *SteeringSystem {
	return &SteeringSystem{}
}

func (s *SteeringSystem) LinkWorld(w *World) {
	s.w = w
	s.movementEntities = w.Em.GetUpdatedEntityList(
		EntityQueryFromComponentBitArray(
			"steering",
			MakeComponentBitArray([]ComponentType{
				POSITION_COMPONENT,
				VELOCITY_COMPONENT,
				MOVEMENTTARGET_COMPONENT,
				STEER_COMPONENT,
				MASS_COMPONENT,
			})))
}

func (s *SteeringSystem) Update(dt_ms float64) {
	for _, e := range s.movementEntities.entities {
		s.Seek(e)
		s.Apply(e)
	}
}

func (s *SteeringSystem) Seek(e *EntityToken) {
	p0 := s.w.Em.Components.Position[e.ID]
	p1 := s.w.Em.Components.MovementTarget[e.ID]
	v := &s.w.Em.Components.Velocity[e.ID]
	maxV := s.w.Em.Components.MaxVelocity[e.ID]
	st := &s.w.Em.Components.Steer[e.ID]
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
	v := &s.w.Em.Components.Velocity[e.ID]
	maxV := s.w.Em.Components.MaxVelocity[e.ID]
	st := &s.w.Em.Components.Steer[e.ID]
	mass := s.w.Em.Components.Mass[e.ID]
	// TODO: define this properly
	maxSteerForce := 3.0
	*st = st.Truncate(maxSteerForce)
	*st = st.Scale(1 / mass)
	*v = v.Add(*st).Truncate(maxV)
}
