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
	entity_manager *EntityManager
	// targetted entities
	physicsEntities *UpdatedEntityList
}

func (s *PhysicsSystem) Init(entity_manager *EntityManager) {
	// take down a reference to entity manager
	s.entity_manager = entity_manager
	// get a regularly updated list of the entities which have physics
	// (position, velocity and hitbox)
	query := NewBitArraySubsetQuery(
		MakeComponentBitArray([]int{
			POSITION_COMPONENT,
			VELOCITY_COMPONENT,
			HITBOX_COMPONENT}))
	s.physicsEntities = s.entity_manager.GetUpdatedActiveList (query, "physical")
}

// apply velocity to position of entities
// NOTE: this is called from Update and is covered by its mutex on the
// components
func (s *PhysicsSystem) applyPhysics (id uint16, dt_ms uint16) {
	// read the position and velocity, using dt to compute dx, dy
	pos := s.entity_manager.Components.Position.Data[id]
	vel := s.entity_manager.Components.Velocity.Data[id]
	dx := vel[0] * int16(dt_ms)
	dy := vel[1] * int16(dt_ms)
	box := s.entity_manager.Components.Hitbox.Data[id]
	// prevent from leaving the world in X
	if pos[0]+dx <
		int16(box[0]/2) {
		pos[0] = int16(box[0]/2)
	} else if pos[0]+dx >
		int16(constant.WINDOW_WIDTH) - int16(box[0]/2) {
		pos[0] = int16(constant.WINDOW_WIDTH) - int16(box[0]/2)
	} else {
		pos[0] += dx
	}
	// prevent from leaving the world in Y
	if pos[1]+dy <
		int16(box[1]/2) {
		pos[1] = int16(box[1]/2)
	} else if pos[1]+dy >
		int16(constant.WINDOW_HEIGHT) - int16(box[1]/2) {
		pos[1] = int16(constant.WINDOW_HEIGHT) - int16(box[1]/2)
	} else {
		pos[1] += dy
	}
	// set the new position which has been computed
	s.entity_manager.Components.Position.Data[id] = pos
}

func (s *PhysicsSystem) Update(dt_ms uint16) {

	s.entity_manager.Components.Position.Mutex.Lock()
	s.entity_manager.Components.Velocity.Mutex.Lock()
	s.entity_manager.Components.Hitbox.Mutex.Lock()
	s.physicsEntities.Mutex.Lock()
	defer s.entity_manager.Components.Position.Mutex.Unlock()
	defer s.entity_manager.Components.Velocity.Mutex.Unlock()
	defer s.entity_manager.Components.Hitbox.Mutex.Unlock()
	defer s.physicsEntities.Mutex.Unlock()

	for _, id := range s.physicsEntities.Entities {
		s.applyPhysics (id, dt_ms)
	}
}
