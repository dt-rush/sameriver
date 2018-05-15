/**
  *
  *
  *
  *
**/

package engine



type LogicSystem struct {
	entity_manager     *EntityManager
	game_event_manager *GameEventManager
	logic_component    *LogicComponent
	active_component   *ActiveComponent
}

func (s *LogicSystem) Init(
	entity_manager *EntityManager,
	game_event_manager *GameEventManager,
	logic_component *LogicComponent,
	active_component *ActiveComponent,
) {

	s.entity_manager = entity_manager
	s.game_event_manager = game_event_manager
	s.logic_component = logic_component
	s.active_component = active_component
}

func (s *LogicSystem) Update(dt_ms int) {
	for _, id := range s.entity_manager.Entities() {
		if s.entity_manager.EntityHasComponent(id, s.logic_component) &&
			s.active_component.Get(id) {
			logic_unit := s.logic_component.Get(id)
			logic_unit.Logic(dt_ms)
		}
	}
}
