/**
  *
  *
  *
  *
**/



package systems


import (

    // "fmt"

    "github.com/dt-rush/donkeys-qquest/engine"
    "github.com/dt-rush/donkeys-qquest/engine/components"
)


type LogicSystem struct {
    entity_manager *engine.EntityManager
    game_event_system *engine.GameEventSystem
    logic_component *components.LogicComponent
}

func (s *LogicSystem) Init (
    entity_manager *engine.EntityManager,
    game_event_system *engine.GameEventSystem,
    logic_component *components.LogicComponent) {

        s.entity_manager = entity_manager
        s.game_event_system = game_event_system
        s.logic_component = logic_component
}

func (s *LogicSystem) Update (dt_ms float64) {
    for _, id := range s.entity_manager.Entities() {
        if s.entity_manager.EntityHasComponent (id, s.logic_component) {
            logic_unit := s.logic_component.Get (id)
            logic_unit.Logic (dt_ms) 
        }
    }
}
