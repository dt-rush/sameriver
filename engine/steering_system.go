package engine

type SteeringSystem struct {
	w                *World
	movementEntities *UpdatedEntityList
}

func NewSteeringSystem() *SteeringSystem {
	return &SteeringSystem{}
}

func (s *SteeringSystem) GetComponentDeps() []string {
	return []string{
		"Vec2D,Position",
		"Vec2D,Velocity",
		"Float64,MaxVelocity",
		"Vec2D,MovementTarget",
		"Vec2D,Steer",
		"Float64,Mass",
	}
}

func (s *SteeringSystem) LinkWorld(w *World) {
	s.w = w
	s.movementEntities = w.em.GetUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"steering",
			w.em.components.BitArrayFromNames(
				[]string{
					"Position",
					"Velocity",
					"MaxVelocity",
					"MovementTarget",
					"Steer",
					"Mass",
				})))
}

func (s *SteeringSystem) Update(dt_ms float64) {
	for _, e := range s.movementEntities.entities {
		s.Seek(e)
		s.Apply(e)
	}
}

func (s *SteeringSystem) Seek(e *Entity) {
	p0 := e.GetVec2D("Position")
	p1 := e.GetVec2D("MovementTarget")
	v := e.GetVec2D("Velocity")
	maxV := e.GetFloat64("MaxVelocity")
	st := e.GetVec2D("Steer")
	desired := p1.Sub(*p0)
	distance := desired.Magnitude()
	desired = desired.Unit()
	// do slowing for arrival behavior
	// TODO: define this properly
	slowingRadius := 30.0
	if distance <= slowingRadius {
		desired = desired.Scale(*maxV * distance / slowingRadius)
	} else {
		desired = desired.Scale(*maxV)
	}
	force := desired.Sub(*v)
	st.Inc(force)
}

func (s *SteeringSystem) Apply(e *Entity) {
	v := e.GetVec2D("Velocity")
	maxV := e.GetFloat64("MaxVelocity")
	st := e.GetVec2D("Steer")
	mass := e.GetFloat64("Mass")
	// TODO: define this properly
	maxSteerForce := 3.0
	*st = st.Truncate(maxSteerForce)
	*st = st.Scale(1 / *mass)
	*v = v.Add(*st).Truncate(*maxV)
}
