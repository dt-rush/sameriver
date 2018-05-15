/**
  *
  *
  *
  *
**/

package engine

type CollisionSystem struct {
	// Reference to entity manager to reach components
	entity_manager *EntityManager
	// targetted entities
	collidableEntities *UpdatedEntityList
	// How the collision system communicates collision events
	game_event_manager *GameEventManager
	// How the collision system gets populated with specific
	// collision detection logics
	collision_logic_collection    map[uint16]CollisionLogic
	collision_logic_ids           map[string]uint16
	collision_logic_active_states map[uint16]bool
	// to generate IDs for collision logic
	id_generator IDGenerator
}

func (s *CollisionSystem) Init(
	entity_manager *EntityManager,
	game_event_manager *GameEventManager) {

	// take down references to entity_manager and game_event_manager
	s.entity_manager = entity_manager
	s.game_event_manager = game_event_manager
	// get a regularly updated list of the entities which are collidable
	// (position and hitbox)
	query := NewBitArraySubsetQuery(
		MakeComponentBitArray([]int{
			POSITION_COMPONENT,
			HITBOX_COMPONENT}))
	s.collidableEntities = s.entity_manager.GetUpdatedActiveList(query, "collidable")
	// initialize collision logic data members
	s.collision_logic_collection = make(map[uint16]CollisionLogic)
	s.collision_logic_ids = make(map[string]uint16)
	s.collision_logic_active_states = make(map[uint16]bool)
}

func (s *CollisionSystem) AddCollisionLogic(name string, logic CollisionLogic) uint16 {

	id := s.id_generator.Gen()
	Logger.Printf("about to add collision logic %s", name)
	s.collision_logic_collection[id] = logic
	s.collision_logic_ids[name] = id
	Logger.Printf("added collision logic %s", name)
	return id
}

func (s *CollisionSystem) SetCollisionLogicActiveState(id uint16, active bool) {
	s.collision_logic_active_states[id] = active
}

// Test collision between two functions
// NOTE: this is called by Update, so it's covered by the mutex on the
// components
func (s *CollisionSystem) TestCollision(i uint16, j uint16) bool {
	// grab component data
	position_component := s.entity_manager.Components.Position
	hitbox_component := s.entity_manager.Components.Hitbox
	box := hitbox_component.Data[i]
	other_box := hitbox_component.Data[j]
	center := position_component.Data[i]
	other_center := position_component.Data[j]
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

func (s *CollisionSystem) Update(dt_ms uint16) {

	s.entity_manager.Components.Position.Mutex.Lock()
	s.entity_manager.Components.Hitbox.Mutex.Lock()
	s.collidableEntities.Mutex.Lock()
	defer s.entity_manager.Components.Position.Mutex.Unlock()
	defer s.entity_manager.Components.Hitbox.Mutex.Unlock()
	defer s.collidableEntities.Mutex.Unlock()

	entities := s.collidableEntities.Entities
	for i := uint16(0); i < uint16(len(entities)); i++ {

		// compare entity at i to all subsequent entities
		// (this way, all entity pairs will be compared once)
		for j := i + 1; j < uint16(len(entities)); j++ {

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
