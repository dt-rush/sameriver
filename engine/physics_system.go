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
	WORLD_WIDTH  int
	WORLD_HEIGHT int
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
	s.WORLD_WIDTH = WORLD_WIDTH
	s.WORLD_HEIGHT = WORLD_HEIGHT
}

// apply velocity to position of entities
// NOTE: this is called from Update and is covered by its mutex on the
// components
func (s *PhysicsSystem) applyPhysics(entity EntityToken, dt_ms uint16) {
	// read the position and velocity, using dt to compute dx, dy
	box := s.em.Components.Box[entity.ID]
	vel := s.em.Components.Velocity[entity.ID]
	dx := int32(vel[0] * float32(dt_ms/4))
	dy := int32(vel[1] * float32(dt_ms/4))
	// prevent from leaving the world in X
	if box.X+dx < 0 {
		box.X = 0
	} else if box.X+dx > WORLD_WIDTH-box.W {
		box.X = int16(WORLD_WIDTH) - int16(box.X/2)
	} else {
		box.X += dx
	}
	// prevent from leaving the world in Y
	if box.Y+dy <
		int16(box[1]/2) {
		box.Y = int16(box[1] / 2)
	} else if box.Y+dy >
		int16(WORLD_HEIGHT)-int16(box[1]/2) {
		box.Y = int16(WORLD_HEIGHT) - int16(box[1]/2)
	} else {
		box.Y += dy
	}
	// set the new position which has been computed
	s.em.Components.Position.Data[entity.ID] = pos
}

func (s *PhysicsSystem) Update(dt_ms uint16) {

	s.physicsEntities.Mutex.Lock()
	defer s.physicsEntities.Mutex.Unlock()

	for _, e := range s.physicsEntities.Entities {
		// apply the physics only if this entity isn't already locked
		// (atomic operations are cheap, so this isn't a bad thing to
		// do for each entity during each Update())
		if s.em.attemptLockEntityOnce(e) {
			s.applyPhysics(e, dt_ms)
			s.em.releaseEntity(e)
		}
	}
}
