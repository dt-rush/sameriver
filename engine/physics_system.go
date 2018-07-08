/**
  *
  *
  *
  *
**/

package engine

type PhysicsSystem struct {
	// to filter, lookup entities
	em *EntityManager
	// targetted entities
	physicsEntities *UpdatedEntityList
	// world dimensions
	WORLD_WIDTH  int32
	WORLD_HEIGHT int32
}

func (s *PhysicsSystem) Init(
	WORLD_WIDTH int,
	WORLD_HEIGHT int,
	em *EntityManager) {
	// take down a reference to entity manager
	s.em = em
	// get a regularly updated list of the entities which have physics
	// (position, velocity and hitbox)
	query := EntityQueryFromComponentBitArray(
		"physical",
		MakeComponentBitArray([]ComponentType{
			BOX_COMPONENT,
			VELOCITY_COMPONENT}))
	s.physicsEntities = s.em.GetUpdatedEntityList(query)
	// set world dimensions
	s.WORLD_WIDTH = int32(WORLD_WIDTH)
	s.WORLD_HEIGHT = int32(WORLD_HEIGHT)
}

func (s *PhysicsSystem) Update(dt_ms int64) {

	// note: there are no function calls in the below, so we won't
	// be preempted while computin physics (this is very good, get it over with)
	for _, e := range s.physicsEntities.Entities {
		// read the position and velocity, using dt to compute dx, dy
		pos := s.em.ComponentsData.Position[e.ID]
		box := s.em.ComponentsData.Box[e.ID]
		vel := s.em.ComponentsData.Velocity[e.ID]
		dx := vel.X * float64(dt_ms)
		dy := vel.Y * float64(dt_ms)
		// motion in x
		if pos.X+dx < 0 {
			// max out on the left
			pos.X = 0
		} else if pos.X+dx > float64(s.WORLD_WIDTH-box.W) {
			// max out on the right
			pos.X = float64(s.WORLD_WIDTH - box.W)
		} else {
			// otherwise move in x freely
			pos.X += dx
		}
		// motion in y
		if pos.Y+dy < float64(box.H) {
			// max out on the bottom
			pos.Y = 0
		} else if pos.Y+dy > float64(s.WORLD_HEIGHT) {
			// max out on the top
			pos.Y = float64(s.WORLD_HEIGHT)
		} else {
			// otherwise move in y freely
			pos.Y += dy
		}
	}
}
