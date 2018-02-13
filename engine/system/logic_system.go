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
    "github.com/dt-rush/donkeys-qquest/engine/component"
)


type LogicSystem struct {
    entity_manager *engine.EntityManager
    game_event_system *engine.GameEventSystem
    logic_component *components.LogicComponent
    active_component *components.ActiveComponent
}

func (s *LogicSystem) Init (
    entity_manager *engine.EntityManager,
    game_event_system *engine.GameEventSystem,
    logic_component *components.LogicComponent,
    active_component *components.ActiveComponent,
    ) {

        s.entity_manager = entity_manager
        s.game_event_system = game_event_system
        s.logic_component = logic_component
        s.active_component = active_component
}

func (s *LogicSystem) Update (dt_ms int) {
    for _, id := range s.entity_manager.Entities() {
        if s.entity_manager.EntityHasComponent (id, s.logic_component) &&
            s.active_component.Get (id) {
            logic_unit := s.logic_component.Get (id)
            logic_unit.Logic (dt_ms)
        }
    }
}
