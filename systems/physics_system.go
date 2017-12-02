/**
  * 
  * 
  * 
  * 
**/



package systems

import (

	"fmt"
	
	"github.com/dt-rush/donkeys-qquest/engine"
	"github.com/dt-rush/donkeys-qquest/components"
	"github.com/dt-rush/donkeys-qquest/constants"
)


type PhysicsSystem struct {
	// to filter, lookup entities
	entity_manager *engine.EntityManager
	// components this will use
	active_component *components.ActiveComponent
	position_component *components.PositionComponent
	velocity_component *components.VelocityComponent
}


func (s *PhysicsSystem) Init (entity_manager *engine.EntityManager,
	active_component *components.ActiveComponent,
	position_component *components.PositionComponent,
	velocity_component *components.VelocityComponent) {
		
		s.entity_manager = entity_manager
		
		s.active_component = active_component
		s.position_component = position_component
		s.velocity_component = velocity_component	
	}



// TODO TEMP, MOVE TO OWN SYSTEM
func (s *PhysicsSystem) Update (dt_ms float64) {
	// physics
	// TODO
	// separate this into its own system physics_system.go
	// (also checks that you don't leave the screen, but that's just for now)
	for _, id := range s.entity_manager.Entities() {
		// TODO - note that we never preemptively filter entities or
		// query entities or even check each entity, as to whether they
		// have position and velocity, we just assume all do, add this
		// checking using a bitarray component mapper
		if ! s.active_component.Get (id) {
			// don't update inactive entities
			continue
		}
		// apply velocity to position of entities
		pos := s.position_component.Get (id)
		vel := s.velocity_component.Get (id)
		if (id == 0) {
			fmt.Println (vel)
		}
		dx := vel[0] * (dt_ms / 1000.0)
		dy := vel[1] * (dt_ms / 1000.0)

		if pos[0] + dx > 0 && pos[0] + dx < float64 (constants.WIDTH - 20) {
			pos[0] += dx
		}
		if pos[1] + dy > 20 && pos[1] + dy < float64 (constants.HEIGHT) {
			pos[1] += dy
		}
		s.position_component.Set (id, pos)
	}
}
