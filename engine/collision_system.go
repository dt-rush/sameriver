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
	s.collidableEntities = s.entity_manager.GetUpdatedActiveList(
		query, "collidable")
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
	// The way we determine the key to the rate limiters map is a little
	// funky. We first start by comparing entities in a handshake pattern,
	// where we compare 0 with every index after 0, 1, with every index
	// after 1, etc. We then take the ID's at these indexes and put them
	// in a sorted order. The lesser one will be, as a uint16, shifted 16
	// and OR'd with the greater. These are the map keys for a rate limiter
	// corresponding uniquely to the collision between the two ID's.
	// The rate limiter 

	// TODO: make this make sense lol
	// idea: maybe a MAX_ENTITIES x MAX_ENTITIES square array, and when
	// an entity is deactivated, reset the timers on its row/ column?
	// do we even *need* a rate limiter? we want to check collision as often
	// as possible to be accurate, that is, we want to approach the rate
	// of the physics loop (how far can an entity move in each
	// physics loop? this determines how likely a miss is depending on the
	// rate we choose for checking collisions), *but*, we don't want to
	// *spawn events* at that same rate. we want to say, great, we caught a
	// collision, now acting on it will likely take some time, so let's cool
	// off for a bit on this collision, maybe 100ms? 200?
	for ix, id := range entities {
		for jx := ix + 1; jx < uint16(len(entities)); jx++ {
			jd := entities[jx]
			if id < jd

	for i := uint16(0); i < uint16(len(entities)); i++ {

		// compare entity at i to all subsequent entities
		// (this way, all entity pairs will be compared once)
		for j := i + 1; j < uint16(len(entities)); j++ {
			if s.TestCollision(i, j) {
				
				s.game_event_manager.Publish(event_generated)
			}
		}
	}
}
