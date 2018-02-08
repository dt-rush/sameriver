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
    "github.com/dt-rush/donkeys-qquest/engine/components"
)


type CollisionLogic struct {

    // NOTE: in the below, i and j are entity ID's

    Selector func (i int,
        j int,
        em *engine.EntityManager) bool

    EventGenerator func (i int,
        j int,
        em *engine.EntityManager) engine.GameEvent
}

type CollisionSystem struct {

    // To filter, lookup entities
    entity_manager *engine.EntityManager
    // Components this system will use
    active_component *components.ActiveComponent
    position_component *components.PositionComponent
    hitbox_component *components.HitboxComponent
    // How the collision system communicates collision events
    game_event_system *engine.GameEventSystem
    // How the collision system gets populated with specific
    // collision detection logics
    collision_logic_collection map[int] CollisionLogic
    collision_logic_ids map[string] int
    collision_logic_active_states map[int] bool
    // to generate IDs for collision logic
    id_system engine.IDSystem
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

        s.collision_logic_collection = make (map[int] CollisionLogic)
        s.collision_logic_ids = make (map[string] int)
        s.collision_logic_active_states = make (map[int] bool)

    }

func (s *CollisionSystem) AddCollisionLogic (name string, logic CollisionLogic) int {

    id := s.id_system.Gen()
    fmt.Printf ("about to add collision logic %s\n", name)
    s.collision_logic_collection [id] = logic
    s.collision_logic_ids [name] = id
    fmt.Printf ("added collision logic %s\n", name)
    return id
}

func (s *CollisionSystem) SetCollisionLogicActiveState (id int, active bool) {
    s.collision_logic_active_states [id] = active
}

func (s *CollisionSystem) TestCollision (i int, j int) bool {

    box := s.hitbox_component.Get (i)
    other_box := s.hitbox_component.Get (j)
    center := s.position_component.Get (i)
    other_center := s.position_component.Get (j)

    // find the distance between the X centers
    dxabs := center[0] - other_center[0]
    if dxabs < 0 {
        dxabs *= -1
    }
    // find the distance between the Y centers
    dyabs := center[1] - other_center[1]
    if dyabs < 0 {
        dyabs *= -1
    }

    // if the sum of the widths is greater than the x distance, collision (same for y)
    collision := dxabs * 2 < (box[0] + other_box[0]) &&
        dyabs * 2 < (box[1] + other_box[1])

    return collision
}

func (s *CollisionSystem) EntityIsCollidable (i int) bool {

    // TODO? factor out the "get all active components
    // with hitbox and position"
    // logic like this with a usage of the tag system or
    // of the entity-to-component one-to-many bitarray system

    return (s.active_component.Has (i) &&
        s.active_component.Get (i)) &&
        (s.position_component.Has (i) &&
        s.hitbox_component.Has (i))

}

func (s *CollisionSystem) Update (dt_ms float64) {

    entities := s.entity_manager.Entities()
    for i := 0; i < len (entities); i++ {
        entity_i := entities [i]
        if ! s.EntityIsCollidable (entity_i) {
            continue
        }
        // compare entity at i to all subsequent entities
        // (this way, all entity pairs will be compared once)
        for j := i + 1; j < len (entities); j++ {
            entity_j := entities [j]
            if ! s.EntityIsCollidable (entity_j) {
                continue
            }

            for collision_logic_id, collision_logic := range s.collision_logic_collection {
                // if this collision logic is active,
                // and the entities i and j match the selector,
                // and there is a collision,
                //    then emit an event according to the eventgenerator
                if s.collision_logic_active_states [collision_logic_id] &&
                    collision_logic.Selector (i, j, s.entity_manager) &&
                    s.TestCollision (i, j) {
                    event_generated := collision_logic.EventGenerator (i, j, s.entity_manager)
                    s.game_event_system.Publish (event_generated)
                }
            }
        }
    }
}

