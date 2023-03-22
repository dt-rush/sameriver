package sameriver

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
		"Vec2D,Acceleration", // TODO: use acceleration in steering
		"Float64,MaxVelocity",
		"Vec2D,MovementTarget",
		"Vec2D,Steer",
		"Float64,Mass",
	}
}

func (s *SteeringSystem) LinkWorld(w *World) {
	s.w = w
	s.movementEntities = w.GetUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"steering",
			w.em.components.BitArrayFromIDs(
				[]ComponentID{
					POSITION, VELOCITY, ACCELERATION,
					MAXVELOCITY, MOVEMENTTARGET, STEER, MASS,
				})))
}

func (s *SteeringSystem) Update(dt_ms float64) {
	for _, e := range s.movementEntities.entities {
		s.Seek(e)
		s.Apply(e)
	}
}

func (s *SteeringSystem) Seek(e *Entity) {
	p0 := e.GetVec2D(POSITION)
	p1 := e.GetVec2D(MOVEMENTTARGET)
	v := e.GetVec2D(VELOCITY)
	maxV := e.GetFloat64(MAXVELOCITY)
	st := e.GetVec2D(STEER)
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
	v := e.GetVec2D(VELOCITY)
	maxV := e.GetFloat64(MAXVELOCITY)
	st := e.GetVec2D(STEER)
	mass := e.GetFloat64(MASS)
	// TODO: define this properly
	maxSteerForce := 3.0
	*st = st.Truncate(maxSteerForce)
	*st = st.Scale(1 / *mass)
	*v = v.Add(*st).Truncate(*maxV)
}

func (s *SteeringSystem) Expand(n int) {
	// nil?
}
