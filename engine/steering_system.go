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
	s.movementEntities = w.em.GetUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"steering",
			MakeComponentBitArray([]ComponentType{
				POSITION_COMPONENT,
				VELOCITY_COMPONENT,
				MOVEMENTTARGET_COMPONENT,
				STEER_COMPONENT,
				MASS_COMPONENT,
			})))
}

func (s *SteeringSystem) Update() {
	for _, e := range s.movementEntities.entities {
		s.Seek(e)
		s.Apply(e)
	}
}

func (s *SteeringSystem) Seek(e *Entity) {
	p0 := s.w.em.components.Position[e.ID]
	p1 := s.w.em.components.MovementTarget[e.ID]
	v := &s.w.em.components.Velocity[e.ID]
	maxV := s.w.em.components.MaxVelocity[e.ID]
	st := &s.w.em.components.Steer[e.ID]
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

func (s *SteeringSystem) Apply(e *Entity) {
	v := &s.w.em.components.Velocity[e.ID]
	maxV := s.w.em.components.MaxVelocity[e.ID]
	st := &s.w.em.components.Steer[e.ID]
	mass := s.w.em.components.Mass[e.ID]
	// TODO: define this properly
	maxSteerForce := 3.0
	*st = st.Truncate(maxSteerForce)
	*st = st.Scale(1 / mass)
	*v = v.Add(*st).Truncate(maxV)
}
