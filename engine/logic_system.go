/**
  *
  *
  *
  *
**/

package engine


// I have chicken wing sauce on my hands


type LogicSystem struct {
	// Reference to entity manager to reach logic component
	entity_manager     *EntityManager
	// targetted entities
	logicEntities *UpdatedEntityList
}

func (s *LogicSystem) Init(entity_manager *EntityManager) {
	s.entity_manager = entity_manager
	// get a regularly updated list of the entities which have logic component
	query := NewBitArraySubsetQuery (
		MakeComponentBitArray([]int{LOGIC_COMPONENT}))
	s.logicEntities = s.entity_manager.GetUpdatedActiveList (query, "logic-bearing")
}

func (s *LogicSystem) Update(dt_ms uint16) {
	s.entity_manager.Components.Logic.Mutex.Lock()
	s.logicEntities.Mutex.Lock()
	defer s.entity_manager.Components.Logic.Mutex.Unlock()
	defer s.logicEntities.Mutex.Unlock()

	for _, id := range s.logicEntities.Entities {
		s.entity_manager.Components.Logic.Data[id].Logic(dt_ms)
	}
}
