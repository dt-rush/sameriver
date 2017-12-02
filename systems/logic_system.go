/**
  * 
  * 
  * 
  * 
**/



package systems


import (
	//	"fmt"

	"github.com/dt-rush/donkeys-qquest/engine"
)


type LogicSystem struct {
	
	entity_manager *engine.EntityManager
	
	game_event_system *engine.GameEventSystem

	funcs [](func (float64))
}


func (s *LogicSystem) Init (entity_manager *engine.EntityManager,
	game_event_system *engine.GameEventSystem) {
		
		s.entity_manager = entity_manager
		s.game_event_system = game_event_system 
}


func (s *LogicSystem) RunLogic (f func (float64)) {
	s.funcs = append (s.funcs, f)
}



func (s *LogicSystem) Update (dt_ms float64) {
	for _, f := range (s.funcs) {
		f (dt_ms)
	}
}



