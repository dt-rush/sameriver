/**
  * 
  * 
  * 
  * 
**/



package systems


import (
	"fmt"
	"time"
	"math/rand"

	"github.com/dt-rush/donkeys-qquest/engine"
	"github.com/dt-rush/donkeys-qquest/components"
	"github.com/dt-rush/donkeys-qquest/constants"
)


type CollisionSystem struct {
	// To filter, lookup entities
	entity_manager *engine.EntityManager
	// Components this will use
	active_component *components.ActiveComponent
	position_component *components.PositionComponent
	hitbox_component *components.HitboxComponent
	// How the collision system communicates to the game
	game_event_system *engine.GameEventSystem
}

func (s *CollisionSystem) Init (entity_manager *engine.EntityManager,
	active_component *components.ActiveComponent,
	position_component *components.PositionComponent,
	hitbox_component *components.HitboxComponent,
	game_event_system *engine.GameEventSystem) {
		
		s.entity_manager = entity_manager
		
		s.active_component = active_component
		s.position_component = position_component
		s.hitbox_component = hitbox_component
		
		s.game_event_system = game_event_system 
	}

func (s *CollisionSystem) TestCollision (i int, j int) bool {
	box := s.hitbox_component.Get (i)
	other_box := s.hitbox_component.Get (j)
	
	center := s.position_component.Get (i)
	center[0] += box[0]/2
	center[1] -= box[1]/2
	other_center := s.position_component.Get (j)
	other_center[0] += other_box[0]/2
	other_center[1] -= other_box[1]/2
	
	collision := center[0] < other_center[0] + other_box[0] &&
		center[0] + box[0] > other_center[0] &&
		center[1] < other_center[1] + other_box[1] &&
		box[1] + center[1] > other_center[1]

	return collision
}

func (s *CollisionSystem) Update (dt_ms float64) {


	// NOTA BENE: we take [0] assuming unique tag
	player_id := s.entity_manager.GetTagEntities ("player") [0]
	donkey_id := s.entity_manager.GetTagEntities ("donkey") [0]
	
	for i_index, i := range s.entity_manager.Entities() {
		// TODO factor out the "get all active components
		// with hitbox and position"
		// logic like this with a usage of the tag system or
		// of the entity-to-component one-to-many bitarray system
		collidable := s.position_component.Has (i) &&
			s.hitbox_component.Has (i) &&
			s.active_component.Has (i)

		if ! s.active_component.Get (i) {
			// inactive entities can't collide with anything
			return
		}

		// note we're looping a second time through the system-allocated IDs
		// to check against every other box (starting from all those after this one,
		// handshake-theorem -style)
		for j_index, j := range s.entity_manager.Entities() {
			// we want all entities *after* this one!
			// so fail-fast through any checks for entities
			// where j_index < i_index
			if j_index < i_index { continue }

			// test if other collidable
			other_collidable := s.position_component.Has (j) &&
				s.hitbox_component.Has (j) &&
				s.active_component.Has (j)

			if ! s.active_component.Get (j) {
				// this entity is inactive, continue early
				continue
			}

			// actual collision rectangle logic (assuming axis-aligned)
			if collidable && other_collidable {


				// TODO refactor
				// check donkey-player collision
				// NOTE: we have to check both orders because we don't know
				// who will be i or j in the handshake as ID's are added to a bag of ID's which
				// may only come out in a given order by coincidence assuring that, for example,
				// the player were always i and the donkey j, never reaching the donkey first
				// via i to compare collisions with a player on j
				selector := (i == player_id && j == donkey_id ||
					i == donkey_id && j == player_id)
				
				if selector && s.TestCollision (i, j) {
					
					s.game_event_system.Publish (constants.GAME_EVENT_DONKEY_CAUGHT)

					fmt.Println ("GOT A DONKEY!")
					
					// TODO: also simultaneously set invisible?
					// TODO (possibly): separate "visible" component from "active"?
					s.active_component.Set (donkey_id, false)
					// sleep 5 seconds
					go func() {
						time.Sleep (time.Second * 5)
						// respawn the donkey
						donkey_pos := s.position_component.Get (donkey_id)
						donkey_pos [0] = rand.Float64() * float64 (constants.WIDTH - 20) + 20
						donkey_pos [1] = rand.Float64() * float64 (constants.HEIGHT - 20) + 20
						s.active_component.Set (donkey_id, true)
					}()
				}

				// TODO refactor
				
				// check flame-player collision
			
				flame_ids := s.entity_manager.GetTagEntities ("flame")
				j_is_flame := false
				i_is_flame := false
				for _, flame_id := range (flame_ids) {
					if j == flame_id {
						j_is_flame = true
					}
					if i == flame_id {
						i_is_flame = true
					}
				}

				selector = i == player_id && j_is_flame ||
					j == player_id && i_is_flame

				if selector && s.TestCollision (i, j) {

					if constants.DEBUG_COLLISION {
						fmt.Println (s.position_component.Get (i), s.position_component.Get (j))
						fmt.Println (s.hitbox_component.Get (i), s.hitbox_component.Get (j))
					}
					
					s.game_event_system.Publish (constants.GAME_EVENT_FLAME_HIT_PLAYER)
				}
			}
		}
	}	
}

