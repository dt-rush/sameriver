/**
  *
  *
  *
  *
**/

package engine

import (
	"github.com/dt-rush/donkeys-qquest/engine"
)

type CollisionLogic struct {

	// NOTE: in the below, i and j are entity ID's

	Selector func(i int,
		j int,
		em *engine.EntityManager) bool

	EventGenerator func(i int,
		j int,
		em *engine.EntityManager) engine.GameEvent
}

type CollisionSystem struct {

	// To filter, lookup entities
	entity_manager *engine.EntityManager
	// query watcher to see spawn/despawn events and update targetted entities
	spawnWatcher QueryWatcher
	// targetted entities
	CollidableEntities []int
	// How the collision system communicates collision events
	game_event_manager *engine.GameEventManager
	// How the collision system gets populated with specific
	// collision detection logics
	collision_logic_collection    map[int]CollisionLogic
	collision_logic_ids           map[string]int
	collision_logic_active_states map[int]bool
	// to generate IDs for collision logic
	id_generator engine.IDGenerator
}

func (s *CollisionSystem) Init(
	entity_manager *engine.EntityManager,
	game_event_manager *engine.GameEventManager) {

	s.entity_manager = entity_manager
	s.game_event_manager = game_event_manager
	s.setSpawnWatcher()

	s.collision_logic_collection = make(map[int]CollisionLogic)
	s.collision_logic_ids = make(map[string]int)
	s.collision_logic_active_states = make(map[int]bool)
	s.CollidableEntities = make([]int)
}

func (s *CollisionSystem) setSpawnWatcher() {
	query := engine.MakeComponentQuery(
		engine.POSITION_COMPONENT,
		engine.HITBOX_COMPONENT)
}

func (s *CollisionSystem) AddCollisionLogic(name string, logic CollisionLogic) int {

	id := s.id_generator.Gen()
	engine.Logger.Printf("about to add collision logic %s", name)
	s.collision_logic_collection[id] = logic
	s.collision_logic_ids[name] = id
	engine.Logger.Printf("added collision logic %s", name)
	return id
}

func (s *CollisionSystem) SetCollisionLogicActiveState(id int, active bool) {
	s.collision_logic_active_states[id] = active
}

func (s *CollisionSystem) TestCollision(i int, j int) bool {
	// grab component data
	box := s.hitbox_component.Get(i)
	other_box := s.hitbox_component.Get(j)
	center := s.position_component.Get(i)
	other_center := s.position_component.Get(j)
	// find the distance between the X and Y centers
	// NOTE: "abs" is for absolute value
	dxabs := center[0] - other_center[0]
	if dxabs < 0 {
		dxabs *= -1
	}
	dyabs := center[1] - other_center[1]
	if dyabs < 0 {
		dyabs *= -1
	}
	// if the sum of the widths is greater than twice the x distance,
	// collision has occurred (same for y)
	return (dxabs*2 < (box[0]+other_box[0]) &&
		dyabs*2 < (box[1]+other_box[1]))
}

func (s *CollisionSystem) Update(dt_ms int) {

	for i := 0; i < len(c.CollidableEntities); i++ {
		entity_i := entities[i]
		if !s.EntityIsCollidable(entity_i) {
			continue
		}
		// compare entity at i to all subsequent entities
		// (this way, all entity pairs will be compared once)
		for j := i + 1; j < len(entities); j++ {
			entity_j := entities[j]
			if !s.EntityIsCollidable(entity_j) {
				continue
			}

			for collision_logic_id, collision_logic := range s.collision_logic_collection {
				// if this collision logic is active,
				// and the entities i and j match the selector,
				// and there is a collision,
				//    then emit an event according to the eventgenerator
				if s.collision_logic_active_states[collision_logic_id] &&
					collision_logic.Selector(i, j, s.entity_manager) &&
					s.TestCollision(i, j) {
					event_generated := collision_logic.EventGenerator(i, j, s.entity_manager)
					s.game_event_manager.Publish(event_generated)
				}
			}
		}
	}
}
