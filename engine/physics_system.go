/**
  *
  *
  *
  *
**/

package engine

import (
	"github.com/dt-rush/donkeys-qquest/constant"
)

type PhysicsSystem struct {
	// to filter, lookup entities
	em *EntityManager
	// targetted entities
	physicsEntities *UpdatedEntityList
}

func (s *PhysicsSystem) Init(em *EntityManager) {
	// take down a reference to entity manager
	s.em = em
	// get a regularly updated list of the entities which have physics
	// (position, velocity and hitbox)
	query := NewBitArraySubsetQuery(
		MakeComponentBitArray([]int{
			POSITION_COMPONENT,
			VELOCITY_COMPONENT,
			HITBOX_COMPONENT}))
	s.physicsEntities = s.em.GetUpdatedActiveEntityList(query, "physical")
}

// apply velocity to position of entities
// NOTE: this is called from Update and is covered by its mutex on the
// components
func (s *PhysicsSystem) applyPhysics(id uint16, dt_ms uint16) {
	// read the position and velocity, using dt to compute dx, dy
	pos := s.em.Components.Position.Data[id]
	vel := s.em.Components.Velocity.Data[id]
	dx := int16(vel[0] * float32(dt_ms/4))
	dy := int16(vel[1] * float32(dt_ms/4))
	box := s.em.Components.Hitbox.Data[id]
	// prevent from leaving the world in X
	if pos[0]+dx <
		int16(box[0]/2) {
		pos[0] = int16(box[0] / 2)
	} else if pos[0]+dx >
		int16(constant.WINDOW_WIDTH)-int16(box[0]/2) {
		pos[0] = int16(constant.WINDOW_WIDTH) - int16(box[0]/2)
	} else {
		pos[0] += dx
	}
	// prevent from leaving the world in Y
	if pos[1]+dy <
		int16(box[1]/2) {
		pos[1] = int16(box[1] / 2)
	} else if pos[1]+dy >
		int16(constant.WINDOW_HEIGHT)-int16(box[1]/2) {
		pos[1] = int16(constant.WINDOW_HEIGHT) - int16(box[1]/2)
	} else {
		pos[1] += dy
	}
	// set the new position which has been computed
	s.em.Components.Position.Data[id] = pos
}

func (s *PhysicsSystem) Update(dt_ms uint16) {

	s.em.Components.Position.Mutex.Lock()
	s.em.Components.Velocity.Mutex.Lock()
	s.em.Components.Hitbox.Mutex.Lock()
	s.physicsEntities.Mutex.Lock()
	defer s.em.Components.Position.Mutex.Unlock()
	defer s.em.Components.Velocity.Mutex.Unlock()
	defer s.em.Components.Hitbox.Mutex.Unlock()
	defer s.physicsEntities.Mutex.Unlock()

	for _, e := range s.physicsEntities.Entities {
		// apply the physics only if this entity is not held for modification
		// (atomic operations are cheap, so this isn't a bad thing to
		// do for each entity during each Update())
		if s.em.holdEntity(uint16(e.ID)) {
			s.applyPhysics(uint16(e.ID), dt_ms)
			s.em.releaseEntity(uint16(e.ID))
		}
	}
}
