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
		box := s.em.Components.Box[e.ID]
		vel := s.em.Components.Velocity[e.ID]
		dx := int32(vel.X * float32(dt_ms/4))
		dy := int32(vel.Y * float32(dt_ms/4))
		// motion in x
		if box.X+dx < 0 {
			// max out on the left
			box.X = 0
		} else if box.X+dx > s.WORLD_WIDTH-box.W {
			// max out on the right
			box.X = s.WORLD_WIDTH - box.W
		} else {
			// otherwise move in x freely
			box.X += dx
		}
		// motion in y
		if box.Y+dy < box.H {
			// max out on the bottom
			box.Y = 0
		} else if box.Y+dy > s.WORLD_HEIGHT {
			// max out on the top
			box.Y = s.WORLD_HEIGHT
		} else {
			// otherwise move in y freely
			box.Y += dy
		}
	}
}
