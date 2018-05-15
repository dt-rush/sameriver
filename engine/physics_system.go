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
	// component this will use
	active_component   *ActiveComponent
	position_component *PositionComponent
	velocity_component *VelocityComponent
	hitbox_component   *HitboxComponent
}

func (s *PhysicsSystem) Init(entity_manager *EntityManager,
	active_component *ActiveComponent,
	position_component *PositionComponent,
	velocity_component *VelocityComponent,
	hitbox_component *HitboxComponent) {

	s.entity_manager = entity_manager
	s.active_component = active_component
	s.position_component = position_component
	s.velocity_component = velocity_component
	s.hitbox_component = hitbox_component
}

func (s *PhysicsSystem) Update(dt_ms int) {
	for _, id := range s.entity_manager.Entities() {
		// TODO - note that we never preemptively filter entities or
		// query entities or even check each entity, as to whether they
		// have position and velocity, we just assume all do, add this
		// checking using a bitarray component mapper

		// TODO: also consider a way of defining tags which apply
		// based automatically on whether an entity has a set of
		// component, so we can retrieve a list of all entities (ID's)
		// which have position and velocity using a certain name, like
		// "has_physics"
		if !s.active_component.Get(id) {
			// don't update inactive entities
			continue
		}
		// apply velocity to position of entities
		pos := s.position_component.Get(id)
		vel := s.velocity_component.Get(id)

		dx := vel[0] * (float64(dt_ms) / 1000.0)
		dy := vel[1] * (float64(dt_ms) / 1000.0)

		box := s.hitbox_component.Get(id)

		if pos[0]+dx <
			box[0]/2 {
			pos[0] = box[0] / 2
		} else if pos[0]+dx >
			float64(constant.WINDOW_WIDTH)-box[0]/2 {
			pos[0] = float64(constant.WINDOW_WIDTH) - box[0]/2
		} else {
			pos[0] += dx
		}

		if pos[1]+dy <
			box[1]/2 {
			pos[1] = box[1] / 2
		} else if pos[1]+dy >
			float64(constant.WINDOW_HEIGHT)-box[1]/2 {
			pos[1] = float64(constant.WINDOW_HEIGHT) - box[1]/2
		} else {
			pos[1] += dy
		}
		s.position_component.Set(id, pos)
	}
}
