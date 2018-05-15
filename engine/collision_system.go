/**
  *
  *
  *
  *
**/

package engine


import (
	"sync"
)


type CollisionLogic struct {

	// NOTE: in the below, i and j are entity ID's

	Selector func(i int,
		j int,
		em *EntityManager) bool

	EventGenerator func(i int,
		j int,
		em *EntityManager) GameEvent
}

type CollisionSystem struct {

	// To filter, lookup entities
	entity_manager *EntityManager
	position_component *PositionComponent
	hitbox_component *HitboxComponent
	// query watcher to see activate/deactivate events from entity_manager and
	// update targetted entities list
	activeWatcher QueryWatcher
	// targetted entities
	collidableEntities UpdatedEntityList
	// How the collision system communicates collision events
	game_event_manager *GameEventManager
	// How the collision system gets populated with specific
	// collision detection logics
	collision_logic_collection    map[int]CollisionLogic
	collision_logic_ids           map[string]int
	collision_logic_active_states map[int]bool
	// to generate IDs for collision logic
	id_generator IDGenerator
}

func (s *CollisionSystem) Init(
	entity_manager *EntityManager,
	game_event_manager *GameEventManager) {

	// take down references to entity and game event managers
	s.entity_manager = entity_manager
	s.game_event_manager = game_event_manager
	// take down reference to components needed
	s.position_component = s.entity_manager.Components.Position
	s.hitbox_component = s.entity_manager.Components.Hitbox
	// get a regularly updated list of the entities which are collidable
	// (position and hitbox)
	query := MakeComponentQuery([]int{
		POSITION_COMPONENT,
		HITBOX_COMPONENT})
	s.collidableEntities = s.entity_manager.GetUpdatedActiveList (query)
	// initialize collision logic data members
	s.collision_logic_collection = make(map[int]CollisionLogic)
	s.collision_logic_ids = make(map[string]int)
	s.collision_logic_active_states = make(map[int]bool)
}

func (s *CollisionSystem) AddCollisionLogic(name string, logic CollisionLogic) int {

	id := s.id_generator.Gen()
	Logger.Printf("about to add collision logic %s", name)
	s.collision_logic_collection[id] = logic
	s.collision_logic_ids[name] = id
	Logger.Printf("added collision logic %s", name)
	return id
}

func (s *CollisionSystem) SetCollisionLogicActiveState(id int, active bool) {
	s.collision_logic_active_states[id] = active
}

func (s *CollisionSystem) TestCollision(i int, j int) bool {
	// NOTE: this is called by Update, so it's covered by the mutex on the
	// components

	// grab component data
	box := s.hitbox_component.Data[i]
	other_box := s.hitbox_component.Data[j]
	center := s.position_component.Data[i]
	other_center := s.position_component.Data[j]
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
	return (dxabs*2 < int16(box[0]+other_box[0]) &&
		dyabs*2 < int16(box[1]+other_box[1]))
}

func (s *CollisionSystem) Update(dt_ms int) {

	s.position_component.Mutex.Lock()
	s.hitbox_component.Mutex.Lock()
	s.collidableEntities.Mutex.Lock()
	defer s.position_component.Mutex.Unlock()
	defer s.hitbox_component.Mutex.Unlock()
	defer s.collidableEntities.Mutex.Unlock()

	entities := s.collidableEntities.Entities
	for i := 0; i < len(entities); i++ {
		entity_i := entities[i]

		// compare entity at i to all subsequent entities
		// (this way, all entity pairs will be compared once)
		for j := i + 1; j < len(entities); j++ {
			entity_j := entities[j]

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
